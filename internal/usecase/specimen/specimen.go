package specimenuc

import (
	"context"
	"time"

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
	return p.SpecimenRepo.FindOne(id)
}

func (p SpecimenUseCase) Create(req *entity.Specimen) error {
	req.Barcode = p.generateBarcode(context.Background(), req)

	return p.SpecimenRepo.Create(req)
}

func (p SpecimenUseCase) generateBarcode(ctx context.Context, req *entity.Specimen) string {
	return time.Now().Format("20060102150405")
}

func (p SpecimenUseCase) Update(req *entity.Specimen) error {
	return p.SpecimenRepo.Update(req)
}

func (p SpecimenUseCase) Delete(id int) error {
	return p.SpecimenRepo.Delete(id)
}
