package nuha_simrs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
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

	today := time.Now().Format("2006-01-02")

	req := LabListRequest{
		SessionID: c.SessionID,
		ValidFrom: today,
		ValidTo:   today,
	}

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

	var processedCount, failedCount int
	for idx, labReg := range result.Response.List {
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
		Barcode:             "",
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
			// preserve simrs index from incoming package detail
			simrsIdx := detail.Index
			testTypes = append(testTypes, entity.WorkOrderCreateRequestTestType{
				TestTypeID:   int64(testType.ID),
				TestTypeCode: testType.Code,
				SpecimenType: testType.GetFirstType(),
				PackageID:    &packageID,
				SimrsIndex:   &simrsIdx,
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

func (c *SIMRSNuha) SendWorkOrderResults(
	ctx context.Context,
	barcodeSIMRS string,
	insertedUser string,
	insertedIP string,
) error {
	maskedUser := insertedUser
	if len(maskedUser) > 4 {
		maskedUser = maskedUser[:4] + "..."
	}
	slog.InfoContext(ctx, "SendWorkOrderResults called",
		"barcode_simrs", barcodeSIMRS,
		"inserted_user_masked", maskedUser,
		"inserted_ip", insertedIP,
	)

	labNumberStr := strings.TrimPrefix(barcodeSIMRS, "NUHA-")
	labNumber := 0
	if _, err := fmt.Sscanf(labNumberStr, "%d", &labNumber); err != nil || labNumber == 0 {
		slog.ErrorContext(ctx, "Invalid SIMRS barcode format",
			"barcode_simrs", barcodeSIMRS,
		)
		return fmt.Errorf("invalid SIMRS barcode format: %s", barcodeSIMRS)
	}

	wo, err := c.workOrderRepo.GetBySIMRSBarcode(ctx, barcodeSIMRS)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to fetch work order for SendWorkOrderResults",
			"barcode_simrs", barcodeSIMRS,
			"error", err,
		)
		return fmt.Errorf("failed to get work order: %w", err)
	}

	slog.InfoContext(ctx, "Work order fetched for SendWorkOrderResults",
		"work_order_id", wo.ID,
		"work_order_barcode", wo.Barcode,
		"simrs_barcode", wo.BarcodeSIMRS,
		"patient_id", wo.PatientID,
		"medical_record_number", wo.MedicalRecordNumber,
		"specimen_count", len(wo.Specimen),
		"simrs_sent_status", wo.SIMRSSentStatus,
		"verified_status", wo.VerifiedStatus,
	)

	return c.sendWorkOrderResults(ctx, wo, labNumber, insertedUser, insertedIP)
}

func (c *SIMRSNuha) SendWorkOrderResultsByID(
	ctx context.Context,
	workOrderIDStr string,
	insertedUser string,
	insertedIP string,
) error {
	slog.InfoContext(ctx, "SendWorkOrderResultsByID called",
		"work_order_id_str", workOrderIDStr,
		"inserted_user", insertedUser,
		"inserted_ip", insertedIP,
	)

	var workOrderID int64
	if _, err := fmt.Sscanf(workOrderIDStr, "%d", &workOrderID); err != nil {
		slog.ErrorContext(ctx, "Invalid work order ID format",
			"work_order_id_str", workOrderIDStr,
			"error", err,
		)
		return fmt.Errorf("invalid work order ID: %s", workOrderIDStr)
	}

	wo, err := c.workOrderRepo.FindOneForResult(workOrderID)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get work order by ID",
			"work_order_id", workOrderID,
			"error", err,
		)
		return fmt.Errorf("failed to get work order: %w", err)
	}

	if wo.BarcodeSIMRS == "" {
		slog.ErrorContext(ctx, "Work order missing SIMRS barcode",
			"work_order_id", workOrderID,
		)
		return fmt.Errorf("work order does not have SIMRS barcode, cannot send to Nuha")
	}

	labNumberStr := strings.TrimPrefix(wo.BarcodeSIMRS, "NUHA-")
	labNumber := 0
	if _, err := fmt.Sscanf(labNumberStr, "%d", &labNumber); err != nil || labNumber == 0 {
		slog.ErrorContext(ctx, "Invalid SIMRS barcode on work order",
			"barcode_simrs", wo.BarcodeSIMRS,
			"work_order_id", workOrderID,
			"error", err,
		)
		return fmt.Errorf("invalid SIMRS barcode format: %s", wo.BarcodeSIMRS)
	}

	slog.InfoContext(ctx, "Work order ready to send to Nuha",
		"work_order_id", workOrderID,
		"barcode_simrs", wo.BarcodeSIMRS,
		"lab_number", labNumber,
		"specimen_count", len(wo.Specimen),
	)

	return c.sendWorkOrderResults(ctx, wo, labNumber, insertedUser, insertedIP)
}

