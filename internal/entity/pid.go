package entity

// PID (Patient Identification) segment (used in ORM requests)
type PID struct {
	PatientID   string `json:"patient_id" hl7:"3"`
	PatientName string `json:"patient_name" hl7:"5"`
	DateOfBirth string `json:"date_of_birth" hl7:"7"`
	Gender      string `json:"gender" hl7:"8"`
	Address     string `json:"address" hl7:"11"`
	PhoneNumber string `json:"phone_number" hl7:"13"`
}
