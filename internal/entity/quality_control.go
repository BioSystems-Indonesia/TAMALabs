package entity

import "time"

type QualityControl struct {
	ID               int       `json:"id" gorm:"primaryKey"`
	DeviceID         int       `json:"device_id" gorm:"not null"`
	TestTypeID       int       `json:"test_type_id" gorm:"not null"`
	QCLevel          int       `json:"qc_level" gorm:"not null"` // 1, 2, or 3
	LotNumber        string    `json:"lot_number" gorm:"not null"`
	RefValue         float64   `json:"ref_value" gorm:"not null"`
	SDValue          float64   `json:"sd_value" gorm:"not null"`
	MeasuredValue    float64   `json:"measured_value" gorm:"not null"`
	CVValue          float64   `json:"cv_value" gorm:"not null"` // (SD/REF)*100
	Result           string    `json:"result" gorm:"not null"`   // Pass or Fail
	Operator         string    `json:"operator" gorm:"not null"`
	DeviceIdentifier string    `json:"device_identifier" gorm:"column:device_identifier"`
	MessageControlID string    `json:"message_control_id" gorm:"column:message_control_id"` // HL7 message ID for deduplication
	CreatedAt        time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relations
	Device   *Device   `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	TestType *TestType `json:"test_type,omitempty" gorm:"foreignKey:TestTypeID"`
}

func (QualityControl) TableName() string {
	return "quality_controls"
}

type GetManyRequestQualityControl struct {
	GetManyRequest
	DeviceID   *int `query:"device_id"`
	TestTypeID *int `query:"test_type_id"`
	QCLevel    *int `query:"qc_level"`
}
