package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	patientrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/patient"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/hl_seven"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
)

var (
	// RestAppSet is a Wire provider set that provides a RestServer.
	restAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		hlsRepo.NewRepository,
		patientrepo.NewPatientRepository,
		hlsUC.NewUsecase,
		patientuc.NewPatientUseCase,
		rest.NewHlSevenHandler,
		rest.NewHealthCheckHandler,
		rest.NewPatientHandler,
		provideTCP,
		provideHandler,
		provideRest,
	)
)
