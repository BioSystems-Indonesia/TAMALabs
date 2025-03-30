package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ObservationResultRepository  *observation_result.Repository
	ObservationRequestRepository *observation_request.Repository
	SpecimenRepository           *specimen.Repository
	WorkOrderRepository          *workOrderrepo.WorkOrderRepository
	BA400                        *ba400.Repository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	observationResultRepository *observation_result.Repository,
	observationRequestRepository *observation_request.Repository,
	specimenRepository *specimen.Repository,
	workOrderRepository *workOrderrepo.WorkOrderRepository,
	ba400 *ba400.Repository,
) *Usecase {
	return &Usecase{
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
		SpecimenRepository:           specimenRepository,
		WorkOrderRepository:          workOrderRepository,
		BA400:                        ba400,
	}
}
