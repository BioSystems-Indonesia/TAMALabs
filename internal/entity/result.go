package entity

import (
	"fmt"
	"log"
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
	Specimen
	TestResult map[string][]ResultTest `json:"test_result"`
}

type UpdateManyResultTestReq struct {
	Data []ResultTest `json:"data"`
}

type DeleteResultBulkReq struct {
	IDs []int64 `json:"ids"`
}

// ResultTest
type ResultTest struct {
	ID             int64          `json:"id"`
	TestTypeID     int64          `json:"test_type_id"`
	Test           string         `json:"test"`
	Result         *float64       `json:"result"`
	Unit           string         `json:"unit"`
	Category       string         `json:"category"`
	Abnormal       AbnormalResult `json:"abnormal"`
	ReferenceRange string         `json:"reference_range"`
	UpdatedAt      string         `json:"created_at"`

	History []ResultTest `json:"history"`
}

// CreateEmpty why we create empty result test? because we need the placeholder for the result test
// Result test can be filled manualy in frontend or from the result observation
// When we fill we need to show to the user, what testCode that we are filling
// What are the unit they will fill and what is the reference range
func (r ResultTest) CreateEmpty(request ObservationRequest) ResultTest {
	return ResultTest{
		ID:             0,
		Test:           request.TestCode,
		Result:         nil,
		TestTypeID:     int64(request.TestType.ID),
		Unit:           request.TestType.Unit,
		Category:       request.TestType.Category,
		ReferenceRange: fmt.Sprintf("%.2f - %.2f", request.TestType.LowRefRange, request.TestType.HighRefRange),
		UpdatedAt:      request.UpdatedAt.Format(time.RFC3339),
		Abnormal:       NormalResult,
		History:        []ResultTest{},
	}
}

func (r ResultTest) FromObservationResult(observation ObservationResult) ResultTest {
	resultTest := ResultTest{
		ID:             observation.ID,
		Test:           observation.Code,
		TestTypeID:     int64(observation.TestType.ID),
		Unit:           observation.TestType.Unit,
		Category:       observation.TestType.Category,
		ReferenceRange: fmt.Sprintf("%.2f - %.2f", observation.TestType.LowRefRange, observation.TestType.HighRefRange),
		UpdatedAt:      observation.UpdatedAt.Format(time.RFC3339),

		// Result, Abnormal will be filled below
		Result:   nil,
		Abnormal: NormalResult,
		// History will be filled by FillHistory,
	}

	// prevents panic
	if len(observation.Values) < 1 {
		log.Printf("values from observation %d is empty or negative", observation.ID)
		return resultTest
	}

	// only process the first value, if the observation have multiple values we need to handle it later
	result, err := strconv.ParseFloat(observation.Values[0], 64)
	if err != nil {
		log.Printf("parse observation.Values from observation %d failed: %v", observation.ID, err)
		return resultTest
	}

	// convert result into configuration units
	resultOrig := result
	if resultTest.Unit != "%" {
		// TODO make the ConvertCompount return resultOrig if err
		result, err = util.ConvertCompoundUnit(resultOrig, observation.Unit, observation.TestType.Unit)
		if err != nil {
			log.Printf(
				"convert compound unit from observation %d failed: convert %f from %s to %s: %v",
				observation.ID,
				resultOrig,
				observation.Unit,
				observation.TestType.Unit,
				err,
			)
			result = resultOrig
		}
	}

	resultTest.Result = &result

	resultTest.Abnormal = NormalResult
	if result <= observation.TestType.LowRefRange {
		resultTest.Abnormal = LowResult
	} else if result >= observation.TestType.HighRefRange {
		resultTest.Abnormal = HighResult
	}

	return resultTest
}

func (r ResultTest) FillHistory(history []ObservationResult) ResultTest {
	histories := make([]ResultTest, len(history))
	for i, h := range history {
		result := h.GetFirstValue()

		histories[i] = ResultTest{
			ID:             h.ID,
			Test:           h.Code,
			Result:         &result,
			TestTypeID:     int64(h.TestType.ID),
			Unit:           h.Unit,
			Category:       h.TestType.Category,
			ReferenceRange: fmt.Sprintf("%.2f - %.2f", h.TestType.LowRefRange, h.TestType.HighRefRange),
			UpdatedAt:      h.UpdatedAt.Format(time.RFC3339),
		}
	}

	r.History = histories

	return r
}

func (r ResultTest) ToObservationResult() (ObservationResult, error) {
	var lowRefRange, highRefRange float64

	_, err := fmt.Sscanf(r.ReferenceRange, "%f - %f", &lowRefRange, &highRefRange)
	if err != nil {
		return ObservationResult{}, err
	}

	updated_at, _ := time.Parse(time.RFC3339, r.UpdatedAt)

	values := make([]string, 0)
	if r.Result != nil {
		values = append(values, fmt.Sprintf("%f", *r.Result))
	}

	return ObservationResult{
		ID:   r.ID,
		Code: r.Test,
		TestType: TestType{
			ID:   int(r.TestTypeID),
			Name: r.Test, Unit: r.Unit, Category: r.Category,
			LowRefRange: lowRefRange, HighRefRange: highRefRange,
		},
		Values:    values,
		UpdatedAt: updated_at,
	}, nil
}

type AbnormalResult int32

const (
	NormalResult AbnormalResult = 0
	HighResult   AbnormalResult = 1
	LowResult    AbnormalResult = 2
)

type ResultGetManyRequest struct {
	GetManyRequest
	WorkOrderIDs    []int64  `query:"work_order_ids"`
	PatientIDs      []int64  `query:"patient_ids"`
	HasResult       bool     `query:"has_result"`
	WorkOrderStatus []string `query:"work_order_status"`
}
