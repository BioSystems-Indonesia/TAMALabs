package license

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
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

type ExternalActivateRequest struct {
	MachineID   string            `json:"machine_id"`
	LicenseCode string            `json:"license_code"`
	Meta        map[string]string `json:"meta"`
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

var payload struct {
	MachineID  string `json:"machine_id"`
	IssuedAt   string `json:"issued_at"`
	ExpiresAt  string `json:"expires_at"`
	DeviceMeta struct {
		Arch     string `json:"arch"`
		Hostname string `json:"hostname"`
		OS       string `json:"os"`
	} `json:"device_meta"`
}

func NewLicense(pubLoader PublicKeyLoader, fileLoader LicenseFileLoader, pubKeyPath, licensePath string) *License {
	return &License{
		pubLoader:   pubLoader,
		fileLoader:  fileLoader,
		pubKeyPath:  pubKeyPath,
		licensePath: licensePath,
	}
}

type LicenseInformation struct {
	MachineInformation string `json:"machine_information"`
	IssueDate          string `json:"issue_date"`
	LicenseType        string `json:"license_type"`
	ExpirationDate     string `json:"expiration_date"`
}

func verifySignature(pub *rsa.PublicKey, payload, sigB64 string) error {
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return err
	}
	hash := sha256.Sum256([]byte(payload))
	return rsa.VerifyPKCS1v15(pub, crypto.SHA256, hash[:], sig)
}

func (l *License) Check() (*LicenseInformation, error) {
	pubKey, err := l.pubLoader.LoadPublicKey(l.pubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	// Check for revocation/expiration markers under AppData/Local/TAMALabs/license
	// Use same logic as database path - prefer LOCALAPPDATA over ProgramData
	var licenseDir string
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData != "" {
		licenseDir = filepath.Join(localAppData, "TAMALabs", "license")
	} else {
		appData := os.Getenv("APPDATA")
		if appData != "" {
			licenseDir = filepath.Join(appData, "TAMALabs", "license")
		} else if runtime.GOOS == "windows" {
			// Fallback to ProgramData for backward compatibility
			prog := os.Getenv("ProgramData")
			if prog == "" {
				prog = `C:\ProgramData`
			}
			licenseDir = filepath.Join(prog, "TAMALabs", "license")
		}
	}

	revokedPath := filepath.Join(licenseDir, "revoked.json")
	expiredPath := filepath.Join(licenseDir, "expired.json")

	if _, err := os.Stat(revokedPath); err == nil {
		return nil, fmt.Errorf("license has been revoked by server. please contact support or reactivate")
	}

	if _, err := os.Stat(expiredPath); err == nil {
		return nil, fmt.Errorf("license has expired. please renew your license")
	}

	data, err := l.fileLoader.LoadFile(l.licensePath)
	if err != nil {
		return nil, fmt.Errorf("license not found. please activate first")
	}

	var lic LicenseFile
	if err := json.Unmarshal(data, &lic); err != nil {
		return nil, fmt.Errorf("invalid license file: %w", err)
	}

	if err := json.Unmarshal([]byte(lic.Payload), &payload); err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	machineId, err := util.GenerateMachineID()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %w", err)
	}

	payload.MachineID = machineId
	newPayloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode payload: %w", err)
	}

	lic.Payload = string(newPayloadBytes)

	if err := verifySignature(pubKey, lic.Payload, lic.Signature); err != nil {
		return nil, fmt.Errorf("license verification failed: %w", err)
	}

	var payload LicensePayload
	if err := json.Unmarshal([]byte(lic.Payload), &payload); err != nil {
		return nil, fmt.Errorf("invalid payload: %w", err)
	}

	if time.Now().After(payload.ExpiresAt) {
		return nil, fmt.Errorf("license expired at %v", payload.ExpiresAt)
	}

	// Create and return license information
	licenseType := "Standard" // Default
	if payload.Meta != nil {
		if lt, exists := payload.Meta["license_type"]; exists {
			licenseType = lt
		}
	}

	licenseInfo := &LicenseInformation{
		MachineInformation: payload.MachineID,
		IssueDate:          payload.IssuedAt.Format("2006-01-02 15:04:05"),
		LicenseType:        licenseType,
		ExpirationDate:     payload.ExpiresAt.Format("2006-01-02 15:04:05"),
	}

	return licenseInfo, nil
}
