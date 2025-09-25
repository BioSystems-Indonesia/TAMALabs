package diestro

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"go.bug.st/serial"
)

var (
	reDelim       = regexp.MustCompile(`={5,}`) // delimiter flexible: =====
	reNamedEq     = regexp.MustCompile(`(?i)^([A-Za-z]{1,8})\s*=\s*([0-9]+(?:\.[0-9]+)?)\s*([a-zA-Zμ%/]+)$`)
	reNamedPrefix = regexp.MustCompile(`(?i)^([A-Za-z]{1,8})\s+([0-9]+(?:\.[0-9]+)?)\s*([a-zA-Zμ%/]+)$`)
	reValueLine   = regexp.MustCompile(`^([0-9]+(?:\.[0-9]+)?)\s*([a-zA-Zμ%/]+)$`)
	reDate        = regexp.MustCompile(`\d{4}/\d{2}/\d{2}`)
	reTime        = regexp.MustCompile(`\d{2}:\d{2}:\d{2}`)
	reOnlyDigits  = regexp.MustCompile(`^\d{3,10}$`)
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

	h.processBuffer()
}

func (h *Handler) processBuffer() {
	locs := reDelim.FindAllStringIndex(h.buffer, -1)
	if len(locs) == 0 {
		return
	}

	lastEnd := locs[len(locs)-1][1]
	reportsStr := h.buffer[:lastEnd]
	remainder := h.buffer[lastEnd:]
	parts := reDelim.Split(reportsStr, -1)

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		cleaned := cleanReport(p)
		results, err := parseReportTolerant(cleaned)
		if err != nil {
			slog.Error("parse error", "err", err)
			continue
		}
		for _, r := range results {
			fmt.Printf("Parsed result: %+v\n", r)
			h.analyzerUseCase.ProcessDiestro(context.Background(), r)
		}
	}

	h.buffer = remainder
}

