package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/dto"
	"github.com/BioSystems-Indonesia/integration-services-lis/internal/helper"
	repositories "github.com/BioSystems-Indonesia/integration-services-lis/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GeneratePublicLinkResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   string `json:"data"`
}

type ObservationResultSyncService struct {
	repo       repositories.ObservationResultRepository
	user       repositories.UserRepository
	client     *http.Client
	baseUrl    string
	labKeyPath string
	keysDir    string
}

func NewObservationResultSyncService(
	repo repositories.ObservationResultRepository,
	user repositories.UserRepository,
	baseUrl string,
	labKeyPath string,
	keysDir string,
) *ObservationResultSyncService {
	return &ObservationResultSyncService{
		repo: repo,
		user: user,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl:    baseUrl,
		labKeyPath: labKeyPath,
		keysDir:    keysDir,
	}
}

func (s *ObservationResultSyncService) SyncObservationResult(ctx context.Context) error {
	labKey, err := helper.LoadLabKey(s.labKeyPath)
	if err != nil {
		return err
	}

	privateKey, err := helper.LoadPrivateKey(filepath.Join(s.keysDir, labKey.KeyId, "private.pem"))
	if err != nil {
		return err
	}

	specimens, err := s.repo.FindAll(ctx)
	if err != nil {
		return err
	}

	if len(specimens) == 0 {
		return nil
	}

	syncTime := time.Now()
	var syncedObservationResultIds []int

	for _, specimen := range specimens {
		var items []dto.ObservationResultItem
		for _, item := range specimen.ObservationResult {
			var createdBy string
			switch item.CreatedBy {
			case -2:
				createdBy = "Analyzer"
			case -1:
				createdBy = "Unknown"
			default:
				user, _ := s.user.FindById(ctx, int(item.CreatedBy))
				createdBy = user.Fullname
			}

			flag := ""
			if len(item.Values) > 0 && item.ReferenceRange != "" {
				flag = helper.DetermineFlagFromReference(item.Values[0], item.ReferenceRange)
			}

			if flag == "" && len(item.AbnormalFlag) > 0 {
				flag = item.AbnormalFlag[0]
			}

			var category string
			if item.TestTypeID != nil {
				for _, req := range specimen.ObservationRequest {
					if req.TestTypeID != nil && *req.TestTypeID == *item.TestTypeID {
						category = req.TestType.Category
						break
					}
				}
			}
			if category == "" {
				for _, req := range specimen.ObservationRequest {
					if req.TestCode == item.Code {
						category = req.TestType.Category
						break
					}
				}
			}

			var value string
			if len(item.Values) > 0 {
				value = item.Values[0]
			}

			items = append(items, dto.ObservationResultItem{
				Id:       fmt.Sprintf("%sItemTEST%s%d", labKey.LabId, item.Code, item.SpecimenID),
				Category: category,
				Code:     fmt.Sprintf("%s|%s", labKey.LabId, item.Code),
				Value:    value,
				Flag:     flag,
				AddedBy:  createdBy,
			})

		}

		isVerified := false
		if specimen.WorkOrder.VerifiedStatus == "VERIFIED" {
			isVerified = true
		}
		payload := dto.ObservationResultRequest{
			Id:             fmt.Sprintf("%sORD%s", labKey.LabId, specimen.WorkOrder.Barcode),
			LabID:          labKey.LabId,
			OrderId:        specimen.WorkOrder.Barcode,
			CollectionDate: specimen.CollectionDate,
			ComplatedAt:    specimen.WorkOrder.UpdatedAt,
			IsVerified:     isVerified,
			Patient: dto.PatientRequest{
				PatientID:           fmt.Sprintf("%sPAT%d", labKey.LabId, specimen.Patient.ID),
				LabID:               labKey.LabId,
				FirstName:           specimen.Patient.FirstName,
				LastName:            specimen.Patient.LastName,
				Gender:              string(specimen.Patient.Sex),
				Birthdate:           specimen.Patient.Birthdate,
				Address:             specimen.Patient.Address,
				MedicalRecordNumber: specimen.WorkOrder.MedicalRecordNumber,
			},
			Items: items,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("https://%s/protected/observation-result/", s.baseUrl),
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

		fmt.Println(string(body))
		resp, err := s.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		respBody, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Collect all observation_result IDs from this specimen
			for _, obsResult := range specimen.ObservationResult {
				syncedObservationResultIds = append(syncedObservationResultIds, int(obsResult.ID))
			}
		}
		fmt.Println("Response:", string(respBody))
	}

	// Only update last_sync for successfully synced observation results
	if len(syncedObservationResultIds) > 0 {
		if err := s.repo.UpdateLastSync(ctx, syncedObservationResultIds, syncTime); err != nil {
			return fmt.Errorf("update last_sync failed: %w", err)
		}
	}

	return nil

}

