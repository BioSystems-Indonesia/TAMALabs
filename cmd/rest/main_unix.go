//go:build !windows

package main

import "fmt"

// showErrorMessage prints the error to the console on non-Windows systems.
func showErrorMessage(title, message string) {
	fmt.Printf("Error: %s\n%s\n", title, message)
}
