package dto

import "time"

type PatientRequest struct {
	PatientID           string    `json:"patient_id"`
	LabID               string    `json:"lab_id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Gender              string    `json:"gender"`
	Birthdate           time.Time `json:"birthdate"`
	Address             string    `json:"address"`
	MedicalRecordNumber string    `json:"medical_record_number"`
}
