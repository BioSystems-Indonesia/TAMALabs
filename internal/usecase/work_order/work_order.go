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

func (p WorkOrderUseCase) FindOneByID(id int64) (entity.WorkOrder, error) {
	return p.workOrderRepo.FindOne(id)
}

func (p WorkOrderUseCase) Create(req *entity.WorkOrder) error {
	req.Status = entity.WorkOrderStatusNew

	return p.workOrderRepo.Create(req)
}

func (p WorkOrderUseCase) AddSpeciment(workOrderID int64, req *entity.WorkOrderAddSpeciment) (entity.WorkOrder, error) {
	err := p.workOrderRepo.AddSpeciment(workOrderID, req)
	if err != nil {
		return entity.WorkOrder{}, err
	}

	return p.workOrderRepo.FindOne(workOrderID)
}

func (p WorkOrderUseCase) Update(req *entity.WorkOrder) error {
	return p.workOrderRepo.Update(req)
}

func (p WorkOrderUseCase) Delete(id int64) error {
	return p.workOrderRepo.Delete(id)
}
