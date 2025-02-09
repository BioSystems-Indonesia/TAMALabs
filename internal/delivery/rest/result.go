package rest

import (
	"net/http"
	"strconv"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase/result"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
)

type ResultHandler struct {
	cfg           *config.Schema
	resultUsecase *result.Usecase
}

func NewResultHandler(cfg *config.Schema, resultUsecase *result.Usecase) *ResultHandler {
	return &ResultHandler{
		cfg:           cfg,
		resultUsecase: resultUsecase,
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	result, err := h.resultUsecase.ResultDetail(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) CreateResult(c echo.Context) error {
	req := entity.ObservationResultCreate{}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.resultUsecase.CreateResult(c.Request().Context(), req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) DeleteResult(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	result, err := h.resultUsecase.DeleteResult(c.Request().Context(), id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) DeleteResultBulk(c echo.Context) error {
	req := entity.DeleteResultBulkReq{}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.resultUsecase.DeleteResultBulk(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ResultHandler) UpdateResult(c echo.Context) error {
	req := entity.UpdateManyResultTestReq{}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.resultUsecase.UpdateResult(c.Request().Context(), req.Data)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
