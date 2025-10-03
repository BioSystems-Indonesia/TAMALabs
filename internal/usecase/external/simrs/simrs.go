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

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	simrsrepo "github.com/oibacidem/lims-hl-seven/internal/repository/external/simrs"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	testType "github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	workOrder "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/result"
	workOrderUC "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
)

// Usecase handles SIMRS integration business logic
type Usecase struct {
	simrsRepo     *simrsrepo.Repository
	workOrderRepo *workOrder.WorkOrderRepository
	workOrderUC   *workOrderUC.WorkOrderUseCase
	patientRepo   *patientrepo.PatientRepository
	testTypeRepo  *testType.Repository
	cfg           *config.Schema
	resultUC      *result.Usecase
}

// NewUsecase creates a new SIMRS usecase instance
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

// SyncAllRequest synchronizes all lab requests from SIMRS to internal LIMS
func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	// Check if SIMRS repository is properly initialized
	if u.simrsRepo == nil {
		return fmt.Errorf("SIMRS repository is not initialized - SIMRS integration may be disabled or misconfigured")
	}

	slog.Info("Starting SIMRS lab request synchronization")

	// Get all lab requests from SIMRS database
	simrsLabRequests, err := u.simrsRepo.GetAllLabRequests(ctx)
	if err != nil {
		return fmt.Errorf("failed to get SIMRS lab requests: %w", err)
	}

	slog.Info(fmt.Sprintf("Found %d lab requests from SIMRS", len(simrsLabRequests)))

	// Track successfully processed lab requests and patients for deletion
	var processedNoOrders []string
	var processedPatientIDs []string

	for _, simrsLabRequest := range simrsLabRequests {
		err := u.processLabRequest(ctx, simrsLabRequest)
		if err != nil {
			slog.Error("Failed to process lab request", "no_order", simrsLabRequest.NoOrder, "error", err)
			continue
		}
		// Add to processed lists for deletion
		processedNoOrders = append(processedNoOrders, simrsLabRequest.NoOrder)
		processedPatientIDs = append(processedPatientIDs, simrsLabRequest.PatientID)
	}

	// Delete successfully processed lab requests from SIMRS
	if len(processedNoOrders) > 0 {
		err = u.simrsRepo.DeleteProcessedLabRequests(ctx, processedNoOrders)
		if err != nil {
			slog.Error("Failed to delete processed lab requests", "error", err)
			// Don't return error here as the main sync was successful
		}
	}

	// Delete successfully processed patients from SIMRS
	if len(processedPatientIDs) > 0 {
		// Remove duplicates from patient IDs
		uniquePatientIDs := removeDuplicates(processedPatientIDs)

		err = u.simrsRepo.DeleteProcessedPatients(ctx, uniquePatientIDs)
		if err != nil {
			slog.Error("Failed to delete processed patients", "error", err)
			// Don't return error here as the main sync was successful
		}
	}

	slog.Info("SIMRS lab request synchronization completed", "processed", len(processedNoOrders), "total", len(simrsLabRequests))
	return nil
}

// processLabRequest processes a single lab request from SIMRS
func (u *Usecase) processLabRequest(ctx context.Context, simrsLabRequest entity.SimrsLabRequest) error {
	// Check if work order already exists by checking barcode_simrs field
	// Since we don't have GetWorkOrderByBarcodeSimrs method, we'll use a different approach
	// For now, we'll always create new work orders, but in production you might want to add this check

	// Get patient information from SIMRS
	simrsPatient, err := u.simrsRepo.GetPatientByID(ctx, simrsLabRequest.PatientID)
	if err != nil {
		return fmt.Errorf("failed to get SIMRS patient: %w", err)
	}

	// Create or update patient in internal LIMS
	patient, err := u.createOrUpdatePatient(ctx, simrsPatient)
	if err != nil {
		return fmt.Errorf("failed to create/update patient: %w", err)
	}

	// Parse parameter requests
	var paramCodes []string
	if err := json.Unmarshal([]byte(simrsLabRequest.ParamRequest), &paramCodes); err != nil {
		return fmt.Errorf("failed to parse param_request: %w", err)
	}

	// Create work order
	err = u.createWorkOrder(ctx, simrsLabRequest, patient.ID, paramCodes)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	slog.Info("Successfully processed lab request", "no_order", simrsLabRequest.NoOrder, "patient_id", simrsLabRequest.PatientID)
	return nil
}

