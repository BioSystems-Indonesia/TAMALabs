package rest

import (
	"net/http"
	"strconv"

	test_template_uc "github.com/oibacidem/lims-hl-seven/internal/usecase/test_template"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type TestTemplateHandler struct {
	cfg                 *config.Schema
	testTemplateUsecase *test_template_uc.Usecase
}

func NewTestTemplateHandler(cfg *config.Schema, TestTemplateUsecase *test_template_uc.Usecase) *TestTemplateHandler {
	return &TestTemplateHandler{cfg: cfg, testTemplateUsecase: TestTemplateUsecase}
}

func (h *TestTemplateHandler) ListTestTemplate(c echo.Context) error {
	var req entity.TestTemplateGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.testTemplateUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, resp)
}

func (h *TestTemplateHandler) GetOneTestTemplate(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	testTemplate, err := h.testTemplateUsecase.FindOneByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testTemplate)
}

func (h *TestTemplateHandler) CreateTestTemplate(c echo.Context) error {
	var req entity.TestTemplate
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	user := entity.GetEchoContextUser(c)
	req.CreatedBy = user.ID
	req.LastUpdatedBy = user.ID

	testTemplate, err := h.testTemplateUsecase.Create(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, testTemplate)
}

func (h *TestTemplateHandler) UpdateTestTemplate(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	var req entity.TestTemplate
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	user := entity.GetEchoContextUser(c)

	req.ID = id
	req.LastUpdatedBy = user.ID
	testTemplate, err := h.testTemplateUsecase.Update(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testTemplate)
}

func (h *TestTemplateHandler) DeleteTestTemplate(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req, err := h.testTemplateUsecase.FindOneByID(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	testTemplate, err := h.testTemplateUsecase.Delete(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, testTemplate)
}
