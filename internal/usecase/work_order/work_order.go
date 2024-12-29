package workOrderuc

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	workOrderrepo "github.com/oibacidem/lims-hl-seven/internal/repository/sql/work_order"
)

type WorkOrderUseCase struct {
	cfg           *config.Schema
	workOrderRepo *workOrderrepo.WorkOrderRepository
	validate      *validator.Validate
}

func NewWorkOrderUseCase(
	cfg *config.Schema,
	workOrderRepo *workOrderrepo.WorkOrderRepository,
	validate *validator.Validate,
) *WorkOrderUseCase {
	return &WorkOrderUseCase{cfg: cfg, workOrderRepo: workOrderRepo, validate: validate}
}

func (p WorkOrderUseCase) FindAll(
	ctx context.Context, req *entity.WorkOrderGetManyRequest,
) ([]entity.WorkOrder, error) {
	return p.workOrderRepo.FindAll(ctx, req)
}

func (p WorkOrderUseCase) FindManyByID(
	ctx context.Context, id []int64,
) ([]entity.WorkOrder, error) {
	return p.workOrderRepo.FindManyByID(ctx, id)
}

func (p WorkOrderUseCase) FindOneByID(id int64) (entity.WorkOrder, error) {
	return p.workOrderRepo.FindOne(id)
}

func (p WorkOrderUseCase) Create(req *entity.WorkOrder) error {
	req.Status = entity.WorkOrderStatusNew

	return p.workOrderRepo.Create(req)
}

func (p WorkOrderUseCase) AddTest(req *entity.WorkOrder) error {
	return p.workOrderRepo.AddTest(req)
}

func (p WorkOrderUseCase) Delete(id int64) error {
	return p.workOrderRepo.Delete(id)
}

func (p WorkOrderUseCase) DeleteTest(workOrderID int64, patientID int64) error {
	return p.workOrderRepo.DeleteTest(workOrderID, patientID)
}

func (p WorkOrderUseCase) Update(workOrder *entity.WorkOrder) error {
	return p.workOrderRepo.Update(workOrder)
}
