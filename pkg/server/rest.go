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
	GetClient() *echo.Echo
}

// Rest is a struct that contains the rest server
type Rest struct {
	Port   string
	Client *echo.Echo
}

func (r *Rest) Serve() {
	// Start server in a goroutine so that it doesn't block.
	go func() {
		if err := r.Client.Start("localhost:" + r.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error starting server: %v", err)
		}
	}()
	log.Printf("Server started at %s", r.Port)

	// Wait for interrupt signal to gracefully shut down the server with a timeout.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

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
	return &Rest{
		Port:   port,
		Client: e,
	}
}
