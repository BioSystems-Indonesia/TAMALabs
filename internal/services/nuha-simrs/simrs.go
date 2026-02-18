package nuha_simrs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	testTypeRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrder "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	workOrderUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order"
)

type SIMRSNuha struct {
	BaseURL       string
	SessionID     string
	workOrderRepo *workOrder.WorkOrderRepository
	workOrderUC   *workOrderUC.WorkOrderUseCase
	patientRepo   *patientrepo.PatientRepository
	testTypeRepo  *testTypeRepo.Repository
}

func NewSIMRSNuha(
	baseURL string,
	sessionID string,
	workOrderRepo *workOrder.WorkOrderRepository,
	workOrderUC *workOrderUC.WorkOrderUseCase,
	patientrepo *patientrepo.PatientRepository,
	testTypeRepo *testTypeRepo.Repository,
) *SIMRSNuha {
	return &SIMRSNuha{
		BaseURL:       baseURL,
		SessionID:     sessionID,
		workOrderRepo: workOrderRepo,
		workOrderUC:   workOrderUC,
		patientRepo:   patientrepo,
		testTypeRepo:  testTypeRepo,
	}
}

func (c *SIMRSNuha) GetLabOrder(ctx context.Context) error {
	nuhaService := NewNuhaService(c.BaseURL)

	// Get today's date in YYYY-MM-DD format
	today := time.Now().Format("2006-01-02")

	req := LabListRequest{
		SessionID: c.SessionID,
		ValidFrom: today,
		ValidTo:   today,
	}

	// Mask session ID for logs (avoid leaking full session token)
	maskedSession := c.SessionID
	if len(maskedSession) > 6 {
		maskedSession = maskedSession[:6] + "..."
	}

	slog.InfoContext(ctx, "Starting GetLabOrder",
		"base_url", c.BaseURL,
		"session_id_masked", maskedSession,
		"valid_from", req.ValidFrom,
		"valid_to", req.ValidTo,
	)

	// Log masked request payload for easier debugging (session id masked above)
	logReq := req
	logReq.SessionID = maskedSession
	if payloadBytes, err := json.Marshal(logReq); err == nil {
		slog.InfoContext(ctx, "Nuha GetLabList request payload (masked)", "request_payload", string(payloadBytes))
	} else {
		slog.WarnContext(ctx, "Failed to marshal Nuha GetLabList request payload for logging", "error", err)
	}

	start := time.Now()
	result, err := nuhaService.GetLabList(ctx, req)
	durationMs := time.Since(start).Milliseconds()
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get lab list from Nuha SIMRS",
			"error", err,
			"duration_ms", durationMs,
			"session_id_masked", maskedSession,
		)
		return fmt.Errorf("failed to get lab list: %w", err)
	}

	count := 0
	if result != nil {
		count = len(result.Response.List)
	}

	slog.InfoContext(ctx, "Retrieved lab orders from Nuha SIMRS",
		"count", count,
		"date", today,
		"duration_ms", durationMs,
	)

	// Log response metadata if present
	if result != nil {
		slog.InfoContext(ctx, "Nuha response metadata",
			"metadata_code", result.Metadata.Code,
			"metadata_message", result.Metadata.Message,
		)
	}

	if count == 0 {
		slog.WarnContext(ctx, "No lab orders returned from Nuha SIMRS", "date", today)
		return nil
	}

	// Process entries and collect stats
	var processedCount, failedCount int
	for idx, labReg := range result.Response.List {
		// Build concise test summary (limit to first 8 tests)
		testSummaries := make([]string, 0, len(labReg.TestList))
		for i, t := range labReg.TestList {
			if i >= 8 {
				break
			}
			testSummaries = append(testSummaries, fmt.Sprintf("%s(%d)", t.TestName, t.TestID))
		}

		slog.InfoContext(ctx, "Processing lab registration from Nuha SIMRS",
			"index", idx,
			"lab_number", labReg.LabNumber,
			"order_date", labReg.OrderDate.Format("2006-01-02 15:04:05"),
			"patient_name", labReg.PatientName,
			"mrn", labReg.MedicalRecordNo,
			"birthdate", labReg.BirthDate.Format("2006-01-02"),
			"gender", labReg.Gender,
			"age", labReg.AgeDescription,
			"room", labReg.Room,
			"is_cito", labReg.IsCITO,
			"test_count", len(labReg.TestList),
			"tests_sample", strings.Join(testSummaries, ", "),
		)

		if err := c.processLabRegistration(ctx, labReg); err != nil {
			failedCount++
			slog.ErrorContext(ctx, "Failed to process lab registration",
				"lab_number", labReg.LabNumber,
				"patient_name", labReg.PatientName,
				"mrn", labReg.MedicalRecordNo,
				"error", err,
			)
			continue
		}

		processedCount++
	}

	slog.InfoContext(ctx, "GetLabOrder processing summary",
		"total_retrieved", count,
		"processed_success", processedCount,
		"processed_failed", failedCount,
		"date", today,
	)

	return nil
}

