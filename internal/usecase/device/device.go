package deviceuc

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	serverrepo "github.com/oibacidem/lims-hl-seven/internal/repository/server"
	a15 "github.com/oibacidem/lims-hl-seven/internal/repository/smb/A15"
	devicerepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/work_order/runner"
	"github.com/oibacidem/lims-hl-seven/internal/util"
)

type DeviceUseCase struct {
	cfg            *config.Schema
	deviceRepo     *devicerepo.DeviceRepository
	runnerStrategy *runner.Strategy
	ba400          *ba400.Ba400
	a15            *a15.A15
	serverRepo     *serverrepo.ControllerRepository
}

func NewDeviceUseCase(
	cfg *config.Schema,
	deviceRepo *devicerepo.DeviceRepository,
	runnerStrategy *runner.Strategy,
	ba400 *ba400.Ba400,
	a15 *a15.A15,
	serverRepo *serverrepo.ControllerRepository,
) *DeviceUseCase {
	return &DeviceUseCase{cfg: cfg, deviceRepo: deviceRepo, runnerStrategy: runnerStrategy, ba400: ba400, a15: a15, serverRepo: serverRepo}
}

func (p DeviceUseCase) FindAll(
	ctx context.Context, req *entity.GetManyRequestDevice,
) (entity.PaginationResponse[entity.Device], error) {
	return p.deviceRepo.FindAll(ctx, req)
}

func (p DeviceUseCase) FindOneByID(ctx context.Context, id int64) (entity.Device, error) {
	return p.deviceRepo.FindOne(id)
}

func (p DeviceUseCase) Create(ctx context.Context, req *entity.Device) error {
	if req.ReceivePort == "" {
		req.ReceivePort = strconv.Itoa(p.RandomPort())
	}

	device, err := p.deviceRepo.FindOneByReceivePort(req.ReceivePort)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return err
	}

	if device.ID != 0 && device.ID != req.ID {
		return entity.NewUserError(entity.UserErrorDeviceAlreadyExistsReceivePort,
			fmt.Sprintf("device already exists with receive port %d must be unique", req.ReceivePort))
	}

	err = p.deviceRepo.Create(req)
	if err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	_, err = p.serverRepo.StartNewServer(ctx, *req)
	if err != nil {
		return fmt.Errorf("failed to start new server: %w", err)
	}

	return nil
}

func (p DeviceUseCase) Update(ctx context.Context, req *entity.Device) error {
	if req.ReceivePort == "" {
		req.ReceivePort = strconv.Itoa(p.RandomPort())
	}

	device, err := p.deviceRepo.FindOneByReceivePort(req.ReceivePort)
	if err != nil && !errors.Is(err, entity.ErrNotFound) {
		return err
	}

	if device.ID != 0 && device.ID != req.ID {
		return entity.NewUserError(entity.UserErrorDeviceAlreadyExistsReceivePort,
			fmt.Sprintf("device already exists with receive port %s must be unique", req.ReceivePort))
	}

	err = p.deviceRepo.Update(req)
	if err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	_, err = p.serverRepo.StartNewServer(ctx, *req)
	if err != nil {
		return fmt.Errorf("failed to start new server: %w", err)
	}

	return nil
}

func (p DeviceUseCase) Delete(ctx context.Context, id int) error {
	err := p.deviceRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return p.serverRepo.StopServerByDeviceID(ctx, id)
}

func (p *DeviceUseCase) GetDeviceConnection(ctx context.Context, id int) entity.DeviceConnectionResponse {
	device, err := p.deviceRepo.FindOne(int64(id))
	if err != nil {
		return p.buildDeviceNotFoundResponse(id)
	}

	senderErrChan := make(chan error, 1)
	receiverErrChan := make(chan error, 1)

	go p.checkSenderConnection(ctx, device, senderErrChan)
	go p.checkReceiverConnection(ctx, device, receiverErrChan)

	senderErr, receiverErr := p.waitForResults(ctx, senderErrChan, receiverErrChan)

	return p.buildConnectionResponse(id, senderErr, receiverErr)
}

