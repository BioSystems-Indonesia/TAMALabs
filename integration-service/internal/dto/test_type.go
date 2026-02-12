package dto

type TestTypeRequest struct {
	Id            string `json:"id"`
	LabId         string `json:"lab_id"`
	SpecimentType string `json:"speciment_type"`
	CategoryName  string `json:"category_name"`
	Code          string `json:"code"`
	Unit          string `json:"unit"`
	Ref           string `json:"ref"`
	Decimal       int    `json:"decimal"`
}
