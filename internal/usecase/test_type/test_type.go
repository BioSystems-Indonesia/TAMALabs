package test_type

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
)

type Usecase struct {
	repository *test_type.Repository
}

func NewUsecase(testTypeRepository *test_type.Repository) *Usecase {
	return &Usecase{repository: testTypeRepository}
}

func (u *Usecase) FindAll(ctx context.Context, req *entity.TestTypeGetManyRequest) (entity.PaginationResponse[entity.TestType], error) {
	return u.repository.FindAll(ctx, req)
}

func (u *Usecase) FindOneByID(ctx context.Context, id int) (entity.TestType, error) {
	return u.repository.FindOneByID(ctx, id)
}

func (u *Usecase) Create(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	return u.repository.Create(ctx, req)
}

func (u *Usecase) Update(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	return u.repository.Update(ctx, req)
}

func (u *Usecase) Delete(ctx context.Context, req *entity.TestType) (entity.TestType, error) {
	return u.repository.Delete(ctx, req)
}
