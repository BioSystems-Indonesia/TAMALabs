package helper

import (
	"strconv"
	"strings"
)

func ParseReferenceRange(refRange string) (float64, float64, error) {
	refRange = strings.TrimSpace(refRange)

	parts := strings.Split(refRange, "-")
	if len(parts) != 2 {
		return 0, 0, nil // Return zeros if format is invalid
	}

	minVal, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return 0, 0, err
	}

	maxVal, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return 0, 0, err
	}

	return minVal, maxVal, nil
}

func DetermineFlagFromReference(value string, referenceRange string) string {
	if referenceRange == "" {
		return ""
	}

	val, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return "" // If value is not numeric, return empty flag
	}

	minVal, maxVal, err := ParseReferenceRange(referenceRange)
	if err != nil {
		return "" // If parsing fails, return empty flag
	}

	if val < minVal {
		return "L"
	} else if val > maxVal {
		return "H"
	}

	return ""
}
