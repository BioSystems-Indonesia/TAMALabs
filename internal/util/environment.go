package util

import (
	"os"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
)

func IsDevelopment() bool {
	return os.Getenv(constant.ENVKey) == string(constant.EnvDevelopment)
}
