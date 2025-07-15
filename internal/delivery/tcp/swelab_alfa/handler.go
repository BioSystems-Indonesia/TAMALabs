package swelabalfa

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp"
	"github.com/oibacidem/lims-hl-seven/pkg/panics"
)

type Handler struct {
	analyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		analyzerUsecase: analyzerUsecase,
	}
}

func (h *Handler) Handle(conn *net.TCPConn) {
	ctx := context.Background()

	defer conn.Close()
	defer panics.RecoverPanic(ctx)

	mc := mllp.NewClient(conn)
	b, err := mc.ReadAll()
	if err != nil {
		if !errors.Is(err, io.EOF) {
			slog.Error(err.Error())
		}
	}

	res, err := h.handleMessage(ctx, string(b))
	if err != nil {
		slog.Error(err.Error())
	}

	if res != "" {
		resLog := strings.ReplaceAll(res, "\r", "\n")
		slog.Info(fmt.Sprintf("ack message: %s", resLog))
	}

	if err := mc.Write([]byte(res)); err != nil {
		slog.Error(err.Error())
	}
}

func (h *Handler) handleMessage(ctx context.Context, message string) (string, error) {
	// don't do anything if the message is empty
	if message == "" {
		return "", nil
	}

	logMsg := strings.ReplaceAll(message, "\r", "\n")
	slog.Info("received message", "message", logMsg)

	msgByte := []byte(message)
	headerDecoder := hl7.NewDecoder(h251.Registry, &hl7.DecodeOption{HeaderOnly: true})
	header, err := headerDecoder.Decode(msgByte)
	if err != nil {
		return "", fmt.Errorf("decode header failed: %w", err)
	}

	switch m := header.(type) {
	case h251.ORM_O01:
		return h.ORMO01(ctx, m, msgByte)
	case h251.ORU_R01:
		return h.ORUR01(ctx, m, msgByte)
	}

	return "", fmt.Errorf("unknown message type %T", header)
}
