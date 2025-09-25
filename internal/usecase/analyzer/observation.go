package analyzer

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

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

	device, err := u.DeviceRepository.FindOne(1)
	if err != nil {
		return err
	}

	err = ba400.NewBa400().Send(ctx, &entity.SendPayloadRequest{
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
	device, err := u.DeviceRepository.FindOne(1)
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
		err := ba400.NewBa400().Send(ctx, &entity.SendPayloadRequest{
			Patients: []entity.Patient{workOrder.Patient},
			Device:   device,
			Urgent:   false,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Usecase) ProcessORMO01(ctx context.Context, data entity.ORM_O01) ([]entity.Specimen, error) {
	var specimens []entity.Specimen
	var errs []error
	for _, o := range data.Orders {
		s, err := u.SpecimenRepository.FindByBarcode(ctx, o.Barcode)
		if err != nil {
			errs = append(errs, err)
		}
		specimens = append(specimens, s)
	}

	return specimens, errors.Join(errs...)
}

func (u *Usecase) ProcessORUR01(ctx context.Context, data entity.ORU_R01) error {
	var errs []error
	var specimens []entity.Specimen
	for _, p := range data.Patient {
		specimens = append(specimens, p.Specimen...)
	}

	uniqueWorkOrder := map[int64]struct{}{}
	for i, s := range specimens {
		spEntity, err := u.SpecimenRepository.FindByBarcode(ctx, s.Barcode)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		for j := range specimens[i].ObservationResult {
			specimens[i].ObservationResult[j].SpecimenID = int64(spEntity.ID)
		}
		err = u.ObservationResultRepository.CreateMany(ctx, specimens[i].ObservationResult)
		if err != nil {
			errs = append(errs, err)
			continue
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
		slog.Error("error processing ORU_R01", "error", err, "specimen", specimenIDs)
	}

	return nil
}

func (u *Usecase) ProcessCoax(ctx context.Context, data entity.CoaxTestResult) error {
	specimen, err := u.SpecimenRepository.FindByBarcode(ctx, data.DeviceID)
	if err != nil {
		return fmt.Errorf("specimen not found: %w", err)
	}

	observationResult := entity.ObservationResult{
		ID:             0,
		SpecimenID:     int64(specimen.ID),
		TestCode:       data.TestName,
		Description:    data.TestName,
		Values:         entity.JSONStringArray{data.Value},
		Type:           data.TestName,
		Unit:           data.Unit,
		ReferenceRange: data.Reference,
		Date:           u.parseCoaxDate(data.Date),
		AbnormalFlag:   entity.JSONStringArray{},
		Comments:       data.Flags,
		Picked:         false,
	}

	err = u.ObservationResultRepository.Create(ctx, &observationResult)
	if err != nil {
		return fmt.Errorf("error creating observation result: %w", err)
	}

	workOrder := specimen.WorkOrder
	workOrder.Status = entity.WorkOrderStatusCompleted
	err = u.WorkOrderRepository.Update(&workOrder)
	if err != nil {
		return fmt.Errorf("error updating work order: %w", err)
	}

	return nil
}

// paraseCoaxDate parse 2025/05/28 to 2025-05-28
func (u *Usecase) parseCoaxDate(date string) time.Time {
	date = strings.ReplaceAll(date, "/", "-")
	parsed, err := time.Parse(time.DateOnly, date)
	if err != nil {
		slog.Error("error parsing date", "error", err, "date", date)
		return time.Time{}
	}
	return parsed
}

func (u *Usecase) ProcessDiestro(ctx context.Context, data entity.DiestroResult) error {
	speciment, err := u.SpecimenRepository.FindByBarcode(ctx, data.PatientID)
	if err != nil {
		slog.Error("specimen not found", "barcode", data.PatientID, "error", err)
		return err
	}

	observation := entity.ObservationResult{
		SpecimenID: int64(speciment.ID),
		TestCode:   data.TestName,
		Values:     []string{fmt.Sprintf("%.2f", data.Value)},
		Unit:       data.Unit,
		Date:       data.Timestamp,
	}

	fmt.Println(observation)

	err = u.ObservationResultRepository.Create(ctx, &observation)
	if err != nil {
		slog.Error("failed to create observation result", "specimen_id", speciment.ID, "test_code", data.TestName, "error", err)
	}
	return nil
}
