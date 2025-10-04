package verifyu120

import (
	"bufio"
	"context"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"go.bug.st/serial"
)

var (
	reID     = regexp.MustCompile(`ID\s*:\s*(\d+)`)
	reDate   = regexp.MustCompile(`Date\s*:\s*(.+)`)
	reResult = regexp.MustCompile(`^\s*\*?([A-Z]{2,4}|pH)\s+(.+)$`) // Allow leading spaces, matches parameter lines with optional *
)

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
	slog.Debug("Read serial data", "n", n)
	h.buffer += data

	h.parseUrineResult()
}

func (h *Handler) parseUrineResult() {
	scanner := bufio.NewScanner(strings.NewReader(h.buffer))

	var patientID string
	var ts time.Time

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Parse patient ID
		if m := reID.FindStringSubmatch(line); len(m) > 1 {
			patientID = m[1]
			continue
		}

		// Parse date
		if m := reDate.FindStringSubmatch(line); len(m) > 1 {
			dateStr := strings.TrimSpace(m[1])
			// Try different date formats
			layouts := []string{
				"02-01-2006 15:04 pm",
				"02-01-2006 15:04",
				"01-02-2006 15:04 pm",
				"01-02-2006 15:04",
			}
			for _, layout := range layouts {
				if parsed, err := time.Parse(layout, dateStr); err == nil {
					ts = parsed
					break
				}
			}
			continue
		}

		// Skip non-result lines
		if strings.HasPrefix(line, "Operator:") ||
			strings.HasPrefix(line, "No.") ||
			strings.Contains(line, "param |") ||
			strings.Contains(line, "hasil |") {
			continue
		}

		// Parse result lines using regex
		if m := reResult.FindStringSubmatch(line); len(m) > 1 {
			testName := m[1]
			resultPart := strings.TrimSpace(m[2])

			value, valueStr, unit := h.parseResultValue(resultPart)

			slog.Debug("Parsed urine result",
				"testName", testName,
				"valueStr", valueStr,
				"value", value,
				"unit", unit,
				"originalLine", line)

			h.analyzerUseCase.ProcessVerifyU120(context.Background(), entity.VerifyResult{
				PatientID:  patientID,
				TestName:   testName,
				SampleType: "URI",
				Value:      value,
				ValueStr:   valueStr, // Store original string value
				Unit:       unit,
				Timestamp:  ts,
			})
		}
	}
}

// parseResultValue parses the result part and extracts value, original string, and unit
func (h *Handler) parseResultValue(resultPart string) (float64, string, string) {
	// Clean up the result part
	resultPart = strings.TrimSpace(resultPart)

	// Split by whitespace to analyze components
	parts := strings.Fields(resultPart)
	if len(parts) == 0 {
		return 0, "", ""
	}

	// Handle cases where first part is "-" but there are numeric values after
	if parts[0] == "-" && len(parts) > 1 {
		// Look for numeric value after the "-"
		for i := 1; i < len(parts); i++ {
			if val, err := strconv.ParseFloat(parts[i], 64); err == nil {
				unit := ""
				if i+1 < len(parts) && parts[i+1] != "neg" {
					unit = parts[i+1]
				}
				return val, parts[i], unit
			}
		}
	}

	// Handle purely negative results
	if parts[0] == "-" && (len(parts) == 1 || (len(parts) > 1 && parts[len(parts)-1] == "neg")) {
		return 0, "neg", ""
	}

	// Case 1: Single value (e.g., "6.0" for pH)
	if len(parts) == 1 {
		if val, err := strconv.ParseFloat(parts[0], 64); err == nil {
			return val, parts[0], ""
		}
		// Non-numeric single value
		return 0, parts[0], ""
	}

	// Case 2: Two parts - could be "1+" "0.3" or "3.5" "umol/L"
	if len(parts) == 2 {
		// Check if first part is qualitative (contains +, -, etc.)
		if strings.Contains(parts[0], "+") || strings.Contains(parts[0], "-") {
			// Qualitative result like "1+", try to parse second part as numeric
			if val, err := strconv.ParseFloat(parts[1], 64); err == nil {
				return val, parts[0], "" // Use qualitative as string value
			}
			return 0, parts[0], ""
		}

		// Check if second part is unit
		if val, err := strconv.ParseFloat(parts[0], 64); err == nil {
			return val, parts[0], parts[1]
		}

		// Both non-numeric
		return 0, strings.Join(parts, " "), ""
	}

	// Case 3: Three or more parts (e.g., "1+" "0.3" "g/L" or "3.5" "umol/L")
	if len(parts) >= 3 {
		// Try to find numeric value and unit
		for i, part := range parts {
			if val, err := strconv.ParseFloat(part, 64); err == nil {
				// Found numeric value, check for unit
				unit := ""
				if i+1 < len(parts) && parts[i+1] != "neg" {
					unit = parts[i+1]
				}
				// Use first part as string representation if it's qualitative
				valueStr := part
				if i > 0 && (strings.Contains(parts[0], "+") || strings.Contains(parts[0], "-")) {
					valueStr = parts[0]
				}
				return val, valueStr, unit
			}
		}

		// No numeric value found, treat as qualitative
		return 0, parts[0], ""
	}

	return 0, resultPart, ""
}
