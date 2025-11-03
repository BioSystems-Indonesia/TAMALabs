package util

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Conversion factors relative to base units
var (
	// Mass units (base unit: g)
	massConversions = map[string]float64{
		"kg":  1e3,   // 1 kg = 1000 g
		"hg":  1e2,   // 1 hg = 100 g
		"dag": 1e1,   // 1 dag = 10 g
		"g":   1,     // 1 g = 1 g
		"dg":  1e-1,  // 1 dg = 0.1 g
		"cg":  1e-2,  // 1 cg = 0.01 g
		"mg":  1e-3,  // 1 mg = 0.001 g
		"µg":  1e-6,  // 1 µg = 1e-6 g
		"ug":  1e-6,  // Alternative representation for µg
		"ng":  1e-9,  // 1 ng = 1e-9 g
		"pg":  1e-12, // 1 pg = 1e-12 g
		"fg":  1e-15, // 1 fg = 1e-15 g
	}

	// Volume units (base unit: L)
	volumeConversions = map[string]float64{
		"kL":  1e3,   // 1 kL = 1000 L
		"hL":  1e2,   // 1 hL = 100 L
		"daL": 1e1,   // 1 daL = 10 L
		"L":   1,     // 1 L = 1 L
		"dL":  1e-1,  // 1 dL = 0.1 L
		"cL":  1e-2,  // 1 cL = 0.01 L
		"mL":  1e-3,  // 1 mL = 0.001 L
		"µL":  1e-6,  // 1 µL = 1e-6 L
		"uL":  1e-6,  // Alternative representation for µL
		"nL":  1e-9,  // 1 nL = 1e-9 L
		"pL":  1e-12, // 1 pL = 1e-12 L
		"fL":  1e-15, // 1 fL = 1e-15 L
	}

	// Enzyme activity units (base unit: U/L)
	enzymeConversions = map[string]float64{
		"kU":  1e3,   // 1 kU = 1000 U
		"hU":  1e2,   // 1 hU = 100 U
		"daU": 1e1,   // 1 daU = 10 U
		"U":   1,     // 1 U = 1 U
		"dU":  1e-1,  // 1 dU = 0.1 U
		"cU":  1e-2,  // 1 cU = 0.01 U
		"mU":  1e-3,  // 1 mU = 0.001 U
		"µU":  1e-6,  // 1 µU = 1e-6 U
		"uU":  1e-6,  // Alternative representation for µU
		"nU":  1e-9,  // 1 nU = 1e-9 U
		"pU":  1e-12, // 1 pU = 1e-12 U
		"fU":  1e-15, // 1 fU = 1e-15 U
	}

	// Count units (base unit: /µL)
	countConversions = map[string]float64{
		"G/µL":    1e9,  // 1 G/µL = 1e9 /µL
		"M/µL":    1e6,  // 1 M/µL = 1e6 /µL
		"K/µL":    1e3,  // 1 K/µL = 1e3 /µL
		"10^3/µL": 1e3,  // 1 10^3/µL = 1e3 /µL (same as K/µL)
		"10^6/µL": 1e6,  // 1 10^6/µL = 1e6 /µL (same as M/µL)
		"10^9/µL": 1e9,  // 1 10^9/µL = 1e9 /µL (same as G/µL)
		"/µL":     1,    // 1 /µL = 1 /µL
		"/mL":     1e-3, // 1 /mL = 1e-3 /µL
		"/L":      1e-6, // 1 /L = 1e-6 /µL
	}
)

// normalizeUnit replaces escape sequences or alternative representations with standard units
func normalizeUnit(unit string) string {
	// Replace escape sequences with actual characters
	unit = strings.ReplaceAll(unit, `\Z03BC`, "µ") // Replace \Z03BC with µ
	unit = strings.ReplaceAll(unit, `\u00B5`, "µ") // Replace \u00B5 with µ
	unit = strings.ReplaceAll(unit, "ug", "µg")    // Replace ug with µg
	unit = strings.ReplaceAll(unit, "uL", "µL")    // Replace uL with µL
	unit = strings.ReplaceAll(unit, "uU", "µU")    // Replace uU with µU

	// Handle alternative representations for scientific notation
	unit = strings.ReplaceAll(unit, "10^3/uL", "10^3/µL") // Replace 10^3/uL with 10^3/µL
	unit = strings.ReplaceAll(unit, "10^6/uL", "10^6/µL") // Replace 10^6/uL with 10^6/µL
	unit = strings.ReplaceAll(unit, "10^9/uL", "10^9/µL") // Replace 10^9/uL with 10^9/µL

	return unit
}

// parseCompoundUnit splits a compound unit into numerator and denominator
func parseCompoundUnit(unit string) (string, string) {
	// Handle count units specially - these should be treated as single units
	if isCountUnit(unit) {
		return unit, "" // Treat count units as single units
	}

	parts := strings.Split(unit, "/")
	if len(parts) == 1 {
		return parts[0], "" // No denominator
	}
	return parts[0], parts[1]
}

