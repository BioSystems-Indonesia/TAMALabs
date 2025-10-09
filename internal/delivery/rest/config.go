package rest

import (
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	configuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/config"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/labstack/echo/v4"
)

type ConfigHandler struct {
	cfg      *config.Schema
	configUC *configuc.ConfigUseCase
}

func NewConfigHandler(cfg *config.Schema, resultUsecase *configuc.ConfigUseCase) *ConfigHandler {
	return &ConfigHandler{
		cfg:      cfg,
		configUC: resultUsecase,
	}
}

func (h *ConfigHandler) ListConfig(c echo.Context) error {
	var req entity.ConfigGetManyRequest
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.configUC.FindAll(
		c.Request().Context(),
		&req,
	)
	if err != nil {
		return handleError(c, err)
	}

	return successPaginationResponse(c, result)
}

func (h *ConfigHandler) GetConfig(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return handleError(c, entity.ErrBadRequest)
	}

	result, err := h.configUC.FindOneByID(c.Request().Context(), key)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ConfigHandler) EditConfig(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return handleError(c, entity.ErrBadRequest)
	}

	var req entity.Config
	if err := bindAndValidate(c, &req); err != nil {
		return handleError(c, err)
	}

	result, err := h.configUC.Edit(
		c.Request().Context(),
		key,
		req.Value,
	)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}
