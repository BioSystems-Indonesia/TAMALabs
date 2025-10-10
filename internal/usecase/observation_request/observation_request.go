package observation_requestuc

import (
	"context"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/observation_request"
	"github.com/go-playground/validator/v10"
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
) (entity.PaginationResponse[entity.ObservationRequest], error) {
	return p.ObservationRequestRepo.FindAll(ctx, req)
}

func (p ObservationRequestUseCase) FindOneByID(ctx context.Context, id int64) (entity.ObservationRequest, error) {
	return p.ObservationRequestRepo.FindOne(ctx, id)
}

func (p ObservationRequestUseCase) BulkUpdate(ctx context.Context, request []entity.ObservationRequest) error {
	return p.ObservationRequestRepo.BulkUpdate(ctx, request)
}