// convertSimpleUnit converts a simple unit (no denominator)
func convertSimpleUnit(value float64, fromUnit, toUnit string) (float64, error) {
	fromUnit = normalizeUnit(fromUnit)
	toUnit = normalizeUnit(toUnit)

	// If units are the same, no conversion needed
	if fromUnit == toUnit {
		return value, nil
	}

	// Determine the conversion map based on the unit type
	var conversionMap map[string]float64
	switch {
	case isMassUnit(fromUnit) || isMassUnit(toUnit):
		conversionMap = massConversions
	case isVolumeUnit(fromUnit) || isVolumeUnit(toUnit):
		conversionMap = volumeConversions
	case isEnzymeUnit(fromUnit) || isEnzymeUnit(toUnit):
		conversionMap = enzymeConversions
	case isCountUnit(fromUnit) || isCountUnit(toUnit):
		conversionMap = countConversions
	default:
		return 0, errors.New("unsupported unit type")
	}

	// Check if units exist in the conversion map
	if _, ok := conversionMap[fromUnit]; !ok {
		return 0, fmt.Errorf("fromUnit '%s' not found in conversion map", fromUnit)
	}
	if _, ok := conversionMap[toUnit]; !ok {
		return 0, fmt.Errorf("toUnit '%s' not found in conversion map", toUnit)
	}

	// Convert fromUnit to base unit, then to toUnit
	baseValue := value * conversionMap[fromUnit] // Convert to base unit
	result := baseValue / conversionMap[toUnit]  // Convert from base unit to toUnit
	return result, nil
}

// ConvertCompoundUnit converts a compound unit (e.g., mg/dL, U/L)
func ConvertCompoundUnit(value float64, fromUnit, toUnit string) (float64, error) {
	// If units are the same, no conversion needed
	if fromUnit == toUnit {
		return value, nil
	}

	// Parse the compound units
	fromNum, fromDen := parseCompoundUnit(fromUnit)
	toNum, toDen := parseCompoundUnit(toUnit)

	// Convert the numerator
	numValue, err := convertSimpleUnit(value, fromNum, toNum)
	if err != nil {
		return 0, fmt.Errorf("numerator conversion failed: %v", err)
	}

	// Convert the denominator (if it exists)
	if fromDen != "" && toDen != "" {
		denValue, err := convertSimpleUnit(1, fromDen, toDen)
		if err != nil {
			return 0, fmt.Errorf("denominator conversion failed: %v", err)
		}
		numValue /= denValue
	} else if fromDen != "" || toDen != "" {
		return 0, errors.New("mismatched compound units")
	}

	return numValue, nil
}

// Helper functions to determine unit type
func isMassUnit(unit string) bool {
	_, ok := massConversions[unit]
	return ok
}

func isVolumeUnit(unit string) bool {
	_, ok := volumeConversions[unit]
	return ok
}

func isEnzymeUnit(unit string) bool {
	_, ok := enzymeConversions[unit]
	return ok
}

func isCountUnit(unit string) bool {
	_, ok := countConversions[unit]
	return ok
}

// Helper function to compare floating-point numbers with a tolerance
func almostEqual(a, b, tolerance float64) bool {
	return abs(a-b) <= tolerance
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// ConvertReferenceRange converts a reference range by multiplying both values by a factor
func ConvertReferenceRange(refRange string, factor float64) string {
	return ConvertReferenceRangeWithDecimal(refRange, factor, 2)
}

// ConvertReferenceRangeWithDecimal converts a reference range by multiplying both values by a factor
// and formats with specified decimal places
func ConvertReferenceRangeWithDecimal(refRange string, factor float64, decimal int) string {
	if refRange == "" {
		return ""
	}

	// Use regex to match format like "4.0 - 10.0" or "150 - 450"
	parts := strings.Split(refRange, "-")
	if len(parts) != 2 {
		return refRange // Return original if format doesn't match
	}

	lowStr := strings.TrimSpace(parts[0])
	highStr := strings.TrimSpace(parts[1])

	lowValue, err := strconv.ParseFloat(lowStr, 64)
	if err != nil {
		return refRange // Return original if parsing fails
	}

	highValue, err := strconv.ParseFloat(highStr, 64)
	if err != nil {
		return refRange // Return original if parsing fails
	}

	// Convert values
	convertedLow := lowValue * factor
	convertedHigh := highValue * factor

	// Ensure decimal is not negative
	if decimal < 0 {
		decimal = 0
	}

	// Format with appropriate decimal places
	if decimal == 0 || (convertedLow == float64(int(convertedLow)) && convertedHigh == float64(int(convertedHigh))) {
		// If decimal is 0 or both are whole numbers, format without decimal places
		return fmt.Sprintf("%.0f - %.0f", convertedLow, convertedHigh)
	} else {
		// Otherwise, format with specified decimal places
		return fmt.Sprintf("%.*f - %.*f", decimal, convertedLow, decimal, convertedHigh)
	}
}
