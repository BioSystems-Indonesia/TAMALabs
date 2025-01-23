package app

import "C"
import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	configrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/config"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	resultRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/result"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_template"
	testTypeRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	configuc "github.com/oibacidem/lims-hl-seven/internal/usecase/config"
	observation_requestuc "github.com/oibacidem/lims-hl-seven/internal/usecase/observation_request"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	resultUC "github.com/oibacidem/lims-hl-seven/internal/usecase/result"
	specimenuc "github.com/oibacidem/lims-hl-seven/internal/usecase/specimen"
	test_template_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/test_template"
	testTypeUC "github.com/oibacidem/lims-hl-seven/internal/usecase/test_type"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
)

var restUsecaseSet = wire.NewSet(
	testTypeUC.NewUsecase,
	resultUC.NewUsecase,
	wire.Bind(new(usecase.Result), new(*resultUC.Usecase)),
)

var restRepositorySet = wire.NewSet(
	testTypeRepo.NewRepository,
	observation_result.NewRepository,
	wire.Bind(new(repository.ObservationResult), new(*observation_result.Repository)),
	observation_request.NewRepository,
	wire.Bind(new(repository.ObservationRequest), new(*observation_request.Repository)),
	resultRepo.NewRepository,
	wire.Bind(new(repository.Result), new(*resultRepo.Repository)),
)

var (
	// RestAppSet is a Wire provider set that provides a RestServer.
	restAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		provideConfig,
		provideCache,

		restRepositorySet,
		hlsRepo.NewRepository,
		patientrepo.NewPatientRepository,
		workOrderrepo.NewWorkOrderRepository,
		specimen.NewRepository,
		configrepo.NewRepository,
		test_template.NewRepository,

		restUsecaseSet,
		hlsUC.NewUsecase,
		patientuc.NewPatientUseCase,
		specimenuc.NewSpecimenUseCase,
		workOrderuc.NewWorkOrderUseCase,
		observation_requestuc.NewObservationRequestUseCase,
		configuc.NewConfigUseCase,
		test_template_uc.NewUsecase,

		rest.NewHlSevenHandler,
		rest.NewHealthCheckHandler,
		rest.NewPatientHandler,
		rest.NewSpecimenHandler,
		rest.NewWorkOrderHandler,
		rest.NewFeatureListHandler,
		rest.NewObservationRequestHandler,
		rest.NewDeviceHandler,
		rest.NewTestTypeHandler,
		rest.NewResultHandler,
		rest.NewConfigHandler,
		rest.NewTestTemplateHandler,
		provideTCP,
		provideRestHandler,
		provideRestServer,
	)
)
