package models

import "time"

type WorkOrder struct {
	ID                     int64
	Status                 string
	PatientID              int64
	DeviceID               int64
	Barcode                string
	VerifiedStatus         string
	CreatedBy              int64
	LastUpdatedBy          int64
	CreatedAt              time.Time
	UpdatedAt              time.Time
	BarcodeSIMRS           string
	MedicalRecordNumber    string
	VisitNumber            string
	SpecimenCollectionDate *time.Time
	ResultReleaseDate      *time.Time
	Diagnosis              string
}
