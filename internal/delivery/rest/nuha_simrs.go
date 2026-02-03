package rest

import (
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/middleware"
	nuha_simrs "github.com/BioSystems-Indonesia/TAMALabs/internal/services/nuha-simrs"
	"github.com/labstack/echo/v4"
)

type NuhaSIMRSHandler struct {
	simrsNuha             *nuha_simrs.SIMRSNuha
	integrationMiddleware *middleware.IntegrationCheckMiddleware
}

func NewNuhaSIMRSHandler(
	simrsNuha *nuha_simrs.SIMRSNuha,
	integrationMiddleware *middleware.IntegrationCheckMiddleware,
) *NuhaSIMRSHandler {
	return &NuhaSIMRSHandler{
		simrsNuha:             simrsNuha,
		integrationMiddleware: integrationMiddleware,
	}
}

// SyncLabOrders fetches lab orders from Nuha SIMRS and creates work orders
func (h *NuhaSIMRSHandler) SyncLabOrders(c echo.Context) error {
	ctx := c.Request().Context()

	err := h.simrsNuha.GetLabOrder(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Successfully synced lab orders from Nuha SIMRS",
	})
}

func (h *NuhaSIMRSHandler) RegisterRoutes(g *echo.Group) {
	nuha := g.Group("/nuha-simrs")

	// Apply Nuha SIMRS integration check middleware
	nuha.Use(h.integrationMiddleware.CheckNuhaEnabled())

	// Register routes
	nuha.POST("/sync-orders", h.SyncLabOrders)
}
