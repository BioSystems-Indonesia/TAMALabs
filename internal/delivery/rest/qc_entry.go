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
	qc.POST("/results/manual", h.CreateManualQCResult)
	qc.PUT("/entries/:id/select-result", h.UpdateSelectedQCResult)
	qc.GET("/statistics", h.GetQCStatistics)
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

func (h *QCEntryHandler) CreateManualQCResult(c echo.Context) error {
	var req entity.CreateManualQCResultRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	// Get logged-in user from context
	claims := entity.GetEchoContextUser(c)
	if claims.ID == 0 {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	result, err := h.qcUsecase.CreateManualQCResult(c.Request().Context(), &req, claims.Fullname)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Manual QC result created successfully",
		"data":    result,
	})
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

func (h *QCEntryHandler) GetQCStatistics(c echo.Context) error {
	ctx := c.Request().Context()

	deviceIDStr := c.QueryParam("device_id")
	if deviceIDStr == "" {
		return handleError(c, entity.ErrBadRequest.WithInternal(nil))
	}

	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	stats, err := h.qcUsecase.GetQCStatistics(ctx, deviceID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *QCEntryHandler) UpdateSelectedQCResult(c echo.Context) error {
	ctx := c.Request().Context()

	// Get QC entry ID from path
	qcEntryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	// Parse request body
	var req struct {
		QCLevel  int `json:"qc_level" validate:"required,min=1,max=3"`
		ResultID int `json:"result_id" validate:"required"`
	}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	// Update selected result (method selection dilakukan di frontend via filter)
	if err := h.qcUsecase.UpdateSelectedQCResult(ctx, qcEntryID, req.QCLevel, req.ResultID); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Selected QC result updated successfully",
	})
}
