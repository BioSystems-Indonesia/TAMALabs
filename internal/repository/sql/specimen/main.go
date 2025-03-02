package specimen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

type Repository struct {
	db    *gorm.DB
	cfg   *config.Schema
	cache *cache.Cache
}

func NewRepository(db *gorm.DB, cfg *config.Schema, cache *cache.Cache) *Repository {
	r := &Repository{db: db, cfg: cfg, cache: cache}
	err := r.SyncBarcodeSequence(context.Background())
	if err != nil {
		panic(err)
	}

	return r
}

func (r Repository) FindAll(
	ctx context.Context, req *entity.SpecimenGetManyRequest,
) (entity.PaginationResponse[entity.Specimen], error) {
	db := r.db.WithContext(ctx).Preload("ObservationRequest")
	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{})

	if req.PatientID != 0 {
		db = db.Where("patient_id = ?", req.PatientID)
	}

	return sql.GetWithPaginationResponse[entity.Specimen](db, req.GetManyRequest)
}

func (r Repository) FindAllForResult(
	ctx context.Context, req *entity.ResultGetManyRequest,
) (entity.PaginationResponse[entity.Specimen], error) {
	db := r.db.WithContext(ctx).
		Preload("ObservationResult").
		Preload("ObservationResult.TestType").
		Preload("ObservationRequest").
		Preload("Patient").
		Preload("WorkOrder")

	if len(req.WorkOrderIDs) > 0 {
		db = db.Where("specimens.order_id in (?)", req.WorkOrderIDs)
	}
	if len(req.PatientIDs) > 0 {
		db = db.Where("specimens.patient_id in (?)", req.PatientIDs)
	}

	if len(req.WorkOrderStatus) > 0 {
		db = db.Joins("join work_orders on specimens.order_id = work_orders.id and work_orders.status in (?)", req.WorkOrderStatus)
	}

	if req.HasResult {
		subQuery := r.db.Table("specimens").Select("specimens.id").
			Joins("join observation_results on specimens.id = observation_results.specimen_id").
			Where("observation_results.id is not null")
		db = db.Where("specimens.id in (?)", subQuery)
	}

	db = sql.ProcessGetMany(db, req.GetManyRequest, sql.Modify{
		TableName: "specimens",
	})

	return sql.GetWithPaginationResponse[entity.Specimen](db, req.GetManyRequest)
}

func (r Repository) FindOne(ctx context.Context, id int64) (entity.Specimen, error) {
	var specimen entity.Specimen
	err := r.db.
		Where("id = ?", id).
		Preload("ObservationRequest").
		Preload("ObservationRequest.TestType").
		Preload("ObservationResult").
		Preload("ObservationResult.TestType").
		Preload("Patient").
		Preload("WorkOrder").
		First(&specimen).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Specimen{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Specimen{}, fmt.Errorf("error finding Specimen: %w", err)
	}

	return specimen, nil
}

func (r Repository) FindByBarcode(ctx context.Context, barcode string) (entity.Specimen, error) {
	var specimen entity.Specimen
	err := r.db.Debug().WithContext(ctx).
		Where("barcode = ?", barcode).
		Preload("ObservationResult").
		Preload("ObservationResult.TestType").
		Preload("WorkOrder").
		First(&specimen).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Specimen{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Specimen{}, fmt.Errorf("error finding Specimen: %w", err)
	}

	return specimen, nil
}

func (r *Repository) GenerateBarcode(ctx context.Context, specimenType entity.SpecimenType) string {
	seq := r.GetBarcodeSequence(ctx)
	seqPadding := fmt.Sprintf("%06d", seq) // Prints to stdout '000012'

	return fmt.Sprintf("%s_%s%s", specimenType.Code(), time.Now().Format("20060102"), seqPadding)
}

func (r *Repository) GetBarcodeSequence(ctx context.Context) int64 {
	seq, ok := r.cache.Get(constant.KeySpecimenBarcodeSequence)
	if !ok {
		now := time.Now()
		tomorrowMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)
		expire := tomorrowMidnight.Sub(now)
		r.cache.Set(constant.KeySpecimenBarcodeSequence, int64(1), expire)

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

func (r *Repository) IncrementBarcodeSequence(ctx context.Context) error {
	err := r.cache.Increment(constant.KeySpecimenBarcodeSequence, int64(1))
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SyncBarcodeSequence(ctx context.Context) error {
	now := time.Now()
	currentDayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	tomorrowMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.Local)

	var count int64
	err := r.db.Model(entity.Specimen{}).
		Where("created_at >= ? and created_at < ?", currentDayMidnight, tomorrowMidnight).
		Count(&count).Error
	if err != nil {
		return err
	}

	expire := tomorrowMidnight.Sub(now)
	r.cache.Set(constant.KeySpecimenBarcodeSequence, int64(count+1), expire)

	return nil
}
