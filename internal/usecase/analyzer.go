package usecase

import (
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// Analyzer is an interface for Analyzer usecase
type Analyzer interface {
	ProcessOULR22(data entity.OUL_R22) error
}
