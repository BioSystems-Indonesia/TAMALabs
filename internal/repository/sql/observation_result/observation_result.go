package observation_result

import (
	"context"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
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
	return r.DB.Create(data).Error
}

func (r *Repository) CreateMany(ctx context.Context, data []entity.ObservationResult) error {
	return r.DB.Create(data).Error
}

func (r *Repository) FindByID(ctx context.Context, id int64) (result entity.ObservationResult, err error) {
	err = r.DB.First(&result, id).Error
	return
}

func (r *Repository) FindHistory(ctx context.Context, input entity.ObservationResult) (results []entity.ObservationResult, err error) {
	err = r.DB.
		Where("specimen_id = ?", input.SpecimenID).
		Where("code = ?", input.Code).
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

func (r *Repository) DeleteBulk(context context.Context, ids []int64) (entity.ObservationResult, error) {
	var observationResult entity.ObservationResult
	err := r.DB.Delete(&observationResult, "id in ?", ids).Error
	if err != nil {
		return entity.ObservationResult{}, err
	}

	return observationResult, nil
}
