package postrun

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	workOrderrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
)

type RunAction struct {
	workOrderRepo *workOrderrepo.WorkOrderRepository
}

func NewRunAction(
	workOrderRepo *workOrderrepo.WorkOrderRepository,
) *RunAction {
	return &RunAction{workOrderRepo: workOrderRepo}
}

func (w RunAction) PostRun(ctx context.Context, req *entity.WorkOrderRunRequest) error {
	var errs []error
	for _, workOrderID := range req.WorkOrderIDs {
		errUpdate := updateWorkOrderStatus(ctx, w.workOrderRepo, workOrderID, req.GetDevice(), entity.WorkOrderStatusPending)
		if errUpdate != nil {
			errs = append(errs, errUpdate)
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("failed to update work order status %w", errors.Join(errs...))
	}

	return nil
}
