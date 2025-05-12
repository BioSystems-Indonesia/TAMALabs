package panics

import (
	"context"
	"fmt"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

// CapturePanic is a function that captures a panic and logs it.
func CapturePanic(ctx context.Context, f func()) {
	defer func() {
		if r := recover(); r != nil {
			// A panic occurred. Log the recovered panic value (the error).
			slog.ErrorContext(ctx, fmt.Sprintf("Recovered from panic in goroutine: %v\n", r))

			// Log the stack trace for debugging purposes.
			// debug.PrintStack() prints the stack trace of the current goroutine to standard error.
			slog.ErrorContext(ctx, "Stack trace:")
			debug.PrintStack()
		}
	}()

	f()
}
