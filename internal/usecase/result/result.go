package result

import (
	"context"
	"fmt"
	"sort"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	specimenRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	workOrderRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/util"
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

func (u *Usecase) UpdateResult(ctx context.Context, data []entity.ResultTest) ([]entity.ResultTest, error) {
	testResults, err := u.resultRepository.UpdateResultTest(ctx, data)
	if err != nil {
		return []entity.ResultTest{}, err
	}

	return testResults, nil
}

func (u *Usecase) CreateResult(
	ctx context.Context, data entity.ObservationResultCreate,
) (entity.ObservationResult, error) {
	testType, err := u.testTypeRepository.FindOneByCode(ctx, data.Code)
	if err != nil {
		return entity.ObservationResult{}, err
	}

	input := entity.ObservationResult{
		SpecimenID:     data.SpecimenID,
		Code:           data.Code,
		Values:         []string{fmt.Sprintf("%v", data.Value)},
		Unit:           testType.Unit,
		ReferenceRange: fmt.Sprintf("%.2f - %.2f", testType.LowRefRange, testType.HighRefRange),
		TestType:       testType,
	}

	err = u.resultRepository.Create(ctx, &input)
	if err != nil {
		return entity.ObservationResult{}, err
	}

	return input, nil
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

func (u *Usecase) groupResultInCategory(tests []entity.ResultTest) map[string][]entity.ResultTest {
	var resultTestsCategory = map[string][]entity.ResultTest{}
	for _, resultTest := range tests {
		categoryName := resultTest.Category

		categoryTestResults, ok := resultTestsCategory[categoryName]
		if ok {
			resultTestsCategory[categoryName] = append(categoryTestResults, resultTest)
		} else {
			resultTestsCategory[categoryName] = []entity.ResultTest{resultTest}
		}
	}

	return resultTestsCategory
}

func (u *Usecase) processResultDetail(specimen entity.Specimen) []entity.ResultTest {
	var resultTestsCode = map[string][]entity.ResultTest{}
	var tests = make([]entity.ResultTest, len(specimen.ObservationResult))
	for i, observation := range specimen.ObservationResult {
		resultTest := entity.ResultTest{}.FromObservationResult(observation)

		code := resultTest.Test
		codeTestResults, ok := resultTestsCode[code]
		if ok {
			resultTestsCode[code] = append(codeTestResults, resultTest)
		} else {
			resultTestsCode[code] = []entity.ResultTest{resultTest}
		}

		tests[i] = resultTest
	}

	// Sort by updated_at
	sort.Slice(tests, func(i, j int) bool {
		if tests[i].Test != tests[j].Test {
			return tests[i].Test < tests[j].Test
		}

		return tests[i].UpdatedAt > tests[j].UpdatedAt
	})

	// Remove duplicates code, only get the first one
	tests = util.RemoveDuplicatesFromStruct(tests, func(test entity.ResultTest) string {
		return test.Test
	})

	// Add history to each test
	for i := range tests {
		history := resultTestsCode[tests[i].Test]
		sort.Slice(history, func(i, j int) bool {
			return history[i].UpdatedAt > history[j].UpdatedAt
		})
		tests[i].History = history
	}

	return tests
}

func (u *Usecase) DeleteResult(context context.Context, id int64) (entity.ObservationResult, error) {
	return u.resultRepository.Delete(context, id)
}

func (u *Usecase) DeleteResultBulk(context context.Context, req *entity.DeleteResultBulkReq) (entity.ObservationResult, error) {
	return u.resultRepository.DeleteBulk(context, req.IDs)
}
