package entity

// OUL_R22 is a struct that contains the HL7 message, patient details, specimen details, observation request, and list of observations.
type OUL_R22 struct {
	HL7Message         HL7Message         `json:"hl7_message"`         // Metadata from MSH
	Patient            Patient            `json:"patient"`             // Patient details from PID
	Specimen           Specimen           `json:"specimen"`            // Specimen details from SPM
	ObservationRequest ObservationRequest `json:"observation_request"` // Observation request from OBR
	Observations       []Observation      `json:"observations"`        // List of observations from OBX
}
