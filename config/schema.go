package config

// Schema is a struct that contains the configuration schema.
type Schema struct {
	Port     string `validate:"required"`
	LogLevel string `validate:"required"`

	// Runtime
	Version  string `validate:"required"`
	Revision string
}
