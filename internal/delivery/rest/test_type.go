package rest

import (
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type TestTypeHandler struct {
	cfg             *config.Schema
	TestTypeUsecase usecase.TestType
}

func NewTestTypeHandler(cfg *config.Schema, TestTypeUsecase usecase.TestType) *TestTypeHandler {
	return &TestTypeHandler{cfg: cfg, TestTypeUsecase: TestTypeUsecase}
}

func (h *TestTypeHandler) ListTestType(c echo.Context) error {
	var req entity.TestTypeGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	TestTypes, err := h.TestTypeUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(TestTypes)))
	return c.JSON(http.StatusOK, TestTypes)
}
