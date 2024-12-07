package rest

import (
	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
)

type HealthCheckHandler struct {
	cfg *config.Schema
}

func NewHealthCheckHandler(cfg *config.Schema) *HealthCheckHandler {
	return &HealthCheckHandler{cfg: cfg}
}

func (h HealthCheckHandler) Ping(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status":   "OK",
		"name":     h.cfg.Name,
		"version":  h.cfg.Version,
		"revision": h.cfg.Revision,
		"logLevel": h.cfg.LogLevel,
	})
}
