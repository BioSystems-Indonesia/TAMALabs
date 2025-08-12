package a15

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"

	"github.com/oibacidem/lims-hl-seven/pkg/mllp"
	"github.com/oibacidem/lims-hl-seven/pkg/panics"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(conn *net.TCPConn) {
	ctx := context.Background()

	defer panics.RecoverPanic(ctx)
	mc := mllp.NewClient(conn)
	for {
		message, err := mc.ReadAll()
		if err != nil {
			slog.Error("error reading mllp message", "error", err)
			if errors.Is(err, io.EOF) {
				break
			}
			return
		}
		if len(message) == 0 {
			break
		}

		slog.Info("read mllp message", "message", string(message))
	}

}
