package response911

import (
	"fmt"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"log/slog"
)

// ParseResponse911 parses an assembled Response911 (ASTM-like) message and
// converts it into entity.ORU_R01 containing barcode and observation results.
// Expected input: records separated by CR ("\r") or LF. Record examples:
// 4O|1|<barcode>|...
// 5R|1|^^^^UA|4.21|mg/dL|...
// 4L|1|N
func ParseResponse911(raw string) (*entity.ORU_R01, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("empty message")
	}

	// Normalize separators and split into lines
	norm := strings.ReplaceAll(raw, "\r\n", "\r")
	norm = strings.ReplaceAll(norm, "\n", "\r")
	lines := strings.Split(norm, "\r")

	var barcode string
	var results []entity.ObservationResult

	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}

		parts := strings.Split(ln, "|")
		if len(parts) == 0 {
			continue
		}

		recordID := parts[0] // e.g. "4O", "5R", "4L"
		if len(recordID) < 2 {
			continue
		}

		recType := recordID[1]
		switch recType {
		case 'O':
			// Order record - sample/barcode usually in field 3 (index 2)
			if len(parts) > 2 && barcode == "" {
				barcode = strings.TrimSpace(parts[2])
			}
		case 'R':
			// Result record - test id in field 3, value in field 4, unit in field 5
			if len(parts) > 3 {
				testField := parts[2]
				valueField := parts[3]
				unit := ""
				if len(parts) > 4 {
					unit = parts[4]
				}

				// extract test code (use last segment after '^')
				code := ""
				if testField != "" {
					segs := strings.Split(testField, "^")
					for i := len(segs) - 1; i >= 0; i-- {
						if segs[i] != "" {
							code = segs[i]
							break
						}
					}
				}

				if code == "" {
					// fallback: use whole testField
					code = testField
				}

				// sanitize value
				valueStr := strings.TrimSpace(valueField)

				or := entity.ObservationResult{
					TestCode:       code,
					Description:    "",
					Values:         entity.JSONStringArray{valueStr},
					Type:           code,
					Unit:           unit,
					ReferenceRange: "",
					Date:           time.Now(),
					AbnormalFlag:   entity.JSONStringArray{},
					Comments:       "",
					Picked:         false,
				}
				results = append(results, or)
			}
		default:
			// ignore other records (C, Z, L handled elsewhere)
		}
	}

	if barcode == "" {
		slog.Warn("Response911 parser: barcode not found in message", "raw", raw)
		// still return parsed results without barcode
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found in message")
	}

	specimen := entity.Specimen{
		ReceivedDate:      time.Now(),
		Barcode:           barcode,
		ObservationResult: results,
	}

	oru := entity.ORU_R01{
		MSH: entity.MSH{},
		Patient: []entity.Patient{
			{
				Specimen: []entity.Specimen{specimen},
			},
		},
	}

	return &oru, nil
}
