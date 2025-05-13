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
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
	"github.com/oibacidem/lims-hl-seven/pkg/panics"
	"gorm.io/gorm"
)

type WorkOrderHandler struct {
	cfg              *config.Schema
	workOrderUsecase *workOrderuc.WorkOrderUseCase
	deviceUsecase    *deviceuc.DeviceUseCase
	patientUsecase   *patientuc.PatientUseCase
	db               *gorm.DB
}

func NewWorkOrderHandler(
	cfg *config.Schema, workOrderUsecase *workOrderuc.WorkOrderUseCase, db *gorm.DB,
	patientUsecase *patientuc.PatientUseCase,
	deviceUsecase *deviceuc.DeviceUseCase,
) *WorkOrderHandler {
	return &WorkOrderHandler{cfg: cfg, workOrderUsecase: workOrderUsecase, db: db, patientUsecase: patientUsecase, deviceUsecase: deviceUsecase}
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
	return h.runWorkOrder(c, constant.WorkOrderRunActionRun)
}

func (h *WorkOrderHandler) CancelOrder(c echo.Context) error {
	return h.runWorkOrder(c, constant.WorkOrderRunActionCancel)
}

func (h *WorkOrderHandler) runWorkOrder(c echo.Context, action constant.WorkOrderRunAction) error {
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

	// TODO: Change this to strategy pattern
	switch action {
	case constant.WorkOrderRunActionRun:
		querySpeciment := h.db.Model(entity.Specimen{}).Where("order_id in (?)", req.WorkOrderIDs).Select("id")
		err = h.db.Model(entity.ObservationRequest{}).Where("specimen_id in (?)", querySpeciment).
			Update("result_status", constant.ResultStatusSpecimenPending).Error
		if err != nil {
			return handleErrorSSE(c, writer, err)
		}
	case constant.WorkOrderRunActionCancel:
		querySpeciment := h.db.Model(entity.Specimen{}).Where("order_id in (?)", req.WorkOrderIDs).Select("id")
		err = h.db.Model(entity.ObservationRequest{}).Where("specimen_id in (?)", querySpeciment).
			Update("result_status", constant.ResultStatusDelete).Error
		if err != nil {
			return handleErrorSSE(c, writer, err)
		}
	}

	patients, err := h.patientUsecase.FindManyByWorkOrderID(ctx, req.WorkOrderIDs)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	device, err := h.deviceUsecase.FindByID(ctx, req.DeviceID)
	if err != nil {
		return handleErrorSSE(c, writer, fmt.Errorf("error finding device: %w", err))
	}

	sendDone := make(chan error, 1)
	go panics.CapturePanic(ctx, func() {
		// TODO: Change this to strategy pattern using device type
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
			errCancel := h.runWorkOrder(c, constant.WorkOrderRunActionCancel)
			if errCancel != nil {
				return handleErrorSSE(c, writer, errCancel, map[string]interface{}{
					"error_cancel": errCancel.Error(),
				})
			}

			return handleErrorSSE(c, writer, err)
		}

		for _, workOrderID := range req.WorkOrderIDs {
			// TODO: Change this to strategy pattern
			switch action {
			case constant.WorkOrderRunActionRun:
				err := h.updateWorkOrderStatus(c, writer, workOrderID, device, entity.WorkOrderStatusPending)
				if err != nil {
					return handleErrorSSE(c, writer, err)
				}
			case constant.WorkOrderRunActionCancel:
				err := h.updateWorkOrderStatus(c, writer, workOrderID, device, entity.WorkOrderCancelled)
				if err != nil {
					return handleErrorSSE(c, writer, err)
				}
			}
		}

		_, err = writer.Write([]byte(entity.NewWorkOrderStreamingResponse(100, entity.WorkOrderStreamingResponseStatusDone)))
		if err != nil {
			return handleErrorSSE(c, writer, err)
		}

		return c.NoContent(http.StatusOK)
	case <-ctx.Done():
		for _, workOrderID := range req.WorkOrderIDs {
			err := h.updateWorkOrderStatus(c, writer, workOrderID, device, entity.WorkOrderStatusIncompleteSend)
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}
		}

		errCancel := h.runWorkOrder(c, constant.WorkOrderRunActionCancel)
		if errCancel != nil {
			return handleErrorSSE(c, writer, errCancel, map[string]interface{}{
				"error_cancel": errCancel.Error(),
			})
		}

		return handleErrorSSE(c, writer, ctx.Err())
	}
}

func (h *WorkOrderHandler) updateWorkOrderStatus(c echo.Context, writer http.ResponseWriter, workOrderID int64, device entity.Device, status entity.WorkOrderStatus) error {
	workOrder, err := h.workOrderUsecase.FindOneByID(workOrderID)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	workOrder.Status = status
	err = h.workOrderUsecase.UpsertDevice(workOrderID, int64(device.ID))
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	err = h.workOrderUsecase.Update(&workOrder)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}
	return nil
}
