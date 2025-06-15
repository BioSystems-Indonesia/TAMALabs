package usecase

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type WorkOrderPreRunner interface {
	PreRun(ctx context.Context, req *entity.WorkOrderRunRequest) error
}

type WorkOrderDeviceSender interface {
	Send(ctx context.Context, req *entity.SendPayloadRequest) error
	CheckConnection(ctx context.Context, device entity.Device) error
}

type WorkOrderPostRunner interface {
	PostRun(ctx context.Context, req *entity.WorkOrderRunRequest) error
}