func cleanReport(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	reRuns := regexp.MustCompile(`(?:=|-){3,}`)
	s = reRuns.ReplaceAllString(s, "\n")

	var b strings.Builder
	for _, r := range s {
		if r == '\n' || !unicode.IsControl(r) {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	s = b.String()

	reBad := regexp.MustCompile(`[^\p{L}\p{N}\s=./:%μ%+\-()]`)
	s = reBad.ReplaceAllString(s, " ")

	lines := strings.Split(s, "\n")
	for i, L := range lines {
		lines[i] = strings.TrimSpace(regexp.MustCompile(`\s+`).ReplaceAllString(L, " "))
	}
	return strings.Join(lines, "\n")
}

func parseReportTolerant(cleaned string) ([]entity.DiestroResult, error) {
	lines := strings.Split(cleaned, "\n")

	headerTests := parseHeader(lines)

	patientID := ""
	for i, l := range lines {
		low := strings.ToLower(l)
		if strings.HasPrefix(low, "patient") {
			if idx := strings.Index(l, ":"); idx >= 0 {
				cand := strings.TrimSpace(l[idx+1:])
				if reOnlyDigits.MatchString(cand) {
					patientID = cand
					break
				}
			}
			if i+1 < len(lines) {
				cand := strings.TrimSpace(lines[i+1])
				if reOnlyDigits.MatchString(cand) {
					patientID = cand
					break
				}
			}
		}
	}
	if patientID == "" {
		for _, l := range lines {
			ll := strings.ToLower(strings.TrimSpace(l))
			if strings.HasPrefix(ll, "mem") {
				continue
			}
			if reOnlyDigits.MatchString(strings.TrimSpace(l)) {
				patientID = strings.TrimSpace(l)
				break
			}
		}
	}

	var ts time.Time
	foundDate := ""
	for i, l := range lines {
		if reDate.MatchString(l) {
			foundDate = reDate.FindString(l)
			var timeStr string
			if reTime.MatchString(l) {
				timeStr = reTime.FindString(l)
			} else {
				reShort := regexp.MustCompile(`\d{2}:\d{2}`)
				if reShort.MatchString(l) {
					timeStr = reShort.FindString(l) + ":00"
				} else if i+1 < len(lines) && reShort.MatchString(lines[i+1]) {
					timeStr = reShort.FindString(lines[i+1]) + ":00"
				}
			}
			if timeStr == "" {
				timeStr = "00:00:00"
			}
			parseStr := foundDate + " " + timeStr
			t, err := time.Parse("2006/01/02 15:04:05", parseStr)
			if err == nil {
				ts = t
			}
			break
		}
	}
	if foundDate == "" {
		ts = time.Time{}
	}

	results := []entity.DiestroResult{}
	foundNames := map[string]bool{}

	type unnamedVal struct {
		Val  float64
		Unit string
	}
	var unnamed []unnamedVal

	// helper to skip obviously-non-measure lines
	skipIfContains := func(s string) bool {
		low := strings.ToLower(s)
		checks := []string{"analyzer", "laboratorio", "electrolyte", "report", "mem", "patient", "sn:", "====", "-----"}
		for _, c := range checks {
			if strings.Contains(low, c) {
				return true
			}
		}
		return false
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if skipIfContains(line) {
			continue
		}

		// strict named with "=": e.g. "Na=141.7mmol" or "Cl = 106.4 mmol"
		if m := reNamedEq.FindStringSubmatch(line); m != nil {
			val, err := strconv.ParseFloat(m[2], 64)
			if err != nil {
				continue
			}
			name := strings.TrimSpace(m[1])
			unit := strings.TrimSpace(m[3])
			results = append(results, entity.DiestroResult{
				PatientID:  patientID,
				TestName:   name,
				SampleType: "SER",
				Value:      val,
				Unit:       unit,
				Timestamp:  ts,
			})
			foundNames[name] = true
			continue
		}

		// alternative named: "Na 141.7 mmol"
		if m := reNamedPrefix.FindStringSubmatch(line); m != nil {
			val, err := strconv.ParseFloat(m[2], 64)
			if err != nil {
				continue
			}
			name := strings.ToUpper(strings.TrimSpace(m[1]))
			unit := strings.TrimSpace(m[3])
			results = append(results, entity.DiestroResult{
				PatientID:  patientID,
				TestName:   name,
				SampleType: "SER",
				Value:      val,
				Unit:       unit,
				Timestamp:  ts,
			})
			foundNames[name] = true
			continue
		}

		// value-only whole-line: e.g. "106.4mmol" or "106.4 mmol"
		if m := reValueLine.FindStringSubmatch(line); m != nil {
			val, err := strconv.ParseFloat(m[1], 64)
			if err != nil {
				continue
			}
			unit := strings.TrimSpace(m[2])
			unnamed = append(unnamed, unnamedVal{Val: val, Unit: unit})
			continue
		}

		// otherwise ignore (likely SN, version, mixed tokens)
	}

	// assign unnamed values to headerTests order (exclude already-found names)
	remainingNames := []string{}
	for _, h := range headerTests {
		H := strings.ToUpper(strings.TrimSpace(h))
		if H == "" {
			continue
		}
		if !foundNames[H] {
			remainingNames = append(remainingNames, H)
		}
	}

	for i, u := range unnamed {
		var name string
		if i < len(remainingNames) {
			name = remainingNames[i]
		} else {
			name = fmt.Sprintf("ANON%02d", i+1)
		}
		results = append(results, entity.DiestroResult{
			PatientID:  patientID,
			TestName:   name,
			SampleType: "SER",
			Value:      u.Val,
			Unit:       u.Unit,
			Timestamp:  ts,
		})
	}

	return results, nil
}

func parseHeader(lines []string) []string {
	for _, l := range lines {
		if strings.Contains(l, "*") {
			parts := strings.Split(l, "*")
			var out []string
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "" {
					continue
				}
				// only letters/digits prefix
				sb := strings.Builder{}
				for _, r := range p {
					if unicode.IsLetter(r) || unicode.IsDigit(r) {
						sb.WriteRune(r)
					} else {
						break
					}
				}
				if s := sb.String(); s != "" {
					out = append(out, s)
				}
			}
			if len(out) > 0 {
				return out
			}
		}
	}
	for _, l := range lines {
		words := strings.Fields(l)
		if len(words) >= 2 && len(words) <= 6 {
			ok := true
			for _, w := range words {
				w2 := strings.Trim(w, " *")
				if len(w2) > 6 || len(w2) < 1 {
					ok = false
					break
				}
			}
			if ok {
				return words
			}
		}
	}
	return nil
}
