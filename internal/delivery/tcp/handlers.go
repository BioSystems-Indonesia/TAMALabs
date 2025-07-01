// TODO: Move the handler to strategy pattern
package tcp

import (
	"context"
	"fmt"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
)

// OULR22 handles the OULR22 message.
func (h *HlSevenHandler) OULR22(ctx context.Context, m h251.OUL_R22, message []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID
	d := hl7.NewDecoder(h251.Registry, nil)
	message = h.deleteSegment(message, "ORC")
	msg, err := d.Decode(message)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	oul22 := msg.(h251.OUL_R22)
	data, err := MapOULR22ToEntity(&oul22)
	if err != nil {
		return "", fmt.Errorf("mapping failed: %w", err)
	}

	err = h.analyzerUsecase.ProcessOULR22(ctx, data)
	if err != nil {
		return "", fmt.Errorf("process failed: %w", err)
	}

	msh := h251.MSH{
		HL7:                  h251.HL7Name{},
		FieldSeparator:       "|",
		EncodingCharacters:   "^~\\&",
		SendingApplication:   m.MSH.SendingApplication,
		SendingFacility:      m.MSH.SendingFacility,
		ReceivingApplication: SimpleHD(constant.ThisApplication),
		ReceivingFacility:    SimpleHD(constant.ThisFacility), // TODO maybe need device location
		DateTimeOfMessage:    time.Now(),
		Security:             "",
		MessageType: h251.MSG{
			HL7:              h251.HL7Name{},
			MessageCode:      "ACK",
			TriggerEvent:     "R22",
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

	return createACKMessage(ackMsg)
}

func (h *HlSevenHandler) QBPQ11(ctx context.Context, m h251.QBP_Q11, message []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID
	qbp11, err := h.qbpDecoder(message)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	msg, err := MapQBPQ11ToEntity(&qbp11)
	if err != nil {
		return "", fmt.Errorf("mapping failed: %w", err)
	}

	err = h.analyzerUsecase.ProcessQBPQ11(ctx, msg)
	if err != nil {
		return "", fmt.Errorf("process failed: %w", err)
	}

	var msh *h251.MSH
	msh.MessageControlID = msgControlID
	msh.MessageControlID = msgControlID
	msh.MessageType = h251.MSG{
		HL7:              h251.HL7Name{},
		MessageCode:      "ACK",
		TriggerEvent:     "Q11",
		MessageStructure: "ACK",
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

	return createACKMessage(ackMsg)
}

func createACKMessage(msg h251.ACK) (string, error) {
	e := hl7.NewEncoder(nil)
	ack, err := e.Encode(msg)
	if err != nil {
		return "", err
	}

	//bbLog := bytes.ReplaceAll(bb, []byte{'\r'}, []byte{'\n'})
	//log.Println("Sending message: ", string(bbLog))

	return string(ack), nil
}

func SimpleHD(id string) *h251.HD {
	return &h251.HD{
		HL7:             h251.HL7Name{},
		NamespaceID:     id,
		UniversalID:     "",
		UniversalIDType: "",
	}
}
