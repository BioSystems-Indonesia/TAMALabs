package workOrderrepo

import (
	"context"
	"errors"
	"fmt"
	"slices"
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

func (r WorkOrderRepository) FindByStatus(ctx context.Context, status entity.WorkOrderStatus) ([]entity.WorkOrder, error) {
	var workOrder []entity.WorkOrder
	err := r.db.Where("status = ?", status).
		Preload("Patient").
		Preload("Patient.Specimen").
		Preload("Patient.Specimen.ObservationResult").
		Find(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return workOrder, nil
}

func (r WorkOrderRepository) FindOne(id int64) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("id = ?", id).
		Preload("Patient").
		Preload("Patient.Specimen", "order_id = ?", id).
		Preload("Patient.Specimen.ObservationRequest").
		First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

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

func (r WorkOrderRepository) AddTest(workOrder *entity.WorkOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := r.deleteUnusedRelation(tx, workOrder)
		if err != nil {
			return err
		}

		err = r.upsertRelation(tx, workOrder)
		if err != nil {
			return err
		}

		err = tx.Model(entity.WorkOrder{}).Where("id = ?", workOrder.ID).
			Update("updated_at", time.Now()).Error
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

	selectedPatient := util.Filter(oldWorkOrder.Patient, func(patient entity.Patient) bool {
		return slices.Contains(workOrder.PatientIDs, patient.ID)
	})
	selectedPatientObservationRequestCodes := util.Unique(util.Flatten(
		util.Map(selectedPatient, func(patient entity.Patient) []string {
			return util.Flatten(util.Map(patient.Specimen, func(specimen entity.Specimen) []string {
				return util.Map(specimen.ObservationRequest, func(observationRequest entity.ObservationRequest) string {
					return observationRequest.TestCode
				})
			}))
		})))

	toDeleteObservationRequest, _ := util.CompareSlices(
		selectedPatientObservationRequestCodes,
		workOrder.ObservationRequests,
	)
	for _, observationRequestID := range toDeleteObservationRequest {
		var specimentIDs []int64
		err = tx.Model(entity.Specimen{}).Where("order_id = ? and patient_id in (?)", workOrder.ID, workOrder.PatientIDs).
			Pluck("id", &specimentIDs).Error
		if err != nil {
			return fmt.Errorf("error finding specimen work_order:%d: %w", workOrder.ID, err)
		}

		err = tx.Model(&entity.ObservationRequest{}).
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

		testTypes, err := r.getTestType(trx, workOrder.ObservationRequests)
		if err != nil {
			return fmt.Errorf("error getting test type: %w", err)
		}

		for _, testType := range testTypes {
			observationRequest := entity.ObservationRequest{
				TestCode:        testType.Code,
				TestDescription: testType.Name,
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
				"testCode":        testType.Code,
				"testDescription": testType.Name,
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

func (r WorkOrderRepository) getTestType(trx *gorm.DB, observationRequest []string) ([]entity.TestType, error) {
	var testTypes []entity.TestType
	err := trx.Where("code in (?)", observationRequest).Find(&testTypes).Error
	if err != nil {
		return nil, err
	}

	return testTypes, nil
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

func (r WorkOrderRepository) DeleteTest(workOrderID int64, patientID int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Where("patient_id = ? and work_order_id = ?", patientID, workOrderID).
			Delete(&entity.WorkOrderPatient{}).Error
		if err != nil {
			return fmt.Errorf("error updating workOrderPatient: %w", err)
		}

		err = tx.Where("order_id = ? and patient_id = ?", workOrderID, patientID).
			Delete(&entity.Specimen{}).Error
		if err != nil {
			return fmt.Errorf("error updating specimen: %w", err)
		}

		queryGetSpecimenID := tx.Model(entity.Specimen{}).Where("order_id = ? and patient_id = ?", workOrderID, patientID).
			Select("id")
		err = tx.
			Where("specimen_id in (?)", queryGetSpecimenID).
			Delete(&entity.ObservationRequest{}).Error
		if err != nil {
			return fmt.Errorf("error updating observationRequest: %w", err)
		}

		return nil
	})
}

func (r WorkOrderRepository) Update(workOrder *entity.WorkOrder) error {
	res := r.db.Save(workOrder)
	if res.Error != nil {
		return fmt.Errorf("error deleting workOrder: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
