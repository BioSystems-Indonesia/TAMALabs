package entity

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
)

type Result struct {
	Specimen
}

type ResultPaginationResponse struct {
	Data []Result `json:"data"`
	// PaginationResponse
}

type ResultDetail struct {
	WorkOrder
	TestResult map[string][]TestResult `json:"test_result"`
	PrevID     int64                   `json:"prev_id"`
	NextID     int64                   `json:"next_id"`
}

type UpdateManyResultTestReq struct {
	Data []TestResult `json:"data"`
}

type DeleteResultBulkReq struct {
	IDs []int64 `json:"ids"`
}

// TestResult are representation of result that will be manipulated and managed
// in our LIS. it is used to store or add the test manually or to show
// the test to others
type TestResult struct {
	ID                     int64          `json:"id"`
	SpecimenID             int64          `json:"specimen_id"`
	AliasCode              string         `json:"alias_code"`
	TestTypeID             int64          `json:"test_type_id"`
	Test                   string         `json:"test"`
	Result                 string         `json:"result"`
	FormattedResult        *string        `json:"formatted_result"`
	Unit                   string         `json:"unit"`
	Category               string         `json:"category"`
	SpecimenType           string         `json:"specimen_type"`
	Abnormal               AbnormalResult `json:"abnormal"`
	ReferenceRange         string         `json:"reference_range"`
	ComputedReferenceRange string         `json:"computed_reference_range,omitempty"` // Always use TestType.GetReferenceRange() for proper decimal formatting
	CreatedAt              string         `json:"created_at"`
	Picked                 bool           `json:"picked"`
	TestType               TestType       `json:"test_type"`
	CreatedBy              Admin          `json:"created_by" validate:"-"`

	History []TestResult `json:"history"`

	// Calculated fields (not stored in database)
	EGFR *EGFRCalculation `json:"egfr,omitempty"`
}

func (r TestResult) GetResult() string {
	if r.Result == "" {
		return ""
	}

	return r.Result
}

func (r TestResult) GetFormattedResult() string {
	if r.FormattedResult == nil {
		return ""
	}

	return *r.FormattedResult
}

// CreateEmpty why we create empty result test? because we need the placeholder for the result test
// Result test can be filled manualy in frontend or from the result observation
// When we fill we need to show to the user, what testCode that we are filling
// What are the unit they will fill and what is the reference range
func (r TestResult) CreateEmpty(request ObservationRequest, patient *Patient) TestResult {
	decimal := request.TestType.Decimal
	if decimal < 0 {
		decimal = 0
	}

	// Use patient-specific reference range if patient data is available
	var computedRefRange string
	if patient != nil {
		age := patient.GetAge()
		gender := patient.GetGenderString()
		computedRefRange = request.TestType.GetReferenceRangeForPatient(&age, gender)
	} else {
		computedRefRange = request.TestType.GetReferenceRange()
	}

	// Use TestType.Name for display, fallback to TestCode if Name is empty
	testName := request.TestType.Name
	if testName == "" {
		testName = request.TestCode
	}

	return TestResult{
		ID:                     0,
		SpecimenID:             request.SpecimenID,
		Test:                   testName, // Use TestType.Name (GDP/GDS) for display, not TestCode (GLUCOSE)
		Result:                 "",
		TestTypeID:             int64(request.TestType.ID),
		Unit:                   request.TestType.Unit,
		Category:               request.TestType.Category,
		ReferenceRange:         computedRefRange,
		ComputedReferenceRange: computedRefRange,
		CreatedAt:              request.UpdatedAt.Format(time.RFC3339),
		Abnormal:               NoDataResult,
		Picked:                 false,
		History:                []TestResult{},
		TestType:               request.TestType,
		CreatedBy:              Admin{},
	}
}

