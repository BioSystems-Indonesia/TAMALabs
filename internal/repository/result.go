package repository

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type Result interface {
	FindAll(ctx context.Context) ([]entity.WorkOrder, error)
}
