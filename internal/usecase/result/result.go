package result

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	specimenRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	workOrderRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
)

type Usecase struct {
	resultRepository    *observation_result.Repository
	workOrderRepository *workOrderRepo.WorkOrderRepository
	specimenRepository  *specimenRepo.Repository
	testTypeRepository  *test_type.Repository
}

func NewUsecase(
	resultRepository *observation_result.Repository,
	workOrderRepository *workOrderRepo.WorkOrderRepository,
	specimenRepository *specimenRepo.Repository,
	testTypeRepository *test_type.Repository,
) *Usecase {
	return &Usecase{
		resultRepository:    resultRepository,
		workOrderRepository: workOrderRepository,
		specimenRepository:  specimenRepository,
		testTypeRepository:  testTypeRepository,
	}
}

func (u *Usecase) Results(
	ctx context.Context,
	req *entity.ResultGetManyRequest,
) (entity.PaginationResponse[entity.Specimen], error) {
	resp, err := u.specimenRepository.FindAllForResult(ctx, req)
	if err != nil {
		return entity.PaginationResponse[entity.Specimen]{}, err
	}

	for i := range resp.Data {
		resp.Data[i].TestResult = u.processResultDetail(resp.Data[i])
	}

	return resp, nil
}

func (u *Usecase) ResultDetail(ctx context.Context, specimenID int64) (entity.ResultDetail, error) {
	specimen, err := u.specimenRepository.FindOne(ctx, specimenID)
	if err != nil {
		return entity.ResultDetail{}, err
	}

	return entity.ResultDetail{
		Specimen:   specimen,
		TestResult: u.groupResultInCategory(u.processResultDetail(specimen)),
	}, nil
}

// PutTestResult will create ObservationResult
// set the value, unit and everyting else
// and will prepend it in TestResult.History
func (u *Usecase) PutTestResult(ctx context.Context, result entity.TestResult) (entity.TestResult, error) {
	oldResult := result

	obs := entity.ObservationResult{
		SpecimenID:     result.SpecimenID,
		Code:           result.Test,
		Unit:           result.Unit,
		ReferenceRange: result.ReferenceRange,
	}

	if result.Result != nil {
		obs.Values = append(obs.Values, fmt.Sprintf("%f", *result.Result))
	}

	err := u.resultRepository.Create(ctx, &obs)
	if err != nil {
		return result, err
	}

	obs.TestType, err = u.testTypeRepository.FindOneByCode(ctx, obs.Code)
	if err != nil {
		log.Printf("cannot fill test type for result %d: %v", obs.ID, err)
	}

	result = result.FromObservationResult(obs)
	// Hack so the front end is not add first then replace
	// TODO maybe need to find better way than this
	if oldResult.ID != 0 {
		result.ID = oldResult.ID
	}

	history, err := u.resultRepository.FindHistory(ctx, obs)
	if err != nil {
		log.Printf("cannot get history: %v", err)
	}

	result = result.FillHistory(history)

	return result, nil
}

func (u *Usecase) DeleteTestResult(context context.Context, id int64) (entity.ObservationResult, error) {
	return u.resultRepository.Delete(context, id)
}

func (u *Usecase) groupResultInCategory(tests []entity.TestResult) map[string][]entity.TestResult {
	var resultTestsCategory = map[string][]entity.TestResult{}
	for _, resultTest := range tests {
		categoryName := resultTest.Category

		categoryTestResults, ok := resultTestsCategory[categoryName]
		if ok {
			resultTestsCategory[categoryName] = append(categoryTestResults, resultTest)
		} else {
			resultTestsCategory[categoryName] = []entity.TestResult{resultTest}
		}
	}

	return resultTestsCategory
}

func (u *Usecase) processResultDetail(specimen entity.Specimen) []entity.TestResult {
	var tests = make([]entity.TestResult, len(specimen.ObservationRequest))

	// create the placeholder first
	for i, request := range specimen.ObservationRequest {
		tests[i] = entity.TestResult{}.CreateEmpty(request)
	}

	// sort by test code
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].Test < tests[j].Test
	})

	// prepare the data that will be filled into the placeholder
	// prepare by grouping into code. But before that, sort by updated_at
	// the latest updated_at will be the first element
	sort.Slice(specimen.ObservationResult, func(i, j int) bool {
		return specimen.ObservationResult[i].UpdatedAt.After(specimen.ObservationResult[j].UpdatedAt)
	})

	// ok final step to create the order data
	testResults := map[string][]entity.ObservationResult{}
	for _, observation := range specimen.ObservationResult {
		// TODO check whether this will create chaos in order or not
		testResults[observation.Code] = append(testResults[observation.Code], observation)
	}

	// fill the placeholder
	for i, test := range tests {
		newTest := test
		history := testResults[test.Test]
		if len(history) > 0 {
			newTest = newTest.FromObservationResult(history[0])
		}
		newTest = newTest.FillHistory(history)

		// or should be like this or we can just use the above code
		tests[i] = newTest
	}

	return tests
}
