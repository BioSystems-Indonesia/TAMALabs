package util

import "testing"

func TestGenerateRandomDigits(t *testing.T) {
	tests := []struct {
		name      string
		n         int
		wantErr   bool
		expectLen int
	}{
		{
			name:      "Generate 10 digits",
			n:         10,
			wantErr:   false,
			expectLen: 10,
		},
		{
			name:      "Generate 0 digits",
			n:         0,
			wantErr:   true,
			expectLen: 0,
		},
		{
			name:      "Generate 1 digit",
			n:         1,
			wantErr:   false,
			expectLen: 1,
		},
		{
			name:      "Generate negative digits",
			n:         -5,
			wantErr:   true,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomDigits(tt.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomDigits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.expectLen {
				t.Errorf("GenerateRandomDigits() length = %d, want %d", len(got), tt.expectLen)
			}

			// Ensure all characters are digits
			if !tt.wantErr {
				for _, c := range got {
					if c < '0' || c > '9' {
						t.Errorf("GenerateRandomDigits() contains non-digit character: %c", c)
					}
				}
			}
		})
	}
}
