package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_result"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	BA400Repository              *ba400.Repository
	SpecimenRepository           *specimen.Repository
	ObservationResultRepository  *observation_result.Repository
	ObservationRequestRepository *observation_request.Repository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	ba400Repository *ba400.Repository,
	specimenRepository *specimen.Repository,
	observationResultRepository *observation_result.Repository,
	observationRequestRepository *observation_request.Repository,
) *Usecase {
	return &Usecase{
		BA400Repository:              ba400Repository,
		SpecimenRepository:           specimenRepository,
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
	}
}
