package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
)

type tcpDeviceServer struct {
	device entity.Device
	server repository.TCPServerController
}

type TCPServer struct {
	handlerStrategy repository.DeviceTCPHandlerStrategy

	serverDeviceMap map[int]tcpDeviceServer
}

func NewTCPServerRepository(
	handlerStrategy repository.DeviceTCPHandlerStrategy,
	allDevices []entity.Device,
) *TCPServer {
	t := &TCPServer{
		serverDeviceMap: map[int]tcpDeviceServer{},
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

func (r *TCPServer) StartNewServer(
	ctx context.Context,
	device entity.Device,
) (repository.TCPServerController, error) {
	sd, ok := r.serverDeviceMap[device.ID]
	if ok {
		if !r.needRestartServer(sd, device) {
			return sd.server, nil
		}

		err := sd.server.Stop()
		if err != nil {
			return nil, fmt.Errorf("failed to stop server for device %d: %w", device.ID, err)
		}
	}

	deviceServer, err := r.startServer(device)
	if err != nil {
		return nil, fmt.Errorf("failed to start new server for device %d: %w", device.ID, err)
	}
	r.serverDeviceMap[device.ID] = deviceServer

	return deviceServer.server, nil
}

func (r *TCPServer) needRestartServer(serverDevice tcpDeviceServer, newDevice entity.Device) bool {
	if serverDevice.server.State() == constant.ServerStateStopped {
		return true
	}

	if serverDevice.device.ReceivePort != newDevice.ReceivePort {
		return true
	}

	return false
}

func (r *TCPServer) startServer(device entity.Device) (tcpDeviceServer, error) {
	h, err := r.handlerStrategy.ChooseDeviceHandler(device)
	if err != nil {
		if errors.Is(err, entity.ErrDeviceTypeNotSupport) {
			return tcpDeviceServer{
				server: NewDummyServer(),
				device: device,
			}, nil
		}

		return tcpDeviceServer{}, fmt.Errorf("failed to choose device handler for device %d: %w", device.ID, err)
	}

	port := strconv.Itoa(device.ReceivePort)
	if port == "" {
		return tcpDeviceServer{}, fmt.Errorf("device %d has no receive port", device.ID)
	}

	server := server.NewTCP(port)
	server.SetHandler(h)

	errStart := server.Start()
	if errStart != nil {
		return tcpDeviceServer{}, fmt.Errorf("failed to start server for device %d: %w", device.ID, errStart)
	}
	go server.Serve()

	return tcpDeviceServer{
		server: server,
		device: device,
	}, nil
}

func (r *TCPServer) GetAllServerState() map[int]constant.ServerState {
	mapState := make(map[int]constant.ServerState)
	for deviceID, sd := range r.serverDeviceMap {
		mapState[deviceID] = sd.server.State()
	}
	return mapState
}

func (r *TCPServer) GetServerStateByDeviceID(deviceID int) constant.ServerState {
	sd, ok := r.serverDeviceMap[deviceID]
	if !ok {
		return constant.ServerStateStopped
	}
	return sd.server.State()
}

func (r *TCPServer) StopServerByDeviceID(ctx context.Context, deviceID int) error {
	sd, ok := r.serverDeviceMap[deviceID]
	if !ok {
		return fmt.Errorf("server for device %d not found", deviceID)
	}
	return sd.server.Stop()
}

func (r *TCPServer) DeleteServerByDeviceID(ctx context.Context, deviceID int) error {
	_, ok := r.serverDeviceMap[deviceID]
	if !ok {
		return fmt.Errorf("server for device %d not found", deviceID)
	}

	err := r.StopServerByDeviceID(ctx, deviceID)
	if err != nil {
		return fmt.Errorf("failed to stop server for device %d: %w", deviceID, err)
	}

	delete(r.serverDeviceMap, deviceID)
	return nil
}
