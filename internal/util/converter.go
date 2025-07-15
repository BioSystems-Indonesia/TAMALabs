package util

import (
	"errors"
	"fmt"
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
)

// normalizeUnit replaces escape sequences or alternative representations with standard units
func normalizeUnit(unit string) string {
	// Replace escape sequences with actual characters
	unit = strings.ReplaceAll(unit, `\Z03BC`, "µ") // Replace \Z03BC with µ
	unit = strings.ReplaceAll(unit, `\u00B5`, "µ") // Replace \u00B5 with µ
	unit = strings.ReplaceAll(unit, "ug", "µg")    // Replace ug with µg
	unit = strings.ReplaceAll(unit, "uL", "µL")    // Replace uL with µL
	unit = strings.ReplaceAll(unit, "uU", "µU")    // Replace uU with µU
	return unit
}

// parseCompoundUnit splits a compound unit into numerator and denominator
func parseCompoundUnit(unit string) (string, string) {
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
