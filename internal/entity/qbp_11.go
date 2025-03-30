package entity

// QBP_Q11 is a struct that contains the HL7 message.
type QBP_Q11 struct {
	Msh     MSH    `json:"msh"`
	Barcode string `json:"barcode"`
}
