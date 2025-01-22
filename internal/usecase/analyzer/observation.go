package analyzer

import (
	"context"
	"errors"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/util"

	"log/slog"
)

func (u *Usecase) ProcessOULR22(ctx context.Context, data entity.OUL_R22) error {
	specimens := data.Specimens
	var errs []error
	uniqueWorkOrder := map[int64]struct{}{}
	for i := range specimens {
		spEntities, err := u.SpecimenRepository.FindByBarcode(ctx, specimens[i].HL7ID)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		var spEntity entity.Specimen

		if len(spEntities) > 0 {
			spEntity = spEntities[0]
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
