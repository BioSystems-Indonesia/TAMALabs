package result

import (
	"context"
	"fmt"
	"strconv"

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

	return resp, nil
}

func (u *Usecase) ResultDetail(ctx context.Context, specimenID int64) (entity.ResultDetail, error) {
	specimen, err := u.specimenRepository.FindOne(ctx, specimenID)
	if err != nil {
		return entity.ResultDetail{}, err
	}

	return entity.ResultDetail{
		Specimen:   specimen,
		TestResult: u.processResultDetail(specimen),
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

func (u *Usecase) processResultDetail(specimen entity.Specimen) map[string][]entity.ResultTest {
	var resultTestsCategory = map[string][]entity.ResultTest{}
	for _, observation := range specimen.ObservationResult {
		resultTest := entity.ResultTest{}.FromObservationResult(observation)
		// only process the first value, if the observation have multiple values we need to handle it later
		result, err := strconv.ParseFloat(resultTest.Result, 64)
		if err != nil {
			continue
		}

		resultTest.Abnormal = entity.NormalResult
		if result <= observation.TestType.LowRefRange {
			resultTest.Abnormal = entity.LowResult
		} else if result >= observation.TestType.HighRefRange {
			resultTest.Abnormal = entity.HighResult
		}

		categoryName := observation.TestType.Category
		testResults, ok := resultTestsCategory[categoryName]
		if ok {
			resultTestsCategory[categoryName] = append(testResults, resultTest)
		} else {
			resultTestsCategory[categoryName] = []entity.ResultTest{resultTest}
		}
	}

	return resultTestsCategory
}

func (u *Usecase) DeleteResult(context context.Context, id int64) (entity.ObservationResult, error) {
	return u.resultRepository.Delete(context, id)
}
