package specimentuc

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	specimentrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/speciment"
)

type SpecimentUseCase struct {
	cfg           *config.Schema
	specimentRepo *specimentrepo.SpecimentRepository
	validate      *validator.Validate
}

func NewSpecimentUseCase(
	cfg *config.Schema,
	specimentRepo *specimentrepo.SpecimentRepository,
	validate *validator.Validate,
) *SpecimentUseCase {
	return &SpecimentUseCase{cfg: cfg, specimentRepo: specimentRepo, validate: validate}
}

func (p SpecimentUseCase) FindAll(
	ctx context.Context, req *entity.SpecimentGetManyRequest,
) ([]entity.Speciment, error) {
	return p.specimentRepo.FindAll(ctx, req)
}

func (p SpecimentUseCase) FindOneByID(id int64) (entity.Speciment, error) {
	return p.specimentRepo.FindOne(id)
}

func (p SpecimentUseCase) Create(req *entity.Speciment) error {
	req.Barcode = p.generateBarcode(context.Background(), req)

	return p.specimentRepo.Create(req)
}

func (p SpecimentUseCase) generateBarcode(ctx context.Context, req *entity.Speciment) string {
	return time.Now().Format("20060102150405")
}

func (p SpecimentUseCase) Update(req *entity.Speciment) error {
	return p.specimentRepo.Update(req)
}

func (p SpecimentUseCase) Delete(id int64) error {
	return p.specimentRepo.Delete(id)
}
