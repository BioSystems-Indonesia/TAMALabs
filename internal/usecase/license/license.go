package license

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

type LicensePayload struct {
	MachineID string            `json:"machine_id"`
	IssuedAt  time.Time         `json:"issued_at"`
	ExpiresAt time.Time         `json:"expires_at"`
	Meta      map[string]string `json:"device_meta"`
}

type LicenseFile struct {
	Payload   string `json:"payload"`
	Signature string `json:"signature"`
}

type PublicKeyLoader interface {
	LoadPublicKey(path string) (*rsa.PublicKey, error)
}

type LicenseFileLoader interface {
	LoadFile(path string) ([]byte, error)
}

type License struct {
	pubLoader   PublicKeyLoader
	fileLoader  LicenseFileLoader
	pubKeyPath  string
	licensePath string
}

func NewLicense(pubLoader PublicKeyLoader, fileLoader LicenseFileLoader, pubKeyPath, licensePath string) *License {
	return &License{
		pubLoader:   pubLoader,
		fileLoader:  fileLoader,
		pubKeyPath:  pubKeyPath,
		licensePath: licensePath,
	}
}

func verifySignature(pub *rsa.PublicKey, payload, sigB64 string) error {
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(payload))
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig)
}

func (l *License) Check() error {
	data, err := l.fileLoader.LoadFile(l.licensePath)
	if err != nil {
		return fmt.Errorf("license not found. please activate first: %w", err)
	}

	var lic LicenseFile
	if err := json.Unmarshal(data, &lic); err != nil {
		return fmt.Errorf("invalid license file: %w", err)
	}

	pubKey, err := l.pubLoader.LoadPublicKey(l.pubKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}

	if err := verifySignature(pubKey, lic.Payload, lic.Signature); err != nil {
		return fmt.Errorf("license verification failed: %w", err)
	}

	var payload LicensePayload
	if err := json.Unmarshal([]byte(lic.Payload), &payload); err != nil {
		return fmt.Errorf("invalid payload: %w", err)
	}

	if time.Now().After(payload.ExpiresAt) {
		return fmt.Errorf("license expired at %v", payload.ExpiresAt)
	}

	return nil
}
