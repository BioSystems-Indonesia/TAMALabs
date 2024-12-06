package app

import (
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

func provideTCP(config *config.Schema) *hl_seven.TCP {
	tcpEr := hl_seven.NewTCP(config)

	return tcpEr
}

func provideRest(config *config.Schema, handlers *rest.Handler) server.RestServer {
	serv := server.NewRest("8080")
	rest.RegisterRoutes(serv.GetClient(), handlers)
	return serv
}

func provideHandler(
	hlSevenHandler *rest.HlSevenHandler,
) *rest.Handler {
	return &rest.Handler{
		HlSevenHandler: hlSevenHandler,
	}
}
