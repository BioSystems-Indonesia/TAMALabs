package workOrderrepo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql/specimen"
	"github.com/oibacidem/lims-hl-seven/internal/util"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WorkOrderRepository struct {
	db            *gorm.DB
	cfg           *config.Schema
	specimentRepo *specimen.Repository
	cache         *cache.Cache
}

func NewWorkOrderRepository(db *gorm.DB, cfg *config.Schema, specimentRepo *specimen.Repository, cache *cache.Cache) *WorkOrderRepository {
	r := &WorkOrderRepository{db: db, cfg: cfg, specimentRepo: specimentRepo, cache: cache}

	err := r.SyncBarcodeSequence(context.Background())
	if err != nil {
		panic(err)
	}

	return r
}

func (r *WorkOrderRepository) FindAllForResult(ctx context.Context, req *entity.ResultGetManyRequest) (entity.PaginationResponse[entity.WorkOrder], error) {
	db := r.db.WithContext(ctx).
		Preload("Patient").
		Preload("Specimen").
		Preload("Specimen.ObservationRequest").
		Preload("Specimen.ObservationRequest.TestType").
		Preload("Specimen.ObservationResult").
		Preload("Specimen.ObservationResult.TestType").
		Preload("Doctors").
		Preload("Analyzers")

	if len(req.PatientIDs) > 0 {
		db = db.Where("work_orders.patient_id in (?)", req.PatientIDs)
	}

	if len(req.BarcodeIDs) > 0 {
		db = db.Where("work_orders.barcode in (?)", req.BarcodeIDs)
	}

	if len(req.WorkOrderStatus) > 0 {
		db = db.Where("work_orders.status in (?)", req.WorkOrderStatus)
	}

	if len(req.DoctorIDs) > 0 {
		subQuery := r.db.Table("work_order_doctors").Select("work_order_doctors.work_order_id").
			Where("work_order_doctors.admin_id in (?)", req.DoctorIDs)
		db = db.Where("work_orders.id in (?)", subQuery)
	}

	if !req.CreatedAtStart.IsZero() {
		db = db.Where("work_orders.created_at >= ?", req.CreatedAtStart.Add(-24*time.Hour))
	}

	if !req.CreatedAtEnd.IsZero() {
		db = db.Where("work_orders.created_at <= ?", req.CreatedAtEnd.Add(3*24*time.Hour))
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

	resp, err := sql.GetWithPaginationResponse[entity.WorkOrder](db, req.GetManyRequest)
	if err != nil {
		return entity.PaginationResponse[entity.WorkOrder]{}, fmt.Errorf("error finding workOrders: %w", err)
	}

	for i := range resp.Data {
		resp.Data[i].FillData()
	}

	return resp, nil
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
		Preload("Doctors").
		Preload("Analyzers").
		First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	workOrder.FillData()
	return workOrder, nil
}

func (r WorkOrderRepository) FindAll(
	ctx context.Context, req *entity.WorkOrderGetManyRequest,
) (entity.PaginationResponse[entity.WorkOrder], error) {
	db := r.db.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})

	// Prioritize ID filter
	if len(req.GetManyRequest.ID) == 0 {
		if !req.CreatedAtStart.IsZero() {
			db = db.Where("work_orders.created_at >= ?", req.CreatedAtStart.Add(-24*time.Hour))
		}

		if !req.CreatedAtEnd.IsZero() {
			db = db.Where("work_orders.created_at <= ?", req.CreatedAtEnd.Add(24*time.Hour))
		}
	}

	if len(req.BarcodeIds) > 0 {
		db = db.Where("work_orders.barcode in (?)", req.BarcodeIds)
	}

	if len(req.SpecimenIDs) > 0 {
		db = db.Joins("join work_order_specimens on work_order_Specimens.work_order_id = work_orders.id and work_order_Specimens.Specimen_id in (?)", req.SpecimenIDs)
	}

	if len(req.PatientIDs) > 0 {
		db = db.Where("work_orders.patient_id in (?)", req.PatientIDs)
	}

	db = db.
		Preload("Patient").
		Preload("Specimen").
		Preload("Specimen.ObservationRequest").
		Preload("Devices").
		Preload("Doctors").
		Preload("Analyzers")

	resp, err := sql.GetWithPaginationResponse[entity.WorkOrder](db, req.GetManyRequest)
	if err != nil {
		return entity.PaginationResponse[entity.WorkOrder]{}, fmt.Errorf("error finding workOrders: %w", err)
	}

	for i := range resp.Data {
		resp.Data[i].FillData()
	}

	return resp, nil
}

