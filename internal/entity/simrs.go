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
	case PositiveResult:
		return SimrsFlagHigh // Positive could be considered as abnormal/high
	case NegativeResult:
		return SimrsFlagNormal // Negative could be considered as normal
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

// SimrsOrderRequest represents the incoming order request from SIMRS via API
type SimrsOrderRequest struct {
	Order SimrsOrder `json:"order"`
}

type SimrsOrder struct {
	PID SimrsPID `json:"pid"`
	OBR SimrsOBR `json:"obr"`
}

type SimrsPID struct {
	Pname               string `json:"pname"`    // Patient name
	Sex                 string `json:"sex"`      // M/F
	BirthDt             string `json:"birth_dt"` // Format: dd.MM.yyyy (21.07.1985)
	MedicalRecordNumber string `json:"no_rm"`    // Medical Record Number (No RM) (optional)

}

type SimrsOBR struct {
	OrderLab  string  `json:"order_lab"`  // Lab order number/barcode
	OrderTest []int   `json:"order_test"` // Array of test type IDs (primary key)
	Doctor    []int64 `json:"doctor"`     // Array of doctor IDs (optional)
	Analyst   []int64 `json:"analyst"`    // Array of analyzer IDs (optional)
}

// SimrsResultResponse represents the response for GET result endpoint
type SimrsResultResponse struct {
	Response SimrsResultResponseData `json:"response"`
	Result   SimrsResultOBX          `json:"result"`
}

type SimrsResultResponseData struct {
	Sample SimrsResultSample `json:"sampel"` // "sampel" is intentional (Indonesian)
}

type SimrsResultSample struct {
	ResultTest []SimrsResultTest `json:"result_test"`
}

type SimrsResultTest struct {
	Loinc       string `json:"loinc"`
	TestID      int    `json:"test_id"`      // Test code
	NamaTest    string `json:"nama_test"`    // Test name
	Hasil       string `json:"hasil"`        // Result value
	NilaiNormal string `json:"nilai_normal"` // Reference range
	Satuan      string `json:"satuan"`       // Unit
	Flag        string `json:"flag"`         // H/L/N
}

type SimrsResultOBX struct {
	OBX SimrsResultOBXData `json:"obx"`
}

type SimrsResultOBXData struct {
	OrderLab string `json:"order_lab"`
}
