package entity

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/util"
)

type WorkOrderStatus string

const (
	WorkOrderStatusNew            WorkOrderStatus = "NEW"
	WorkOrderStatusIncompleteSend WorkOrderStatus = "INCOMPLETE_SEND"
	WorkOrderStatusPending        WorkOrderStatus = "PENDING"
	WorkOrderCancelled            WorkOrderStatus = "CANCELLED"
	WorkOrderStatusCompleted      WorkOrderStatus = "SUCCESS"
)

type WorkOrderVerifiedStatus string

const (
	WorkOrderVerifiedStatusPending  WorkOrderVerifiedStatus = "PENDING"
	WorkOrderVerifiedStatusVerified WorkOrderVerifiedStatus = "VERIFIED"
	WorkOrderVerifiedStatusRejected WorkOrderVerifiedStatus = "REJECTED"
)

type WorkOrderCreateRequest struct {
	PatientID       int64                            `json:"patient_id" validate:"required"`
	TestTypes       []WorkOrderCreateRequestTestType `json:"test_types" validate:"required,min=1"`
	CreatedBy       int64                            `json:"created_by" validate:"required"`
	DoctorIDs       []int64                          `json:"doctor_ids" gorm:"-"`
	AnalyzerIDs     []int64                          `json:"analyzer_ids" gorm:"-"`
	TestTemplateIDs []int64                          `json:"test_template_ids" gorm:"-"`

	Barcode      string `json:"barcode" gorm:"column:barcode;index:work_order_barcode,unique"`
	BarcodeSIMRS string `json:"barcode_simrs" gorm:"column:barcode_simrs;index:work_order_barcode_simrs"`
}

type WorkOrderCreateRequestTestType struct {
	TestTypeID   int64  `json:"test_type_id" validate:"required"`
	TestTypeCode string `json:"test_type_code" validate:"required"`
	SpecimenType string `json:"specimen_type" validate:"required"`
}

type WorkOrder struct {
	ID                 int64           `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	Status             WorkOrderStatus `json:"status" gorm:"not null" validate:"work-order-status"`
	PatientID          int64           `json:"patient_id" gorm:"type:not null;default:0"`
	DeviceIDDeprecated int64           `json:"device_id" gorm:"column:device_id;type:not null;default:0"`
	//nolint:lll // tag cannot be shorter
	Barcode        string    `json:"barcode" gorm:"column:barcode;type:varchar(255);default:'';index:work_order_barcode,unique"`
	BarcodeSIMRS   string    `json:"barcode_simrs" gorm:"column:barcode_simrs;type:varchar(255);default:''"`
	VerifiedStatus string    `json:"verified_status" gorm:"column:verified_status;type:varchar(255);default:''"`
	CreatedBy      int64     `json:"created_by" gorm:"column:created_by;type:bigint;default:0"`
	LastUpdatedBy  int64     `json:"last_updated_by" gorm:"column:last_updated_by;type:bigint;default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"index:work_order_created_at"`
	UpdatedAt      time.Time `json:"updated_at" gorm:""`

	DoctorIDs       []int64 `json:"doctor_ids" gorm:"-"`
	AnalyzerIDs     []int64 `json:"analyzer_ids" gorm:"-"`
	TestTemplateIDs []int64 `json:"test_template_ids" gorm:"-"`

	Patient          Patient        `json:"patient" gorm:"foreignKey:PatientID;->" validate:"-"`
	Specimen         []Specimen     `json:"specimen_list,omitempty" gorm:"foreignKey:OrderID;->" validate:"-"`
	Devices          []Device       `json:"devices" gorm:"many2many:work_order_devices;->" validate:"-"`
	CreatedByUser    Admin          `json:"created_by_user" gorm:"foreignKey:CreatedBy;->" validate:"-"`
	LastUpdateByUser Admin          `json:"last_updated_by_user" gorm:"foreignKey:LastUpdatedBy;->" validate:"-"`
	Doctors          []Admin        `json:"doctors" gorm:"many2many:work_order_doctors;->" validate:"-"`
	Analyzers        []Admin        `json:"analyzers" gorm:"many2many:work_order_analyzers;->" validate:"-"`
	TestTemplates    []TestTemplate `json:"test_template" gorm:"many2many:work_order_test_templates;->" validate:"-"`

	TestResult        []TestResult `json:"test_result" gorm:"-"`
	TotalRequest      int64        `json:"total_request" gorm:"-"`
	TotalResultFilled int64        `json:"total_result_filled" gorm:"-"`
	PercentComplete   float64      `json:"percent_complete" gorm:"-"`
	HaveCompleteData  bool         `json:"have_complete_data" gorm:"-"`
}

