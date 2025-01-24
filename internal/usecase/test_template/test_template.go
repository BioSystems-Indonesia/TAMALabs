package test_template_uc

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_template"
)

type Usecase struct {
	repository *test_template.Repository
}

func NewUsecase(testTypeRepository *test_template.Repository) *Usecase {
	return &Usecase{repository: testTypeRepository}
}

func (u *Usecase) FindAll(ctx context.Context, req *entity.TestTemplateGetManyRequest) (entity.PaginationResponse[entity.TestTemplate], error) {
	return u.repository.FindAll(ctx, req)
}

func (u *Usecase) FindOneByID(ctx context.Context, id int) (entity.TestTemplate, error) {
	return u.repository.FindOneByID(ctx, id)
}

func (u *Usecase) Create(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	return u.repository.Create(ctx, req)
}

func (u *Usecase) Update(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	return u.repository.Update(ctx, req)
}

func (u *Usecase) Delete(ctx context.Context, req *entity.TestTemplate) (entity.TestTemplate, error) {
	return u.repository.Delete(ctx, req)
}
