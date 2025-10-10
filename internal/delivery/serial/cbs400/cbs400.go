package cbs400

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"go.bug.st/serial"
)

var (
	// CBS400 specific patterns
	reResultLine = regexp.MustCompile(`^r(\d+)\s+([\d\s.]+)$`) // matches "r371                   3.91 141.2 106.5  1.21  1.23  2.46  7.44"
	reNumbers    = regexp.MustCompile(`\d+(?:\.\d+)?`)         // extracts decimal numbers
)

// CBS400 parameter order: K, Na, Cl, iCa, nCa, TCa, pH
var cbs400Parameters = []struct {
	Name string
	Unit string
}{
	{"K", "mmol/L"},
	{"Na", "mmol/L"},
	{"Cl", "mmol/L"},
	{"iCa", "mmol/L"},
	{"nCa", "mmol/L"},
	{"TCa", "mmol/L"},
	{"pH", ""},
}

type Handler struct {
	analyzerUseCase usecase.Analyzer
	buffer          string
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{analyzerUseCase: analyzerUsecase}
}

func (h *Handler) Handle(port serial.Port) {
	buf := make([]byte, 1024)
	n, err := port.Read(buf)
	if err != nil {
		slog.Error("Error reading serial data", "error", err)
		return
	}

	data := string(buf[:n])
	slog.Debug("Read CBS400 serial data", "n", n, "data", data)
	h.buffer += data

	h.processBuffer()
}

func (h *Handler) processBuffer() {
	lines := strings.Split(h.buffer, "\n")

	// Keep the last incomplete line in buffer
	if len(lines) > 0 && !strings.HasSuffix(h.buffer, "\n") {
		h.buffer = lines[len(lines)-1]
		lines = lines[:len(lines)-1]
	} else {
		h.buffer = ""
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		results, err := parseCBS400Line(line)
		if err != nil {
			slog.Debug("Parse error", "line", line, "error", err)
			continue
		}

		for _, result := range results {
			fmt.Printf("Parsed CBS400 result: %+v\n", result)
			h.analyzerUseCase.ProcessCBS400(context.Background(), result)
		}
	}
}

func parseCBS400Line(line string) ([]entity.CBS400Result, error) {
	// Check if line matches CBS400 result pattern
	matches := reResultLine.FindStringSubmatch(line)
	if len(matches) != 3 {
		return nil, fmt.Errorf("line does not match CBS400 result pattern: %s", line)
	}

	patientID := matches[1]
	valuesStr := matches[2]

	// Extract all numeric values
	numberStrings := reNumbers.FindAllString(valuesStr, -1)
	if len(numberStrings) < len(cbs400Parameters) {
		return nil, fmt.Errorf("insufficient values found: expected %d, got %d", len(cbs400Parameters), len(numberStrings))
	}

	var results []entity.CBS400Result
	timestamp := time.Now()

	// Map values to parameters in order: K, Na, Cl, iCa, nCa, TCa, pH
	for i, param := range cbs400Parameters {
		if i >= len(numberStrings) {
			break
		}

		value, err := strconv.ParseFloat(numberStrings[i], 64)
		if err != nil {
			slog.Warn("Failed to parse value", "value", numberStrings[i], "parameter", param.Name)
			continue
		}

		// Basic validation for reasonable ranges
		if isValidCBS400Value(param.Name, value) {
			result := entity.CBS400Result{
				PatientID:  patientID,
				TestName:   param.Name,
				SampleType: "SER", // Serum sample type
				Value:      value,
				Unit:       param.Unit,
				Timestamp:  timestamp,
			}
			results = append(results, result)
		} else {
			slog.Warn("Value out of expected range", "parameter", param.Name, "value", value)
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no valid results parsed from line: %s", line)
	}

	return results, nil
}

// isValidCBS400Value performs basic validation on CBS400 parameter values
func isValidCBS400Value(paramName string, value float64) bool {
	switch paramName {
	case "K": // Potassium
		return value >= 1.0 && value <= 10.0
	case "Na": // Sodium
		return value >= 100.0 && value <= 200.0
	case "Cl": // Chloride
		return value >= 50.0 && value <= 150.0
	case "iCa": // Ionized Calcium
		return value >= 0.5 && value <= 3.0
	case "nCa": // Normalized Calcium
		return value >= 0.5 && value <= 3.0
	case "TCa": // Total Calcium
		return value >= 1.0 && value <= 5.0
	case "pH": // pH
		return value >= 6.0 && value <= 8.5
	default:
		// Unknown parameter, allow it through
		return true
	}
}
