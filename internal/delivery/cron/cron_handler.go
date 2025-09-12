package cron

import (
	"context"

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
