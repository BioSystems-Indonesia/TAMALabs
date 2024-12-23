package sql

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type Specimen interface {
	Create(ctx context.Context, data *entity.Specimen) error
}
