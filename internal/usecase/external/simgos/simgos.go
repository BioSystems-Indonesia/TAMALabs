package simgosuc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	simgosrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/external/simgos"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	testType "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrder "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/result"
	workOrderUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order"
)

type Usecase struct {
	simgosRepo    *simgosrepo.Repository
	workOrderRepo *workOrder.WorkOrderRepository
	workOrderUC   *workOrderUC.WorkOrderUseCase
	patientRepo   *patientrepo.PatientRepository
	testTypeRepo  *testType.Repository
	cfg           *config.Schema
	resultUC      *result.Usecase
	lastDSN       string // Track last used DSN to detect changes
}

func NewUsecase(simgosRepo *simgosrepo.Repository, workOrderRepo *workOrder.WorkOrderRepository,
	workOrderUC *workOrderUC.WorkOrderUseCase, patientRepo *patientrepo.PatientRepository,
	testTypeRepo *testType.Repository, cfg *config.Schema, resultUC *result.Usecase) *Usecase {
	lastDSN := ""
	if simgosRepo != nil {
		lastDSN = cfg.SimgosDatabaseDSN
	}
	return &Usecase{
		simgosRepo:    simgosRepo,
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
	currentDSN := u.cfg.SimgosDatabaseDSN

	// If DSN changed or repository is nil, reinitialize
	if u.simgosRepo == nil || u.lastDSN != currentDSN {
		if currentDSN == "" {
			u.simgosRepo = nil
			u.lastDSN = ""
			return fmt.Errorf("Database Sharing database DSN is not configured")
		}

		slog.Info("Reinitializing Database Sharing repository", "dsn_changed", u.lastDSN != currentDSN, "repo_was_nil", u.simgosRepo == nil)

		db, err := simgosrepo.NewDB(currentDSN)
		if err != nil {
			return fmt.Errorf("failed to create Database Sharing database connection: %w", err)
		}

		u.simgosRepo = simgosrepo.NewRepository(db)
		u.lastDSN = currentDSN
		slog.Info("Database Sharing repository reinitialized successfully")
	}

	return nil
}

// SyncAllRequest fetches NEW orders from Database Sharing and creates work orders in LIS
func (u *Usecase) SyncAllRequest(ctx context.Context) error {
	// Ensure repository is initialized and up-to-date
	if err := u.ensureRepository(); err != nil {
		return fmt.Errorf("Database Sharing repository initialization failed: %w", err)
	}

	slog.Info("Starting Database Sharing lab request synchronization")

	// Fetch all NEW orders from Database Sharing
	newOrders, err := u.simgosRepo.GetNewLabOrders(ctx)
	if err != nil {
		return fmt.Errorf("failed to get new lab orders from Database Sharing: %w", err)
	}

	slog.Info(fmt.Sprintf("Found %d new lab orders from Database Sharing", len(newOrders)))

	var processedCount int
	var errorCount int

	for _, order := range newOrders {
		err := u.processNewOrder(ctx, order)
		if err != nil {
			slog.Error("Failed to process lab order", "no_lab_order", order.NoLabOrder, "error", err)
			errorCount++
			continue
		}

		// Update order status to PENDING after successful processing
		err = u.simgosRepo.UpdateOrderStatus(ctx, order.NoLabOrder, entity.SimgosStatusPending)
		if err != nil {
			slog.Error("Failed to update order status to PENDING", "no_lab_order", order.NoLabOrder, "error", err)
			errorCount++
			continue
		}

		processedCount++
		slog.Info("Successfully processed and updated order status",
			"no_lab_order", order.NoLabOrder,
			"status", entity.SimgosStatusPending)
	}

	slog.Info("Database Sharing lab request synchronization completed",
		"total", len(newOrders),
		"processed", processedCount,
		"errors", errorCount)

	return nil
}

// processNewOrder processes a single new order from Database Sharing
func (u *Usecase) processNewOrder(ctx context.Context, order entity.SimgosLabOrder) error {
	// Get order details
	orderDetails, err := u.simgosRepo.GetOrderDetailsByNoLabOrder(ctx, order.NoLabOrder)
	if err != nil {
		return fmt.Errorf("failed to get order details: %w", err)
	}

	if len(orderDetails) == 0 {
		return fmt.Errorf("no order details found for no_lab_order: %s", order.NoLabOrder)
	}

	// Find or create patient
	patient, err := u.findOrCreatePatient(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to find/create patient: %w", err)
	}

	// Create work order
	err = u.createWorkOrder(ctx, order, orderDetails, patient.ID)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	slog.Info("Successfully processed new order",
		"no_lab_order", order.NoLabOrder,
		"patient_id", patient.ID,
		"test_count", len(orderDetails))

	return nil
}

// findOrCreatePatient finds an existing patient or creates a new one from Database Sharing order data
func (u *Usecase) findOrCreatePatient(ctx context.Context, order entity.SimgosLabOrder) (*entity.Patient, error) {
	// Split patient name into first and last name
	var firstName, lastName string
	nameParts := strings.Fields(order.PatientName)
	if len(nameParts) > 0 {
		firstName = nameParts[0]
	}
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}

	// Convert sex
	sex := entity.SimgosSex(order.Sex).ToPatientSex()

	// Try to find existing patient by Medical Record Number (No RM)
	if order.NoRM != "" {
		var existingPatient entity.Patient
		err := u.patientRepo.FindByMedicalRecordNumber(ctx, order.NoRM, &existingPatient)
		if err == nil {
			// Patient found, update information
			slog.Info("Patient found by Medical Record Number",
				"patient_id", existingPatient.ID,
				"no_rm", order.NoRM)

			existingPatient.FirstName = firstName
			existingPatient.LastName = lastName
			existingPatient.Sex = sex
			existingPatient.Birthdate = order.BirthDate
			existingPatient.UpdatedAt = time.Now()

			err = u.patientRepo.Update(&existingPatient)
			if err != nil {
				return nil, fmt.Errorf("error updating patient: %w", err)
			}

			return &existingPatient, nil
		}
	}

	// Create new patient
	patient := entity.Patient{
		MedicalRecordNumber: order.NoRM,
		FirstName:           firstName,
		LastName:            lastName,
		Sex:                 sex,
		Birthdate:           order.BirthDate,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	result, err := u.patientRepo.FirstOrCreate(&patient)
	if err != nil {
		return nil, fmt.Errorf("error creating patient: %w", err)
	}

	slog.Info("Patient created from Database Sharing order",
		"patient_id", result.ID,
		"name", order.PatientName,
		"no_rm", order.NoRM)

	return &result, nil
}

// createWorkOrder creates a work order from Database Sharing order and details
func (u *Usecase) createWorkOrder(ctx context.Context, order entity.SimgosLabOrder, orderDetails []entity.SimgosOrderDetail, patientID int64) error {
	// Map parameter codes to test types
	var testTypes []entity.WorkOrderCreateRequestTestType
	var notFoundCodes []string

	for _, detail := range orderDetails {
		// Find test type by code
		testType, err := u.testTypeRepo.FindOneByCode(ctx, detail.ParameterCode)
		if err != nil {
			slog.Warn("Test type not found", "parameter_code", detail.ParameterCode)
			notFoundCodes = append(notFoundCodes, detail.ParameterCode)
			continue
		}

		testTypes = append(testTypes, entity.WorkOrderCreateRequestTestType{
			TestTypeID:   int64(testType.ID),
			TestTypeCode: testType.Code,
			SpecimenType: testType.GetFirstType(),
		})
	}

	if len(testTypes) == 0 {
		return fmt.Errorf("no valid test types found for parameter codes: %v", notFoundCodes)
	}

	if len(notFoundCodes) > 0 {
		slog.Warn("Some parameter codes not found in test types", "codes", notFoundCodes)
	}

	// Check if work order with same Database Sharing barcode already exists
	existingWorkOrder, err := u.workOrderRepo.GetBySIMRSBarcode(ctx, order.NoLabOrder)
	if err == nil {
		// Work order exists, log and skip
		slog.Info("Work order with Database Sharing barcode already exists, skipping creation",
			"order_id", existingWorkOrder.ID,
			"no_lab_order", order.NoLabOrder)
		return nil
	}

	// Work order doesn't exist, create new one
	if err != entity.ErrNotFound {
		return fmt.Errorf("error checking existing work order: %w", err)
	}

	workOrderReq := &entity.WorkOrderCreateRequest{
		PatientID:           patientID,
		TestTypes:           testTypes,
		CreatedBy:           1, // System created
		BarcodeSIMRS:        order.NoLabOrder,
		MedicalRecordNumber: order.NoRM,
		Barcode:             "", // Will be generated by workOrderUC.Create
	}

	_, err = u.workOrderUC.Create(workOrderReq)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	slog.Info("Work order created from Database Sharing",
		"no_lab_order", order.NoLabOrder,
		"test_count", len(testTypes),
		"patient_id", patientID)

	return nil
}

// SyncAllResult syncs completed work order results back to Database Sharing
func (u *Usecase) SyncAllResult(ctx context.Context, workOrderIDs []int64) error {
	// Ensure repository is initialized and up-to-date
	if err := u.ensureRepository(); err != nil {
		return fmt.Errorf("Database Sharing repository initialization failed: %w", err)
	}

	slog.Info("Starting Database Sharing result synchronization")

	// Get work orders from today only
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	workOrders, err := u.workOrderRepo.FindAll(ctx, &entity.WorkOrderGetManyRequest{
		GetManyRequest: entity.GetManyRequest{
			CreatedAtStart: startOfDay,
			CreatedAtEnd:   endOfDay,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to get work orders: %w", err)
	}

	slog.Info(fmt.Sprintf("Found %d work orders to sync (today only)", len(workOrders.Data)))

	var syncedCount int
	var errorCount int

	for _, workOrder := range workOrders.Data {
		// Skip if no SIMRS barcode (not from Database Sharing)
		if workOrder.BarcodeSIMRS == "" {
			continue
		}

		err := u.syncResultByWorkOrder(ctx, workOrder.ID)
		if err != nil {
			slog.Error("Failed to sync result", "work_order_id", workOrder.ID, "error", err)
			errorCount++
			continue
		}

		syncedCount++
	}

	slog.Info("Database Sharing result synchronization completed",
		"total_checked", len(workOrders.Data),
		"synced", syncedCount,
		"errors", errorCount)

	return nil
}

// syncResultByWorkOrder syncs results for a single work order to Database Sharing
func (u *Usecase) syncResultByWorkOrder(ctx context.Context, workOrderID int64) error {
	// Get work order with results
	workOrder, err := u.resultUC.ResultDetail(ctx, workOrderID)
	if err != nil {
		return fmt.Errorf("failed to get work order results: %w", err)
	}

	// Skip if no SIMRS barcode
	if workOrder.BarcodeSIMRS == "" {
		slog.Debug("Work order does not have Database Sharing barcode", "order_id", workOrderID)
		return nil
	}

	// Fill test result details
	workOrder.FillTestResultDetail(false)

	if len(workOrder.TestResult) == 0 {
		slog.Debug("No test results found", "order_id", workOrderID)
		return nil
	}

	// Check if order exists in Database Sharing
	exists, err := u.simgosRepo.CheckOrderExists(ctx, workOrder.BarcodeSIMRS)
	if err != nil {
		return fmt.Errorf("failed to check order existence: %w", err)
	}

	if !exists {
		slog.Debug("Order not found in Database Sharing, skipping", "no_lab_order", workOrder.BarcodeSIMRS)
		return nil
	}

	// Get current order status
	simgosOrder, err := u.simgosRepo.GetLabOrderByNoLabOrder(ctx, workOrder.BarcodeSIMRS)
	if err != nil {
		return fmt.Errorf("failed to get Database Sharing order: %w", err)
	}

	// Only sync if status is PENDING (order has been fetched by LIS)
	if simgosOrder.Status != string(entity.SimgosStatusPending) {
		slog.Debug("Order status is not PENDING, skipping",
			"no_lab_order", workOrder.BarcodeSIMRS,
			"status", simgosOrder.Status)
		return nil
	}

	// Prepare order detail updates
	var orderDetailUpdates []entity.SimgosOrderDetail

	for _, testResultGroup := range workOrder.TestResult {
		for _, testResult := range testResultGroup {
			// Skip if no result
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

			// Create order detail update
			orderDetail := entity.SimgosOrderDetail{
				NoLabOrder:     workOrder.BarcodeSIMRS,
				ParameterCode:  testType.Code,
				ResultValue:    testResult.Result,
				Unit:           testResult.Unit,
				ReferenceRange: testResult.ReferenceRange,
				Flag:           string(entity.NewSimgosFlag(testResult)),
			}

			orderDetailUpdates = append(orderDetailUpdates, orderDetail)
		}
	}

	if len(orderDetailUpdates) == 0 {
		slog.Debug("No valid results to sync", "order_id", workOrderID)
		return nil
	}

	// Update order details in Database Sharing
	err = u.simgosRepo.BatchUpdateOrderDetails(ctx, orderDetailUpdates)
	if err != nil {
		return fmt.Errorf("failed to update order details: %w", err)
	}

	// Check if all order details are completed
	totalDetails, err := u.simgosRepo.CountOrderDetailsByNoLabOrder(ctx, workOrder.BarcodeSIMRS)
	if err != nil {
		return fmt.Errorf("failed to count order details: %w", err)
	}

	completedDetails, err := u.simgosRepo.CountCompletedOrderDetailsByNoLabOrder(ctx, workOrder.BarcodeSIMRS)
	if err != nil {
		return fmt.Errorf("failed to count completed order details: %w", err)
	}

	// If all details are completed, update order status to LIS_SUCCESS
	if completedDetails >= totalDetails {
		err = u.simgosRepo.UpdateOrderStatus(ctx, workOrder.BarcodeSIMRS, entity.SimgosStatusLISSuccess)
		if err != nil {
			return fmt.Errorf("failed to update order status to LIS_SUCCESS: %w", err)
		}

		slog.Info("All order details completed, updated status to LIS_SUCCESS",
			"no_lab_order", workOrder.BarcodeSIMRS,
			"total_details", totalDetails,
			"completed_details", completedDetails)
	} else {
		slog.Info("Order details partially completed",
			"no_lab_order", workOrder.BarcodeSIMRS,
			"total_details", totalDetails,
			"completed_details", completedDetails)
	}

	slog.Info("Successfully synced lab results to Database Sharing",
		"order_id", workOrderID,
		"no_lab_order", workOrder.BarcodeSIMRS,
		"results_count", len(orderDetailUpdates))

	return nil
}

// TestConnection tests the Database Sharing database connection
func (u *Usecase) TestConnection(ctx context.Context) error {
	if u.simgosRepo == nil {
		return fmt.Errorf("Database Sharing repository is not initialized")
	}

	return u.simgosRepo.TestConnection(ctx)
}
