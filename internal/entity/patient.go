package entity

import (
	"database/sql"
	"time"
)

type PatientSex string

const (
	PatientSexMale    PatientSex = "M"
	PatientSexFemale  PatientSex = "F"
	PatientSexUnknown PatientSex = "U"
)

func (p PatientSex) String() string {
	switch p {
	case PatientSexMale:
		return "Male"
	case PatientSexFemale:
		return "Female"
	case PatientSexUnknown:
		return "Unknown"
	default:
		return "Undefined"
	}
}

func NewPatientSexFromKhanza(khanzaPatientSex KhanzaPatientSex) PatientSex {
	switch khanzaPatientSex {
	case KhanzaPatientSexMale:
		return PatientSexMale
	case KhanzaPatientSexFemale:
		return PatientSexFemale
	default:
		return PatientSexUnknown
	}
}

type Patient struct {
	ID          int64          `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	SIMRSPID    sql.NullString `json:"simrs_pid" gorm:"null;column:simrs_pid;uniqueIndex:idx_patient_simrs_pid"`
	FirstName   string         `json:"first_name" gorm:"not null" validate:"required"`
	LastName    string         `json:"last_name" gorm:"not null" validate:""`
	Birthdate   time.Time      `json:"birthdate" gorm:"not null" validate:"required"`
	Sex         PatientSex     `json:"sex" gorm:"not null" validate:"required,sex"`
	PhoneNumber string         `json:"phone_number" gorm:"not null" validate:""`
	Location    string         `json:"location" gorm:"not null" validate:""`
	Address     string         `json:"address" gorm:"not null" validate:""`
	CreatedAt   time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"not null"`

	Specimen []Specimen `json:"specimen_list" gorm:"foreignKey:PatientID;->" validate:"-"`
}

type GetManyRequestPatient struct {
	GetManyRequest

	BirthDate time.Time `query:"birthdate"`
}
