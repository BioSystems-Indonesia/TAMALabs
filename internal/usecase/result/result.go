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
		u.changeStatusIfNeeded(ctx, &resp.Data[i])
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
		// Try to find by alias_code if not found by code
		obs.TestType, err = u.testTypeRepository.FindOneByAliasCode(ctx, obs.TestCode)
		if err != nil {
			slog.Info("cannot fill test type for result", "id", obs.ID, "error", err)
		}
	}

	// Get specimen type for this observation result
	specimen, err := u.specimenRepository.FindOne(ctx, obs.SpecimenID)
	specimenType := ""
	if err != nil {
		slog.Info("cannot find specimen for result", "specimen_id", obs.SpecimenID, "error", err)
	} else {
		specimenType = specimen.Type
	}

	result = result.FromObservationResult(obs, specimenType)
	// Hack so the front end is not add first then replace
	// TODO maybe need to find better way than this
	if oldResult.ID != 0 {
		result.ID = oldResult.ID
	}

	history, err := u.resultRepository.FindHistory(ctx, obs)
	if err != nil {
		slog.Info("cannot get history", "error", err)
	}

	// Create specimen types map for history
	specimenTypes := make(map[int64]string)
	specimenTypes[obs.SpecimenID] = specimenType

	result = result.FillHistory(history, specimenTypes)

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

	// Create a map from specimen ID to specimen type for quick lookup
	specimenTypes := make(map[int64]string)

	for _, specimen := range workOrder.Specimen {
		specimenTypes[int64(specimen.ID)] = specimen.Type
		allObservationRequests = append(allObservationRequests, specimen.ObservationRequest...)
		allObservationResults = append(allObservationResults, specimen.ObservationResult...)
	}

	allTests := make([]entity.TestResult, len(allObservationRequests))
	// create the placeholder first
	for i, request := range allObservationRequests {
		// Find the corresponding specimen for this request
		var correspondingSpecimen entity.Specimen
		for _, specimen := range workOrder.Specimen {
			if int64(specimen.ID) == request.SpecimenID {
				correspondingSpecimen = specimen
				break
			}
		}
		allTests[i] = entity.TestResult{}.CreateEmpty(request, correspondingSpecimen)
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
		// Create a composite key: testCode + specimenType for more precise matching
		// Get specimen type for this observation
		specimenType := specimenTypes[observation.SpecimenID]
		compositeKey := fmt.Sprintf("%s:%s", observation.TestCode, specimenType)

		testResults[compositeKey] = append(testResults[compositeKey], observation)
	}

	// fill the placeholder
	totalResultFilled := 0
	for i, test := range allTests {
		newTest := test
		// Create composite key for matching: testCode + specimenType
		compositeKey := fmt.Sprintf("%s:%s", test.Test, test.SpecimenType)
		history := testResults[compositeKey]

		if len(history) > 0 {
			// Pick the latest history or the manually picked one
			pickedTest := history[0]
			for _, v := range history {
				if v.Picked {
					pickedTest = v
					break
				}
			}

			// Get specimen type for this observation result
			specimenType := specimenTypes[pickedTest.SpecimenID]
			newTest = newTest.FromObservationResult(pickedTest, specimenType)
		}
		newTest = newTest.FillHistory(history, specimenTypes)

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

func (u *Usecase) changeStatusIfNeeded(ctx context.Context, workOrder *entity.WorkOrder) {
	if !u.needToChangeStatus(workOrder) {
		return
	}

	workOrder.Status = entity.WorkOrderStatusCompleted

	go func(id int64) {
		err := u.workOrderRepository.ChangeStatus(ctx, id, entity.WorkOrderStatusCompleted)

		if err != nil {
			slog.ErrorContext(ctx, "failed to change status", "workOrderID", id, "error", err)
			return
		}
	}(workOrder.ID)
}

func (u *Usecase) needToChangeStatus(workOrder *entity.WorkOrder) bool {
	if workOrder.Status == entity.WorkOrderCancelled {
		return false
	}

	if workOrder.Status == entity.WorkOrderStatusCompleted {
		return false
	}

	if workOrder.PercentComplete < 1 {
		return false
	}

	return true
}
