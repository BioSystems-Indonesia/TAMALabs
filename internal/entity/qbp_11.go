package entity

// QBP_Q11 is a struct that contains the HL7 message.
type QBP_Q11 struct {
	Msh MSH `json:"msh"`
	QPD QPD `json:"qpd"`
}

// QPD is a struct that contains the HL7 message.
type QPD struct {
	QueryTag string `json:"query_tag"`
	Barcode  string `json:"barcode"`
}
