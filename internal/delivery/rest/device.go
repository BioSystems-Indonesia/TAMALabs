package rest

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"gorm.io/gorm"
)

type DeviceHandler struct {
	db *gorm.DB
}

func (h *DeviceHandler) ListDevices(c echo.Context) error {
	var devices []entity.Device
	if err := h.db.Find(&devices).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	c.Response().Header().Set(entity.HeaderXTotalCount, strconv.Itoa(len(devices)))
	return c.JSON(http.StatusOK, devices)
}

func (h *DeviceHandler) CreateDevice(c echo.Context) error {
	device := new(entity.Device)
	if err := c.Bind(device); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.db.Create(device).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, device)
}

func (h *DeviceHandler) GetDevice(c echo.Context) error {
	var device entity.Device
	if err := h.db.First(&device, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	return c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) UpdateDevice(c echo.Context) error {
	device := new(entity.Device)
	if err := c.Bind(device); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := h.db.Save(device).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, device)
}

func (h *DeviceHandler) DeleteDevice(c echo.Context) error {
	if err := h.db.Delete(&entity.Device{}, c.Param("id")).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func NewDeviceHandler(db *gorm.DB) *DeviceHandler {
	return &DeviceHandler{db: db}
}
