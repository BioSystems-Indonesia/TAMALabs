package usecase

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

type DeviceSender interface {
	Send(ctx context.Context, req *entity.SendPayloadRequest) error
	CheckConnection(ctx context.Context, device entity.Device) error
}
