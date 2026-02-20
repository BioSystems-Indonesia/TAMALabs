// services/keypair_service.go
package services

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
	repositories "github.com/BioSystems-Indonesia/integration-services-lis/internal/repository"
)

type KeyPairService struct {
	repo repositories.KeyPairRepository
}

func NewKeyPairService(repo repositories.KeyPairRepository) *KeyPairService {
	return &KeyPairService{repo: repo}
}

func (s *KeyPairService) EnsureKeyPair(keyID string) (*models.KeyPair, error) {
	if s.repo.Exists(keyID) {
		return s.repo.Load(keyID)
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	})

	publicBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	publicPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicBytes,
	})

	pair := &models.KeyPair{
		PrivateKeyPEM: privatePEM,
		PublicKeyPEM:  publicPEM,
	}

	if err := s.repo.Save(keyID, pair); err != nil {
		return nil, err
	}

	return pair, nil
}

func (s *KeyPairService) SendPublicKey(baseUrl string, keyDir string, labId string, apiKey string) error {
	pubKeyPath := filepath.Join(keyDir, "public.pem")
	file, err := os.Open(pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "public.pem")
	if err != nil {
		return fmt.Errorf("failed to create multipart file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("failed to copy public key: %w", err)
	}

	if err := writer.Close(); err != nil {
		return err
	}

	url := fmt.Sprintf("https://%s/unauthenticated/public-key/%s", baseUrl, labId)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", apiKey)
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload public key failed (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
