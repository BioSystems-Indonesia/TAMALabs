package app

import (
	"github.com/google/wire"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
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
	wire.Bind(new(repository.ObservationResult), new(*observation_result.Repository)),
	observation_request.NewRepository,
	wire.Bind(new(repository.ObservationRequest), new(*observation_request.Repository)),
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
		provideTCPHandler,
		provideTCPServer,
	)
)
