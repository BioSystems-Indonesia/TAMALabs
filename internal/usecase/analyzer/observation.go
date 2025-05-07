package analyzer

import (
	"context"
	"errors"

	"log/slog"

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
func (u *Usecase) ProcessQBPQ11(ctx context.Context, data entity.QBP_Q11) error {
	if data.QPD.Barcode != "" {
		return u.withBarcode(ctx, data.QPD.Barcode)
	}
	return u.withoutBarcode(ctx)
}

// withBarcode processes the QBP_Q11 message with a barcode.
func (u *Usecase) withBarcode(ctx context.Context, barcode string) error {
	speciment, err := u.SpecimenRepository.FindByBarcode(ctx, barcode)
	if err != nil {
		return err
	}

	device, err := u.DeviceRepository.FindByID(ctx, 1)
	if err != nil {
		return err
	}

	err = ba400.SendToBA400(ctx, &entity.SendPayloadRequest{
		Patients: []entity.Patient{speciment.Patient},
		Device:   device,
		Urgent:   false,
	})
	if err != nil {
		return err
	}

	return nil
}

// withoutBarcode processes the QBP_Q11 message without a barcode.
func (u *Usecase) withoutBarcode(ctx context.Context) error {
	device, err := u.DeviceRepository.FindByID(ctx, 1)
	if err != nil {
		return err
	}
	// find work order with status WorkOrderStatusNew
	workOrders, err := u.WorkOrderRepository.FindByStatus(ctx, entity.WorkOrderStatusNew)
	if err != nil {
		return err
	}

	if len(workOrders) == 0 {
		return errors.New("no work order found")
	}

	for _, workOrder := range workOrders {
		ba400.SendToBA400(ctx, &entity.SendPayloadRequest{
			Patients: []entity.Patient{workOrder.Patient},
			Device:   device,
			Urgent:   false,
		})
	}

	return nil
}
