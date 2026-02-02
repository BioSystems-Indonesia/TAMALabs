package patientrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql"
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

func (r PatientRepository) FindAll(
	ctx context.Context,
	req *entity.GetManyRequestPatient,
) (entity.PaginationResponse[entity.Patient], error) {
	db := r.db.WithContext(ctx)
	db = sql.ProcessGetMany(db, req.GetManyRequest,
		sql.Modify{
			ProcessSearch: func(db *gorm.DB, query string) *gorm.DB {
				return db.Where("first_name || ' ' || last_name like ?", "%"+query+"%")
			},
		})

	if !req.BirthDate.IsZero() {
		db = db.Where("date(birthdate) = ?", req.BirthDate.Format(time.DateOnly))
	}

	return sql.GetWithPaginationResponse[entity.Patient](db, req.GetManyRequest)
}

func (r PatientRepository) FindManyByWorkOrderID(
	ctx context.Context,
	workOrderIDs []int64,
) ([]entity.Patient, error) {
	var patientIDs []int64
	err := r.db.WithContext(ctx).Model(&entity.WorkOrder{}).
		Where("id in (?)", workOrderIDs).
		Pluck("patient_id", &patientIDs).Error
	if err != nil {
		return nil, err
	}

	var patients []entity.Patient
	err = r.db.WithContext(ctx).Where("id in (?)", patientIDs).
		Preload("Specimen", "order_id in (?)", workOrderIDs).
		Preload("Specimen.ObservationRequest").
		Preload("Specimen.ObservationRequest.TestType").
		Find(&patients).Error
	if err != nil {
		return nil, err
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

func (r PatientRepository) FindByID(ctx context.Context, id int64) (*entity.Patient, error) {
	var patient entity.Patient
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, entity.ErrNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error finding patient: %w", err)
	}

	return &patient, nil
}

func (r PatientRepository) FindByPatientID(ctx context.Context, patientID string) (*entity.Patient, error) {
	var patient entity.Patient
	err := r.db.WithContext(ctx).Where("patient_id = ?", patientID).First(&patient).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("error finding patient by patient_id: %w", err)
	}

	return &patient, nil
}

func (r PatientRepository) Create(patient *entity.Patient) error {
	return r.db.Create(patient).Error
}

func (r PatientRepository) CreateWithContext(ctx context.Context, patient *entity.Patient) error {
	return r.db.WithContext(ctx).Create(patient).Error
}

func (r PatientRepository) CreateManyFromSIMRS(patients []entity.Patient) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "simrs_pid"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name", "birthdate", "sex", "address"}),
	}).Create(patients).Error
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

func (r PatientRepository) IsExists(patient *entity.Patient) (bool, error) {
	var c int64
	err := r.db.
		Model(&entity.Patient{}).
		Where("first_name = ?", patient.FirstName).
		Where("last_name = ?", patient.LastName).
		Where("date(birthdate) = ?", patient.Birthdate.Format("2006-01-02")).
		Count(&c).Error

	return c > 0, err
}

func (r PatientRepository) FirstOrCreate(patient *entity.Patient) (entity.Patient, error) {
	err := r.db.
		Model(&entity.Patient{}).
		Where("first_name = ?", patient.FirstName).
		Where("last_name = ?", patient.LastName).
		Where("date(birthdate) = ?", patient.Birthdate.Format("2006-01-02")).
		FirstOrCreate(patient).
		Error

	return *patient, err
}

func (r PatientRepository) FindByMedicalRecordNumber(ctx context.Context, medicalRecordNumber string, patient *entity.Patient) error {
	err := r.db.WithContext(ctx).
		Where("medical_record_number = ? AND medical_record_number != ''", medicalRecordNumber).
		First(patient).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.ErrNotFound
	}

	return err
}