func (c *SIMRSNuha) sendWorkOrderResults(
	ctx context.Context,
	workOrder entity.WorkOrder,
	labNumber int,
	insertedUser string,
	insertedIP string,
) error {
	totalObsResults := 0
	for _, sp := range workOrder.Specimen {
		totalObsResults += len(sp.ObservationResult)
	}

	slog.InfoContext(ctx, "Sending work order results to Nuha SIMRS",
		"work_order_id", workOrder.ID,
		"work_order_barcode", workOrder.Barcode,
		"barcode_simrs", workOrder.BarcodeSIMRS,
		"lab_number", labNumber,
		"patient_id", workOrder.Patient.ID,
		"patient_mrn", workOrder.MedicalRecordNumber,
		"work_order_status", workOrder.Status,
		"simrs_sent_status", workOrder.SIMRSSentStatus,
		"verified_status", workOrder.VerifiedStatus,
		"specimen_count", len(workOrder.Specimen),
		"total_observation_results", totalObsResults,
		"inserted_user", insertedUser,
		"inserted_ip", insertedIP,
	)

	if len(workOrder.Specimen) == 0 {
		return fmt.Errorf("work order has no specimens to send")
	}

	var batchItems []BatchInsertResultItem
	skippedCount := 0
	skippedMissingTestID := 0
	index := 1

	for specIdx, specimen := range workOrder.Specimen {
		slog.InfoContext(ctx, "Processing specimen",
			"specimen_index", specIdx,
			"specimen_id", specimen.ID,
			"observation_result_count", len(specimen.ObservationResult))

// map test_type -> package_id and test_type -> simrs index (if available)
			testTypeToPackageID := make(map[int]*int)
			testTypeToSimrsIndex := make(map[int]*int)
			for _, obsReq := range specimen.ObservationRequest {
				if obsReq.TestTypeID != nil {
					testTypeToPackageID[*obsReq.TestTypeID] = obsReq.PackageID
					// preserve simrs index from observation request (nullable)
					testTypeToSimrsIndex[*obsReq.TestTypeID] = obsReq.SimrsIndex
			}
		}

		for obsIdx, obsResult := range specimen.ObservationResult {
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

			testIDFromAlias := 0
			aliasCode := strings.TrimSpace(obsResult.TestType.AliasCode)
			if aliasCode != "" {
				norm := strings.Trim(aliasCode, `"' `)
				if norm == "0" || strings.EqualFold(norm, "null") || strings.EqualFold(norm, "nil") {
					slog.DebugContext(ctx, "Skipping alias code because it's null/zero",
						"alias_code_raw", aliasCode,
					)
				} else {
					if id, err := strconv.Atoi(norm); err == nil {
						testIDFromAlias = id
					} else {
						re := regexp.MustCompile(`\d+`)
						m := re.FindString(aliasCode)
						if m != "" {
							if id2, err2 := strconv.Atoi(m); err2 == nil {
								testIDFromAlias = id2
								slog.InfoContext(ctx, "Parsed numeric part from alias code",
									"alias_code_raw", aliasCode,
									"extracted_digits", m,
									"parsed_test_id", id2,
								)
							} else {
								slog.WarnContext(ctx, "Failed to parse digits extracted from alias code",
									"alias_code", aliasCode,
									"extracted", m,
									"error", err2,
								)
							}
						} else {
							slog.DebugContext(ctx, "AliasCode contains no digits; cannot derive TestID",
								"alias_code", aliasCode,
							)
						}
					}
				}
			}

			valueStr := obsResult.Values[0]
			abnormalFlag := determineAbnormalFlag(valueStr, obsResult.ReferenceRange)

			testName := strings.TrimSpace(obsResult.TestType.Name)
			if testName == "" {
				if obsResult.TestCode != "" {
					testName = obsResult.TestCode
				} else if obsResult.Description != "" {
					testName = obsResult.Description
				} else {
					testName = fmt.Sprintf("unknown_test_%d", obsResult.ID)
				}
				slog.WarnContext(ctx, "ObservationResult has no TestType.Name; using fallback",
					"observation_result_id", obsResult.ID,
					"test_code", obsResult.TestCode,
					"resolved_test_name", testName,
					"test_type_id", obsResult.TestType.ID,
				)
			}

			slog.InfoContext(ctx, "Preparing observation result for Nuha",
				"specimen_id", specimen.ID,
				"observation_result_id", obsResult.ID,
				"test_name", obsResult.TestType.Name,
				"resolved_test_name", testName,
				"test_type_id", obsResult.TestType.ID,
				"alias_code", aliasCode,
				"parsed_test_id", testIDFromAlias,
				"value", valueStr,
				"unit", obsResult.Unit,
				"reference_range", obsResult.ReferenceRange,
				"abnormal_flag", abnormalFlag,
			)

			resultText := ""
			if obsResult.Comments != "" {
				resultText = obsResult.Comments
			}

			packageID := 0
			if pkgID, ok := testTypeToPackageID[obsResult.TestType.ID]; ok && pkgID != nil {
				packageID = *pkgID
			}

// prefer SIMRS (Nuha) index saved in ObservationRequest if available
				itemIndex := index
				if simIdx, ok := testTypeToSimrsIndex[obsResult.TestType.ID]; ok && simIdx != nil {
					itemIndex = *simIdx
				}

				item := BatchInsertResultItem{
					LabNumber:      labNumber,
					TestName:       testName,
					Result:         valueStr,
					ReferenceRange: obsResult.ReferenceRange,
					Abnormal:       abnormalFlag,
					Unit:           obsResult.Unit,
					TestID:         testIDFromAlias,
					PackageID:      packageID,
					Index:          itemIndex,
				ResultText:     resultText,
				InsertedUser:   insertedUser,
				InsertedIP:     insertedIP,
			}

			if item.TestID == 0 {
				skippedMissingTestID++
				skippedCount++
				slog.WarnContext(ctx, "Skipping batch item due missing TestID (alias not mapped)",
					"observation_result_id", obsResult.ID,
					"resolved_test_name", item.TestName,
					"test_code", obsResult.TestCode,
				)
				continue
			}

			batchItems = append(batchItems, item)
			if index <= 5 {
				if b, err := json.Marshal(item); err == nil {
					slog.DebugContext(ctx, "Batch item preview", "item_json", string(b))
				}
			}
			index++
		}
	}

	slog.InfoContext(ctx, "Collected test results for batch sending",
		"work_order_id", workOrder.ID,
		"batch_count", len(batchItems),
		"skipped_count", skippedCount,
		"skipped_missing_test_id", skippedMissingTestID,
		"total_observation_results", totalObsResults,
	)

	if len(batchItems) == 0 {
		if skippedCount > 0 {
			return fmt.Errorf("no valid results to send to Nuha SIMRS (all %d results were skipped - empty values)", skippedCount)
		}
		return fmt.Errorf("no observation results found in work order")
	}

	nuhaService := NewNuhaService(c.BaseURL)

	previewCount := 5
	var preview []BatchInsertResultItem
	if len(batchItems) <= previewCount {
		preview = batchItems
	} else {
		preview = batchItems[:previewCount]
	}
	if b, err := json.Marshal(preview); err == nil {
		slog.InfoContext(ctx, "Individual-insert payload preview (first items, truncated)",
			"preview_json", string(b),
			"total_items", len(batchItems),
		)
	}

	missingTestIDCount := 0
	missingSamples := make([]string, 0, 3)
	for i, it := range batchItems {
		if it.TestID == 0 {
			missingTestIDCount++
			if len(missingSamples) < 3 {
				missingSamples = append(missingSamples, fmt.Sprintf("%s(index=%d)", it.TestName, i))
			}
		}
	}
	if missingTestIDCount > 0 {
		slog.WarnContext(ctx, "Some items have missing TestID (alias not mapped)",
			"missing_count", missingTestIDCount,
			"samples", strings.Join(missingSamples, ", "),
		)
	}

	startSend := time.Now()
	sentCount := 0
	failedCount := 0
	failedSamples := make([]string, 0)

	for i, it := range batchItems {
		req := InsertResultRequest{
			SessionID:      c.SessionID,
			LabNumber:      it.LabNumber,
			TestName:       it.TestName,
			Result:         it.Result,
			Unit:           it.Unit,
			ReferenceRange: it.ReferenceRange,
			Abnormal:       it.Abnormal,
			Description:    "",
			Notes:          "",
			TestID:         it.TestID,
			ResultText:     "",
			PackageID:      it.PackageID,
			Spacing:        "",
			Index:          it.Index,
			InsertedUser:   it.InsertedUser,
			InsertedIP:     it.InsertedIP,
		}

		if _, err := nuhaService.InsertResult(ctx, req); err != nil {
			failedCount++
			failedSamples = append(failedSamples, fmt.Sprintf("%s(index=%d): %v", it.TestName, i, err))
			slog.ErrorContext(ctx, "Failed to send single result to Nuha",
				"work_order_id", workOrder.ID,
				"index", i,
				"test_name", it.TestName,
				"error", err,
			)
			continue
		}

		sentCount++
	}

	durationMs := time.Since(startSend).Milliseconds()

	if failedCount > 0 {
		_ = c.workOrderRepo.UpdateSIMRSSentStatus(ctx, workOrder.ID, "FAILED", nil)
		slog.ErrorContext(ctx, "One or more individual result inserts failed",
			"work_order_id", workOrder.ID,
			"sent_count", sentCount,
			"failed_count", failedCount,
			"duration_ms", durationMs,
			"failed_samples_preview", strings.Join(failedSamples, ", "),
		)
		return fmt.Errorf("failed to send %d/%d results to Nuha SIMRS", failedCount, len(batchItems))
	}

	slog.InfoContext(ctx, "All individual results sent to Nuha SIMRS successfully",
		"work_order_id", workOrder.ID,
		"sent_count", sentCount,
		"duration_ms", durationMs,
	)

	now := time.Now()
	if err := c.workOrderRepo.UpdateSIMRSSentStatus(ctx, workOrder.ID, "SENT", &now); err != nil {
		slog.ErrorContext(ctx, "Failed to update SIMRS sent status", "error", err)
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

func determineAbnormalFlag(resultValue, referenceRange string) string {
	if resultValue == "" || referenceRange == "" {
		return "1"
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(resultValue), 64)
	if err != nil {
		return "1"
	}

	refRange := strings.TrimSpace(referenceRange)
	refRange = strings.ReplaceAll(refRange, " ", "")

	if strings.Contains(refRange, "-") && !strings.HasPrefix(refRange, "<") && !strings.HasPrefix(refRange, ">") {
		parts := strings.Split(refRange, "-")
		if len(parts) == 2 {
			min, err1 := strconv.ParseFloat(parts[0], 64)
			max, err2 := strconv.ParseFloat(parts[1], 64)

			if err1 == nil && err2 == nil {
				if value < min || value > max {
					return "2"
				}
				return "1"
			}
		}
	}

	if strings.HasPrefix(refRange, "<") {
		maxStr := strings.TrimPrefix(refRange, "<")
		maxStr = strings.TrimPrefix(maxStr, "=")
		maxStr = strings.TrimSpace(maxStr)

		if max, err := strconv.ParseFloat(maxStr, 64); err == nil {
			if strings.Contains(refRange, "=") {
				if value > max {
					return "2"
				}
			} else {
				if value >= max {
					return "2"
				}
			}
			return "1"
		}
	}

	if strings.HasPrefix(refRange, ">") {
		minStr := strings.TrimPrefix(refRange, ">")
		minStr = strings.TrimPrefix(minStr, "=")
		minStr = strings.TrimSpace(minStr)

		if min, err := strconv.ParseFloat(minStr, 64); err == nil {
			if strings.Contains(refRange, "=") {
				if value < min {
					return "2"
				}
			} else {
				if value <= min {
					return "2"
				}
			}
			return "1"
		}
	}

	if strings.Contains(strings.ToLower(refRange), "to") {
		parts := strings.Split(strings.ToLower(refRange), "to")
		if len(parts) == 2 {
			min, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			max, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

			if err1 == nil && err2 == nil {
				if value < min || value > max {
					return "2"
				}
				return "1"
			}
		}
	}

	return "1"
}
