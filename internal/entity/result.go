package entity

import (
	"fmt"
	"log/slog"
	"strconv"
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
	Result          *float64       `json:"result"`
	FormattedResult *float64       `json:"formatted_result"`
	Unit            string         `json:"unit"`
	Category        string         `json:"category"`
	SpecimenType    string         `json:"specimen_type"`
	Abnormal        AbnormalResult `json:"abnormal"`
	ReferenceRange  string         `json:"reference_range"`
	CreatedAt       string         `json:"created_at"`
	Picked          bool           `json:"picked"`
	TestType        TestType       `json:"test_type"`

	History []TestResult `json:"history"`

	// Calculated fields (not stored in database)
	EGFR *EGFRCalculation `json:"egfr,omitempty"`
}

func (r TestResult) GetResult() float64 {
	if r.Result == nil {
		return 0
	}

	return *r.Result
}

func (r TestResult) GetFormattedResult() float64 {
	if r.FormattedResult == nil {
		return 0
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
		Result:         nil,
		TestTypeID:     int64(request.TestType.ID),
		Unit:           request.TestType.Unit,
		Category:       request.TestType.Category,
		ReferenceRange: fmt.Sprintf("%.*f - %.*f", decimal, request.TestType.LowRefRange, decimal, request.TestType.HighRefRange),
		CreatedAt:      request.UpdatedAt.Format(time.RFC3339),
		Abnormal:       NoDataResult,
		Picked:         false,
		History:        []TestResult{},
		TestType:       request.TestType,
	}
}

func (r TestResult) FromObservationResult(observation ObservationResult, specimenType string) TestResult {
	resultTest := TestResult{
		ID:             observation.ID,
		SpecimenID:     observation.SpecimenID,
		Test:           observation.TestCode,
		TestTypeID:     int64(observation.TestType.ID),
		Unit:           observation.TestType.Unit,
		Category:       observation.TestType.Category,
		ReferenceRange: fmt.Sprintf("%.*f - %.*f", observation.TestType.Decimal, observation.TestType.LowRefRange, observation.TestType.Decimal, observation.TestType.HighRefRange),
		CreatedAt:      observation.CreatedAt.Format(time.RFC3339),
		Picked:         observation.Picked,
		TestType:       observation.TestType,

		// Result, Abnormal will be filled below
		Result:   nil,
		Abnormal: NoDataResult,
		// History will be filled by FillHistory,
	}

	// prevents panic
	if len(observation.Values) < 1 {
		slog.Info("values from observation is empty or negative", "id", observation.ID)
		return resultTest
	}

	// only process the first value, if the observation have multiple values we need to handle it later
	result, err := strconv.ParseFloat(observation.Values[0], 64)
	if err != nil {
		slog.Info("parse observation.Values from observation failed", "id", observation.ID, "error", err)
		return resultTest
	}

	// convert result into configuration units
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

	resultTest.Result = &result
	if observation.TestType.Decimal < 1 {
		observation.TestType.Decimal = 2
	}

	formattedResult, err := strconv.ParseFloat(fmt.Sprintf("%.*f", observation.TestType.Decimal, result), 64)
	if err != nil {
		resultTest.FormattedResult = resultTest.Result
	}
	resultTest.FormattedResult = &formattedResult

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
		result := h.GetFirstValue()
		specimenType := specimenTypes[h.SpecimenID]

		decimal := h.TestType.Decimal
		if decimal < 0 {
			decimal = 0
		}

		histories[i] = TestResult{
			ID:             h.ID,
			SpecimenID:     h.SpecimenID,
			Test:           h.TestCode,
			Result:         &result,
			TestTypeID:     int64(h.TestType.ID),
			Unit:           h.Unit,
			Category:       h.TestType.Category,
			ReferenceRange: fmt.Sprintf("%.*f - %.*f", decimal, h.TestType.LowRefRange, decimal, h.TestType.HighRefRange),
			CreatedAt:      h.CreatedAt.Format(time.RFC3339),
			Picked:         h.Picked,
			TestType:       h.TestType,
			SpecimenType:   specimenType,
		}
	}

	r.History = histories

	return r
}

type AbnormalResult int32

const (
	NormalResult AbnormalResult = 0
	HighResult   AbnormalResult = 1
	LowResult    AbnormalResult = 2
	NoDataResult AbnormalResult = 3
)

type ResultGetManyRequest struct {
	GetManyRequest
	BarcodeIDs      []int64  `query:"barcode_ids"`
	PatientIDs      []int64  `query:"patient_ids"`
	DoctorIDs       []int64  `query:"doctor_ids"`
	HasResult       bool     `query:"has_result"`
	WorkOrderStatus []string `query:"work_order_status"`
}