func (c *SIMRSNuha) processLabRegistration(ctx context.Context, labReg LabRegistration) error {
	patient, err := c.findOrCreatePatient(ctx, labReg)
	if err != nil {
		return fmt.Errorf("failed to find/create patient: %w", err)
	}

	slog.InfoContext(ctx, "Patient processed",
		"patient_id", patient.ID,
		"name", patient.FirstName+" "+patient.LastName,
		"mr_number", labReg.MedicalRecordNo)

	testTypes, err := c.parseTestTypes(ctx, labReg.TestList)
	if err != nil {
		return fmt.Errorf("failed to parse test types: %w", err)
	}

	if len(testTypes) == 0 {
		slog.WarnContext(ctx, "No valid test types found for lab registration", "lab_number", labReg.LabNumber)
		return nil
	}

	labNumberStr := fmt.Sprintf("NUHA-%d", labReg.LabNumber)
	existingWorkOrder, err := c.workOrderRepo.GetBySIMRSBarcode(ctx, labNumberStr)
	if err == nil {
		// Work order exists, update it
		slog.InfoContext(ctx, "Work order already exists, updating",
			"work_order_id", existingWorkOrder.ID,
			"lab_number", labNumberStr)
		return c.updateWorkOrder(ctx, existingWorkOrder, patient, testTypes, labReg)
	}

	workOrderReq := &entity.WorkOrderCreateRequest{
		PatientID:           patient.ID,
		TestTypes:           testTypes,
		CreatedBy:           int64(constant.CreatedBySystem),
		BarcodeSIMRS:        labNumberStr,
		MedicalRecordNumber: labReg.MedicalRecordNo,
		Barcode:             "", // Will be generated by system
		DoctorIDs:           []int64{},
		AnalyzerIDs:         []int64{},
	}

	workOrderID, err := c.workOrderUC.Create(workOrderReq)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	slog.InfoContext(ctx, "Work order created successfully",
		"work_order_id", workOrderID,
		"lab_number", labNumberStr,
		"patient_id", patient.ID,
		"test_count", len(testTypes))

	return nil
}

