package usecase

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

type WorkOrderPreRunner interface {
	PreRun(ctx context.Context, req *entity.WorkOrderRunRequest) error
}

type WorkOrderPostRunner interface {
	PostRun(ctx context.Context, req *entity.WorkOrderRunRequest) error
}
