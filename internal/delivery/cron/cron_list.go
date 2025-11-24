package cron

import (
	"context"
	"log/slog"
)

type CronJob struct {
	Name        string
	Description string
	Schedule    string
	Execute     func(ctx context.Context) error
}

func GetAllJob(h *CronHandler) []CronJob {
	// Get backup schedule from config
	backupSchedule := h.getBackupSchedule()

	ctx := context.Background()

	// Base jobs that are always active
	jobs := []CronJob{
		{
			Name:        "daily_cleanup",
			Description: "Daily cleanup task to prevent memory leaks and reset daily sequences",
			Schedule:    "0 0 1 * * *", // Run at 1 AM every day
			Execute:     h.DailyCleanup,
		},
		{
			Name:        "license_heartbeat",
			Description: "Send periodic heartbeat to license server to check license status",
			Schedule:    "0 0 */12 * * *", // Run every 12 hours (at minute 0, second 0)
			Execute:     h.LicenseHeartbeat,
		},
		{
			Name:        "database_backup",
			Description: "Backup database with configurable schedule",
			Schedule:    backupSchedule,
			Execute:     h.BackupDB,
		},
	}

	// Conditionally add Database Sharing jobs only if enabled and selected
	simgosEnabled, err := h.configUC.Get(ctx, "SimgosIntegrationEnabled")
	if err == nil && simgosEnabled == "true" {
		selectedSimrs, err := h.configUC.Get(ctx, "SelectedSimrs")
		if err == nil && selectedSimrs == "simgos" {
			slog.Info("Adding Database Sharing cron jobs to scheduler")
			jobs = append(jobs,
				CronJob{
					Name:        "sync_all_request_SIMRS",
					Description: "Synchronizes all requests from Database Sharing",
					Schedule:    "*/8 * * * * *", // Run every 8 seconds
					Execute:     h.SyncAllRequestSIMGOS,
				},
				CronJob{
					Name:        "sync_all_result_SIMRS",
					Description: "Synchronizes all results to Database Sharing",
					Schedule:    "*/10 * * * * *", // Run every 10 seconds
					Execute:     h.SyncAllResultSIMGOS,
				},
			)
		} else {
			slog.Info("Database Sharing cron jobs not added: not selected", "selected", selectedSimrs)
		}
	} else {
		slog.Info("Database Sharing cron jobs not added: Integration not enabled", "enabled", simgosEnabled)
	}

	// TODO: Add SIMRS jobs conditionally when needed
	// simrsEnabled, err := h.configUC.Get(ctx, "SimrsIntegrationEnabled")
	// if err == nil && simrsEnabled == "true" { ... }

	return jobs
}
