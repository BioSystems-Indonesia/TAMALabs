package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
)

var (
	// TCPAppSet is a Wire provider set that provides a TCPServer.
	tcpAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		provideCache,
		hlsRepo.NewRepository,
		observation_result.NewRepository,
		observation_request.NewRepository,
		specimen.NewRepository,
		hlsUC.NewUsecase,
		tcp.NewHlSevenHandler,
		provideTCP,
		provideTCPHandler,
		provideTCPServer,
	)
)
