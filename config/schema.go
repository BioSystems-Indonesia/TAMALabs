package config

// Schema is a struct that contains the configuration schema.
type Schema struct {
	Name           string `validate:"required"`
	Port           string `validate:"required"`
	FrontendFolder string `validate:"required"`
	Version        string `validate:"required"`
	LogLevel       string `validate:"required"`
	Revision       string
}
