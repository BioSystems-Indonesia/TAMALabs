package workOrderrepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
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

func (r *WorkOrderRepository) FindAllForResult(ctx context.Context, req *entity.ResultGetManyRequest) (entity.PaginationResponse[entity.WorkOrder], error) {
	db := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("Specimen").
		Preload("Specimen.ObservationRequest").
		Preload("Specimen.ObservationRequest.TestType").
		Preload("Specimen.ObservationResult").
		Preload("Specimen.ObservationResult.TestType")

	if len(req.PatientIDs) > 0 {
		db = db.Where("work_orders.patient_id in (?)", req.PatientIDs)
	}

	if len(req.WorkOrderStatus) > 0 {
		db = db.Joins("work_orders.status in (?)", req.WorkOrderStatus)
	}

	if !req.CreatedAtStart.IsZero() {
		db = db.Where("work_orders.created_at >= ?", req.CreatedAtStart.Add(-24*time.Hour))
	}

	if !req.CreatedAtEnd.IsZero() {
		db = db.Where("work_orders.created_at <= ?", req.CreatedAtEnd.Add(24*time.Hour))
	}

	if req.HasResult {
		subQuery := r.db.Table("specimens").Select("specimens.order_id").
			Joins("join observation_results on specimens.id = observation_results.specimen_id").
			Where("observation_results.id is not null")
		db = db.Where("work_orders.id in (?)", subQuery)
	}

	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		TableName: "work_orders",
	})

	return sql.GetWithPaginationResponse[entity.WorkOrder](db, req.GetManyRequest)
}

func (r WorkOrderRepository) FindOneForResult(id int64) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("id = ?", id).
		Preload("Patient").
		Preload("Specimen").
		Preload("Specimen.ObservationRequest").
		Preload("Specimen.ObservationRequest.TestType").
		Preload("Specimen.ObservationResult").
		Preload("Specimen.ObservationResult.TestType").
		First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	return workOrder, nil
}

func (r WorkOrderRepository) FindAll(ctx context.Context, req *entity.WorkOrderGetManyRequest) (entity.PaginationResponse[entity.WorkOrder], error) {
	db := r.db.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})

	if !req.CreatedAtStart.IsZero() {
		db = db.Where("work_orders.created_at >= ?", req.CreatedAtStart.Add(-24*time.Hour))
	}

	if !req.CreatedAtEnd.IsZero() {
		db = db.Where("work_orders.created_at <= ?", req.CreatedAtEnd.Add(24*time.Hour))
	}

	if req.Query != "" {
		db = db.Where("work_orders.id like ? or work_orders.created_at like ?", "%"+req.Query+"%", "%"+req.Query+"%")
	}

	if len(req.SpecimenIDs) > 0 {
		db = db.Joins("join work_order_specimens on work_order_Specimens.work_order_id = work_orders.id and work_order_Specimens.Specimen_id in (?)", req.SpecimenIDs)
	}

	if len(req.PatientIDs) > 0 {
		db = db.Joins("join work_order_patients on work_order_patients.work_order_id = work_orders.id and work_order_patients.patient_id in (?)", req.PatientIDs)
	}

	db = db.Debug().
		Preload("Patient").
		Preload("Specimen").
		Preload("Specimen.ObservationRequest")

	return sql.GetWithPaginationResponse[entity.WorkOrder](db, req.GetManyRequest)
}

func (r WorkOrderRepository) FindOne(id int64) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("id = ?", id).
		Preload("Patient").
		Preload("Patient.Specimen", "order_id = ?", id).
		Preload("Patient.Specimen.ObservationRequest").
		Preload("Patient.Specimen.ObservationRequest.TestType").
		First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	return workOrder, nil
}

func (r WorkOrderRepository) Create(req *entity.WorkOrderCreateRequest) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var patient entity.Patient
		err := tx.Find(&patient, req.PatientID).Error
		if err != nil {
			return fmt.Errorf("error finding patient: %w", err)
		}

		workOrder = entity.WorkOrder{
			Status:    entity.WorkOrderStatusNew,
			PatientID: req.PatientID,
		}
		err = tx.Save(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error creating workOrder: %w", err)
		}

		err = r.upsertRelation(tx, req, &patient, &workOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return workOrder, err
}

