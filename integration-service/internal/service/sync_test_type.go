package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/dto"
	"github.com/BioSystems-Indonesia/integration-services-lis/internal/helper"
	repositories "github.com/BioSystems-Indonesia/integration-services-lis/internal/repository"
	"github.com/google/uuid"
)

type SyncTestType struct {
	repo       repositories.TestTypeRepository
	client     *http.Client
	targetURL  string
	labKeyPath string
	keysDir    string
}

func NewSyncTestType(
	repo repositories.TestTypeRepository,
	targetURL string,
	labKeyPath string,
	keysDir string,
) *SyncTestType {
	return &SyncTestType{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		targetURL:  targetURL,
		labKeyPath: labKeyPath,
		keysDir:    keysDir,
	}
}

func (s *SyncTestType) SyncTestType(ctx context.Context) error {
	labKey, err := helper.LoadLabKey(s.labKeyPath)
	if err != nil {
		return err
	}

	privateKey, err := helper.LoadPrivateKey(filepath.Join(s.keysDir, labKey.KeyId, "private.pem"))
	if err != nil {
		return err
	}

	testTypes, err := s.repo.FindAll(ctx)
	if err != nil {
		return err
	}

	if len(testTypes) == 0 {
		return nil
	}

	syncTime := time.Now()
	var syncedTestTypeIds []int

	for _, testType := range testTypes {
		unit := testType.Unit
		if unit == "" {
			unit = "-"
		}

		payload := dto.TestTypeRequest{
			Id:            fmt.Sprintf("%s%s", labKey.LabId, testType.Code),
			LabId:         labKey.LabId,
			SpecimentType: testType.Type,
			CategoryName:  testType.Category,
			Code:          fmt.Sprintf("%s|%s", labKey.LabId, testType.Code),
			Unit:          unit,
			Ref:           fmt.Sprintf("%f - %f", testType.LowRefRange, testType.HighRefRange),
			Decimal:       testType.Decimal,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			s.targetURL,
			bytes.NewBuffer(body),
		)

		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		method := req.Method
		bodyHash := sha256.Sum256(body)
		bodyHashBase64 := base64.StdEncoding.EncodeToString(bodyHash[:])
		path := req.URL.Path
		nonce := uuid.NewString()

		stringToSign := strings.Join([]string{
			method,
			path,
			nonce,
			bodyHashBase64,
		}, "\n")

		signature, err := helper.Sign(privateKey, stringToSign)
		if err != nil {
			return fmt.Errorf("sing payload failed: %w", err)
		}

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Key-Id", labKey.KeyId)
		req.Header.Add("X-Nonce", nonce)
		req.Header.Add("X-Signature", signature)

		resp, err := s.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		respBody, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			syncedTestTypeIds = append(syncedTestTypeIds, int(testType.ID))
		}
		fmt.Println("Response:", string(respBody))
	}

	// Only update last_sync for successfully synced test types
	if len(syncedTestTypeIds) > 0 {
		if err := s.repo.UpdateLastSync(ctx, syncedTestTypeIds, syncTime); err != nil {
			return fmt.Errorf("update last_sync failed: %w", err)
		}
	}

	return nil
}
