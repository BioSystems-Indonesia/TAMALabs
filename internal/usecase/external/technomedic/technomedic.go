package technomedic

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	externalEntity "github.com/BioSystems-Indonesia/TAMALabs/internal/entity/external"
	adminRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/admin"
	patientRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	subcategoryRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/sub_category"
	testTypeRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrderRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
)

type Usecase struct {
	testTypeRepo    *testTypeRepo.Repository
	adminRepo       *adminRepo.AdminRepository
	patientRepo     *patientRepo.PatientRepository
	workOrderRepo   *workOrderRepo.WorkOrderRepository
	subCategoryRepo *subcategoryRepo.Repository
	barcodeUC       usecase.BarcodeGenerator
}

func NewUsecase(
	testTypeRepo *testTypeRepo.Repository,
	adminRepo *adminRepo.AdminRepository,
	patientRepo *patientRepo.PatientRepository,
	workOrderRepo *workOrderRepo.WorkOrderRepository,
	subCategoryRepo *subcategoryRepo.Repository,
	barcodeUC usecase.BarcodeGenerator,
) *Usecase {
	return &Usecase{
		testTypeRepo:    testTypeRepo,
		adminRepo:       adminRepo,
		patientRepo:     patientRepo,
		workOrderRepo:   workOrderRepo,
		subCategoryRepo: subCategoryRepo,
		barcodeUC:       barcodeUC,
	}
}

// GetTestTypes returns all available test types
func (u *Usecase) GetTestTypes(ctx context.Context) ([]externalEntity.TechnoMedicTestType, error) {
	result, err := u.testTypeRepo.FindAll(ctx, &entity.TestTypeGetManyRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get test types: %w", err)
	}

	testTypes := make([]externalEntity.TechnoMedicTestType, len(result.Data))
	for i, tt := range result.Data {
		testTypes[i] = externalEntity.TechnoMedicTestType{
			ID:           strconv.Itoa(tt.ID),
			Code:         tt.Code,
			Name:         tt.Name,
			Category:     tt.Category,
			SubCategory:  tt.SubCategory,
			SpecimenType: tt.GetFirstType(),
			Unit:         tt.Unit,
		}
	}

	return testTypes, nil
}

// GetSubCategories returns all unique sub-categories
func (u *Usecase) GetSubCategories(ctx context.Context) ([]externalEntity.TechnoMedicSubCategory, error) {
	result, err := u.subCategoryRepo.FindAll(ctx, &entity.SubCategoryGetManyRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get sub-categories: %w", err)
	}

	subCategories := make([]externalEntity.TechnoMedicSubCategory, len(result.Data))
	for i, sc := range result.Data {
		subCategories[i] = externalEntity.TechnoMedicSubCategory{
			ID:          strconv.Itoa(sc.ID),
			Name:        sc.Name,
			Category:    sc.Category,
			Description: sc.Description,
		}
	}

	return subCategories, nil
}

// GetTestTypesBySubCategory returns all test types for a sub-category
func (u *Usecase) GetTestTypesBySubCategory(ctx context.Context, subCategoryID int) ([]externalEntity.TechnoMedicTestType, error) {
	testTypes, err := u.subCategoryRepo.GetTestTypesBySubCategoryID(ctx, subCategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test types by sub-category: %w", err)
	}

	result := make([]externalEntity.TechnoMedicTestType, len(testTypes))
	for i, tt := range testTypes {
		result[i] = externalEntity.TechnoMedicTestType{
			ID:           strconv.Itoa(tt.ID),
			Code:         tt.Code,
			Name:         tt.Name,
			Category:     tt.Category,
			SubCategory:  tt.SubCategory,
			SpecimenType: tt.GetFirstType(),
			Unit:         tt.Unit,
		}
	}

	return result, nil
}

// GetDoctors returns all doctors
func (u *Usecase) GetDoctors(ctx context.Context) ([]externalEntity.TechnoMedicDoctor, error) {
	doctors, err := u.adminRepo.FindAllByRole(ctx, entity.RoleDoctor)
	if err != nil {
		return nil, fmt.Errorf("failed to get doctors: %w", err)
	}

	result := make([]externalEntity.TechnoMedicDoctor, len(doctors))
	for i, doctor := range doctors {
		result[i] = externalEntity.TechnoMedicDoctor{
			ID:       doctor.ID,
			Fullname: doctor.Fullname,
			Username: doctor.Username,
			IsActive: doctor.IsActive,
		}
	}

	return result, nil
}

