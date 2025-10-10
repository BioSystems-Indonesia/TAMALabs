package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
	"go.bug.st/serial"
)

type DeviceServer struct {
	device entity.Device
	server server.Controller
}

type ControllerRepository struct {
	handlerStrategy repository.DeviceServerStrategy

	serverDeviceMap map[int]DeviceServer
}

func NewControllerRepository(
	handlerStrategy repository.DeviceServerStrategy,
	allDevices []entity.Device,
) *ControllerRepository {
	t := &ControllerRepository{
		serverDeviceMap: map[int]DeviceServer{},
		handlerStrategy: handlerStrategy,
	}

	for _, device := range allDevices {
		_, err := t.StartNewServer(context.Background(), device)
		if err != nil {
			slog.Error("failed to start new server for device", "device_id", device.ID, "error", err)
		}
	}

	return t
}

func (r *ControllerRepository) StartNewServer(
	ctx context.Context,
	device entity.Device,
) (server.Controller, error) {
	sd, ok := r.serverDeviceMap[device.ID]
	if ok {
		if !r.needRestartServer(sd, device) {
			return sd.server, nil
		}

		err := sd.server.Stop()
		if err != nil {
			slog.Error("failed to stop server for device", "device_id", device.ID, "error", err)
		}
	}

	deviceServer, err := r.startServer(device)
	if err != nil {
		return nil, fmt.Errorf("failed to start new server for device %d: %w", device.ID, err)
	}
	r.serverDeviceMap[device.ID] = deviceServer

	return deviceServer.server, nil
}

func (r *ControllerRepository) needRestartServer(serverDevice DeviceServer, newDevice entity.Device) bool {
	if serverDevice.server.State() == constant.ServerStateStopped {
		return true
	}

	if serverDevice.device.ReceivePort != newDevice.ReceivePort {
		return true
	}

	return false
}

func (r *ControllerRepository) startServer(device entity.Device) (DeviceServer, error) {
	s, err := r.handlerStrategy.ChooseDeviceServer(device)
	if err != nil {
		if errors.Is(err, entity.ErrDeviceTypeNotSupport) {
			return DeviceServer{
				server: NewDummyServer(),
				device: device,
			}, nil
		}

		return DeviceServer{}, fmt.Errorf("failed to choose device handler for device %d: %w", device.ID, err)
	}

	errStart := s.Start()
	if errStart != nil {
		return DeviceServer{}, fmt.Errorf("failed to start server for device %d: %w", device.ID, errStart)
	}
	go s.Serve()

	return DeviceServer{
		server: s,
		device: device,
	}, nil
}

func (r *ControllerRepository) GetAllServerState() map[int]constant.ServerState {
	mapState := make(map[int]constant.ServerState)
	for deviceID, sd := range r.serverDeviceMap {
		mapState[deviceID] = sd.server.State()
	}
	return mapState
}

func (r *ControllerRepository) GetServerStateByDeviceID(deviceID int) constant.ServerState {
	sd, ok := r.serverDeviceMap[deviceID]
	if !ok {
		return constant.ServerStateStopped
	}
	return sd.server.State()
}

func (r *ControllerRepository) StopServerByDeviceID(ctx context.Context, deviceID int) error {
	sd, ok := r.serverDeviceMap[deviceID]
	if !ok {
		return fmt.Errorf("server for device %d not found", deviceID)
	}
	return sd.server.Stop()
}

func (r *ControllerRepository) DeleteServerByDeviceID(ctx context.Context, deviceID int) error {
	_, ok := r.serverDeviceMap[deviceID]
	if !ok {
		return fmt.Errorf("server for device %d not found", deviceID)
	}

	err := r.StopServerByDeviceID(ctx, deviceID)
	if err != nil {
		slog.Error("failed to stop server for device", "device_id", deviceID, "error", err)
	}

	delete(r.serverDeviceMap, deviceID)
	return nil
}

func (r *ControllerRepository) GetAllSerialPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	return ports, nil
}
