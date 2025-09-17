package patientuc

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
)

type PatientUseCase struct {
	cfg           *config.Schema
	patientRepo   *patientrepo.PatientRepository
	workOrderRepo *workOrderrepo.WorkOrderRepository
	validate      *validator.Validate
}

func NewPatientUseCase(
	cfg *config.Schema,
	patientRepo *patientrepo.PatientRepository,
	workOrderRepo *workOrderrepo.WorkOrderRepository,
	validate *validator.Validate,
) *PatientUseCase {
	return &PatientUseCase{
		cfg:           cfg,
		patientRepo:   patientRepo,
		workOrderRepo: workOrderRepo,
		validate:      validate,
	}
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

// GetPatientResultHistory returns the result history of a patient
func (p *PatientUseCase) GetPatientResultHistory(
	ctx context.Context,
	id int64,
	req entity.GetPatientRecordHistoryRequest,
) (entity.GetPatientResultHistoryResponse, error) {
	patient, err := p.patientRepo.FindOne(id)
	if err != nil {
		return entity.GetPatientResultHistoryResponse{}, err
	}

	workOrders, err := p.workOrderRepo.FindAllForResult(ctx, &entity.ResultGetManyRequest{
		PatientIDs: []int64{id},
		GetManyRequest: entity.GetManyRequest{
			CreatedAtStart: req.StartDate,
			CreatedAtEnd:   req.EndDate,
		},
	})
	if err != nil {
		return entity.GetPatientResultHistoryResponse{}, fmt.Errorf("failed to find work orders for patient %d: %w", id, err)
	}

	workOrders.Data = p.fillResultDetail(workOrders.Data)
	workOrders.Data = p.fillCalculatedResult(ctx, workOrders.Data)

	allTestResults := make([]entity.TestResult, 0)
	for i := range workOrders.Data {
		allTestResults = append(allTestResults, workOrders.Data[i].TestResult...)
	}

	return entity.GetPatientResultHistoryResponse{
		Patient:    patient,
		TestResult: allTestResults,
	}, nil
}

func (*PatientUseCase) fillResultDetail(workOrders []entity.WorkOrder) []entity.WorkOrder {
	for i := range workOrders {
		workOrders[i].FillResultDetail(entity.ResultDetailOption{
			HideEmpty:   true,
			HideHistory: true,
		})
	}

	return workOrders
}

func (p *PatientUseCase) fillCalculatedResult(ctx context.Context, workOrders []entity.WorkOrder) []entity.WorkOrder {
	for i := range workOrders {
		workOrders[i].CalculateEGFRForResults(ctx)
	}
	return workOrders
}
