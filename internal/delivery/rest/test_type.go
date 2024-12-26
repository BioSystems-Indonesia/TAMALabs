package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/test_type"
)

type TestTypeHandler struct {
	cfg             *config.Schema
	TestTypeUsecase *test_type.Usecase
}

func NewTestTypeHandler(cfg *config.Schema, TestTypeUsecase *test_type.Usecase) *TestTypeHandler {
	return &TestTypeHandler{cfg: cfg, TestTypeUsecase: TestTypeUsecase}
}

func (h *TestTypeHandler) FindTestType(c echo.Context) error {
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
