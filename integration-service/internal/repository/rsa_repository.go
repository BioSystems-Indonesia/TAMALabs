package repositories

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
)

type KeyPairRepository interface {
	Save(keyID string, pair *models.KeyPair) error
	Load(keyID string) (*models.KeyPair, error)
	Exists(keyID string) bool
}

type FileSystemKeyPairRepository struct {
	baseDir string
}

func NewFileSystemKeyPairRepository(baseDir string) *FileSystemKeyPairRepository {
	return &FileSystemKeyPairRepository{baseDir: baseDir}
}

func (r *FileSystemKeyPairRepository) Save(keyID string, pair *models.KeyPair) error {
	dir := filepath.Join(r.baseDir, keyID)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "private.pem"), pair.PrivateKeyPEM, 0600); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "public.pem"), pair.PublicKeyPEM, 0644); err != nil {
		return err
	}

	return nil
}

func (r *FileSystemKeyPairRepository) Load(keyID string) (*models.KeyPair, error) {
	dir := filepath.Join(r.baseDir, keyID)

	privatePem, err := os.ReadFile(filepath.Join(dir, "private.pem"))
	if err != nil {
		return nil, fmt.Errorf("private key not found")
	}

	publicPem, err := os.ReadFile(filepath.Join(dir, "public.pem"))
	if err != nil {
		return nil, fmt.Errorf("public key not found")
	}

	return &models.KeyPair{
		PrivateKeyPEM: privatePem,
		PublicKeyPEM:  publicPem,
	}, nil
}

func (r *FileSystemKeyPairRepository) Exists(keyID string) bool {
	privatePath := filepath.Join(r.baseDir, keyID, "private.pem")
	_, err := os.Stat(privatePath)
	return err == nil
}
