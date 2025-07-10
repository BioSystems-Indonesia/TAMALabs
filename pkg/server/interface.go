package server

import (
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/constant"
)

// Controller is a common interface for both TCP and Serial servers
type Controller interface {
	SetPort(port string)
	Start() error
	State() constant.ServerState
	Serve()
	Stop() error
}

// Shared constants for server implementations
const disconnectTimeout = 10 * time.Second
