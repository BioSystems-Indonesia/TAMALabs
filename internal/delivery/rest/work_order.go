package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
)

type WorkOrderHandler struct {
	cfg              *config.Schema
	workOrderUsecase *workOrderuc.WorkOrderUseCase
}

func NewWorkOrderHandler(cfg *config.Schema, workOrderUsecase *workOrderuc.WorkOrderUseCase) *WorkOrderHandler {
	return &WorkOrderHandler{cfg: cfg, workOrderUsecase: workOrderUsecase}
}

func (h *WorkOrderHandler) FindWorkOrders(c echo.Context) error {
	var req entity.WorkOrderGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	workOrders, err := h.workOrderUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(workOrders)))
	return c.JSON(http.StatusOK, workOrders)
}

func (h *WorkOrderHandler) GetOneWorkOrder(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	workOrder, err := h.workOrderUsecase.FindOneByID(id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, workOrder)
}

func (h *WorkOrderHandler) DeleteWorkOrder(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	if err := h.workOrderUsecase.Delete(id); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.WorkOrder{
		ID: id,
	})
}

func (h *WorkOrderHandler) UpdateWorkOrder(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.WorkOrder{ID: id}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.workOrderUsecase.Update(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, req)
}

func (h *WorkOrderHandler) CreateWorkOrder(c echo.Context) error {
	var req entity.WorkOrder
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.workOrderUsecase.Create(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}

func (h Handler) RunWorkOrder(c echo.Context) error {
	var req entity.WorkOrderRunRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	workOrders, err := h.workOrderUsecase.FindManyByID(c.Request().Context(), req.WorkOrderID)
	if err != nil {
		return handleError(c, err)
	}

	if len(workOrders) == 0 {
		return handleError(c, entity.ErrNotFound.WithInternal(errors.New("work order not found")))
	}

	// TODO: Support multiple work order by merge the patient
	err = ba400.SendToBA400(c.Request().Context(), workOrders[0].Patient)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
