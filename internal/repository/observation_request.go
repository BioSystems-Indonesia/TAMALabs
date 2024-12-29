package repository

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type ObservationRequest interface {
	Create(ctx context.Context, data *entity.ObservationRequest) error
	CreateMany(ctx context.Context, data []entity.ObservationRequest) error
	FindAll(ctx context.Context, req *entity.ObservationRequestGetManyRequest) ([]entity.ObservationRequest, error)
	FindOne(ctx context.Context, id int64) (entity.ObservationRequest, error)
}
