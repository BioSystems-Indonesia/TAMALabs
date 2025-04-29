package workOrderuc

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
)

type WorkOrderUseCase struct {
	cfg                *config.Schema
	workOrderRepo      *workOrderrepo.WorkOrderRepository
	validate           *validator.Validate
	barcodeGeneratorUC usecase.BarcodeGenerator
}

func NewWorkOrderUseCase(
	cfg *config.Schema,
	workOrderRepo *workOrderrepo.WorkOrderRepository,
	validate *validator.Validate,
	barcodeGeneratorUC usecase.BarcodeGenerator,
) *WorkOrderUseCase {
	return &WorkOrderUseCase{cfg: cfg, workOrderRepo: workOrderRepo, validate: validate, barcodeGeneratorUC: barcodeGeneratorUC}
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
