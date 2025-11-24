package entity

import (
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type TestType struct {
	ID               int     `json:"id" gorm:"primaryKey"`
	Name             string  `json:"name"`
	Code             string  `json:"code" gorm:"unique"`
	AliasCode        string  `json:"alias_code" gorm:"index:test_type_alias_code"`
	LoincCode        string  `json:"loinc_code" gorm:"index:test_type_loinc_code"`
	Unit             string  `json:"unit"`
	LowRefRange      float64 `json:"low_ref_range"`
	HighRefRange     float64 `json:"high_ref_range"`
	NormalRefString  string  `json:"normal_ref_string" gorm:"column:normal_ref_string"` // New field for string reference values
	Decimal          int     `json:"decimal"`
	Category         string  `json:"category"`
	SubCategory      string  `json:"sub_category"`
	Description      string  `json:"description"`
	IsCalculatedTest bool    `json:"is_calculated_test" gorm:"column:is_calculated_test;default:false"`
	DeviceID         *int    `json:"device_id" gorm:"index:test_type_device_id"` // Deprecated: Use Devices instead

	Type              []TestTypeSpecimenType   `json:"types" gorm:"-"`
	Device            *Device                  `json:"device,omitempty" gorm:"foreignKey:DeviceID"` // Deprecated: Use Devices instead
	Devices           []Device                 `json:"devices,omitempty" gorm:"many2many:test_type_devices;"`
	SpecificRefRanges []SpecificReferenceRange `json:"specific_ref_ranges,omitempty" gorm:"-"` // Not stored directly, converted to/from JSON
	AlternativeCodes  []string                 `json:"alternative_codes,omitempty" gorm:"-"`   // Alternative codes for the same test from different devices
	// TypeDB is a specimen type separated by comma
	TypeDB              string `json:"-" gorm:"column:type"`
	SpecificRefRangesDB string `json:"-" gorm:"column:specific_ref_ranges"` // JSON string in database
	AlternativeCodesDB  string `json:"-" gorm:"column:alternative_codes"`   // JSON array in database
}

// MarshalJSON custom JSON marshaling to include device_ids
func (t *TestType) MarshalJSON() ([]byte, error) {
	deviceIDs := make([]int, len(t.Devices))
	for i, device := range t.Devices {
		deviceIDs[i] = device.ID
	}

	// Create anonymous struct to avoid recursion
	return json.Marshal(&struct {
		ID                int                      `json:"id"`
		Name              string                   `json:"name"`
		Code              string                   `json:"code"`
		AliasCode         string                   `json:"alias_code"`
		LoincCode         string                   `json:"loinc_code"`
		Unit              string                   `json:"unit"`
		LowRefRange       float64                  `json:"low_ref_range"`
		HighRefRange      float64                  `json:"high_ref_range"`
		NormalRefString   string                   `json:"normal_ref_string"`
		Decimal           int                      `json:"decimal"`
		Category          string                   `json:"category"`
		SubCategory       string                   `json:"sub_category"`
		Description       string                   `json:"description"`
		IsCalculatedTest  bool                     `json:"is_calculated_test"`
		DeviceID          *int                     `json:"device_id,omitempty"`
		Type              []TestTypeSpecimenType   `json:"types"`
		Device            *Device                  `json:"device,omitempty"`
		Devices           []Device                 `json:"devices,omitempty"`
		DeviceIDs         []int                    `json:"device_ids,omitempty"`
		SpecificRefRanges []SpecificReferenceRange `json:"specific_ref_ranges,omitempty"`
	}{
		ID:                t.ID,
		Name:              t.Name,
		Code:              t.Code,
		AliasCode:         t.AliasCode,
		LoincCode:         t.LoincCode,
		Unit:              t.Unit,
		LowRefRange:       t.LowRefRange,
		HighRefRange:      t.HighRefRange,
		NormalRefString:   t.NormalRefString,
		Decimal:           t.Decimal,
		Category:          t.Category,
		SubCategory:       t.SubCategory,
		Description:       t.Description,
		IsCalculatedTest:  t.IsCalculatedTest,
		DeviceID:          t.DeviceID,
		Type:              t.Type,
		Device:            t.Device,
		Devices:           t.Devices,
		DeviceIDs:         deviceIDs,
		SpecificRefRanges: t.SpecificRefRanges,
	})
}

func (t *TestType) GetFirstType() string {
	if len(t.Type) == 0 {
		return "SER"
	}

	return t.Type[0].Type
}

// HasCode checks if the given code matches the main Code or any AlternativeCodes (case-insensitive)
func (t *TestType) HasCode(searchCode string) bool {
	searchCode = strings.ToLower(strings.TrimSpace(searchCode))

	// Check main code
	if strings.ToLower(strings.TrimSpace(t.Code)) == searchCode {
		return true
	}

	// Check alternative codes
	for _, altCode := range t.AlternativeCodes {
		if strings.ToLower(strings.TrimSpace(altCode)) == searchCode {
			return true
		}
	}

	return false
}

// TestTypeSpecimenType is a specimen type struct, we need it to conform with react-admin standard
// for ordered input
type TestTypeSpecimenType struct {
	Type string `json:"type"`
}

// SpecificReferenceRange represents a reference range for specific criteria (age/gender)
type SpecificReferenceRange struct {
	Gender          *string  `json:"gender,omitempty"`            // "M", "F", or null for all
	AgeMin          *float64 `json:"age_min,omitempty"`           // minimum age in years
	AgeMax          *float64 `json:"age_max,omitempty"`           // maximum age in years
	LowRefRange     *float64 `json:"low_ref_range,omitempty"`     // numeric low value
	HighRefRange    *float64 `json:"high_ref_range,omitempty"`    // numeric high value
	NormalRefString *string  `json:"normal_ref_string,omitempty"` // or string value
}

// MatchesCriteria checks if this range matches patient criteria
func (s *SpecificReferenceRange) MatchesCriteria(age *float64, gender *string) bool {
	// Check gender
	if s.Gender != nil && gender != nil && *s.Gender != *gender {
		return false
	}

	// Check age range
	if age != nil {
		if s.AgeMin != nil && *age < *s.AgeMin {
			return false
		}
		if s.AgeMax != nil && *age > *s.AgeMax {
			return false
		}
	}

	return true
}

// GetReferenceRange formats the reference range as string
func (s *SpecificReferenceRange) GetReferenceRange(decimal int) string {
	if s.NormalRefString != nil && *s.NormalRefString != "" {
		return *s.NormalRefString
	}

	if s.LowRefRange != nil && s.HighRefRange != nil {
		if decimal < 0 {
			decimal = 0
		}
		lowStr := fmt.Sprintf("%.*f", decimal, *s.LowRefRange)
		highStr := fmt.Sprintf("%.*f", decimal, *s.HighRefRange)
		return fmt.Sprintf("%s - %s", lowStr, highStr)
	}

	return ""
}

// GetDescription returns a human-readable description
func (s *SpecificReferenceRange) GetDescription() string {
	var parts []string

	if s.Gender != nil {
		switch *s.Gender {
		case "M":
			parts = append(parts, "Male")
		case "F":
			parts = append(parts, "Female")
		}
	} else {
		parts = append(parts, "All Genders")
	}

	if s.AgeMin != nil && s.AgeMax != nil {
		parts = append(parts, fmt.Sprintf("%.0f-%.0f years", *s.AgeMin, *s.AgeMax))
	} else if s.AgeMin != nil {
		parts = append(parts, fmt.Sprintf("≥%.0f years", *s.AgeMin))
	} else if s.AgeMax != nil {
		parts = append(parts, fmt.Sprintf("≤%.0f years", *s.AgeMax))
	} else {
		parts = append(parts, "All Ages")
	}

	return strings.Join(parts, ", ")
}

func (t *TestType) BeforeCreate(tx *gorm.DB) error {
	types := make([]string, len(t.Type))
	for i, specimenType := range t.Type {
		types[i] = specimenType.Type
	}
	t.TypeDB = strings.Join(types, ",")

	// Serialize SpecificRefRanges to JSON
	if len(t.SpecificRefRanges) > 0 {
		jsonData, err := json.Marshal(t.SpecificRefRanges)
		if err != nil {
			return err
		}
		t.SpecificRefRangesDB = string(jsonData)
	}

	// Serialize AlternativeCodes to JSON
	if len(t.AlternativeCodes) > 0 {
		jsonData, err := json.Marshal(t.AlternativeCodes)
		if err != nil {
			return err
		}
		t.AlternativeCodesDB = string(jsonData)
	}

	return nil
}

func (t *TestType) BeforeUpdate(tx *gorm.DB) error {
	types := make([]string, len(t.Type))
	for i, specimenType := range t.Type {
		types[i] = specimenType.Type
	}
	t.TypeDB = strings.Join(types, ",")

	// Serialize SpecificRefRanges to JSON
	if len(t.SpecificRefRanges) > 0 {
		jsonData, err := json.Marshal(t.SpecificRefRanges)
		if err != nil {
			return err
		}
		t.SpecificRefRangesDB = string(jsonData)
	} else {
		t.SpecificRefRangesDB = ""
	}

	// Serialize AlternativeCodes to JSON
	if len(t.AlternativeCodes) > 0 {
		jsonData, err := json.Marshal(t.AlternativeCodes)
		if err != nil {
			return err
		}
		t.AlternativeCodesDB = string(jsonData)
	} else {
		t.AlternativeCodesDB = ""
	}

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

	// Deserialize SpecificRefRanges from JSON
	if t.SpecificRefRangesDB != "" {
		var ranges []SpecificReferenceRange
		if err := json.Unmarshal([]byte(t.SpecificRefRangesDB), &ranges); err == nil {
			t.SpecificRefRanges = ranges
		}
	}

	// Deserialize AlternativeCodes from JSON
	if t.AlternativeCodesDB != "" {
		var codes []string
		if err := json.Unmarshal([]byte(t.AlternativeCodesDB), &codes); err == nil {
			t.AlternativeCodes = codes
		}
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
	if t.IsNumericReference() {
		decimal := t.Decimal
		if decimal < 0 {
			decimal = 0 // Default to 0 decimal places (whole numbers) if not set
		}

		// Format with specified decimal places
		// This preserves trailing zeros (e.g., "1.90" instead of "1.9")
		lowStr := fmt.Sprintf("%.*f", decimal, t.LowRefRange)
		highStr := fmt.Sprintf("%.*f", decimal, t.HighRefRange)

		return fmt.Sprintf("%s - %s", lowStr, highStr)
	}
	return t.NormalRefString
}

// GetSpecificReferenceRangeForPatient returns the matching specific reference range struct
// for a patient based on age and gender. Returns nil if no specific range matches.
func (t *TestType) GetSpecificReferenceRangeForPatient(age *float64, gender *string) *SpecificReferenceRange {
	// Try to find a matching specific reference range
	for i := range t.SpecificRefRanges {
		if t.SpecificRefRanges[i].MatchesCriteria(age, gender) {
			return &t.SpecificRefRanges[i]
		}
	}
	return nil
}

// GetReferenceRangeForPatient returns the appropriate reference range for a patient
// based on age and gender. If a specific range exists, use it; otherwise, use the default.
func (t *TestType) GetReferenceRangeForPatient(age *float64, gender *string) string {
	// Try to find a matching specific reference range
	specificRange := t.GetSpecificReferenceRangeForPatient(age, gender)
	if specificRange != nil {
		rangeStr := specificRange.GetReferenceRange(t.Decimal)
		if rangeStr != "" {
			return rangeStr
		}
	}

	// Fallback to default reference range
	defaultRange := t.GetReferenceRange()
	return defaultRange
}

// GetAllReferenceRanges returns all reference ranges including the default one
func (t *TestType) GetAllReferenceRanges() []map[string]interface{} {
	ranges := make([]map[string]interface{}, 0)

	// Add default range
	defaultRange := map[string]interface{}{
		"description": "Default (All)",
		"gender":      nil,
		"age_min":     nil,
		"age_max":     nil,
		"range":       t.GetReferenceRange(),
	}
	ranges = append(ranges, defaultRange)

	// Add specific ranges
	for i, refRange := range t.SpecificRefRanges {
		ranges = append(ranges, map[string]interface{}{
			"index":       i,
			"description": refRange.GetDescription(),
			"gender":      refRange.Gender,
			"age_min":     refRange.AgeMin,
			"age_max":     refRange.AgeMax,
			"range":       refRange.GetReferenceRange(t.Decimal),
		})
	}

	return ranges
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
