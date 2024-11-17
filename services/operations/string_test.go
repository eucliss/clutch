package operations

import (
	"strconv"
	"testing"
)

func TestReplace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		value    string
		expected string
	}{
		{
			name:     "basic replacement",
			input:    "original",
			value:    "replaced",
			expected: "replaced",
		},
		{
			name:     "empty input",
			input:    "",
			value:    "replaced",
			expected: "replaced",
		},
		{
			name:     "empty value",
			input:    "original",
			value:    "",
			expected: "",
		},
		{
			name:     "special characters",
			input:    "original",
			value:    "!@#$%^",
			expected: "!@#$%^",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Replace(tt.input, tt.value)
			if result != tt.expected {
				t.Errorf("Replace() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRandomInt(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		upperLimit  string
		lowerLimit  string
		shouldError bool
	}{
		{
			name:        "valid range",
			input:       "original",
			upperLimit:  "10",
			lowerLimit:  "1",
			shouldError: false,
		},
		{
			name:        "equal limits",
			input:       "original",
			upperLimit:  "5",
			lowerLimit:  "3",
			shouldError: false,
		},
		{
			name:        "equal limits",
			input:       "original",
			upperLimit:  "50",
			lowerLimit:  "30",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RandomInt(tt.input, tt.upperLimit, tt.lowerLimit)

			if tt.shouldError {
				if result != tt.input {
					t.Errorf("RandomInt() with invalid input should return original input, got %v", result)
				}
				return
			}

			// Convert result to int for validation
			resultInt, err := strconv.Atoi(result)
			if err != nil {
				t.Errorf("RandomInt() returned invalid integer: %v", result)
				return
			}

			// Convert limits to int for comparison
			upper, _ := strconv.Atoi(tt.upperLimit)
			lower, _ := strconv.Atoi(tt.lowerLimit)

			// Check if result is within bounds
			if resultInt < lower || resultInt >= upper {
				t.Errorf("RandomInt() = %v, want value between %v and %v", resultInt, lower, upper)
			}
		})
	}
}

func TestRandomIntMultipleCalls(t *testing.T) {
	// Test that multiple calls produce different results
	results := make(map[string]bool)
	iterations := 100
	upperLimit := "1000"
	lowerLimit := "1"

	for i := 0; i < iterations; i++ {
		result := RandomInt("test", upperLimit, lowerLimit)
		results[result] = true
	}

	// Check if we got different results (at least 50% unique values)
	if len(results) < iterations/2 {
		t.Errorf("RandomInt() not producing enough unique values, got %v unique values out of %v iterations", len(results), iterations)
	}
}
