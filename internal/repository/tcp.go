package repository

import (
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type TCPServerController interface {
	SetPort(port string)
	Start() error
	State() constant.ServerState
	Serve()
	Stop() error
}

type DeviceTCPHandlerStrategy interface {
	ChooseDeviceHandler(device entity.Device) (server.TCPHandler, error)
}
