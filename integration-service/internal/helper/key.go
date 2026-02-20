package helper

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/models"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM")
	}

	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	keyIfc, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := keyIfc.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA private key")
	}

	return key, nil
}

func Sign(privateKey *rsa.PrivateKey, data string) (string, error) {
	hash := sha256.Sum256([]byte(data))

	sig, err := rsa.SignPKCS1v15(
		rand.Reader,
		privateKey,
		crypto.SHA256,
		hash[:],
	)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sig), nil
}

func ReadJSONMaybeBase64(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	trimmed := bytes.TrimSpace(data)

	if len(trimmed) > 0 && trimmed[0] == '{' {
		return json.Unmarshal(trimmed, v)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(trimmed))
	if err != nil {
		return fmt.Errorf("invalid json or base64")
	}

	return json.Unmarshal(decoded, v)
}

func LoadLabKey(path string) (*models.LabKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read lab key file failed: %w", err)
	}

	raw := strings.TrimSpace(string(data))

	if strings.HasPrefix(raw, `"`) {
		raw, err = strconv.Unquote(raw)
		if err != nil {
			return nil, fmt.Errorf("unquote lab key failed: %w", err)
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, fmt.Errorf("base64 decode lab key failed: %w", err)
	}

	var labKey models.LabKey
	if err := json.Unmarshal(decoded, &labKey); err != nil {
		return nil, fmt.Errorf("unmarshal lab key failed: %w", err)
	}

	return &labKey, nil
}
