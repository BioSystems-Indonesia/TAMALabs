package entity

import "time"

type ObservationRequest struct {
	ID              int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TestCode        string    `json:"test_code"`
	TestDescription string    `json:"test_description"`
	RequestedDate   time.Time `json:"requested_date"`
	ResultStatus    string    `json:"result_status"`
	SpecimenID      int64     `json:"specimen_id" gorm:"not null;index:observation_request_uniq,unique,priority:2" validate:"required"`
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`
}

type ObservationResult struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecimenID     int64     `json:"specimen_id" gorm:"not null;index:observation_result_uniq,unique,priority:2" validate:"required"`
	Code           string    `json:"code"`
	Description    string    `json:"description"`
	Values         []string  `json:"values" gorm:"type:json"` // Using JSON for the slice
	Type           string    `json:"type"`
	Unit           string    `json:"unit"`
	ReferenceRange string    `json:"reference_range"`
	Date           time.Time `json:"date"`
	AbnormalFlag   []string  `json:"abnormal_flag" gorm:"type:json"` // Using JSON for the slice
	Comments       string    `json:"comments"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`
}

type ObservationRequestGetManyRequest struct {
	GetManyRequest
}
