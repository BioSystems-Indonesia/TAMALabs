package entity

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/util"
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
	ID              int64          `json:"id"`
	SpecimenID      int64          `json:"specimen_id"`
	AliasCode       string         `json:"alias_code"`
	TestTypeID      int64          `json:"test_type_id"`
	Test            string         `json:"test"`
	Result          string         `json:"result"`
	FormattedResult *string        `json:"formatted_result"`
	Unit            string         `json:"unit"`
	Category        string         `json:"category"`
	SpecimenType    string         `json:"specimen_type"`
	Abnormal        AbnormalResult `json:"abnormal"`
	ReferenceRange  string         `json:"reference_range"`
	CreatedAt       string         `json:"created_at"`
	Picked          bool           `json:"picked"`
	TestType        TestType       `json:"test_type"`
	CreatedBy       Admin          `json:"created_by" validate:"-"`

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
func (r TestResult) CreateEmpty(request ObservationRequest) TestResult {
	decimal := request.TestType.Decimal
	if decimal < 0 {
		decimal = 0
	}

	return TestResult{
		ID:             0,
		SpecimenID:     request.SpecimenID,
		Test:           request.TestCode,
		Result:         "",
		TestTypeID:     int64(request.TestType.ID),
		Unit:           request.TestType.Unit,
		Category:       request.TestType.Category,
		ReferenceRange: request.TestType.GetReferenceRange(),
		CreatedAt:      request.UpdatedAt.Format(time.RFC3339),
		Abnormal:       NoDataResult,
		Picked:         false,
		History:        []TestResult{},
		TestType:       request.TestType,
		CreatedBy:      Admin{},
	}
}

func (r TestResult) FromObservationResult(observation ObservationResult, specimenType string) TestResult {
	// Use existing reference range if available, otherwise generate from TestType
	var referenceRange string
	if observation.ReferenceRange != "" {
		referenceRange = observation.ReferenceRange
	} else if observation.TestType.ID != 0 {
		// Use TestType's GetReferenceRange method to determine appropriate reference range
		referenceRange = observation.TestType.GetReferenceRange()
	} else {
		referenceRange = "" // Empty reference range for invalid/missing TestType
	}

	resultTest := TestResult{
		ID:             observation.ID,
		SpecimenID:     observation.SpecimenID,
		Test:           observation.TestCode,
		TestTypeID:     int64(observation.TestType.ID),
		Unit:           observation.TestType.Unit,
		Category:       observation.TestType.Category,
		ReferenceRange: referenceRange,
		CreatedAt:      observation.CreatedAt.Format(time.RFC3339),
		Picked:         observation.Picked,
		TestType:       observation.TestType,
		CreatedBy:      observation.CreatedByAdmin,

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

	resultStr := strconv.FormatFloat(result, 'f', 2, 64)
	resultTest.Result = resultStr

	if observation.TestType.Decimal < 1 {
		observation.TestType.Decimal = 2
	}

	formattedResult, err := strconv.ParseFloat(fmt.Sprintf("%.*f", observation.TestType.Decimal, result), 64)
	if err != nil {
		resultTest.FormattedResult = &resultTest.Result
	} else {
		formattedResultStr := strconv.FormatFloat(formattedResult, 'f', 2, 64)
		resultTest.FormattedResult = &formattedResultStr
	}

	resultTest.Abnormal = NormalResult
	if result <= observation.TestType.LowRefRange {
		resultTest.Abnormal = LowResult
	} else if result >= observation.TestType.HighRefRange {
		resultTest.Abnormal = HighResult
	}

	return resultTest
}

func (r TestResult) FillHistory(history []ObservationResult, specimenTypes map[int64]string) TestResult {
	histories := make([]TestResult, len(history))
	for i, h := range history {
		// Use string value to preserve qualitative results like "negative", "1+", etc.
		resultStr := h.GetFirstValueAsString()
		specimenType := specimenTypes[h.SpecimenID]

		decimal := h.TestType.Decimal
		if decimal < 0 {
			decimal = 0
		}

		histories[i] = TestResult{
			ID:             h.ID,
			SpecimenID:     h.SpecimenID,
			Test:           h.TestCode,
			Result:         resultStr, // Use string value directly
			TestTypeID:     int64(h.TestType.ID),
			Unit:           h.Unit,
			Category:       h.TestType.Category,
			ReferenceRange: h.TestType.GetReferenceRange(),
			CreatedAt:      h.CreatedAt.Format(time.RFC3339),
			Picked:         h.Picked,
			TestType:       h.TestType,
			SpecimenType:   specimenType,
			CreatedBy:      h.CreatedByAdmin,
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
