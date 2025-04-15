package device

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{DB: db}
}

func (r *Repository) FindByID(ctx context.Context, id int64) (entity.Device, error) {
	var device entity.Device
	tx := r.DB.First(&device, id)
	if tx.Error != nil {
		return entity.Device{}, tx.Error
	}
	return device, nil
}
