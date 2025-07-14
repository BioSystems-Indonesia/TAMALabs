package repository

import (
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type DeviceServerStrategy interface {
	ChooseDeviceServer(device entity.Device) (server.Controller, error)
}

type DeviceTCPHandlerStrategy interface {
	ChooseDeviceTCPHandler(device entity.Device) (server.TCPHandler, error)
}

type DeviceSerialHandlerStrategy interface {
	ChooseDeviceSerialHandler(device entity.Device) (server.SerialHandler, error)
}
