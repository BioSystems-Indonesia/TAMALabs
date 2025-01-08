package rest

import (
	"net/http"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/internal/usecase/test_type"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type TestTypeHandler struct {
	cfg             *config.Schema
	testTypeUsecase *test_type.Usecase
}

func NewTestTypeHandler(cfg *config.Schema, TestTypeUsecase *test_type.Usecase) *TestTypeHandler {
	return &TestTypeHandler{cfg: cfg, testTypeUsecase: TestTypeUsecase}
}

func (h *TestTypeHandler) ListTestType(c echo.Context) error {
	var req entity.TestTypeGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	TestTypes, err := h.testTypeUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(TestTypes)))
	return c.JSON(http.StatusOK, TestTypes)
}

func (h *TestTypeHandler) GetOneTestType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	testType, err := h.testTypeUsecase.FindOneByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) CreateTestType(c echo.Context) error {
	var req entity.TestType
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	testType, err := h.testTypeUsecase.Create(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, testType)
}

func (h *TestTypeHandler) UpdateTestType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	var req entity.TestType
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	req.ID = id
	testType, err := h.testTypeUsecase.Update(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}
