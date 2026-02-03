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
	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
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
	lastDSN       string // Track last used DSN to detect changes
}

func NewUsecase(simrsRepo *simrsrepo.Repository, workOrderRepo *workOrder.WorkOrderRepository,
	workOrderUC *workOrderUC.WorkOrderUseCase, patientRepo *patientrepo.PatientRepository, testTypeRepo *testType.Repository, cfg *config.Schema, resultUC *result.Usecase) *Usecase {
	lastDSN := ""
	if simrsRepo != nil {
		lastDSN = cfg.SimrsDatabaseDSN
	}
	return &Usecase{
		simrsRepo:     simrsRepo,
		workOrderRepo: workOrderRepo,
		workOrderUC:   workOrderUC,
		patientRepo:   patientRepo,
		testTypeRepo:  testTypeRepo,
		cfg:           cfg,
		resultUC:      resultUC,
		lastDSN:       lastDSN,
	}
}

// ensureRepository checks if repository needs to be reinitialized due to DSN change
func (u *Usecase) ensureRepository() error {
	currentDSN := u.cfg.SimrsDatabaseDSN

	// If DSN changed or repository is nil, reinitialize
	if u.simrsRepo == nil || u.lastDSN != currentDSN {
		if currentDSN == "" {
			u.simrsRepo = nil
			u.lastDSN = ""
			return fmt.Errorf("SIMRS database DSN is not configured")
		}

		slog.Info("Reinitializing SIMRS repository", "dsn_changed", u.lastDSN != currentDSN, "repo_was_nil", u.simrsRepo == nil)

		db, err := simrsrepo.NewDB(currentDSN)
		if err != nil {
			return fmt.Errorf("failed to create SIMRS database connection: %w", err)
		}

		u.simrsRepo = simrsrepo.NewRepository(db)
		u.lastDSN = currentDSN
		slog.Info("SIMRS repository reinitialized successfully")
	}

	return nil
}

