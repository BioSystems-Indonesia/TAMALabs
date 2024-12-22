package observation_requestuc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
)

type ObservationRequestUseCase struct {
	cfg                    *config.Schema
	ObservationRequestRepo *observation_request.Repository
	validate               *validator.Validate
}

func NewObservationRequestUseCase(
	cfg *config.Schema,
	ObservationRequestRepo *observation_request.Repository,
	validate *validator.Validate,
) *ObservationRequestUseCase {
	return &ObservationRequestUseCase{cfg: cfg, ObservationRequestRepo: ObservationRequestRepo, validate: validate}
}

func (p ObservationRequestUseCase) FindAll(
	ctx context.Context, req *entity.ObservationRequestGetManyRequest,
) ([]entity.ObservationRequest, error) {
	return p.ObservationRequestRepo.FindAll(ctx, req)
}

func (p ObservationRequestUseCase) FindOneByID(id int64) (entity.ObservationRequest, error) {
	return p.ObservationRequestRepo.FindOne(id)
}
