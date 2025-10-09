package cron

import (
	"context"
	"log/slog"
	"runtime"

	khanzauc "github.com/oibacidem/lims-hl-seven/internal/usecase/external/khanza"
)

// SIMRSUsecase interface to avoid import cycle
type SIMRSUsecase interface {
	SyncAllRequest(ctx context.Context) error
	SyncAllResult(ctx context.Context, workOrderIDs []int64) error
}

type CronHandler struct {
	khanzaUC *khanzauc.Usecase
	simrsUC  SIMRSUsecase
}

func NewCronHandler(khanzaUC *khanzauc.Usecase, simrsUC SIMRSUsecase) *CronHandler {
	return &CronHandler{
		khanzaUC: khanzaUC,
		simrsUC:  simrsUC,
	}
}

func (c *CronHandler) SyncAllRequestSIMRS(ctx context.Context) error {
	slog.Info("Starting SIMRS sync all request cron job")

	err := c.simrsUC.SyncAllRequest(ctx)
	if err != nil {
		slog.Error("Failed to sync all requests to SIMRS", "error", err)
		return err
	}

	slog.Info("Successfully completed SIMRS sync all request cron job")
	return nil
}
func (c *CronHandler) SyncAllResultSIMRS(ctx context.Context) error {
	slog.Info("Starting SIMRS sync all result cron job")

	err := c.simrsUC.SyncAllResult(ctx, []int64{})
	if err != nil {
		slog.Error("Failed to sync all results to SIMRS", "error", err)
		return err
	}

	slog.Info("Successfully completed SIMRS sync all result cron job")
	return nil
}

func (c *CronHandler) SyncAllRequest(ctx context.Context) error {
	err := c.khanzaUC.SyncAllRequest(ctx)
	if err != nil {
		return err
	}

	return nil

}

func (c *CronHandler) SyncAllResult(ctx context.Context) error {
	err := c.khanzaUC.SyncAllResult(ctx, []int64{})
	if err != nil {
		return err
	}

	return nil
}

func (c *CronHandler) DailyCleanup(ctx context.Context) error {
	slog.Info("Starting daily cleanup task")

	runtime.GC()

	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	slog.Info("Daily cleanup completed",
		"goroutines", runtime.NumGoroutine(),
		"memory_alloc_mb", m.Alloc/1024/1024,
		"total_alloc_mb", m.TotalAlloc/1024/1024,
		"num_gc", m.NumGC,
	)

	return nil
}