func (wo *WorkOrder) GetFirstDoctor() Admin {
	if len(wo.Doctors) > 0 {
		return wo.Doctors[0]
	}

	return Admin{}
}

type ResultDetailOption struct {
	HideEmpty   bool
	HideHistory bool
}

// CalculateEGFRForResults calculates eGFR for creatinine results in the work order
func (wo *WorkOrder) CalculateEGFRForResults(ctx context.Context) {
	// Find creatinine results
	for i := range wo.TestResult {
		testResult := &wo.TestResult[i]

		// Check if this is a creatinine test
		if wo.isCreatinineTest(testResult.Test) && testResult.Result != nil {
			// Get patient information
			if wo.Patient.ID == 0 {
				continue // Skip if no patient data
			}

			// Calculate age
			age := util.CalculateAge(wo.Patient.Birthdate)

			// Convert creatinine to mg/dL if needed
			creatinineValue := *testResult.Result
			if testResult.Unit != "mg/dL" {
				convertedValue, err := util.ConvertCreatinineUnit(creatinineValue, testResult.Unit, "mg/dL")
				if err != nil {
					slog.WarnContext(ctx, "Failed to convert creatinine unit for eGFR calculation",
						"test_result_id", testResult.ID,
						"unit", testResult.Unit,
						"error", err)
					continue
				}
				creatinineValue = convertedValue
			}

			// Calculate eGFR using CKD-EPI formula
			var sex util.PatientSex
			if wo.Patient.Sex == PatientSexMale {
				sex = util.PatientSexMale
			} else {
				sex = util.PatientSexFemale
			}

			egfrResult := util.CalculateEGFRCKDEPI(creatinineValue, age, sex)

			// Add eGFR to the test result
			testResult.EGFR = &EGFRCalculation{
				Value:    egfrResult.Value,
				Formula:  egfrResult.Formula,
				Unit:     egfrResult.Unit,
				Category: egfrResult.Category,
			}

			// Also add eGFR to history results if they exist
			for j := range testResult.History {
				if testResult.History[j].Result != nil {
					historyCreatinine := *testResult.History[j].Result
					if testResult.History[j].Unit != "mg/dL" {
						convertedValue, err := util.ConvertCreatinineUnit(historyCreatinine, testResult.History[j].Unit, "mg/dL")
						if err != nil {
							continue
						}
						historyCreatinine = convertedValue
					}

					historyEGFR := util.CalculateEGFRCKDEPI(historyCreatinine, age, sex)
					testResult.History[j].EGFR = &EGFRCalculation{
						Value:    historyEGFR.Value,
						Formula:  historyEGFR.Formula,
						Unit:     historyEGFR.Unit,
						Category: historyEGFR.Category,
					}
				}
			}
		}
	}
}

// isCreatinineTest checks if the test code represents a creatinine test
func (wo *WorkOrder) isCreatinineTest(testCode string) bool {
	testCodeUpper := strings.ToUpper(testCode)
	creatinineCodes := []string{
		"CREATININE", "KREATININ",
	}

	for _, code := range creatinineCodes {
		if testCodeUpper == code {
			return true
		}
	}

	return false
}

func (wo *WorkOrder) FillResultDetail(opt ResultDetailOption) {
	var allObservationRequests []ObservationRequest
	var allObservationResults []ObservationResult
	for _, specimen := range wo.Specimen {
		allObservationRequests = append(allObservationRequests, specimen.ObservationRequest...)
		allObservationResults = append(allObservationResults, specimen.ObservationResult...)
	}

	allTests := make([]TestResult, len(allObservationRequests))
	// create the placeholder first
	for i, request := range allObservationRequests {
		allTests[i] = TestResult{}.CreateEmpty(request)
	}

	// sort by test code
	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Test < allTests[j].Test
	})

	// prepare the data that will be filled into the placeholder
	// prepare by grouping into code. But before that, sort by updated_at
	// the latest updated_at will be the first element
	sort.Slice(allObservationResults, func(i, j int) bool {
		return allObservationResults[i].CreatedAt.After(allObservationResults[j].CreatedAt)
	})

	// ok final step to create the order data
	testResults := map[string][]ObservationResult{}
	for _, observation := range allObservationResults {
		// TODO check whether this will create chaos in order or not
		testResults[observation.TestCode] = append(testResults[observation.TestCode], observation)
	}

	// fill the placeholder
	totalResultFilled := wo.pickDefaultResult(allTests, testResults, opt)

	wo.TotalRequest = int64(len(allObservationRequests))
	wo.TotalResultFilled = int64(totalResultFilled)
	wo.HaveCompleteData = len(allObservationRequests) == totalResultFilled
	if len(allObservationRequests) != 0 {
		wo.PercentComplete = float64(totalResultFilled) / float64(len(allObservationRequests))
	}

	if opt.HideEmpty {
		var filteredTests []TestResult
		for _, test := range allTests {
			if test.Result == nil || *test.Result == 0 {
				continue
			}
			filteredTests = append(filteredTests, test)
		}
		allTests = filteredTests
	}

	wo.TestResult = allTests
}

