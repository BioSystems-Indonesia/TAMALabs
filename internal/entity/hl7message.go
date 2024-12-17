package entity

type HL7Message struct {
	ID                   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	MessageControlID     string `json:"message_control_id"`    // MSH-10
	SendingApplication   string `json:"sending_application"`   // MSH-3
	SendingFacility      string `json:"sending_facility"`      // MSH-4
	ReceivingApplication string `json:"receiving_application"` // MSH-5
	ReceivingFacility    string `json:"receiving_facility"`    // MSH-6
	MessageType          string `json:"message_type"`          // MSH-9
	MessageVersion       string `json:"message_version"`       // MSH-12
	MessageDate          string `json:"message_date"`          // MSH-7
}