func (r TestResult) FromObservationResult(observation ObservationResult, specimenType string, patient *Patient) TestResult {
	// Always compute reference range based on patient data if available
	var referenceRange string
	var computedRefRange string
	var lowRefRange, highRefRange float64

	if observation.TestType.ID != 0 {
		// Use patient-specific reference range if available
		if patient != nil {
			age := patient.GetAge()
			gender := patient.GetGenderString()

			// Get the specific reference range for this patient
			specificRange := observation.TestType.GetSpecificReferenceRangeForPatient(&age, gender)
			if specificRange != nil && specificRange.LowRefRange != nil && specificRange.HighRefRange != nil {
				// Use specific range values
				lowRefRange = *specificRange.LowRefRange
				highRefRange = *specificRange.HighRefRange
				referenceRange = observation.TestType.GetReferenceRangeForPatient(&age, gender)
			} else {
				// Fallback to global range
				lowRefRange = observation.TestType.LowRefRange
				highRefRange = observation.TestType.HighRefRange
				referenceRange = observation.TestType.GetReferenceRange()
			}
			computedRefRange = referenceRange
		} else {
			// No patient data, use global range
			lowRefRange = observation.TestType.LowRefRange
			highRefRange = observation.TestType.HighRefRange
			referenceRange = observation.TestType.GetReferenceRange()
			computedRefRange = referenceRange
		}
	} else {
		referenceRange = "" // Empty reference range for invalid/missing TestType
		computedRefRange = ""
		lowRefRange = 0
		highRefRange = 0
	}

	resultTest := TestResult{
		ID:                     observation.ID,
		SpecimenID:             observation.SpecimenID,
		Test:                   observation.TestCode,
		TestTypeID:             int64(observation.TestType.ID),
		Unit:                   observation.TestType.Unit,
		Category:               observation.TestType.Category,
		ReferenceRange:         referenceRange,
		ComputedReferenceRange: computedRefRange,
		CreatedAt:              observation.CreatedAt.Format(time.RFC3339),
		Picked:                 observation.Picked,
		TestType:               observation.TestType,
		CreatedBy:              observation.CreatedByAdmin,

		// Result, Abnormal will be filled below
		Result:   "",
		Abnormal: NoDataResult,
		// History will be filled by FillHistory,
	}

	// prevents panic
	if len(observation.Values) < 1 {
		slog.Info("values from observation is empty or negative", "id", observation.ID)
		return resultTest
	}

	// Get the first value as string directly
	valueStr := observation.Values[0]
	resultTest.Result = valueStr

	// Try to parse as float for unit conversion and abnormal checking
	result, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		// If it's not a number (like "3+", "positive", etc), check for specific patterns
		slog.Info("observation value is not numeric, checking for patterns", "id", observation.ID, "value", valueStr)
		resultTest.FormattedResult = &valueStr

		// Check for positive/negative patterns
		if strings.Contains(strings.ToLower(valueStr), "+") {
			resultTest.Abnormal = PositiveResult
		} else if strings.Contains(strings.ToLower(valueStr), "-") ||
			strings.Contains(strings.ToLower(valueStr), "negatif") ||
			strings.Contains(strings.ToLower(valueStr), "negative") {
			resultTest.Abnormal = NegativeResult
		} else {
			resultTest.Abnormal = NoDataResult // Can't determine abnormal status
		}
		return resultTest
	}

	// convert result into configuration units (only for numeric values)
	resultOrig := result
	if resultTest.Unit != "%" {
		// TODO make the ConvertCompount return resultOrig if err
		result, err = util.ConvertCompoundUnit(resultOrig, observation.Unit, observation.TestType.Unit)
		if err != nil {
			slog.Warn(
				"convert compound unit from observation failed",
				"id", observation.ID,
				"result", resultOrig,
				"unit", observation.Unit,
				"test_type_unit", observation.TestType.Unit,
				"error", err,
			)
			result = resultOrig
		}
	}

	// Use TestType decimal setting for all formatting
	decimal := observation.TestType.Decimal
	if decimal < 0 {
		decimal = 0 // Default to 0 (whole numbers) if negative
	}

	// Format result string with proper decimal places
	resultStr := strconv.FormatFloat(result, 'f', decimal, 64)
	resultTest.Result = resultStr

	// Format for display - use the same decimal precision
	formattedResultStr := fmt.Sprintf("%.*f", decimal, result)
	resultTest.FormattedResult = &formattedResultStr

	// Determine abnormal status using the correct reference range values
	resultTest.Abnormal = NormalResult
	if result <= lowRefRange {
		resultTest.Abnormal = LowResult
	} else if result >= highRefRange {
		resultTest.Abnormal = HighResult
	}

	return resultTest
}

func (r TestResult) FillHistory(history []ObservationResult, specimenTypes map[int64]string, patient *Patient) TestResult {
	histories := make([]TestResult, len(history))
	for i, h := range history {
		// Use string value to preserve qualitative results like "negative", "1+", etc.
		resultStr := h.GetFirstValueAsString()
		specimenType := specimenTypes[h.SpecimenID]

		decimal := h.TestType.Decimal
		if decimal < 0 {
			decimal = 0
		}

		// Try to format numeric values with proper decimal places
		var formattedResult *string
		if numValue, err := strconv.ParseFloat(resultStr, 64); err == nil {
			// It's a numeric value, format with proper decimals
			formatted := fmt.Sprintf("%.*f", decimal, numValue)
			formattedResult = &formatted
		} else {
			// It's a qualitative value (like "negative", "1+"), keep as is
			formattedResult = &resultStr
		}

		// Compute reference range with proper decimal formatting - use patient-specific if available
		var computedRefRange string
		if patient != nil {
			age := patient.GetAge()
			gender := patient.GetGenderString()
			computedRefRange = h.TestType.GetReferenceRangeForPatient(&age, gender)
		} else {
			computedRefRange = h.TestType.GetReferenceRange()
		}

		histories[i] = TestResult{
			ID:                     h.ID,
			SpecimenID:             h.SpecimenID,
			Test:                   h.TestCode,
			Result:                 resultStr, // Use string value directly
			FormattedResult:        formattedResult,
			TestTypeID:             int64(h.TestType.ID),
			Unit:                   h.Unit,
			Category:               h.TestType.Category,
			ReferenceRange:         computedRefRange,
			ComputedReferenceRange: computedRefRange,
			CreatedAt:              h.CreatedAt.Format(time.RFC3339),
			Picked:                 h.Picked,
			TestType:               h.TestType,
			SpecimenType:           specimenType,
			CreatedBy:              h.CreatedByAdmin,
		}
	}

	r.History = histories

	return r
}

type AbnormalResult int32

const (
	NormalResult   AbnormalResult = 0
	HighResult     AbnormalResult = 1
	LowResult      AbnormalResult = 2
	NoDataResult   AbnormalResult = 3
	PositiveResult AbnormalResult = 4
	NegativeResult AbnormalResult = 5
)

type ResultGetManyRequest struct {
	GetManyRequest
	BarcodeIDs      []int64  `query:"barcode_ids"`
	PatientIDs      []int64  `query:"patient_ids"`
	DoctorIDs       []int64  `query:"doctor_ids"`
	HasResult       bool     `query:"has_result"`
	WorkOrderStatus []string `query:"work_order_status"`
}