func (c *SIMRSNuha) findOrCreatePatient(ctx context.Context, labReg LabRegistration) (entity.Patient, error) {
	if labReg.MedicalRecordNo != "" {
		var existingPatient entity.Patient
		err := c.patientRepo.FindByMedicalRecordNumber(ctx, labReg.MedicalRecordNo, &existingPatient)
		if err == nil {
			return existingPatient, nil
		}
	}

	firstName, lastName := parsePatientName(labReg.PatientName)

	sex := entity.PatientSexUnknown
	switch strings.ToUpper(labReg.Gender) {
	case "L", "LAKI-LAKI", "M", "MALE":
		sex = entity.PatientSexMale
	case "P", "PEREMPUAN", "F", "FEMALE":
		sex = entity.PatientSexFemale
	}

	patient := &entity.Patient{
		FirstName:           firstName,
		LastName:            lastName,
		Birthdate:           labReg.BirthDate.Time,
		Sex:                 sex,
		MedicalRecordNumber: labReg.MedicalRecordNo,
		Address:             labReg.Address,
		PhoneNumber:         "",
		Location:            labReg.Room,
	}

	result, err := c.patientRepo.FirstOrCreate(patient)
	if err != nil {
		return entity.Patient{}, fmt.Errorf("failed to create patient: %w", err)
	}

	if result.MedicalRecordNumber == "" && labReg.MedicalRecordNo != "" {
		result.MedicalRecordNumber = labReg.MedicalRecordNo
		if err := c.patientRepo.Update(&result); err != nil {
			slog.WarnContext(ctx, "Failed to update patient medical record number", "error", err)
		}
	}

	return result, nil
}

func (c *SIMRSNuha) parseTestTypes(ctx context.Context, testList []LabTest) ([]entity.WorkOrderCreateRequestTestType, error) {
	var testTypes []entity.WorkOrderCreateRequestTestType
	seenTests := make(map[int]bool)

	for _, labTest := range testList {
		if labTest.TestType == "p" && len(labTest.TestDetails) > 0 {
			slog.InfoContext(ctx, "Processing package test",
				"package_name", labTest.TestName,
				"package_id", labTest.TestID,
				"detail_count", len(labTest.TestDetails))

			for _, detail := range labTest.TestDetails {
				if seenTests[detail.TestID] {
					continue
				}
				seenTests[detail.TestID] = true

				testIDStr := fmt.Sprintf("%d", detail.TestID)
				testType, err := c.testTypeRepo.FindOneByAliasCode(ctx, testIDStr)
				if err != nil {
					testType, err = c.testTypeRepo.FindOneByCode(ctx, detail.TestName)
					if err != nil {
						slog.WarnContext(ctx, "Test type not found in system (from package)",
							"package_name", labTest.TestName,
							"test_name", detail.TestName,
							"test_id", detail.TestID,
							"searched_alias_code", testIDStr)
						continue
					}
					slog.InfoContext(ctx, "Test type found by code/name (from package)",
						"package_name", labTest.TestName,
						"test_name", detail.TestName,
						"test_id", detail.TestID,
						"test_type_id", testType.ID)
				} else {
					slog.InfoContext(ctx, "Test type found by alias code (from package)",
						"package_name", labTest.TestName,
						"test_name", detail.TestName,
						"test_id", detail.TestID,
						"test_type_id", testType.ID,
						"alias_code", testType.AliasCode)
				}

				packageID := detail.PackageID
				testTypes = append(testTypes, entity.WorkOrderCreateRequestTestType{
					TestTypeID:   int64(testType.ID),
					TestTypeCode: testType.Code,
					SpecimenType: testType.GetFirstType(),
					PackageID:    &packageID,
				})
			}
		} else {
			if seenTests[labTest.TestID] {
				continue
			}
			seenTests[labTest.TestID] = true

			testIDStr := fmt.Sprintf("%d", labTest.TestID)

			testType, err := c.testTypeRepo.FindOneByAliasCode(ctx, testIDStr)
			if err != nil {
				testType, err = c.testTypeRepo.FindOneByCode(ctx, labTest.TestName)
				if err != nil {
					slog.WarnContext(ctx, "Test type not found in system",
						"test_name", labTest.TestName,
						"test_id", labTest.TestID,
						"searched_alias_code", testIDStr,
						"searched_by", "alias_code (TestID) and code (TestName)")
					continue
				}
				slog.InfoContext(ctx, "Test type found by code/name (not by alias code)",
					"test_name", labTest.TestName,
					"test_id", labTest.TestID,
					"test_type_id", testType.ID,
					"test_type_code", testType.Code)
			} else {
				slog.InfoContext(ctx, "Test type found by alias code (TestID)",
					"test_name", labTest.TestName,
					"test_id", labTest.TestID,
					"test_type_id", testType.ID,
					"test_type_code", testType.Code,
					"alias_code", testType.AliasCode)
			}

			testTypes = append(testTypes, entity.WorkOrderCreateRequestTestType{
				TestTypeID:   int64(testType.ID),
				TestTypeCode: testType.Code,
				SpecimenType: testType.GetFirstType(),
				PackageID:    nil,
			})
		}
	}

	return testTypes, nil
}

