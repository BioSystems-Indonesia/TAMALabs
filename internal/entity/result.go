package entity

import (
	"fmt"
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

type ResultTest struct {
	ID             int64          `json:"id"`
	TestTypeID     int64          `json:"test_type_id"`
	Test           string         `json:"test"`
	Result         string         `json:"result"`
	Unit           string         `json:"unit"`
	Category       string         `json:"category"`
	Abnormal       AbnormalResult `json:"abnormal"`
	ReferenceRange string         `json:"reference_range"`
	UpdatedAt      string         `json:"created_at"`

	History []ResultTest `json:"history"`
}

func (r ResultTest) FromObservationResult(observation ObservationResult) ResultTest {
	resultTest := ResultTest{
		ID:             observation.ID,
		Test:           observation.Code,
		Result:         observation.Values[0],
		TestTypeID:     int64(observation.TestType.ID),
		Unit:           observation.TestType.Unit,
		Category:       observation.TestType.Category,
		ReferenceRange: fmt.Sprintf("%.2f - %.2f", observation.TestType.LowRefRange, observation.TestType.HighRefRange),
		UpdatedAt:      observation.UpdatedAt.Format(time.RFC3339),
	}

	// only process the first value, if the observation have multiple values we need to handle it later
	var result float64
	resultOrig, err := strconv.ParseFloat(resultTest.Result, 64)
	if err != nil {
		return resultTest
	}
	// convert result into configuration units
	if resultTest.Unit != "%" {
		result, err = util.ConvertCompoundUnit(resultOrig, observation.Unit, observation.TestType.Unit)
		if err != nil {
			result = resultOrig
		}
	}

	resultTest.Abnormal = NormalResult
	if result <= observation.TestType.LowRefRange {
		resultTest.Abnormal = LowResult
	} else if result >= observation.TestType.HighRefRange {
		resultTest.Abnormal = HighResult
	}

	return resultTest
}

func (r ResultTest) ToObservationResult() (ObservationResult, error) {
	var lowRefRange, highRefRange float64

	_, err := fmt.Sscanf(r.ReferenceRange, "%f - %f", &lowRefRange, &highRefRange)
	if err != nil {
		return ObservationResult{}, err
	}

	updated_at, _ := time.Parse(time.RFC3339, r.UpdatedAt)

	return ObservationResult{
		ID:   r.ID,
		Code: r.Test,
		TestType: TestType{
			ID:   int(r.TestTypeID),
			Name: r.Test, Unit: r.Unit, Category: r.Category,
			LowRefRange: lowRefRange, HighRefRange: highRefRange,
		},
		Values:    []string{r.Result},
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
