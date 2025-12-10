package rest

import (
	"net/http"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/cron"
	"github.com/labstack/echo/v4"
)

// CronHandler handles cron job related HTTP requests
type CronHandler struct {
	cronManager CronManagerInterface
}

// CronManagerInterface defines the interface for cron manager
type CronManagerInterface interface {
	GetAllJobs() []cron.CronJobInfo
	EnableJob(jobName string) error
	DisableJob(jobName string) error
	UpdateBackupSchedule() error
	ReloadJobs() error
}

// NewCronHandler creates a new cron handler
func NewCronHandler(cronManager CronManagerInterface) *CronHandler {
	return &CronHandler{
		cronManager: cronManager,
	}
}

// GetAllJobs returns all cron jobs
func (h *CronHandler) GetAllJobs(c echo.Context) error {
	jobs := h.cronManager.GetAllJobs()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": jobs,
	})
}

// EnableJob enables a specific cron job
func (h *CronHandler) EnableJob(c echo.Context) error {
	jobName := c.Param("name")
	if jobName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "job name is required",
		})
	}

	if err := h.cronManager.EnableJob(jobName); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "job enabled successfully",
	})
}

// DisableJob disables a specific cron job
func (h *CronHandler) DisableJob(c echo.Context) error {
	jobName := c.Param("name")
	if jobName == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "job name is required",
		})
	}

	// Prevent disabling critical jobs
	if jobName == "license_heartbeat" {
		return c.JSON(http.StatusForbidden, map[string]interface{}{
			"error": "license_heartbeat is a critical job and cannot be disabled",
		})
	}

	if err := h.cronManager.DisableJob(jobName); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "job disabled successfully",
	})
}

// UpdateBackupSchedule updates the backup job schedule
func (h *CronHandler) UpdateBackupSchedule(c echo.Context) error {
	if err := h.cronManager.UpdateBackupSchedule(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "backup schedule updated successfully",
	})
}

// ReloadJobs reloads all cron jobs based on current configuration
func (h *CronHandler) ReloadJobs(c echo.Context) error {
	if err := h.cronManager.ReloadJobs(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "cron jobs reloaded successfully",
	})
}
