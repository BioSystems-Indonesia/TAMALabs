package entity

import "time"

type Observation struct {
	Request ObservationRequest  `json:"request"`
	Result  []ObservationResult `json:"result"`
}

type ObservationRequest struct {
	ID              int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecimenID      int       `json:"specimen_id" gorm:"not null;index:observation_request_uniq,unique,priority:2" validate:"required"`                // Foreign key linking to Specimen
	OrderID         string    `json:"order_id" gorm:"not null;index:observation_request_uniq,unique,priority:1"`                                       // OBR-2
	TestCode        string    `json:"test_code" gorm:"not null;index:observation_request_uniq,unique,priority:3" validate:"required,observation-type"` // OBR-4
	TestDescription string    `json:"test_description" gorm:"not null"`                                                                                // OBR-4
	RequestedDate   time.Time `json:"requested_date" gorm:"not null"`                                                                                  // OBR-7
	ResultStatus    string    `json:"result_status" gorm:"not null"`                                                                                   // OBR-25
	CreatedAt       time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"not null"`

	Order    WorkOrder `json:"order" gorm:"foreignKey:OrderID;->" validate:"-"`
	Specimen Specimen  `json:"specimen" gorm:"foreignKey:SpecimenID;->" validate:"-"`
}

type ObservationResult struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	RequestID      int       `json:"request_id"`      // Foreign key linking to ObservationRequest
	Code           string    `json:"code"`            // OBX-3
	Description    string    `json:"description"`     // OBX-3
	Value          []string  `json:"value"`           // OBX-5
	Type           string    `json:"type"`            // OBX-2
	Unit           string    `json:"unit"`            // OBX-6
	ReferenceRange string    `json:"reference_range"` // OBX-7
	Date           time.Time `json:"date"`            // OBX-14
	AbnormalFlag   []string  `json:"abnormal_flag"`   // OBX-8
	Comments       string    `json:"comments"`        // OBX-16
}

type ObservationRequestGetManyRequest struct {
	GetManyRequest
}
