package observation_request

import (
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{db: db, cfg: cfg}
}

func (r *Repository) Create(data *entity.ObservationRequest) error {
	return r.db.Save(data).Error
}

func (r *Repository) CreateMany(data *[]entity.ObservationRequest) error {
	return r.db.Save(data).Error
}
