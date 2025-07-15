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
	// raw := `MSH|^~\&|2|3|LIS|PC|20250711113400||ORU^R01|4|P|2.3.1||||||UNICODEPID|7||10||||1|OBR|5||000000000000|2^3|||20250711113300||||||||3|||OBX|1|NM|WBC||0.0|x10^9/L|3.5-10.0|L|||F||||||OBX|1|NM|LY%||**.*|%|20.0-40.0||||F||||||OBX|1|NM|MO%||**.*|%|1.0-15.0||||F||||||OBX|1|NM|GR%||**.*|%|50.0-70.0||||F||||||OBX|1|NM|LY#||**.*|x10^9/L|0.6-4.1||||F||||||OBX|1|NM|MO#||**.*|x10^9/L|0.1-1.8||||F||||||OBX|1|NM|GR#||**.*|x10^9/L|2.0-7.8||||F||||||OBX|1|NM|RBC||0.00|x10^12/L|3.50-6.00|L|||F||||||OBX|1|NM|HGB||0.0|g/dL|11.0-17.5|L|||F||||||OBX|1|NM|HCT||0.0|%|35.0-54.0|L|||F||||||OBX|1|NM|MCV||**.*|fL|80.0-100.0||||F||||||OBX|1|NM|MCH||**.*|Pg|26.0-34.0||||F||||||OBX|1|NM|MCHC||**.*|g/dL|31.5-36.0||||F||||||OBX|1|NM|RDW_CV||**.*|%|11.0-16.0||||F||||||OBX|1|NM|RDW_SD||**.*|fL|35.0-56.0||||F||||||OBX|1|NM|PLT||0|x10^9/L|100-350|L|||F||||||OBX|1|NM|MPV||**.*|fL|6.5-12.0||||F||||||OBX|1|NM|PDW||**.*|fL|9.0-17.0||||F||||||OBX|1|NM|PCT||0.00|%|0.10-0.28|L|||F||||||OBX|1|NM|P_LCR||0.0|%|11.0-45.0|L|||F||||||OBX|1|NM|P_LCC||0|x10^9/L|11-135|L|||F||||||OBX|1|NM|WBCHistogram^LeftLine||17||||||F||||||OBX|1|NM|WBCHistogram^RightLine||39||||||F||||||OBX|1|NM|WBCHistogram^MiddleLine||70||||||F||||||OBX|1|ED|WBCHistogram||3^Histogram^32Byte^HEX^00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000||||||F||||||OBX|1|NM|RBCHistogram^LeftLine||5||||||F||||||OBX|1|NM|RBCHistogram^RightLine||15||||||F||||||OBX|1|ED|RBCHistogram||3^Histogram^32Byte^HEX^00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000||||||F||||||OBX|1|NM|PLTHistogram^LeftLine||4||||||F||||||OBX|1|NM|PLTHistogram^RightLine||148||||||F||||||OBX|1|ED|PLTHistogram||3^Histogram^32Byte^HEX^0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000||||||F||||||exit status 0xc000013a`

	// h.handleMessage(context.Background(), raw)

	buf := make([]byte, 1024)
	var timer *time.Timer
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

		// Reset timer setiap ada data baru
		if timer != nil {
			timer.Stop()
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
