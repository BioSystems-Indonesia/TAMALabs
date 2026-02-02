package sub_category

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{db: db, cfg: cfg}
}

// FindAll returns all sub-categories with pagination
func (r *Repository) FindAll(
	ctx context.Context,
	req *entity.SubCategoryGetManyRequest,
) (entity.PaginationResponse[entity.SubCategory], error) {
	db := r.db.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
			return db.Where("name LIKE ? OR code LIKE ?", "%"+query+"%", "%"+query+"%")
		},
	})

	if req.Category != "" {
		db = db.Where("category = ?", req.Category)
	}

	return sql.GetWithPaginationResponse[entity.SubCategory](db, req.GetManyRequest)
}

// FindByID finds a sub-category by ID
func (r *Repository) FindByID(ctx context.Context, id int) (*entity.SubCategory, error) {
	var subCategory entity.SubCategory
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&subCategory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding sub-category: %w", err)
	}
	return &subCategory, nil
}

// FindByName finds a sub-category by name
func (r *Repository) FindByName(ctx context.Context, name string) (*entity.SubCategory, error) {
	var subCategory entity.SubCategory
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&subCategory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error finding sub-category: %w", err)
	}
	return &subCategory, nil
}

// Create creates a new sub-category
func (r *Repository) Create(ctx context.Context, subCategory *entity.SubCategory) error {
	return r.db.WithContext(ctx).Create(subCategory).Error
}

// Update updates a sub-category
func (r *Repository) Update(ctx context.Context, subCategory *entity.SubCategory) error {
	return r.db.WithContext(ctx).Save(subCategory).Error
}

// Delete deletes a sub-category
func (r *Repository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&entity.SubCategory{}, id).Error
}

// GetTestTypesBySubCategoryID returns all test types for a sub-category
func (r *Repository) GetTestTypesBySubCategoryID(ctx context.Context, subCategoryID int) ([]entity.TestType, error) {
	// First, get the sub-category to get its name for fallback
	subCategory, err := r.FindByID(ctx, subCategoryID)
	if err != nil {
		return nil, fmt.Errorf("error finding sub-category: %w", err)
	}

	var testTypes []entity.TestType

	// Query using sub_category_id OR sub_category string (for backward compatibility)
	err = r.db.WithContext(ctx).
		Preload("Devices").
		Preload("SubCategoryDetail").
		Where("sub_category_id = ? OR sub_category = ?", subCategoryID, subCategory.Name).
		Find(&testTypes).Error

	if err != nil {
		return nil, fmt.Errorf("error finding test types by sub-category: %w", err)
	}

	return testTypes, nil
}
