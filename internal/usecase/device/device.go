package deviceuc

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/device"
)

type DeviceUseCase struct {
	deviceRepo *device.Repository
}

func NewDeviceUseCase(repo *device.Repository) *DeviceUseCase {
	return &DeviceUseCase{
		deviceRepo: repo,
	}
}

func (u *DeviceUseCase) FindByID(ctx context.Context, id int64) (entity.Device, error) {
	device, err := u.deviceRepo.FindByID(ctx, id)
	if err != nil {
		return entity.Device{}, err
	}

	return device, nil
}
