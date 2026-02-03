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
		// If not found, return default value for known keys
		if err.Error() == "record not found" || err.Error() == "error finding Config: record not found" {
			defaultValue := getDefaultConfigValue(key)
			result = entity.Config{
				ID:    key,
				Value: defaultValue,
			}
			return c.JSON(http.StatusOK, result)
		}
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, result)
}

// getDefaultConfigValue returns default value for config keys
func getDefaultConfigValue(key string) string {
	defaults := map[string]string{
		"NuhaIntegrationEnabled":   "false",
		"NuhaBaseURL":              "https://api.nuha-simrs.example.com",
		"NuhaSessionID":            "",
		"SimrsIntegrationEnabled":  "false",
		"SimgosIntegrationEnabled": "false",
		"KhanzaIntegrationEnabled": "false",
		"BackupScheduleType":       "daily",
		"BackupInterval":           "6",
		"BackupTime":               "02:00",
	}

	if val, ok := defaults[key]; ok {
		return val
	}
	return ""
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
