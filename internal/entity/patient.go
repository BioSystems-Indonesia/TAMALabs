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
	ID                  int64          `json:"id,omitempty" gorm:"primaryKey;autoIncrement"`
	SIMRSPID            sql.NullString `json:"simrs_pid" gorm:"null;column:simrs_pid;uniqueIndex:idx_patient_simrs_pid"`
	MedicalRecordNumber string         `json:"medical_record_number" gorm:"column:medical_record_number;type:varchar(255);default:'';index:idx_patient_medical_record_number"`
	FirstName           string         `json:"first_name" gorm:"not null" validate:"required"`
	LastName            string         `json:"last_name" gorm:"not null" validate:""`
	Birthdate           time.Time      `json:"birthdate" gorm:"not null" validate:"required"`
	Sex                 PatientSex     `json:"sex" gorm:"not null" validate:"required,sex"`
	PhoneNumber         string         `json:"phone_number" gorm:"not null" validate:""`
	Location            string         `json:"location" gorm:"not null" validate:""`
	Address             string         `json:"address" gorm:"not null" validate:""`
	CreatedAt           time.Time      `json:"created_at" gorm:"not null"`
	UpdatedAt           time.Time      `json:"updated_at" gorm:"not null"`

	Specimen []Specimen `json:"specimen_list" gorm:"foreignKey:PatientID;->" validate:"-"`
}

// GetAge calculates the patient's age in years (as float64 for precision)
func (p *Patient) GetAge() float64 {
	now := time.Now()
	years := now.Year() - p.Birthdate.Year()

	// Adjust if birthday hasn't occurred this year yet
	if now.YearDay() < p.Birthdate.YearDay() {
		years--
	}

	// Calculate fractional part (for more precision if needed)
	daysSinceBirthday := now.YearDay() - p.Birthdate.YearDay()
	if daysSinceBirthday < 0 {
		daysInYear := 365
		if isLeapYear(now.Year()) {
			daysInYear = 366
		}
		daysSinceBirthday += daysInYear
	}

	fraction := float64(daysSinceBirthday) / 365.25
	return float64(years) + fraction
}

// GetGenderString returns gender as string pointer for reference range matching
func (p *Patient) GetGenderString() *string {
	if p.Sex == PatientSexMale {
		gender := "M"
		return &gender
	}
	if p.Sex == PatientSexFemale {
		gender := "F"
		return &gender
	}
	return nil // Unknown gender
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

type GetManyRequestPatient struct {
	GetManyRequest

	BirthDate time.Time `query:"birthdate"`
}

type GetPatientRecordHistoryRequest struct {
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`
}

type GetPatientResultHistoryResponse struct {
	Patient    Patient      `json:"patient"`
	TestResult []TestResult `json:"test_result"`
}
