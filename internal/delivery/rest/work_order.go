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
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
	"github.com/oibacidem/lims-hl-seven/pkg/panics"
	"gorm.io/gorm"
)

type WorkOrderHandler struct {
	cfg              *config.Schema
	workOrderUsecase *workOrderuc.WorkOrderUseCase
	patientUsecase   *patientuc.PatientUseCase
	db               *gorm.DB
}

func NewWorkOrderHandler(
	cfg *config.Schema, workOrderUsecase *workOrderuc.WorkOrderUseCase, db *gorm.DB,
	patientUsecase *patientuc.PatientUseCase,
) *WorkOrderHandler {
	return &WorkOrderHandler{cfg: cfg, workOrderUsecase: workOrderUsecase, db: db, patientUsecase: patientUsecase}
}

func (h *WorkOrderHandler) FindWorkOrders(c echo.Context) error {
	var req entity.WorkOrderGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	resp, err := h.workOrderUsecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, resp)
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
	var req entity.WorkOrderCreateRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	workOrder, err := h.workOrderUsecase.Create(&req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, workOrder)
}

func (h *WorkOrderHandler) EditWorkOrder(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	var req entity.WorkOrderCreateRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	workOrder, err := h.workOrderUsecase.Edit(int(id), &req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, workOrder)
}

func (h *WorkOrderHandler) RunWorkOrder(c echo.Context) error {
	ctx := c.Request().Context()
	writer, flusher, err := createSSEWriter(c)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}
	defer flusher.Flush()

	var req entity.WorkOrderRunRequest
	if err = bindAndValidate(c, &req); err != nil {
		return handleErrorSSE(c, writer, err)
	}

	querySpeciment := h.db.Model(entity.Specimen{}).Where("order_id in (?)", req.WorkOrderIDs).Select("id")
	err = h.db.Model(entity.ObservationRequest{}).Where("specimen_id in (?)", querySpeciment).
		Update("result_status", constant.ResultStatusSpecimenPending).Error
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	patients, err := h.patientUsecase.FindManyByWorkOrderID(ctx, req.WorkOrderIDs)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	device := entity.Device{}
	tx := h.db.First(&device, req.DeviceID)
	if tx.Error != nil {
		return handleErrorSSE(c, writer, fmt.Errorf("error finding device: %w", tx.Error))
	}

	sendDone := make(chan error, 1)
	go panics.CapturePanic(ctx, func() {
		err := ba400.SendToBA400(ctx, &entity.SendPayloadRequest{
			Patients: patients,
			Device:   device,
			Urgent:   req.Urgent,

			Writer:  writer,
			Flusher: flusher,
		})
		if err != nil {
			sendDone <- fmt.Errorf("error sending to ba400: %w", err)
			return
		}

		sendDone <- nil
	})

	select {
	case err := <-sendDone:
		if err != nil {
			return handleErrorSSE(c, writer, err)
		}

		for _, workOrderID := range req.WorkOrderIDs {
			workOrder, err := h.workOrderUsecase.FindOneByID(workOrderID)
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}

			workOrder.Status = entity.WorkOrderStatusPending
			err = h.workOrderUsecase.UpsertDevice(workOrderID, int64(device.ID))
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}

			err = h.workOrderUsecase.Update(&workOrder)
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}
		}

		_, err = writer.Write([]byte(entity.NewWorkOrderStreamingResponse(100, entity.WorkOrderStreamingResponseStatusDone)))
		if err != nil {
			return handleErrorSSE(c, writer, err)
		}

		return c.NoContent(http.StatusOK)
	case <-ctx.Done():
		for _, workOrderID := range req.WorkOrderIDs {
			workOrder, err := h.workOrderUsecase.FindOneByID(workOrderID)
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}

			workOrder.Status = entity.WorkOrderStatusIncompleteSend
			err = h.workOrderUsecase.Update(&workOrder)
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}
		}

		return handleErrorSSE(c, writer, ctx.Err())
	}
}

func (h *WorkOrderHandler) CancelOrder(c echo.Context) error {
	var req entity.WorkOrderCancelRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	// ctx := c.Request().Context()
	writer, flusher, err := createSSEWriter(c)
	if err != nil {
		return handleError(c, err)
	}

	querySpeciment := h.db.Model(entity.Specimen{}).Where("order_id = ?", req.WorkOrderID).Select("id")
	err = h.db.Model(entity.ObservationRequest{}).Where("specimen_id in (?)", querySpeciment).
		Update("result_status", constant.ResultStatusDelete).Error
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

	workOrders.Status = entity.WorkOrderCancelled
	err = ba400.SendToBA400(c.Request().Context(), &entity.SendPayloadRequest{
		Patients: []entity.Patient{workOrders.Patient},
		Device:   device,
		Urgent:   false,
		Writer:   writer,
		Flusher:  flusher,
	})
	if err != nil {
		return handleError(c, err)
	}

	err = h.workOrderUsecase.Update(&workOrders)
	if err != nil {
		return handleError(c, err)
	}

	flusher.Flush()

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
	})
}
