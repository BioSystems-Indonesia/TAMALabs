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

	resp, err := h.testTypeUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, resp)
}

func (h *TestTypeHandler) ListTestTypeFilter(c echo.Context) error {
	var req entity.TestTypeGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.testTypeUsecase.ListAllFilter(
		c.Request().Context(),
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, resp)
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

func (h *TestTypeHandler) GetOneTestTypeByCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return handleError(c, entity.ErrBadRequest.WithInternal(nil))
	}

	testType, err := h.testTypeUsecase.FindOneByCode(c.Request().Context(), code)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) GetOneTestTypeByAliasCode(c echo.Context) error {
	aliasCode := c.Param("alias_code")
	if aliasCode == "" {
		return handleError(c, entity.ErrBadRequest.WithInternal(nil))
	}

	testType, err := h.testTypeUsecase.FindOneByAliasCode(c.Request().Context(), aliasCode)
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

func (h *TestTypeHandler) DeleteTestType(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req, err := h.testTypeUsecase.FindOneByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	testType, err := h.testTypeUsecase.Delete(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testType)
}

func (h *TestTypeHandler) UploadBulkTestType(c echo.Context) error {
	mf, err := c.FormFile("file")
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	f, err := mf.Open()
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}
	defer f.Close()

	err = h.testTypeUsecase.BulkCreate(c.Request().Context(), f)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
