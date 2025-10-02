package ncc3300

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
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
				h.handleMessage(context.Background(), rawHl7)
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
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "|OBX", "\nOBX")
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "|OBR", "\nOBR")
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "|PID", "\nPID")
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "UNICODEPID", "UNICODE\nPID")
	rawWithNewline = strings.ReplaceAll(rawWithNewline, "MSH|", "\nMSH|")

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
		} else if strings.HasPrefix(line, "MSH") || strings.HasPrefix(line, "OBR") || strings.HasPrefix(line, "PID") {
			if strings.HasPrefix(line, "PID") {
				parts := strings.Split(line, "|")
				if len(parts) > 7 && parts[7] == "1" {
					parts[7] = ""
					line = strings.Join(parts, "|")
				}
			}
			resultArray = append(resultArray, line)
		}
	}

	msgByte := []byte{}

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

	for _, result := range resultArray {
		msgByte = append(msgByte, []byte(result+"\r")...)
	}

	headerDecoder := hl7.NewDecoder(h251.Registry, &hl7.DecodeOption{HeaderOnly: true})
	header, err := headerDecoder.Decode(msgByte)
	if err != nil {
		return "", fmt.Errorf("decode header failed: %w", err)
	}

	switch m := header.(type) {
	case h251.ORU_R01:
		return h.ORUR01(ctx, m, msgByte)
	}

	return "", fmt.Errorf("unknown message type %T", header)
}
