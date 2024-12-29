package analyzer

import (
	"github.com/oibacidem/lims-hl-seven/internal/repository"
)

// Usecase is a struct handle HLSeven
type Usecase struct {
	ObservationResultRepository  repository.ObservationResult
	ObservationRequestRepository repository.ObservationRequest
}

// NewUsecase returns a new HLSeven
func NewUsecase(
	observationResultRepository repository.ObservationResult,
	observationRequestRepository repository.ObservationRequest,
) *Usecase {
	return &Usecase{
		ObservationResultRepository:  observationResultRepository,
		ObservationRequestRepository: observationRequestRepository,
	}
}
