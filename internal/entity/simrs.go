package entity

import "time"

// SimrsPatient represents the patients table entity in SIMRS database
type SimrsPatient struct {
	ID        int64     `json:"id" db:"id"`
	PatientID string    `json:"patient_id" db:"patient_id"`
	FirstName string    `json:"first_name" db:"first_name"`
	LastName  string    `json:"last_name" db:"last_name"`
	Birthdate time.Time `json:"birthdate" db:"birthdate"`
	Gender    string    `json:"gender" db:"gender"` // 'L' or 'P'
	Address   string    `json:"address" db:"address"`
	Phone     string    `json:"phone" db:"phone"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// SimrsLabRequest represents the lab_requests table entity in SIMRS database
type SimrsLabRequest struct {
	ID           int64     `json:"id" db:"id"`
	NoOrder      string    `json:"no_order" db:"no_order"`
	PatientID    string    `json:"patient_id" db:"patient_id"`
	ParamRequest string    `json:"param_request" db:"param_request"` // JSON string containing array of parameter codes
	RequestedBy  string    `json:"requested_by" db:"requested_by"`
	RequestedAt  time.Time `json:"requested_at" db:"requested_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// SimrsLabResult represents the lab_results table entity in SIMRS database
type SimrsLabResult struct {
	ID          int64     `json:"id" db:"id"`
	NoOrder     string    `json:"no_order" db:"no_order"`
	ParamCode   string    `json:"param_code" db:"param_code"`
	ResultValue string    `json:"result_value" db:"result_value"`
	Unit        string    `json:"unit" db:"unit"`
	RefRange    string    `json:"ref_range" db:"ref_range"`
	Flag        string    `json:"flag" db:"flag"` // 'H', 'L', 'N'
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// SimrsGender represents gender enum for SIMRS
type SimrsGender string

const (
	SimrsGenderMale   SimrsGender = "L"
	SimrsGenderFemale SimrsGender = "P"
)

// SimrsFlag represents result flag enum for SIMRS
type SimrsFlag string

const (
	SimrsFlagHigh   SimrsFlag = "H"
	SimrsFlagLow    SimrsFlag = "L"
	SimrsFlagNormal SimrsFlag = "N"
)

// NewSimrsFlag converts internal TestResult to SimrsFlag
func NewSimrsFlag(result TestResult) SimrsFlag {
	switch result.Abnormal {
	case HighResult:
		return SimrsFlagHigh
	case LowResult:
		return SimrsFlagLow
	case NormalResult:
		return SimrsFlagNormal
	case NoDataResult:
		return SimrsFlagNormal
	default:
		return SimrsFlagNormal
	}
}

// NewSimrsGender converts internal PatientSex to SimrsGender
func NewSimrsGender(sex PatientSex) SimrsGender {
	switch sex {
	case PatientSexMale:
		return SimrsGenderMale
	case PatientSexFemale:
		return SimrsGenderFemale
	default:
		return SimrsGenderMale
	}
}

// ToPatientSex converts SimrsGender to internal PatientSex
func (g SimrsGender) ToPatientSex() PatientSex {
	switch g {
	case SimrsGenderMale:
		return PatientSexMale
	case SimrsGenderFemale:
		return PatientSexFemale
	default:
		return PatientSexMale
	}
}
