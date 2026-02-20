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
		Where(`
			EXISTS (
				SELECT 1 FROM observation_results o
				WHERE o.specimen_id = specimens.id
				AND (
					o.last_sync IS NULL
					OR o.updated_at IS NULL
					OR o.updated_at > o.last_sync
				)
			)
		`).
		Preload(
			"ObservationResult",
			`
			(
				(
					last_sync IS NULL
					OR updated_at IS NULL
					OR updated_at > last_sync
				)
				AND (
					picked = ?
					OR NOT EXISTS (
						SELECT 1 FROM observation_results o2
						WHERE o2.specimen_id = observation_results.specimen_id
						AND (
							o2.last_sync IS NULL
							OR o2.updated_at IS NULL
							OR o2.updated_at > o2.last_sync
						)
						AND o2.picked = 1
					)
				)
			)
		`, true).
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
