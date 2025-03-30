package analyzer

import (
	"context"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
)

// ProcessQuery is a usecase for processing query message.
func (u *Usecase) ProcessQuery(ctx context.Context) error {
	// do something
	// find work order with WorkOrderStatusNew
	workOrders, err := u.WorkOrderRepository.FindByStatus(ctx, entity.WorkOrderStatusNew)
	if err != nil {
		return err
	}

	for i := range workOrders {
		ba400.NewOML_O33(workOrders[i].Patient, workOrders[i].Device, false)
	}
	return nil
}

// ProcessQueryAll is a usecase for processing query all message.
func (u *Usecase) ProcessQueryAll(ctx context.Context, barcode string) error {
	// do something
	// find work order with WorkOrderStatusNew and barcode
	return nil
}
