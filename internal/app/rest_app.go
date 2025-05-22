package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/middleware"
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
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

var restUsecaseSet = wire.NewSet(
	testTypeUC.NewUsecase,
	resultUC.NewUsecase,
	barcodeGeneratorUC.NewUsecase,
	wire.Bind(new(usecase.BarcodeGenerator), new(*barcodeGeneratorUC.Usecase)),
)

var restRepositorySet = wire.NewSet(
	testTypeRepo.NewRepository,
	observation_result.NewRepository,
	observation_request.NewRepository,
	device.NewRepository,
	daily_sequence.NewRepository,
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
		adminrepo.NewAdminRepository,
		unit.NewRepository,
		rolerepo.NewRoleRepository,

		restUsecaseSet,
		tcpUsecaseSet,

		patientuc.NewPatientUseCase,
		specimenuc.NewSpecimenUseCase,
		workOrderuc.NewWorkOrderUseCase,
		observation_requestuc.NewObservationRequestUseCase,
		configuc.NewConfigUseCase,
		test_template_uc.NewUsecase,
		admin_uc.NewAdminUsecase,
		auth_uc.NewAuthUseCase,
		unitUC.NewUnitUseCase,
		deviceuc.NewDeviceUseCase,
		role_uc.NewRoleUsecase,

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
		wire.Struct(new(rest.DeviceHandler), "*"),

		middleware.NewJWTMiddleware,

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
