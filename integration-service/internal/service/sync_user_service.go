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
	"strconv"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/dto"
	"github.com/BioSystems-Indonesia/integration-services-lis/internal/helper"
	repositories "github.com/BioSystems-Indonesia/integration-services-lis/internal/repository"
	"github.com/google/uuid"
)

type UserSyncService struct {
	repo       repositories.UserRepository
	client     *http.Client
	targetURL  string
	labKeyPath string
	keysDir    string
}

func NewUserSyncService(
	repo repositories.UserRepository,
	targetURL string,
	labKeyPath string,
	keysDir string,
) *UserSyncService {
	return &UserSyncService{
		repo: repo,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		targetURL:  targetURL,
		labKeyPath: labKeyPath,
		keysDir:    keysDir,
	}
}

func (s *UserSyncService) SyncUser(ctx context.Context) error {
	labKey, err := helper.LoadLabKey(s.labKeyPath)
	if err != nil {
		return err
	}

	privateKey, err := helper.LoadPrivateKey(filepath.Join(s.keysDir, labKey.KeyId, "private.pem"))
	if err != nil {
		return err
	}

	admins, err := s.repo.FindAll(ctx)
	if err != nil {
		return err
	}

	if len(admins) == 0 {
		return nil
	}

	syncTime := time.Now()
	var syncedUserIds []int

	for _, admin := range admins {
		var status string
		if admin.IsActive {
			status = "ACTIVE"
		}
		payload := dto.UserRequest{
			UserID:   fmt.Sprintf("USRID%s", strconv.Itoa(int(admin.ID))),
			LabID:    labKey.LabId,
			Fullname: admin.Fullname,
			Username: admin.Username,
			Password: admin.PasswordHash,
			Role:     admin.Roles[0].ID,
			Status:   status,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal payload failed: %w", err)
		}
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			s.targetURL,
			bytes.NewBuffer(body),
		)

		if err != nil {
			return fmt.Errorf("create request failed: %w", err)
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
			return fmt.Errorf("sign payload failed: %w", err)
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

		respBody, err := ioutil.ReadAll(resp.Body)
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			syncedUserIds = append(syncedUserIds, int(admin.ID))
		}
		fmt.Println("Response:", string(respBody))
	}

	// Only update last_sync for successfully synced users
	if len(syncedUserIds) > 0 {
		if err := s.repo.UpdateLastSync(ctx, syncedUserIds, syncTime); err != nil {
			return fmt.Errorf("update last_sync failed: %w", err)
		}
	}

	return nil
}
