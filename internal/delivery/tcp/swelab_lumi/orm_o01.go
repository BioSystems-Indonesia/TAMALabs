package swelablumi

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp/common"
	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
)

func (h *Handler) ORMO01(ctx context.Context, m h251.ORM_O01, message []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID

	ormO01, err := h.decodeORMO01(message)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	specimens, err := h.analyzerUsecase.ProcessORMO01(ctx, ormO01)
	if err != nil {
		return "", fmt.Errorf("process failed: %w", err)
	}
	if len(specimens) == 0 {
		return "", errors.New("error on process, no specimen found")
	}

	orrO02 := h.createORRO02(ormO01, msgControlID, specimens)

	return common.Encode(orrO02)
}

func (h *Handler) createORRO02(
	ormO01 entity.ORM_O01,
	msgControlID h251.ST,
	specimens []entity.Specimen,
) h251.ORR_O02 {
	msh := h.createMSH(ormO01.Msh, msgControlID)
	msa := h251.MSA{
		AcknowledgmentCode: "AA",
		MessageControlID:   msh.MessageControlID,
		TextMessage:        "Message accepted",
	}
	orders := h.createORRO02Order(specimens)
	response := h251.ORR_O02_Response{
		HL7: h251.HL7Name{},
		Patient: &h251.ORR_O02_Patient{
			HL7: h251.HL7Name{},
			PID: common.EncodeToPID(specimens[0].Patient),
		},
		Order: orders,
	}
	orrO02 := h251.ORR_O02{
		HL7:      h251.HL7Name{},
		MSH:      msh,
		MSA:      &msa,
		Response: &response,
		ERR:      nil,
	}
	return orrO02
}

func (h *Handler) decodeORMO01(message []byte) (entity.ORM_O01, error) {
	d := hl7.NewDecoder(h251.Registry, nil)
	msg, err := d.Decode(message)
	if err != nil {
		return entity.ORM_O01{}, fmt.Errorf("decode failed: %w", err)
	}

	m, ok := msg.(h251.ORM_O01)
	if !ok {
		return entity.ORM_O01{}, fmt.Errorf("invalid message type, expected ORM_O01, got %T", msg)
	}

	ormO01, err := h.MapORMO01ToEntity(&m)
	if err != nil {
		return entity.ORM_O01{}, fmt.Errorf("mapping failed: %w", err)
	}

	return ormO01, nil
}

func (h *Handler) createORRO02Order(specimens []entity.Specimen) []h251.ORR_O02_Order {
	orders := []h251.ORR_O02_Order{}
	for _, s := range specimens {
		for i, or := range s.ObservationRequest {
			date := time.Now()
			o := h251.ORR_O02_Order{
				HL7: h251.HL7Name{},
				ORC: &h251.ORC{
					OrderControl:      h251.ID(constant.OrderControlNodeAF),
					PlacerOrderNumber: &h251.EI{EntityIdentifier: s.Barcode},
					// FillerOrderNumber:     &h251.EI{EntityIdentifier: s.Barcode},
					DateTimeOfTransaction: date,
				},
				OrderDetailSegment: &h251.ORR_O02_OrderDetailSegment{
					HL7: h251.HL7Name{},
					OBR: &h251.OBR{
						SetID:             strconv.Itoa(i + 1),
						PlacerOrderNumber: &h251.EI{EntityIdentifier: s.Barcode},
						// FillerOrderNumber: &h251.EI{EntityIdentifier: s.Barcode},
						UniversalServiceIdentifier: h251.CE{
							Identifier: or.TestCode,
							Text:       or.TestCode,
						},
						Priority:            string(entity.PriorityR),
						ObservationDateTime: time.Now(),
						RequestedDateTime:   time.Now(),
					},
				},
			}
			orders = append(orders, o)
		}
	}
	return orders
}

func (h *Handler) createMSH(m entity.MSH, msgControlID h251.ST) *h251.MSH {
	msh := &h251.MSH{
		HL7:                  h251.HL7Name{},
		FieldSeparator:       "|",
		EncodingCharacters:   "^~\\&",
		SendingApplication:   common.SimpleHD(m.SendingApplication),
		SendingFacility:      common.SimpleHD(m.SendingFacility),
		ReceivingApplication: common.SimpleHD(constant.ThisApplication),
		ReceivingFacility:    common.SimpleHD(constant.ThisFacility),
		DateTimeOfMessage:    time.Now(),
		Security:             "",
		MessageType: h251.MSG{
			HL7:          h251.HL7Name{},
			MessageCode:  "ORR",
			TriggerEvent: "O02",
		},
		MessageControlID:                    msgControlID,
		ProcessingID:                        h251.PT{ProcessingID: "P"},
		VersionID:                           h251.VID{VersionID: "2.3.1"},
		SequenceNumber:                      "",
		ContinuationPointer:                 "",
		AcceptAcknowledgmentType:            "ER",
		ApplicationAcknowledgmentType:       "AL",
		CountryCode:                         "ID",
		CharacterSet:                        []string{"UTF-8"},
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
	return msh
}

func (h *Handler) MapORMO01ToEntity(msg *h251.ORM_O01) (entity.ORM_O01, error) {
	msh := common.MapMSHToEntity(msg.MSH)

	orders := []entity.Order_ORM_O01{}
	for _, o := range msg.Order {
		barcode := o.ORC.FillerOrderNumber
		orders = append(orders, entity.Order_ORM_O01{
			Barcode: barcode.EntityIdentifier,
		})
	}

	return entity.ORM_O01{
		Msh:    msh,
		Orders: orders,
	}, nil
}
