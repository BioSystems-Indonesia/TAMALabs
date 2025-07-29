package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// RestServer is an interface for rest server
type RestServer interface {
	Serve()
	Stop() error
	GetClient() *echo.Echo
}

// Rest is a struct that contains the rest server
type Rest struct {
	Port   string
	Client *echo.Echo
	ctx    context.Context
	cancel context.CancelFunc
}

func (r *Rest) Serve() {
	// Create context for graceful shutdown
	r.ctx, r.cancel = context.WithCancel(context.Background())

	// Start server in a goroutine so that it doesn't block.
	errChan := make(chan error, 1)
	go func() {
		if err := r.Client.Start("0.0.0.0:" + r.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("error starting server: %w", err)
		}
	}()
	slog.Info("Server started at", slog.String("port", r.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Wait for error, interrupt signal or context cancellation
	select {
	case <-r.ctx.Done():
		slog.Info("Context cancelled, shutting down server...")
		panic("context cancelled")
	case err := <-errChan:
		slog.Error("Error starting server", slog.String("error", err.Error()))
		panic(fmt.Sprintf("error starting server: %v", err))
	case <-quit:
		slog.Info("Interrupt signal received, shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Client.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
	}
}

func (r *Rest) GetClient() *echo.Echo {
	return r.Client
}

// Stop gracefully stops the server
func (r *Rest) Stop() error {
	if r.cancel != nil {
		r.cancel()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.Client.Shutdown(ctx)
}

type Validator struct {
	v *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.v.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

// NewRest returns a new rest server
func NewRest(port string, validate *validator.Validate) RestServer {
	e := echo.New()
	e.Validator = &Validator{v: validate}
	e.HideBanner = true
	return &Rest{
		Port:   port,
		Client: e,
	}
}
