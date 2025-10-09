package cbs400

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCBS400Line(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		wantResults int
		checkFirst  bool
		firstParam  string
		firstValue  float64
		firstUnit   string
	}{
		{
			name:        "Valid CBS400 result line",
			input:       "r371                   3.91 141.2 106.5  1.21  1.23  2.46  7.44",
			wantErr:     false,
			wantResults: 7,
			checkFirst:  true,
			firstParam:  "K",
			firstValue:  3.91,
			firstUnit:   "mmol/L",
		},
		{
			name:        "Valid CBS400 result line with different patient ID",
			input:       "r123                   4.12 140.5 105.2  1.18  1.20  2.38  7.42",
			wantErr:     false,
			wantResults: 7,
			checkFirst:  true,
			firstParam:  "K",
			firstValue:  4.12,
			firstUnit:   "mmol/L",
		},
		{
			name:    "Invalid line - not CBS400 format",
			input:   "invalid line format",
			wantErr: true,
		},
		{
			name:    "Invalid line - missing values",
			input:   "r123   3.91 141.2",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := parseCBS400Line(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.wantResults)

			if tt.checkFirst && len(results) > 0 {
				first := results[0]
				assert.Equal(t, tt.firstParam, first.TestName)
				assert.Equal(t, tt.firstValue, first.Value)
				assert.Equal(t, tt.firstUnit, first.Unit)
				assert.Equal(t, "SER", first.SampleType)
				assert.NotZero(t, first.Timestamp)
			}
		})
	}
}

func TestCBS400Parameters(t *testing.T) {
	// Test that we have the correct parameters in the right order
	expectedParams := []string{"K", "Na", "Cl", "iCa", "nCa", "TCa", "pH"}
	expectedUnits := []string{"mmol/L", "mmol/L", "mmol/L", "mmol/L", "mmol/L", "mmol/L", ""}

	assert.Len(t, cbs400Parameters, len(expectedParams))

	for i, param := range cbs400Parameters {
		assert.Equal(t, expectedParams[i], param.Name, "Parameter %d name mismatch", i)
		assert.Equal(t, expectedUnits[i], param.Unit, "Parameter %d unit mismatch", i)
	}
}

func TestProcessBuffer(t *testing.T) {
	// This would require a mock usecase to properly test
	// For now, we'll just ensure the function doesn't panic
	handler := &Handler{
		analyzerUseCase: nil, // This would be a mock in real tests
		buffer:          "",
	}

	// Test with empty buffer
	handler.processBuffer()

	// Test with sample data
	handler.buffer = "r371                   3.91 141.2 106.5  1.21  1.23  2.46  7.44\n"
	// Note: This will log errors due to nil usecase, but shouldn't panic
}

func TestIsValidCBS400Value(t *testing.T) {
	tests := []struct {
		param    string
		value    float64
		expected bool
	}{
		// Valid values
		{"K", 3.91, true},
		{"Na", 141.2, true},
		{"Cl", 106.5, true},
		{"iCa", 1.21, true},
		{"nCa", 1.23, true},
		{"TCa", 2.46, true},
		{"pH", 7.44, true},

		// Invalid values - too low
		{"K", 0.5, false},
		{"Na", 50.0, false},
		{"Cl", 30.0, false},
		{"iCa", 0.2, false},
		{"nCa", 0.3, false},
		{"TCa", 0.5, false},
		{"pH", 5.0, false},

		// Invalid values - too high
		{"K", 15.0, false},
		{"Na", 250.0, false},
		{"Cl", 200.0, false},
		{"iCa", 5.0, false},
		{"nCa", 4.0, false},
		{"TCa", 8.0, false},
		{"pH", 9.0, false},

		// Unknown parameter (should be valid)
		{"UNKNOWN", 999.0, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%.2f", tt.param, tt.value), func(t *testing.T) {
			result := isValidCBS400Value(tt.param, tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
