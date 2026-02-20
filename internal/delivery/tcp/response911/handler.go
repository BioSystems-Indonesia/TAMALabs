package response911

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/panics"
)

// Handler receives ASTM-like frames from Response 911 devices, assembles
// fragmented frames into a complete message, parses barcode/test/value and
// forwards results to usecase.ProcessORUR01.
// Connection behavior mirrors Abbott (plain TCP, no MLLP); an ACK (0x06) is sent back.
type Handler struct {
	analyzerUsecase usecase.Analyzer
	agg             *aggregator
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{analyzerUsecase: analyzerUsecase, agg: newAggregator(analyzerUsecase.ProcessORUR01, 1*time.Second)}
}

func (h *Handler) Handle(conn *net.TCPConn) {
	ctx := context.Background()
	defer panics.RecoverPanic(ctx)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buf := make([]byte, 4096)

	// connBuf holds raw bytes from the socket (may contain partial frames)
	var connBuf string
	// msgBuilder accumulates frame payloads (without STX/ETX/checksum) until an L record is seen
	var msgBuilder strings.Builder

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Response911 connection closed")
				break
			}
			slog.Error("error reading Response911 message", "error", err)
			return
		}

		if n == 0 {
			continue
		}

		chunk := string(buf[:n])
		connBuf += chunk

		// Log the raw chunk for debugging (replace CR with newline for readability)
		display := strings.ReplaceAll(chunk, "\r", "\n")
		slog.Info(fmt.Sprintf("Response911 received data (%d bytes):\n%s", n, display))

		// Extract complete frames (STX ... ETX). Each frame typically: STX + payload + ETX + checksum + CR
		for {
			stx := strings.Index(connBuf, "\x02")
			if stx == -1 {
				break
			}
			etxRel := strings.Index(connBuf[stx+1:], "\x03")
			if etxRel == -1 {
				break // wait for more data
			}
			etx := stx + 1 + etxRel

			// payload between STX and ETX
			frame := connBuf[stx+1 : etx]
			// append to assembled message
			msgBuilder.WriteString(frame)
			// ensure records are separated by CR
			if !strings.HasSuffix(frame, "\r") {
				msgBuilder.WriteString("\r")
			}

			// remove processed portion (we skip checksum+CR by removing only up to ETX)
			connBuf = connBuf[etx+1:]

			// Check if this frame contains an L record (terminator)
			// Look for lines where second character is 'L' (e.g. "4L|1|N")
			assembled := msgBuilder.String()
			lines := strings.Split(assembled, "\r")
			messageComplete := false
			for _, ln := range lines {
				ln = strings.TrimSpace(ln)
				if ln == "" {
					continue
				}
				if len(ln) >= 2 && ln[1] == 'L' {
					messageComplete = true
					break
				}
			}

			if messageComplete {
				rawMsg := strings.Trim(msgBuilder.String(), "\r\n")
				// Parse assembled ASTM-like message
				oru, err := ParseResponse911(rawMsg)
				if err != nil {
					slog.Error("error parsing Response911 message", "error", err, "raw", rawMsg)
				} else {
					// aggregate results per-barcode and flush after debounce timeout
					h.agg.Add(oru)

					// log extracted summary (barcode, first test)
					if len(oru.Patient) > 0 && len(oru.Patient[0].Specimen) > 0 && len(oru.Patient[0].Specimen[0].ObservationResult) > 0 {
						s := oru.Patient[0].Specimen[0]
						tr := s.ObservationResult[0]
						slog.Info("Response911 parsed and queued",
							"barcode", s.Barcode,
							"test_code", tr.TestCode,
							"value", tr.Values[0],
						)
					}
				}

				// reset builder for next message
				msgBuilder.Reset()
			}
		}

		// Send ACK (keep same behavior as Abbott)
		ack := []byte{0x06}
		if _, err := conn.Write(ack); err != nil {
			slog.Error("error sending ACK to Response911", "error", err)
		}
	}
}
