package workOrderrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkOrderRepository struct {
	db            *gorm.DB
	cfg           *config.Schema
	specimentRepo *specimen.Repository
}

func NewWorkOrderRepository(db *gorm.DB, cfg *config.Schema, specimentRepo *specimen.Repository) *WorkOrderRepository {
	return &WorkOrderRepository{db: db, cfg: cfg, specimentRepo: specimentRepo}
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

func (r WorkOrderRepository) FindManyByID(ctx context.Context, id []int64) ([]entity.WorkOrder, error) {
	var workOrders []entity.WorkOrder
	err := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("Patient.Specimen", "order_id in (?)", id).
		Preload("Patient.Specimen.ObservationRequest").
		Find(&workOrders, "id in (?)", id).Error
	if err != nil {
		return nil, fmt.Errorf("error finding workOrders: %w", err)
	}
	return workOrders, nil
}
func (r WorkOrderRepository) FindOne(id int64) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("id = ?", id).
		Preload("Specimen").
		Preload("Specimen.ObservationRequest").
		Preload("Specimen.Patient").
		Preload("Patient").First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	workOrder.PatientIDs = make([]int64, len(workOrder.Patient))
	for i, patient := range workOrder.Patient {
		workOrder.PatientIDs[i] = patient.ID
	}
	workOrder.PatientIDs = util.Unique(workOrder.PatientIDs)

	workOrder.ObservationRequests = make([]string, 0)
	for _, specimen := range workOrder.Specimen {
		for _, observationRequest := range specimen.ObservationRequest {
			workOrder.ObservationRequests = append(workOrder.ObservationRequests, observationRequest.TestCode)
		}
	}
	workOrder.ObservationRequests = util.Unique(workOrder.ObservationRequests)

	return workOrder, nil
}

func (r WorkOrderRepository) Create(workOrder *entity.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(workOrder).Error
		if err != nil {
			return err
		}

		err = r.upsertRelation(tx, workOrder)
		if err != nil {
			return err
		}
		return nil
	})
}

func (r WorkOrderRepository) Update(workOrder *entity.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := r.deleteUnusedRelation(tx, workOrder)
		if err != nil {
			return err
		}

		err = r.upsertRelation(tx, workOrder)
		if err != nil {
			return err
		}

		err = tx.Save(workOrder).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r WorkOrderRepository) deleteUnusedRelation(tx *gorm.DB, workOrder *entity.WorkOrder) error {
	oldWorkOrder, err := r.FindOne(workOrder.ID)
	if err != nil {
		return err
	}

	toDeletePatient, _ := util.CompareSlices(oldWorkOrder.PatientIDs, workOrder.PatientIDs)
	for _, patientID := range toDeletePatient {
		var specimens []entity.Specimen
		err = tx.Find(&specimens, "order_id = ? AND patient_id = ?", workOrder.ID, patientID).Error
		if err != nil {
			return fmt.Errorf("error finding specimen work_order:%d patient:%d: %w", workOrder.ID, patientID, err)
		}

		for _, s := range specimens {
			err = tx.Delete(&entity.Specimen{}, "id = ?", s.ID).Error
			if err != nil {
				return fmt.Errorf("error deleting specimen specimen:%d: %w", s.ID, err)
			}

			var observationRequests []entity.ObservationRequest
			err = tx.Delete(&observationRequests, "specimen_id = ?", s.ID).Error
			if err != nil {
				return fmt.Errorf("error deleting observationRequest specimen:%d: %w", s.ID, err)
			}
		}

		err = tx.Delete(&entity.WorkOrderPatient{}, "work_order_id = ? AND patient_id = ?", workOrder.ID, patientID).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderPatient work_order:%d patient:%d: %w", workOrder.ID, patientID, err)
		}
	}

	toDeleteObservationRequest, _ := util.CompareSlices(
		oldWorkOrder.ObservationRequests,
		workOrder.ObservationRequests,
	)
	for _, observationRequestID := range toDeleteObservationRequest {
		var specimentIDs []int64
		err = tx.Model(entity.Specimen{}).Where("order_id = ?", workOrder.ID).Pluck("id", &specimentIDs).Error
		if err != nil {
			return fmt.Errorf("error finding specimen work_order:%d: %w", workOrder.ID, err)
		}

		err := tx.Model(&entity.ObservationRequest{}).
			Where("specimen_id in (?) AND test_code = ?", specimentIDs, observationRequestID).
			Delete(&entity.ObservationRequest{}).Error
		if err != nil {
			return fmt.Errorf("error deleting observationRequest %v: %w", observationRequestID, err)
		}
	}

	return nil
}

const defaultSerumType = "SER"

func (r WorkOrderRepository) upsertRelation(trx *gorm.DB, workOrder *entity.WorkOrder) error {
	for _, patientID := range workOrder.PatientIDs {
		workOrderPatient := entity.WorkOrderPatient{
			WorkOrderID: workOrder.ID,
			PatientID:   patientID,
		}
		workOrderPatientQuery := trx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&workOrderPatient)
		if workOrderPatientQuery.Error != nil {
			return workOrderPatientQuery.Error
		}

		specimen := entity.Specimen{
			PatientID:      int(patientID),
			OrderID:        int(workOrder.ID),
			Type:           defaultSerumType, // TODO: Change it so it not be hardcoded
			Barcode:        r.specimentRepo.GenerateBarcode(trx.Statement.Context),
			CollectionDate: time.Now().Format(time.RFC3339),
		}
		specimenQuery := trx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&specimen)
		err := specimenQuery.Error
		if err != nil {
			return err
		}

		log.Debugj(map[string]interface{}{
			"message":     "specimen insert",
			"patientID":   patientID,
			"specimenID":  specimen.ID,
			"type":        defaultSerumType,
			"rowAffected": specimenQuery.RowsAffected,
		})

		if specimenQuery.RowsAffected != 0 {
			err := r.specimentRepo.IncrementBarcodeSequence(trx.Statement.Context)
			if err != nil {
				return err
			}
		}

		if specimenQuery.RowsAffected == 0 {
			err = trx.Where("patient_id = ? AND order_id = ? AND type = ?", patientID, workOrder.ID, defaultSerumType).
				Find(&specimen).Error
			if err != nil {
				return err
			}

			log.Debugj(map[string]interface{}{
				"message":     "specimen find",
				"patientID":   patientID,
				"specimenID":  specimen.ID,
				"type":        defaultSerumType,
				"rowAffected": specimenQuery.RowsAffected,
			})
		}

		for _, observationRequestID := range workOrder.ObservationRequests {
			observationType, ok := entity.TableObservationType.Find(observationRequestID)
			if !ok {
				return fmt.Errorf("observation request: %w", entity.ErrBadRequest)
			}

			observationRequest := entity.ObservationRequest{
				TestCode:        observationType.ID,
				TestDescription: observationType.Name,
				SpecimenID:      int64(specimen.ID),
				RequestedDate:   time.Now(),
			}

			observationRequestQuery := trx.Clauses(clause.OnConflict{DoNothing: true}).Create(&observationRequest)
			err := observationRequestQuery.Error
			if err != nil {
				return err
			}

			log.Debugj(map[string]interface{}{
				"message":         "observation request insert",
				"testCode":        observationType.ID,
				"testDescription": observationType.Name,
				"patientID":       patientID,
				"specimenID":      specimen.ID,
				"rowAffected":     observationRequestQuery.RowsAffected,
			})
			if observationRequestQuery.RowsAffected == 0 {
				continue
			}
		}
	}

	return nil
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
