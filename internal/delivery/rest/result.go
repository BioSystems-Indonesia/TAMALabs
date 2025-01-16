package rest

import (
	"net/http"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/internal/entity"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
)

type ResultHandler struct {
	cfg           *config.Schema
	resultUsecase usecase.Result
}

func NewResultHandler(cfg *config.Schema, resultUsecase usecase.Result) *ResultHandler {
	return &ResultHandler{
		cfg:           cfg,
		resultUsecase: resultUsecase,
	}
}

func (h *ResultHandler) ListResult(c echo.Context) error {
	results, err := h.resultUsecase.Results(c.Request().Context(), nil)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(results)))
	return c.JSON(http.StatusOK, results)
}

func (h *ResultHandler) GetResult(c echo.Context) error {
	barcode := c.Param("barcode")
	result, err := h.resultUsecase.ResultDetail(c.Request().Context(), barcode)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
