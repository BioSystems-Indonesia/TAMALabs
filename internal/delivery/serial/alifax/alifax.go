package alifax

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"go.bug.st/serial"
)

const (
	STX = 0x02 // Start of Text
	ETX = 0x03 // End of Text
)

// Data represents the parsed serial data structure
type Data struct {
	Command           string // e.g., "R"
	WorkstationNumber string // e.g., "01"
	PatientID         string // 15 characters
	RackNo            string // 2 characters
	Position          string // 2 characters
	Cycle             string // 2 characters
	Result            string // 4 characters
	Checksum          string // 1 character
}

type Handler struct {
	analyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		analyzerUsecase: analyzerUsecase,
	}
}

func (h *Handler) Handle(port serial.Port) {
	var buffer []byte
	for {
		readBuffer := make([]byte, 128)
		n, err := port.Read(readBuffer)
		if err != nil {
			log.Printf("Error reading from serial port: %v", err)
			continue
		}

		if n > 0 {
			// Append received data to buffer
			buffer = append(buffer, readBuffer[:n]...)

			// Process complete messages in buffer
			buffer = h.processBuffer(buffer)
		}

		// Small delay to prevent excessive CPU usage
		time.Sleep(10 * time.Millisecond)
	}
}

func (h *Handler) processBuffer(buffer []byte) []byte {
	ctx := context.Background()
	for {
		// Find STX
		stxIndex := findSTX(buffer)
		if stxIndex == -1 {
			// No STX found, clear buffer
			return buffer[:0]
		}

		// Remove data before STX
		if stxIndex > 0 {
			buffer = buffer[stxIndex:]
		}

		// Find ETX after STX
		etxIndex := findETX(buffer[1:]) // Start searching after STX
		if etxIndex == -1 {
			// No ETX found, keep buffer as is and wait for more data
			return buffer
		}

		// Adjust ETX index (add 1 because we started searching from index 1)
		etxIndex += 1

		// Extract message (excluding STX and ETX)
		message := buffer[1:etxIndex]

		// Parse the message
		parsed, err := parseMessage(message)
		if err != nil {
			slog.Error("error parsing", "err", err, "raw", string(message))
		} else {
			slog.Info("parsed data",
				"command", parsed.Command,
				"workstationNumber", parsed.WorkstationNumber,
				"patientID", parsed.PatientID,
				"rackNo", parsed.RackNo,
				"position", parsed.Position,
				"cycle", parsed.Cycle,
				"result", parsed.Result,
				"checksum", parsed.Checksum,
			)
			err := h.analyzerUsecase.ProcessORUR01(ctx, toORU(parsed))
			if err != nil {
				slog.Error("Error processing ORUR01", "err", err)
			}
		}

		// Remove processed message from buffer
		buffer = buffer[etxIndex+1:]

		// If buffer is empty, break
		if len(buffer) == 0 {
			break
		}
	}

	return buffer
}

func findSTX(data []byte) int {
	for i, b := range data {
		if b == STX {
			return i
		}
	}
	return -1
}

func findETX(data []byte) int {
	for i, b := range data {
		if b == ETX {
			return i
		}
	}
	return -1
}

func parseMessage(message []byte) (Data, error) {
	messageStr := string(message)

	// Expected message format:
	// R (1) 01 (2) + PatientID (15) + RackNo (2) + Position (2) + Cycle/Bayer (2) + Result (4) + Checksum (1)
	// Total expected length: 29 characters

	if len(messageStr) < 29 {
		return Data{}, fmt.Errorf("message too short: expected at least 31 characters, got %d", len(messageStr))
	}

	parsed := Data{
		Command:           messageStr[0:1],                     // R01
		WorkstationNumber: messageStr[1:3],                     // R01
		PatientID:         strings.TrimSpace(messageStr[3:18]), // 15 chars (trim spaces)
		RackNo:            messageStr[18:20],                   // 2 chars
		Position:          messageStr[20:22],                   // 2 chars
		Cycle:             messageStr[22:24],                   // 2 chars
		Result:            messageStr[24:28],                   // 4 chars
		Checksum:          messageStr[28:29],                   // 1 char
	}

	// Validate command
	if parsed.Command != "R" {
		return parsed, fmt.Errorf("invalid command: expected 'R01', got '%s'", parsed.Command)
	}

	// Verify checksum if needed
	if !verifyChecksum(messageStr[:29], parsed.Checksum) {
		log.Printf("Warning: Checksum verification failed")
	}

	return parsed, nil
}

// TODO verify to read Checksum
func verifyChecksum(data string, expectedChecksum string) bool {
	// Simple XOR checksum calculation
	var checksum byte
	for _, b := range []byte(data) {
		checksum ^= b
	}

	// Convert checksum to ASCII character
	calculatedChecksum := string(rune(checksum))

	return calculatedChecksum == expectedChecksum
}

func displayParsedData(data Data) {
	fmt.Println("\n=== Received Data ===")
	fmt.Printf("Command:     %s\n", data.Command)
	fmt.Printf("Workstation: %s\n", data.WorkstationNumber)
	fmt.Printf("Patient ID:  %s\n", data.PatientID)
	fmt.Printf("Rack No:     %s\n", data.RackNo)
	fmt.Printf("Position:    %s\n", data.Position)
	fmt.Printf("Cycle:       %s\n", data.Cycle)
	fmt.Printf("Result:      %s\n", data.Result)
	fmt.Printf("Checksum:    %s\n", data.Checksum)
	fmt.Printf("Timestamp:   %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println("=====================")
}

func toORU(d Data) entity.ORU_R01 {
	var number string = "0"
	var description string = ""
	if positiveNumber(d.Result) {
		number = toNumber(d.Result)
	} else {
		description = setDescription(d.Result)
	}
	return entity.ORU_R01{
		MSH: entity.MSH{},
		Patient: []entity.Patient{
			{
				Specimen: []entity.Specimen{
					{
						ReceivedDate: time.Now(),
						Barcode:      "WBL" + d.PatientID,
						ObservationResult: []entity.ObservationResult{
							{
								TestCode:       "ELD",
								Description:    description,
								Values:         entity.JSONStringArray{number},
								Type:           "ELD",
								Unit:           "mm/h",
								ReferenceRange: "",
								Date:           time.Now(),
								Comments:       description,
								Picked:         false,
							},
						},
					},
				},
			},
		},
	}
}

func setDescription(s string) string {
	switch s {
	case "-001":
		return "NF"
	case "-002":
		return "NP"
	case "-004":
		return "NP"

	default:
		return ""
	}
}

func toNumber(s string) string {
	return strings.TrimLeft(s, "0")
}

func positiveNumber(s string) bool {
	if len(s) > 0 {
		return s[0] == '0'
	}
	return false
}
