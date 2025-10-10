package app

import (
	"crypto/rand"
	"fmt"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

var seedConfig = []entity.Config{
	{
		ID:    "Port",
		Value: "8322",
	},
	{
		ID:    "LogLevel",
		Value: "info",
	},
	{
		ID:    "SigningKey",
		Value: string(generateRandomSigningKey(32)),
	},
	{
		ID:    "KhanzaIntegrationEnabled",
		Value: "false",
	},
	{
		ID:    "KhanzaBridgeDatabaseDSN",
		Value: "root:secret@tcp(localhost:3306)/khanza?parseTime=true",
	},
	{
		ID:    "KhanzaMainDatabaseDSN",
		Value: "root:secret@tcp(localhost:3306)/khanza?parseTime=true",
	},
	{
		ID:    "SimrsIntegrationEnabled",
		Value: "false",
	},
	{
		ID:    "SimrsDatabaseDSN",
		Value: "root:secret@tcp(localhost:3306)/simrs_db?charset=utf8mb4&parseTime=True&loc=Local",
	},
}

// GenerateRandomSigningKey creates a cryptographically secure random signing key.
//
// Recommended key length for HMAC-SHA algorithms (like HS256, HS384, HS512)
// is at least 32 bytes (256 bits). For stronger algorithms or longer security,
// you might consider longer keys.
//
// The function returns the key as a byte slice and an error if key generation fails.
func generateRandomSigningKey(keyLength int) []byte {
	if keyLength <= 0 {
		panic(fmt.Errorf("key length must be positive"))
	}

	// Create a byte slice of the specified length to hold the random key
	randomBytes := make([]byte, keyLength)

	// Read random bytes from crypto/rand.Reader.
	// crypto/rand.Reader is a cryptographically secure random number generator.
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(fmt.Errorf("failed to generate random key: %w", err))
	}

	return randomBytes
}
