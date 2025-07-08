package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	analyxpanca "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_panca"
	analyxtrias "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_trias"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	"github.com/oibacidem/lims-hl-seven/internal/middleware"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	a15 "github.com/oibacidem/lims-hl-seven/internal/repository/smb/A15"
	adminrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/admin"
	configrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/config"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/daily_sequence"
	device "github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	rolerepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/role"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_template"
	testTypeRepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/test_type"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/unit"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/server"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	admin_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/admin"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	auth_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/auth"
	barcodeGeneratorUC "github.com/oibacidem/lims-hl-seven/internal/usecase/barcode_generator"
	configuc "github.com/oibacidem/lims-hl-seven/internal/usecase/config"
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
	observation_requestuc "github.com/oibacidem/lims-hl-seven/internal/usecase/observation_request"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	resultUC "github.com/oibacidem/lims-hl-seven/internal/usecase/result"
	role_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/role"
	specimenuc "github.com/oibacidem/lims-hl-seven/internal/usecase/specimen"
	test_template_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/test_template"
	testTypeUC "github.com/oibacidem/lims-hl-seven/internal/usecase/test_type"
	unitUC "github.com/oibacidem/lims-hl-seven/internal/usecase/unit"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/work_order/runner"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/work_order/runner/postrun"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/work_order/runner/prerun"
)

var restUsecaseSet = wire.NewSet(
	testTypeUC.NewUsecase,
	resultUC.NewUsecase,
	barcodeGeneratorUC.NewUsecase,
	wire.Bind(new(usecase.BarcodeGenerator), new(*barcodeGeneratorUC.Usecase)),
	patientuc.NewPatientUseCase,
	specimenuc.NewSpecimenUseCase,
	prerun.NewRunAction,
	prerun.NewCancelAction,
	postrun.NewCancelAction,
	postrun.NewRunAction,
	postrun.NewIncompleteSendAction,
	runner.NewStrategy,
	workOrderuc.NewWorkOrderUseCase,
	observation_requestuc.NewObservationRequestUseCase,
	configuc.NewConfigUseCase,
	test_template_uc.NewUsecase,
	admin_uc.NewAdminUsecase,
	auth_uc.NewAuthUseCase,
	unitUC.NewUnitUseCase,
	deviceuc.NewDeviceUseCase,
	role_uc.NewRoleUsecase,
	analyzer.NewUsecase,
	wire.Bind(new(usecase.Analyzer), new(*analyzer.Usecase)),
)

var restRepositorySet = wire.NewSet(
	testTypeRepo.NewRepository,
	observation_result.NewRepository,
	observation_request.NewRepository,
	device.NewDeviceRepository,
	daily_sequence.NewRepository,
	patientrepo.NewPatientRepository,
	workOrderrepo.NewWorkOrderRepository,
	specimen.NewRepository,
	configrepo.NewRepository,
	test_template.NewRepository,
	adminrepo.NewAdminRepository,
	unit.NewRepository,
	rolerepo.NewRoleRepository,
	hlsRepo.NewBa400,
	a15.NewA15,
	server.NewTCPServerRepository,
	provideAllDevices,
)

var tcpHandlerSet = wire.NewSet(
	tcp.NewHlSevenHandler,
	analyxtrias.NewHandler,
	analyxpanca.NewHandler,
	swelabalfa.NewHandler,
	tcp.NewDeviceStrategy,
	wire.Bind(new(repository.DeviceTCPHandlerStrategy), new(*tcp.DeviceStrategy)),
)

var restMiddlewareSet = wire.NewSet(
	middleware.NewJWTMiddleware,
)

var restHandlerSet = wire.NewSet(
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
	rest.NewAdminHandler,
	rest.NewAuthHandler,
	rest.NewUnitHandler,
	rest.NewRoleHandler,
	rest.NewDeviceHandler,
	rest.NewTestTemplateHandler,
	rest.NewServerControllerHandler,
	rest.NewLogHandler,
)

var (
	// RestAppSet is a Wire provider set that provides a RestServer.
	restAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		provideConfig,
		provideCache,
		restRepositorySet,
		restUsecaseSet,
		restMiddlewareSet,
		restHandlerSet,
		tcpHandlerSet,
		provideRestHandler,
		provideRestServer,
	)
)
