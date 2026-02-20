package edan_h30

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/panics"
	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
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
			slog.Error("Error reading from connection", "error", err)
		}
		return
	}

	message := string(b)

	res, err := h.handleMessage(ctx, message)
	if err != nil {
		slog.Error("Error handling message", "error", err)
	}

	if res != "" {
		resLog := strings.ReplaceAll(res, "\r", "\n")
		slog.Info(fmt.Sprintf("ack message: %s", resLog))
	}

	if err := mc.Write([]byte(res)); err != nil {
		slog.Error("Error writing response", "error", err)
	}
}

func (h *Handler) handleMessage(ctx context.Context, message string) (string, error) {
	// don't do anything if the message is empty
	if message == "" {
		return "", nil
	}

	// Clean up the message for EDAN H30
	// Convert literal \n to actual newlines
	message = strings.ReplaceAll(message, "\\n", "\n")

	// Log the cleaned message
	logMsg := strings.ReplaceAll(message, "\r", "\n")
	slog.Info("EDAN H30: Cleaned message received", "message", logMsg)

	msgByte := []byte(message)
	headerDecoder := hl7.NewDecoder(h251.Registry, &hl7.DecodeOption{HeaderOnly: true})
	header, err := headerDecoder.Decode(msgByte)
	if err != nil {
		slog.Error("EDAN H30: Header decode failed", "error", err, "message_length", len(message))
		return "", fmt.Errorf("decode header failed: %w", err)
	}

	switch m := header.(type) {
	case h251.ORU_R01:
		slog.Info("EDAN H30: Processing ORU_R01 message")
		return h.ORUR01(ctx, m, msgByte)
	}

	slog.Error("EDAN H30: Unknown message type", "type", fmt.Sprintf("%T", header))
	return "", fmt.Errorf("unknown message type %T", header)
}
