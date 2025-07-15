package panics

import (
	"context"
	"log/slog"
	"runtime/debug"
)

// CapturePanic is a function that captures a panic and logs it.
func CapturePanic(ctx context.Context, f func()) {
	defer RecoverPanic(ctx)

	f()
}

func RecoverPanic(ctx context.Context) {
	if r := recover(); r != nil {
		// A panic occurred. Log the recovered panic value (the error).
		slog.ErrorContext(
			ctx,
			"Recovered from panic in",
			"panic",
			r,
			"stack",
			string(debug.Stack()),
		)
	}
}
