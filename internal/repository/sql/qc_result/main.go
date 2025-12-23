package qcresultrepo

import (
	"context"
	"math"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type QCResultRepository struct {
	db *gorm.DB
}

func NewQCResultRepository(db *gorm.DB) *QCResultRepository {
	return &QCResultRepository{db: db}
}

func (r *QCResultRepository) Create(ctx context.Context, result *entity.QCResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

func (r *QCResultRepository) GetMany(ctx context.Context, req entity.GetManyRequestQCResult) ([]entity.QCResult, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.QCResult{})

	if req.QCEntryID != nil {
		db = db.Where("qc_results.qc_entry_id = ?", *req.QCEntryID)
	}

	if req.Method != nil {
		db = db.Where("qc_results.method = ?", *req.Method)
	}

	// If filtering by device or test type, join qc_entries once and apply where clauses
	if req.DeviceID != nil || req.TestTypeID != nil {
		db = db.Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id")
		if req.DeviceID != nil {
			db = db.Where("qc_entries.device_id = ?", *req.DeviceID)
		}
		if req.TestTypeID != nil {
			db = db.Where("qc_entries.test_type_id = ?", *req.TestTypeID)
		}
	}

	// Count before pagination
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})

	// Fetch results with preload
	var results []entity.QCResult
	if err := db.Preload("QCEntry.Device").Preload("QCEntry.TestType").
		Order("qc_results.created_at DESC").
		Find(&results).Error; err != nil {
		return nil, 0, err
	}

	// Calculate result count for each entry
	entryCountMap := make(map[int]int)
	for i := range results {
		if results[i].QCEntryID > 0 {
			if _, exists := entryCountMap[results[i].QCEntryID]; !exists {
				var count int64
				r.db.Model(&entity.QCResult{}).Where("qc_entry_id = ?", results[i].QCEntryID).Count(&count)
				entryCountMap[results[i].QCEntryID] = int(count)
			}
			results[i].ResultCount = entryCountMap[results[i].QCEntryID]
		}
	}

	return results, total, nil
}

func (r *QCResultRepository) GetByEntryID(ctx context.Context, entryID int) ([]entity.QCResult, error) {
	var results []entity.QCResult
	err := r.db.WithContext(ctx).
		Where("qc_results.qc_entry_id = ?", entryID).
		Order("qc_results.created_at ASC").
		Find(&results).Error

	return results, err
}

func (r *QCResultRepository) GetByEntryIDAndMethod(ctx context.Context, entryID int, method string) ([]entity.QCResult, error) {
	var results []entity.QCResult
	err := r.db.WithContext(ctx).
		Where("qc_results.qc_entry_id = ? AND qc_results.method = ?", entryID, method).
		Order("qc_results.created_at ASC").
		Find(&results).Error

	return results, err
}

func (r *QCResultRepository) CalculateStatistics(ctx context.Context, entryID int) (mean float64, sd float64, count int, err error) {
	var results []entity.QCResult
	err = r.db.WithContext(ctx).
		Where("qc_entry_id = ?", entryID).
		Find(&results).Error

	if err != nil {
		return 0, 0, 0, err
	}

	count = len(results)
	if count == 0 {
		return 0, 0, 0, nil
	}

	// Calculate mean
	var sum float64
	for _, r := range results {
		sum += r.MeasuredValue
	}
	mean = sum / float64(count)

	// Calculate SD using formula: sqrt(Σ(x - mean)² / (n - 1))
	if count > 1 {
		var variance float64
		for _, r := range results {
			diff := r.MeasuredValue - mean
			variance += diff * diff
		}
		sd = math.Sqrt(variance / float64(count-1))
	} else {
		sd = 0
	}

	return mean, sd, count, nil
}

func (r *QCResultRepository) GetCountByLevel(ctx context.Context, deviceID int, testTypeID *int) (map[string]interface{}, error) {
	db := r.db.WithContext(ctx).Model(&entity.QCResult{}).
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		Where("qc_entries.device_id = ?", deviceID)

	// Optional filter by test type
	if testTypeID != nil {
		db = db.Where("qc_entries.test_type_id = ?", *testTypeID)
	}

	var level1Count, level2Count, level3Count int64

	// Count QC Level 1
	if err := db.Where("qc_entries.qc_level = ?", 1).Count(&level1Count).Error; err != nil {
		return nil, err
	}

	// Reset the query for each count
	db = r.db.WithContext(ctx).Model(&entity.QCResult{}).
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		Where("qc_entries.device_id = ?", deviceID)
	if testTypeID != nil {
		db = db.Where("qc_entries.test_type_id = ?", *testTypeID)
	}

	// Count QC Level 2
	if err := db.Where("qc_entries.qc_level = ?", 2).Count(&level2Count).Error; err != nil {
		return nil, err
	}

	// Reset the query for level 3
	db = r.db.WithContext(ctx).Model(&entity.QCResult{}).
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		Where("qc_entries.device_id = ?", deviceID)
	if testTypeID != nil {
		db = db.Where("qc_entries.test_type_id = ?", *testTypeID)
	}

	// Count QC Level 3
	if err := db.Where("qc_entries.qc_level = ?", 3).Count(&level3Count).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"level_1_count": level1Count,
		"level_2_count": level2Count,
		"level_3_count": level3Count,
		"total_count":   level1Count + level2Count + level3Count,
	}, nil
}
