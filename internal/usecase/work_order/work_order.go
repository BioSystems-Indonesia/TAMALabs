package workOrderuc

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/work_order/runner"
	"github.com/oibacidem/lims-hl-seven/pkg/panics"
)

type WorkOrderUseCase struct {
	cfg                *config.Schema
	workOrderRepo      *workOrderrepo.WorkOrderRepository
	validate           *validator.Validate
	barcodeGeneratorUC usecase.BarcodeGenerator
	patientUsecase     *patientuc.PatientUseCase
	deviceUsecase      *deviceuc.DeviceUseCase
	runnerStrategy     *runner.Strategy
}

func NewWorkOrderUseCase(
	cfg *config.Schema,
	workOrderRepo *workOrderrepo.WorkOrderRepository,
	validate *validator.Validate,
	barcodeGeneratorUC usecase.BarcodeGenerator,
	patientUsecase *patientuc.PatientUseCase,
	deviceUsecase *deviceuc.DeviceUseCase,
	runnerStrategy *runner.Strategy,
) *WorkOrderUseCase {
	return &WorkOrderUseCase{
		cfg:                cfg,
		workOrderRepo:      workOrderRepo,
		validate:           validate,
		barcodeGeneratorUC: barcodeGeneratorUC,
		patientUsecase:     patientUsecase,
		deviceUsecase:      deviceUsecase,
		runnerStrategy:     runnerStrategy,
	}
}

func (p WorkOrderUseCase) FindAll(
	ctx context.Context, req *entity.WorkOrderGetManyRequest,
) (entity.PaginationResponse[entity.WorkOrder], error) {
	return p.workOrderRepo.FindAll(ctx, req)
}

func (p WorkOrderUseCase) FindOneByID(id int64) (entity.WorkOrder, error) {
	return p.workOrderRepo.FindOne(id)
}

func (p WorkOrderUseCase) Create(req *entity.WorkOrderCreateRequest) (entity.WorkOrder, error) {
	barcode, err := p.barcodeGeneratorUC.NextOrderBarcode(context.Background())
	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("failed to p.barcodeGeneratorUC.NextOrderBarcode %w", err)
	}
	req.Barcode = barcode
	return p.workOrderRepo.Create(req)
}

func (p WorkOrderUseCase) Edit(id int, req *entity.WorkOrderCreateRequest) (entity.WorkOrder, error) {
	return p.workOrderRepo.Edit(id, req)
}

func (p WorkOrderUseCase) Delete(id int64) error {
	return p.workOrderRepo.Delete(id)
}

func (p WorkOrderUseCase) Update(workOrder *entity.WorkOrder) error {
	return p.workOrderRepo.Update(workOrder)
}

func (p WorkOrderUseCase) UpsertDevice(workOrderID int64, deviceID int64) error {
	return p.workOrderRepo.UpsertDevice(workOrderID, deviceID)
}

func (p WorkOrderUseCase) RunWorkOrderAsync(
	ctx context.Context,
	req *entity.WorkOrderRunRequest,
	action constant.WorkOrderRunAction,
) (<-chan entity.WorkOrderRunStreamMessage, error) {
	ch := make(chan entity.WorkOrderRunStreamMessage)
	req.SetProgressWriter(ch)

	go panics.CapturePanic(ctx, func() {
		defer close(ch)

		errChan := make(chan error, 1)
		go panics.CapturePanic(ctx, func() {
			errChan <- p.runWorkOrder(ctx, req, action)
		})

		select {
		case <-ctx.Done():
			slog.ErrorContext(ctx, fmt.Sprintf("Work Order Run Canceled: %v", ctx.Err()))

			postRun, err := p.runnerStrategy.ChoosePostRunner(ctx, constant.WorkOrderRunActionIncompleteSend)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Work Order Run Failed: %v", err))
			}

			err = postRun.PostRun(ctx, req)
			if err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Work Order Run Failed: %v", err))
			}

			ch <- entity.WorkOrderRunStreamMessage{
				Error:  ctx.Err(),
				IsDone: true,
			}
			return
		case errSend := <-errChan:
			if errSend != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Work Order Run Error: %v", errSend))
				ch <- entity.WorkOrderRunStreamMessage{
					Error:  errSend,
					IsDone: true,
				}
				return
			}

			ch <- entity.WorkOrderRunStreamMessage{
				Percentage: 100,
				Status:     entity.WorkOrderStreamingResponseStatusDone,
				IsDone:     true,
			}
			return
		}
	})

	return ch, nil
}

func (p WorkOrderUseCase) runWorkOrder(
	ctx context.Context,
	req *entity.WorkOrderRunRequest,
	action constant.WorkOrderRunAction,
) error {
	preRunner, err := p.runnerStrategy.ChoosePreRunner(ctx, action)
	if err != nil {
		return fmt.Errorf("failed to p.runnerStrategy.ChoosePreRunner %w", err)
	}

	err = preRunner.PreRun(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to p.runnerStrategy.PreRun %w", err)
	}

	err = p.fillData(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to p.fillData %w", err)
	}

	payload := &entity.SendPayloadRequest{
		Patients:       req.GetPatients(),
		Device:         req.GetDevice(),
		Urgent:         req.Urgent,
		ProgressWriter: req.ProgressWriter(),
	}
	sender, err := p.runnerStrategy.ChooseSendRunner(ctx, payload.Device)
	if err != nil {
		return fmt.Errorf("failed to p.runnerStrategy.ChooseSendRunner %w", err)
	}

	err = sender.Send(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed on sending %w", err)
	}

	postRunner, err := p.runnerStrategy.ChoosePostRunner(ctx, action)
	if err != nil {
		return fmt.Errorf("failed to p.runnerStrategy.ChoosePreRunner %w", err)
	}

	err = postRunner.PostRun(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to p.runnerStrategy.PreRun %w", err)
	}

	return nil
}

func (p WorkOrderUseCase) fillData(ctx context.Context, req *entity.WorkOrderRunRequest) error {
	patients, err := p.patientUsecase.FindManyByWorkOrderID(ctx, req.WorkOrderIDs)
	if err != nil {
		return fmt.Errorf("failed to p.patientUsecase.FindManyByWorkOrderID %w", err)
	}

	device, err := p.deviceUsecase.FindOneByID(ctx, req.DeviceID)
	if err != nil {
		return fmt.Errorf("failed to p.deviceUsecase.FindOneByID %w", err)
	}

	req.SetPatients(patients)
	req.SetDevice(device)

	return nil
}
