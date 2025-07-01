package usecase

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type DeviceSender interface {
	Send(ctx context.Context, req *entity.SendPayloadRequest) error
	CheckConnection(ctx context.Context, device entity.Device) error
}
