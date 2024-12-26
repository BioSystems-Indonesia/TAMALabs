package specimen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r Repository) FindAll(ctx context.Context, req *entity.SpecimenGetManyRequest) ([]entity.Specimen, error) {
	var Specimens []entity.Specimen

	db := r.db.WithContext(ctx)
	if len(req.ID) > 0 {
		db = db.Where("id in (?)", req.ID)
	}

	if req.PatientID != 0 {
		db = db.Where("patient_id = ?", req.PatientID)
	}

	if req.Sort != "" {
		db = db.Order(clause.OrderByColumn{
			Column: clause.Column{
				Name: req.Sort,
			},
			Desc: req.IsSortDesc(),
		})
	}

	err := db.Preload("ObservationRequest").Find(&Specimens).Error
	if err != nil {
		return nil, fmt.Errorf("error finding Specimens: %w", err)
	}
	return Specimens, nil
}

func (r Repository) FindOne(ctx context.Context, id int64) (entity.Specimen, error) {
	var Specimen entity.Specimen
	err := r.db.Where("id = ?", id).First(&Specimen).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Specimen{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Specimen{}, fmt.Errorf("error finding Specimen: %w", err)
	}

	return Specimen, nil
}

func (r *Repository) GenerateBarcode(ctx context.Context) (string, error) {
	seq, err := r.GetBarcodeSequence(ctx)
	if err != nil {
		return "", err
	}

	seqPadding := fmt.Sprintf("%06d", seq) // Prints to stdout '000012'

	return fmt.Sprintf("%s%s", time.Now().Format("20060102"), seqPadding), nil
}

func (r *Repository) GetBarcodeSequence(ctx context.Context) (int64, error) {
	seq, ok := r.cache.Get(constant.KeySpecimenBarcodeSequence)
	if !ok {
		return 0, errors.New("barcode sequence not found")
	}

	return seq.(int64), nil
}

func (r *Repository) IncrementBarcodeSequence(ctx context.Context) error {
	err := r.cache.Increment(constant.KeySpecimenBarcodeSequence, 1)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) SyncBarcodeSequence(ctx context.Context) error {
	now := time.Now()
	currentDayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, constant.LocationJakarta)
	tomorrowMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, constant.LocationJakarta)

	var count int64
	err := r.db.Model(entity.Specimen{}).
		Where("created_at >= ? and created_at < ?", currentDayMidnight, tomorrowMidnight).
		Count(&count).Error
	if err != nil {
		return err
	}

	expire := tomorrowMidnight.Sub(now)
	r.cache.Set(constant.KeySpecimenBarcodeSequence, count+1, expire)

	return nil
}
