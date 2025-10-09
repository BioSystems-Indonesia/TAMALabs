package observation_request

import (
	"context"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

func (r *Repository) Create(ctx context.Context, data *entity.ObservationRequest) error {
	return r.DB.Save(data).Error
}

func (r *Repository) CreateMany(ctx context.Context, data []entity.ObservationRequest) error {
	return r.DB.Save(data).Error
}

func (r *Repository) FindAll(
	ctx context.Context, req *entity.ObservationRequestGetManyRequest,
) (entity.PaginationResponse[entity.ObservationRequest], error) {
	db := r.DB.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})

	if len(req.SpecimenID) > 0 {
		db = db.Where("specimen_id in (?)", req.SpecimenID)
	}

	return sql.GetWithPaginationResponse[entity.ObservationRequest](db, req.GetManyRequest)
}

func (r *Repository) FindOne(ctx context.Context, id int64) (entity.ObservationRequest, error) {
	var ObservationRequest entity.ObservationRequest
	err := r.DB.Where("id = ?", id).First(&ObservationRequest).Error
	if err != nil {
		return entity.ObservationRequest{}, fmt.Errorf("error finding ObservationRequest: %w", err)
	}

	return ObservationRequest, nil
}

func (r *Repository) BulkUpdate(ctx context.Context, request []entity.ObservationRequest) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, req := range request {
			if err := tx.Model(&entity.ObservationRequest{}).Where("id = ?", req.ID).
				Updates(req).Error; err != nil {
				return fmt.Errorf("error updating ObservationRequest: %w", err)
			}
		}

		return nil
	})
}
