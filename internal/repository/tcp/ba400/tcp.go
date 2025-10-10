package ba400

import (
	"log/slog"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/tcp"
)

// TCP is a struct handle TCP connection
type TCP struct {
	tcp tcp.TCPer
}

// NewTCP returns a new TCP connection
func NewTCP(cfg *config.Schema) *TCP {
	conn, err := tcp.NewTCP("localhost", 5678, 10*time.Second)
	if err != nil {
		slog.Error("failed to connect tcp", "error", err)
	}
	return &TCP{
		tcp: conn,
	}
}
