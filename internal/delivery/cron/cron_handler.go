package cron

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	// use modernc.org/sqlite (imported in app package) which does not require cgo

	khanzauc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/khanza"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
)

// SIMRSUsecase interface to avoid import cycle
type SIMRSUsecase interface {
	SyncAllRequest(ctx context.Context) error
	SyncAllResult(ctx context.Context, workOrderIDs []int64) error
}

// SIMGOSUsecase interface to avoid import cycle
type SIMGOSUsecase interface {
	SyncAllRequest(ctx context.Context) error
	SyncAllResult(ctx context.Context, workOrderIDs []int64) error
}

// License heartbeat request/response structs
type HeartbeatRequest struct {
	MachineID   string `json:"machine_id"`
	LicenseCode string `json:"license_code"`
}

type HeartbeatResponse struct {
	Error string `json:"error"`
}

// ConfigUsecase interface to get configuration values
type ConfigUsecase interface {
	Get(ctx context.Context, key string) (string, error)
}

type CronHandler struct {
	khanzaUC         *khanzauc.Usecase
	simrsUC          SIMRSUsecase
	simgosUC         SIMGOSUsecase
	configUC         ConfigUsecase
	licenseServerURL string
	machineID        string
}

// licenseDirPaths returns the absolute license directory and commonly used file paths
// Uses AppData\Local instead of ProgramData to avoid admin privilege requirements
func licenseDirPaths() (licenseDir, licensePath, pubKeyPath, revokedPath, expiredPath string) {
	// Prefer LOCALAPPDATA for user-writable application data
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		programRoot := filepath.Join(localAppData, "TAMALabs")
		licenseDir = filepath.Join(programRoot, "license")
	} else {
		// Fallback to APPDATA
		appData := os.Getenv("APPDATA")
		if appData != "" {
			programRoot := filepath.Join(appData, "TAMALabs")
			licenseDir = filepath.Join(programRoot, "license")
		} else if runtime.GOOS == "windows" {
			// Last resort: ProgramData (backward compatibility)
			programData := os.Getenv("ProgramData")
			if programData == "" {
				programData = `C:\ProgramData`
			}
			programRoot := filepath.Join(programData, "TAMALabs")
			licenseDir = filepath.Join(programRoot, "license")
		}
	}

	licensePath = filepath.Join(licenseDir, "license.json")
	pubKeyPath = filepath.Join(licenseDir, "server_public.pem")
	revokedPath = filepath.Join(licenseDir, "revoked.json")
	expiredPath = filepath.Join(licenseDir, "expired.json")
	return
}

var API_KEY = "KJKDANCJSANIUWYR6243UJFOISJFJKVOMV72487YEHFHFHSDVOHF9AMDC9AN9SDN98YE98YEHDIU2Y897873YYY68686487WGDUDUAGYTE8QTEYADIUHADUYW8E8BWTNC8N8NAMDOAIMDAUDUWYAD87NYW7Y7CBT87EY8142164B36248732M87MCIFH8NYRWCM8MYCMUOIDOIADOIDOIUR83YR983Y98328N32C83NYC8732NYC8732Y87Y32NCNSAIHJAOJFOIJFOIQFIUIUNCNHCIUHWV8NRYNV8Y989N9198298YOIJOI090103021313JKJDHAHDJAJASHHAH"

func NewCronHandler(khanzaUC *khanzauc.Usecase, simrsUC SIMRSUsecase, simgosUC SIMGOSUsecase, configUC ConfigUsecase) *CronHandler {
	licenseServerURL := os.Getenv("LICENSE_SERVER_URL")
	if licenseServerURL == "" {
		licenseServerURL = "https://tamalabs.biosystems.id"
	}

	machineID, err := util.GenerateMachineID()
	if err != nil {
		slog.Error("Failed to generate machine ID for heartbeat", "error", err)
		machineID = "unknown"
	}

	return &CronHandler{
		khanzaUC:         khanzaUC,
		simrsUC:          simrsUC,
		simgosUC:         simgosUC,
		configUC:         configUC,
		licenseServerURL: licenseServerURL,
		machineID:        machineID,
	}
}

