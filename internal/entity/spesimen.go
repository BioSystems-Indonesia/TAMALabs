package entity

import "time"

type Specimen struct {
	ID                     int       `json:"id" gorm:"primaryKey;autoIncrement"`
	SpecimenHL7ID          string    `json:"specimen_hl7_id"`          // SPM-2
	PatientID              int       `json:"patient_id"`               // Foreign key linking to Patient
	SpecimenType           string    `json:"specimen_type"`            // SPM-4
	SpecimenCollectionDate string    `json:"specimen_collection_date"` // SPM-17
	SpecimenReceivedDate   time.Time `json:"specimen_received_date"`
	SpecimenSource         string    `json:"specimen_source"`    // SPM-8
	SpecimenCondition      string    `json:"specimen_condition"` // SPM-25
	CollectionMethod       string    `json:"collection_method"`  // SPM-10
	Comments               string    `json:"comments"`           // SPM-26
}
