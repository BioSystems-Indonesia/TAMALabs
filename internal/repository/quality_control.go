package repository

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

type QualityControl interface {
	Create(ctx context.Context, qc *entity.QualityControl) error
	GetMany(ctx context.Context, req entity.GetManyRequestQualityControl) ([]entity.QualityControl, int64, error)
	GetByID(ctx context.Context, id int) (*entity.QualityControl, error)
	GetStatistics(ctx context.Context, deviceID int) (map[string]interface{}, error)
}
