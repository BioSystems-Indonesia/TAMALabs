package license

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FSKeyLoader implements PublicKeyLoader using the local filesystem PEM file.
type FSKeyLoader struct{}

func NewFSKeyLoader() *FSKeyLoader { return &FSKeyLoader{} }

var API_KEY = "KJKDANCJSANIUWYR6243UJFOISJFJKVOMV72487YEHFHFHSDVOHF9AMDC9AN9SDN98YE98YEHDIU2Y897873YYY68686487WGDUDUAGYTE8QTEYADIUHADUYW8E8BWTNC8N8NAMDOAIMDAUDUWYAD87NYW7Y7CBT87EY8142164B36248732M87MCIFH8NYRWCM8MYCMUOIDOIADOIDOIUR83YR983Y98328N32C83NYC8732NYC8732Y87Y32NCNSAIHJAOJFOIJFOIQFIUIUNCNHCIUHWV8NRYNV8Y989N9198298YOIJOI090103021313JKJDHAHDJAJASHHAH"

// downloadPublicKey downloads public key from license server
func (l *FSKeyLoader) downloadPublicKey() ([]byte, error) {
	client := &http.Client{}
	// Get license server URL from environment variable
	licenseServerURL := os.Getenv("LICENSE_SERVER_URL")
	if licenseServerURL == "" {
		licenseServerURL = "http://localhost:8080" // Default fallback
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/pubkey", licenseServerURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to download public key from server: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-API-KEY", API_KEY)

	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to download public key from server: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d when downloading public key", resp.StatusCode)
	}

	// Read the response body
	keyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key response: %v", err)
	}

	return keyData, nil
}

func (l *FSKeyLoader) LoadPublicKey(path string) (*rsa.PublicKey, error) {
	var pemBytes []byte
	var err error

	// Try to read from local file first
	pemBytes, err = os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, try to download from server
		if os.IsNotExist(err) {
			pemBytes, err = l.downloadPublicKey()
			if err != nil {
				return nil, fmt.Errorf("failed to load public key from file and server: %v", err)
			}

			// Ensure directory exists before saving file
			dir := filepath.Dir(path)
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Printf("Warning: failed to create directory %s: %v\n", dir, err)
			}

			// Save downloaded key to local file for future use
			if saveErr := os.WriteFile(path, pemBytes, 0644); saveErr != nil {
				// Log the error but don't fail the operation
				fmt.Printf("Warning: failed to save public key to %s: %v\n", path, saveErr)
			}
		} else {
			return nil, err
		}
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("no pem data")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

// FSFileLoader implements LicenseFileLoader using local filesystem.
type FSFileLoader struct{}

func NewFSFileLoader() *FSFileLoader { return &FSFileLoader{} }

func (l *FSFileLoader) LoadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
