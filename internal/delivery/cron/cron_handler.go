package cron

import (
	"context"
	"log/slog"
	"runtime"

	khanzauc "github.com/oibacidem/lims-hl-seven/internal/usecase/external/khanza"
)

type CronHandler struct {
	khanzaUC *khanzauc.Usecase
}

func NewCronHandler(khanzaUC *khanzauc.Usecase) *CronHandler {
	return &CronHandler{
		khanzaUC: khanzaUC,
	}
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
