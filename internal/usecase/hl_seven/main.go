package hl_seven

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ORMRepository *hl_seven.Repository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	ORMRepository *hl_seven.Repository,
) *Usecase {
	return &Usecase{
		ORMRepository: ORMRepository,
	}
}
