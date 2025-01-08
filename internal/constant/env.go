package constant

const ENVKey = "LIMS_ENV"

type Env string

const (
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
)
