package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/dto"
	"github.com/BioSystems-Indonesia/integration-services-lis/internal/helper"
	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
	"github.com/google/uuid"
)

type LabSyncService struct {
	targetUrl  string
	apiKey     string
	labKeyPath string
}

func NewLabSyncService(targetUrl string, apiKey string, labKeyPath string) *LabSyncService {
	return &LabSyncService{
		targetUrl:  targetUrl,
		apiKey:     apiKey,
		labKeyPath: labKeyPath,
	}
}

func (s *LabSyncService) CreateLab() error {
	data, err := os.ReadFile("lab_info.json")
	if err != nil {
		return fmt.Errorf("read lab_info.json failed: %w", err)
	}

	var cfg models.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	cfg.LabId = uuid.NewString()

	payload, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal payload failed: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.targetUrl,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.apiKey)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}

	var labResponse dto.HTTPResponse[dto.LabResponse]
	if err := json.Unmarshal(body, &labResponse); err != nil {
		return fmt.Errorf("unmarshal response failed: %w", err)
	}

	jsonByte, err := json.MarshalIndent(labResponse.Data, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lab response failed: %w", err)
	}

	err = helper.WriteJSON(s.labKeyPath, jsonByte)
	if err != nil {
		return fmt.Errorf("write lab_key.json failed: %w", err)
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"CreateLab failed: status=%d response=%s",
			resp.StatusCode,
			string(body),
		)
	}

	fmt.Println("CreateLab success:", labResponse)
	return nil
}