func (r WorkOrderRepository) Edit(id int, req *entity.WorkOrderCreateRequest) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var patient entity.Patient
		err := tx.Find(&patient, req.PatientID).Error
		if err != nil {
			return fmt.Errorf("error finding patient: %w", err)
		}

		err = tx.Where("id = ?", id).
			Preload("Patient").
			Preload("Patient.Specimen", "order_id = ?", id).
			Preload("Patient.Specimen.ObservationRequest").
			Preload("Patient.Specimen.ObservationRequest.TestType").
			First(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error finding workOrder: %w", err)
		}
		workOrder.Status = entity.WorkOrderStatusNew
		workOrder.PatientID = req.PatientID

		err = tx.Save(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error creating workOrder: %w", err)
		}

		err = r.deleteUnusedRelation(tx, req, &workOrder)
		if err != nil {
			return err
		}

		err = r.upsertRelation(tx, req, &patient, &workOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return workOrder, err
}

func (r WorkOrderRepository) deleteUnusedRelation(tx *gorm.DB, req *entity.WorkOrderCreateRequest, workOrder *entity.WorkOrder) error {
	oldTestTypeIDs := util.Unique(
		util.Flatten(util.Map(workOrder.Patient.Specimen, func(specimen entity.Specimen) []int64 {
			return util.Map(specimen.ObservationRequest, func(observationRequest entity.ObservationRequest) int64 {
				return int64(observationRequest.TestType.ID)
			})
		})),
	)

	toDeleteTestTypeIDs, _ := util.CompareSlices(
		oldTestTypeIDs,
		req.TestIDs,
	)
	slog.Info("deleteUnusedRelation",
		"oldTestTypeIDs", oldTestTypeIDs,
		"toDeleteObservationRequest", toDeleteTestTypeIDs,
		"testIDs", req.TestIDs,
	)

	var specimentIDs []int64
	err := tx.Model(entity.Specimen{}).Where("order_id = ? and patient_id = ?", workOrder.ID, workOrder.PatientID).
		Pluck("id", &specimentIDs).Error
	if err != nil {
		return fmt.Errorf("error finding specimen work_order:%d: %w", workOrder.ID, err)
	}

	for _, testTypeID := range toDeleteTestTypeIDs {
		var testType entity.TestType
		err = tx.First(&testType, "id = ?", testTypeID).Error
		if err != nil {
			return fmt.Errorf("error finding testType %v: %w", testTypeID, err)
		}

		err = tx.Model(&entity.ObservationRequest{}).
			Where("specimen_id in (?) AND test_code = ?", specimentIDs, testType.Code).
			Delete(&entity.ObservationRequest{}).Error
		if err != nil {
			return fmt.Errorf("error deleting observationRequest %v: %w", testType.Code, err)
		}
	}

	return nil
}

const defaultSerumType = "SER"

