package analyzer

import (
	devicerepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ObservationResultRepository  *observation_result.Repository
	ObservationRequestRepository *observation_request.Repository
	SpecimenRepository           *specimen.Repository
	WorkOrderRepository          *workOrderrepo.WorkOrderRepository
	DeviceRepository             *devicerepo.DeviceRepository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	observationResultRepository *observation_result.Repository,
	observationRequestRepository *observation_request.Repository,
	specimenRepository *specimen.Repository,
	workOrderRepository *workOrderrepo.WorkOrderRepository,
	deviceRepository *devicerepo.DeviceRepository,
) *Usecase {
	return &Usecase{
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
		SpecimenRepository:           specimenRepository,
		WorkOrderRepository:          workOrderRepository,
		DeviceRepository:             deviceRepository,
	}
}