func (c *SIMRSNuha) updateWorkOrder(
	ctx context.Context,
	existingWorkOrder entity.WorkOrder,
	patient entity.Patient,
	newTestTypes []entity.WorkOrderCreateRequestTestType,
	labReg LabRegistration,
) error {
	existingTestTypeIDs := make(map[int64]bool)
	for _, specimen := range existingWorkOrder.Specimen {
		for _, obsReq := range specimen.ObservationRequest {
			existingTestTypeIDs[int64(obsReq.TestType.ID)] = true
		}
	}

	var testsToAdd []entity.WorkOrderCreateRequestTestType
	for _, testType := range newTestTypes {
		if !existingTestTypeIDs[testType.TestTypeID] {
			testsToAdd = append(testsToAdd, testType)
		}
	}

	if len(testsToAdd) == 0 {
		slog.InfoContext(ctx, "No new tests to add to existing work order",
			"work_order_id", existingWorkOrder.ID)
		return nil
	}

	allTestTypes := make([]entity.WorkOrderCreateRequestTestType, 0)
	for _, specimen := range existingWorkOrder.Specimen {
		for _, obsReq := range specimen.ObservationRequest {
			allTestTypes = append(allTestTypes, entity.WorkOrderCreateRequestTestType{
				TestTypeID:   int64(obsReq.TestType.ID),
				TestTypeCode: obsReq.TestType.Code,
				SpecimenType: specimen.Type,
				PackageID:    obsReq.PackageID,
			})
		}
	}
	allTestTypes = append(allTestTypes, testsToAdd...)

	workOrderReq := &entity.WorkOrderCreateRequest{
		PatientID:           patient.ID,
		TestTypes:           allTestTypes,
		CreatedBy:           int64(constant.CreatedBySystem),
		BarcodeSIMRS:        fmt.Sprintf("NUHA-%d", labReg.LabNumber),
		MedicalRecordNumber: labReg.MedicalRecordNo,
		Barcode:             existingWorkOrder.Barcode,
		DoctorIDs:           existingWorkOrder.DoctorIDs,
		AnalyzerIDs:         existingWorkOrder.AnalyzerIDs,
	}

	_, err := c.workOrderRepo.Edit(int(existingWorkOrder.ID), workOrderReq)
	if err != nil {
		return fmt.Errorf("failed to update work order: %w", err)
	}

	slog.InfoContext(ctx, "Work order updated successfully",
		"work_order_id", existingWorkOrder.ID,
		"added_tests", len(testsToAdd),
		"total_tests", len(allTestTypes))

	return nil
}

// SendResultToNuha sends test result to Nuha SIMRS
func (c *SIMRSNuha) SendResultToNuha(
	ctx context.Context,
	labNumber int,
	testID int,
	testName string,
	resultValue string,
	unit string,
	referenceRange string,
	abnormalFlag string,
	resultText string,
	packageID int,
	index int,
	insertedUser string,
	insertedIP string,
) error {
	nuhaService := NewNuhaService(c.BaseURL)

	req := InsertResultRequest{
		SessionID:      c.SessionID,
		LabNumber:      labNumber,
		TestName:       testName,
		Result:         resultValue,
		Unit:           unit,
		ReferenceRange: referenceRange,
		Abnormal:       abnormalFlag,
		Description:    "",
		Notes:          "",
		TestID:         testID,
		ResultText:     resultText,
		PackageID:      packageID,
		Spacing:        "",
		Index:          index,
		InsertedUser:   insertedUser,
		InsertedIP:     insertedIP,
	}

	fmt.Println(req)

	result, err := nuhaService.InsertResult(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to insert result to Nuha SIMRS: %w", err)
	}

	slog.InfoContext(ctx, "Result sent to Nuha SIMRS successfully",
		"lab_number", labNumber,
		"test_id", testID,
		"test_name", testName,
		"status", result.Response.Status,
		"message", result.Response.Message)

	return nil
}

