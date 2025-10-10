package entity

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type TestType struct {
	ID               int     `json:"id" gorm:"primaryKey"`
	Name             string  `json:"name"`
	Code             string  `json:"code" gorm:"unique"`
	AliasCode        string  `json:"alias_code" gorm:"index:test_type_alias_code"`
	Unit             string  `json:"unit"`
	LowRefRange      float64 `json:"low_ref_range"`
	HighRefRange     float64 `json:"high_ref_range"`
	NormalRefString  string  `json:"normal_ref_string" gorm:"column:normal_ref_string"` // New field for string reference values
	Decimal          int     `json:"decimal"`
	Category         string  `json:"category"`
	SubCategory      string  `json:"sub_category"`
	Description      string  `json:"description"`
	IsCalculatedTest bool    `json:"is_calculated_test" gorm:"column:is_calculated_test;default:false"`
	DeviceID         *int    `json:"device_id" gorm:"index:test_type_device_id"`

	Type   []TestTypeSpecimenType `json:"types" gorm:"-"`
	Device *Device                `json:"device,omitempty" gorm:"foreignKey:DeviceID"`
	// TypeDB is a specimen type separated by comma
	TypeDB string `json:"-" gorm:"column:type"`
}

func (t *TestType) GetFirstType() string {
	if len(t.Type) == 0 {
		return "SER"
	}

	return t.Type[0].Type
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

// IsNumericReference checks if this test type uses numeric reference ranges
// Returns true if low_ref_range and high_ref_range are not both zero
func (t *TestType) IsNumericReference() bool {
	return t.LowRefRange != 0 || t.HighRefRange != 0
}

// IsStringReference checks if this test type uses string reference values
// Returns true if normal_ref_string is not empty and numeric ranges are both zero
func (t *TestType) IsStringReference() bool {
	return t.NormalRefString != "" && !t.IsNumericReference()
}

// GetReferenceRange returns the appropriate reference range based on type
func (t *TestType) GetReferenceRange() string {
	// Debug logging

	if t.IsNumericReference() {
		decimal := t.Decimal
		if decimal < 0 {
			decimal = 2
		}
		result := fmt.Sprintf("%.*f - %.*f", decimal, t.LowRefRange, decimal, t.HighRefRange)
		return result
	}
	return t.NormalRefString
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
	DeviceID      *int     `query:"device_id"`
}