func (wo *WorkOrder) pickDefaultResult(
	allTests []TestResult,
	testResults map[string][]ObservationResult,
	opt ResultDetailOption,
) int {
	totalResultFilled := 0
	for i, test := range allTests {
		newTest := test
		history := testResults[test.Test]
		if len(history) > 0 {
			// Pick the latest history or the manually picked one
			pickedTest := history[0]
			for _, v := range history {
				if v.Picked {
					pickedTest = v
					break
				}
			}
			newTest = newTest.FromObservationResult(pickedTest, test.SpecimenType)

		}
		newTest = wo.fillTestHistory(newTest, history, opt)

		// or should be like this or we can just use the above code
		allTests[i] = newTest

		// count the filled result
		if newTest.Result != nil {
			totalResultFilled++
		}
	}
	return totalResultFilled
}

func (wo *WorkOrder) fillTestHistory(
	test TestResult,
	history []ObservationResult,
	opt ResultDetailOption,
) TestResult {
	if opt.HideHistory {
		return test
	}

	specimenTypes := make(map[int64]string)

	for _, specimen := range wo.Specimen {
		specimenTypes[int64(specimen.ID)] = specimen.Type
	}

	return test.FillHistory(history, specimenTypes)
}

func (wo *WorkOrder) FillData() {
	var doctorIDs []int64
	for _, d := range wo.Doctors {
		doctorIDs = append(doctorIDs, d.ID)
	}

	var analyzerIDs []int64
	for _, a := range wo.Analyzers {
		analyzerIDs = append(analyzerIDs, a.ID)
	}

	var testTemplateIDs []int64
	for _, t := range wo.TestTemplates {
		testTemplateIDs = append(testTemplateIDs, int64(t.ID))
	}

	wo.DoctorIDs = doctorIDs
	wo.AnalyzerIDs = analyzerIDs
	wo.TestTemplateIDs = testTemplateIDs
}

// TODO make all code in usecase to use this to fill TestResult
func (w *WorkOrder) FillTestResultDetail(hideEmpty bool) {
	var allObservationRequests []ObservationRequest
	var allObservationResults []ObservationResult

	// Create a map from specimen ID to specimen type for quick lookup
	specimenTypes := make(map[int64]string)

	for _, specimen := range w.Specimen {
		specimenTypes[int64(specimen.ID)] = specimen.Type
		allObservationRequests = append(allObservationRequests, specimen.ObservationRequest...)
		allObservationResults = append(allObservationResults, specimen.ObservationResult...)
	}

	allTests := make([]TestResult, len(allObservationRequests))
	// create the placeholder first
	for i, request := range allObservationRequests {
		// Find the corresponding specimen for this request
		// var correspondingSpecimen Specimen
		// for _, specimen := range w.Specimen {
		// 	if int64(specimen.ID) == request.SpecimenID {
		// 		correspondingSpecimen = specimen
		// 		break
		// 	}
		// }
		allTests[i] = TestResult{}.CreateEmpty(request)
	}

	// sort by test code
	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Test < allTests[j].Test
	})

	// prepare the data that will be filled into the placeholder
	// prepare by grouping into code. But before that, sort by updated_at
	// the latest updated_at will be the first element
	sort.Slice(allObservationResults, func(i, j int) bool {
		return allObservationResults[i].CreatedAt.After(allObservationResults[j].CreatedAt)
	})

	// ok final step to create the order data
	testResults := map[string][]ObservationResult{}
	for _, observation := range allObservationResults {
		// Create a composite key: testCode + specimenType for more precise matching
		// Get specimen type for this observation
		specimenType := specimenTypes[observation.SpecimenID]
		compositeKey := fmt.Sprintf("%s:%s", observation.TestCode, specimenType)

		testResults[compositeKey] = append(testResults[compositeKey], observation)
	}

	// fill the placeholder
	totalResultFilled := 0
	for i, test := range allTests {
		newTest := test
		// Create composite key for matching: testCode + specimenType
		compositeKey := fmt.Sprintf("%s:%s", test.Test, test.SpecimenType)
		history := testResults[compositeKey]

		if len(history) > 0 {
			// Pick the latest history or the manually picked one
			pickedTest := history[0]
			for _, v := range history {
				if v.Picked {
					pickedTest = v
					break
				}
			}

			// Get specimen type for this observation result
			specimenType := specimenTypes[pickedTest.SpecimenID]
			newTest = newTest.FromObservationResult(pickedTest, specimenType)
		}
		newTest = newTest.FillHistory(history, specimenTypes)

		// or should be like this or we can just use the above code
		allTests[i] = newTest

		// count the filled result
		if newTest.Result != nil {
			totalResultFilled++
		}
	}

	w.TotalRequest = int64(len(allObservationRequests))
	w.TotalResultFilled = int64(totalResultFilled)
	w.HaveCompleteData = len(allObservationRequests) == totalResultFilled
	if len(allObservationRequests) != 0 {
		w.PercentComplete = float64(totalResultFilled) / float64(len(allObservationRequests))
	}

	if hideEmpty {
		var filteredTests []TestResult
		for _, test := range allTests {
			if test.Result == nil || *test.Result == 0 {
				continue
			}
			filteredTests = append(filteredTests, test)
		}
		allTests = filteredTests
	}

	w.TestResult = allTests
}

