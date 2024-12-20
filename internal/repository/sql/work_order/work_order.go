package workOrderrepo

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

	if len(req.SpecimenIDs) > 0 {
		db = db.Joins("join work_order_Specimens on work_order_Specimens.work_order_id = work_orders.id and work_order_Specimens.Specimen_id in (?)", req.SpecimenIDs)
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
	err := r.db.Where("id = ?", id).Preload("Specimens").Preload("Specimens.Patient").First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	workOrder.SpecimenIDs = make([]int64, len(workOrder.Specimens))
	for i, Specimen := range workOrder.Specimens {
		workOrder.SpecimenIDs[i] = int64(Specimen.ID)
	}

	return workOrder, nil
}

func (r WorkOrderRepository) Create(workOrder *entity.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(workOrder).Error
		if err != nil {
			return err
		}

		var specimens []entity.Specimen
		for _, patientID := range workOrder.PatientIds {
			speciment := entity.Specimen{
				PatientID:      int(patientID),
				Type:           "SER", // TODO: Change it so it not be hardcoded
				CollectionDate: time.Now().Format(time.RFC3339),
			}
			err := tx.Create(&speciment).Error
			if err != nil {
				return err
			}

			specimens = append(specimens, speciment)
		}

		var observationRequests []entity.ObservationRequest
		for _, specimen := range specimens {
			workOrderSpecimen := entity.WorkOrderSpecimen{
				WorkOrderID: workOrder.ID,
				SpecimenID:  int64(specimen.ID),
			}
			err := tx.Create(&workOrderSpecimen).Error
			if err != nil {
				return err
			}

			for _, observationRequestID := range workOrder.ObservationRequests {
				observationType, ok := entity.TableObservationType.Find(observationRequestID)
				if !ok {
					return fmt.Errorf("observation request: %w", entity.ErrBadRequest)
				}

				observationRequest := entity.ObservationRequest{
					TestCode:        observationType.ID,
					TestDescription: observationType.Name,
					SpecimenID:      specimen.ID,
					OrderID:         strconv.Itoa(int(workOrder.ID)),
				}
				err := tx.Create(&observationRequest).Error
				if err != nil {
					return err
				}

				observationRequests = append(observationRequests, observationRequest)
			}
		}

		return nil
	})
}

func (r WorkOrderRepository) AddSpecimen(workOrderID int64, req *entity.WorkOrderAddSpecimen) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, SpecimenID := range req.SpecimenIDs {
			workOrderSpecimen := entity.WorkOrderSpecimen{
				WorkOrderID: workOrderID,
				SpecimenID:  SpecimenID,
			}
			err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&workOrderSpecimen).Error
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

		err = tx.Delete(&entity.WorkOrderSpecimen{}, "work_order_id = ?", workOrder.ID).Error
		if err != nil {
			return err
		}

		for _, SpecimenID := range workOrder.SpecimenIDs {
			workOrderSpecimen := entity.WorkOrderSpecimen{
				WorkOrderID: workOrder.ID,
				SpecimenID:  SpecimenID,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}
			err := tx.Create(&workOrderSpecimen).Error
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
