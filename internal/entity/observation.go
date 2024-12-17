package entity

import "time"

type ObservationRequest struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecimenID      int       `json:"specimen_id"`      // Foreign key linking to Specimen
	OrderID         string    `json:"order_id"`         // OBR-2
	TestCode        string    `json:"test_code"`        // OBR-4
	TestDescription string    `json:"test_description"` // OBR-4
	RequestedDate   time.Time `json:"requested_date"`   // OBR-7
	ResultStatus    string    `json:"result_status"`    // OBR-25

}

type Observation struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestID      int       `json:"request_id"`      // Foreign key linking to ObservationRequest
	Code           string    `json:"code"`            // OBX-3
	Description    string    `json:"description"`     // OBX-3
	Value          string    `json:"value"`           // OBX-5
	Unit           string    `json:"unit"`            // OBX-6
	ReferenceRange string    `json:"reference_range"` // OBX-7
	Date           time.Time `json:"date"`            // OBX-14
	AbnormalFlag   string    `json:"abnormal_flag"`   // OBX-8
	Comments       string    `json:"comments"`        // OBX-16
}
