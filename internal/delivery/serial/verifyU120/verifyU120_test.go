package verifyu120

import (
	"context"
	"testing"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func TestHandler_parseResultValue(t *testing.T) {
	h := &Handler{}

	tests := []struct {
		name          string
		input         string
		expectedValue float64
		expectedStr   string
		expectedUnit  string
	}{
		{
			name:          "Negative result with neg",
			input:         "-              neg",
			expectedValue: 0,
			expectedStr:   "neg",
			expectedUnit:  "",
		},
		{
			name:          "Simple negative",
			input:         "-",
			expectedValue: 0,
			expectedStr:   "neg",
			expectedUnit:  "",
		},
		{
			name:          "Qualitative positive with numeric value and unit",
			input:         "1+       0.3    g/L",
			expectedValue: 0.3,
			expectedStr:   "1+",
			expectedUnit:  "g/L",
		},
		{
			name:          "Qualitative strong positive with numeric value and unit",
			input:         "3+       200 Ery/uL",
			expectedValue: 200,
			expectedStr:   "3+",
			expectedUnit:  "Ery/uL",
		},
		{
			name:          "Numeric value with unit",
			input:         "3.5 umol/L",
			expectedValue: 3.5,
			expectedStr:   "3.5",
			expectedUnit:  "umol/L",
		},
		{
			name:          "Simple numeric value",
			input:         "6.0",
			expectedValue: 6.0,
			expectedStr:   "6.0",
			expectedUnit:  "",
		},
		{
			name:          "Specific gravity",
			input:         "1.025",
			expectedValue: 1.025,
			expectedStr:   "1.025",
			expectedUnit:  "",
		},
		{
			name:          "Ketone with qualitative and numeric",
			input:         "3+       8.0 mmol/L",
			expectedValue: 8.0,
			expectedStr:   "3+",
			expectedUnit:  "mmol/L",
		},
		{
			name:          "Just qualitative result",
			input:         "1+",
			expectedValue: 0,
			expectedStr:   "1+",
			expectedUnit:  "",
		},
		{
			name:          "Glucose negative",
			input:         "-              neg",
			expectedValue: 0,
			expectedStr:   "neg",
			expectedUnit:  "",
		},
		{
			name:          "URO with dash then numeric",
			input:         "-       3.5 umol/L",
			expectedValue: 3.5,
			expectedStr:   "3.5",
			expectedUnit:  "umol/L",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, valueStr, unit := h.parseResultValue(tt.input)

			if value != tt.expectedValue {
				t.Errorf("parseResultValue() value = %v, want %v", value, tt.expectedValue)
			}
			if valueStr != tt.expectedStr {
				t.Errorf("parseResultValue() valueStr = %v, want %v", valueStr, tt.expectedStr)
			}
			if unit != tt.expectedUnit {
				t.Errorf("parseResultValue() unit = %v, want %v", unit, tt.expectedUnit)
			}
		})
	}
}

func TestRegexPatterns(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		regex    string
		expected bool
	}{
		{
			name:     "Normal parameter",
			input:    "LEU       -              neg",
			regex:    `^\s*\*?([A-Z]{2,4}|pH)\s+(.+)$`,
			expected: true,
		},
		{
			name:     "Abnormal parameter with asterisk",
			input:    "*PRO      1+       0.3    g/L",
			regex:    `^\s*\*?([A-Z]{2,4}|pH)\s+(.+)$`,
			expected: true,
		},
		{
			name:     "pH parameter",
			input:    " pH           6.0",
			regex:    `^\s*\*?([A-Z]{2,4}|pH)\s+(.+)$`,
			expected: true,
		},
		{
			name:     "ID line",
			input:    " ID:1502250001",
			regex:    `ID\s*:\s*(\d+)`,
			expected: true,
		},
		{
			name:     "Date line",
			input:    " Date:09-28-2025 08:51 pm",
			regex:    `Date\s*:\s*(.+)`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var matched bool
			switch tt.regex {
			case `^\s*\*?([A-Z]{2,4}|pH)\s+(.+)$`:
				matched = reResult.MatchString(tt.input)
			case `ID\s*:\s*(\d+)`:
				matched = reID.MatchString(tt.input)
			case `Date\s*:\s*(.+)`:
				matched = reDate.MatchString(tt.input)
			}

			if matched != tt.expected {
				t.Errorf("Regex match for %s = %v, want %v", tt.input, matched, tt.expected)
			}
		})
	}
}

