package coax

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"go.bug.st/serial"
)

type Handler struct {
	analyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		analyzerUsecase: analyzerUsecase,
	}
}

func (h *Handler) Handle(port serial.Port) {
	buf := make([]byte, 1024)
	var buffer string // Accumulate partial lines

	n, err := port.Read(buf)
	if err != nil {
		slog.Error("Error reading serial data", "error", err)
		return
	}

	slog.Debug("Read serial data", "n", n)
	if n == 0 {
		return
	}

	data := string(buf[:n])
	buffer += data

	for {
		newlineIndex := strings.Index(buffer, "\n")
		if newlineIndex < 0 {
			break
		}

		line := strings.TrimSpace(buffer[:newlineIndex])
		buffer = buffer[newlineIndex+1:]

		if line == "" {
			continue
		}

		result, err := parseLine(line)
		if err != nil {
			slog.Error("Parse error", "error", err, "line", line)
			continue
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			slog.Error("JSON error", "error", err)
			continue
		}

		slog.Info("Parsed result coax", "result", string(jsonData))

		err = h.analyzerUsecase.ProcessCoax(context.Background(), result)
		if err != nil {
			slog.Error("Error processing coax test result", "error", err)
			continue
		}
	}
}

func parseLine(line string) (entity.CoaxTestResult, error) {
	// Remove control characters and trim spaces
	cleaned := strings.Trim(line, "\x02\x03")
	fields := strings.Fields(cleaned)

	// Debug: Print cleaned fields
	slog.Info("Cleaned fields", "fields", fields)

	// Safely handle fields (some lines might have fewer columns)
	result := entity.CoaxTestResult{}
	if len(fields) > 0 {
		result.RecordType = fields[0] // R
	}
	if len(fields) > 1 {
		result.DeviceID = fields[1] // 0202001465
	}
	if len(fields) > 2 {
		result.Status = fields[2] // 1
	}
	if len(fields) > 3 {
		result.Date = fields[3] // 2025/05/28
	}
	if len(fields) > 4 {
		result.Time = fields[4] // 15:46
	}
	if len(fields) > 5 {
		result.TestType = fields[5] // 1
	}
	if len(fields) > 6 {
		result.TestName = fields[7] // APTT (skip fields[6] which is "1")
	}
	if len(fields) > 7 {
		result.Value = fields[8] // 69.7
	}
	if len(fields) > 8 {
		result.Unit = fields[9] // s
	}
	if len(fields) > 9 {
		result.Reference = fields[10] // 57
	}
	if len(fields) > 10 {
		result.Flags = fields[11] // NR
	}
	if len(fields) > 11 {
		result.Extra = fields[12:] // [24]
	}

	return result, nil
}
