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
) ([]entity.Specimen, error) {
	return p.SpecimenRepo.FindAll(ctx, req)
}

func (p SpecimenUseCase) FindOneByID(id int64) (entity.Specimen, error) {
	return p.SpecimenRepo.FindOne(context.TODO(), id)
}

func (p SpecimenUseCase) Create(req *entity.Specimen) error {
	req.Barcode = entity.GenerateBarcode()

	return p.SpecimenRepo.Create(context.TODO(), req)
}

func (p SpecimenUseCase) Update(req *entity.Specimen) error {
	return p.SpecimenRepo.Update(context.TODO(), req)
}

func (p SpecimenUseCase) Delete(id int) error {
	return p.SpecimenRepo.Delete(id)
}
