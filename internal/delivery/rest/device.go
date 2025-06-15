package rest

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	deviceuc "github.com/oibacidem/lims-hl-seven/internal/usecase/device"
	"github.com/oibacidem/lims-hl-seven/pkg/panics"
)

const defaultTimeout = 5 * time.Second

type DeviceHandler struct {
	usecase *deviceuc.DeviceUseCase
}

func NewDeviceHandler(usecase *deviceuc.DeviceUseCase) *DeviceHandler {
	return &DeviceHandler{usecase}
}

func (h *DeviceHandler) RegisterRoute(device *echo.Group) {
	device.GET("", h.ListDevices)
	device.GET("/connection", h.GetDeviceConnection)
	device.POST("", h.CreateDevice)
	device.GET("/:id", h.GetDevice)
	device.PUT("/:id", h.UpdateDevice)
	device.DELETE("/:id", h.DeleteDevice)
}

func (h *DeviceHandler) ListDevices(c echo.Context) error {
	var req entity.GetManyRequestDevice
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.usecase.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, result)
}

func (h *DeviceHandler) CreateDevice(c echo.Context) error {
	var req entity.Device
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.usecase.Create(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusCreated, req)
}

func (h *DeviceHandler) GetDevice(c echo.Context) error {
	ctx := c.Request().Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	device, err := h.usecase.FindOneByID(ctx, id)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) UpdateDevice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	req := entity.Device{ID: int(id)}
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	if err := h.usecase.Update(&req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, req)
}

func (h *DeviceHandler) DeleteDevice(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return handleError(c, entity.ErrBadRequest.WithInternal(err))
	}

	if err := h.usecase.Delete(int(id)); err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, entity.Device{
		ID: int(id),
	})
}

func (h *DeviceHandler) GetDeviceConnection(c echo.Context) error {
	w, f, err := createSSEWriter(c)
	if err != nil {
		return handleErrorSSE(c, w, err)
	}
	f.Flush()

	var req entity.DeviceConnectionRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	timeout := time.Duration(req.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = defaultTimeout
	}

	deviceConnection := make(chan entity.DeviceConnectionResponse, len(req.DeviceIDs))
	for _, id := range req.DeviceIDs {
		idInput := id
		go panics.CapturePanic(c.Request().Context(), func() {
			deviceConnection <- h.handleConnection(c.Request().Context(), timeout, idInput)
		})
	}

	// Loop exactly the number of times as there are devices.
	for i := 0; i < len(req.DeviceIDs); i++ {
		select {
		case res := <-deviceConnection:
			message := entity.NewDeviceConnectionMessage(res.DeviceID, res.Message, res.Status)
			if _, err := w.Write([]byte(message)); err != nil {
				return err
			}
			f.Flush()
		case <-c.Request().Context().Done():
			return nil
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *DeviceHandler) handleConnection(
	ctx context.Context,
	timeout time.Duration,
	idInput int,
) entity.DeviceConnectionResponse {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errChan := make(chan error)
	go panics.CapturePanic(ctx, func() {
		errCon := h.usecase.GetDeviceConnection(ctx, idInput)
		errChan <- errCon
	})

	select {
	case <-ctx.Done():
		return entity.DeviceConnectionResponse{
			DeviceID: idInput,
			Message:  "Connection timeout",
			Status:   entity.DeviceConnectionStatusDisconnected,
		}
	case errCon := <-errChan:
		if errCon != nil {
			if errors.Is(errCon, entity.ErrDeviceTypeNotSupport) {
				return entity.DeviceConnectionResponse{
					DeviceID: idInput,
					Message:  "Device type not supported yet",
					Status:   entity.DeviceConnectionStatusNotSupported,
				}
			}

			return entity.DeviceConnectionResponse{
				DeviceID: idInput,
				Message:  errCon.Error(),
				Status:   entity.DeviceConnectionStatusDisconnected,
			}
		}

		return entity.DeviceConnectionResponse{
			DeviceID: idInput,
			Message:  "Connection successful",
			Status:   entity.DeviceConnectionStatusConnected,
		}
	}
}
