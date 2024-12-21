package entity

// OUL_R22 is a struct that contains the HL7 message, patient details, specimen details, observation request, and list of observations.
type OUL_R22 struct {
	Msh       MSH        `json:"msh"`       // Metadata from MSH
	Patient   Patient    `json:"patient"`   // Patient details from PID
	Specimens []Specimen `json:"specimens"` // Specimens details from SPM
}
