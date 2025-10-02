package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/cron"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/alifax"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/coax"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/serial/diestro"
	ncc3300 "github.com/oibacidem/lims-hl-seven/internal/delivery/serial/ncc_3300"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/a15"
	analyxpanca "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_panca"
	analyxtrias "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/analyx_trias"
	ncc61 "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/neomedika_ncc61"
	swelabalfa "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_alfa"
	swelablumi "github.com/oibacidem/lims-hl-seven/internal/delivery/tcp/swelab_lumi"
	"github.com/oibacidem/lims-hl-seven/internal/middleware"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/internal/repository/rest/a15rest"
	"github.com/oibacidem/lims-hl-seven/internal/repository/server"
	smbA15 "github.com/oibacidem/lims-hl-seven/internal/repository/smb/A15"
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
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	admin_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/admin"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	auth_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/auth"
	barcodeGeneratorUC "github.com/oibacidem/lims-hl-seven/internal/usecase/barcode_generator"
	configuc "github.com/oibacidem/lims-hl-seven/internal/usecase/config"
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
	externaluc "github.com/oibacidem/lims-hl-seven/internal/usecase/external"
	khanzauc "github.com/oibacidem/lims-hl-seven/internal/usecase/external/khanza"
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
	khanzauc.NewUsecase,
	externaluc.NewUsecase,
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
	smbA15.NewA15,
	a15rest.NewA15,
	server.NewControllerRepository,
	provideKhanzaRepository,
	provideAllDevices,
)

var tcpHandlerSet = wire.NewSet(
	a15.NewHandler,
	coax.NewHandler,
	diestro.NewHandler,
	ncc3300.NewHandler,
	alifax.NewHandler,
	tcp.NewHlSevenHandler,
	analyxtrias.NewHandler,
	analyxpanca.NewHandler,
	swelabalfa.NewHandler,
	swelablumi.NewHandler,
	ncc61.NewHandler,
	delivery.NewDeviceServerStrategy,
	wire.Bind(new(repository.DeviceServerStrategy), new(*delivery.DeviceServerStrategy)),
)

var restMiddlewareSet = wire.NewSet(
	middleware.NewJWTMiddleware,
)

var restHandlerSet = wire.NewSet(
	rest.NewHlSevenHandler,
	rest.NewHealthCheckHandler,
	rest.NewHealthHandler,
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
	rest.NewExternalHandler,
	rest.NewKhanzaExternalHandler,
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
		cron.NewCronHandler,
		cron.NewCronManager,
		provideRestHandler,
		provideRestServer,
	)

	canalHandlerSet = wire.NewSet(
		provideDB,
		provideConfig,
		provideValidator,
		provideCache,

		workOrderrepo.NewWorkOrderRepository,
		patientrepo.NewPatientRepository,
		testTypeRepo.NewRepository,
		observation_result.NewRepository,
		observation_request.NewRepository,
		specimen.NewRepository,
		daily_sequence.NewRepository,
		provideKhanzaRepository,

		resultUC.NewUsecase,
		barcodeGeneratorUC.NewUsecase,
		wire.Bind(new(usecase.BarcodeGenerator), new(*barcodeGeneratorUC.Usecase)),

		provideCanalHandler,
	)
)
