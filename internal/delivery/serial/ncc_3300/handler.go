package ncc3300

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"go.bug.st/serial"
)

type Handler struct {
	AnalyzerUsecase usecase.Analyzer
}

func NewHandler(analyzerUsecase usecase.Analyzer) *Handler {
	return &Handler{
		AnalyzerUsecase: analyzerUsecase,
	}
}

var rawHl7 string

func (h *Handler) Handle(port serial.Port) {
	buf := make([]byte, 1024)
	var timer *time.Timer

	defer func() {
		if timer != nil {
			timer.Stop()
		}
	}()

	for {
		n, err := port.Read(buf)
		if err != nil || n == 0 {
			slog.Error("Error reading serial data", "error", err)
			continue
		}

		buffer := buf[:n]
		for _, b := range buffer {
			if b >= 32 && b <= 126 {
				rawHl7 += string(b)
			}
		}

		// Reset timer setiap ada data baru dengan proper cleanup
		if timer != nil {
			if !timer.Stop() {
				// Drain the channel if timer already fired
				select {
				case <-timer.C:
				default:
				}
			}
		}

		timer = time.AfterFunc(300*time.Millisecond, func() {
			if rawHl7 != "" {
				slog.Info("Processing HL7 message", "length", len(rawHl7))
				ack, err := h.handleMessage(context.Background(), rawHl7)
				if err != nil {
					slog.Error("Failed to handle message", "error", err)
				} else {
					slog.Info("Message processed successfully", "ack", ack)
				}
				rawHl7 = ""
			}
		})
	}
}

func (h *Handler) handleMessage(ctx context.Context, message string) (string, error) {
	// don't do anything if the message is empty
	if message == "" {
		return "", nil
	}

	rawWithNewline := message

	// Remove all \r first to normalize
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "\r", "\n")

	// If already has newlines, we're good
	if !strings.Contains(rawWithNewline, "\n") {
		// No newlines found, need to split by segment names
		// Use negative lookahead pattern by inserting newline before segment names
		rawWithNewline = strings.ReplaceAll(rawWithNewline, "MSH|", "\nMSH|")
		rawWithNewline = strings.ReplaceAll(rawWithNewline, "PID|", "\nPID|")
		rawWithNewline = strings.ReplaceAll(rawWithNewline, "PV1|", "\nPV1|")
		rawWithNewline = strings.ReplaceAll(rawWithNewline, "OBR|", "\nOBR|")
		rawWithNewline = strings.ReplaceAll(rawWithNewline, "OBX|", "\nOBX|")
		rawWithNewline = strings.ReplaceAll(rawWithNewline, "UNICODEPID|", "UNICODE\nPID|")
	}

	// Clean up double newlines and trim
	rawWithNewline = strings.TrimLeft(rawWithNewline, "\n")

	resultArray := []string{}
	lines := strings.Split(rawWithNewline, "\n")
	obxCounter := 1
	for _, line := range lines {
		if strings.HasPrefix(line, "OBX") && !strings.Contains(line, "Histogram") {
			parts := strings.Split(line, "|")
			if len(parts) > 1 {
				parts[1] = fmt.Sprintf("%d", obxCounter)
				line = strings.Join(parts, "|")
				obxCounter++
			}
			resultArray = append(resultArray, line)
		} else if strings.HasPrefix(line, "MSH") || strings.HasPrefix(line, "OBR") || strings.HasPrefix(line, "PID") || strings.HasPrefix(line, "PV1") {
			if strings.HasPrefix(line, "PID") {
				parts := strings.Split(line, "|")
				// Fix invalid birthdate - if it's just "1", clear it
				if len(parts) > 7 && parts[7] == "1" {
					parts[7] = ""
				}
				// Fix missing or invalid gender - if field 8 is empty, "1", or contains invalid chars
				if len(parts) > 8 {
					gender := strings.TrimSpace(parts[8])
					// Valid gender: M, F, O, U, A, N
					if gender != "M" && gender != "F" && gender != "O" && gender != "U" && gender != "A" && gender != "N" {
						parts[8] = ""
					}
				}
				line = strings.Join(parts, "|")
			}
			resultArray = append(resultArray, line)
		}
	}

	msgByte := []byte{}

	// Check if OBR exists
	hasOBR := false
	for _, line := range resultArray {
		if strings.HasPrefix(line, "OBR") {
			hasOBR = true
			break
		}
	}

	for i, line := range resultArray {
		if strings.HasPrefix(line, "PID") {
			// Insert PV1 after PID
			pv1 := "PV1|1||||||||||||||||||||"
			newArray := make([]string, len(resultArray)+1)
			copy(newArray, resultArray[:i+1])
			newArray[i+1] = pv1
			copy(newArray[i+2:], resultArray[i+1:])
			resultArray = newArray
			break
		}
	}

	// If no OBR exists, add a dummy OBR before first OBX
	if !hasOBR {
		// Try to get barcode from PID segment
		barcode := ""
		for _, line := range resultArray {
			if strings.HasPrefix(line, "PID") {
				parts := strings.Split(line, "|")
				// PID field 2 is Patient ID (barcode)
				if len(parts) > 2 && parts[2] != "" {
					barcode = parts[2]
					break
				}
			}
		}

		for i, line := range resultArray {
			if strings.HasPrefix(line, "OBX") {
				// Insert dummy OBR before first OBX with barcode
				obr := fmt.Sprintf("OBR|1||%s|CBC^Complete Blood Count|||%s", barcode, time.Now().Format("20060102150405"))
				newArray := make([]string, len(resultArray)+1)
				copy(newArray, resultArray[:i])
				newArray[i] = obr
				copy(newArray[i+1:], resultArray[i:])
				resultArray = newArray
				break
			}
		}
	}

	for _, result := range resultArray {
		msgByte = append(msgByte, []byte(result+"\r")...)
	}

	headerDecoder := hl7.NewDecoder(h251.Registry, &hl7.DecodeOption{HeaderOnly: true})
	header, err := headerDecoder.Decode(msgByte)
	if err != nil {
		slog.Error("decode header failed", "error", err, "msgByte", string(msgByte))
		return "", fmt.Errorf("decode header failed: %w", err)
	}

	slog.Info(string(msgByte))

	switch m := header.(type) {
	case h251.ORU_R01:
		ack, err := h.ORUR01(ctx, m, msgByte)
		if err != nil {
			slog.Error("ORUR01 processing failed", "error", err)
			return "", err
		}
		return ack, nil
	}

	slog.Error("unknown message type", "type", fmt.Sprintf("%T", header))
	return "", fmt.Errorf("unknown message type %T", header)
}
