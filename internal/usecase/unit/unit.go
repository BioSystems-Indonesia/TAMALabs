package unit

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	unitRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/unit"
)

type UnitUseCase struct {
	cfg      *config.Schema
	unitRepo *unitRepo.Repository
	validate *validator.Validate
}

func NewUnitUseCase(
	cfg *config.Schema,
	unitRepo *unitRepo.Repository,
	validate *validator.Validate,
) *UnitUseCase {
	return &UnitUseCase{cfg: cfg, unitRepo: unitRepo, validate: validate}
}

func (p UnitUseCase) FindAll(
	ctx context.Context, req *entity.UnitGetManyRequest,
) (entity.PaginationResponse[entity.Unit], error) {
	return p.unitRepo.FindAll(ctx, req)
}
