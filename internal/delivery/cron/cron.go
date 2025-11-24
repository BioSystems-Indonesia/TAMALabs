package cron

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/BioSystems-Indonesia/TAMALabs/pkg/panics"
	"github.com/robfig/cron/v3"
)

// CronJobInfo represents information about a cron job
type CronJobInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	Active      bool   `json:"active"`
}

// ConfigChecker defines interface to check configuration
type ConfigChecker interface {
	Get(ctx context.Context, key string) (string, error)
}

// CronManager manages all cron jobs
type CronManager struct {
	cron          *cron.Cron
	h             *CronHandler
	jobs          []CronJob
	entries       map[string]cron.EntryID // Track cron entries by job name
	disabledJobs  map[string]bool         // Track disabled jobs
	configChecker ConfigChecker           // Check configuration for conditional jobs
}

// NewCronManager creates a new cron manager
func NewCronManager(h *CronHandler, configChecker ConfigChecker) *CronManager {
	return &CronManager{
		cron:          cron.New(cron.WithSeconds()),
		h:             h,
		entries:       make(map[string]cron.EntryID),
		disabledJobs:  make(map[string]bool),
		configChecker: configChecker,
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
	cm.jobs = jobs // Store jobs for later access
	for _, job := range jobs {
		entryID, err := cm.cron.AddFunc(job.Schedule, wrapFunction(job))
		if err != nil {
			return fmt.Errorf("failed to add cron job '%s': %w", job.Name, err)
		}
		cm.entries[job.Name] = entryID
		slog.Info("Registered cron job", "job", job.Name, "schedule", job.Schedule)
	}

	return nil
}

func (cm *CronManager) Stop() error {
	cm.cron.Stop()
	return nil
}

// ReloadJobs reloads all cron jobs based on current configuration
// This removes all existing jobs and re-registers them
func (cm *CronManager) ReloadJobs() error {
	slog.Info("Reloading cron jobs based on current configuration")

	// Remove all existing job entries
	for jobName, entryID := range cm.entries {
		cm.cron.Remove(entryID)
		slog.Info("Removed cron job", "job", jobName)
	}

	// Clear the maps
	cm.entries = make(map[string]cron.EntryID)
	cm.disabledJobs = make(map[string]bool)

	// Re-register jobs based on current configuration
	if err := cm.RegisterJob(); err != nil {
		return fmt.Errorf("failed to re-register jobs: %w", err)
	}

	slog.Info("Cron jobs reloaded successfully")
	return nil
}

// GetAllJobs returns information about all registered cron jobs
func (cm *CronManager) GetAllJobs() []CronJobInfo {
	var jobInfos []CronJobInfo
	for _, job := range cm.jobs {
		isActive := cm.isJobActive(job.Name)
		jobInfos = append(jobInfos, CronJobInfo{
			Name:        job.Name,
			Description: job.Description,
			Schedule:    job.Schedule,
			Active:      isActive,
		})
	}
	return jobInfos
}

// isJobActive checks if a job is currently active
func (cm *CronManager) isJobActive(jobName string) bool {
	// Check if manually disabled
	if cm.disabledJobs[jobName] {
		return false
	}

	// Check if conditionally disabled based on configuration
	if cm.shouldJobBeInactive(jobName) {
		return false
	}

	// Check if registered in cron
	_, exists := cm.entries[jobName]
	return exists
}

// shouldJobBeInactive checks if a job should be inactive based on configuration
func (cm *CronManager) shouldJobBeInactive(jobName string) bool {
	// Database Sharing jobs should only be active when Database Sharing is selected
	if jobName == "sync_all_request_simrs" || jobName == "sync_all_result_simrs" {
		ctx := context.Background()

		// Check if Database Sharing integration is enabled
		simgosEnabled, err := cm.configChecker.Get(ctx, "SimgosIntegrationEnabled")
		if err != nil {
			slog.Debug("shouldJobBeInactive: Failed to get SimgosIntegrationEnabled", "job", jobName, "error", err)
			return true // Job should be inactive
		}
		if simgosEnabled != "true" {
			slog.Debug("shouldJobBeInactive: Database Sharing integration not enabled", "job", jobName, "enabled", simgosEnabled)
			return true // Job should be inactive
		}

		// Check if Database Sharing is selected as the active SIMRS
		selectedSimrs, err := cm.configChecker.Get(ctx, "SelectedSimrs")
		if err != nil {
			slog.Debug("shouldJobBeInactive: Failed to get SelectedSimrs", "job", jobName, "error", err)
			return true // Job should be inactive
		}
		if selectedSimrs != "simgos" {
			slog.Debug("shouldJobBeInactive: Database Sharing not selected", "job", jobName, "selected", selectedSimrs)
			return true // Job should be inactive
		}

		slog.Debug("shouldJobBeInactive: Database Sharing job should be active", "job", jobName)
		return false // Job should be active
	}

	return false
}

// EnableJob enables a specific cron job
func (cm *CronManager) EnableJob(jobName string) error {
	// Check if job exists
	var targetJob *CronJob
	for i := range cm.jobs {
		if cm.jobs[i].Name == jobName {
			targetJob = &cm.jobs[i]
			break
		}
	}
	if targetJob == nil {
		return fmt.Errorf("job '%s' not found", jobName)
	}

	// Check if already enabled
	if !cm.disabledJobs[jobName] {
		return nil // Already enabled
	}

	// Add the job back to cron
	entryID, err := cm.cron.AddFunc(targetJob.Schedule, wrapFunction(*targetJob))
	if err != nil {
		return fmt.Errorf("failed to enable job '%s': %w", jobName, err)
	}

	cm.entries[jobName] = entryID
	delete(cm.disabledJobs, jobName)
	slog.Info("Enabled cron job", "job", jobName)
	return nil
}

// DisableJob disables a specific cron job
func (cm *CronManager) DisableJob(jobName string) error {
	// Prevent disabling critical jobs
	if jobName == "license_heartbeat" {
		return fmt.Errorf("cannot disable critical job 'license_heartbeat'")
	}

	// Check if job exists
	entryID, exists := cm.entries[jobName]
	if !exists {
		return fmt.Errorf("job '%s' not found or not registered", jobName)
	}

	// Remove the job from cron
	cm.cron.Remove(entryID)
	delete(cm.entries, jobName)
	cm.disabledJobs[jobName] = true
	slog.Info("Disabled cron job", "job", jobName)
	return nil
}

// UpdateBackupSchedule updates the database backup job with a new schedule
func (cm *CronManager) UpdateBackupSchedule() error {
	jobName := "database_backup"

	// Find the backup job
	var targetJob *CronJob
	for i := range cm.jobs {
		if cm.jobs[i].Name == jobName {
			targetJob = &cm.jobs[i]
			break
		}
	}
	if targetJob == nil {
		return fmt.Errorf("backup job not found")
	}

	// Get new schedule
	newSchedule := cm.h.getBackupSchedule()

	// If schedule hasn't changed, no need to update
	if targetJob.Schedule == newSchedule {
		return nil
	}

	// Remove old job if it exists
	if entryID, exists := cm.entries[jobName]; exists {
		cm.cron.Remove(entryID)
		delete(cm.entries, jobName)
	}

	// Update the job schedule
	targetJob.Schedule = newSchedule

	// Re-register the job with new schedule
	entryID, err := cm.cron.AddFunc(newSchedule, wrapFunction(*targetJob))
	if err != nil {
		return fmt.Errorf("failed to update backup schedule: %w", err)
	}

	cm.entries[jobName] = entryID
	slog.Info("Updated backup schedule", "job", jobName, "schedule", newSchedule)
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
				"cron job failed",
				"error", err,
				"job", job.Name,
				"schedule", job.Schedule,
				"description", job.Description,
			)
			return
		}
		// Note: Success logging is handled by individual job handlers
		// to distinguish between actual execution and skipped jobs
	}
}
