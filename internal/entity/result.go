package entity

import "fmt"

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

type ResultTest struct {
	ID             int64          `json:"id"`
	Test           string         `json:"test"`
	Result         string         `json:"result"`
	Unit           string         `json:"unit"`
	Category       string         `json:"category"`
	Abnormal       AbnormalResult `json:"abnormal"`
	ReferenceRange string         `json:"reference_range"`
}

func (r ResultTest) FromObservationResult(observation ObservationResult) ResultTest {
	return ResultTest{
		ID:             observation.ID,
		Test:           observation.TestType.Name,
		Result:         observation.Values[0],
		Unit:           observation.TestType.Unit,
		Category:       observation.TestType.Category,
		ReferenceRange: fmt.Sprintf("%.2f - %.2f", observation.TestType.LowRefRange, observation.TestType.HighRefRange),
	}
}

func (r ResultTest) ToObservationResult() (ObservationResult, error) {
	var lowRefRange, highRefRange float64

	_, err := fmt.Sscanf(r.ReferenceRange, "%f - %f", &lowRefRange, &highRefRange)
	if err != nil {
		return ObservationResult{}, err
	}

	return ObservationResult{
		ID:   r.ID,
		Code: r.Test,
		TestType: TestType{
			Name: r.Test, Unit: r.Unit, Category: r.Category,
			LowRefRange: lowRefRange, HighRefRange: highRefRange,
		},
		Values: []string{r.Result},
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
