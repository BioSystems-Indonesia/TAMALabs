package observation_result

import (
	"context"
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
	return r.DB.Save(data).Error
}

func (r *Repository) CreateMany(ctx context.Context, data []entity.ObservationResult) error {
	return r.DB.Create(data).Error
}
