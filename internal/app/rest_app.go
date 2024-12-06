package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/hl_seven"
)

var (
	// RestAppSet is a Wire provider set that provides a RestServer.
	restAppSet = wire.NewSet(
		hlsRepo.NewRepository,
		hlsUC.NewUsecase,
		rest.NewHlSevenHandler,
		provideTCP,
		provideHandler,
		provideRest,
	)
)
