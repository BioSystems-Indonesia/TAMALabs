package analyzer

import (
	devicerepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/device"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/observation_request"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/observation_result"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/specimen"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrderrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ObservationResultRepository  *observation_result.Repository
	ObservationRequestRepository *observation_request.Repository
	SpecimenRepository           *specimen.Repository
	WorkOrderRepository          *workOrderrepo.WorkOrderRepository
	DeviceRepository             *devicerepo.DeviceRepository
	TestTypeRepository           *test_type.Repository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	observationResultRepository *observation_result.Repository,
	observationRequestRepository *observation_request.Repository,
	specimenRepository *specimen.Repository,
	workOrderRepository *workOrderrepo.WorkOrderRepository,
	deviceRepository *devicerepo.DeviceRepository,
	testTypeRepository *test_type.Repository,
) *Usecase {
	return &Usecase{
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
		SpecimenRepository:           specimenRepository,
		WorkOrderRepository:          workOrderRepository,
		DeviceRepository:             deviceRepository,
		TestTypeRepository:           testTypeRepository,
	}
}
