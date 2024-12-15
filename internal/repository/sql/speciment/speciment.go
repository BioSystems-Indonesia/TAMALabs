package specimentrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SpecimentRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewSpecimentRepository(db *gorm.DB, cfg *config.Schema) *SpecimentRepository {
	return &SpecimentRepository{db: db, cfg: cfg}
}

func (r SpecimentRepository) FindAll(ctx context.Context, req *entity.SpecimenGetManyRequest) ([]entity.Specimen, error) {
	var speciments []entity.Specimen

	db := r.db.WithContext(ctx)
	if len(req.ID) > 0 {
		db = db.Where("id in (?)", req.ID)
	}

	if req.Query != "" {
		db = db.Where("description like ?", req.Query+"%")
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

	err := db.Find(&speciments).Error
	if err != nil {
		return nil, fmt.Errorf("error finding speciments: %w", err)
	}
	return speciments, nil
}

func (r SpecimentRepository) FindOne(id int64) (entity.Specimen, error) {
	var speciment entity.Specimen
	err := r.db.Where("id = ?", id).First(&speciment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Specimen{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Specimen{}, fmt.Errorf("error finding speciment: %w", err)
	}

	return speciment, nil
}

func (r SpecimentRepository) Create(speciment *entity.Specimen) error {
	return r.db.Create(speciment).Error
}

func (r SpecimentRepository) Update(speciment *entity.Specimen) error {
	res := r.db.Save(speciment).Error
	if res != nil {
		return fmt.Errorf("error updating speciment: %w", res)
	}

	return nil
}

func (r SpecimentRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Specimen{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting speciment: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
