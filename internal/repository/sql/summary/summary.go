package summaryrepo

import (
	"context"
	"log/slog"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"gorm.io/gorm"
)

type SummaryRepository struct {
	db *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) *SummaryRepository {
	return &SummaryRepository{db: db}
}

func (r *SummaryRepository) GetWorkTrendSummary(ctx context.Context) []entity.Summary {
	db := r.db.WithContext(ctx)

	var results []entity.Summary

	// Anchor window to today (don't base on MAX(created_at) which may be far in the past)
	now := time.Now()

	for i := 6; i >= 0; i-- {
		d := now.AddDate(0, 0, -i)
		results = append(results, entity.Summary{
			Name:  d.Format("2006-01-02"),
			Total: 0,
		})
	}

	query := `
		SELECT 
			substr(created_at, 1, 10) AS date, 
			COUNT(*) AS total
		FROM work_orders
		WHERE substr(created_at, 1, 10) >= ?
		GROUP BY substr(created_at, 1, 10)
	`

	startDate := now.AddDate(0, 0, -6).Format("2006-01-02")

	rows, err := db.Raw(query, startDate).Rows()
	if err != nil {
		slog.Debug("failed to query work_orders summary", "error", err)
		return results
	}
	defer rows.Close()

	var date string
	var total int64
	m := make(map[string]int64)

	for rows.Next() {
		if err := rows.Scan(&date, &total); err == nil {
			m[date] = total
		}
	}

	for i := range results {
		if v, ok := m[results[i].Name]; ok {
			results[i].Total = int(v)
		}
	}

	return results
}

func (r *SummaryRepository) GetAbnormalSummary(ctx context.Context) []entity.Summary {
	var totalAbnormal int64
	var totalNormal int64
	db := r.db.WithContext(ctx)
	// calculate start date for 7-day window (including today)
	now := time.Now()
	startDate := now.AddDate(0, 0, -6).Format("2006-01-02")

	// Count abnormal results in the last 7 days
	err := db.Raw(`
		SELECT COUNT(*) 
		FROM observation_results AS o
		LEFT JOIN test_types AS t 
			ON t.code = o.code
		WHERE 
			substr(o.created_at, 1, 10) >= ?
			AND json_valid(o."values") = 1
			AND (
				CAST(json_extract(o."values", '$[0]') AS REAL) < IFNULL(t.low_ref_range, -999999)
				OR CAST(json_extract(o."values", '$[0]') AS REAL) > IFNULL(t.high_ref_range, 999999)
			)
	`, startDate).Scan(&totalAbnormal).Error

	if err != nil {
		return nil
	}

	// Count normal results in the last 7 days
	err = db.Raw(`
		SELECT COUNT(*) 
		FROM observation_results AS o
		LEFT JOIN test_types AS t 
			ON t.code = o.code
		WHERE 
			substr(o.created_at, 1, 10) >= ?
			AND json_valid(o."values") = 1
			AND (
				CAST(json_extract(o."values", '$[0]') AS REAL) >= IFNULL(t.low_ref_range, -999999)
				AND CAST(json_extract(o."values", '$[0]') AS REAL) <= IFNULL(t.high_ref_range, 999999)
			)
	`, startDate).Scan(&totalNormal).Error
	if err != nil {
		return nil
	}
	return []entity.Summary{
		{
			Name:  "Abnormal",
			Total: int(totalAbnormal),
		},
		{
			Name:  "Normal",
			Total: int(totalNormal),
		},
	}
}

func (r *SummaryRepository) GetMostOrderedTest(ctx context.Context) []entity.Summary {
	db := r.db.WithContext(ctx)
	var items []entity.Summary

	// limit to last 7 days (including today)
	now := time.Now()
	startDate := now.AddDate(0, 0, -6).Format("2006-01-02")

	// Query utama: hitung berdasarkan test_code (last 7 days)
	err := db.Table("observation_requests").
		Select("test_code AS name, COUNT(*) AS total").
		Where("substr(created_at, 1, 10) >= ?", startDate).
		Group("test_code").
		Order("total DESC").
		Limit(10).
		Scan(&items).Error

	if err != nil {
		return nil
	}

	// Jika hasil kosong, fallback berdasarkan test_types

	if len(items) == 0 {
		// Fallback: use test_types.name if available (also limited to last 7 days)
		err = db.Table("observation_requests").
			Select(`
				COALESCE(test_types.name, observation_requests.test_code) AS name,
				COUNT(*) AS total`).
			Joins(`LEFT JOIN test_types
					ON test_types.code = observation_requests.test_code`).
			Where("substr(observation_requests.created_at, 1, 10) >= ?", startDate).
			Group("COALESCE(test_types.name, observation_requests.test_code)").
			Order("total DESC").
			Limit(10).
			Scan(&items).Error

		if err != nil {
			return nil
		}
	}

	return items
}

