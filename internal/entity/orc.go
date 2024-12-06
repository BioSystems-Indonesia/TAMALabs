package entity

// ORC (Common Order) segment (used in ORM requests)
type ORC struct {
	OrderControl      string `json:"order_control" hl7:"1"`
	OrderID           string `json:"order_id" hl7:"2"`
	PlacerOrderNumber string `json:"placer_order_number" hl7:"3"`
	OrderStatus       string `json:"order_status" hl7:"5"`
	OrderPriority     string `json:"order_priority" hl7:"6"`
	OrderDateTime     string `json:"order_date_time" hl7:"9"`
}
