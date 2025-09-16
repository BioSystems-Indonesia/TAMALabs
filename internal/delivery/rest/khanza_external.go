package rest

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	khanzauc "github.com/oibacidem/lims-hl-seven/internal/usecase/external/khanza"
)

type KhanzaExternalHandler struct {
	usecase *khanzauc.Usecase
}

func NewKhanzaExternalHandler(usecase *khanzauc.Usecase) *KhanzaExternalHandler {
	return &KhanzaExternalHandler{
		usecase: usecase,
	}
}

func (h *KhanzaExternalHandler) ProcessRequest(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return handleError(c, err)
	}

	err = h.usecase.ProcessRequest(c.Request().Context(), body)
	if err != nil {
		return handleError(c, err)
	}

	return c.String(http.StatusOK, "Request berhasil tersimpan.Silahkan cek aplikasi LIS")
}

func (h *KhanzaExternalHandler) GetResult(c echo.Context) error {
	id := c.Param("id")

	body, err := h.usecase.GetResult(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, body)
}

func (h *KhanzaExternalHandler) RegisterRoutes(g *echo.Group) {
	khanza := g.Group("/khanza")
	khanza.POST("/order", h.ProcessRequest)
	khanza.GET("/result/:user/:key/:id", h.GetResult)
}
