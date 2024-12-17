package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	specimentrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/speciment"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	specimentuc "github.com/oibacidem/lims-hl-seven/internal/usecase/speciment"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
)

var (
	// RestAppSet is a Wire provider set that provides a RestServer.
	restAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		hlsRepo.NewRepository,
		patientrepo.NewPatientRepository,
		specimentrepo.NewSpecimentRepository,
		workOrderrepo.NewWorkOrderRepository,
		observation_request.NewRepository,
		observation.NewRepository,
		specimen.NewRepository,
		hlsUC.NewUsecase,
		patientuc.NewPatientUseCase,
		specimentuc.NewSpecimentUseCase,
		workOrderuc.NewWorkOrderUseCase,
		rest.NewHlSevenHandler,
		rest.NewHealthCheckHandler,
		rest.NewPatientHandler,
		rest.NewSpecimentHandler,
		rest.NewWorkOrderHandler,
		rest.NewFeatureListHandler,
		provideTCP,
		provideRestHandler,
		provideRestServer,
	)
)
