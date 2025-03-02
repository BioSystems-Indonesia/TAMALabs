package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	configrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/config"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_template"
	testTypeRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/unit"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	configuc "github.com/oibacidem/lims-hl-seven/internal/usecase/config"
	observation_requestuc "github.com/oibacidem/lims-hl-seven/internal/usecase/observation_request"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	resultUC "github.com/oibacidem/lims-hl-seven/internal/usecase/result"
	specimenuc "github.com/oibacidem/lims-hl-seven/internal/usecase/specimen"
	test_template_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/test_template"
	testTypeUC "github.com/oibacidem/lims-hl-seven/internal/usecase/test_type"
	unitUC "github.com/oibacidem/lims-hl-seven/internal/usecase/unit"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

var restUsecaseSet = wire.NewSet(
	testTypeUC.NewUsecase,
	resultUC.NewUsecase,
)

var restRepositorySet = wire.NewSet(
	testTypeRepo.NewRepository,
	observation_result.NewRepository,
	observation_request.NewRepository,
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
		unit.NewRepository,

		restUsecaseSet,
		tcpUsecaseSet,

		patientuc.NewPatientUseCase,
		specimenuc.NewSpecimenUseCase,
		workOrderuc.NewWorkOrderUseCase,
		observation_requestuc.NewObservationRequestUseCase,
		configuc.NewConfigUseCase,
		test_template_uc.NewUsecase,
		unitUC.NewUnitUseCase,

		rest.NewHlSevenHandler,
		rest.NewHealthCheckHandler,
		rest.NewPatientHandler,
		rest.NewSpecimenHandler,
		rest.NewWorkOrderHandler,
		rest.NewFeatureListHandler,
		rest.NewObservationRequestHandler,
		rest.NewTestTypeHandler,
		rest.NewResultHandler,
		rest.NewConfigHandler,
		rest.NewUnitHandler,
		wire.Struct(new(rest.DeviceHandler), "*"),

		rest.NewTestTemplateHandler,
		provideTCP,

		tcp.NewHlSevenHandler,
		provideTCPServer,
		wire.Bind(new(rest.TCPServerController), new(*server.TCP)),
		wire.Struct(new(rest.ServerControllerHandler), "*"),

		provideRestHandler,
		provideRestServer,
	)
)
