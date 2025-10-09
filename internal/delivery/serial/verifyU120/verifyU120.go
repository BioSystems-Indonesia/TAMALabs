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
	reID     = regexp.MustCompile(`ID\s*:\s*(\d+)`)                    // Matches " ID:251005004"
	reDate   = regexp.MustCompile(`Date\s*:\s*(.+)`)                   // Matches " Date:05-10-2025 15:01"
	reResult = regexp.MustCompile(`^\s*\*?\s*([A-Z]{2,4}|pH)\s+(.+)$`) // Even more flexible pattern
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
	slog.Info("Read serial data",
		"bytesRead", n,
		"data", data,
		"currentBufferLength", len(h.buffer))

	h.buffer += data

	slog.Info("Buffer after append",
		"newBufferLength", len(h.buffer),
		"bufferContent", h.buffer)

	// Parse immediately but use deduplication logic to prevent duplicates
	h.parseUrineResult()
}

func (h *Handler) parseUrineResult() {
	// Check if we have received the complete message (ending with ETX)
	if !strings.Contains(h.buffer, "\u0003") {
		slog.Debug("Waiting for complete message (ETX not found)")
		return
	}

	// Clean buffer from control characters
	cleanBuffer := strings.ReplaceAll(h.buffer, "\u0002", "")   // Remove STX
	cleanBuffer = strings.ReplaceAll(cleanBuffer, "\u0003", "") // Remove ETX
	cleanBuffer = strings.ReplaceAll(cleanBuffer, "\r", "")     // Remove CR, keep LF

	// Check if buffer contains a complete result by scanning lines
	hasID := false
	hasResult := false

	scanner := bufio.NewScanner(strings.NewReader(cleanBuffer))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if reID.MatchString(line) {
			hasID = true
		}
		if reResult.MatchString(line) {
			hasResult = true
		}

		// Early exit if we found both
		if hasID && hasResult {
			break
		}
	}

	slog.Info("parseUrineResult called",
		"bufferLength", len(h.buffer),
		"cleanBufferLength", len(cleanBuffer),
		"hasID", hasID,
		"hasResult", hasResult)

	// Only proceed if we have both ID and at least one result
	if !hasID || !hasResult {
		slog.Debug("Buffer incomplete, waiting for more data")
		return
	}

	scanner = bufio.NewScanner(strings.NewReader(cleanBuffer))

	var patientID string
	var ts time.Time
	var results []entity.VerifyResult // Collect all results for batch processing

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		slog.Debug("Processing line", "line", line)

		// Parse patient ID
		if m := reID.FindStringSubmatch(line); len(m) > 1 {
			patientID = m[1]
			slog.Info("Found patient ID", "patientID", patientID)
			continue
		}

		// Parse date
		if m := reDate.FindStringSubmatch(line); len(m) > 1 {
			dateStr := strings.TrimSpace(m[1])
			// Try different date formats
			layouts := []string{
				"02-01-2006 15:04",
				"02-01-2006 15:04 pm",
				"01-02-2006 15:04",
				"01-02-2006 15:04 pm",
			}
			for _, layout := range layouts {
				if parsed, err := time.Parse(layout, dateStr); err == nil {
					ts = parsed
					slog.Info("Found timestamp", "timestamp", ts, "dateStr", dateStr)
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
			slog.Debug("Skipping non-result line", "line", line)
			continue
		}

		// Parse result lines using regex
		if m := reResult.FindStringSubmatch(line); len(m) > 1 {
			testName := m[1]
			resultPart := strings.TrimSpace(m[2])

			value, valueStr := h.parseResultValueOnly(resultPart)

			slog.Info("Processing urine result",
				"testName", testName,
				"patientID", patientID,
				"originalResultPart", resultPart,
				"extractedValueStr", valueStr,
				"extractedValue", value,
				"originalLine", line)

			// Add to results collection instead of processing immediately
			if valueStr == "-" {
				valueStr = "NEGATIF"
			} else if valueStr == "+" {
				valueStr = "NORMAL (+)"
			}
			results = append(results, entity.VerifyResult{
				PatientID:  patientID,
				TestName:   testName,
				SampleType: "URI",
				Value:      value,
				ValueStr:   valueStr,
				Unit:       "",
				Timestamp:  ts,
			})
		} else {
			slog.Debug("Line did not match result regex", "line", line)
		}
	}

	// Process all results as a batch
	if len(results) > 0 {
		err := h.analyzerUseCase.ProcessVerifyU120Batch(context.Background(), results)
		if err != nil {
			slog.Error("Failed to process VerifyU120 batch", "error", err, "resultCount", len(results))
		}
	}

	// Clear buffer after processing to prevent reprocessing
	h.buffer = ""
	slog.Info("Finished processing urine results",
		"patientID", patientID,
		"processedCount", len(results),
		"timestamp", ts)
}

// parseResultValueOnly parses the result part and extracts only the value (no unit)
func (h *Handler) parseResultValueOnly(resultPart string) (float64, string) {
	// Clean up the result part
	resultPart = strings.TrimSpace(resultPart)

	// Split by whitespace to analyze components
	parts := strings.Fields(resultPart)
	if len(parts) == 0 {
		return 0, ""
	}

	// Case 1: Single "-" means negative
	if len(parts) == 1 && parts[0] == "-" {
		return 0, "-"
	}

	// Case 2: Single value (e.g., "6.0" for pH, "1.030" for SG)
	if len(parts) == 1 {
		if val, err := strconv.ParseFloat(parts[0], 64); err == nil {
			return val, parts[0]
		}
		// Non-numeric single value (like qualitative results)
		return 0, parts[0]
	}

	// Case 3: First part is qualitative with + (like "+", "+-")
	if strings.Contains(parts[0], "+") {
		// For qualitative results like "+", "+-", use the qualitative value as the main value
		return 0, parts[0]
	}

	// Case 4: First part is "-" followed by anything (like "- 3.5 umol/L" or "- neg")
	// This means the result is NEGATIVE regardless of the numeric value that follows
	if parts[0] == "-" {
		// Always return "-" as the value for negative results
		return 0, "-"
	}

	// Case 5: Numeric value first (like "3.5 umol/L") - but only if not preceded by "-"
	if val, err := strconv.ParseFloat(parts[0], 64); err == nil {
		return val, parts[0]
	}

	// Case 6: Return first part as string value
	return 0, parts[0]
}
