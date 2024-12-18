package specimentuc

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	specimentrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/speciment"
)

type SpecimenUseCase struct {
	cfg           *config.Schema
	specimentRepo *specimentrepo.SpecimentRepository
	validate      *validator.Validate
}

func NewSpecimenUseCase(
	cfg *config.Schema,
	specimentRepo *specimentrepo.SpecimentRepository,
	validate *validator.Validate,
) *SpecimenUseCase {
	return &SpecimenUseCase{cfg: cfg, specimentRepo: specimentRepo, validate: validate}
}

func (p SpecimenUseCase) FindAll(
	ctx context.Context, req *entity.SpecimentGetManyRequest,
) ([]entity.Specimen, error) {
	return p.specimentRepo.FindAll(ctx, req)
}

func (p SpecimenUseCase) FindOneByID(id int64) (entity.Specimen, error) {
	return p.specimentRepo.FindOne(id)
}

func (p SpecimenUseCase) Create(req *entity.Specimen) error {
	req.Barcode = p.generateBarcode(context.Background(), req)

	return p.specimentRepo.Create(req)
}

func (p SpecimenUseCase) generateBarcode(ctx context.Context, req *entity.Specimen) string {
	return time.Now().Format("20060102150405")
}

func (p SpecimenUseCase) Update(req *entity.Specimen) error {
	return p.specimentRepo.Update(req)
}

func (p SpecimenUseCase) Delete(id int64) error {
	return p.specimentRepo.Delete(id)
}
