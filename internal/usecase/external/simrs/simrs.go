package simrsuc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	simrsrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/external/simrs"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	testType "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrder "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/result"
	workOrderUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order"
)

type Usecase struct {
	simrsRepo     *simrsrepo.Repository
	workOrderRepo *workOrder.WorkOrderRepository
	workOrderUC   *workOrderUC.WorkOrderUseCase
	patientRepo   *patientrepo.PatientRepository
	testTypeRepo  *testType.Repository
	cfg           *config.Schema
	resultUC      *result.Usecase
}

func NewUsecase(simrsRepo *simrsrepo.Repository, workOrderRepo *workOrder.WorkOrderRepository,
	workOrderUC *workOrderUC.WorkOrderUseCase, patientRepo *patientrepo.PatientRepository, testTypeRepo *testType.Repository, cfg *config.Schema, resultUC *result.Usecase) *Usecase {
	return &Usecase{
		simrsRepo:     simrsRepo,
		workOrderRepo: workOrderRepo,
		workOrderUC:   workOrderUC,
		patientRepo:   patientRepo,
		testTypeRepo:  testTypeRepo,
		cfg:           cfg,
		resultUC:      resultUC,
	}
}

func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	if u.simrsRepo == nil {
		return fmt.Errorf("SIMRS repository is not initialized - SIMRS integration may be disabled or misconfigured")
	}

	slog.Info("Starting SIMRS lab request synchronization")

	simrsLabRequests, err := u.simrsRepo.GetAllLabRequests(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SIMRS lab requests: %w", err)
	}

	slog.Info(fmt.Sprintf("Found %d lab requests from SIMRS", len(simrsLabRequests)))

	var processedNoOrders []string
	var processedPatientIDs []string

	for _, simrsLabRequest := range simrsLabRequests {
		err := u.processLabRequest(ctx, simrsLabRequest)
		if err != nil {
			slog.Error("Failed to process lab request", "no_order", simrsLabRequest.NoOrder, "error", err)
			continue
		}
		processedNoOrders = append(processedNoOrders, simrsLabRequest.NoOrder)
		processedPatientIDs = append(processedPatientIDs, simrsLabRequest.PatientID)
	}

	if len(processedNoOrders) > 0 {
		err = u.simrsRepo.DeleteProcessedLabRequests(ctx, processedNoOrders)
		if err != nil {
			slog.Error("Failed to delete processed lab requests", "error", err)
		}
	}

	if len(processedPatientIDs) > 0 {
		uniquePatientIDs := removeDuplicates(processedPatientIDs)

		err = u.simrsRepo.DeleteProcessedPatients(ctx, uniquePatientIDs)
		if err != nil {
			slog.Error("Failed to delete processed patients", "error", err)
		}
	}

	slog.Info("SIMRS lab request synchronization completed", "processed", len(processedNoOrders), "total", len(simrsLabRequests))
	return nil
}

func (u *Usecase) processLabRequest(ctx context.Context, simrsLabRequest entity.SimrsLabRequest) error {
	simrsPatient, err := u.simrsRepo.GetPatientByID(ctx, simrsLabRequest.PatientID)
	if err != nil {
		return fmt.Errorf("failed to get SIMRS patient: %w", err)
	}

	patient, err := u.createOrUpdatePatient(ctx, simrsPatient)
	if err != nil {
		return fmt.Errorf("failed to create/update patient: %w", err)
	}

	var paramCodes []string
	if err := json.Unmarshal([]byte(simrsLabRequest.ParamRequest), &paramCodes); err != nil {
		return fmt.Errorf("failed to parse param_request: %w", err)
	}

	err = u.createWorkOrder(ctx, simrsLabRequest, patient.ID, paramCodes)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	slog.Info("Successfully processed lab request", "no_order", simrsLabRequest.NoOrder, "patient_id", simrsLabRequest.PatientID)
	return nil
}