func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	// Ensure repository is initialized and up-to-date
	if err := u.ensureRepository(); err != nil {
		return fmt.Errorf("SIMRS repository initialization failed: %w", err)
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
	// Ensure repository is initialized and up-to-date
	if err := u.ensureRepository(); err != nil {
		return fmt.Errorf("SIMRS repository initialization failed: %w", err)
	}

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

// ProcessOrder handles incoming lab order from SIMRS via API
func (u *Usecase) ProcessOrder(ctx context.Context, rawRequest []byte) error {
	var request entity.SimrsOrderRequest
	var err error

	slog.InfoContext(ctx, "debug simrs process order via API", "rawRequest", string(rawRequest))

	err = json.Unmarshal(rawRequest, &request)
	if err != nil {
		return fmt.Errorf("simrs.ProcessOrder json.Unmarshal failed: %w\nbody: %s", err, string(rawRequest))
	}

	// Validate request
	if err := validateOrderRequest(request); err != nil {
		return fmt.Errorf("simrs.ProcessOrder validation failed: %w", err)
	}

	// Find or create patient
	patient, err := u.findOrCreatePatientFromAPI(ctx, request)
	if err != nil {
		return fmt.Errorf("simrs.ProcessOrder findOrCreatePatient failed: %w", err)
	}

	// Create work order
	if err := u.createWorkOrderFromAPI(ctx, request, patient); err != nil {
		return fmt.Errorf("simrs.ProcessOrder createWorkOrder failed: %w", err)
	}

	slog.InfoContext(ctx, "SIMRS order processed successfully", "order_lab", request.Order.OBR.OrderLab, "patient_id", patient.ID)
	return nil
}

func validateOrderRequest(req entity.SimrsOrderRequest) error {
	if req.Order.PID.Pname == "" {
		return errors.New("patient name is required")
	}
	if req.Order.PID.Sex == "" {
		return errors.New("patient sex is required")
	}
	if req.Order.PID.BirthDt == "" {
		return errors.New("patient birth date is required")
	}
	if req.Order.OBR.OrderLab == "" {
		return errors.New("order lab number is required")
	}
	if len(req.Order.OBR.OrderTest) == 0 {
		return errors.New("at least one test is required")
	}
	return nil
}

func (u *Usecase) findOrCreatePatientFromAPI(ctx context.Context, req entity.SimrsOrderRequest) (entity.Patient, error) {
	// Check if patient repository is initialized
	if u.patientRepo == nil {
		return entity.Patient{}, fmt.Errorf("patient repository is not initialized - this may occur when using SIMRS API mode without database configuration")
	}

	// Split patient name
	var firstName, lastName string
	nameParts := splitName(req.Order.PID.Pname)
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = nameParts[1]
	}

	sex := entity.PatientSexUnknown
	switch req.Order.PID.Sex {
	case "M":
		sex = entity.PatientSexMale
	case "F":
		sex = entity.PatientSexFemale
	}

	// Parse birthdate - format: dd.MM.yyyy (21.07.1985)
	birthdate, err := time.Parse("02.01.2006", req.Order.PID.BirthDt)
	if err != nil {
		return entity.Patient{}, fmt.Errorf("cannot parse birth date (expected format dd.MM.yyyy): %w", err)
	}

	// First, try to find existing patient by Medical Record Number if provided
	if req.Order.PID.MedicalRecordNumber != "" {
		var existingPatient entity.Patient
		err := u.patientRepo.FindByMedicalRecordNumber(ctx, req.Order.PID.MedicalRecordNumber, &existingPatient)
		if err == nil {
			// Patient found by Medical Record Number, update the info
			slog.InfoContext(ctx, "Patient found by Medical Record Number",
				"patient_id", existingPatient.ID,
				"medical_record_number", req.Order.PID.MedicalRecordNumber)

			// Update patient information
			existingPatient.FirstName = firstName
			existingPatient.LastName = lastName
			existingPatient.Sex = sex
			existingPatient.Birthdate = birthdate
			existingPatient.UpdatedAt = time.Now()

			err = u.patientRepo.Update(&existingPatient)
			if err != nil {
				return entity.Patient{}, fmt.Errorf("error updating patient: %w", err)
			}

			return existingPatient, nil
		}
		// If not found by MRM, continue to create new patient
		slog.InfoContext(ctx, "Patient not found by Medical Record Number, will create new",
			"medical_record_number", req.Order.PID.MedicalRecordNumber)
	}

	// Create new patient with Medical Record Number
	patient := entity.Patient{
		MedicalRecordNumber: req.Order.PID.MedicalRecordNumber,
		FirstName:           firstName,
		LastName:            lastName,
		Sex:                 sex,
		Birthdate:           birthdate,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Use FirstOrCreate to find or create patient based on name and birthdate
	result, err := u.patientRepo.FirstOrCreate(&patient)
	if err != nil {
		return entity.Patient{}, fmt.Errorf("error finding/creating patient: %w", err)
	}

	if result.ID == 0 {
		return result, errors.New("patient id is zero, id is not filled on DB")
	}

	slog.InfoContext(ctx, "Patient found/created from API", "patient_id", result.ID, "name", req.Order.PID.Pname)
	return result, nil
}

func (u *Usecase) createWorkOrderFromAPI(ctx context.Context, req entity.SimrsOrderRequest, patient entity.Patient) error {
	// Check if required repositories are initialized
	if u.testTypeRepo == nil {
		return fmt.Errorf("test type repository is not initialized - this may occur when using SIMRS API mode without database configuration")
	}
	if u.workOrderRepo == nil {
		return fmt.Errorf("work order repository is not initialized - this may occur when using SIMRS API mode without database configuration")
	}
	if u.workOrderUC == nil {
		return fmt.Errorf("work order usecase is not initialized - this may occur when using SIMRS API mode without database configuration")
	}

	// Get test types by ID
	var testTypes []entity.WorkOrderCreateRequestTestType
	var notFoundIDs []int

	for _, testTypeID := range req.Order.OBR.OrderTest {
		// Find test type by ID
		testType, err := u.testTypeRepo.FindOneByID(ctx, testTypeID)
		if err != nil {
			slog.WarnContext(ctx, "test type not found", "id", testTypeID)
			notFoundIDs = append(notFoundIDs, testTypeID)
			continue
		}

		testTypes = append(testTypes, entity.WorkOrderCreateRequestTestType{
			TestTypeID:   int64(testType.ID),
			TestTypeCode: testType.Code,
			SpecimenType: testType.GetFirstType(),
		})
	}

	if len(testTypes) == 0 {
		return fmt.Errorf("no valid test types found for IDs: %v", req.Order.OBR.OrderTest)
	}

	if len(notFoundIDs) > 0 {
		slog.WarnContext(ctx, "some test type IDs not found", "ids", notFoundIDs)
	}

	// Check if work order with same SIMRS barcode already exists
	existingWorkOrder, err := u.workOrderRepo.GetBySIMRSBarcode(ctx, req.Order.OBR.OrderLab)
	if err == nil {
		// Work order exists, update with new test types
		slog.InfoContext(ctx, "Work order with SIMRS barcode exists, updating with new tests",
			"order_id", existingWorkOrder.ID,
			"simrs_barcode", req.Order.OBR.OrderLab,
			"existing_tests", len(existingWorkOrder.Specimen),
			"new_tests", len(testTypes))

		// Get existing test type IDs
		existingTestTypeIDs := make(map[int64]bool)
		for _, specimen := range existingWorkOrder.Specimen {
			for _, obsReq := range specimen.ObservationRequest {
				existingTestTypeIDs[int64(obsReq.TestType.ID)] = true
			}
		}

		// Filter out test types that already exist
		var newTestTypes []entity.WorkOrderCreateRequestTestType
		for _, testType := range testTypes {
			if !existingTestTypeIDs[testType.TestTypeID] {
				newTestTypes = append(newTestTypes, testType)
			}
		}

		if len(newTestTypes) == 0 {
			slog.InfoContext(ctx, "No new test types to add, all tests already exist in work order",
				"order_id", existingWorkOrder.ID,
				"simrs_barcode", req.Order.OBR.OrderLab)
			return nil
		}

		// Merge existing and new test types
		allTestTypes := make([]entity.WorkOrderCreateRequestTestType, 0)
		for _, specimen := range existingWorkOrder.Specimen {
			for _, obsReq := range specimen.ObservationRequest {
				allTestTypes = append(allTestTypes, entity.WorkOrderCreateRequestTestType{
					TestTypeID:   int64(obsReq.TestType.ID),
					TestTypeCode: obsReq.TestType.Code,
					SpecimenType: specimen.Type,
				})
			}
		}
		allTestTypes = append(allTestTypes, newTestTypes...)

		// Merge doctor and analyzer IDs
		mergedDoctorIDs := append(existingWorkOrder.DoctorIDs, req.Order.OBR.Doctor...)
		mergedAnalyzerIDs := append(existingWorkOrder.AnalyzerIDs, req.Order.OBR.Analyst...)

		// Remove duplicates
		mergedDoctorIDs = uniqueInt64(mergedDoctorIDs)
		mergedAnalyzerIDs = uniqueInt64(mergedAnalyzerIDs)

		workOrderReq := &entity.WorkOrderCreateRequest{
			PatientID:           patient.ID,
			TestTypes:           allTestTypes,
			CreatedBy:           int64(constant.CreatedBySystem),
			BarcodeSIMRS:        req.Order.OBR.OrderLab,
			MedicalRecordNumber: req.Order.PID.MedicalRecordNumber,
			Barcode:             existingWorkOrder.Barcode,
			DoctorIDs:           mergedDoctorIDs,
			AnalyzerIDs:         mergedAnalyzerIDs,
		}

		_, err := u.workOrderRepo.Edit(int(existingWorkOrder.ID), workOrderReq)
		if err != nil {
			return fmt.Errorf("failed to update work order: %w", err)
		}

		slog.InfoContext(ctx, "Work order updated with new tests from API",
			"order_id", existingWorkOrder.ID,
			"simrs_barcode", req.Order.OBR.OrderLab,
			"added_test_count", len(newTestTypes),
			"total_test_count", len(allTestTypes),
			"patient_id", patient.ID)
		return nil
	}

	// Work order doesn't exist, create new one
	if err != entity.ErrNotFound {
		return fmt.Errorf("error checking existing work order: %w", err)
	}

	workOrderReq := &entity.WorkOrderCreateRequest{
		PatientID:           patient.ID,
		TestTypes:           testTypes,
		CreatedBy:           int64(constant.CreatedBySystem), // System created from API
		BarcodeSIMRS:        req.Order.OBR.OrderLab,
		MedicalRecordNumber: req.Order.PID.MedicalRecordNumber,
		Barcode:             "",                    // Will be generated by workOrderUC.Create
		DoctorIDs:           req.Order.OBR.Doctor,  // Doctor IDs from request
		AnalyzerIDs:         req.Order.OBR.Analyst, // Analyzer IDs from request
	}

	_, err = u.workOrderUC.Create(workOrderReq)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	slog.InfoContext(ctx, "Work order created from API",
		"simrs_barcode", req.Order.OBR.OrderLab,
		"test_count", len(testTypes),
		"patient_id", patient.ID,
		"doctor_count", len(req.Order.OBR.Doctor),
		"analyzer_count", len(req.Order.OBR.Analyst))
	return nil
}

// uniqueInt64 removes duplicate int64 values from a slice
func uniqueInt64(slice []int64) []int64 {
	keys := make(map[int64]bool)
	result := []int64{}
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	return result
}

// splitName splits a full name into first and last name
func splitName(fullName string) []string {
	// Simple split by space
	// Can be enhanced with more sophisticated logic if needed
	parts := []string{}
	if fullName != "" {
		// Split by space and filter empty strings
		for _, part := range []rune(fullName) {
			if part == ' ' {
				continue
			}
		}
		// Simple approach: split by first space
		spaceIdx := -1
		for i, ch := range fullName {
			if ch == ' ' {
				spaceIdx = i
				break
			}
		}
		if spaceIdx > 0 {
			parts = append(parts, fullName[:spaceIdx])
			if spaceIdx+1 < len(fullName) {
				parts = append(parts, fullName[spaceIdx+1:])
			}
		} else {
			parts = append(parts, fullName)
		}
	}
	return parts
}

// GetResult retrieves lab result by order lab number
func (u *Usecase) GetResult(ctx context.Context, orderLab string) (entity.SimrsResultResponse, error) {
	slog.InfoContext(ctx, "debug simrs get result", "orderLab", orderLab)

	// Check if work order repository is initialized
	if u.workOrderRepo == nil {
		return entity.SimrsResultResponse{}, fmt.Errorf("work order repository is not initialized - this may occur when using SIMRS API mode without database configuration")
	}

	// Get work order by SIMRS barcode
	workOrder, err := u.workOrderRepo.GetBySIMRSBarcode(ctx, orderLab)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return entity.SimrsResultResponse{}, fmt.Errorf("order not found")
		}
		return entity.SimrsResultResponse{}, fmt.Errorf("error getting work order: %w", err)
	}

	// Fill test result details
	workOrder.FillTestResultDetail(false)

	// Build response from test results
	resultTests := make([]entity.SimrsResultTest, 0, len(workOrder.TestResult))
	for _, testResult := range workOrder.TestResult {
		// Get result value
		hasil := testResult.GetResult()

		// Parse and format numeric results
		if hasil != "" && testResult.TestType.Decimal >= 0 {
			if resultFloat, err := strconv.ParseFloat(hasil, 64); err == nil {
				hasil = strconv.FormatFloat(resultFloat, 'f', testResult.TestType.Decimal, 64)
			}
		}

		// Determine flag
		flag := ""
		khanzaFlag := entity.NewKhanzaFlag(testResult)
		switch khanzaFlag {
		case entity.FlagHigh, entity.FlagHighHigh:
			flag = "H"
		case entity.FlagLow, entity.FlagLowLow:
			flag = "L"
		default:
			flag = "N"
		}

		resultTests = append(resultTests, entity.SimrsResultTest{
			Loinc:       testResult.TestType.LoincCode,
			TestID:      testResult.TestType.ID,
			NamaTest:    testResult.TestType.Name,
			Hasil:       hasil,
			NilaiNormal: testResult.ReferenceRange,
			Satuan:      testResult.Unit,
			Flag:        flag,
		})
	}

	response := entity.SimrsResultResponse{
		Response: entity.SimrsResultResponseData{
			Sample: entity.SimrsResultSample{
				ResultTest: resultTests,
			},
		},
		Result: entity.SimrsResultOBX{
			OBX: entity.SimrsResultOBXData{
				OrderLab: orderLab,
			},
		},
	}

	slog.InfoContext(ctx, "debug simrs get result", "response", response)
	return response, nil
}

// DeleteOrder deletes a work order by SIMRS barcode
func (u *Usecase) DeleteOrder(ctx context.Context, orderLab string) error {
	slog.InfoContext(ctx, "deleting order by SIMRS barcode", "orderLab", orderLab)

	// Check if work order repository is initialized
	if u.workOrderRepo == nil {
		return fmt.Errorf("work order repository is not initialized - this may occur when using SIMRS API mode without database configuration")
	}

	// Get work order by SIMRS barcode
	workOrder, err := u.workOrderRepo.GetBySIMRSBarcode(ctx, orderLab)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return fmt.Errorf("order not found")
		}
		return fmt.Errorf("error getting work order: %w", err)
	}

	// Delete the work order
	err = u.workOrderRepo.Delete(int64(workOrder.ID))
	if err != nil {
		return fmt.Errorf("error deleting work order: %w", err)
	}

	slog.InfoContext(ctx, "order deleted successfully", "orderLab", orderLab, "workOrderID", workOrder.ID)
	return nil
}
