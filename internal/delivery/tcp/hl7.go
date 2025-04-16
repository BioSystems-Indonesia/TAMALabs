package tcp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
	"strings"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp"
)

// HlSevenHandler is a struct that contains the handler of the REST server.
type HlSevenHandler struct {
	AnalyzerUsecase usecase.Analyzer
}

// NewHlSevenHandler creates a new instance of HlSevenHandler.
func NewHlSevenHandler(analyzerUsecase usecase.Analyzer) *HlSevenHandler {
	return &HlSevenHandler{
		AnalyzerUsecase: analyzerUsecase,
	}
}

func (h *HlSevenHandler) Handle(conn *net.TCPConn) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			log.Println("panic: recovered in processing connection", r)
		}
	}()

	mc := mllp.NewClient(conn)
	b, err := mc.ReadAll()
	if err != nil {
		if err != io.EOF {
			log.Println(err)
		}
	}

	res, err := h.HL7Handler(context.Background(), string(b))
	if err != nil {
		log.Println(err)
	}

	mc.Write([]byte(res))
}

// HL7Handler handles the HL7 message.
func (h *HlSevenHandler) HL7Handler(ctx context.Context, message string) (string, error) {
	if message != "" {
		logMsg := strings.ReplaceAll(message, "\r", "\n")
		log.Println("received message: ", logMsg)
	}

	// don't do anything if the message is empty
	if message == "" {
		return "", nil
	}

	msgByte := []byte(message)
	headerDecoder := hl7.NewDecoder(h251.Registry, &hl7.DecodeOption{HeaderOnly: true})
	header, err := headerDecoder.Decode(msgByte)
	if err != nil {
		return "", fmt.Errorf("decode header failed: %w", err)
	}

	switch m := header.(type) {
	case h251.OUL_R22:
		return h.OULR22(ctx, m, msgByte)
	case h251.QBP_Q11:
		return h.QBPQ11(ctx, m, msgByte)
	}

	return "", fmt.Errorf("unknown message type")
}

func (h *HlSevenHandler) qbpDecoder(message []byte) (h251.QBP_Q11, error) {
	d := hl7.NewDecoder(h251.Registry, nil)

	// get QPD segment
	qpd := h.getSegment(message, "QPD")

	// manually decode QPD segment
	qpds := strings.Split(string(qpd), "|")
	if len(qpds) < 3 {
		return h251.QBP_Q11{}, fmt.Errorf("QPD segment is not complete")
	}
	if len(qpds) == 3 {
		qpds = append(qpds, "")
	}
	UserParametersInSuccessiveFields := h251.VARIES(qpds[3])
	manualQPD := h251.QPD{
		HL7: h251.HL7Name{},
		MessageQueryName: h251.CE{
			HL7: h251.HL7Name{},
			// EntityIdentifier: qpds[0],
			// NamespaceID:      qpds[1],
			// UniversalID:      qpds[2],
			// UniversalIDType:  qpds[3],
		},
		QueryTag:                         qpds[2],
		UserParametersInSuccessiveFields: &UserParametersInSuccessiveFields,
	}

	// delete QPD segment
	message = h.deleteSegment(message, "QPD")

	msg, err := d.Decode(message)
	if err != nil {
		return h251.QBP_Q11{}, fmt.Errorf("decode failed: %w", err)
	}

	qbp11 := msg.(h251.QBP_Q11)
	qbp11.QPD = &manualQPD
	return qbp11, nil

}

func (h *HlSevenHandler) deleteSegment(message []byte, seg string) []byte {
	// Normalize line endings to \n
	message = bytes.ReplaceAll(message, []byte("\r"), []byte("\n"))
	lines := bytes.Split(message, []byte("\n"))

	var filteredLines [][]byte
	for _, line := range lines {
		if !bytes.HasPrefix(line, []byte(seg)) {
			filteredLines = append(filteredLines, line)
		}
	}

	return bytes.Join(filteredLines, []byte("\n"))
}

func (h *HlSevenHandler) getSegment(message []byte, seg string) []byte {
	// Normalize line endings to \n
	message = bytes.ReplaceAll(message, []byte("\r"), []byte("\n"))
	lines := bytes.Split(message, []byte("\n"))

	var filteredLines [][]byte
	for _, line := range lines {
		if bytes.HasPrefix(line, []byte(seg)) {
			filteredLines = append(filteredLines, line)
		}
	}

	return bytes.Join(filteredLines, []byte("\n"))
}

func simpleHD(id string) *h251.HD {
	return &h251.HD{
		HL7:             h251.HL7Name{},
		NamespaceID:     id,
		UniversalID:     "",
		UniversalIDType: "",
	}
}
