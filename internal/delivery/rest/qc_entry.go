package rest

import (
	"net/http"
	"strconv"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/labstack/echo/v4"
)

type QCEntryHandler struct {
	qcUsecase usecase.QualityControl
}

func NewQCEntryHandler(qcUsecase usecase.QualityControl) *QCEntryHandler {
	return &QCEntryHandler{
		qcUsecase: qcUsecase,
	}
}

func (h *QCEntryHandler) RegisterRoute(qc *echo.Group) {
	qc.POST("/entries", h.CreateQCEntry)
	qc.GET("/entries", h.ListQCEntries)
	qc.GET("/results", h.ListQCResults)
}

func (h *QCEntryHandler) CreateQCEntry(c echo.Context) error {
	var req entity.CreateQCEntryRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	entry, err := h.qcUsecase.CreateQCEntry(c.Request().Context(), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, entry)
}

func (h *QCEntryHandler) ListQCEntries(c echo.Context) error {
	ctx := c.Request().Context()

	var req entity.GetManyRequestQCEntry
	if err := c.Bind(&req); err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	entries, total, err := h.qcUsecase.GetQCEntries(ctx, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.PaginationResponse[entity.QCEntry]{
		Data:  entries,
		Total: total,
	})
}

func (h *QCEntryHandler) ListQCResults(c echo.Context) error {
	ctx := c.Request().Context()

	var req entity.GetManyRequestQCResult
	if err := c.Bind(&req); err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	// Parse qc_entry_id if provided
	if entryIDStr := c.QueryParam("qc_entry_id"); entryIDStr != "" {
		entryID, err := strconv.Atoi(entryIDStr)
		if err == nil {
			req.QCEntryID = &entryID
		}
	}

	results, total, err := h.qcUsecase.GetQCResults(ctx, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.PaginationResponse[entity.QCResult]{
		Data:  results,
		Total: total,
	})
}
