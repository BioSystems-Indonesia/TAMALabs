package util

import (
	"os"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
)

func IsDevelopment() bool {
	return os.Getenv(constant.ENVKey) == string(constant.EnvDevelopment)
}
