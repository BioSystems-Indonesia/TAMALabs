package test_type

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

// NewRepository creates a new test type repository
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

// FindAll returns all test types
func (r *Repository) FindAll(ctx context.Context, req *entity.TestTypeGetManyRequest) ([]entity.TestType, error) {
	var data []entity.TestType
	if err := r.DB.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
