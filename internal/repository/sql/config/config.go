package configrepo

import (
	"context"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

func (r *Repository) FindAll(ctx context.Context, req *entity.ConfigGetManyRequest) (entity.ConfigPaginationResponse, error) {
	var configs []entity.Config

	db := r.DB.WithContext(ctx)
	if len(req.ID) > 0 {
		db = db.Where("key in (?)", req.ID)
	}

	if req.Sort != "" {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: req.Sort,
			},
			Desc: req.IsSortDesc(),
		})
	}

	err := db.Find(&configs).Error
	if err != nil {
		return entity.ConfigPaginationResponse{}, fmt.Errorf("error finding Configs: %w", err)
	}

	var total int64
	err = db.Model(&entity.Config{}).Count(&total).Error
	if err != nil {
		return entity.ConfigPaginationResponse{}, fmt.Errorf("error counting Configs: %w", err)
	}

	return entity.ConfigPaginationResponse{
		Data: configs,
		PaginationResponse: entity.PaginationResponse{
			Total: total,
		},
	}, nil
}

func (r *Repository) FindOne(ctx context.Context, id string) (entity.Config, error) {
	var config entity.Config
	err := r.DB.WithContext(ctx).First(&config, "id = ?", id).Error
	if err != nil {
		return entity.Config{}, fmt.Errorf("error finding Config: %w", err)
	}
	return config, nil
}

func (r *Repository) Edit(ctx context.Context, id string, value string) (entity.Config, error) {
	config, err := r.FindOne(ctx, id)
	if err != nil {
		return entity.Config{}, fmt.Errorf("error finding Config: %w", err)
	}

	config.Value = value
	err = r.DB.WithContext(ctx).Save(&config).Error
	if err != nil {
		return entity.Config{}, fmt.Errorf("error updating Config: %w", err)
	}
	return config, nil
}
