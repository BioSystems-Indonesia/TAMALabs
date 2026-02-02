package sub_category

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	subcategoryrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/sub_category"
)

type Usecase struct {
	repo *subcategoryrepo.Repository
}

func NewUsecase(repo *subcategoryrepo.Repository) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

// FindAll returns all sub-categories
func (u *Usecase) FindAll(ctx context.Context, req *entity.SubCategoryGetManyRequest) (entity.PaginationResponse[entity.SubCategory], error) {
	return u.repo.FindAll(ctx, req)
}

// FindByID finds a sub-category by ID
func (u *Usecase) FindByID(ctx context.Context, id int) (*entity.SubCategory, error) {
	return u.repo.FindByID(ctx, id)
}

// GetTestTypesBySubCategoryID returns all test types for a sub-category
func (u *Usecase) GetTestTypesBySubCategoryID(ctx context.Context, subCategoryID int) ([]entity.TestType, error) {
	return u.repo.GetTestTypesBySubCategoryID(ctx, subCategoryID)
}
