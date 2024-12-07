package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	hlsRepo "github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
	hlsUC "github.com/oibacidem/lims-hl-seven/internal/usecase/hl_seven"
)

var (
	// TCPAppSet is a Wire provider set that provides a TCPServer.
	tcpAppSet = wire.NewSet(
		hlsRepo.NewRepository,
		hlsUC.NewUsecase,
		tcp.NewHlSevenHandler,
		provideTCP,
		provideTCPHandler,
		provideTCPServer,
	)
)
