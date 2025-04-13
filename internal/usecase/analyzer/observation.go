package analyzer

import (
	"context"
	"errors"

	"log/slog"

	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/internal/util"
)

// ProcessOULR22 processes the OUL_R22 message.
func (u *Usecase) ProcessOULR22(ctx context.Context, data entity.OUL_R22) error {
	specimens := data.Specimens
	var errs []error
	uniqueWorkOrder := map[int64]struct{}{}
	for i := range specimens {
		spEntity, err := u.SpecimenRepository.FindByBarcode(ctx, specimens[i].Barcode)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		if specimens[i].ObservationResult != nil {
			for j := range specimens[i].ObservationResult {
				specimens[i].ObservationResult[j].SpecimenID = int64(spEntity.ID)
			}
			err := u.ObservationResultRepository.CreateMany(ctx, specimens[i].ObservationResult)
			if err != nil {
				errs = append(errs, err)
				continue
			}
		}

		_, ok := uniqueWorkOrder[spEntity.WorkOrder.ID]
		if !ok {
			workOrder := spEntity.WorkOrder
			workOrder.Status = entity.WorkOrderStatusCompleted
			err := u.WorkOrderRepository.Update(&workOrder)
			if err != nil {
				errs = append(errs, err)
				continue
			}

			uniqueWorkOrder[spEntity.WorkOrder.ID] = struct{}{}
		}
	}

	err := errors.Join(errs...)
	if err != nil {
		specimenIDs := util.Map(specimens, func(s entity.Specimen) int {
			return s.ID
		})
		slog.Error("error processing OUL_R22", "error", err, "specimen", specimenIDs)
	}

	return nil
}

// ProcessQBPQ11 processes the QBP_Q11 message.
func (u *Usecase) ProcessQBPQ11(ctx context.Context, data entity.QBP_Q11) ([]h251.OML_O33, error) {
	if data.QPD.Barcode != "" {
		omlO33, err := u.withBarcode(ctx, data.QPD.Barcode)
		if err != nil {
			return nil, err
		}
		return []h251.OML_O33{omlO33}, nil
	}
	return u.withoutBarcode(ctx)
}

// withBarcode processes the QBP_Q11 message with a barcode.
func (u *Usecase) withBarcode(ctx context.Context, barcode string) (h251.OML_O33, error) {
	speciment, err := u.SpecimenRepository.FindByBarcode(ctx, barcode)
	if err != nil {
		return h251.OML_O33{}, err
	}

	device, err := u.DeviceRepository.FindByID(ctx, 1)
	if err != nil {
		return h251.OML_O33{}, err
	}

	o := ba400.NewOML_O33(speciment.Patient, device, false)
	if err != nil {
		return h251.OML_O33{}, err
	}

	return o, nil
}

// withoutBarcode processes the QBP_Q11 message without a barcode.
func (u *Usecase) withoutBarcode(ctx context.Context) ([]h251.OML_O33, error) {
	device, err := u.DeviceRepository.FindByID(ctx, 1)
	if err != nil {
		return nil, err
	}
	// find work order with status WorkOrderStatusNew
	workOrders, err := u.WorkOrderRepository.FindByStatus(ctx, entity.WorkOrderStatusNew)
	if err != nil {
		return nil, err
	}

	if len(workOrders) == 0 {
		return nil, errors.New("no work order found")
	}

	var omlO33s []h251.OML_O33

	for _, workOrder := range workOrders {
		o := ba400.NewOML_O33(workOrder.Patient, device, false)
		omlO33s = append(omlO33s, o)
	}

	return omlO33s, nil
}
