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

	// SIMGOS Configuration
	SimgosIntegrationEnabled string `validate:"-"`
	SimgosDatabaseDSN        string `validate:"-"`

	// TechnoMedic Configuration
	TechnoMedicIntegrationEnabled string `validate:"-"`

	// Backup Configuration
	BackupScheduleType string `validate:"-"` // "interval" or "daily"
	BackupInterval     string `validate:"-"` // hours between backups (for interval type)
	BackupTime         string `validate:"-"` // time of day for backup (for daily type) HH:MM format

	// Runtime
	Version  string `validate:"required"`
	Revision string
}
