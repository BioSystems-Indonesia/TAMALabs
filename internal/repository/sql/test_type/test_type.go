package test_type

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

// NewRepository creates a new test type repository
func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{DB: db, cfg: cfg}
}

// FindAll returns all test types.
func (r *Repository) FindAll(
	ctx context.Context, req *entity.TestTypeGetManyRequest,
) (entity.PaginationResponse[entity.TestType], error) {
	db := r.DB
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
			return db.Where("name like ?", "%"+query+"%").
				Or("code like ?", "%"+query+"%").
				Or("description like ?", "%"+query+"%")
		},
	})

	if req.Code != "" {
		db = db.Where("code like ?", "%"+req.Code+"%")
	}

	if len(req.Categories) != 0 {
		db = db.Where("category in (?)", req.Categories)
	}

	if len(req.SubCategories) != 0 {
		db = db.Where("sub_category in (?)", req.SubCategories)
	}

	return sql.GetWithPaginationResponse[entity.TestType](db, req.GetManyRequest)
}

func (r *Repository) FindOneByID(ctx context.Context, id int) (entity.TestType, error) {
	var data entity.TestType
	if err := r.DB.First(&data, id).Error; err != nil {
		return entity.TestType{}, err
	}
	return data, nil
}

func (r *Repository) Create(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	if err := r.DB.Create(req).Error; err != nil {
		return entity.TestType{}, err
	}
	return *req, nil
}

func (r *Repository) Update(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	if err := r.DB.Save(req).Error; err != nil {
		return entity.TestType{}, err
	}
	return *req, nil
}

func (r *Repository) Delete(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	if err := r.DB.Delete(req).Error; err != nil {
		return entity.TestType{}, err
	}
	return *req, nil
}
