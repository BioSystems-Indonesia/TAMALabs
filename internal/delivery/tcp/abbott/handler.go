package abbott

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/panics"
)

type Handler struct {
	analyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		analyzerUsecase: analyzerUsecase,
	}
}

// Handle receives data from Abbott device and processes it
func (h *Handler) Handle(conn *net.TCPConn) {
	ctx := context.Background()

	defer panics.RecoverPanic(ctx)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	buffer := make([]byte, 4096)

	for {
		n, err := reader.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				slog.Info("Abbott connection closed")
				break
			}
			slog.Error("error reading Abbott message", "error", err)
			return
		}

		if n == 0 {
			continue
		}

		// Get the actual data received
		message := buffer[:n]
		rawData := string(message)

		// Display received data (replace \r with \n for better readability)
		displayMessage := strings.ReplaceAll(rawData, "\r", "\n")
		slog.Info(fmt.Sprintf("Abbott received data (%d bytes):\n%s", n, displayMessage))

		// Parse Abbott data
		abbottMessage, err := ParseAbbottData(rawData)
		if err != nil {
			slog.Error("error parsing Abbott data", "error", err)
		} else if abbottMessage != nil {
			// Process the parsed data
			err = h.analyzerUsecase.ProcessAbbott(ctx, *abbottMessage)
			if err != nil {
				slog.Error("error processing Abbott data", "error", err)
			} else {
				slog.Info("Successfully processed Abbott data",
					"patient_id", abbottMessage.SampleInfo.PatientID,
					"sample_id", abbottMessage.SampleInfo.SampleID,
					"test_count", len(abbottMessage.TestResults))
			}
		}

		// Send ACK response (simple acknowledgment)
		ack := []byte{0x06} // ACK character
		if _, err := conn.Write(ack); err != nil {
			slog.Error("error sending ACK to Abbott", "error", err)
		}
	}
}
