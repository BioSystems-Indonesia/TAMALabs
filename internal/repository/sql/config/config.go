package configrepo

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

func (r *Repository) FindAll(
	ctx context.Context, req *entity.ConfigGetManyRequest,
) (entity.PaginationResponse[entity.Config], error) {
	db := r.DB.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				return db.Where("key like ?", "%"+query+"%")
			},
		})

	return sql.GetWithPaginationResponse[entity.Config](db, req.GetManyRequest)
}

func (r *Repository) FindOne(ctx context.Context, id string) (entity.Config, error) {
	var config entity.Config
	err := r.DB.WithContext(ctx).First(&config, "id = ?", id).Error
	if err != nil {
		return entity.Config{}, fmt.Errorf("error finding Config: %w", err)
	}
	return config, nil
}

// Get retrieves a config value by key (alias for FindOne with value extraction)
func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	config, err := r.FindOne(ctx, key)
	if err != nil {
		return "", err
	}
	return config.Value, nil
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
