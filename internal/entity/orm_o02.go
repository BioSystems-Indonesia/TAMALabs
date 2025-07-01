package entity

type ORM_O01 struct {
	Msh    MSH `json:"msh"`
	Orders []Order_ORM_O01
}

// QPD is a struct that contains the HL7 message.
type Order_ORM_O01 struct {
	Barcode string `json:"barcode"`
}
