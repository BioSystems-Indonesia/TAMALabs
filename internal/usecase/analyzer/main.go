package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ObservationResultRepository  repository.ObservationResult
	ObservationRequestRepository repository.ObservationRequest
	SpecimenRepository           *specimen.Repository
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	observationResultRepository repository.ObservationResult,
	observationRequestRepository repository.ObservationRequest,
	specimenRepository *specimen.Repository,
) *Usecase {
	return &Usecase{
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
		SpecimenRepository:           specimenRepository,
	}
}
