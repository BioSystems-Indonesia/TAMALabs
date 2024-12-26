package sql

import (
	"context"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type TestType interface {
	FindAll(ctx context.Context, req *entity.TestTypeGetManyRequest) ([]entity.TestType, error)
}