// SendWorkOrderResults sends all test results from a work order to Nuha SIMRS
func (c *SIMRSNuha) SendWorkOrderResults(
	ctx context.Context,
	barcodeSIMRS string,
	insertedUser string,
	insertedIP string,
) error {
	// Extract lab number from SIMRS barcode (NUHA-12345 -> 12345)
	labNumberStr := strings.TrimPrefix(barcodeSIMRS, "NUHA-")
	labNumber := 0
	if _, err := fmt.Sscanf(labNumberStr, "%d", &labNumber); err != nil || labNumber == 0 {
		return fmt.Errorf("invalid SIMRS barcode format: %s", barcodeSIMRS)
	}

	// Get work order by SIMRS barcode
	wo, err := c.workOrderRepo.GetBySIMRSBarcode(ctx, barcodeSIMRS)
	if err != nil {
		return fmt.Errorf("failed to get work order: %w", err)
	}

	return c.sendWorkOrderResults(ctx, wo, labNumber, insertedUser, insertedIP)
}

// SendWorkOrderResultsByID sends all test results from a work order to Nuha SIMRS by work order ID
func (c *SIMRSNuha) SendWorkOrderResultsByID(
	ctx context.Context,
	workOrderIDStr string,
	insertedUser string,
	insertedIP string,
) error {
	// Parse work order ID
	var workOrderID int64
	if _, err := fmt.Sscanf(workOrderIDStr, "%d", &workOrderID); err != nil {
		return fmt.Errorf("invalid work order ID: %s", workOrderIDStr)
	}

	// Get work order by ID
	wo, err := c.workOrderRepo.FindOneForResult(workOrderID)
	if err != nil {
		return fmt.Errorf("failed to get work order: %w", err)
	}

	// Check if barcode_simrs exists
	if wo.BarcodeSIMRS == "" {
		return fmt.Errorf("work order does not have SIMRS barcode, cannot send to Nuha")
	}

	// Extract lab number from SIMRS barcode (NUHA-12345 -> 12345)
	labNumberStr := strings.TrimPrefix(wo.BarcodeSIMRS, "NUHA-")
	labNumber := 0
	if _, err := fmt.Sscanf(labNumberStr, "%d", &labNumber); err != nil || labNumber == 0 {
		return fmt.Errorf("invalid SIMRS barcode format: %s", wo.BarcodeSIMRS)
	}

	return c.sendWorkOrderResults(ctx, wo, labNumber, insertedUser, insertedIP)
}

