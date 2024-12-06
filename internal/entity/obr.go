package entity

// OBR (Observation Request) segment (used in ORM requests)
type OBR struct {
	SetID              string `json:"set_id" hl7:"1"`
	PlacerOrderNumber  string `json:"placer_order_number" hl7:"2"`
	UniversalServiceID string `json:"universal_service_id" hl7:"4"`
	RequestDateTime    string `json:"request_date_time" hl7:"7"`
}
