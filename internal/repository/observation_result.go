package repository

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type ObservationResult interface {
	Create(ctx context.Context, data *entity.ObservationResult) error
	CreateMany(ctx context.Context, data []entity.ObservationResult) error
}
