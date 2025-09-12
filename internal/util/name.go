package util

import "strings"

// SplitName splits a full name into first and last name components
func SplitName(name string) (firstName, lastName string) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", ""
	}

	parts := strings.Fields(name)

	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	case 2:
		return parts[0], parts[1]
	default:
		return parts[0], strings.Join(parts[1:], " ")
	}
}
