package observation_result

import (
	"context"
	"fmt"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

func (r *Repository) Create(ctx context.Context, data *entity.ObservationResult) error {
	if data.CreatedBy == 0 {
		data.CreatedBy = int64(constant.CreatedBySystem)
	}

	return r.DB.Create(data).Error
}

func (r *Repository) CreateMany(ctx context.Context, data []entity.ObservationResult) error {
	for i := range data {
		if data[i].CreatedBy == 0 {
			data[i].CreatedBy = int64(constant.CreatedBySystem)
		}
	}

	return r.DB.Create(data).Error
}

func (r *Repository) Exists(ctx context.Context, specimenID int64, code string, date time.Time, firstValue string) (bool, error) {
	var count int64

	err := r.DB.WithContext(ctx).Raw(
		"SELECT COUNT(*) FROM observation_results WHERE specimen_id = ? AND code = ? AND date = ? AND (json_extract(values, '$[0]') = ? OR values LIKE ?)",
		specimenID, code, date, firstValue, "%\""+firstValue+"\"%",
	).Scan(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (result entity.ObservationResult, err error) {
	err = r.DB.First(&result, id).Error
	return
}

func (r *Repository) FindHistory(ctx context.Context, input entity.ObservationResult) (results []entity.ObservationResult, err error) {
	err = r.DB.Preload("CreatedByAdmin").Preload("TestType").
		Where("specimen_id = ?", input.SpecimenID).
		Where("code = ?", input.TestCode).
		Order("created_at DESC").
		Find(&results).Error
	return
}

func (r *Repository) Delete(context context.Context, id int64) (entity.ObservationResult, error) {
	var observationResult entity.ObservationResult
	err := r.DB.Where("id = ?", id).First(&observationResult).Error
	if err != nil {
		return entity.ObservationResult{}, err
	}

	err = r.DB.Delete(&observationResult).Error
	if err != nil {
		return entity.ObservationResult{}, err
	}

	return observationResult, nil
}

func (r *Repository) PickObservationResult(ctx context.Context, id int64) (entity.ObservationResult, error) {
	var observationResult entity.ObservationResult

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id = ?", id).First(&observationResult).Error
		if err != nil {
			return fmt.Errorf("get observation result error: %w", err)
		}

		// Update all other observation result picked to false
		err = tx.Model(&entity.ObservationResult{}).
			Where("specimen_id = ? AND code = ? AND id != ?", observationResult.SpecimenID, observationResult.TestCode, id).
			Update("picked", false).
			Error
		if err != nil {
			return fmt.Errorf("failed to update observation result: %w", err)
		}

		// Update the picked result to true
		err = tx.Model(&entity.ObservationResult{}).
			Where("id = ?", id).Update("picked", true).Error
		if err != nil {
			return fmt.Errorf("failed to update observation result pick: %w", err)
		}

		return nil
	})

	observationResult.Picked = true
	return observationResult, err
}

func (r *Repository) ApproveResult(context context.Context, workOrderID int64) error {
	db := r.DB.Model(&entity.WorkOrder{})
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id =?", workOrderID).
			Update("verified_status", entity.WorkOrderVerifiedStatusVerified).
			Error
		if err != nil {
			return fmt.Errorf("failed to update verified status: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RejectResult(context context.Context, workOrderID int64) error {
	db := r.DB.Model(&entity.WorkOrder{})
	err := db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("id =?", workOrderID).
			Update("verified_status", entity.WorkOrderVerifiedStatusRejected).
			Error
		if err != nil {
			return fmt.Errorf("failed to update verified status: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
