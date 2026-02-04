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
	qcentryrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/qc_entry"
	qcresultrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/qc_result"
	qualitycontrolrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/quality_control"
	rolerepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/role"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/specimen"
	subcategoryrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/sub_category"
	summaryrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/summary"
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
	simgosuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/simgos"
	simrsuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/simrs"
	observation_requestuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/observation_request"
	patientuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/patient"
	quality_control_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/quality_control"
	resultUC "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/result"
	role_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/role"
	specimenuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/specimen"
	summary_uc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/summary"
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
	provideSimrsUsecase,
	wire.Bind(new(cron.SIMRSUsecase), new(*simrsuc.Usecase)),
	provideSimgosUsecase,
	wire.Bind(new(cron.SIMGOSUsecase), new(*simgosuc.Usecase)),
	externaluc.NewUsecase,
	provideLicenseService,
	summary_uc.NewSummaryUsecase,
	quality_control_uc.NewQualityControlUsecase,
	wire.Bind(new(usecase.QualityControl), new(*quality_control_uc.QualityControlUsecase)),
)

var restRepositorySet = wire.NewSet(
	testTypeRepo.NewRepository,
	observation_result.NewRepository,
	observation_request.NewRepository,
	device.NewDeviceRepository,
	qualitycontrolrepo.NewQualityControlRepository,
	wire.Bind(new(repository.QualityControl), new(*qualitycontrolrepo.QualityControlRepository)),
	qcentryrepo.NewQCEntryRepository,
	wire.Bind(new(repository.QCEntry), new(*qcentryrepo.QCEntryRepository)),
	qcresultrepo.NewQCResultRepository,
	wire.Bind(new(repository.QCResult), new(*qcresultrepo.QCResultRepository)),
	daily_sequence.NewRepository,
	patientrepo.NewPatientRepository,
	workOrderrepo.NewWorkOrderRepository,
	specimen.NewRepository,
	configrepo.NewRepository,
	test_template.NewRepository,
	adminrepo.NewAdminRepository,
	unit.NewRepository,
	rolerepo.NewRoleRepository,
	subcategoryrepo.NewRepository,
	hlsRepo.NewBa400,
	smbA15.NewA15,
	a15rest.NewA15,
	server.NewControllerRepository,
	provideKhanzaRepository,
	provideSimrsRepository,
	provideSimgosRepository,
	provideAllDevices,
	summaryrepo.NewSummaryRepository,
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
	middleware.NewIntegrationCheckMiddleware,
	provideIntegrationCheckConfig,
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
	rest.NewQCEntryHandler,
	rest.NewLogHandler,
	rest.NewExternalHandler,
	rest.NewKhanzaExternalHandler,
	rest.NewSimrsExternalHandler,
	rest.NewTechnoMedicHandler,
	rest.NewLicenseHandler,
	provideCronHandler,
	provideTechnoMedicUsecase,
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
		provideConfigUsecaseForCron,
		provideConfigCheckerForCron,
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
