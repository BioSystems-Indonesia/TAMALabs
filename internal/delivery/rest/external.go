package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	externaluc "github.com/oibacidem/lims-hl-seven/internal/usecase/external"
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
}

// SyncAllResults syncs all completed work order results to external system
func (h *ExternalHandler) SyncAllResults(c echo.Context) error {
	var req entity.ExternalSyncAllResultsRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	err := h.usecase.SyncAllResult(c.Request().Context(), req.OrderIDs)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "All results synced successfully"})
}
