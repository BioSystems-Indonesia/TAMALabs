package middleware

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ConfigGetter interface to get configuration values
type ConfigGetter interface {
	Get(ctx context.Context, key string) (string, error)
}

// IntegrationCheckMiddleware checks if specific integration is enabled
type IntegrationCheckMiddleware struct {
	configGetter ConfigGetter
}

// NewIntegrationCheckMiddleware creates a new integration check middleware
func NewIntegrationCheckMiddleware(configGetter ConfigGetter) *IntegrationCheckMiddleware {
	return &IntegrationCheckMiddleware{
		configGetter: configGetter,
	}
}

// CheckSimrsEnabled returns middleware that checks if SIMRS integration is enabled
func (m *IntegrationCheckMiddleware) CheckSimrsEnabled() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//ctx := c.Request().Context()

			//// Check if SIMRS integration is enabled
			//simrsEnabled, err := m.configGetter.Get(ctx, "SimrsIntegrationEnabled")
			//if err != nil || simrsEnabled != "true" {
			//	return echo.NewHTTPError(http.StatusForbidden, "SIMRS integration is not enabled")
			//}
			//
			//// Check if SIMRS is selected as active integration (either "simrs" or "simrs-api")
			//selectedSimrs, err := m.configGetter.Get(ctx, "SelectedSimrs")
			//if err != nil || (selectedSimrs != "simrs" && selectedSimrs != "simrs-api") {
			//	return echo.NewHTTPError(http.StatusForbidden, "SIMRS is not the active integration")
			//}

			return next(c)
		}
	}
}

// CheckTechnoMedicEnabled returns middleware that checks if TechnoMedic integration is enabled
func (m *IntegrationCheckMiddleware) CheckTechnoMedicEnabled() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Check if TechnoMedic integration is enabled
			enabled, err := m.configGetter.Get(ctx, "TechnoMedicIntegrationEnabled")
			if err != nil || enabled != "true" {
				return echo.NewHTTPError(http.StatusForbidden, "TechnoMedic integration is not enabled")
			}

			return next(c)
		}
	}
}

// CheckSimgosEnabled returns middleware that checks if Database Sharing integration is enabled
func (m *IntegrationCheckMiddleware) CheckSimgosEnabled() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Check if Database Sharing integration is enabled
			enabled, err := m.configGetter.Get(ctx, "SimgosIntegrationEnabled")
			if err != nil || enabled != "true" {
				return echo.NewHTTPError(http.StatusForbidden, "Database Sharing integration is not enabled")
			}

			// Check if Database Sharing is selected as active integration
			selectedSimrs, err := m.configGetter.Get(ctx, "SelectedSimrs")
			if err != nil || selectedSimrs != "simgos" {
				return echo.NewHTTPError(http.StatusForbidden, "Database Sharing is not the active integration")
			}

			return next(c)
		}
	}
}

// CheckKhanzaEnabled returns middleware that checks if Khanza integration is enabled
func (m *IntegrationCheckMiddleware) CheckKhanzaEnabled() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Check if Khanza integration is enabled
			khanzaEnabled, err := m.configGetter.Get(ctx, "KhanzaIntegrationEnabled")
			if err != nil || khanzaEnabled != "true" {
				return echo.NewHTTPError(http.StatusForbidden, "Khanza integration is not enabled")
			}

			// Check if Khanza is selected as active integration
			selectedSimrs, err := m.configGetter.Get(ctx, "SelectedSimrs")
			if err != nil || selectedSimrs != "khanza" {
				return echo.NewHTTPError(http.StatusForbidden, "Khanza is not the active integration")
			}

			return next(c)
		}
	}
}
