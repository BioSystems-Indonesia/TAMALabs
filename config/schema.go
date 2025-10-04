package config

// Schema is a struct that contains the configuration schema.
type Schema struct {
	Port       string `validate:"required"`
	LogLevel   string `validate:"required"`
	SigningKey string `validate:"required"`

	// MySQL Configuration
	KhanzaIntegrationEnabled string `validate:"-"`
	KhanzaBridgeDatabaseDSN  string `validate:"-"`
	KhanzaMainDatabaseDSN    string `validate:"-"`

	// SIMRS Configuration
	SimrsIntegrationEnabled string `validate:"-"`
	SimrsDatabaseDSN        string `validate:"-"`

	// Runtime
	Version  string `validate:"required"`
	Revision string
}
