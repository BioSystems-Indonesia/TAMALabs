package test_type

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
)

type Usecase struct {
	TestTypeRepository repository.TestType
}

func NewUsecase(testTypeRepository repository.TestType) *Usecase {
	return &Usecase{TestTypeRepository: testTypeRepository}
}

func (u *Usecase) FindAll(ctx context.Context, req *entity.TestTypeGetManyRequest) ([]entity.TestType, error) {
	return u.TestTypeRepository.FindAll(ctx, req)
}
