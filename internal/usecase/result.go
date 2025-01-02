package usecase

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type Result interface {
	Results(ctx context.Context, req *entity.ResultGetManyRequest) ([]entity.Result, error)
	ResultDetail(ctx context.Context, barcode string) (entity.Result, error)
}
