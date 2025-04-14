package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	device "github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	analyzerUC "github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
)

var tcpUsecaseSet = wire.NewSet(
	analyzerUC.NewUsecase,
	wire.Bind(new(usecase.Analyzer), new(*analyzerUC.Usecase)),
)

var tcpRepositorySet = wire.NewSet(
	observation_result.NewRepository,
	observation_request.NewRepository,
	device.NewRepository,
)

var (
	// TCPAppSet is a Wire provider set that provides a TCPServer.
	tcpAppSet = wire.NewSet(
		provideValidator,
		provideDB,
		provideConfig,
		provideCache,

		tcpRepositorySet,
		specimen.NewRepository,
		workOrderrepo.NewWorkOrderRepository,

		tcpUsecaseSet,

		tcp.NewHlSevenHandler,
		provideTCP,
		provideTCPServer,
	)
)