// sendWorkOrderResults is the actual implementation that sends results to Nuha
func (c *SIMRSNuha) sendWorkOrderResults(
	ctx context.Context,
	workOrder entity.WorkOrder,
	labNumber int,
	insertedUser string,
	insertedIP string,
) error {
	slog.InfoContext(ctx, "Sending work order results to Nuha SIMRS",
		"work_order_id", workOrder.ID,
		"barcode_simrs", workOrder.BarcodeSIMRS,
		"lab_number", labNumber,
		"inserted_user", insertedUser,
		"specimen_count", len(workOrder.Specimen))

	// Validate that work order has specimens
	if len(workOrder.Specimen) == 0 {
		return fmt.Errorf("work order has no specimens to send")
	}

	// Collect all test results into batch
	var batchItems []BatchInsertResultItem
	skippedCount := 0
	index := 1

	// Loop through all specimens and observation results
	for specIdx, specimen := range workOrder.Specimen {
		slog.InfoContext(ctx, "Processing specimen",
			"specimen_index", specIdx,
			"specimen_id", specimen.ID,
			"observation_result_count", len(specimen.ObservationResult))

		// Create a map from test_type_id to package_id from observation requests
		testTypeToPackageID := make(map[int]*int)
		for _, obsReq := range specimen.ObservationRequest {
			if obsReq.TestTypeID != nil {
				testTypeToPackageID[*obsReq.TestTypeID] = obsReq.PackageID
			}
		}

		for obsIdx, obsResult := range specimen.ObservationResult {
			// Skip if no result value
			if len(obsResult.Values) == 0 || obsResult.Values[0] == "" {
				skippedCount++
				slog.WarnContext(ctx, "Skipping observation result - no value",
					"observation_result_id", obsResult.ID,
					"test_name", obsResult.TestType.Name,
					"values_length", len(obsResult.Values))
				continue
			}

			slog.InfoContext(ctx, "Processing observation result",
				"observation_index", obsIdx,
				"observation_result_id", obsResult.ID,
				"test_name", obsResult.TestType.Name,
				"value", obsResult.Values[0])

			// Get the test type to map TestID (using alias code)
			testIDFromAlias := 0
			if obsResult.TestType.AliasCode != "" {
				if _, err := fmt.Sscanf(obsResult.TestType.AliasCode, "%d", &testIDFromAlias); err != nil {
					slog.WarnContext(ctx, "Failed to parse alias code as test ID",
						"alias_code", obsResult.TestType.AliasCode,
						"test_name", obsResult.TestType.Name)
				}
			}

			// Determine abnormal flag based on reference range and result value
			abnormalFlag := determineAbnormalFlag(obsResult.Values[0], obsResult.ReferenceRange)

			slog.InfoContext(ctx, "Determined abnormal flag",
				"test_name", obsResult.TestType.Name,
				"value", obsResult.Values[0],
				"reference_range", obsResult.ReferenceRange,
				"abnormal_flag", abnormalFlag)

			// Get result value (use first value)
			resultValue := obsResult.Values[0]

			// Create result text if comment exists
			resultText := ""
			if obsResult.Comments != "" {
				resultText = obsResult.Comments
			}

			// Get package ID from observation request mapping
			packageID := 0
			if pkgID, ok := testTypeToPackageID[obsResult.TestType.ID]; ok && pkgID != nil {
				packageID = *pkgID
			}

			// Add to batch
			batchItems = append(batchItems, BatchInsertResultItem{
				LabNumber:      labNumber,
				TestName:       obsResult.TestType.Name,
				Result:         resultValue,
				ReferenceRange: obsResult.ReferenceRange,
				Abnormal:       abnormalFlag,
				Unit:           obsResult.Unit,
				TestID:         testIDFromAlias,
				PackageID:      packageID,
				Index:          index,
				ResultText:     resultText,
				InsertedUser:   insertedUser,
				InsertedIP:     insertedIP,
			})
			index++
		}
	}

	slog.InfoContext(ctx, "Collected test results for batch sending",
		"work_order_id", workOrder.ID,
		"batch_count", len(batchItems),
		"skipped_count", skippedCount)

	// Check if there are no results to send
	if len(batchItems) == 0 {
		if skippedCount > 0 {
			return fmt.Errorf("no valid results to send to Nuha SIMRS (all %d results were skipped - empty values)", skippedCount)
		}
		return fmt.Errorf("no observation results found in work order")
	}

	// Send batch request
	nuhaService := NewNuhaService(c.BaseURL)
	batchReq := BatchInsertResultRequest{
		SessionID: c.SessionID,
		Data:      batchItems,
	}

	result, err := nuhaService.BatchInsertResults(ctx, batchReq)
	if err != nil {
		// All failed - update status as failed
		_ = c.workOrderRepo.UpdateSIMRSSentStatus(ctx, workOrder.ID, "FAILED", nil)
		return fmt.Errorf("failed to send batch results to Nuha SIMRS: %w", err)
	}

	slog.InfoContext(ctx, "Batch results sent to Nuha SIMRS successfully",
		"work_order_id", workOrder.ID,
		"sent_count", len(batchItems),
		"status", result.Response.Status,
		"message", result.Response.Message)

	// All results sent successfully - update status
	now := time.Now()
	if err := c.workOrderRepo.UpdateSIMRSSentStatus(ctx, workOrder.ID, "SENT", &now); err != nil {
		slog.ErrorContext(ctx, "Failed to update SIMRS sent status", "error", err)
		// Don't return error here as the main operation (sending results) succeeded
	}

	return nil
}

