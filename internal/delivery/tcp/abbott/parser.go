package abbott

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

// ParseAbbottData parses raw Abbott device data
func ParseAbbottData(rawData string) (*entity.AbbottMessage, error) {
	// Replace \r\n or \r with \n for consistent line splitting
	rawData = strings.ReplaceAll(rawData, "\r\n", "\n")
	rawData = strings.ReplaceAll(rawData, "\r", "\n")
	
	lines := strings.Split(rawData, "\n")

	message := &entity.AbbottMessage{
		DeviceInfo:  entity.AbbottDeviceInfo{},
		SampleInfo:  entity.AbbottSampleInfo{},
		TestResults: []entity.AbbottTestResult{},
		Timestamp:   time.Now(),
	}

	var date, timeStr string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "RESULT" {
			continue
		}

		parts := strings.Split(line, ";")
		if len(parts) == 0 {
			continue
		}

		key := strings.TrimSpace(parts[0])

		switch key {
		case "DATE":
			if len(parts) > 1 {
				date = strings.TrimSpace(parts[1])
			}
		case "TIME":
			if len(parts) > 1 {
				timeStr = strings.TrimSpace(parts[1])
			}
		case "MODE":
			if len(parts) > 1 {
				message.DeviceInfo.Mode = strings.TrimSpace(parts[1])
			}
		case "UNIT":
			if len(parts) > 1 {
				message.DeviceInfo.Unit = strings.TrimSpace(parts[1])
			}
		case "SEQ":
			if len(parts) > 1 {
				message.SampleInfo.SeqNumber = strings.TrimSpace(parts[1])
			}
		case "SID":
			if len(parts) > 1 {
				message.SampleInfo.SampleID = strings.TrimSpace(parts[1])
			}
		case "PID":
			if len(parts) > 1 {
				message.SampleInfo.PatientID = strings.TrimSpace(parts[1])
				slog.Info("Parsed PID", "value", message.SampleInfo.PatientID)
			}
		case "ID":
			if len(parts) > 1 {
				message.SampleInfo.PatientName = strings.TrimSpace(parts[1])
			}
		case "TYPE":
			if len(parts) > 1 {
				message.SampleInfo.Type = strings.TrimSpace(parts[1])
			}
		case "TEST":
			if len(parts) > 1 {
				// Test name, tidak perlu disimpan karena akan ada di hasil individual
			}
		case "OPERATOR":
			if len(parts) > 1 {
				message.DeviceInfo.Operator = strings.TrimSpace(parts[1])
			}
		default:
			// Skip CURVE, HISTOGRAM, and device header
			if strings.Contains(key, "CURVE") || strings.Contains(key, "HISTOGRAM") ||
				strings.Contains(key, "EMERALD") || key == "FSE" {
				continue
			}

			// Parse test results (e.g., WBC, RBC, HGB, etc.)
			// Format: TestCode;Value;Unit/Extra;Flag;RefMin1;RefMin2;RefMax1;RefMax2
			if len(parts) >= 8 && !strings.Contains(key, "CURVE") && !strings.Contains(key, "HISTOGRAM") {
				result := entity.AbbottTestResult{
					TestCode: key,
					Value:    strings.TrimSpace(parts[1]),
					Unit:     "",                          // Abbott doesn't send unit in this format
					Flag:     strings.TrimSpace(parts[3]), // Flag is at index 3
					RefMin:   strings.TrimSpace(parts[4]),
					RefMax:   strings.TrimSpace(parts[6]),
				}

				// Only add if we have a value
				if result.Value != "" {
					message.TestResults = append(message.TestResults, result)
				}
			}
		}
	}

	// Set date and time
	message.SampleInfo.Date = date
	message.SampleInfo.Time = timeStr

	slog.Info("Abbott parsing completed",
		"patient_id", message.SampleInfo.PatientID,
		"sample_id", message.SampleInfo.SampleID,
		"test_count", len(message.TestResults))

	// Parse timestamp
	if date != "" && timeStr != "" {
		timestamp, err := parseAbbottDateTime(date, timeStr)
		if err != nil {
			slog.Warn("failed to parse Abbott timestamp", "error", err, "date", date, "time", timeStr)
		} else {
			message.Timestamp = timestamp
		}
	}

	return message, nil
}

// parseAbbottDateTime parses Abbott date/time format (DD/MM/YYYY and HH:MM:SS)
func parseAbbottDateTime(date, timeStr string) (time.Time, error) {
	dateTimeStr := fmt.Sprintf("%s %s", date, timeStr)
	// Try format: DD/MM/YYYY HH:MM:SS
	timestamp, err := time.Parse("02/01/2006 15:04:05", dateTimeStr)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// ConvertToAbbottResults converts AbbottMessage to slice of AbbottResult
func ConvertToAbbottResults(message *entity.AbbottMessage) []entity.AbbottResult {
	results := make([]entity.AbbottResult, 0, len(message.TestResults))

	for _, testResult := range message.TestResults {
		result := entity.AbbottResult{
			PatientID: message.SampleInfo.PatientID,
			SampleID:  message.SampleInfo.SampleID,
			TestName:  testResult.TestCode,
			Value:     testResult.Value,
			Unit:      testResult.Unit,
			Flag:      testResult.Flag,
			RefMin:    testResult.RefMin,
			RefMax:    testResult.RefMax,
			Timestamp: message.Timestamp,
		}
		results = append(results, result)
	}

	return results
}
