package analyxtrias

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h231"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp/common"
)

func (h *Handler) ORMO01(ctx context.Context, m h231.ORM_O01, message []byte) (string, error) {
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
	msgControlID h231.ST,
	specimens []entity.Specimen,
) h231.ORR_O02 {
	msh := h.createMSH(ormO01.Msh, msgControlID)
	msa := h231.MSA{
		AcknowledgementCode: "AA",
		MessageControlID:    msh.MessageControlID,
		TextMessage:         "Message accepted",
	}
	orders := h.createORRO02Order(specimens)
	response := h231.ORR_O02_Response{
		HL7: h231.HL7Name{},
		Patient: &h231.ORR_O02_Patient{
			HL7: h231.HL7Name{},
			PID: EncodeToPID(specimens[0].Patient),
		},
		Order: orders,
	}
	orrO02 := h231.ORR_O02{
		HL7:      h231.HL7Name{},
		MSH:      msh,
		MSA:      &msa,
		Response: &response,
		ERR:      nil,
	}
	return orrO02
}

func (h *Handler) decodeORMO01(message []byte) (entity.ORM_O01, error) {
	d := hl7.NewDecoder(h231.Registry, nil)
	msg, err := d.Decode(message)
	if err != nil {
		return entity.ORM_O01{}, fmt.Errorf("decode failed: %w", err)
	}

	m, ok := msg.(h231.ORM_O01)
	if !ok {
		return entity.ORM_O01{}, fmt.Errorf("invalid message type, expected ORM_O01, got %T", msg)
	}

	ormO01, err := h.MapORMO01ToEntity(&m)
	if err != nil {
		return entity.ORM_O01{}, fmt.Errorf("mapping failed: %w", err)
	}

	return ormO01, nil
}

func (h *Handler) createORRO02Order(specimens []entity.Specimen) []h231.ORR_O02_Order {
	orders := []h231.ORR_O02_Order{}
	for _, s := range specimens {
		for i, or := range s.ObservationRequest {
			o := h231.ORR_O02_Order{
				HL7: h231.HL7Name{},
				ORC: &h231.ORC{
					OrderControl:          h231.ID(constant.OrderControlNodeAF),
					PlacerOrderNumber:     &h231.EI{EntityIdentifier: strconv.Itoa(i + 1)},
					FillerOrderNumber:     &h231.EI{EntityIdentifier: s.Barcode},
					DateTimeOfTransaction: time.Now(),
				},
				OrderDetailSegment: &h231.ORR_O02_OrderDetailSegment{
					HL7: h231.HL7Name{},
					OBR: &h231.OBR{
						SetID: strconv.Itoa(i + 1),
						// PlacerOrderNumber: &h231.EI{EntityIdentifier: strconv.Itoa(i + 1)},
						FillerOrderNumber: &h231.EI{EntityIdentifier: s.Barcode},
						UniversalServiceID: h231.CE{
							Identifier: or.TestCode,
							// Text:       or.TestCode,
						},
						// Priority:            string(entity.PriorityR),
						RequestedDateTime:        time.Now(),
						ObservationDateTime:      time.Now(),
						SpecimenReceivedDateTime: time.Now(),
						DiagnosticServSectID:     "HM",
						PrincipalResultInterpreter: &h231.NDL{
							OPName: &h231.CN{
								GivenName:  s.WorkOrder.GetFirstDoctor().GetFirstName(),
								FamilyName: s.WorkOrder.GetFirstDoctor().GetLastName(),
							},
						},
					},
				},
			}
			orders = append(orders, o)
		}
	}
	return orders
}

func (h *Handler) MapORMO01ToEntity(msg *h231.ORM_O01) (entity.ORM_O01, error) {
	msh := MapMSHToEntity(msg.MSH)

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

func EncodeToPID(in entity.Patient) *h231.PID {
	id := strconv.FormatInt(in.ID, 10)
	return &h231.PID{
		HL7:   h231.HL7Name{},
		SetID: "1",
		// PatientID:             h231.CX{ID: id},
		PatientIdentifierList: []h231.CX{{ID: id}},
		// AlternatePatientID:    []h231.CX{},
		PatientName: []h231.XPN{{
			FamilyNameLastNamePrefix: in.LastName,
			GivenName:                in.FirstName,
		}},
		// MotherSMaidenName: []h231.XPN{},
		DateTimeOfBirth: in.Birthdate,
		Sex:             in.Sex.String(),
	}
}
