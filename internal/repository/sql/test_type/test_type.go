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
	query := r.DB
	if len(req.ID) != 0 {
		query = query.Where("id in (?)", req.ID)
	}

	if req.Query != "" {
		query = query.Where("name like ?", "%"+req.Query+"%").
			Or("code like ?", "%"+req.Query+"%").
			Or("description like ?", "%"+req.Query+"%")
	}

	if req.Search != "" {
		query = query.Where("name like ?", "%"+req.Search+"%").
			Or("code like ?", "%"+req.Search+"%").
			Or("description like ?", "%"+req.Search+"%")
	}

	if req.Code != "" {
		query = query.Where("code like ?", "%"+req.Code+"%")
	}

	if len(req.Categories) != 0 {
		query = query.Where("category in (?)", req.Categories)
	}

	if len(req.SubCategories) != 0 {
		query = query.Where("sub_category in (?)", req.SubCategories)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
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
