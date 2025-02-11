package util

import (
	"testing"
)

// TestConvertSimpleUnit tests conversions for simple units (no denominator)
func TestConvertSimpleUnit(t *testing.T) {
	tests := []struct {
		value    float64
		fromUnit string
		toUnit   string
		expected float64
	}{
		// Mass conversions
		{1, "kg", "g", 1000},  // 1 kg = 1000 g
		{1, "g", "mg", 1000},  // 1 g = 1000 mg
		{1, "mg", "µg", 1000}, // 1 mg = 1000 µg
		{1, "µg", "ng", 1000}, // 1 µg = 1000 ng
		{1, "ng", "pg", 1000}, // 1 ng = 1000 pg
		{1, "pg", "fg", 1000}, // 1 pg = 1000 fg

		// Volume conversions
		{1, "kL", "L", 1000},  // 1 kL = 1000 L
		{1, "L", "mL", 1000},  // 1 L = 1000 mL
		{1, "mL", "µL", 1000}, // 1 mL = 1000 µL
		{1, "µL", "nL", 1000}, // 1 µL = 1000 nL
		{1, "nL", "pL", 1000}, // 1 nL = 1000 pL
		{1, "pL", "fL", 1000}, // 1 pL = 1000 fL
	}

	for _, test := range tests {
		result, err := convertSimpleUnit(test.value, test.fromUnit, test.toUnit)
		if err != nil {
			t.Errorf("convertCompoundUnit(%v, %s, %s) failed: %v", test.value, test.fromUnit, test.toUnit, err)
		} else if !almostEqual(result, test.expected, 1e-9) {
			t.Errorf("convertCompoundUnit(%v, %s, %s) = %v; expected %v", test.value, test.fromUnit, test.toUnit, result, test.expected)
		}
	}
}

// TestConvertCompoundUnit tests conversions for compound units (e.g., mg/dL, U/L)
func TestConvertCompoundUnit(t *testing.T) {
	tests := []struct {
		value    float64
		fromUnit string
		toUnit   string
		expected float64
	}{
		// Mass/Volume conversions
		{100, "mg/dL", "g/L", 1},        // 100 mg/dL = 1 g/L
		{1, "g/L", "mg/dL", 100},        // 1 g/L = 100 mg/dL
		{10, "µg/mL", "ng/dL", 1000000}, // 10 µg/mL = 1,000,000 ng/dL

		// Enzyme activity conversions
		{
			50, "U/L", "mU/mL", 50, // 50 U/L = 50 mU/mL
		},
		{
			5, "kU/L", "U/mL", 5, // 5 kU/L = 5 U/mL
		},
		{1, "kU/L", "U/L", 1000},  // 1 kU/L = 1000 U/L
		{1, "U/L", "mU/L", 1000},  // 1 U/L = 1000 mU/L
		{1, "mU/L", "µU/L", 1000}, // 1 mU/L = 1000 µU/L
		{1, "µU/L", "nU/L", 1000}, // 1 µU/L = 1000 nU/L
		{1, "nU/L", "pU/L", 1000}, // 1 nU/L = 1000 pU/L
		{1, "pU/L", "fU/L", 1000}, // 1 pU/L = 1000 fU/L

		// Volume conversions
		{2, "dL", "mL", 200}, // 2 dL = 200 mL
		{1, "L", "mL", 1000}, // 1 L = 1000 mL
	}

	for _, test := range tests {
		result, err := ConvertCompoundUnit(test.value, test.fromUnit, test.toUnit)
		if err != nil {
			t.Errorf("convertCompoundUnit(%v, %s, %s) failed: %v", test.value, test.fromUnit, test.toUnit, err)
		} else if !almostEqual(result, test.expected, 1e-9) {
			t.Errorf("convertCompoundUnit(%v, %s, %s) = %v; expected %v", test.value, test.fromUnit, test.toUnit, result, test.expected)
		}
	}
}

// TestNormalizeUnit tests the normalization of units
func TestNormalizeUnit(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`\Z03BCg`, "µg"},  // \Z03BCg → µg
		{`\u00B5g`, "µg"},  // \u00B5g → µg
		{"ug", "µg"},       // ug → µg
		{"uL", "µL"},       // uL → µL
		{"uU", "µU"},       // uU → µU
		{"mg/dL", "mg/dL"}, // No change
	}

	for _, test := range tests {
		result := normalizeUnit(test.input)
		if result != test.expected {
			t.Errorf("normalizeUnit(%s) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

// TestErrorHandling tests error cases
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		value    float64
		fromUnit string
		toUnit   string
	}{
		{1, "invalid", "g"}, // Invalid fromUnit
		{1, "g", "invalid"}, // Invalid toUnit
		{1, "mg/dL", "g"},   // Mismatched compound units
		{1, "g", "mg/dL"},   // Mismatched compound units
	}

	for _, test := range tests {
		_, err := ConvertCompoundUnit(test.value, test.fromUnit, test.toUnit)
		if err == nil {
			t.Errorf("convertCompoundUnit(%v, %s, %s) should have failed", test.value, test.fromUnit, test.toUnit)
		}
	}
}
