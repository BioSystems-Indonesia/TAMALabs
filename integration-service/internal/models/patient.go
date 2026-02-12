package models

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

type Patient struct {
	ID                  int64
	FirstName           string
	LastName            string
	Birthdate           time.Time
	Sex                 PatientSex
	PhoneNumber         string
	Location            string
	Address             string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	SimrsPid            sql.NullString
	MedicalRecordNumber string
}
