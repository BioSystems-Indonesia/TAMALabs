package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ObservationResultRepository  repository.ObservationResult
	ObservationRequestRepository repository.ObservationRequest
	SpecimenRepository           *specimen.Repository
	WorkOrderRepository          *workOrderrepo.WorkOrderRepository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	observationResultRepository repository.ObservationResult,
	observationRequestRepository repository.ObservationRequest,
	specimenRepository *specimen.Repository,
	workOrderRepository *workOrderrepo.WorkOrderRepository,
) *Usecase {
	return &Usecase{
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
		SpecimenRepository:           specimenRepository,
		WorkOrderRepository:          workOrderRepository,
	}
}
