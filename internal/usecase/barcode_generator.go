package usecase

import (
	"context"
)

type BarcodeGenerator interface {
	NextOrderBarcode(ctx context.Context) (string, error)
}
