package specimen

import (
	"context"
	"errors"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewRepository(db *gorm.DB, cfg *config.Schema) *Repository {
	return &Repository{db: db, cfg: cfg}
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

	err := db.Find(&Specimens).Error
	if err != nil {
		return nil, fmt.Errorf("error finding Specimens: %w", err)
	}
	return Specimens, nil
}

func (r Repository) FindOne(id int64) (entity.Specimen, error) {
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

func (r Repository) Create(Specimen *entity.Specimen) error {
	return r.db.Create(Specimen).Error
}

func (r Repository) Update(Specimen *entity.Specimen) error {
	res := r.db.Save(Specimen).Error
	if res != nil {
		return fmt.Errorf("error updating Specimen: %w", res)
	}

	return nil
}

func (r Repository) Delete(id int) error {
	res := r.db.Delete(&entity.Specimen{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting Specimen: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
