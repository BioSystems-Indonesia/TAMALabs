package test_type

import (
	"context"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	subcategoryrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/sub_category"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
)

type Usecase struct {
	repository      *test_type.Repository
	subCategoryRepo *subcategoryrepo.Repository
}

func NewUsecase(testTypeRepository *test_type.Repository, subCategoryRepo *subcategoryrepo.Repository) *Usecase {
	return &Usecase{
		repository:      testTypeRepository,
		subCategoryRepo: subCategoryRepo,
	}
}

func (u *Usecase) FindAll(ctx context.Context, req *entity.TestTypeGetManyRequest) (entity.PaginationResponse[entity.TestType], error) {
	return u.repository.FindAll(ctx, req)
}

func (u *Usecase) ListAllFilter(ctx context.Context) (entity.TestTypeFilter, error) {
	return u.repository.FindAllFilter(ctx)
}

func (u *Usecase) FindOneByID(ctx context.Context, id int) (entity.TestType, error) {
	return u.repository.FindOneByID(ctx, id)
}

func (u *Usecase) FindOneByCode(ctx context.Context, code string) (entity.TestType, error) {
	return u.repository.FindOneByCode(ctx, code)
}

func (u *Usecase) FindOneByAliasCode(ctx context.Context, aliasCode string) (entity.TestType, error) {
	return u.repository.FindOneByAliasCode(ctx, aliasCode)
}

func (u *Usecase) Create(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	// Handle sub_category: ensure it exists in sub_categories table
	if req.SubCategory != "" {
		subCategoryID, err := u.ensureSubCategoryExists(ctx, req.SubCategory, req.Category)
		if err != nil {
			// Log error but continue - fallback to just using string
			// This maintains backward compatibility
		} else {
			req.SubCategoryID = &subCategoryID
		}
	}

	req.UpdatedAt = time.Now()

	return u.repository.Create(ctx, req)
}

func (u *Usecase) Update(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	// Handle sub_category: ensure it exists in sub_categories table
	if req.SubCategory != "" {
		subCategoryID, err := u.ensureSubCategoryExists(ctx, req.SubCategory, req.Category)
		if err != nil {
			// Log error but continue - fallback to just using string
		} else {
			req.SubCategoryID = &subCategoryID
		}
	}

	return u.repository.Update(ctx, req)
}

// ensureSubCategoryExists checks if sub-category exists, creates if not, and returns the ID
func (u *Usecase) ensureSubCategoryExists(ctx context.Context, name string, category string) (int, error) {
	// Try to find existing sub-category
	subCategory, err := u.subCategoryRepo.FindByName(ctx, name)
	if err == nil {
		// Found existing
		return subCategory.ID, nil
	}

	// Not found, create new sub-category
	// Generate code from name (first 3 letters uppercase)
	code := strings.ToUpper(name)
	if len(code) > 3 {
		code = code[:3]
	}

	newSubCategory := &entity.SubCategory{
		Name:        name,
		Code:        code,
		Category:    category,
		Description: "",
	}

	err = u.subCategoryRepo.Create(ctx, newSubCategory)
	if err != nil {
		return 0, err
	}

	return newSubCategory.ID, nil
}

func (u *Usecase) Delete(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	return u.repository.Delete(ctx, req)
}

func (u *Usecase) FindAllSimple(ctx context.Context) ([]entity.TestType, error) {
	return u.repository.FindAllSimple(ctx)
}