func (r *SummaryRepository) GetTestTypeDistribution(ctx context.Context) []entity.Summary {
	db := r.db.WithContext(ctx)

	var items []entity.Summary

	// limit to last 7 days (including today)
	now := time.Now()
	startDate := now.AddDate(0, 0, -6).Format("2006-01-02")

	err := db.Table("observation_requests").
		Select(`
			COALESCE(test_types.category, 'Uncategorized') AS name,
			COUNT(*) AS total`).
		Joins(`LEFT JOIN test_types
				ON test_types.code = observation_requests.test_code`).
		Where("substr(observation_requests.created_at, 1, 10) >= ?", startDate).
		Group(`COALESCE(test_types.category, 'Uncategorized')`).
		Order("total DESC").
		Scan(&items).Error

	if err != nil {
		return nil
	}

	if len(items) == 0 {
		items = []entity.Summary{{Name: "Uncategorized", Total: 0}}
	}

	return items
}

func (r *SummaryRepository) GetAgeGroupDistribution(ctx context.Context) []entity.Summary {
	// Age group calculation depends on patient.birthdate column; try a DB-side computation if possible
	db := r.db.WithContext(ctx)

	groups := []entity.Summary{
		{Name: "0-9", Total: 0},
		{Name: "10-19", Total: 0},
		{Name: "20-29", Total: 0},
		{Name: "30-39", Total: 0},
		{Name: "40-49", Total: 0},
		{Name: "50+", Total: 0},
	}

	// Attempt to compute age on DB if birthdate exists
	// This query will vary by DB; try a commonly-supported expression using date_part or julianday fallback is hard.
	// We'll instead try to select birthdate and compute age in Go to be safer.
	type p struct {
		Birthdate *time.Time
	}
	var rows []p
	_ = db.Table("patients").Select("birthdate").Find(&rows).Error
	for _, r := range rows {
		if r.Birthdate == nil {
			continue
		}
		age := int(time.Since(*r.Birthdate).Hours() / 24 / 365)
		switch {
		case age < 10:
			groups[0].Total++
		case age < 20:
			groups[1].Total++
		case age < 30:
			groups[2].Total++
		case age < 40:
			groups[3].Total++
		case age < 50:
			groups[4].Total++
		default:
			groups[5].Total++
		}
	}

	return groups
}

func (r *SummaryRepository) GetGenderDistribution(ctx context.Context) []entity.Summary {
	db := r.db.WithContext(ctx)
	type g struct {
		Sex string
		Cnt int64
	}

	var out []g
	_ = db.Table("patients").Select("sex as sex, count(*) as cnt").Group("sex").Scan(&out)

	var items []entity.Summary
	for _, v := range out {
		name := v.Sex
		if name == "" {
			name = "unknown"
		}
		items = append(items, entity.Summary{Name: name, Total: int(v.Cnt)})
	}
	if len(items) == 0 {
		items = []entity.Summary{{Name: "male", Total: 0}, {Name: "female", Total: 0}}
	}
	return items
}

func (r *SummaryRepository) GetSummary(ctx context.Context) entity.SummaryCardResponse {
	db := r.db.WithContext(ctx)
	var resp entity.SummaryCardResponse

	today := time.Now().Format("2006-01-02")

	db.Raw(`
		SELECT COUNT(*) 
		FROM work_orders 
		WHERE substr(created_at, 1, 10) = ?
	`, today).Scan(&resp.TotalWorkOrders)

	db.Raw(`
		SELECT COUNT(*) 
		FROM work_orders 
		WHERE status = 'SUCCESS' 
		  AND substr(created_at, 1, 10) = ?
	`, today).Scan(&resp.CompletedWorkOrders)

	db.Raw(`
		SELECT COUNT(*) 
		FROM work_orders 
		WHERE status = 'PENDING' 
		  AND substr(created_at, 1, 10) = ?
	`, today).Scan(&resp.PendingWorkOrders)

	db.Raw(`
		SELECT COUNT(*) 
		FROM work_orders 
		WHERE status = 'NEW' 
		  AND substr(created_at, 1, 10) = ?
	`, today).Scan(&resp.IncomplateWorkOrders)

	db.Raw(`
		SELECT COUNT(*) 
		FROM observation_results 
		WHERE substr(created_at, 1, 10) = ?
	`, today).Scan(&resp.TotalTest)

	db.Raw(`SELECT COUNT(*) FROM devices`).Scan(&resp.DevicesConnected)

	db.Raw(`SELECT COUNT(*) FROM patients`).Scan(&resp.TotalPatients)

	db.Raw(`SELECT COUNT(*) FROM test_types`).Scan(&resp.TotalTestParameters)

	return resp
}