type WorkOrderDoctor struct {
	WorkOrderID int64 `json:"work_order_id" gorm:"primaryKey" validate:"required"`
	AdminID     int64 `json:"admin_id" gorm:"primaryKey" validate:"required"`
}

type WorkOrderAnalyzer struct {
	WorkOrderID int64 `json:"work_order_id" gorm:"primaryKey" validate:"required"`
	AdminID     int64 `json:"admin_id" gorm:"primaryKey" validate:"required"`
}

type WorkOrderTestTemplate struct {
	WorkOrderID    int64 `json:"work_order_id" gorm:"primaryKey" validate:"required"`
	TestTemplateID int64 `json:"test_template_id" gorm:"primaryKey" validate:"required"`
}

type WorkOrderDevice struct {
	WorkOrderID int64     `json:"work_order_id" gorm:"not null;index:work_order_device_uniq,unique" validate:"required"`
	DeviceID    int64     `json:"device_id" gorm:"not null;index:work_order_device_uniq,unique" validate:"required"`
	CreatedAt   time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"not null"`
}

type WorkOrderRunRequest struct {
	DeviceID     int64   `json:"device_id" gorm:"-" validate:"required"`
	WorkOrderIDs []int64 `json:"work_order_ids" gorm:"-" validate:"required"`
	Urgent       bool    `json:"urgent" gorm:"-"`

	progressWriter chan WorkOrderRunStreamMessage
	patients       []Patient
	device         Device
}

func (w *WorkOrderRunRequest) SetPatients(patients []Patient) {
	w.patients = patients
}

func (w *WorkOrderRunRequest) GetPatients() []Patient {
	return w.patients
}

func (w *WorkOrderRunRequest) SetDevice(device Device) {
	w.device = device
}

func (w *WorkOrderRunRequest) GetDevice() Device {
	return w.device
}

func (w *WorkOrderRunRequest) ProgressWriter() chan WorkOrderRunStreamMessage {
	if w.progressWriter == nil {
		w.progressWriter = make(chan WorkOrderRunStreamMessage)
	}

	return w.progressWriter
}

func (w *WorkOrderRunRequest) SetProgressWriter(progress chan WorkOrderRunStreamMessage) {
	w.progressWriter = progress
}

type WorkOrderGetManyRequest struct {
	GetManyRequest
	BarcodeIds  []int64 `query:"barcode_ids"`
	PatientID   int64   `query:"patient_id"`
	PatientIDs  []int64 `query:"patient_ids"`
	SpecimenIDs []int64 `query:"specimen_ids"`
}

type WorkOrderStreamingResponseStatus string

const (
	WorkOrderStreamingResponseStatusDone       WorkOrderStreamingResponseStatus = "DONE"
	WorkOrderStreamingResponseStatusInProgress WorkOrderStreamingResponseStatus = "IN_PROGRESS"
)

type WorkOrderStreamingResponse string

func NewWorkOrderStreamingResponse(percentage float64, status WorkOrderStreamingResponseStatus) string {
	return fmt.Sprintf("data: percentage=%d&status=%s\n\n", int(percentage), status)
}

type SendPayloadRequest struct {
	Patients []Patient
	Device   Device
	Urgent   bool

	ProgressWriter chan WorkOrderRunStreamMessage
}

// WorkOrderRunStreamMessage represents a message sent from the use case to the controller.
type WorkOrderRunStreamMessage struct {
	Percentage float64
	Status     WorkOrderStreamingResponseStatus
	Error      error
	IsDone     bool
}
