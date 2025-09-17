package util

import (
	"errors"
	"math"
	"time"
)

// PatientSex represents patient gender for eGFR calculation
type PatientSex string

const (
	PatientSexMale   PatientSex = "M"
	PatientSexFemale PatientSex = "F"
)

// EGFRResult represents the calculated eGFR result
type EGFRResult struct {
	Value    float64 `json:"value"`
	Formula  string  `json:"formula"`
	Unit     string  `json:"unit"`
	Category string  `json:"category"`
}

var ErrUnsupportedCreatinineUnit = errors.New("unsupported creatinine unit")

// CalculateAge calculates age in years from birthdate
func CalculateAge(birthdate time.Time) float64 {
	now := time.Now()
	years := now.Sub(birthdate).Hours() / 24 / 365.25
	return years
}

// CalculateEGFRCKDEPI calculates eGFR using CKD-EPI formula (2021 equation)
// This is the most current and widely accepted formula
func CalculateEGFRCKDEPI(creatinine float64, age float64, sex PatientSex) EGFRResult {
	var egfr float64
	var sexFactor float64

	// Set sex-specific factors based on CKD-EPI 2021 equation
	if sex == PatientSexFemale {
		sexFactor = 1.012
	} else {
		sexFactor = 1.0
	}

	// Convert mg/dL to mg/dL if needed (assuming input is in mg/dL)
	// If input is in µmol/L, divide by 88.4 to convert to mg/dL

	// Calculate eGFR using CKD-EPI 2021 formula
	var alpha, kappa float64
	if sex == PatientSexFemale {
		alpha = -0.241
		kappa = 0.7
	} else {
		alpha = -0.302
		kappa = 0.9
	}

	minRatio := math.Min(creatinine/kappa, 1.0)
	maxRatio := math.Max(creatinine/kappa, 1.0)

	egfr = 142 * math.Pow(minRatio, alpha) * math.Pow(maxRatio, -1.200) * math.Pow(0.9938, age) * sexFactor

	category := categorizeEGFR(egfr)

	return EGFRResult{
		Value:    math.Round(egfr*100) / 100, // Round to 2 decimal places
		Formula:  "CKD-EPI 2021",
		Unit:     "mL/min/1.73m²",
		Category: category,
	}
}

// CalculateEGFRMDRD calculates eGFR using MDRD formula (legacy, but still used)
func CalculateEGFRMDRD(creatinine float64, age float64, sex PatientSex) EGFRResult {
	var egfr float64
	var sexFactor float64

	if sex == PatientSexFemale {
		sexFactor = 0.742
	} else {
		sexFactor = 1.0
	}

	// MDRD formula: 175 × (Scr)^-1.154 × (Age)^-0.203 × (0.742 if female)
	egfr = 175 * math.Pow(creatinine, -1.154) * math.Pow(age, -0.203) * sexFactor

	category := categorizeEGFR(egfr)

	return EGFRResult{
		Value:    math.Round(egfr*100) / 100,
		Formula:  "MDRD",
		Unit:     "mL/min/1.73m²",
		Category: category,
	}
}

// categorizeEGFR categorizes eGFR value according to CKD stages
func categorizeEGFR(egfr float64) string {
	switch {
	case egfr >= 90:
		return "Normal or High (≥90)"
	case egfr >= 60:
		return "Mildly Decreased (60-89)"
	case egfr >= 45:
		return "Mild to Moderately Decreased (45-59)"
	case egfr >= 30:
		return "Moderately to Severely Decreased (30-44)"
	case egfr >= 15:
		return "Severely Decreased (15-29)"
	default:
		return "Kidney Failure (<15)"
	}
}

// ConvertCreatinineUnit converts creatinine from one unit to another
func ConvertCreatinineUnit(value float64, fromUnit, toUnit string) (float64, error) {
	// Normalize units
	fromUnit = normalizeUnit(fromUnit)
	toUnit = normalizeUnit(toUnit)

	if fromUnit == toUnit {
		return value, nil
	}

	// Convert to mg/dL as base unit
	var valueInMgDL float64
	switch fromUnit {
	case "µmol/L", "umol/L":
		valueInMgDL = value / 88.4 // Convert µmol/L to mg/dL
	case "mg/dL":
		valueInMgDL = value
	default:
		return 0, ErrUnsupportedCreatinineUnit
	}

	// Convert from mg/dL to target unit
	switch toUnit {
	case "µmol/L", "umol/L":
		return valueInMgDL * 88.4, nil // Convert mg/dL to µmol/L
	case "mg/dL":
		return valueInMgDL, nil
	default:
		return 0, ErrUnsupportedCreatinineUnit
	}
}