// GetAnalysts returns all analysts/analyzers
func (u *Usecase) GetAnalysts(ctx context.Context) ([]externalEntity.TechnoMedicAnalyst, error) {
	analysts, err := u.adminRepo.FindAllByRole(ctx, entity.RoleAnalyst)
	if err != nil {
		return nil, fmt.Errorf("failed to get analysts: %w", err)
	}

	result := make([]externalEntity.TechnoMedicAnalyst, len(analysts))
	for i, analyst := range analysts {
		result[i] = externalEntity.TechnoMedicAnalyst{
			ID:       analyst.ID,
			Fullname: analyst.Fullname,
			Username: analyst.Username,
			IsActive: analyst.IsActive,
		}
	}

	return result, nil
}

// CreateOrder creates a new order from TechnoMedic
func (u *Usecase) CreateOrder(ctx context.Context, req *externalEntity.TechnoMedicOrderRequest) error {
	slog.Info("Creating order from TechnoMedic",
		"no_order", req.NoOrder,
		"patient_id", req.Patient.PatientID,
	)

	// Parse birthdate
	birthdate, err := time.Parse("2006-01-02", req.Patient.Birthdate)
	if err != nil {
		return fmt.Errorf("invalid birthdate format, expected YYYY-MM-DD: %w", err)
	}

	// Find or create patient
	patient, err := u.findOrCreatePatient(ctx, req.Patient, birthdate)
	if err != nil {
		return fmt.Errorf("failed to find or create patient: %w", err)
	}

	// Get test types based on param_request or sub_category_request
	testTypeIDs, err := u.getTestTypesForOrder(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to get test types: %w", err)
	}

	if len(testTypeIDs) == 0 {
		return errors.New("no valid test types found for the request")
	}

	// Convert test type IDs to WorkOrderCreateRequestTestType
	testTypes := []entity.WorkOrderCreateRequestTestType{}
	for _, ttID := range testTypeIDs {
		testType, err := u.testTypeRepo.FindOneByID(ctx, int(ttID))
		if err != nil {
			slog.Warn("Failed to find test type", "id", ttID, "error", err)
			continue
		}
		testTypes = append(testTypes, entity.WorkOrderCreateRequestTestType{
			TestTypeID:   int64(testType.ID),
			TestTypeCode: testType.Code,
			SpecimenType: testType.GetFirstType(),
		})
	}

	if len(testTypes) == 0 {
		return errors.New("no valid test types could be loaded")
	}

	// Parse requested_at
	var requestedAt time.Time
	if req.RequestedAt != "" {
		requestedAt, err = time.Parse("2006-01-02 15:04:05", req.RequestedAt)
		if err != nil {
			// Try alternative format
			requestedAt, err = time.Parse("2006-01-02T15:04:05Z", req.RequestedAt)
			if err != nil {
				return fmt.Errorf("invalid requested_at format: %w", err)
			}
		}
	} else {
		requestedAt = time.Now()
	}

	// Generate barcode
	barcode, err := u.barcodeUC.NextOrderBarcode(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate barcode: %w", err)
	}

	// Create work order request
	createReq := &entity.WorkOrderCreateRequest{
		MedicalRecordNumber: patient.MedicalRecordNumber,
		PatientID:           patient.ID,
		Barcode:             barcode,
		TestTypes:           testTypes,
		CreatedBy:           -1, // System-generated from TechnoMedic
	}

	// Create work order
	_, err = u.workOrderRepo.Create(createReq)
	if err != nil {
		return fmt.Errorf("failed to create work order: %w", err)
	}

	// Update the work order to add barcode_simrs and source
	// We need to get the created work order and update it
	workOrder, err := u.workOrderRepo.FindOneByBarcode(ctx, barcode)
	if err != nil {
		return fmt.Errorf("failed to find created work order: %w", err)
	}

	workOrder.BarcodeSIMRS = req.NoOrder
	workOrder.Source = "technomedic"
	workOrder.CreatedAt = requestedAt

	err = u.workOrderRepo.Update(&workOrder)
	if err != nil {
		return fmt.Errorf("failed to update work order with order number: %w", err)
	}

	slog.Info("Successfully created order from TechnoMedic",
		"work_order_id", workOrder.ID,
		"barcode", barcode,
		"test_types_count", len(testTypes),
	)

	return nil
}

