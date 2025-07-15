package entity

type Unit struct {
	Base  string `json:"base" gorm:"uniqueIndex:units_base_value_uniq"`
	Value string `json:"value" gorm:"uniqueIndex:units_base_value_uniq"`
}

type UnitGetManyRequest struct {
	GetManyRequest

	Base  string `query:"base"`
	Value string `query:"value"`
}
