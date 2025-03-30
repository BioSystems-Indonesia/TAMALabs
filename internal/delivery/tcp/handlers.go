package tcp

import (
	"context"
	"fmt"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
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

	err = h.AnalyzerUsecase.ProcessOULR22(ctx, data)
	if err != nil {
		return "", fmt.Errorf("process failed: %w", err)
	}

	var msh *h251.MSH
	msh.MessageControlID = msgControlID
	msh.MessageType = h251.MSG{
		HL7:              h251.HL7Name{},
		MessageCode:      "ACK",
		TriggerEvent:     "R22",
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

	if msg.Barcode == "" {
		return "", fmt.Errorf("barcode is empty")
	} else {

	}

	var msh *h251.MSH
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
