package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
	"gorm.io/gorm"
)

type WorkOrderHandler struct {
	cfg              *config.Schema
	workOrderUsecase *workOrderuc.WorkOrderUseCase
	db               *gorm.DB
}

func NewWorkOrderHandler(cfg *config.Schema, workOrderUsecase *workOrderuc.WorkOrderUseCase, db *gorm.DB) *WorkOrderHandler {
	return &WorkOrderHandler{cfg: cfg, workOrderUsecase: workOrderUsecase, db: db}
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

func (h *WorkOrderHandler) CreateWorkOrder(c echo.Context) error {
	var req entity.WorkOrder
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.workOrderUsecase.Create(&req); err != nil {
		return handleError(c, err)
	}

	workOrder, err := h.workOrderUsecase.FindOneByID(req.ID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, workOrder)
}

func (h *WorkOrderHandler) AddTestWorkOrder(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.WorkOrder{ID: id}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.workOrderUsecase.AddTest(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *WorkOrderHandler) RunWorkOrder(c echo.Context) error {
	var req entity.WorkOrderRunRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	querySpeciment := h.db.Model(entity.Specimen{}).Where("order_id = ?", req.WorkOrderID).Select("id")
	err := h.db.Model(entity.ObservationRequest{}).Where("specimen_id in (?)", querySpeciment).
		Update("result_status", constant.ResultStatusSpecimenPending).Error
	if err != nil {
		return handleError(c, err)
	}

	workOrders, err := h.workOrderUsecase.FindOneByID(req.WorkOrderID)
	if err != nil {
		return handleError(c, err)
	}

	device := entity.Device{}
	tx := h.db.First(&device, req.DeviceID)
	if tx.Error != nil {
		return handleError(c, fmt.Errorf("error finding device: %w", tx.Error))
	}

	err = ba400.SendToBA400(c.Request().Context(), workOrders.Patient, device)
	if err != nil {
		return handleError(c, err)
	}

	workOrders.Status = entity.WorkOrderStatusPending
	workOrders.DeviceID = int64(device.ID)
	err = h.workOrderUsecase.Update(&workOrders)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}

func (h *WorkOrderHandler) DeleteTestWorkOrder(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	patientID, err := strconv.ParseInt(c.Param("patient_id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	err = h.workOrderUsecase.DeleteTest(id, patientID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.WorkOrder{
		ID: id,
	})
}

func (h *WorkOrderHandler) CancelOrder(c echo.Context) error {
	var req entity.WorkOrderCancelRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	querySpeciment := h.db.Model(entity.Specimen{}).Where("order_id = ?", req.WorkOrderID).Select("id")
	err := h.db.Model(entity.ObservationRequest{}).Where("specimen_id in (?)", querySpeciment).
		Update("result_status", constant.ResultStatusDelete).Error
	if err != nil {
		return handleError(c, err)
	}

	workOrders, err := h.workOrderUsecase.FindOneByID(req.WorkOrderID)
	if err != nil {
		return handleError(c, err)
	}

	device := entity.Device{}
	tx := h.db.First(&device, workOrders.DeviceID)
	if tx.Error != nil {
		return handleError(c, fmt.Errorf("error finding device: %w", tx.Error))
	}

	workOrders.Status = entity.WorkOrderCancelled
	err = ba400.SendToBA400(c.Request().Context(), workOrders.Patient, device)
	if err != nil {
		return handleError(c, err)
	}

	err = h.workOrderUsecase.Update(&workOrders)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