func (p *DeviceUseCase) buildDeviceNotFoundResponse(id int) entity.DeviceConnectionResponse {
	return entity.DeviceConnectionResponse{
		DeviceID: id,
		Sender: entity.DeviceConnectionStatusResponse{
			Message: "Device not found",
			Status:  entity.DeviceConnectionStatusDisconnected,
		},
		Receiver: entity.DeviceConnectionStatusResponse{
			Message: "Device not found",
			Status:  entity.DeviceConnectionStatusDisconnected,
		},
	}
}

func (p *DeviceUseCase) checkSenderConnection(ctx context.Context, device entity.Device, errChan chan<- error) {
	defer close(errChan)

	sender, err := p.ChooseDeviceSender(ctx, device)
	if err != nil {
		errChan <- err
		return
	}

	errChan <- sender.CheckConnection(ctx, device)
}

func (p *DeviceUseCase) checkReceiverConnection(
	ctx context.Context,
	device entity.Device,
	errChan chan<- error,
) {
	defer close(errChan)

	state := p.serverRepo.GetServerStateByDeviceID(device.ID)
	switch state {
	case constant.ServerStateConnect:
		errChan <- nil
	case constant.ServerStateServing:
		errChan <- entity.ErrDeviceNotConnected
	case constant.ServerStateNoServer:
		errChan <- entity.ErrDeviceTypeNotSupport
	case constant.ServerStateStopped:
		errChan <- errors.New("server is not serving")
	default:
		errChan <- errors.New("server is in unknown state")
	}
}

//nolint:nonamedreturns // much more readable this way
func (p *DeviceUseCase) waitForResults(
	ctx context.Context,
	senderChan, receiverChan <-chan error,
) (senderErr, receiverErr error) {
	var senderReceived, receiverReceived bool

	// Wait for both results or timeout
	for !senderReceived || !receiverReceived {
		select {
		case err, ok := <-senderChan:
			if ok {
				senderErr = err
				senderReceived = true
			}
		case err, ok := <-receiverChan:
			if ok {
				receiverErr = err
				receiverReceived = true
			}
		case <-ctx.Done():
			// Handle timeout
			if !senderReceived {
				senderErr = errors.New("sender check timeout")
			}
			if !receiverReceived {
				receiverErr = errors.New("receiver check timeout")
			}
			return senderErr, receiverErr
		}
	}

	return senderErr, receiverErr
}

func (p *DeviceUseCase) buildConnectionResponse(id int, senderErr, receiverErr error) entity.DeviceConnectionResponse {
	response := entity.DeviceConnectionResponse{DeviceID: id}

	response.Sender = p.buildStatusResponse(senderErr)
	response.Receiver = p.buildStatusResponse(receiverErr)

	return response
}

func (p *DeviceUseCase) buildStatusResponse(err error) entity.DeviceConnectionStatusResponse {
	if err == nil {
		return entity.DeviceConnectionStatusResponse{
			Message: "Connection successful",
			Status:  entity.DeviceConnectionStatusConnected,
		}
	}

	if errors.Is(err, entity.ErrDeviceNotConnected) {
		return entity.DeviceConnectionStatusResponse{
			Message: "Device not connected to LIS, standby mode",
			Status:  entity.DeviceConnectionStatusStandby,
		}
	}

	if errors.Is(err, entity.ErrDeviceTypeNotSupport) {
		return entity.DeviceConnectionStatusResponse{
			Message: "This device does not support this feature",
			Status:  entity.DeviceConnectionStatusNotSupported,
		}
	}

	return entity.DeviceConnectionStatusResponse{
		Message: err.Error(),
		Status:  entity.DeviceConnectionStatusDisconnected,
	}
}

func (p *DeviceUseCase) ChooseDeviceSender(ctx context.Context, device entity.Device) (usecase.DeviceSender, error) {
	switch device.Type {
	case entity.DeviceTypeBA400, entity.DeviceTypeBA200, entity.DeviceTypeOther:
		return p.ba400, nil
	case entity.DeviceTypeA15:
		return p.a15, nil
	default:
		return nil, entity.ErrDeviceTypeNotSupport
	}
}

func (p *DeviceUseCase) RandomPort() int {
	minPortNumber := 1024
	maxPortNumber := 10000

	const maxAttempts = 10
	for i := 0; i < maxAttempts; i++ {
		port := util.RandomNumber(minPortNumber, maxPortNumber)
		if !entity.IsImportantWindowsPort(port) {
			return port
		}
	}

	return 0
}
