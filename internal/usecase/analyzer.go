package usecase

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// Analyzer is an interface for Analyzer usecase
type Analyzer interface {
	ProcessOULR22(ctx context.Context, data entity.OUL_R22) error
	ProcessQBPQ11(ctx context.Context, data entity.QBP_Q11) error
	ProcessORMO01(ctx context.Context, data entity.ORM_O01) ([]entity.Specimen, error)
	ProcessORUR01(ctx context.Context, data entity.ORU_R01) error
	ProcessCoax(ctx context.Context, data entity.CoaxTestResult) error
}