// getBackupSchedule returns the cron schedule for database backup based on config
func (c *CronHandler) getBackupSchedule() string {
	ctx := context.Background()
	defaultSchedule := "0 0 2 * * *" // Default: 2 AM daily

	scheduleType, err := c.configUC.Get(ctx, "BackupScheduleType")
	if err != nil || scheduleType == "" {
		slog.Debug("Using default backup schedule", "schedule", defaultSchedule)
		return defaultSchedule
	}

	if scheduleType == "interval" {
		// Interval-based backup (every N hours)
		intervalStr, err := c.configUC.Get(ctx, "BackupInterval")
		if err != nil || intervalStr == "" {
			return defaultSchedule
		}

		interval, err := strconv.Atoi(intervalStr)
		if err != nil || interval < 1 || interval > 24 {
			slog.Warn("Invalid backup interval, using default", "interval", intervalStr)
			return defaultSchedule
		}

		// Format: "0 0 */N * * *" - every N hours at minute 0
		return fmt.Sprintf("0 0 */%d * * *", interval)
	} else if scheduleType == "daily" {
		// Daily at specific time
		timeStr, err := c.configUC.Get(ctx, "BackupTime")
		if err != nil || timeStr == "" {
			return defaultSchedule
		}

		// Parse time (HH:MM format)
		parts := strings.Split(timeStr, ":")
		if len(parts) != 2 {
			slog.Warn("Invalid backup time format, using default", "time", timeStr)
			return defaultSchedule
		}

		hour, err1 := strconv.Atoi(parts[0])
		minute, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil || hour < 0 || hour > 23 || minute < 0 || minute > 59 {
			slog.Warn("Invalid backup time values, using default", "time", timeStr)
			return defaultSchedule
		}

		// Format: "0 MM HH * * *" - daily at HH:MM
		return fmt.Sprintf("0 %d %d * * *", minute, hour)
	}

	return defaultSchedule
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

func (c *CronHandler) SyncAllRequestSIMGOS(ctx context.Context) error {
	// Check if Database Sharing integration is enabled
	simgosEnabled, err := c.configUC.Get(ctx, "SimgosIntegrationEnabled")
	if err != nil {
		slog.Info("Database Sharing sync skipped: Failed to get integration status", "error", err)
		return nil
	}

	if simgosEnabled != "true" {
		slog.Info("Database Sharing sync skipped: Integration not enabled", "enabled", simgosEnabled)
		return nil
	}

	// Check if Database Sharing is selected
	selectedSimrs, err := c.configUC.Get(ctx, "SelectedSimrs")
	if err != nil {
		slog.Info("Database Sharing sync skipped: Failed to get SelectedSimrs", "error", err)
		return nil
	}

	if selectedSimrs != "simgos" {
		slog.Info("Database Sharing sync skipped: not selected", "selected", selectedSimrs)
		return nil
	}

	slog.Info("Starting Database Sharing sync all request cron job")

	if c.simgosUC == nil {
		slog.Info("Database Sharing sync skipped: Usecase is nil")
		return nil
	}

	err = c.simgosUC.SyncAllRequest(ctx)
	if err != nil {
		slog.Error("Failed to sync all requests from Database Sharing", "error", err)
		return err
	}

	slog.Info("Successfully completed Database Sharing sync all request cron job")
	return nil
}

func (c *CronHandler) SyncAllResultSIMGOS(ctx context.Context) error {
	// Check if Database Sharing integration is enabled
	simgosEnabled, err := c.configUC.Get(ctx, "SimgosIntegrationEnabled")
	if err != nil {
		slog.Info("Database Sharing result sync skipped: Failed to get integration status", "error", err)
		return nil
	}

	if simgosEnabled != "true" {
		slog.Info("Database Sharing result sync skipped: Integration not enabled", "enabled", simgosEnabled)
		return nil
	}

	// Check if Database Sharing is selected
	selectedSimrs, err := c.configUC.Get(ctx, "SelectedSimrs")
	if err != nil {
		slog.Info("Database Sharing result sync skipped: Failed to get SelectedSimrs", "error", err)
		return nil
	}

	if selectedSimrs != "simgos" {
		slog.Info("Database Sharing result sync skipped: not selected", "selected", selectedSimrs)
		return nil
	}

	slog.Info("Starting Database Sharing sync all result cron job")

	if c.simgosUC == nil {
		slog.Info("Database Sharing result sync skipped: Usecase is nil")
		return nil
	}

	err = c.simgosUC.SyncAllResult(ctx, []int64{})
	if err != nil {
		slog.Error("Failed to sync all results to Database Sharing", "error", err)
		return err
	}

	slog.Info("Successfully completed Database Sharing sync all result cron job")
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

// LicenseHeartbeat sends periodic heartbeat to license server
func (c *CronHandler) LicenseHeartbeat(ctx context.Context) error {
	client := &http.Client{Timeout: 30 * time.Second}
	slog.Debug("Performing license heartbeat")

	// Read license file content
	_, licensePath, _, _, _ := licenseDirPaths()
	licenseData, err := os.ReadFile(licensePath)
	if err != nil {
		slog.Warn("License file not found, skipping heartbeat", "error", err)
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(licenseData, &data); err != nil {
		slog.Warn("Failed to decode license file", "error", err)
		return nil
	}

	// Check if license_code exists
	licenseCodeRaw, exists := data["license_code"]
	if !exists {
		slog.Warn("License code not found in license file")
		return nil
	}

	licenseCode, ok := licenseCodeRaw.(string)
	if !ok {
		slog.Warn("License code is not a string")
		return nil
	}

	slog.Debug("Sending heartbeat", "machine_id", c.machineID, "license_code", licenseCode)

	// Prepare heartbeat request
	heartbeatReq := HeartbeatRequest{
		MachineID:   c.machineID,
		LicenseCode: licenseCode,
	}

	// Send heartbeat
	jsonData, err := json.Marshal(heartbeatReq)
	if err != nil {
		slog.Error("Failed to marshal heartbeat request", "error", err)
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/heartbeat", c.licenseServerURL),
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		slog.Warn("License heartbeat failed - server unreachable",
			"error", err,
			"server", c.licenseServerURL,
			"strategy", "continue_with_local_license")
		return nil // Don't fail the cron job
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", API_KEY)

	resp, err := client.Do(request)
	if err != nil {
		slog.Warn("License heartbeat failed - server unreachable",
			"error", err,
			"server", c.licenseServerURL,
			"strategy", "continue_with_local_license")
		return nil
	}

	// Parse response - simple string response
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read heartbeat response", "error", err)
		return nil // Don't fail the cron job
	}
	defer resp.Body.Close()

	// Clean the response string
	bodyStr := strings.TrimSpace(string(bodyBytes))
	// Remove quotes if present
	bodyStr = strings.Trim(bodyStr, `"`)

	slog.Info("Heartbeat response", "msg", bodyStr)

	var hr struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	parsed := false
	if err := json.Unmarshal(bodyBytes, &hr); err == nil {
		parsed = true
	} else {
		var inner string
		if err2 := json.Unmarshal(bodyBytes, &inner); err2 == nil {
			if err3 := json.Unmarshal([]byte(inner), &hr); err3 == nil {
				parsed = true
			}
		}
	}

	if parsed {
		lowerMsg := strings.ToLower(hr.Message)
		if hr.Code >= 400 && (strings.Contains(lowerMsg, "device not found") || strings.Contains(lowerMsg, "revoked") || strings.Contains(lowerMsg, "mismatch")) {
			slog.Error("License heartbeat failed (structured) - revocation detected", "code", hr.Code, "message", hr.Message)
			c.handleLicenseRevoked()
		} else {
			slog.Debug("License heartbeat successful (structured)", "status", hr.Status, "message", hr.Message)
		}
	} else {
		switch strings.ToLower(bodyStr) {
		case "missing required fields":
			slog.Warn("License heartbeat failed - missing required fields", "response", bodyStr)
		case "device not found":
			slog.Error("License heartbeat failed - device not found", "response", bodyStr)
			c.handleLicenseRevoked()
		case "license mismatch":
			slog.Error("License heartbeat failed - license mismatch", "response", bodyStr)
			c.handleLicenseRevoked()
		case "device revoked":
			slog.Error("License heartbeat failed - device revoked", "response", bodyStr)
			c.handleLicenseRevoked()
		case "ok":
			slog.Debug("License heartbeat successful", "status", bodyStr)
		default:
			slog.Debug("License heartbeat successful", "status", bodyStr)
		}
	}

	return nil
}

func (c *CronHandler) handleLicenseRevoked() {
	slog.Error("CRITICAL: License has been REVOKED by server. Removing local license files.")

	runtime.GC()
	time.Sleep(100 * time.Millisecond) // Give a moment for GC
	_, licensePath, pubKeyPath, revokedPath, _ := licenseDirPaths()

	c.removeFileWithRetry(licensePath)
	c.removeFileWithRetry(pubKeyPath)

	revocationData := map[string]interface{}{
		"revoked_at": time.Now().Unix(),
		"reason":     "license_revoked_by_server",
		"action":     "contact_support_or_reactivate",
	}

	if data, err := json.MarshalIndent(revocationData, "", "  "); err == nil {
		if err := os.MkdirAll(filepath.Dir(revokedPath), 0755); err != nil {
			slog.Warn("Failed create revoked dir", "error", err)
		}
		if err := os.WriteFile(revokedPath, data, 0644); err != nil {
			slog.Error("Failed to create revoked.json", "error", err)
		} else {
			slog.Info("Successfully created revoked.json")
		}
	}

	slog.Error("Application will need license reactivation to continue")
}

func (c *CronHandler) removeFileWithRetry(filePath string) {
	maxRetries := 3
	retryDelay := 200 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		err := os.Remove(filePath)
		if err == nil {
			slog.Info("Successfully removed file", "file", filePath, "attempt", i+1)
			return
		}

		if os.IsNotExist(err) {
			slog.Info("File already does not exist", "file", filePath)
			return
		}

		slog.Warn("Failed to remove file, retrying...",
			"file", filePath,
			"attempt", i+1,
			"max_retries", maxRetries,
			"error", err.Error())

		if i < maxRetries-1 {
			time.Sleep(retryDelay)
			retryDelay *= 2 // Exponential backoff
		}
	}

	slog.Error("Failed to remove file after all retries", "file", filePath)
}

func (c *CronHandler) BackupDB(ctx context.Context) error {
	// Use LOCALAPPDATA for consistency with license and database paths
	var programRoot string
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		programRoot = filepath.Join(localAppData, "TAMALabs")
	} else {
		// Fallback to APPDATA
		appData := os.Getenv("APPDATA")
		if appData != "" {
			programRoot = filepath.Join(appData, "TAMALabs")
		} else if runtime.GOOS == "windows" {
			// Last resort: ProgramData (backward compatibility)
			programData := os.Getenv("ProgramData")
			if programData == "" {
				programData = `C:\ProgramData`
			}
			programRoot = filepath.Join(programData, "TAMALabs")
		} else {
			return fmt.Errorf("unable to determine application data directory")
		}
	}

	srcDB := filepath.Join(programRoot, "database", "TAMALabs.db")
	backupDir := filepath.Join(programRoot, "backup")

	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup folder: %w", err)
	}

	backupFile := filepath.Join(backupDir, fmt.Sprintf("TAMALabs-%s.db", time.Now().Format("20060102_150405")))

	// use the "sqlite" driver provided by modernc.org/sqlite (no cgo required)
	db, err := sql.Open("sqlite", srcDB)
	if err != nil {
		return fmt.Errorf("failed to open source DB: %w", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("VACUUM INTO '%s';", backupFile))
	if err != nil {
		return fmt.Errorf("backup failed: %w", err)
	}

	return nil
}