func (r WorkOrderRepository) upsertRelation(
	trx *gorm.DB,
	req *entity.WorkOrderCreateRequest,
	patient *entity.Patient,
	workOrder *entity.WorkOrder,
) error {
	specimenTestTypes, err := r.groupBySpecimenType(trx, req)
	if err != nil {
		return fmt.Errorf("error grouping specimen type: %w", err)
	}

	for specimenType, testTypes := range specimenTestTypes {
		specimen := entity.Specimen{
			PatientID:      int(req.PatientID),
			OrderID:        int(workOrder.ID),
			Type:           string(specimenType),
			Barcode:        r.specimentRepo.GenerateBarcode(trx.Statement.Context, specimenType),
			CollectionDate: time.Now().Format(time.RFC3339),
		}
		specimenQuery := trx.Clauses(clause.OnConflict{
			// DoNothing: true,
		}).Debug().Create(&specimen)
		err = specimenQuery.Error
		if err != nil {
			return err
		}

		slog.Info("specimen insert",
			"patientID", patient.ID,
			"specimenID", specimen.ID,
			"workOrderID", specimen.OrderID,
			"type", defaultSerumType,
			"rowAffected", specimenQuery.RowsAffected,
		)

		if specimenQuery.RowsAffected != 0 {
			err := r.specimentRepo.IncrementBarcodeSequence(trx.Statement.Context)
			if err != nil {
				return err
			}
		}

		if specimenQuery.RowsAffected == 0 {
			err = trx.Where("patient_id = ? AND order_id = ? AND type = ?", patient.ID, workOrder.ID, defaultSerumType).
				First(&specimen).Error
			if err != nil {
				return fmt.Errorf("error finding specimen: %w", err)
			}

			slog.Info("specimen find",
				"patientID", patient.ID,
				"specimenID", specimen.ID,
				"type", defaultSerumType,
				"rowAffected", specimenQuery.RowsAffected,
			)
		}

		for _, testType := range testTypes {
			observationRequest := entity.ObservationRequest{
				TestCode:        testType.Code,
				TestDescription: testType.Name,
				SpecimenID:      int64(specimen.ID),
				RequestedDate:   time.Now(),
			}

			observationRequestQuery := trx.Clauses(clause.OnConflict{DoNothing: false}).Create(&observationRequest)
			err := observationRequestQuery.Error
			if err != nil {
				return err
			}

			slog.Info("observation request insert",
				"testCode", testType.Code,
				"testDescription", testType.Name,
				"patientID", patient.ID,
				"specimenID", specimen.ID,
				"rowAffected", observationRequestQuery.RowsAffected,
			)
			if observationRequestQuery.RowsAffected == 0 {
				continue
			}
		}
	}

	return nil
}

func (r WorkOrderRepository) groupBySpecimenType(trx *gorm.DB, req *entity.WorkOrderCreateRequest) (map[entity.SpecimenType][]entity.TestType, error) {
	testTypes, err := r.getTestType(trx, req.TestIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting test type: %w", err)
	}

	specimenTypes := make(map[entity.SpecimenType][]entity.TestType)
	for _, testType := range testTypes {
		specimenType := entity.SpecimenType(testType.Type)
		if specimenType == "" {
			specimenType = entity.SpecimenTypeSER
		}

		specimenTypes[specimenType] = append(
			specimenTypes[specimenType], testType,
		)
	}

	return specimenTypes, nil
}

func (r WorkOrderRepository) getTestType(trx *gorm.DB, observationRequest []int64) ([]entity.TestType, error) {
	var testTypes []entity.TestType
	err := trx.Where("id in (?)", observationRequest).Find(&testTypes).Error
	if err != nil {
		return nil, err
	}

	return testTypes, nil
}

func (r WorkOrderRepository) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var workOrder entity.WorkOrder
		err := tx.Where("id = ?", id).
			First(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error finding workOrder: %w", err)
		}

		queryGetSpecimenID := tx.Model(entity.Specimen{}).
			Where("order_id = ? and patient_id = ?", workOrder.ID, workOrder.PatientID).
			Select("id")
		err = tx.
			Where("specimen_id in (?)", queryGetSpecimenID).
			Delete(&entity.ObservationRequest{}).Error
		if err != nil {
			return fmt.Errorf("error deleting observationRequest: %w", err)
		}

		err = tx.Where("order_id = ? and patient_id = ?", workOrder.ID, workOrder.PatientID).
			Delete(&entity.Specimen{}).Error
		if err != nil {
			return fmt.Errorf("error deleting specimen: %w", err)
		}

		err = tx.Where("work_order_id = ?", workOrder.ID).
			Delete(&entity.WorkOrderDevice{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderDevice: %w", err)
		}

		err = tx.Where("id = ?", workOrder.ID).
			Delete(&entity.WorkOrder{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrder: %w", err)
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

func (r WorkOrderRepository) UpsertDevice(workOrderID int64, deviceID int64) error {
	workOrderDevice := entity.WorkOrderDevice{
		WorkOrderID: workOrderID,
		DeviceID:    deviceID,
	}
	res := r.db.Model(&entity.WorkOrderDevice{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&workOrderDevice)
	if res.Error != nil {
		return fmt.Errorf("error upserting workOrderDevice: %w", res.Error)
	}

	return nil
}
