package postrun

import (
	"context"
	"errors"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	workOrderrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
)

type IncompleteSendAction struct {
	workOrderRepo *workOrderrepo.WorkOrderRepository
}

func NewIncompleteSendAction(
	workOrderRepo *workOrderrepo.WorkOrderRepository,
) *IncompleteSendAction {
	return &IncompleteSendAction{workOrderRepo: workOrderRepo}
}

func (w IncompleteSendAction) PostRun(ctx context.Context, req *entity.WorkOrderRunRequest) error {
	var errs []error
	for _, workOrderID := range req.WorkOrderIDs {
		errUpdate := updateWorkOrderStatus(ctx, w.workOrderRepo, workOrderID, req.GetDevice(), entity.WorkOrderStatusIncompleteSend)
		if errUpdate != nil {
			errs = append(errs, errUpdate)
		}
	}
	if len(errs) != 0 {
		return fmt.Errorf("failed to update work order status %w", errors.Join(errs...))
	}

	return nil
}
