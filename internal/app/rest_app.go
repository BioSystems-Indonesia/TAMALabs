package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	specimentrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/speciment"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/hl_seven"
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
