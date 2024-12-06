package entity

// MSH (Message Header) segment for both request and response
type MSH struct {
	FieldSeparator       string `json:"field_separator" hl7:"1"`
	SendingApplication   string `json:"sending_application" hl7:"2"`
	SendingFacility      string `json:"sending_facility" hl7:"3"`
	ReceivingApplication string `json:"receiving_application" hl7:"4"`
	ReceivingFacility    string `json:"receiving_facility" hl7:"5"`
	MessageDateTime      string `json:"message_date_time" hl7:"6"`
	MessageType          string `json:"message_type" hl7:"8"`
	MessageControlID     string `json:"message_control_id" hl7:"9"`
	ProcessingID         string `json:"processing_id" hl7:"10"`
	VersionID            string `json:"version_id" hl7:"11"`
}