func parsePatientName(fullName string) (firstName, lastName string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "Unknown", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

// determineAbnormalFlag determines if a result value is normal, low, or high based on reference range
// Returns: "1" (Normal), "2" (Abnormal - High or Low)
func determineAbnormalFlag(resultValue, referenceRange string) string {
	if resultValue == "" || referenceRange == "" {
		return "1" // Default to normal if no data
	}

	// Try to parse result value as float
	value, err := strconv.ParseFloat(strings.TrimSpace(resultValue), 64)
	if err != nil {
		// If result is not numeric, return normal
		return "1"
	}

	// Clean reference range
	refRange := strings.TrimSpace(referenceRange)
	refRange = strings.ReplaceAll(refRange, " ", "")

	// Handle different reference range formats
	// Format: "min-max" (e.g., "10-20", "3.5-5.5")
	if strings.Contains(refRange, "-") && !strings.HasPrefix(refRange, "<") && !strings.HasPrefix(refRange, ">") {
		parts := strings.Split(refRange, "-")
		if len(parts) == 2 {
			min, err1 := strconv.ParseFloat(parts[0], 64)
			max, err2 := strconv.ParseFloat(parts[1], 64)

			if err1 == nil && err2 == nil {
				if value < min || value > max {
					return "2" // Abnormal (Low or High)
				}
				return "1" // Normal
			}
		}
	}

	// Format: "< max" or "<max" (e.g., "< 5", "<10")
	if strings.HasPrefix(refRange, "<") {
		maxStr := strings.TrimPrefix(refRange, "<")
		maxStr = strings.TrimPrefix(maxStr, "=")
		maxStr = strings.TrimSpace(maxStr)

		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			if strings.Contains(refRange, "=") {
				// <= case
				if value > max {
					return "2" // Abnormal (High)
				}
			} else {
				// < case
				if value >= max {
					return "2" // Abnormal (High)
				}
			}
			return "1" // Normal
		}
	}

	// Format: "> min" or ">min" (e.g., "> 100", ">50")
	if strings.HasPrefix(refRange, ">") {
		minStr := strings.TrimPrefix(refRange, ">")
		minStr = strings.TrimPrefix(minStr, "=")
		minStr = strings.TrimSpace(minStr)

		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			if strings.Contains(refRange, "=") {
				// >= case
				if value < min {
					return "2" // Abnormal (Low)
				}
			} else {
				// > case
				if value <= min {
					return "2" // Abnormal (Low)
				}
			}
			return "1" // Normal
		}
	}

	// Format: "min to max" (e.g., "10 to 20")
	if strings.Contains(strings.ToLower(refRange), "to") {
		parts := strings.Split(strings.ToLower(refRange), "to")
		if len(parts) == 2 {
			min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

			if err1 == nil && err2 == nil {
				if value < min || value > max {
					return "2" // Abnormal (Low or High)
				}
				return "1" // Normal
			}
		}
	}

	// If format not recognized, return normal
	return "1"
}
