package postrun

import (
	"context"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	workOrderrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
)

func updateWorkOrderStatus(
	ctx context.Context,
	repo *workOrderrepo.WorkOrderRepository,
	workOrderID int64,
	device entity.Device,
	status entity.WorkOrderStatus,
) error {
	workOrder, err := repo.FindOne(workOrderID)
	if err != nil {
		return fmt.Errorf("failed to p.FindOneByID %w", err)
	}

	workOrder.Status = status
	err = repo.UpsertDevice(workOrderID, int64(device.ID))
	if err != nil {
		return fmt.Errorf("failed to p.UpsertDevice %w", err)
	}

	err = repo.Update(&workOrder)
	if err != nil {
		return fmt.Errorf("failed to p.Update %w", err)
	}
	return nil
}
