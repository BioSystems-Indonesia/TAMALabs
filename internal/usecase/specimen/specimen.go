package specimenuc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
)

type SpecimenUseCase struct {
	cfg          *config.Schema
	SpecimenRepo *specimen.Repository
	validate     *validator.Validate
}

func NewSpecimenUseCase(
	cfg *config.Schema,
	SpecimenRepo *specimen.Repository,
	validate *validator.Validate,
) *SpecimenUseCase {
	return &SpecimenUseCase{cfg: cfg, SpecimenRepo: SpecimenRepo, validate: validate}
}

func (p SpecimenUseCase) FindAll(
	ctx context.Context, req *entity.SpecimenGetManyRequest,
) (entity.PaginationResponse[entity.Specimen], error) {
	return p.SpecimenRepo.FindAll(ctx, req)
}

func (p SpecimenUseCase) FindOneByID(id int64) (entity.Specimen, error) {
	return p.SpecimenRepo.FindOne(context.TODO(), id)
}

func (p SpecimenUseCase) FindAllByWorkOrderIDs(ctx context.Context, workOrderIDs []int64) ([]entity.Specimen, error) {
	return p.SpecimenRepo.FindAllByWorkOrderIDs(ctx, workOrderIDs)
}

func (p SpecimenUseCase) BulkUpdateSpecimen(ctx context.Context, specimens []entity.Specimen) error {
	return p.SpecimenRepo.BulkUpdate(ctx, specimens)
}
