package entity

type TestType struct {
	ID           int     `json:"id" gorm:"primaryKey"`
	Name         string  `json:"name"`
	Code         string  `json:"code" gorm:"unique"`
	Unit         string  `json:"unit"`
	LowRefRange  float64 `json:"low_ref_range"`
	HighRefRange float64 `json:"high_ref_range"`
	Category     string  `json:"category"`
	SubCategory  string  `json:"sub_category"`
	Description  string  `json:"description"`
}

type TestTypeFilter struct {
	Categories    []string `json:"categories"`
	SubCategories []string `json:"sub_categories"`
}

type TestTypeGetManyRequest struct {
	GetManyRequest

	Code          string   `query:"code"`
	Categories    []string `query:"categories"`
	SubCategories []string `query:"subCategories"`
}
