package models

type TestType struct {
	ID                int
	Name              string
	Code              string
	AliasCode         string
	Unit              string
	LowRefRange       float64
	HighRefRange      float64
	NormalRefString   string
	Decimal           int
	Category          string
	SubCategory       string
	Description       string
	IsCalculatedTest  bool
	DeviceID          *int
	Type              string
	SpecificRefRanges []string `gorm:"serializer:json"`
	LoincCode         string
	AlternativeCodes  []string `gorm:"serializer:json"`
}
