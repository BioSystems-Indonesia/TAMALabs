package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	BA400Repository       *ba400.Repository
	SpecimenRepository    *specimen.Repository
	ObservationRepository *observation.Repository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	ba400Repository *ba400.Repository,
	specimenRepository *specimen.Repository,
	observationRepository *observation.Repository,
) *Usecase {
	return &Usecase{
		BA400Repository:       ba400Repository,
		SpecimenRepository:    specimenRepository,
		ObservationRepository: observationRepository,
	}
}
