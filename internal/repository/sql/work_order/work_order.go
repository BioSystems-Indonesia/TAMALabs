package workOrderrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkOrderRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewWorkOrderRepository(db *gorm.DB, cfg *config.Schema) *WorkOrderRepository {
	return &WorkOrderRepository{db: db, cfg: cfg}
}

func (r WorkOrderRepository) FindAll(ctx context.Context, req *entity.WorkOrderGetManyRequest) ([]entity.WorkOrder, error) {
	var workOrders []entity.WorkOrder

	db := r.db.WithContext(ctx)
	if len(req.ID) > 0 {
		db = db.Where("id in (?)", req.ID)
	}

	if len(req.SpecimentIDs) > 0 {
		db = db.Joins("join work_order_speciments on work_order_speciments.work_order_id = work_orders.id and work_order_speciments.speciment_id in (?)", req.SpecimentIDs)
	}

	if req.Sort != "" {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: req.Sort,
			},
			Desc: req.IsSortDesc(),
		})
	}

	err := db.Find(&workOrders).Error
	if err != nil {
		return nil, fmt.Errorf("error finding workOrders: %w", err)
	}
	return workOrders, nil
}

func (r WorkOrderRepository) FindOne(id int64) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("id = ?", id).Preload("Speciments").Preload("Speciments.Patient").First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	workOrder.SpecimentIDs = make([]int64, len(workOrder.Speciments))
	for i, speciment := range workOrder.Speciments {
		workOrder.SpecimentIDs[i] = speciment.ID
	}

	return workOrder, nil
}

func (r WorkOrderRepository) Create(workOrder *entity.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(workOrder).Error
		if err != nil {
			return err
		}

		for _, specimentID := range workOrder.SpecimentIDs {
			workOrderSpeciment := entity.WorkOrderSpeciment{
				WorkOrderID: workOrder.ID,
				SpecimentID: specimentID,
			}
			err := tx.Create(&workOrderSpeciment).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r WorkOrderRepository) AddSpeciment(workOrderID int64, req *entity.WorkOrderAddSpeciment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, specimentID := range req.SpecimentIDs {
			workOrderSpeciment := entity.WorkOrderSpeciment{
				WorkOrderID: workOrderID,
				SpecimentID: specimentID,
			}
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&workOrderSpeciment).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r WorkOrderRepository) Update(workOrder *entity.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(workOrder).Error
		if err != nil {
			return err
		}

		err = tx.Delete(&entity.WorkOrderSpeciment{}, "work_order_id = ?", workOrder.ID).Error
		if err != nil {
			return err
		}

		for _, specimentID := range workOrder.SpecimentIDs {
			workOrderSpeciment := entity.WorkOrderSpeciment{
				WorkOrderID: workOrder.ID,
				SpecimentID: specimentID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err := tx.Create(&workOrderSpeciment).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r WorkOrderRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.WorkOrder{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting workOrder: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
