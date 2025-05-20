package patientuc

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
)

type PatientUseCase struct {
	cfg         *config.Schema
	patientRepo *patientrepo.PatientRepository
	validate    *validator.Validate
}

func NewPatientUseCase(
	cfg *config.Schema,
	patientRepo *patientrepo.PatientRepository,
	validate *validator.Validate,
) *PatientUseCase {
	return &PatientUseCase{cfg: cfg, patientRepo: patientRepo, validate: validate}
}

func (p PatientUseCase) FindAll(
	ctx context.Context, req *entity.GetManyRequestPatient,
) (entity.PaginationResponse[entity.Patient], error) {
	return p.patientRepo.FindAll(ctx, req)
}

func (p PatientUseCase) FindManyByWorkOrderID(
	ctx context.Context, workOrderIDs []int64,
) ([]entity.Patient, error) {
	return p.patientRepo.FindManyByWorkOrderID(ctx, workOrderIDs)
}

func (p PatientUseCase) FindOneByID(id int64) (entity.Patient, error) {
	return p.patientRepo.FindOne(id)
}

func (p PatientUseCase) Create(req *entity.Patient) error {
	exists, err := p.patientRepo.IsExists(req)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf(
			"patient %s %s with %s already exists",
			req.FirstName,
			req.LastName,
			req.Birthdate.Format("2006-01-02"),
		)
	}

	return p.patientRepo.Create(req)
}

func (p PatientUseCase) Update(req *entity.Patient) error {
	return p.patientRepo.Update(req)
}

func (p PatientUseCase) Delete(id int64) error {
	return p.patientRepo.Delete(id)
}
