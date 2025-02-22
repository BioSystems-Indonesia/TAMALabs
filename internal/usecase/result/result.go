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

// UpdateResult will store the result into histories
// in the handler, it will use PUT
func (u *Usecase) UpdateResult(ctx context.Context, data []entity.TestResult) ([]entity.TestResult, error) {
	// This will create N+1 condition but right now only used for one test
	// TODO but do we really need to update in bulk?
	var testResults = []entity.TestResult{}
	for _, v := range data {
		testType, err := u.testTypeRepository.FindOneByID(ctx, int(v.TestTypeID))
		if err != nil {
			return data, err
		}

		values := make([]string, 0)
		if v.Result != nil {
			values = append(values, fmt.Sprintf("%f", *v.Result))
		}

		input := entity.ObservationResult{
			SpecimenID:     v.SpecimenID,
			Code:           testType.Code,
			Values:         values,
			Unit:           testType.Unit,
			ReferenceRange: fmt.Sprintf("%.2f - %.2f", testType.LowRefRange, testType.HighRefRange),
			TestType:       testType,
		}

		err = u.resultRepository.Create(ctx, &input)
		if err != nil {
			return data, err
		}

		test := entity.TestResult{}
		test = test.FromObservationResult(input)
		// Hack so the front end is not add first then replace
		// TODO maybe need to find better way than this
		if v.ID != 0 {
			test.ID = v.ID
		}

		history, err := u.resultRepository.FindHistory(ctx, input)
		if err != nil {
			log.Printf("cannot get history: ")
		}

		test = test.FillHistory(history)

		testResults = append(testResults, test)
	}

	return testResults, nil
}

func (u *Usecase) CreateResult(
	ctx context.Context, data entity.ObservationResultCreate,
) ([]entity.ObservationResult, error) {
	var results []entity.ObservationResult
	for _, v := range data.Tests {
		testType, err := u.testTypeRepository.FindOneByID(ctx, int(v.TestTypeID))
		if err != nil {
			return []entity.ObservationResult{}, err
		}

		input := entity.ObservationResult{
			SpecimenID:     data.SpecimenID,
			Code:           testType.Code,
			Values:         []string{fmt.Sprintf("%v", v.Value)},
			Unit:           testType.Unit,
			ReferenceRange: fmt.Sprintf("%.2f - %.2f", testType.LowRefRange, testType.HighRefRange),
			TestType:       testType,
		}

		results = append(results, input)
	}

	err := u.resultRepository.CreateMany(ctx, results)
	if err != nil {
		return []entity.ObservationResult{}, err
	}

	return results, nil
}

func (u *Usecase) mapListResult(specimens []entity.Specimen) []entity.Result {
	results := make([]entity.Result, len(specimens))
	for i, specimen := range specimens {
		results[i] = entity.Result{
			Specimen: specimen,
		}
	}

	return results
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
		history := testResults[test.Test]
		if len(history) > 0 {
			test = test.FromObservationResult(history[0])
		}
		test = test.FillHistory(history)

		// or should be like this or we can just use the above code
		tests[i] = test
	}

	return tests
}

func (u *Usecase) DeleteResult(context context.Context, id int64) (entity.ObservationResult, error) {
	return u.resultRepository.Delete(context, id)
}

func (u *Usecase) DeleteResultBulk(context context.Context, req *entity.DeleteResultBulkReq) (entity.ObservationResult, error) {
	return u.resultRepository.DeleteBulk(context, req.IDs)
}
