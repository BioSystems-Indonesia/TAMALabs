package ncc3300

import (
	"context"
	"fmt"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp/common"
	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
)

func (h *Handler) ORUR01(ctx context.Context, m h251.ORU_R01, msgByte []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID

	oruR01, err := h.decodeORUR01(msgByte)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	err = h.AnalyzerUsecase.ProcessORUR01(ctx, oruR01)
	if err != nil {
		return "", fmt.Errorf("process failed: %w", err)
	}

	msh := h.createMSHAck(oruR01.MSH, msgControlID)
	msa := &h251.MSA{
		AcknowledgmentCode: "AA",
		MessageControlID:   msh.MessageControlID,
		TextMessage:        "Message accepted",
	}

	ackMsg := h251.ACK{
		HL7: h251.HL7Name{},
		MSH: msh,
		MSA: msa,
	}

	return common.Encode(ackMsg)
}

func (h *Handler) decodeORUR01(msgByte []byte) (entity.ORU_R01, error) {
	d := hl7.NewDecoder(h251.Registry, nil)
	msg, err := d.Decode(msgByte)
	if err != nil {
		return entity.ORU_R01{}, fmt.Errorf("decode failed: %w", err)
	}

	oul22, ok := msg.(h251.ORU_R01)
	if !ok {
		return entity.ORU_R01{}, fmt.Errorf("invalid message type, expected ORU_R01, got %T", msg)
	}

	data, err := h.MapORUR01ToEntity(&oul22)
	if err != nil {
		return entity.ORU_R01{}, fmt.Errorf("mapping failed: %w", err)
	}

	return data, nil
}

func (h *Handler) MapORUR01ToEntity(msg *h251.ORU_R01) (entity.ORU_R01, error) {
	msh := common.MapMSHToEntity(msg.MSH)

	// Patient Mapping (PID Segment)
	var patient []entity.Patient
	for _, res := range msg.PatientResult {
		p := h.mapORUR01PatientToPatientEntity(res)
		var specimen []entity.Specimen
		for _, o := range res.OrderObservation {
			specimen = append(specimen, h.mapORUR01OrderObservationToSpecimenEntity(o, res))
		}
		p.Specimen = specimen

		patient = append(patient, p)
	}

	return entity.ORU_R01{
		MSH:     msh,
		Patient: patient,
	}, nil
}

func (h *Handler) mapORUR01OrderObservationToSpecimenEntity(s h251.ORU_R01_OrderObservation, res h251.ORU_R01_PatientResult) entity.Specimen {
	specimen := h.mapOBRToSpecimenEntity(s.OBR)
	observationResults := []entity.ObservationResult{}

	for _, o := range s.Observation {
		observationResults = append(observationResults, h.mapOBXToObservationResultEntity(o.OBX))
	}
	specimen.ObservationResult = observationResults

	return specimen
}

func (h *Handler) mapOBRToSpecimenEntity(obr *h251.OBR) entity.Specimen {
	if obr == nil {
		return entity.Specimen{}
	}

	return entity.Specimen{
		Barcode: obr.UniversalServiceIdentifier.Identifier,
	}
}

func (h *Handler) getObservationIdentifier(field h251.CE) (string, string) {
	return field.Identifier, field.Text
}

func (h *Handler) getUnits(field *h251.CE) string {
	if field == nil {
		return ""
	}
	return field.Identifier
}

func (h *Handler) mapObservationValueToValues(values []h251.VARIES) entity.JSONStringArray {
	if values == nil {
		return entity.JSONStringArray{}
	}

	var results entity.JSONStringArray
	for i := range values {
		results = append(results, fmt.Sprintf("%v", values[i]))
	}
	return results
}

func (h *Handler) mapOBXToObservationResultEntity(obx *h251.OBX) entity.ObservationResult {
	if obx == nil {
		return entity.ObservationResult{}
	}

	testCode, description := h.getObservationIdentifier(obx.ObservationIdentifier)

	return entity.ObservationResult{
		TestCode:       testCode,
		Description:    description,
		Values:         h.mapObservationValueToValues(obx.ObservationValue),
		Type:           obx.ValueType,
		Unit:           h.getUnits(obx.Units),
		ReferenceRange: obx.ReferencesRange,
		Date:           obx.DateTimeOfTheObservation,
		AbnormalFlag:   obx.AbnormalFlags,
		Comments:       obx.ObservationResultStatus,
	}
}

func (h *Handler) mapORUR01PatientToPatientEntity(p h251.ORU_R01_PatientResult) entity.Patient {
	if p.Patient == nil {
		return entity.Patient{}
	}

	return common.MapPIDToPatientEntity(p.Patient.PID)
}

func (h *Handler) createMSHAck(m entity.MSH, msgControlID h251.ST) *h251.MSH {
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
			MessageCode:  "ACK",
			TriggerEvent: "R01",
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
