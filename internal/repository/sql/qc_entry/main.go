package qcentryrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type QCEntryRepository struct {
	db *gorm.DB
}

func NewQCEntryRepository(db *gorm.DB) *QCEntryRepository {
	return &QCEntryRepository{db: db}
}

func (r *QCEntryRepository) Create(ctx context.Context, entry *entity.QCEntry) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *QCEntryRepository) Update(ctx context.Context, id int, req *entity.UpdateQCEntryRequest) error {
	updates := make(map[string]interface{})

	if req.LotNumber != "" {
		updates["lot_number"] = req.LotNumber
	}
	if req.TargetMean != nil {
		updates["target_mean"] = *req.TargetMean
	}
	if req.TargetSD != nil {
		updates["target_sd"] = *req.TargetSD
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	result := r.db.WithContext(ctx).Model(&entity.QCEntry{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}

func (r *QCEntryRepository) GetByID(ctx context.Context, id int) (*entity.QCEntry, error) {
	var entry entity.QCEntry
	err := r.db.WithContext(ctx).
		Preload("Device").
		Preload("TestType").
		Where("id = ?", id).
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrNotFound
	}

	return &entry, err
}

func (r *QCEntryRepository) GetMany(ctx context.Context, req entity.GetManyRequestQCEntry) ([]entity.QCEntry, int64, error) {
	db := r.db.WithContext(ctx)

	if req.DeviceID != nil {
		db = db.Where("device_id = ?", *req.DeviceID)
	}
	if req.TestTypeID != nil {
		db = db.Where("test_type_id = ?", *req.TestTypeID)
	}
	if req.QCLevel != nil {
		db = db.Where("qc_level = ?", *req.QCLevel)
	}
	if req.IsActive != nil {
		db = db.Where("is_active = ?", *req.IsActive)
	}

	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})
	db = db.Preload("Device").Preload("TestType")

	var entries []entity.QCEntry
	var total int64

	if err := db.Model(&entity.QCEntry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").Find(&entries).Error; err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *QCEntryRepository) GetActiveEntry(ctx context.Context, deviceID, testTypeID, qcLevel int) (*entity.QCEntry, error) {
	var entry entity.QCEntry
	err := r.db.WithContext(ctx).
		Where("device_id = ? AND test_type_id = ? AND qc_level = ? AND is_active = ?",
			deviceID, testTypeID, qcLevel, true).
		Order("created_at DESC").
		First(&entry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrNotFound
	}

	return &entry, err
}

func (r *QCEntryRepository) DeactivateOldEntries(ctx context.Context, deviceID, testTypeID, qcLevel int) error {
	return r.db.WithContext(ctx).
		Model(&entity.QCEntry{}).
		Where("device_id = ? AND test_type_id = ? AND qc_level = ? AND is_active = ?",
			deviceID, testTypeID, qcLevel, true).
		Update("is_active", false).Error
}

func (r *QCEntryRepository) GetDeviceSummary(ctx context.Context, deviceID int) (*entity.QCSummary, error) {
	summary := &entity.QCSummary{
		DeviceID: deviceID,
	}

	// Get total QC count
	var totalCount int64
	if err := r.db.WithContext(ctx).
		Model(&entity.QCResult{}).
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		Where("qc_entries.device_id = ?", deviceID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}
	summary.TotalQC = int(totalCount) / 2

	// Get this month's QC count
	var monthCount int64
	if err := r.db.WithContext(ctx).
		Model(&entity.QCResult{}).
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		// Use substr on created_at to be robust to ISO timestamps (YYYY-MM-DD...)
		Where("qc_entries.device_id = ? AND substr(qc_results.created_at,1,7) = strftime('%Y-%m','now')", deviceID).
		Count(&monthCount).Error; err != nil {
		return nil, err
	}
	summary.QCThisMonth = int(monthCount) / 2

	// Get last QC result
	var lastResult entity.QCResult
	err := r.db.WithContext(ctx).
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		Where("qc_entries.device_id = ?", deviceID).
		Order("qc_results.created_at DESC").
		First(&lastResult).Error

	if err == nil {
		// Set pointer to last result time
		t := lastResult.CreatedAt
		summary.LastQCDate = &t
		summary.LastQCStatus = lastResult.Result
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Get today's QC status - check which levels have results today
	var todayResults []struct {
		QCLevel int
	}
	if err := r.db.WithContext(ctx).
		Model(&entity.QCResult{}).
		Select("DISTINCT qc_entries.qc_level").
		Joins("JOIN qc_entries ON qc_entries.id = qc_results.qc_entry_id").
		// Use substr to compare date portion (YYYY-MM-DD) to date('now')
		Where("qc_entries.device_id = ? AND substr(qc_results.created_at,1,10) = date('now')", deviceID).
		Scan(&todayResults).Error; err != nil {
		return nil, err
	}

	// default false
	summary.Level1Today = false
	summary.Level2Today = false
	summary.Level3Today = false

	for _, tr := range todayResults {
		switch tr.QCLevel {
		case 1:
			summary.Level1Today = true
		case 2:
			summary.Level2Today = true
		case 3:
			summary.Level3Today = true
		}
	}

	if !summary.Level1Today && !summary.Level2Today && !summary.Level3Today {
		summary.QCTodayStatus = "Not Done"
	} else if summary.Level1Today && summary.Level2Today && summary.Level3Today {
		summary.QCTodayStatus = "Done"
	} else {
		summary.QCTodayStatus = "Partial"
	}

	// Get level completion status - check if device has active entries for each level
	var activeEntries []entity.QCEntry
	if err := r.db.WithContext(ctx).
		Where("device_id = ? AND is_active = ?", deviceID, true).
		Find(&activeEntries).Error; err != nil {
		return nil, err
	}

	for _, entry := range activeEntries {
		switch entry.QCLevel {
		case 1:
			summary.Level1Complete = true
		case 2:
			summary.Level2Complete = true
		case 3:
			summary.Level3Complete = true
		}
	}

	return summary, nil
}
