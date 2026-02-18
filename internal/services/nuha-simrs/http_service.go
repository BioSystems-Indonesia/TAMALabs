package nuha_simrs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type NuhaService struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewNuhaService(baseURL string) *NuhaService {
	return &NuhaService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *NuhaService) GetLabList(
	ctx context.Context,
	req LabListRequest,
) (*LabListResponse, error) {
	url := fmt.Sprintf("%s/api/v1/emr/lab/list-new", s.BaseURL)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := strings.TrimSpace(string(bodyBytes))
		if len(bodyStr) > 2048 {
			bodyStr = bodyStr[:2048] + "...(truncated)"
		}
		return nil, fmt.Errorf("unexpected status code: %d from %s - response body: %s", resp.StatusCode, url, bodyStr)
	}

	var result LabListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}

func (s *NuhaService) InsertResult(
	ctx context.Context,
	req InsertResultRequest,
) (*InsertResultResponse, error) {
	url := fmt.Sprintf("%s/api/v1/emr/lab/insert", s.BaseURL)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	maskedReq := req
	if len(maskedReq.SessionID) > 6 {
		maskedReq.SessionID = maskedReq.SessionID[:6] + "..."
	}
	if pb, err := json.Marshal(maskedReq); err == nil {
		payload := string(pb)
		if len(payload) > 2048 {
			payload = payload[:2048] + "...(truncated)"
		}
		slog.InfoContext(ctx, "Nuha InsertResult request (masked)", "url", url, "payload", payload)
	} else {
		slog.WarnContext(ctx, "Failed to marshal InsertResult request for logging", "error", err)
	}

	resp, err := s.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := strings.TrimSpace(string(bodyBytes))
		if len(bodyStr) > 2048 {
			bodyStr = bodyStr[:2048] + "...(truncated)"
		}
		return nil, fmt.Errorf("unexpected status code: %d from %s - response body: %s", resp.StatusCode, url, bodyStr)
	}

	var result InsertResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}

func (s *NuhaService) BatchInsertResults(
	ctx context.Context,
	req BatchInsertResultRequest,
) (*InsertResultResponse, error) {
	url := fmt.Sprintf("%s/api/v1/emr/lab/insert-bulk", s.BaseURL)

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	maskedBatch := req
	if len(maskedBatch.SessionID) > 6 {
		maskedBatch.SessionID = maskedBatch.SessionID[:6] + "..."
	}
	if len(maskedBatch.Data) > 5 {
		maskedBatch.Data = maskedBatch.Data[:5]
	}
	if pb, err := json.Marshal(maskedBatch); err == nil {
		payload := string(pb)
		if len(payload) > 4096 {
			payload = payload[:4096] + "...(truncated)"
		}
		slog.InfoContext(ctx, "Nuha BatchInsertResults request (masked preview)", "url", url, "data_preview", payload, "original_count", len(req.Data))
	} else {
		slog.WarnContext(ctx, "Failed to marshal BatchInsertResults request for logging", "error", err)
	}

	resp, err := s.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyStr := strings.TrimSpace(string(bodyBytes))
		if len(bodyStr) > 2048 {
			bodyStr = bodyStr[:2048] + "...(truncated)"
		}
		return nil, fmt.Errorf("unexpected status code: %d from %s - response body: %s", resp.StatusCode, url, bodyStr)
	}

	var result InsertResultResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response failed: %w", err)
	}

	return &result, nil
}
