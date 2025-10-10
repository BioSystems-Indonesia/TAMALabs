package repository

import (
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
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