func TestHandler_parseUrineResult_FullSample(t *testing.T) {
	// Mock analyzer use case to capture results
	var capturedResults []entity.VerifyResult
	mockUseCase := &MockAnalyzerUseCase{
		results: &capturedResults,
	}

	h := NewHandler(mockUseCase)

	// Sample data based on the provided format
	sampleData := ` ID:1502250001
 Date:09-28-2025 08:51 pm
 Operator: 01
 No. 000014
 LEU       -              neg
 NIT       -              neg
 URO       -       3.5 umol/L
*PRO      1+       0.3    g/L
 pH           6.0
*BLO      3+       200 Ery/uL
 SG         1.025
*KET      3+       8.0 mmol/L
 BIL       -              neg
 GLU       -              neg`

	h.buffer = sampleData
	h.parseUrineResult()

	// Verify results
	expectedResults := map[string]struct {
		value    float64
		valueStr string
		unit     string
	}{
		"LEU": {0, "neg", ""},
		"NIT": {0, "neg", ""},
		"URO": {3.5, "3.5", "umol/L"},
		"PRO": {0.3, "1+", "g/L"},
		"pH":  {6.0, "6.0", ""},
		"BLO": {200, "3+", "Ery/uL"},
		"SG":  {1.025, "1.025", ""},
		"KET": {8.0, "3+", "mmol/L"},
		"BIL": {0, "neg", ""},
		"GLU": {0, "neg", ""},
	}

	if len(capturedResults) != len(expectedResults) {
		t.Errorf("Expected %d results, got %d", len(expectedResults), len(capturedResults))
	}

	for _, result := range capturedResults {
		expected, exists := expectedResults[result.TestName]
		if !exists {
			t.Errorf("Unexpected test result: %s", result.TestName)
			continue
		}

		if result.Value != expected.value {
			t.Errorf("Test %s: expected value %v, got %v", result.TestName, expected.value, result.Value)
		}

		if result.ValueStr != expected.valueStr {
			t.Errorf("Test %s: expected valueStr %v, got %v", result.TestName, expected.valueStr, result.ValueStr)
		}

		if result.Unit != expected.unit {
			t.Errorf("Test %s: expected unit %v, got %v", result.TestName, expected.unit, result.Unit)
		}

		if result.PatientID != "1502250001" {
			t.Errorf("Test %s: expected patientID 1502250001, got %v", result.TestName, result.PatientID)
		}

		if result.SampleType != "URI" {
			t.Errorf("Test %s: expected sampleType URI, got %v", result.TestName, result.SampleType)
		}
	}
}

// MockAnalyzerUseCase for testing
type MockAnalyzerUseCase struct {
	results *[]entity.VerifyResult
}

func (m *MockAnalyzerUseCase) ProcessOULR22(ctx context.Context, data entity.OUL_R22) error {
	return nil
}

func (m *MockAnalyzerUseCase) ProcessQBPQ11(ctx context.Context, data entity.QBP_Q11) error {
	return nil
}

func (m *MockAnalyzerUseCase) ProcessORMO01(ctx context.Context, data entity.ORM_O01) ([]entity.Specimen, error) {
	return nil, nil
}

func (m *MockAnalyzerUseCase) ProcessORUR01(ctx context.Context, data entity.ORU_R01) error {
	return nil
}

func (m *MockAnalyzerUseCase) ProcessCoax(ctx context.Context, data entity.CoaxTestResult) error {
	return nil
}

func (m *MockAnalyzerUseCase) ProcessDiestro(ctx context.Context, data entity.DiestroResult) error {
	return nil
}

func (m *MockAnalyzerUseCase) ProcessVerifyU120(ctx context.Context, result entity.VerifyResult) error {
	*m.results = append(*m.results, result)
	return nil
}
