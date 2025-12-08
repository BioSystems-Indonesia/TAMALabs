package entity

import "time"

// SimrsLabOrder represents the lab_order table entity in Database Sharing SIMRS database
type SimrsLabOrder struct {
	ID          int       `json:"id" db:"id"`
	NoLabOrder  string    `json:"no_lab_order" db:"no_lab_order"`
	NoRM        string    `json:"no_rm" db:"no_rm"`
	PatientName string    `json:"patient_name" db:"patient_name"`
	BirthDate   time.Time `json:"birth_date" db:"birth_date"`
	Sex         string    `json:"sex" db:"sex"` // 'M' or 'F'
	Doctor      string    `json:"doctor" db:"doctor"`
	Analyst     string    `json:"analyst" db:"analyst"`
	Status      string    `json:"status" db:"status"` // 'NEW', 'PENDING', 'LIS_SUCCESS', 'SIMRS_SUCCESS'
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// SimrsOrderDetail represents the order_detail table entity in Database Sharing SIMRS database
type SimrsOrderDetail struct {
	ID             int       `json:"id" db:"id"`
	NoLabOrder     string    `json:"no_lab_order" db:"no_lab_order"`
	ParameterCode  string    `json:"parameter_code" db:"parameter_code"`
	ParameterName  string    `json:"parameter_name" db:"parameter_name"`
	ResultValue    string    `json:"result_value" db:"result_value"`
	Unit           string    `json:"unit" db:"unit"`
	ReferenceRange string    `json:"reference_range" db:"reference_range"`
	Flag           string    `json:"flag" db:"flag"` // 'H', 'L', 'N', or empty
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// SimrsOrderStatus represents the status of a lab order in Database Sharing workflow
type SimrsOrderStatus string

const (
	SimrsStatusNew          SimrsOrderStatus = "NEW"           // Order inserted by SIMRS
	SimrsStatusPending      SimrsOrderStatus = "PENDING"       // Order fetched by LIS, waiting to be processed
	SimrsStatusLISSuccess   SimrsOrderStatus = "LIS_SUCCESS"   // LIS completed the test results
	SimrsStatusSIMRSSuccess SimrsOrderStatus = "SIMRS_SUCCESS" // SIMRS successfully fetched the results
)

// SimrsSex represents sex/gender enum for Database Sharing
type SimrsSex string

const (
	SimrsSexMale   SimrsSex = "M"
	SimrsSexFemale SimrsSex = "F"
)

// ToPatientSex converts SimrsSex to internal PatientSex
func (s SimrsSex) ToPatientSex() PatientSex {
	switch s {
	case SimrsSexMale:
		return PatientSexMale
	case SimrsSexFemale:
		return PatientSexFemale
	default:
		return PatientSexUnknown
	}
}

// NewSimrsSex converts internal PatientSex to SimrsSex
func NewSimrsSex(sex PatientSex) SimrsSex {
	switch sex {
	case PatientSexMale:
		return SimrsSexMale
	case PatientSexFemale:
		return SimrsSexFemale
	default:
		return SimrsSexMale
	}
}
