package rest

import (
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	externaluc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external"
	"github.com/labstack/echo/v4"
)

// ExternalHandler handles HTTP requests for external integration
type ExternalHandler struct {
	usecase *externaluc.Usecase
}

// NewExternalHandler creates a new ExternalHandler instance
func NewExternalHandler(usecase *externaluc.Usecase) *ExternalHandler {
	return &ExternalHandler{
		usecase: usecase,
	}
}

// RegisterRoutes registers all external routes
func (h *ExternalHandler) RegisterRoutes(router *echo.Group) {
	external := router.Group("/external")
	external.POST("/sync-all-results", h.SyncAllResults)
	external.POST("/sync-all-requests", h.SyncAllRequests)
	external.POST("/simrs/test-connection", h.TestSimrsConnection)
}

// SyncAllResults syncs all completed work order results to external system
func (h *ExternalHandler) SyncAllResults(c echo.Context) error {
	// No need to bind request body, we'll automatically get all completed work orders
	err := h.usecase.SyncAllResult(c.Request().Context(), nil) // Pass nil to sync all completed orders
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "All results synced successfully"})
}

func (h *ExternalHandler) SyncAllRequests(c echo.Context) error {
	var req entity.ExternalSyncAllResultsRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	err := h.usecase.SyncAllRequest(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "All requests synced successfully"})
}

// TestSimrsConnection tests SIMRS database connection
func (h *ExternalHandler) TestSimrsConnection(c echo.Context) error {
	var req struct {
		DSN string `json:"dsn" validate:"required"`
	}

	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	err := h.usecase.TestSimrsConnection(c.Request().Context(), req.DSN)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":   "Connection failed",
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Connection successful"})
}
