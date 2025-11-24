package rest

import (
	"io"
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/middleware"
	simrsuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/simrs"
	"github.com/labstack/echo/v4"
)

type SimrsExternalHandler struct {
	usecase               *simrsuc.Usecase
	integrationMiddleware *middleware.IntegrationCheckMiddleware
}

func NewSimrsExternalHandler(usecase *simrsuc.Usecase, integrationMiddleware *middleware.IntegrationCheckMiddleware) *SimrsExternalHandler {
	return &SimrsExternalHandler{
		usecase:               usecase,
		integrationMiddleware: integrationMiddleware,
	}
}

func (h *SimrsExternalHandler) ProcessOrder(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return handleError(c, err)
	}

	err = h.usecase.ProcessOrder(c.Request().Context(), body)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Order berhasil diterima dan sedang diproses",
	})
}

func (h *SimrsExternalHandler) GetResult(c echo.Context) error {
	orderLab := c.Param("orderlab")
	if orderLab == "" {
		return handleError(c, echo.NewHTTPError(http.StatusBadRequest, "order lab is required"))
	}

	result, err := h.usecase.GetResult(c.Request().Context(), orderLab)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *SimrsExternalHandler) DeleteOrder(c echo.Context) error {
	orderLab := c.Param("orderlab")
	if orderLab == "" {
		return handleError(c, echo.NewHTTPError(http.StatusBadRequest, "order lab is required"))
	}

	err := h.usecase.DeleteOrder(c.Request().Context(), orderLab)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Order berhasil dihapus",
	})
}

func (h *SimrsExternalHandler) RegisterRoutes(g *echo.Group) {
	simrs := g.Group("/his", h.integrationMiddleware.CheckSimrsEnabled())
	simrs.POST("/order", h.ProcessOrder)
	simrs.GET("/result/:orderlab", h.GetResult)
	simrs.DELETE("/order/:orderlab", h.DeleteOrder)
}
