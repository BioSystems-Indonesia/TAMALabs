package rest

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
	observation_requestuc "github.com/oibacidem/lims-hl-seven/internal/usecase/observation_request"
	patientuc "github.com/oibacidem/lims-hl-seven/internal/usecase/patient"
	specimenuc "github.com/oibacidem/lims-hl-seven/internal/usecase/specimen"
	workOrderuc "github.com/oibacidem/lims-hl-seven/internal/usecase/work_order"
	"gorm.io/gorm"
)

type WorkOrderHandler struct {
	cfg                       *config.Schema
	workOrderUsecase          *workOrderuc.WorkOrderUseCase
	deviceUsecase             *deviceuc.DeviceUseCase
	patientUsecase            *patientuc.PatientUseCase
	specimenUsecase           *specimenuc.SpecimenUseCase
	observationRequestUsecase *observation_requestuc.ObservationRequestUseCase
}

func NewWorkOrderHandler(
	cfg *config.Schema, workOrderUsecase *workOrderuc.WorkOrderUseCase, db *gorm.DB,
	patientUsecase *patientuc.PatientUseCase,
	deviceUsecase *deviceuc.DeviceUseCase,
	specimenUsecase *specimenuc.SpecimenUseCase,
	observationRequestUsecase *observation_requestuc.ObservationRequestUseCase,
) *WorkOrderHandler {
	return &WorkOrderHandler{
		cfg: cfg, workOrderUsecase: workOrderUsecase,
		patientUsecase: patientUsecase, deviceUsecase: deviceUsecase,
		specimenUsecase: specimenUsecase, observationRequestUsecase: observationRequestUsecase,
	}
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

const deadline = 15 * time.Second

func (h *WorkOrderHandler) runWorkOrder(c echo.Context, action constant.WorkOrderRunAction) error {
	ctx := c.Request().Context()
	writer, flusher, err := createSSEWriter(c)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	var req entity.WorkOrderRunRequest
	if err = bindAndValidate(c, &req); err != nil {
		return handleErrorSSE(c, writer, err)
	}

	_, err = writer.Write([]byte(entity.NewWorkOrderStreamingResponse(50, entity.WorkOrderStreamingResponseStatusDone)))
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}
	flusher.Flush()

	progressWriter, err := h.workOrderUsecase.RunWorkOrderAsync(ctx, &req, action)
	if err != nil {
		return handleErrorSSE(c, writer, err)
	}

	for {
		select {
		case msg, ok := <-progressWriter:
			if !ok {
				return handleErrorSSE(c, writer, errors.New("process writer closed"))
			}

			slog.Info("message received", slog.Attr{
				Key:   "percentage",
				Value: slog.Float64Value(msg.Percentage),
			}, slog.Attr{
				Key:   "status",
				Value: slog.StringValue(string(msg.Status)),
			}, slog.Attr{
				Key:   "isDone",
				Value: slog.BoolValue(msg.IsDone),
			})

			if msg.Error != nil {
				slog.Info("message error", slog.Attr{
					Key:   "error",
					Value: slog.StringValue(msg.Error.Error()),
				})
				return handleErrorSSE(c, writer, msg.Error)
			}

			_, err = writer.Write([]byte(entity.NewWorkOrderStreamingResponse(msg.Percentage, msg.Status)))
			if err != nil {
				return handleErrorSSE(c, writer, err)
			}
			flusher.Flush()

			if msg.IsDone {
				return c.NoContent(http.StatusOK)
			}
		case <-time.After(deadline):
			err := errors.New("send timeout, please check your connection")
			return handleErrorSSE(c, writer, err)
		}
	}
}

func (h *WorkOrderHandler) GetWorkOrderBarcode(c echo.Context) error {
	barcodes, err := h.workOrderUsecase.FindAllBarcodes(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	resp := make([]entity.Table, len(barcodes))
	for i, barcode := range barcodes {
		resp[i] = entity.Table{
			ID:   barcode,
			Name: barcode,
		}
	}

	return successMany(c, resp)
}
