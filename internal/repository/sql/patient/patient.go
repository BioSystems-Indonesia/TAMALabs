package patientrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PatientRepository struct {
	db  *gorm.DB
	cfg *config.Schema
}

func NewPatientRepository(db *gorm.DB, cfg *config.Schema) *PatientRepository {
	return &PatientRepository{db: db, cfg: cfg}
}

func (r PatientRepository) FindAll(ctx context.Context, req *entity.GetManyRequest) ([]entity.Patient, error) {
	var patients []entity.Patient

	db := r.db.WithContext(ctx)
	db = db.Order(clause.OrderByColumn{
		Column: clause.Column{
			Name: req.Sort,
		},
		Desc: req.IsSortDesc(),
	})

	err := db.Find(&patients).Error
	if err != nil {
		return nil, fmt.Errorf("error finding patients: %w", err)
	}
	return patients, nil
}

func (r PatientRepository) FindOne(id int64) (entity.Patient, error) {
	var patient entity.Patient
	err := r.db.Where("id = ?", id).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.Patient{}, entity.ErrNotFound
	}

	if err != nil {
		return entity.Patient{}, fmt.Errorf("error finding patient: %w", err)
	}

	return patient, nil
}

func (r PatientRepository) Create(patient *entity.Patient) error {
	return r.db.Create(patient).Error
}

func (r PatientRepository) Update(patient *entity.Patient) error {
	res := r.db.Save(patient).Error
	if res != nil {
		return fmt.Errorf("error updating patient: %w", res)
	}

	return nil
}

func (r PatientRepository) Delete(id int64) error {
	res := r.db.Delete(&entity.Patient{ID: id})
	if res.Error != nil {
		return fmt.Errorf("error deleting patient: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return entity.ErrNotFound
	}

	return nil
}
