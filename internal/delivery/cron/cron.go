package cron

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BioSystems-Indonesia/TAMALabs/pkg/panics"
	"github.com/robfig/cron/v3"
)

// CronManager manages all cron jobs
type CronManager struct {
	cron *cron.Cron
	h    *CronHandler
}

// NewCronManager creates a new cron manager
func NewCronManager(h *CronHandler) *CronManager {
	return &CronManager{
		cron: cron.New(cron.WithSeconds()),
		h:    h,
	}
}

func (cm *CronManager) Start() error {
	if err := cm.RegisterJob(); err != nil {
		return fmt.Errorf("failed to register cron job: %w", err)
	}
	cm.cron.Start()

	return nil
}

// RegisterJob registers a new cron job
func (cm *CronManager) RegisterJob() error {
	jobs := GetAllJob(cm.h)
	for _, job := range jobs {
		if _, err := cm.cron.AddFunc(job.Schedule, wrapFunction(job)); err != nil {
			return fmt.Errorf("failed to add cron job '%s': %w", job.Name, err)
		}
		slog.Info("Registered cron job", "job", job.Name, "schedule", job.Schedule)
	}

	return nil
}

func (cm *CronManager) Stop() error {
	cm.cron.Stop()
	return nil
}

func wrapFunction(job CronJob) func() {
	return func() {
		ctx := context.Background()
		defer panics.RecoverPanic(ctx)

		err := job.Execute(ctx)
		if err != nil {
			slog.ErrorContext(
				ctx,
				"error sync all request",
				"error", err,
				"job", job.Name,
				"schedule", job.Schedule,
				"description", job.Description,
			)
			return
		}

		slog.InfoContext(ctx, "cron job executed", "job", job.Name, "schedule", job.Schedule, "description", job.Description)
	}
}
