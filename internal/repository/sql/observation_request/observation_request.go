package observation_request

import (
	"context"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{db: db, cfg: cfg}
}

func (r *Repository) Create(ctx context.Context, data *entity.ObservationRequest) error {
	return r.db.Save(data).Error
}

func (r *Repository) CreateMany(ctx context.Context, data []entity.ObservationRequest) error {
	return r.db.Save(data).Error
}

func (r *Repository) FindAll(ctx context.Context, req *entity.ObservationRequestGetManyRequest) ([]entity.ObservationRequest, error) {
	var ObservationRequests []entity.ObservationRequest

	db := r.db.WithContext(ctx)
	if len(req.ID) > 0 {
		db = db.Where("id in (?)", req.ID)
	}

	if len(req.SpecimenID) > 0 {
		db = db.Where("specimen_id in (?)", req.SpecimenID)
	}

	if req.Sort != "" {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: req.Sort,
			},
			Desc: req.IsSortDesc(),
		})
	}

	err := db.Find(&ObservationRequests).Error
	if err != nil {
		return nil, fmt.Errorf("error finding ObservationRequests: %w", err)
	}
	return ObservationRequests, nil
}

func (r *Repository) FindOne(ctx context.Context, id int64) (entity.ObservationRequest, error) {
	var ObservationRequest entity.ObservationRequest
	err := r.db.Where("id = ?", id).First(&ObservationRequest).Error
	if err != nil {
		return entity.ObservationRequest{}, fmt.Errorf("error finding ObservationRequest: %w", err)
	}

	return ObservationRequest, nil
}
