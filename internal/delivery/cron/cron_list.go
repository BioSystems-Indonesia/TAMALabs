package cron

import "context"

type CronJob struct {
	Name        string
	Description string
	Schedule    string
	Execute     func(ctx context.Context) error
}

func GetAllJob(h *CronHandler) []CronJob {
	return []CronJob{
		{
			Name:        "daily_cleanup",
			Description: "Daily cleanup task to prevent memory leaks and reset daily sequences",
			Schedule:    "0 0 1 * * *", // Run at 1 AM every day
			Execute:     h.DailyCleanup,
		},
		// {
		// 	Name:        "sync_all_request",
		// 	Description: "Synchronizes all requests from external systems",
		// 	Schedule:    "0 */5 * * * *",
		// 	Execute:     h.SyncAllRequest,
		// },
		// {
		// 	Name:        "sync_all_result",
		// 	Description: "Synchronizes all results from external systems",
		// 	Schedule:    "0 */5 * * * *",
		// 	Execute:     h.SyncAllResult,
		// },
		{
			Name:        "sync_all_request_simrs",
			Description: "Synchronizes all requests to SIMRS",
			Schedule:    "*/8 * * * * *", // Run every 8 seconds
			Execute:     h.SyncAllRequestSIMRS,
		},
		{
			Name:        "sync_all_result_simrs",
			Description: "Synchronizes all results to SIMRS",
			Schedule:    "*/10 * * * * *", // Run every 10 seconds
			Execute:     h.SyncAllResultSIMRS,
		},
	}
}
