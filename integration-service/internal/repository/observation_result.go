package repositories

import (
	"context"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
	"gorm.io/gorm"
)

type ObservationResultRepository interface {
	FindAll(ctx context.Context) ([]models.Specimen, error)
	Verify(ctx context.Context, id string, status string) error
	UpdateLastSync(ctx context.Context, observationResultIds []int, syncTime time.Time) error
}

type ObservationResultImpl struct {
	db *gorm.DB
}

func NewObservationResultRepository(db *gorm.DB) ObservationResultRepository {
	return &ObservationResultImpl{
		db: db,
	}
}

func (r ObservationResultImpl) FindAll(ctx context.Context) ([]models.Specimen, error) {
	var specimens []models.Specimen

	err := r.db.WithContext(ctx).
		Table("specimens").
		Joins(`
			JOIN observation_results 
			ON observation_results.specimen_id = specimens.id
		`).
		Where("observation_results.picked = ?", true).
		Where(`
			(
				observation_results.last_sync IS NULL
				OR observation_results.updated_at IS NULL
				OR observation_results.updated_at > observation_results.last_sync
			)
		`).
		Preload(
			"ObservationResult",
			`
				picked = ?
						AND (
						last_sync IS NULL
				OR updated_at IS NULL
				OR updated_at > last_sync
				)
			`,
			true,
		).
		Preload("ObservationRequest").
		Preload("ObservationRequest.TestType").
		Preload("WorkOrder").
		Preload("Patient").
		Find(&specimens).
		Error

	if err != nil {
		return nil, err
	}

	return specimens, nil
}

func (r ObservationResultImpl) UpdateLastSync(ctx context.Context, observationResultIds []int, syncTime time.Time) error {
	if len(observationResultIds) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Table("observation_results").
		Where("id IN ?", observationResultIds).
		UpdateColumn("last_sync", syncTime).
		Error
}

func (r ObservationResultImpl) Verify(ctx context.Context, id string, status string) error {

	err := r.db.WithContext(ctx).Table("work_orders").Where("barcode = ?", id).
		UpdateColumn("verified_status", status).Error
	return err
}
