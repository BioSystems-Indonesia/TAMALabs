package ba400

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/pkg/tcp"
)

// TCP is a struct handle TCP connection
type TCP struct {
	tcp tcp.TCPer
}

// NewTCP returns a new TCP connection
func NewTCP(cfg *config.Schema) *TCP {
	conn, err := tcp.NewTCP("localhost", 5678, 10*time.Second)
	if err != nil {
		log.Error(fmt.Sprintf("failed to connect tcp, error: %s", err.Error()))
	}
	return &TCP{
		tcp: conn,
	}
}
