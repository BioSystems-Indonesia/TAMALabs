package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	observation_requestuc "github.com/oibacidem/lims-hl-seven/internal/usecase/observation_request"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	specimenuc "github.com/oibacidem/lims-hl-seven/internal/usecase/specimen"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
)

var (
	// RestAppSet is a Wire provider set that provides a RestServer.
	restAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		provideCache,
		hlsRepo.NewRepository,
		patientrepo.NewPatientRepository,
		workOrderrepo.NewWorkOrderRepository,
		observation_result.NewRepository,
		observation_request.NewRepository,
		specimen.NewRepository,
		hlsUC.NewUsecase,
		patientuc.NewPatientUseCase,
		specimenuc.NewSpecimenUseCase,
		workOrderuc.NewWorkOrderUseCase,
		observation_requestuc.NewObservationRequestUseCase,
		rest.NewHlSevenHandler,
		rest.NewHealthCheckHandler,
		rest.NewPatientHandler,
		rest.NewSpecimenHandler,
		rest.NewWorkOrderHandler,
		rest.NewFeatureListHandler,
		rest.NewObservationRequestHandler,
		rest.NewDeviceHandler,
		provideTCP,
		provideRestHandler,
		provideRestServer,
	)
)
