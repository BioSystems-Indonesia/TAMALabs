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
	reID   = regexp.MustCompile(`ID\s*:\s*(\d+)`)
	reDate = regexp.MustCompile(`Date\s*:\s*(.+)`)
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

		if m := reID.FindStringSubmatch(line); len(m) > 1 {
			patientID = m[1]
			continue
		}

		if m := reDate.FindStringSubmatch(line); len(m) > 1 {
			dateStr := strings.TrimSpace(m[1])
			parsed, err := time.Parse("02-01-2006 15:04", dateStr)
			if err == nil {
				ts = parsed
			}
			continue
		}

		if strings.HasPrefix(line, "Operator:") || strings.HasPrefix(line, "No.") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		testName := fields[0]
		var valueStr, unit string
		var value float64

		if len(fields) == 2 {
			valueStr = fields[1]
		} else if len(fields) == 3 {
			valueStr = fields[2]
		} else if len(fields) == 4 {
			valueStr = fields[2]
			unit = fields[3]
		}

		if v, err := strconv.ParseFloat(valueStr, 64); err == nil {
			value = v
		} else {
			value = 0
		}

		h.analyzerUseCase.ProcessVerifyU120(context.Background(), entity.VerifyResult{
			PatientID:  patientID,
			TestName:   testName,
			SampleType: "URI",
			Value:      value,
			Unit:       unit,
			Timestamp:  ts,
		})
	}
}
