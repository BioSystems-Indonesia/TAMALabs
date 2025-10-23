package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	licenseuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/license"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
	"github.com/labstack/echo/v4"
)

var API_KEY = "KJKDANCJSANIUWYR6243UJFOISJFJKVOMV72487YEHFHFHSDVOHF9AMDC9AN9SDN98YE98YEHDIU2Y897873YYY68686487WGDUDUAGYTE8QTEYADIUHADUYW8E8BWTNC8N8NAMDOAIMDAUDUWYAD87NYW7Y7CBT87EY8142164B36248732M87MCIFH8NYRWCM8MYCMUOIDOIADOIDOIUR83YR983Y98328N32C83NYC8732NYC8732Y87Y32NCNSAIHJAOJFOIJFOIQFIUIUNCNHCIUHWV8NRYNV8Y989N9198298YOIJOI090103021313JKJDHAHDJAJASHHAH"

type LicenseHandler struct {
	licenseUsecase *licenseuc.License
}

type ActivateRequest struct {
	LicenseCode string `json:"license_code"`
}

type ExternalActivateRequest struct {
	MachineID   string            `json:"machine_id"`
	LicenseCode string            `json:"license_code"`
	Meta        map[string]string `json:"meta"`
}

type Response struct {
	Code   int              `json:"code"`
	Status string           `json:"status"`
	Data   ActivateResponse `json:"data"`
}

type ActivateResponse struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type HeartbeatRequest struct {
	MachineID string            `json:"machine_id"`
	Meta      map[string]string `json:"meta"`
}

type HeartbeatResponse struct {
	Status    string `json:"status"` // "active", "revoked", "expired"
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func NewLicenseHandler(licenseUsecase *licenseuc.License) *LicenseHandler {
	return &LicenseHandler{
		licenseUsecase: licenseUsecase,
	}
}

// CheckLicense checks the current license status
func (h *LicenseHandler) CheckLicense(c echo.Context) error {
	licenseInfo, err := h.licenseUsecase.Check()

	if err != nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"valid":   false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"valid":        true,
		"message":      "License is valid",
		"machine_id":   licenseInfo.MachineInformation,
		"issued_at":    licenseInfo.IssueDate,
		"expires_at":   licenseInfo.ExpirationDate,
		"license_type": licenseInfo.LicenseType,
	})
}

// getSystemInfo collects system information from the server
func getSystemInfo() map[string]string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	return map[string]string{
		"hostname": hostname,
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
	}
}

// removeFileWithRetry attempts to remove a file with retry mechanism
func (h *LicenseHandler) removeFileWithRetry(filePath string) {
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

func (h *LicenseHandler) ActivateLicense(c echo.Context) error {
	client := &http.Client{}
	var req ActivateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	machineID, err := util.GenerateMachineID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate machine id",
		})
	}
	systemInfo := getSystemInfo()

	externalReq := ExternalActivateRequest{
		MachineID:   machineID,
		LicenseCode: req.LicenseCode,
		Meta:        systemInfo,
	}

	licenseServerURL := os.Getenv("LICENSE_SERVER_URL")
	if licenseServerURL == "" {
		licenseServerURL = "http://localhost:8080" // Default fallback
	}

	jsonData, err := json.Marshal(externalReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to prepare request",
		})
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/activate", licenseServerURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create request",
		})
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-Key", API_KEY)

	resp, err := client.Do(request)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{
			"error": "Cannot connect to license server",
		})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("License activation failed: %s", string(bodyBytes)),
		})
	}

	var activateResp Response
	if err := json.NewDecoder(resp.Body).Decode(&activateResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to parse license server response",
		})
	}

	licenseData := map[string]string{
		"license_code": req.LicenseCode,
		"payload":      activateResp.Data.Payload,
		"signature":    activateResp.Data.Signature,
	}

	licenseJSON, err := json.MarshalIndent(licenseData, "", "  ")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to format license data",
		})
	}

	if err := os.WriteFile("license/license.json", licenseJSON, 0644); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save license file",
		})
	}

	h.removeFileWithRetry("license/revoked.json")
	h.removeFileWithRetry("license/expired.json")

	return c.JSON(http.StatusOK, activateResp.Data)
}

func (h *LicenseHandler) RegisterRoutes(g *echo.Group) {
	g.GET("/license/check", h.CheckLicense)
	g.POST("/license/activate", h.ActivateLicense)
}
