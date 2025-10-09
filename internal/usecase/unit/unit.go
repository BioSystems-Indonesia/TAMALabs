package unit

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	unitRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/unit"
	"github.com/go-playground/validator/v10"
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
