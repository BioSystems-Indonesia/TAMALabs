package qualitycontrolrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type QualityControlRepository struct {
	db *gorm.DB
}

func NewQualityControlRepository(db *gorm.DB) *QualityControlRepository {
	return &QualityControlRepository{db: db}
}

func (r *QualityControlRepository) Create(ctx context.Context, qc *entity.QualityControl) error {
	return r.db.WithContext(ctx).Create(qc).Error
}

func (r *QualityControlRepository) GetMany(
	ctx context.Context,
	req entity.GetManyRequestQualityControl,
) ([]entity.QualityControl, int64, error) {
	db := r.db.WithContext(ctx)

	// Apply filters
	if req.DeviceID != nil {
		db = db.Where("device_id = ?", *req.DeviceID)
	}
	if req.TestTypeID != nil {
		db = db.Where("test_type_id = ?", *req.TestTypeID)
	}
	if req.QCLevel != nil {
		db = db.Where("qc_level = ?", *req.QCLevel)
	}

	// Apply pagination and sorting
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})

	// Preload relations
	db = db.Preload("Device").Preload("TestType")

	// Get records
	var qcs []entity.QualityControl
	var total int64

	// Count total
	if err := db.Model(&entity.QualityControl{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting QC records: %w", err)
	}

	// Get paginated results
	if err := db.Order("created_at DESC").Find(&qcs).Error; err != nil {
		return nil, 0, fmt.Errorf("error getting QC records: %w", err)
	}

	return qcs, total, nil
}

func (r *QualityControlRepository) GetByID(ctx context.Context, id int) (*entity.QualityControl, error) {
	var qc entity.QualityControl
	err := r.db.WithContext(ctx).
		Preload("Device").
		Preload("TestType").
		Where("id = ?", id).
		First(&qc).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error finding QC record: %w", err)
	}

	return &qc, nil
}

func (r *QualityControlRepository) GetStatistics(ctx context.Context, deviceID int) (map[string]interface{}, error) {
	var totalQC int64
	var qcThisMonth int64
	var lastQC time.Time
	var passCount int64
	var failCount int64

	db := r.db.WithContext(ctx).Model(&entity.QualityControl{}).Where("device_id = ?", deviceID)

	// Total QC count
	if err := db.Count(&totalQC).Error; err != nil {
		return nil, fmt.Errorf("error counting total QC: %w", err)
	}

	// QC this month
	startOfMonth := time.Now().Truncate(24 * time.Hour)
	startOfMonth = time.Date(startOfMonth.Year(), startOfMonth.Month(), 1, 0, 0, 0, 0, startOfMonth.Location())
	if err := db.Where("created_at >= ?", startOfMonth).Count(&qcThisMonth).Error; err != nil {
		return nil, fmt.Errorf("error counting QC this month: %w", err)
	}

	// Last QC date
	var qc entity.QualityControl
	if err := db.Order("created_at DESC").First(&qc).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error getting last QC: %w", err)
	}
	if !errors.Is(db.Error, gorm.ErrRecordNotFound) {
		lastQC = qc.CreatedAt
	}

	// Pass/Fail counts
	if err := db.Where("result = ?", "Pass").Count(&passCount).Error; err != nil {
		return nil, fmt.Errorf("error counting pass: %w", err)
	}
	if err := db.Where("result = ?", "Fail").Count(&failCount).Error; err != nil {
		return nil, fmt.Errorf("error counting fail: %w", err)
	}

	status := "Pending"
	if totalQC > 0 {
		if passCount > failCount {
			status = "Good"
		} else {
			status = "Warning"
		}
	}

	return map[string]interface{}{
		"total_qc":      totalQC,
		"qc_this_month": qcThisMonth,
		"last_qc":       lastQC,
		"pass_count":    passCount,
		"fail_count":    failCount,
		"status":        status,
	}, nil
}

func (r *QualityControlRepository) GetCountByLevel(ctx context.Context, deviceID int, testTypeID *int) (map[string]interface{}, error) {
	db := r.db.WithContext(ctx).Model(&entity.QualityControl{}).Where("device_id = ?", deviceID)

	// Optional filter by test type
	if testTypeID != nil {
		db = db.Where("test_type_id = ?", *testTypeID)
	}

	var level1Count, level2Count, level3Count int64

	// Count QC Level 1
	if err := db.Where("qc_level = ?", 1).Count(&level1Count).Error; err != nil {
		return nil, fmt.Errorf("error counting level 1: %w", err)
	}

	// Count QC Level 2
	if err := db.Where("qc_level = ?", 2).Count(&level2Count).Error; err != nil {
		return nil, fmt.Errorf("error counting level 2: %w", err)
	}

	// Count QC Level 3
	if err := db.Where("qc_level = ?", 3).Count(&level3Count).Error; err != nil {
		return nil, fmt.Errorf("error counting level 3: %w", err)
	}

	return map[string]interface{}{
		"level_1_count": level1Count,
		"level_2_count": level2Count,
		"level_3_count": level3Count,
		"total_count":   level1Count + level2Count + level3Count,
	}, nil
}
