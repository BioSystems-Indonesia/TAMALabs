package rest

import (
	"io"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/external"
)

type HRISExternalHandler struct {
	UC external.UseCase
}

func (h *HRISExternalHandler) ProcessRequest(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	err = h.UC.ProcessRequest(c.Request().Context(), body)
	if err != nil {
		return err
	}

	_, err = c.Response().Write([]byte("OK"))
	return err
}


func (h *HRISExternalHandler) GetResult(c echo.Context) error {
	id := c.Param("id")

	body, err := h.UC.GetResult(c.Request().Context(), id)
	if err != nil {
		return err
	}

	_, err = c.Response().Write(body)
	return err
}

func (h *HRISExternalHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/order", h.ProcessRequest)
	g.GET("/result/user/key/:id", h.GetResult)
}
