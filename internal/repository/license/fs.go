package license

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// FSKeyLoader implements PublicKeyLoader using the local filesystem PEM file.
type FSKeyLoader struct{}

func NewFSKeyLoader() *FSKeyLoader { return &FSKeyLoader{} }

func (l *FSKeyLoader) LoadPublicKey(path string) (*rsa.PublicKey, error) {
	pemBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
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
