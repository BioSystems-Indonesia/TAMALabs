package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	// Read version.go
	content, err := os.ReadFile("internal/constant/version.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading version.go: %v\n", err)
		os.Exit(1)
	}

	// Extract version using regex
	re := regexp.MustCompile(`const AppVersion = "([^"]+)"`)
	matches := re.FindSubmatch(content)
	if len(matches) < 2 {
		fmt.Fprintf(os.Stderr, "Could not find AppVersion in version.go\n")
		os.Exit(1)
	}

	version := string(matches[1])

	// Write version.ini
	iniContent := fmt.Sprintf("[Version]\nAppVersion=%s\n", version)
	err = os.WriteFile("version.ini", []byte(iniContent), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing version.ini: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated version.ini with version: %s\n", version)
}
