package rest

import (
	"io"
	"net/http"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/internal/usecase/analyzer"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/result"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
)

type ResultHandler struct {
	cfg             *config.Schema
	resultUsecase   *result.Usecase
	analyzerUsecase *analyzer.Usecase
}

func NewResultHandler(cfg *config.Schema, resultUsecase *result.Usecase, analyzerUsecase *analyzer.Usecase) *ResultHandler {
	return &ResultHandler{
		cfg:             cfg,
		resultUsecase:   resultUsecase,
		analyzerUsecase: analyzerUsecase,
	}
}

func (h *ResultHandler) ListResult(c echo.Context) error {
	var req entity.ResultGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.resultUsecase.Results(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, resp)
}

func (h *ResultHandler) GetResult(c echo.Context) error {
	workOrderID, err := strconv.ParseInt(c.Param("work_order_id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	result, err := h.resultUsecase.ResultDetail(c.Request().Context(), workOrderID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) AddTestResult(c echo.Context) error {
	req := entity.TestResult{}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.resultUsecase.PutTestResult(c.Request().Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) DeleteTestResult(c echo.Context) error {
	testResultID, err := strconv.ParseInt(c.Param("test_result_id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	result, err := h.resultUsecase.DeleteTestResult(c.Request().Context(), testResultID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) TooglePickTestResult(c echo.Context) error {
	testResultID, err := strconv.ParseInt(c.Param("test_result_id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	result, err := h.resultUsecase.TooglePickTestResult(c.Request().Context(), testResultID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h Handler) RefreshResult(c echo.Context) error {
	err := h.analyzerUsecase.ProcessA15(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h ResultHandler) ApproveResult(c echo.Context) error {
	workOrderID, err := strconv.ParseInt(c.Param("work_order_id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	admin := entity.GetEchoContextUser(c)
	err = h.resultUsecase.ApproveResult(c.Request().Context(), workOrderID, admin.ID)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h ResultHandler) RejectResult(c echo.Context) error {
	workOrderID, err := strconv.ParseInt(c.Param("work_order_id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	admin := entity.GetEchoContextUser(c)
	err = h.resultUsecase.RejectResult(c.Request().Context(), workOrderID, admin.ID)
	if err != nil {
		return handleError(c, err)
	}

	return c.NoContent(http.StatusOK)
}

func (h ResultHandler) UploadFileA15(c echo.Context) error {
	src, err := c.FormFile("file")
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	file, err := src.Open()
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	h.analyzerUsecase.FileResult(c.Request().Context(), string(data))

	return c.NoContent(http.StatusOK)
}
