package tcp

import (
	"bytes"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
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
	msgByte := []byte(message)
	headerDecoder := hl7.NewDecoder(h251.Registry, &hl7.DecodeOption{HeaderOnly: true})
	header, err := headerDecoder.Decode(msgByte)
	if err != nil {
		return "", err
	}

	switch m := header.(type) {
	case h251.OUL_R22:
		d := hl7.NewDecoder(h251.Registry, nil)
		msgByte = h.deleteSegment(msgByte, "ORC")
		msg, err := d.Decode(msgByte)
		if err != nil {
			return "", err
		}
		oul22 := msg.(h251.OUL_R22)
		data, err := MapOULR22ToEntity(&oul22)
		if err != nil {
			return "", err
		}
		err = h.AnalyzerUsecase.ProcessOULR22(ctx, data)
		if err != nil {
			return "", err
		}
	case h251.OUL_R21:
		log.Println(m)
	}

	if err != nil {
		return "", err
	}

	// MSH segment
	msh := h251.MSH{
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
			TriggerEvent:     "OUL_R22",
			MessageStructure: "ACK",
		},
		MessageControlID:                    uuid.NewString(),
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

	msa := h251.MSA{
		AcknowledgmentCode: "AA",
		MessageControlID:   msh.MessageControlID,
		TextMessage:        "Message accepted",
	}

	ackMsg := h251.ACK{
		HL7: h251.HL7Name{},
		MSH: &msh,
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
	bb = bytes.ReplaceAll(bb, []byte{'\r'}, []byte{'\n'})

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
