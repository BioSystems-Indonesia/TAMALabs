package ncc3300

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"github.com/tarm/serial"
)

type Handler struct {
	AnalyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		AnalyzerUsecase: analyzerUsecase,
	}
}

var rawHl7 string

func (h *Handler) Handle(port *serial.Port) {
	buf := make([]byte, 1024)
	var timer *time.Timer
	for {
		n, err := port.Read(buf)
		if err != nil || n == 0 {
			slog.Error("Error reading serial data", "error", err)
			continue
		}

		buffer := buf[:n]
		for _, b := range buffer {
			if b >= 32 && b <= 126 {
				rawHl7 += string(b)
			}
		}

		// Reset timer setiap ada data baru
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(300*time.Millisecond, func() {
			if rawHl7 != "" {
				h.handleMessage(context.Background(), rawHl7)
				rawHl7 = ""
			}
		})
	}
}

func (h *Handler) handleMessage(ctx context.Context, message string) (string, error) {
	// don't do anything if the message is empty
	if message == "" {
		return "", nil
	}

	// Tambahkan newline sebelum setiap 'OBX'
	// Tambahkan newline sebelum setiap 'OBX', 'OBR', dan 'MSH'
	rawWithNewline := message
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "|OBX", "\nOBX")
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "|OBR", "\nOBR")
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "MSH|", "\nMSH|")
	// Pisahkan berdasarkan newline
	resultArray := []string{}
	lines := strings.Split(rawWithNewline, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OBX") && !strings.Contains(line, "Histogram") {
			resultArray = append(resultArray, line)
		} else if strings.HasPrefix(line, "MSH") || strings.HasPrefix(line, "OBR") {
			resultArray = append(resultArray, line)
		}
	}

	msgByte := []byte{}

	for _, result := range resultArray {
		msgByte = append(msgByte, []byte(result)...)
	}

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
