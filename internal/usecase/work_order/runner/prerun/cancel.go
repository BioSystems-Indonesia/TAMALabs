package prerun

import (
	"context"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/observation_request"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
)

type CancelAction struct {
	observationRequestRepo *observation_request.Repository
	specimenRepo           *specimen.Repository
}

func NewCancelAction(
	observationRequestRepo *observation_request.Repository,
	specimenRepo *specimen.Repository,
) *CancelAction {
	return &CancelAction{observationRequestRepo: observationRequestRepo, specimenRepo: specimenRepo}
}

func (w CancelAction) PreRun(ctx context.Context, req *entity.WorkOrderRunRequest) error {
	specimens, err := w.specimenRepo.FindAllByWorkOrderIDs(ctx, req.WorkOrderIDs)
	if err != nil {
		return fmt.Errorf("error finding specimens: %w", err)
	}

	var observationRequest []entity.ObservationRequest
	for _, s := range specimens {
		observationRequest = append(observationRequest, s.ObservationRequest...)
	}
	for i := range observationRequest {
		observationRequest[i].ResultStatus = string(constant.ResultStatusDelete)
	}

	err = w.observationRequestRepo.BulkUpdate(ctx, observationRequest)
	if err != nil {
		return fmt.Errorf("observation request: %w", err)
	}

	return nil
}