func (s *ObservationResultSyncService) GeneratePublicLink(ctx context.Context, apiKey string, barcode string) (*GeneratePublicLinkResponse, error) {
	log.Printf("🔗 [GeneratePublicLink] Starting public link generation for barcode: %s", barcode)

	labKey, err := helper.LoadLabKey(s.labKeyPath)
	if err != nil {
		log.Printf("❌ [GeneratePublicLink] Failed to load lab key from path '%s': %v", s.labKeyPath, err)
		return nil, err
	}

	resultId := fmt.Sprintf("%sORD%s", labKey.LabId, barcode)
	log.Printf("🔍 [GeneratePublicLink] Generated result ID: %s (LabId: %s, Barcode: %s)", resultId, labKey.LabId, barcode)

	url := fmt.Sprintf("https://%s/unauthenticated/observation-result/generate/%s", s.baseUrl, resultId)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		nil,
	)
	if err != nil {
		log.Printf("❌ [GeneratePublicLink] Failed to create HTTP request for URL '%s': %v", url, err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-API-Key", apiKey)

	log.Printf("📤 [GeneratePublicLink] Sending POST request to: %s", url)

	resp, err := s.client.Do(req)
	if err != nil {
		log.Printf("❌ [GeneratePublicLink] HTTP request failed for URL '%s': %v", url, err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		log.Printf("❌ [GeneratePublicLink] Server returned error status %d for barcode '%s'. Response body: %s", resp.StatusCode, barcode, string(respBody))
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ [GeneratePublicLink] Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("📥 [GeneratePublicLink] Received response body: %s", string(respBody))

	var serverResponse map[string]interface{}
	if err := json.Unmarshal(respBody, &serverResponse); err != nil {
		log.Printf("❌ [GeneratePublicLink] Failed to parse JSON response: %v. Raw response: %s", err, string(respBody))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	data := ""
	if d, ok := serverResponse["data"]; ok {
		data = fmt.Sprintf("%v", d)
		log.Printf("✅ [GeneratePublicLink] Successfully generated public link for barcode '%s': %s", barcode, data)
	} else {
		log.Printf("⚠️ [GeneratePublicLink] Response doesn't contain 'data' field for barcode '%s'. Full response: %+v", barcode, serverResponse)
	}

	response := &GeneratePublicLinkResponse{
		Code:   200,
		Status: "Ok",
		Data:   data,
	}

	return response, nil
}

func (s *ObservationResultSyncService) VerifyObservationResult(ctx context.Context, apiKey string) error {
	labKey, err := helper.LoadLabKey(s.labKeyPath)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("wss://%s/ws/v1/verification/stream?lab_id=%s&api_key=%s", s.baseUrl, labKey.LabId, apiKey)
	log.Println("🛜 Connecting to verification stream")

	conn, _, err := websocket.DefaultDialer.Dial(path, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to websocket: %w", err)
	}
	log.Println("🛜 Connected to verification stream")
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		var verificationRequest dto.VerifyObservationResultRequest
		if err := json.Unmarshal(message, &verificationRequest); err != nil {
			log.Println("Failed to unmarshal verification request:", err)
			continue
		}

		parts := strings.Split(verificationRequest.ResultId, "ORD")
		if len(parts) < 2 {
			continue
		}
		labResult := parts[1]

		log.Printf("Received verification request: %+v\n", verificationRequest)
		if err := s.repo.Verify(ctx, labResult, verificationRequest.Status); err != nil {
			log.Println("Failed to verify observation result:", err)
			continue
		}
		log.Printf("Successfully verified observation result ID: %s with status: %s\n", labResult, verificationRequest.Status)
	}
}
