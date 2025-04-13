package specimen

import (
	"context"
	"errors"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/sql"
	"gorm.io/gorm"
)

type Repository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	r := &Repository{db: db, cfg: cfg}

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
		Preload("Patient").
		First(&specimen).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Specimen{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Specimen{}, fmt.Errorf("error finding Specimen: %w", err)
	}

	return specimen, nil
}