// createOrUpdatePatient creates or updates patient in internal LIMS
func (u *Usecase) createOrUpdatePatient(ctx context.Context, simrsPatient *entity.SimrsPatient) (*entity.Patient, error) {
	// Try to find existing patient by SIMRS PID using FirstOrCreate approach
	patient := &entity.Patient{
		SIMRSPID:    sql.NullString{String: simrsPatient.PatientID, Valid: true},
		FirstName:   simrsPatient.FirstName,
		LastName:    simrsPatient.LastName,
		Birthdate:   simrsPatient.Birthdate,
		Sex:         entity.SimrsGender(simrsPatient.Gender).ToPatientSex(),
		Address:     simrsPatient.Address,
		PhoneNumber: simrsPatient.Phone,
		Location:    simrsPatient.Address, // Use address as location for now
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Use FirstOrCreate to either find existing or create new patient
	result, err := u.patientRepo.FirstOrCreate(patient)
	if err != nil {
		return nil, fmt.Errorf("failed to create/update patient: %w", err)
	}

	return &result, nil
}

// createWorkOrder creates a work order from SIMRS lab request
func (u *Usecase) createWorkOrder(ctx context.Context, simrsLabRequest entity.SimrsLabRequest, patientID int64, paramCodes []string) error {
	// Get test types by codes
	var testTypes []entity.WorkOrderCreateRequestTestType
	for _, paramCode := range paramCodes {
		testType, err := u.testTypeRepo.FindOneByCode(ctx, paramCode)
		if err != nil {
			slog.Warn("Test type not found", "param_code", paramCode)
			continue
		}

		// Get the first specimen type from the test type
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

	// Create work order request
	workOrderReq := &entity.WorkOrderCreateRequest{
		PatientID:    patientID,
		TestTypes:    testTypes,
		CreatedBy:    1, // Default admin user, should be configurable
		BarcodeSIMRS: simrsLabRequest.NoOrder,
		Barcode:      "", // Will be auto-generated by the system
	}

	// Create work order using the usecase's Create method (which auto-generates barcode)
	_, err := u.workOrderUC.Create(workOrderReq)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	return nil
}

// SyncAllResult synchronizes all lab results from internal LIMS to SIMRS
func (u *Usecase) SyncAllResult(ctx context.Context, orderIDs []int64) error {
	orderIDStrings := make([]string, len(orderIDs))
	for i, orderID := range orderIDs {
		orderIDStrings[i] = strconv.Itoa(int(orderID))
	}

	workOrders, err := u.workOrderRepo.FindAll(ctx, &entity.WorkOrderGetManyRequest{
		GetManyRequest: entity.GetManyRequest{
			CreatedAtStart: time.Now().Add(14 * -24 * time.Hour),
			CreatedAtEnd:   time.Now(),
			ID:             orderIDStrings,
		},
	})
	if err != nil {
		return err
	}

	var errs []error
	for _, workOrder := range workOrders.Data {
		fmt.Println(workOrder.ID)
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

// syncResultByOrderID synchronizes lab results for a specific work order
func (u *Usecase) syncResultByOrderID(ctx context.Context, orderID int64) error {
	// Get work order using FindOne method
	workOrder, err := u.resultUC.ResultDetail(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get work order: %w", err)
	}

	// Check if work order has SIMRS barcode
	if workOrder.BarcodeSIMRS == "" {
		slog.Debug("Work order does not have SIMRS barcode", "order_id", orderID)
		return nil
	}

	// Fill test result details for the work order
	workOrder.FillTestResultDetail(false)

	if len(workOrder.TestResult) == 0 {
		slog.Debug("No test results found", "order_id", orderID)
		return nil
	}

	// Convert test results to SIMRS lab results
	var simrsLabResults []entity.SimrsLabResult
	for _, testResultGroup := range workOrder.TestResult {
		// testResultGroup is []entity.TestResult, so we need to iterate through it
		for _, testResult := range testResultGroup {
			// Skip if no result value
			if testResult.Result == nil {
				continue
			}

			// Get test type information using FindOneByID
			testType, err := u.testTypeRepo.FindOneByID(ctx, int(testResult.TestTypeID))
			if err != nil {
				slog.Warn("Failed to get test type", "test_type_id", testResult.TestTypeID)
				continue
			}

			simrsLabResult := entity.SimrsLabResult{
				NoOrder:     workOrder.BarcodeSIMRS,
				ParamCode:   testType.Code,
				ResultValue: fmt.Sprintf("%.2f", *testResult.Result),
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

	// Batch insert lab results to SIMRS
	err = u.simrsRepo.BatchInsertLabResults(ctx, simrsLabResults)
	if err != nil {
		return fmt.Errorf("failed to batch insert lab results: %w", err)
	}

	slog.Info("Successfully synced lab results to SIMRS", "order_id", orderID, "results_count", len(simrsLabResults))
	return nil
}

// GetLabRequestByNoOrder retrieves lab request details from SIMRS
func (u *Usecase) GetLabRequestByNoOrder(ctx context.Context, noOrder string) (*entity.SimrsLabRequest, error) {
	return u.simrsRepo.GetLabRequestByNoOrder(ctx, noOrder)
}

// GetPatientByID retrieves patient information from SIMRS
func (u *Usecase) GetPatientByID(ctx context.Context, patientID string) (*entity.SimrsPatient, error) {
	return u.simrsRepo.GetPatientByID(ctx, patientID)
}

// GetLabResultsByNoOrder retrieves lab results from SIMRS
func (u *Usecase) GetLabResultsByNoOrder(ctx context.Context, noOrder string) ([]entity.SimrsLabResult, error) {
	return u.simrsRepo.GetLabResultsByNoOrder(ctx, noOrder)
}

// removeDuplicates removes duplicate strings from a slice
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