// GetOrder retrieves order details including results
func (u *Usecase) GetOrder(ctx context.Context, noOrder string) (*externalEntity.TechnoMedicGetOrderResponse, error) {
	// Find work order by barcode_simrs
	workOrder, err := u.workOrderRepo.GetBySIMRSBarcode(ctx, noOrder)
	if err != nil {
		if errors.Is(err, entity.ErrNotFound) {
			return nil, fmt.Errorf("order not found: %s", noOrder)
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	// Get patient details
	patient, err := u.patientRepo.FindByID(ctx, workOrder.PatientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient: %w", err)
	}

	// Build response
	fullName := patient.FirstName
	if patient.LastName != "" {
		fullName = patient.FirstName + " " + patient.LastName
	}

	response := &externalEntity.TechnoMedicGetOrderResponse{
		NoOrder:     workOrder.BarcodeSIMRS,
		Status:      string(workOrder.Status),
		RequestedBy: "",
		RequestedAt: workOrder.CreatedAt.Format("2006-01-02 15:04:05"),
		Patient: externalEntity.TechnoMedicPatientRequest{
			PatientID:           patient.MedicalRecordNumber,
			FullName:            fullName,
			Sex:                 string(patient.Sex),
			Address:             patient.Address,
			Birthdate:           patient.Birthdate.Format("2006-01-02"),
			MedicalRecordNumber: patient.MedicalRecordNumber,
			PhoneNumber:         patient.PhoneNumber,
		},
		SubCategories: []externalEntity.SubCategory{},
	}

	// Add completion/verification timestamps if available
	if workOrder.CompletedAt != nil {
		completedAt := workOrder.CompletedAt.Format("2006-01-02 15:04:05")
		response.CompletedAt = &completedAt
	}

	if workOrder.VerifiedAt != nil {
		verifiedAt := workOrder.VerifiedAt.Format("2006-01-02 15:04:05")
		response.VerifiedAt = &verifiedAt
	}

	if workOrder.VerifiedBy != nil {
		response.VerifiedBy = workOrder.VerifiedBy
	}

	// Group results by sub-category, and collect tests without sub-category
	subCatMap := make(map[string]*externalEntity.SubCategory)
	var noSubCategoryResults []externalEntity.Results

	for _, specimen := range workOrder.Specimen {
		for _, obsResult := range specimen.ObservationResult {
			if obsResult.TestType.ID == 0 {
				continue
			}

			// Get the first value from Values array
			value := ""
			if len(obsResult.Values) > 0 {
				value = obsResult.Values[0]
			}

			// Get the first flag from AbnormalFlag array
			flag := ""
			if len(obsResult.AbnormalFlag) > 0 {
				flag = obsResult.AbnormalFlag[0]
			}

			// Create result object
			result := externalEntity.Results{
				ID:           strconv.FormatInt(obsResult.ID, 10),
				Code:         obsResult.TestCode,
				CategoryName: obsResult.TestType.Name,
				Value:        value,
				SpecimenType: obsResult.TestType.GetFirstType(),
				Unit:         obsResult.Unit,
				Ref:          obsResult.ReferenceRange,
				Flag:         flag,
			}

			// Check if test has sub-category
			subCat := obsResult.TestType.SubCategory
			if subCat == "" {
				// No sub-category: add to root level parameters_result
				noSubCategoryResults = append(noSubCategoryResults, result)
			} else {
				// Has sub-category: group by sub-category
				if _, exists := subCatMap[subCat]; !exists {
					// Get the correct sub-category ID
					var subCatID string

					if obsResult.TestType.SubCategoryID != nil {
						// Use sub_category_id if available
						subCatID = strconv.Itoa(*obsResult.TestType.SubCategoryID)
					} else {
						// Fallback: try to find by name
						subCategory, err := u.subCategoryRepo.FindByName(ctx, subCat)
						if err == nil {
							subCatID = strconv.Itoa(subCategory.ID)
						} else {
							// Ultimate fallback: use test type ID (old behavior)
							subCatID = strconv.Itoa(obsResult.TestType.ID)
						}
					}

					subCatMap[subCat] = &externalEntity.SubCategory{
						ID:               subCatID,
						Name:             subCat,
						ParametersResult: []externalEntity.Results{},
					}
				}

				// Add result to sub-category
				subCatMap[subCat].ParametersResult = append(subCatMap[subCat].ParametersResult, result)
			}
		}
	}

	// Convert map to slice
	for _, sc := range subCatMap {
		response.SubCategories = append(response.SubCategories, *sc)
	}

	// Add tests without sub-category to root level
	if len(noSubCategoryResults) > 0 {
		response.ParametersResult = noSubCategoryResults
	}

	return response, nil
}

// Helper functions

func (u *Usecase) findOrCreatePatient(ctx context.Context, patientReq externalEntity.TechnoMedicPatientRequest, birthdate time.Time) (*entity.Patient, error) {
	// Try to find existing patient by medical record number
	var patient entity.Patient
	if patientReq.MedicalRecordNumber != "" {
		err := u.patientRepo.FindByMedicalRecordNumber(ctx, patientReq.MedicalRecordNumber, &patient)
		if err == nil {
			return &patient, nil
		}
	}

	// If not found, create new patient
	// Convert sex string to PatientSex type
	var sex entity.PatientSex
	switch patientReq.Sex {
	case "M":
		sex = entity.PatientSexMale
	case "F":
		sex = entity.PatientSexFemale
	default:
		sex = entity.PatientSexUnknown
	}

	newPatient := &entity.Patient{
		MedicalRecordNumber: patientReq.MedicalRecordNumber,
		FirstName:           patientReq.FullName, // Use FullName as FirstName for simplicity
		LastName:            "",
		Sex:                 sex,
		Address:             patientReq.Address,
		Birthdate:           birthdate,
		PhoneNumber:         patientReq.PhoneNumber,
	}

	err := u.patientRepo.CreateWithContext(ctx, newPatient)
	if err != nil {
		return nil, fmt.Errorf("failed to create patient: %w", err)
	}

	return newPatient, nil
}

func (u *Usecase) getTestTypesForOrder(ctx context.Context, req *externalEntity.TechnoMedicOrderRequest) ([]int64, error) {
	var testTypeIDs []int64

	// 1. Get by test_type_ids (direct IDs)
	if len(req.TestTypeIDs) > 0 {
		testTypeIDs = append(testTypeIDs, req.TestTypeIDs...)
	}

	// 2. Get by sub_category_ids
	if len(req.SubCategoryIDs) > 0 {
		for _, subCatID := range req.SubCategoryIDs {
			testTypes, err := u.subCategoryRepo.GetTestTypesBySubCategoryID(ctx, int(subCatID))
			if err != nil {
				slog.Warn("Failed to find test types by sub-category ID", "sub_category_id", subCatID, "error", err)
				continue
			}
			for _, tt := range testTypes {
				testTypeIDs = append(testTypeIDs, int64(tt.ID))
			}
		}
	}

	// 3. Get by param_request (test type codes)
	if len(req.ParamRequest) > 0 {
		for _, code := range req.ParamRequest {
			result, err := u.testTypeRepo.FindAll(ctx, &entity.TestTypeGetManyRequest{
				Code: code,
			})
			if err != nil {
				slog.Warn("Failed to find test type by code", "code", code, "error", err)
				continue
			}
			for _, tt := range result.Data {
				testTypeIDs = append(testTypeIDs, int64(tt.ID))
			}
		}
	}

	// 4. Get by sub_category_request (sub-category names)
	if len(req.SubCategoryRequest) > 0 {
		result, err := u.testTypeRepo.FindAll(ctx, &entity.TestTypeGetManyRequest{
			SubCategories: req.SubCategoryRequest,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to find test types by sub-category: %w", err)
		}

		for _, tt := range result.Data {
			testTypeIDs = append(testTypeIDs, int64(tt.ID))
		}
	}

	// Remove duplicates
	testTypeIDs = uniqueInt64(testTypeIDs)

	return testTypeIDs, nil
}

func uniqueInt64(slice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