func (u *Usecase) createOrUpdatePatient(ctx context.Context, simrsPatient *entity.SimrsPatient) (*entity.Patient, error) {
	patient := &entity.Patient{
		SIMRSPID:    sql.NullString{String: simrsPatient.PatientID, Valid: true},
		FirstName:   simrsPatient.FirstName,
		LastName:    simrsPatient.LastName,
		Birthdate:   simrsPatient.Birthdate,
		Sex:         entity.SimrsGender(simrsPatient.Gender).ToPatientSex(),
		Address:     simrsPatient.Address,
		PhoneNumber: simrsPatient.Phone,
		Location:    simrsPatient.Address,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := u.patientRepo.FirstOrCreate(patient)
	if err != nil {
		if err.Error() == "constraint failed: UNIQUE constraint failed: patients.simrs_pid (2067)" ||
			err.Error() == "UNIQUE constraint failed: patients.simrs_pid" ||
			err.Error() == "constraint failed: UNIQUE constraint failed: patients.simrs_pid" {

			slog.Info("Patient with SIMRS PID already exists, searching for existing patient", "simrs_pid", simrsPatient.PatientID)

			findReq := &entity.GetManyRequestPatient{
				GetManyRequest: entity.GetManyRequest{
					Start: 0,
					End:   1000, // Get first 1000 patients
				},
			}

			paginatedResult, err := u.patientRepo.FindAll(ctx, findReq)
			if err != nil {
				return nil, fmt.Errorf("failed to search for existing patient: %w", err)
			}

			var existingPatient *entity.Patient
			for _, p := range paginatedResult.Data {
				if p.SIMRSPID.Valid && p.SIMRSPID.String == simrsPatient.PatientID {
					existingPatient = &p
					break
				}
			}

			if existingPatient != nil {
				slog.Info("Found existing patient, updating", "patient_id", existingPatient.ID, "simrs_pid", simrsPatient.PatientID)

				existingPatient.FirstName = simrsPatient.FirstName
				existingPatient.LastName = simrsPatient.LastName
				existingPatient.Birthdate = simrsPatient.Birthdate
				existingPatient.Sex = entity.SimrsGender(simrsPatient.Gender).ToPatientSex()
				existingPatient.Address = simrsPatient.Address
				existingPatient.PhoneNumber = simrsPatient.Phone
				existingPatient.Location = simrsPatient.Address
				existingPatient.UpdatedAt = time.Now()

				err := u.patientRepo.Update(existingPatient)
				if err != nil {
					return nil, fmt.Errorf("failed to update existing patient: %w", err)
				}

				return existingPatient, nil
			}

			return nil, fmt.Errorf("patient with SIMRS PID %s exists but could not be found in search results", simrsPatient.PatientID)
		}

		return nil, fmt.Errorf("failed to create/find patient: %w", err)
	}

	slog.Info("Patient created or found successfully", "patient_id", result.ID, "simrs_pid", simrsPatient.PatientID)
	return &result, nil
}

func (u *Usecase) createWorkOrder(ctx context.Context, simrsLabRequest entity.SimrsLabRequest, patientID int64, paramCodes []string) error {
	var testTypes []entity.WorkOrderCreateRequestTestType
	for _, paramCode := range paramCodes {
		testType, err := u.testTypeRepo.FindOneByCode(ctx, paramCode)
		if err != nil {
			slog.Warn("Test type not found", "param_code", paramCode)
			continue
		}

		specimenType := testType.GetFirstType()

		testTypeReq := entity.WorkOrderCreateRequestTestType{
			TestTypeID:   int64(testType.ID),
			TestTypeCode: testType.Code,
			SpecimenType: specimenType,
		}
		testTypes = append(testTypes, testTypeReq)
	}

	if len(testTypes) == 0 {
		return fmt.Errorf("no valid test types found for param codes: %v", paramCodes)
	}

	workOrderReq := &entity.WorkOrderCreateRequest{
		PatientID:    patientID,
		TestTypes:    testTypes,
		CreatedBy:    1,
		BarcodeSIMRS: simrsLabRequest.NoOrder,
		Barcode:      "",
	}

	_, err := u.workOrderUC.Create(workOrderReq)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	return nil
}

func (u *Usecase) SyncAllResult(ctx context.Context, orderIDs []int64) error {
	orderIDStrings := make([]string, len(orderIDs))
	for i, orderID := range orderIDs {
		orderIDStrings[i] = strconv.Itoa(int(orderID))
	}

	workOrders, err := u.workOrderRepo.FindAll(ctx, &entity.WorkOrderGetManyRequest{
		GetManyRequest: entity.GetManyRequest{
			CreatedAtStart: time.Now().Add(14 * -24 * time.Hour),
			CreatedAtEnd:   time.Now(),
		},
	})
	if err != nil {
		return err
	}

	var errs []error
	for _, workOrder := range workOrders.Data {
		err := u.syncResultByOrderID(ctx, workOrder.ID)
		if err != nil {
			slog.Error("error syncing result", "error", err)
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (u *Usecase) syncResultByOrderID(ctx context.Context, orderID int64) error {
	workOrder, err := u.resultUC.ResultDetail(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get work order: %w", err)
	}

	if workOrder.BarcodeSIMRS == "" {
		slog.Debug("Work order does not have SIMRS barcode", "order_id", orderID)
		return nil
	}

	workOrder.FillTestResultDetail(false)

	if len(workOrder.TestResult) == 0 {
		slog.Debug("No test results found", "order_id", orderID)
		return nil
	}

	var simrsLabResults []entity.SimrsLabResult
	for _, testResultGroup := range workOrder.TestResult {
		for _, testResult := range testResultGroup {
			if testResult.Result == "" {
				continue
			}

			// Skip if TestTypeID is 0 (invalid)
			if testResult.TestTypeID == 0 {
				slog.Warn("Skipping test result with invalid TestTypeID", "test_type_id", testResult.TestTypeID)
				continue
			}

			testType, err := u.testTypeRepo.FindOneByID(ctx, int(testResult.TestTypeID))
			if err != nil {
				slog.Warn("Failed to get test type", "test_type_id", testResult.TestTypeID)
				continue
			}

			simrsLabResult := entity.SimrsLabResult{
				NoOrder:     workOrder.BarcodeSIMRS,
				ParamCode:   testType.Code,
				ResultValue: testResult.Result,
				Unit:        testResult.Unit,
				RefRange:    testResult.ReferenceRange,
				Flag:        string(entity.NewSimrsFlag(testResult)),
				CreatedAt:   time.Now(),
			}

			simrsLabResults = append(simrsLabResults, simrsLabResult)
		}
	}

	if len(simrsLabResults) == 0 {
		slog.Debug("No valid SIMRS lab results to sync", "order_id", orderID)
		return nil
	}

	err = u.simrsRepo.BatchInsertLabResults(ctx, simrsLabResults)
	if err != nil {
		return fmt.Errorf("failed to batch insert lab results: %w", err)
	}

	slog.Info("Successfully synced lab results to SIMRS", "order_id", orderID, "results_count", len(simrsLabResults))
	return nil
}

func (u *Usecase) GetLabRequestByNoOrder(ctx context.Context, noOrder string) (*entity.SimrsLabRequest, error) {
	return u.simrsRepo.GetLabRequestByNoOrder(ctx, noOrder)
}

func (u *Usecase) GetPatientByID(ctx context.Context, patientID string) (*entity.SimrsPatient, error) {
	return u.simrsRepo.GetPatientByID(ctx, patientID)
}

func (u *Usecase) GetLabResultsByNoOrder(ctx context.Context, noOrder string) ([]entity.SimrsLabResult, error) {
	return u.simrsRepo.GetLabResultsByNoOrder(ctx, noOrder)
}

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}

	return result
}
