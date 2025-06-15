package deviceuc

import (
	"context"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	devicerepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/work_order/runner"
)

type DeviceUseCase struct {
	cfg            *config.Schema
	deviceRepo     *devicerepo.DeviceRepository
	runnerStrategy *runner.Strategy
}

func NewDeviceUseCase(
	cfg *config.Schema,
	deviceRepo *devicerepo.DeviceRepository,
	runnerStrategy *runner.Strategy,
) *DeviceUseCase {
	return &DeviceUseCase{cfg: cfg, deviceRepo: deviceRepo, runnerStrategy: runnerStrategy}
}

func (p DeviceUseCase) FindAll(
	ctx context.Context, req *entity.GetManyRequestDevice,
) (entity.PaginationResponse[entity.Device], error) {
	return p.deviceRepo.FindAll(ctx, req)
}

func (p DeviceUseCase) FindOneByID(ctx context.Context, id int64) (entity.Device, error) {
	return p.deviceRepo.FindOne(id)
}

func (p DeviceUseCase) Create(req *entity.Device) error {
	return p.deviceRepo.Create(req)
}

func (p DeviceUseCase) Update(req *entity.Device) error {
	return p.deviceRepo.Update(req)
}

func (p DeviceUseCase) Delete(id int) error {
	return p.deviceRepo.Delete(id)
}

func (p *DeviceUseCase) GetDeviceConnection(ctx context.Context, id int) error {
	device, err := p.deviceRepo.FindOne(int64(id))
	if err != nil {
		return err
	}

	sender, err := p.runnerStrategy.ChooseSendRunner(ctx, device)
	if err != nil {
		return fmt.Errorf("failed to p.runnerStrategy.ChooseSendRunner %w", err)
	}

	return sender.CheckConnection(ctx, device)
}
