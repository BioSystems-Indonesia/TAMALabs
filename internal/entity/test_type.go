package entity

import (
	"strings"

	"gorm.io/gorm"
)

type TestType struct {
	ID               int     `json:"id" gorm:"primaryKey"`
	Name             string  `json:"name"`
	Code             string  `json:"code" gorm:"unique"`
	Unit             string  `json:"unit"`
	LowRefRange      float64 `json:"low_ref_range"`
	HighRefRange     float64 `json:"high_ref_range"`
	Decimal          int     `json:"decimal"`
	Category         string  `json:"category"`
	SubCategory      string  `json:"sub_category"`
	Description      string  `json:"description"`
	IsCalculatedTest bool    `json:"is_calculated_test" gorm:"column:is_calculated_test;default:false"`

	Type []TestTypeSpecimenType `json:"types" gorm:"-"`
	// TypeDB is a specimen type separated by comma
	TypeDB string `json:"-" gorm:"column:type"`
}

// TestTypeSpecimenType is a specimen type struct, we need it to conform with react-admin standard
// for ordered input
type TestTypeSpecimenType struct {
	Type string `json:"type"`
}

func (t *TestType) BeforeCreate(tx *gorm.DB) error {
	types := make([]string, len(t.Type))
	for i, specimenType := range t.Type {
		types[i] = specimenType.Type
	}

	t.TypeDB = strings.Join(types, ",")

	return nil
}

func (t *TestType) BeforeUpdate(tx *gorm.DB) error {
	types := make([]string, len(t.Type))
	for i, specimenType := range t.Type {
		types[i] = specimenType.Type
	}

	t.TypeDB = strings.Join(types, ",")

	return nil
}

func (t *TestType) AfterFind(tx *gorm.DB) error {
	types := strings.Split(t.TypeDB, ",")

	t.Type = make([]TestTypeSpecimenType, 0)
	for _, specimenType := range types {
		if specimenType == "" {
			continue
		}

		t.Type = append(t.Type, TestTypeSpecimenType{Type: specimenType})
	}

	if len(t.Type) == 0 {
		t.Type = append(t.Type, TestTypeSpecimenType{Type: "SER"})
	}

	return nil
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
