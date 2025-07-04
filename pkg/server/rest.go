package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

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
	go func() {
		if err := r.Client.Start("0.0.0.0:" + r.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	log.Printf("Server started at %s", r.Port)

	// Wait for interrupt signal or context cancellation to gracefully shut down the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		log.Println("Interrupt signal received, shutting down server...")
	case <-r.ctx.Done():
		log.Println("Context cancelled, shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.Client.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
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
