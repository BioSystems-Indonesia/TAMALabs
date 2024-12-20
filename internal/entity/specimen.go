package entity

import "time"

type Specimen struct {
	ID             int       `json:"id" gorm:"primaryKey;autoIncrement"`
	HL7ID          string    `json:"specimen_hl7_id" gorm:"not null"`                        // SPM-2
	PatientID      int       `json:"patient_id" gorm:"not null" validate:"required"`         // Foreign key linking to Patient
	Type           string    `json:"type" gorm:"not null" validate:"required,specimen-type"` // SPM-4
	CollectionDate string    `json:"collection_date" gorm:"not null"`                        // SPM-17
	ReceivedDate   time.Time `json:"received_date" gorm:"not null"`
	Source         string    `json:"source" gorm:"not null"`    // SPM-8
	Condition      string    `json:"condition" gorm:"not null"` // SPM-25
	Method         string    `json:"method" gorm:"not null"`    // SPM-10
	Comments       string    `json:"comments" gorm:"not null"`  // SPM-26
	Barcode        string    `json:"barcode" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null"`

	// Relationships
	Patient Patient `json:"patient" gorm:"foreignKey:PatientID;->" validate:"-"`
}

type SpecimenGetManyRequest struct {
	GetManyRequest
	PatientID int64 `query:"patient_id"`
}
