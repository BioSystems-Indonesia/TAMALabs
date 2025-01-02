package result

import (
	"context"
	"fmt"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	specimenRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
)

type Usecase struct {
	resultRepository    repository.Result
	workOrderRepository *workOrderRepo.WorkOrderRepository
	specimenRepository  *specimenRepo.Repository
}

func NewUsecase(
	resultRepository repository.Result,
	workOrderRepository *workOrderRepo.WorkOrderRepository,
	specimenRepository *specimenRepo.Repository,
) *Usecase {
	return &Usecase{
		resultRepository:    resultRepository,
		workOrderRepository: workOrderRepository,
		specimenRepository:  specimenRepository,
	}
}

func (u *Usecase) Results(ctx context.Context, req *entity.ResultGetManyRequest) ([]entity.Result, error) {
	// it should be one table that handle the result of the patient with multiple work orders
	worksOrders, err := u.workOrderRepository.FindByStatus(ctx, entity.WorkOrderStatusNew)
	if err != nil {
		return nil, err
	}

	return u.mapListResult(worksOrders), nil
}

func (u *Usecase) ResultDetail(ctx context.Context, barcode string) (entity.Result, error) {
	// this is bad design because if patient have multiple visit this will query the rest patient visit result

	result := entity.Result{
		Barcode: barcode,
	}

	specimens, err := u.specimenRepository.FindByBarcode(ctx, barcode)
	if err != nil {
		return result, err
	}

	result.Detail = u.processResultDetail(specimens)

	return result, nil
}

func (u *Usecase) mapListResult(worksOrders []entity.WorkOrder) []entity.Result {
	var results []entity.Result
	for _, workOrder := range worksOrders {
		for _, patient := range workOrder.Patient {
			if u.isPatientAlreadyExist(patient.ID, results) {
				continue
			}
			results = append(results, entity.Result{
				ID:          patient.Specimen[0].Barcode,
				Date:        workOrder.CreatedAt,
				Barcode:     patient.Specimen[0].Barcode,
				PatientName: patient.FirstName + " " + patient.LastName,
				PatientID:   patient.ID,
			})
		}
	}
	return results
}

func (u *Usecase) processResultDetail(specimens []entity.Specimen) entity.ResultDetail {
	var resultDetail entity.ResultDetail
	for _, specimen := range specimens {
		for _, observation := range specimen.ObservationResult {
			resultTest := entity.ResultTest{
				Test:           observation.TestType.Name,
				Result:         observation.Values[0],
				Unit:           observation.TestType.Unit,
				Category:       observation.TestType.Category,
				ReferenceRange: fmt.Sprintf("%v - %v", observation.TestType.LowRefRange, observation.TestType.HighRefRange),
			}
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

			switch observation.TestType.Category {
			case "Hematology":
				resultDetail.Hematology = append(resultDetail.Hematology, resultTest)
			case "Biochemistry":
				resultDetail.Biochemistry = append(resultDetail.Biochemistry, resultTest)
			case "Observation":
				resultDetail.Observation = append(resultDetail.Observation, resultTest)
			}
		}
	}
	return resultDetail
}

func (u *Usecase) isPatientAlreadyExist(patientID int64, results []entity.Result) bool {
	for _, result := range results {
		if result.PatientID == patientID {
			return true
		}
	}
	return false
}
