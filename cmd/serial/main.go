package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

type TestResult struct {
	RecordType string   `json:"record_type"`
	DeviceID   string   `json:"device_id"`
	Status     string   `json:"status"`
	Date       string   `json:"date"`
	Time       string   `json:"time"`
	TestType   string   `json:"test_type"`
	TestName   string   `json:"test_name"`
	Value      string   `json:"value"`
	Unit       string   `json:"unit"`
	Reference  string   `json:"reference"`
	Flags      string   `json:"flags"`
	Extra      []string `json:"extra"` // To handle additional fields
}

func main() {
	// Serial port configuration - UPDATED TO 115200 BAUD
	// portName := "/dev/ttyUSB0" // Linux/macOS
	portName := "COM6" // Windows
	baudRate := 115200 // Updated baud rate for Coax Biosystem

	mode := &serial.Mode{
		BaudRate: baudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {
		log.Fatalf("Failed to open serial port: %v", err)
	}
	defer port.Close()

	log.Printf("Listening on %s at %d baud...\n", portName, baudRate)

	// Buffer for incoming data
	buf := make([]byte, 1024)
	var buffer string // Accumulate partial lines

	for {
		n, err := port.Read(buf)
		if err != nil {
			log.Printf("Error reading serial data: %v", err)
			continue
		}

		if n > 0 {
			data := string(buf[:n])
			buffer += data

			fmt.Println(data)
			// Process complete lines
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
					log.Printf("Parse error: %v | Raw: %s", err, line)
					continue
				}
				fmt.Println(buffer)

				jsonData, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					log.Printf("JSON error: %v", err)
					continue
				}

				log.Printf("Parsed result:\n%s\n", jsonData)
			}
		}

		time.Sleep(10 * time.Millisecond) // Reduced delay for high baud rate
	}
}

func parseLine(line string) (TestResult, error) {
	// Remove control characters and trim spaces
	cleaned := strings.Trim(line, "\x02\x03")
	fields := strings.Fields(cleaned) // Split by whitespace to handle irregular spacing

	// Debug: Print cleaned fields
	fmt.Printf("Cleaned fields: %v\n", fields)

	// Safely handle fields (some lines might have fewer columns)
	result := TestResult{}
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
