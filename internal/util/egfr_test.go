package util

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	// Test case: Person born 30 years ago
	birthdate := time.Now().AddDate(-30, 0, 0)
	age := CalculateAge(birthdate)

	if age < 29.9 || age > 30.1 {
		t.Errorf("Expected age around 30, got %f", age)
	}
}

func TestCalculateEGFRCKDEPI(t *testing.T) {
	tests := []struct {
		name       string
		creatinine float64
		age        float64
		sex        PatientSex
		expected   float64
		tolerance  float64
	}{
		{
			name:       "Normal young male",
			creatinine: 1.0,
			age:        25,
			sex:        PatientSexMale,
			expected:   100,
			tolerance:  10,
		},
		{
			name:       "Normal young female",
			creatinine: 0.8,
			age:        25,
			sex:        PatientSexFemale,
			expected:   110,
			tolerance:  15,
		},
		{
			name:       "Elderly male with elevated creatinine",
			creatinine: 2.0,
			age:        70,
			sex:        PatientSexMale,
			expected:   35,
			tolerance:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateEGFRCKDEPI(tt.creatinine, tt.age, tt.sex)

			if result.Value < tt.expected-tt.tolerance || result.Value > tt.expected+tt.tolerance {
				t.Errorf("Expected eGFR around %f±%f, got %f", tt.expected, tt.tolerance, result.Value)
			}

			if result.Formula != "CKD-EPI 2021" {
				t.Errorf("Expected formula 'CKD-EPI 2021', got %s", result.Formula)
			}

			if result.Unit != "mL/min/1.73m²" {
				t.Errorf("Expected unit 'mL/min/1.73m²', got %s", result.Unit)
			}
		})
	}
}

func TestCategorizeEGFR(t *testing.T) {
	tests := []struct {
		egfr     float64
		expected string
	}{
		{100, "Normal or High (≥90)"},
		{75, "Mildly Decreased (60-89)"},
		{50, "Mild to Moderately Decreased (45-59)"},
		{35, "Moderately to Severely Decreased (30-44)"},
		{20, "Severely Decreased (15-29)"},
		{10, "Kidney Failure (<15)"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := categorizeEGFR(tt.egfr)
			if result != tt.expected {
				t.Errorf("Expected category '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestConvertCreatinineUnit(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		fromUnit  string
		toUnit    string
		expected  float64
		tolerance float64
	}{
		{
			name:      "µmol/L to mg/dL",
			value:     88.4,
			fromUnit:  "µmol/L",
			toUnit:    "mg/dL",
			expected:  1.0,
			tolerance: 0.01,
		},
		{
			name:      "mg/dL to µmol/L",
			value:     1.0,
			fromUnit:  "mg/dL",
			toUnit:    "µmol/L",
			expected:  88.4,
			tolerance: 0.1,
		},
		{
			name:      "Same unit",
			value:     1.0,
			fromUnit:  "mg/dL",
			toUnit:    "mg/dL",
			expected:  1.0,
			tolerance: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertCreatinineUnit(tt.value, tt.fromUnit, tt.toUnit)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result < tt.expected-tt.tolerance || result > tt.expected+tt.tolerance {
				t.Errorf("Expected %f±%f, got %f", tt.expected, tt.tolerance, result)
			}
		})
	}
}
