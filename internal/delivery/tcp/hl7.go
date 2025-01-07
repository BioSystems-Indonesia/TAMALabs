package tcp

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/usecase"
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

// HL7Handler handles the HL7 message.
func (h *HlSevenHandler) HL7Handler(ctx context.Context, message string) (string, error) {
	if message != "" {
		logMsg := strings.ReplaceAll(message, "\r", "\n")
		log.Println("Received message: ", logMsg)
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

	msgControlID := ""

	switch m := header.(type) {
	case h251.OUL_R22:
		msgControlID = m.MSH.MessageControlID
		d := hl7.NewDecoder(h251.Registry, nil)
		msgByte = h.deleteSegment(msgByte, "ORC")
		msg, err := d.Decode(msgByte)
		if err != nil {
			return "", fmt.Errorf("decode failed: %w", err)
		}
		oul22 := msg.(h251.OUL_R22)
		data, err := MapOULR22ToEntity(&oul22)
		if err != nil {
			return "", fmt.Errorf("mapping failed: %w", err)
		}
		log.Printf("%#v", data)
		err = h.AnalyzerUsecase.ProcessOULR22(ctx, data)
		if err != nil {
			return "", fmt.Errorf("process failed: %w", err)
		}
	case h251.OUL_R21:
		log.Println(m)
	}

	if err != nil {
		return "", err
	}

	// MSH segment
	// TODO: FIXME
	var msh *h251.MSH
	switch headerNew := header.(type) {
	case h251.OUL_R22:
		msh = headerNew.MSH
		msh.MessageType = h251.MSG{
			HL7:              h251.HL7Name{},
			MessageCode:      "ACK",
			TriggerEvent:     "OUL",
			MessageStructure: "ACK",
		}
	default:
		msh = &h251.MSH{
			HL7:                  h251.HL7Name{},
			FieldSeparator:       "|",
			EncodingCharacters:   "^~\\&",
			SendingApplication:   simpleHD("BioLIS"),
			SendingFacility:      simpleHD("Lab1"),
			ReceivingApplication: simpleHD("BA200"),
			ReceivingFacility:    simpleHD("Lab1"),
			DateTimeOfMessage:    time.Now(),
			Security:             "",
			MessageType: h251.MSG{
				HL7:              h251.HL7Name{},
				MessageCode:      "ACK",
				TriggerEvent:     "OUL",
				MessageStructure: "ACK",
			},
			MessageControlID:                    msgControlID,
			ProcessingID:                        h251.PT{ProcessingID: "P"},
			VersionID:                           h251.VID{VersionID: "2.5.1"},
			SequenceNumber:                      "",
			ContinuationPointer:                 "",
			AcceptAcknowledgmentType:            "ER",
			ApplicationAcknowledgmentType:       "AL",
			CountryCode:                         "ID",
			CharacterSet:                        []string{"UNICODE UTF-8"},
			PrincipalLanguageOfMessage:          &h251.CE{},
			AlternateCharacterSetHandlingScheme: "",
			MessageProfileIdentifier: []h251.EI{
				{
					HL7:              h251.HL7Name{},
					EntityIdentifier: "LAB-28",
					NamespaceID:      "IHE",
					UniversalID:      "",
					UniversalIDType:  "",
				},
			},
		}
	}

	msa := h251.MSA{
		AcknowledgmentCode: "AA",
		MessageControlID:   msh.MessageControlID,
		TextMessage:        "Message accepted",
	}

	ackMsg := h251.ACK{
		HL7: h251.HL7Name{},
		MSH: msh,
		SFT: nil,
		MSA: &msa,
		ERR: nil,
	}

	// Create Encoder with options
	e := hl7.NewEncoder(nil)
	bb, err := e.Encode(ackMsg)
	if err != nil {
		return "", err
	}

	bbLog := bytes.ReplaceAll(bb, []byte{'\r'}, []byte{'\n'})
	log.Println("Sending message: ", string(bbLog))
	//bb = bytes.ReplaceAll(bb, []byte{'\r'}, []byte{'\n'})

	// Encode the message

	return string(bb), nil
}

func (h *HlSevenHandler) deleteSegment(message []byte, seg string) []byte {
	lines := bytes.Split(message, []byte("\n"))

	var filteredLines [][]byte
	for _, line := range lines {
		if !bytes.HasPrefix(line, []byte(seg)) {
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
