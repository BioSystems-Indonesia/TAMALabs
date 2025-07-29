package result

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
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
) (entity.PaginationResponse[entity.WorkOrder], error) {
	resp, err := u.workOrderRepository.FindAllForResult(ctx, req)
	if err != nil {
		return entity.PaginationResponse[entity.WorkOrder]{}, err
	}

	for i := range resp.Data {
		u.fillResultDetail(&resp.Data[i], true)
	}

	return resp, nil
}

func (u *Usecase) ResultDetail(ctx context.Context, workOrderID int64) (entity.ResultDetail, error) {
	workOrder, err := u.workOrderRepository.FindOneForResult(workOrderID)
	if err != nil {
		return entity.ResultDetail{}, err
	}
	u.fillResultDetail(&workOrder, false)

	// Why inverted? preve is get from next and next from get
	// because the default of ordering in front end are from the latest
	// maybe we need to config this..
	prevID, err := u.workOrderRepository.FindNextID(ctx, workOrderID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	nextID, err := u.workOrderRepository.FindPrevID(ctx, workOrderID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
	}

	return entity.ResultDetail{
		WorkOrder:  workOrder,
		TestResult: u.groupResultInCategory(workOrder.TestResult),
		NextID:     prevID,
		PrevID:     nextID,
	}, nil
}

// PutTestResult will create ObservationResult
// set the value, unit and everyting else
// and will prepend it in TestResult.History
func (u *Usecase) PutTestResult(ctx context.Context, result entity.TestResult) (entity.TestResult, error) {
	oldResult := result

	obs := entity.ObservationResult{
		SpecimenID:     result.SpecimenID,
		TestCode:       result.Test,
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

	obs, err = u.TooglePickTestResult(ctx, obs.ID)
	if err != nil {
		return result, fmt.Errorf("failed to toogle pick test result: %w", err)
	}

	obs.TestType, err = u.testTypeRepository.FindOneByCode(ctx, obs.TestCode)
	if err != nil {
		slog.Info("cannot fill test type for result", "id", obs.ID, "error", err)
	}

	result = result.FromObservationResult(obs)
	// Hack so the front end is not add first then replace
	// TODO maybe need to find better way than this
	if oldResult.ID != 0 {
		result.ID = oldResult.ID
	}

	history, err := u.resultRepository.FindHistory(ctx, obs)
	if err != nil {
		slog.Info("cannot get history", "error", err)
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

func (u *Usecase) fillResultDetail(workOrder *entity.WorkOrder, hideEmpty bool) {
	var allObservationRequests []entity.ObservationRequest
	var allObservationResults []entity.ObservationResult
	for _, specimen := range workOrder.Specimen {
		allObservationRequests = append(allObservationRequests, specimen.ObservationRequest...)
		allObservationResults = append(allObservationResults, specimen.ObservationResult...)
	}

	allTests := make([]entity.TestResult, len(allObservationRequests))
	// create the placeholder first
	for i, request := range allObservationRequests {
		allTests[i] = entity.TestResult{}.CreateEmpty(request)
	}

	// sort by test code
	sort.Slice(allTests, func(i, j int) bool {
		return allTests[i].Test < allTests[j].Test
	})

	// prepare the data that will be filled into the placeholder
	// prepare by grouping into code. But before that, sort by updated_at
	// the latest updated_at will be the first element
	sort.Slice(allObservationResults, func(i, j int) bool {
		return allObservationResults[i].CreatedAt.After(allObservationResults[j].CreatedAt)
	})

	// ok final step to create the order data
	testResults := map[string][]entity.ObservationResult{}
	for _, observation := range allObservationResults {
		// TODO check whether this will create chaos in order or not
		testResults[observation.TestCode] = append(testResults[observation.TestCode], observation)
	}

	// fill the placeholder
	totalResultFilled := 0
	for i, test := range allTests {
		newTest := test
		history := testResults[test.Test]
		if len(history) > 0 {
			// Pick the latest history or the manually picked one
			pickedTest := history[0]
			for _, v := range history {
				if v.Picked {
					pickedTest = v
					break
				}
			}

			newTest = newTest.FromObservationResult(pickedTest)
		}
		newTest = newTest.FillHistory(history)

		// or should be like this or we can just use the above code
		allTests[i] = newTest

		// count the filled result
		if newTest.Result != nil {
			totalResultFilled++
		}
	}

	workOrder.TotalRequest = int64(len(allObservationRequests))
	workOrder.TotalResultFilled = int64(totalResultFilled)
	workOrder.HaveCompleteData = len(allObservationRequests) == totalResultFilled
	if len(allObservationRequests) != 0 {
		workOrder.PercentComplete = float64(totalResultFilled) / float64(len(allObservationRequests))
	}

	if hideEmpty {
		var filteredTests []entity.TestResult
		for _, test := range allTests {
			if test.Result == nil || *test.Result == 0 {
				continue
			}
			filteredTests = append(filteredTests, test)
		}
		allTests = filteredTests
	}

	workOrder.TestResult = allTests
}

func (u *Usecase) TooglePickTestResult(ctx context.Context, testResultID int64) (entity.ObservationResult, error) {
	return u.resultRepository.PickObservationResult(ctx, testResultID)
}

func (u *Usecase) ApproveResult(context context.Context, workOrderID int64, adminID int64) error {
	workOrder, err := u.workOrderRepository.FindOne(workOrderID)
	if err != nil {
		return err
	}

	if !slices.Contains(workOrder.DoctorIDs, adminID) {
		return entity.ErrForbidden
	}

	return u.resultRepository.ApproveResult(context, workOrderID)
}

func (u *Usecase) RejectResult(context context.Context, workOrderID int64, adminID int64) error {
	workOrder, err := u.workOrderRepository.FindOne(workOrderID)
	if err != nil {
		return err
	}

	if !slices.Contains(workOrder.DoctorIDs, adminID) {
		return entity.ErrForbidden
	}

	return u.resultRepository.RejectResult(context, workOrderID)
}
