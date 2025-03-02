package entity

type Unit struct {
	Base  string `json:"base"`
	Value string `json:"value"`
}

type UnitGetManyRequest struct {
	GetManyRequest

	Base  string `query:"base"`
	Value string `query:"value"`
}
