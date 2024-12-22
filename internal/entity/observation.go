package entity

import "time"

type Observation struct {
	ID int64 `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID         string    `json:"order_id" gorm:"not null;index:observation_request_uniq,unique,priority:1"`                                       // OBR-2
	SpecimenID           int64               `json:"specimen_id" gorm:"not null;index:observation_request_uniq,unique,priority:2" validate:"required"`                                     // Foreign key linking to Specimen
	ObservationRequestID int64               `json:"observation_request_id"` // Foreign key for ObservationRequest
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`
	Request              ObservationRequest  `json:"request" gorm:"foreignKey:ObservationRequestID"` // Correct foreign key for Request
	Result               []ObservationResult `json:"result" gorm:"foreignKey:ObservationID"`         // Correct foreign key for Result
	Order    WorkOrder `json:"order" gorm:"foreignKey:OrderID;->" validate:"-"`

}

type ObservationRequest struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TestCode        string    `json:"test_code"`
	TestDescription string    `json:"test_description"`
	RequestedDate   time.Time `json:"requested_date"`
	ResultStatus    string    `json:"result_status"`


	ObservationID int64 `json:"observation_id"` // Foreign key linking to Observation
}

type ObservationResult struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	Code           string    `json:"code"`
	Description    string    `json:"description"`
	Values         []string  `json:"values" gorm:"type:json"` // Using JSON for the slice
	Type           string    `json:"type"`
	Unit           string    `json:"unit"`
	ReferenceRange string    `json:"reference_range"`
	Date           time.Time `json:"date"`
	AbnormalFlag   []string  `json:"abnormal_flag" gorm:"type:json"` // Using JSON for the slice
	Comments       string    `json:"comments"`

	ObservationID int64 `json:"observation_id"` // Foreign key linking to Observation
}

type ObservationRequestGetManyRequest struct {
	GetManyRequest
}
