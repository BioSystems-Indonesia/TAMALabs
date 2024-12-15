package entity

import "time"

type Specimen struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	HL7ID          string    `json:"specimen_hl7_id"`          // SPM-2
	PatientID      int       `json:"patient_id"`               // Foreign key linking to Patient
	Type           string    `json:"specimen_type"`            // SPM-4
	CollectionDate string    `json:"specimen_collection_date"` // SPM-17
	ReceivedDate   time.Time `json:"specimen_received_date"`
	Source         string    `json:"specimen_source"`    // SPM-8
	Condition      string    `json:"specimen_condition"` // SPM-25
	Method         string    `json:"collection_method"`  // SPM-10
	Comments       string    `json:"comments"`           // SPM-26

	Barcode   string    `json:"barcode" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	// Relationships
	Patient Patient `json:"patient" gorm:"foreignKey:PatientID" validate:"-"`
}

type SpecimenGetManyRequest struct {
	GetManyRequest
	PatientID int64 `query:"patient_id"`
}