func (r WorkOrderRepository) FindAllBarcodes(ctx context.Context) ([]string, error) {
	var barcodes []string
	err := r.db.WithContext(ctx).Table("work_orders").Select("barcode").
		Order("barcode desc").Find(&barcodes).Error
	if err != nil {
		return nil, fmt.Errorf("error finding workOrders: %w", err)
	}
	return barcodes, nil
}

func (r WorkOrderRepository) FindOne(id int64) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("id = ?", id).
		Preload("Patient").
		Preload("Patient.Specimen", "order_id = ?", id).
		Preload("Patient.Specimen.ObservationRequest").
		Preload("Patient.Specimen.ObservationRequest.TestType").
		Preload("Devices").
		Preload("Doctors").
		Preload("Analyzers").
		Preload("TestTemplates").
		First(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	workOrder.FillData()
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

		verifiedStatus := string(entity.WorkOrderVerifiedStatusPending)
		if len(req.DoctorIDs) == 0 {
			verifiedStatus = string(entity.WorkOrderVerifiedStatusVerified)
		}

		workOrder = entity.WorkOrder{
			Status:          entity.WorkOrderStatusNew,
			VerifiedStatus:  verifiedStatus,
			PatientID:       req.PatientID,
			Barcode:         req.Barcode,
			BarcodeSIMRS:    req.BarcodeSIMRS,
			AnalyzerIDs:     req.AnalyzerIDs,
			DoctorIDs:       req.DoctorIDs,
			TestTemplateIDs: req.TestTemplateIDs,
			CreatedBy:       req.CreatedBy,
			LastUpdatedBy:   req.CreatedBy,
		}
		err = tx.Save(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error creating workOrder: %w", err)
		}

		err = r.upsertRelation(tx, req, &patient, &workOrder)
		if err != nil {
			return fmt.Errorf("error upserting relation: %w", err)
		}

		err = r.IncrementBarcodeSequence(tx.Statement.Context)
		if err != nil {
			return fmt.Errorf("error incrementing barcode sequence: %w", err)
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
			Preload("Doctors").
			Preload("Analyzers").
			Preload("TestTemplates").
			First(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error finding workOrder: %w", err)
		}
		workOrder.Status = entity.WorkOrderStatusNew
		if len(workOrder.DoctorIDs) == 0 {
			workOrder.VerifiedStatus = string(entity.WorkOrderVerifiedStatusVerified)
		} else {
			workOrder.VerifiedStatus = string(entity.WorkOrderVerifiedStatusPending)
		}
		workOrder.PatientID = req.PatientID
		workOrder.Barcode = req.Barcode
		workOrder.BarcodeSIMRS = req.BarcodeSIMRS
		workOrder.LastUpdatedBy = req.CreatedBy
		workOrder.FillData()

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

func (r WorkOrderRepository) deleteUnusedRelation(
	tx *gorm.DB,
	req *entity.WorkOrderCreateRequest,
	workOrder *entity.WorkOrder,
) error {
	toDeleteDoctorIDs, _ := util.CompareSlices(
		workOrder.DoctorIDs,
		req.DoctorIDs,
	)
	for _, doctorID := range toDeleteDoctorIDs {
		err := tx.Where("work_order_id =? AND admin_id =?", workOrder.ID, doctorID).
			Delete(&entity.WorkOrderDoctor{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderDoctor: %w", err)
		}
	}
	slog.Info("deleteUnusedRelation",
		"workOrder.DoctorIDs", workOrder.DoctorIDs,
		"req.DoctorIDs", req.DoctorIDs,
		"toDeleteDoctorIDs", toDeleteDoctorIDs,
	)

	toDeleteAnalyzerIDs, _ := util.CompareSlices(
		workOrder.AnalyzerIDs,
		req.AnalyzerIDs,
	)
	for _, analyzerID := range toDeleteAnalyzerIDs {
		err := tx.Where("work_order_id =? AND admin_id =?", workOrder.ID, analyzerID).
			Delete(&entity.WorkOrderAnalyzer{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderAnalyzer: %w", err)
		}
	}
	slog.InfoContext(tx.Statement.Context, "deleteUnusedRelation",
		"workOrder.AnalyzerIDs", workOrder.AnalyzerIDs,
		"req.AnalyzerIDs", req.AnalyzerIDs,
		"toDeleteAnalyzerIDs", toDeleteAnalyzerIDs,
	)

	toDeleteTestTemplateIDs, _ := util.CompareSlices(
		workOrder.TestTemplateIDs,
		req.TestTemplateIDs,
	)
	for _, testTemplateID := range toDeleteTestTemplateIDs {
		err := tx.Where("work_order_id =? AND test_template_id =?", workOrder.ID, testTemplateID).
			Delete(&entity.WorkOrderTestTemplate{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderTestTemplate: %w", err)
		}
	}

	for _, specimen := range workOrder.Patient.Specimen {
		oldTestTypeIDs := util.Unique(
			util.Map(specimen.ObservationRequest, func(observationRequest entity.ObservationRequest) int64 {
				return int64(observationRequest.TestType.ID)
			}),
		)

		reqSameSpecimen := util.Filter(req.TestTypes, func(testType entity.WorkOrderCreateRequestTestType) bool {
			return testType.SpecimenType == specimen.Type
		})
		testIDs := util.Map(reqSameSpecimen, func(testType entity.WorkOrderCreateRequestTestType) int64 {
			return testType.TestTypeID
		})

		toDeleteTestTypeIDs, _ := util.CompareSlices(
			oldTestTypeIDs,
			testIDs,
		)
		slog.Info("deleteUnusedRelation",
			"oldTestTypeIDs", oldTestTypeIDs,
			"toDeleteObservationRequest", toDeleteTestTypeIDs,
			"testIDs", req.TestTypes,
		)

		for _, testTypeID := range toDeleteTestTypeIDs {
			var testType entity.TestType
			err := tx.First(&testType, "id = ?", testTypeID).Error
			if err != nil {
				return fmt.Errorf("error finding testType %v: %w", testTypeID, err)
			}

			err = tx.Model(&entity.ObservationRequest{}).
				Where("specimen_id = ? AND test_code = ?", specimen.ID, testType.Code).
				Delete(&entity.ObservationRequest{}).Error
			if err != nil {
				return fmt.Errorf("error deleting observationRequest %v: %w", testType.Code, err)
			}
		}

		var observationRequestCount int64
		err := tx.Model(&entity.ObservationRequest{}).
			Where("specimen_id = ?", specimen.ID).
			Count(&observationRequestCount).Error
		if err != nil {
			return fmt.Errorf("error counting observationRequest specimenID: %v: %w", specimen.ID, err)
		}
		if observationRequestCount > 0 {
			continue
		}

		err = tx.Where("id = ?", specimen.ID).
			Delete(&entity.Specimen{}).Error
		if err != nil {
			return fmt.Errorf("error deleting specimen specimenID: %v: %w", specimen.ID, err)
		}
	}

	return nil
}

func (r WorkOrderRepository) upsertRelation(
	trx *gorm.DB,
	req *entity.WorkOrderCreateRequest,
	patient *entity.Patient,
	workOrder *entity.WorkOrder,
) error {
	for _, adminID := range req.DoctorIDs {
		workOrderDoctor := entity.WorkOrderDoctor{
			WorkOrderID: workOrder.ID,
			AdminID:     adminID,
		}
		err := trx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&workOrderDoctor).Error
		if err != nil {
			return fmt.Errorf("error upserting workOrderDoctor: %w", err)
		}
	}

	for _, analyzerID := range req.AnalyzerIDs {
		workOrderAnalyzer := entity.WorkOrderAnalyzer{
			WorkOrderID: workOrder.ID,
			AdminID:     analyzerID,
		}
		err := trx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&workOrderAnalyzer).Error
		if err != nil {
			return fmt.Errorf("error upserting workOrderAnalyzer: %w", err)
		}
	}

	for _, testTemplateID := range req.TestTemplateIDs {
		testTemplate := entity.WorkOrderTestTemplate{
			WorkOrderID:    workOrder.ID,
			TestTemplateID: testTemplateID,
		}

		err := trx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&testTemplate).Error
		if err != nil {
			return fmt.Errorf("error upserting testTemplate: %w", err)
		}
	}

	specimenTestTypes, err := r.groupBySpecimenType(trx, req)
	if err != nil {
		return fmt.Errorf("error grouping specimen type: %w", err)
	}

	for specimenType, testTypes := range specimenTestTypes {
		// Check if barcode is auto-generated (format: YYMMDDNNN)
		// If it's custom barcode, don't add specimen type prefix
		var specimenBarcode string
		autoGeneratedPattern := regexp.MustCompile(`^\d{6}\d{3}$`) // YYMMDDNNN pattern
		if autoGeneratedPattern.MatchString(workOrder.Barcode) {
			// Auto-generated barcode, add specimen type prefix
			specimenBarcode = fmt.Sprintf("%s%s", specimenType, workOrder.Barcode)
		} else {
			// Custom barcode, use as is
			specimenBarcode = workOrder.Barcode
		}

		specimen := entity.Specimen{
			PatientID:      int(req.PatientID),
			OrderID:        int(workOrder.ID),
			Type:           string(specimenType),
			Barcode:        specimenBarcode,
			CollectionDate: time.Now().Format(time.RFC3339),
		}
		specimenQuery := trx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&specimen)
		err = specimenQuery.Error
		if err != nil {
			return err
		}

		slog.Info("specimen insert",
			"patientID", patient.ID,
			"specimenID", specimen.ID,
			"workOrderID", specimen.OrderID,
			"type", specimenType,
			"rowAffected", specimenQuery.RowsAffected,
		)

		if specimenQuery.RowsAffected == 0 {
			err = trx.Where("patient_id = ? AND order_id = ? AND type = ?", patient.ID, workOrder.ID, specimenType).
				First(&specimen).Error
			if err != nil {
				return fmt.Errorf("error finding specimen: %w", err)
			}

			slog.Info("specimen find",
				"patientID", patient.ID,
				"specimenID", specimen.ID,
				"type", specimenType,
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

			observationRequestQuery := trx.Clauses(clause.OnConflict{DoNothing: true}).Create(&observationRequest)
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
	testIDs := util.Map(req.TestTypes, func(testType entity.WorkOrderCreateRequestTestType) int64 {
		return testType.TestTypeID
	})

	testTypes, err := r.getTestType(trx, testIDs)
	if err != nil {
		return nil, fmt.Errorf("error getting test type: %w", err)
	}

	specimenTypes := make(map[entity.SpecimenType][]entity.TestType)
	for _, testType := range req.TestTypes {
		specimenType := entity.SpecimenType(testType.SpecimenType)
		if specimenType == "" {
			specimenType = entity.SpecimenTypeSER
		}

		tt, ok := testTypes[int(testType.TestTypeID)]
		if !ok {
			slog.WarnContext(trx.Statement.Context, "test type not found",
				"testTypeID", testType.TestTypeID)
		}

		specimenTypes[specimenType] = append(
			specimenTypes[specimenType], tt,
		)
	}

	return specimenTypes, nil
}

func (r WorkOrderRepository) getTestType(trx *gorm.DB, observationRequest []int64) (map[int]entity.TestType, error) {
	var testTypes []entity.TestType
	err := trx.Where("id in (?)", observationRequest).Find(&testTypes).Error
	if err != nil {
		return nil, fmt.Errorf("error finding test type: %w", err)
	}

	var testTypeCache = make(map[int]entity.TestType)
	for _, testType := range testTypes {
		testTypeCache[testType.ID] = testType
	}

	return testTypeCache, nil
}

func (r WorkOrderRepository) Delete(id int64) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var workOrder entity.WorkOrder
		err := tx.Where("id = ?", id).
			First(&workOrder).Error
		if err != nil {
			return fmt.Errorf("error finding workOrder: %w", err)
		}
		workOrder.FillData()

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

		err = tx.Where("work_order_id =?", workOrder.ID).
			Delete(&entity.WorkOrderDoctor{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderDoctor: %w", err)
		}

		err = tx.Where("work_order_id =?", workOrder.ID).
			Delete(&entity.WorkOrderAnalyzer{}).Error
		if err != nil {
			return fmt.Errorf("error deleting workOrderAnalyzer: %w", err)
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
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	res := r.db.Model(&entity.WorkOrderDevice{}).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&workOrderDevice)
	if res.Error != nil {
		return fmt.Errorf("error upserting workOrderDevice: %w", res.Error)
	}

	return nil
}

func RandomNumber(n int) string {
	if n <= 0 {
		return "" // Return empty string for non-positive length
	}

	// More efficient way using strings.Builder and bytes/runes
	const digits = "0123456789"
	var sb strings.Builder
	sb.Grow(n) // Pre-allocate capacity for efficiency

	for i := 0; i < n; i++ {
		randomIndex := rand.Intn(len(digits)) // Get random index 0-9
		sb.WriteByte(digits[randomIndex])     // Append the byte corresponding to the digit
	}

	return sb.String()
}

func (r WorkOrderRepository) GetBarcodeSequence(ctx context.Context) int64 {
	seq, ok := r.cache.Get(constant.KeyWorkOrderBarcodeSequence)
	if !ok {
		now := time.Now()
		tomorrowMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
		expire := tomorrowMidnight.Sub(now)
		r.cache.Set(constant.KeyWorkOrderBarcodeSequence, int64(1), expire)

		return 1
	}

	switch seq.(type) {
	case int64:
		return seq.(int64)
	case int:
		return int64(seq.(int))
	default:
		panic(fmt.Sprintf("unknown type: %T", seq))
	}
}

func (r WorkOrderRepository) IncrementBarcodeSequence(ctx context.Context) error {
	_, _, found := r.cache.GetWithExpiration(constant.KeyWorkOrderBarcodeSequence)
	if !found {
		err := r.SyncBarcodeSequence(ctx)
		if err != nil {
			return fmt.Errorf("error syncing barcode sequence: %w", err)
		}

		return nil
	}

	err := r.cache.Increment(constant.KeyWorkOrderBarcodeSequence, int64(1))
	if err != nil {
		return err
	}

	return nil
}

func (r WorkOrderRepository) SyncBarcodeSequence(ctx context.Context) error {
	now := time.Now()
	currentDayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tomorrowMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

	var count int64
	err := r.db.Model(entity.WorkOrder{}).
		Where("created_at >= ? and created_at < ?", currentDayMidnight, tomorrowMidnight).
		Count(&count).Error
	if err != nil {
		return fmt.Errorf("error counting workOrder: %w", err)
	}

	expire := tomorrowMidnight.Sub(now)
	r.cache.Set(constant.KeyWorkOrderBarcodeSequence, int64(count+1), expire)

	return nil
}

func (r WorkOrderRepository) FindOneByBarcode(ctx context.Context, barcode string) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("barcode = ?", barcode).
		Preload("Patient").
		Preload("Patient.Specimen").
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

func (r WorkOrderRepository) FindByStatus(ctx context.Context, status entity.WorkOrderStatus) ([]entity.WorkOrder, error) {
	var workOrders []entity.WorkOrder
	err := r.db.Where("status = ?", status).
		Preload("Devices").
		Preload("Patient").
		Preload("Devices").
		Preload("Patient.Specimen").
		Preload("Patient.Specimen.ObservationRequest").
		Preload("Patient.Specimen.ObservationRequest.TestType").
		Find(&workOrders).Error
	if err != nil {
		return nil, fmt.Errorf("error finding workOrder: %w", err)
	}

	return workOrders, nil
}

func (r WorkOrderRepository) FindNextID(ctx context.Context, currentWorkOrderID int64) (int64, error) {
	return r.findNearestNumber(ctx, "id > ?", currentWorkOrderID, "ASC")
}

func (r WorkOrderRepository) FindPrevID(ctx context.Context, currentWorkOrderID int64) (int64, error) {
	return r.findNearestNumber(ctx, "id < ?", currentWorkOrderID, "DESC")
}

func (r WorkOrderRepository) findNearestNumber(ctx context.Context, where string, curID int64, dir string) (int64, error) {
	var id int64

	err := r.db.WithContext(ctx).
		Model(entity.WorkOrder{}).
		Select("id").
		Where(where, curID).
		Order("id " + dir).
		Limit(1).
		Scan(&id).Error

	return id, err
}

func (r WorkOrderRepository) GetBySIMRSBarcode(ctx context.Context, barcode string) (entity.WorkOrder, error) {
	var workOrder entity.WorkOrder
	err := r.db.Where("barcode_simrs = ?", barcode).
		Take(&workOrder).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.WorkOrder{}, entity.ErrNotFound
	}
	if err != nil {
		return entity.WorkOrder{}, fmt.Errorf("error finding workOrder: %w", err)
	}

	return workOrder, nil
}

func (r WorkOrderRepository) ChangeStatus(ctx context.Context, id int64, status entity.WorkOrderStatus) error {
	res := r.db.WithContext(ctx).Model(&entity.WorkOrder{}).
		Where("id = ?", id).
		Update("status", status)
	if res.Error != nil {
		return fmt.Errorf("error updating workOrder status: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
