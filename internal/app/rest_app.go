package app

import (
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/cron"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/rest"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/serial/alifax"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/serial/cbs400"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/serial/coax"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/serial/diestro"
	ncc3300 "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/serial/ncc_3300"
	verifyu120 "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/serial/verifyU120"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/a15"
	analyxpanca "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/analyx_panca"
	analyxtrias "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/analyx_trias"
	ncc61 "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/neomedika_ncc61"
	swelabalfa "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/swelab_alfa"
	swelablumi "github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/swelab_lumi"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/tcp/wondfo"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/middleware"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/rest/a15rest"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/server"
	smbA15 "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/smb/A15"
	adminrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/admin"
	configrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/daily_sequence"
	device "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/device"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/observation_request"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/observation_result"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	rolerepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/role"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/specimen"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_template"
	testTypeRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/unit"
	workOrderrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	hlsRepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/tcp/ba400"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	admin_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/admin"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/analyzer"
	auth_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/auth"
	barcodeGeneratorUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/barcode_generator"
	configuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/config"
	deviceuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/device"
	externaluc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external"
	khanzauc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/khanza"
	simrsuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/simrs"
	observation_requestuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/observation_request"
	patientuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/patient"
	resultUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/result"
	role_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/role"
	specimenuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/specimen"
	test_template_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/test_template"
	testTypeUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/test_type"
	unitUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/unit"
	workOrderuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order/runner"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order/runner/postrun"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order/runner/prerun"
	"github.com/google/wire"
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
	simrsuc.NewUsecase,
	wire.Bind(new(cron.SIMRSUsecase), new(*simrsuc.Usecase)),
	externaluc.NewUsecase,
	provideLicenseService,
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
	provideSimrsRepository,
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
	wondfo.NewHandler,
	cbs400.NewHandler,
	verifyu120.NewHandler,
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
	rest.NewLicenseHandler,
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
