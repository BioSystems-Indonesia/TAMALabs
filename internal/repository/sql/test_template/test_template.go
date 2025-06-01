package test_template

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
	"gorm.io/gorm"
)

type Repository struct {
	DB  *gorm.DB
	cfg *config.Schema
}

// NewRepository creates a new test type repository.
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

// FindAll returns all test types.
func (r *Repository) FindAll(
	_ context.Context,
	req *entity.TestTemplateGetManyRequest,
) (entity.PaginationResponse[entity.TestTemplate], error) {
	db := r.DB.Preload("CreatedByUser").Preload("LastUpdatedByUser")
	sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
			return db.Where("name like ?", "%"+query+"%").
				Or("description like ?", "%"+query+"%")
		}})

	resp, err := sql.GetWithPaginationResponse[entity.TestTemplate](db, req.GetManyRequest)
	if err != nil {
		return entity.PaginationResponse[entity.TestTemplate]{}, err
	}

	return resp, nil
}

func (r *Repository) FindOneByID(ctx context.Context, id int) (entity.TestTemplate, error) {
	db := r.DB.Preload("CreatedByUser").Preload("LastUpdatedByUser")
	var data entity.TestTemplate
	if err := db.First(&data, id).Error; err != nil {
		return entity.TestTemplate{}, err
	}

	return data, nil
}

func (r *Repository) Create(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(req).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.TestTemplate{}, err
	}

	return *req, nil
}

func (r *Repository) Update(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	err := r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(req).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return entity.TestTemplate{}, err
	}

	return *req, nil
}

func (r *Repository) Delete(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	if err := r.DB.Delete(req).Error; err != nil {
		return entity.TestTemplate{}, err
	}
	return *req, nil
}
