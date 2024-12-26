package test_type

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
)

type Usecase struct {
	TestTypeRepository *test_type.Repository
}

func NewUsecase(testTypeRepository *test_type.Repository) *Usecase {
	return &Usecase{TestTypeRepository: testTypeRepository}
}

func (u *Usecase) FindAll(ctx context.Context, req *entity.TestTypeGetManyRequest) ([]entity.TestType, error) {
	return u.TestTypeRepository.FindAll(ctx, req)
}
