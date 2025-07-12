package app

import (
	"github.com/google/wire"
	ncc3300 "github.com/oibacidem/lims-hl-seven/internal/delivery/serial/ncc_3300"
	devicerepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
	observation_request "github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	observation_result "github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	specimen "github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"
)

var serialNCC3300UsecaseSet = wire.NewSet(
	analyzer.NewUsecase,
	wire.Bind(new(usecase.Analyzer), new(*analyzer.Usecase)),
)

var serialNCC3300RepositorySet = wire.NewSet(
	observation_result.NewRepository,
	observation_request.NewRepository,
	specimen.NewRepository,
	workOrderrepo.NewWorkOrderRepository,
	devicerepo.NewDeviceRepository,
)

var serialNCC3300AppSet = wire.NewSet(
	provideDB,
	provideConfig,
	provideCache,
	serialNCC3300RepositorySet,
	serialNCC3300UsecaseSet,
	ncc3300.NewHandler,
)
