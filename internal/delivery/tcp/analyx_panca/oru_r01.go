package analyxpanca

import (
	"context"
	"fmt"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h231"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp/common"
)

func (h *Handler) ORUR01(ctx context.Context, m h231.ORU_R01, msgByte []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID

	oruR01, err := h.decodeORUR01(msgByte)
	if err != nil {
		return "", fmt.Errorf("decode failed: %w", err)
	}

	err = h.analyzerUsecase.ProcessORUR01(ctx, oruR01)
	if err != nil {
		return "", fmt.Errorf("process failed: %w", err)
	}

	msh := h.createMSHAck(oruR01.MSH, msgControlID)
	msa := &h231.MSA{
		AcknowledgementCode: "AA",
		MessageControlID:    msgControlID,
	}

	ackMsg := h231.ACK{
		HL7: h231.HL7Name{},
		MSH: msh,
		MSA: msa,
	}

	return common.EncodeWithOptions(ackMsg, &hl7.EncodeOption{
		TrimTrailingSeparator: true,
	})
}

func (h *Handler) decodeORUR01(msgByte []byte) (entity.ORU_R01, error) {
	d := hl7.NewDecoder(h231.Registry, nil)
	msg, err := d.Decode(msgByte)
	if err != nil {
		return entity.ORU_R01{}, fmt.Errorf("decode failed: %w", err)
	}

	oul22, ok := msg.(h231.ORU_R01)
	if !ok {
		return entity.ORU_R01{}, fmt.Errorf("invalid message type, expected ORU_R01, got %T", msg)
	}

	data, err := h.MapORUR01ToEntity(&oul22)
	if err != nil {
		return entity.ORU_R01{}, fmt.Errorf("mapping failed: %w", err)
	}

	return data, nil
}

func (h *Handler) MapORUR01ToEntity(msg *h231.ORU_R01) (entity.ORU_R01, error) {
	msh := mapMSHToEntity(msg.MSH)

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

func (h *Handler) mapORUR01OrderObservationToSpecimenEntity(
	s h231.ORU_R01_OrderObservation,
	res h231.ORU_R01_PatientResult,
) entity.Specimen {
	specimen := h.mapPIDtoSpecimenEntity(res.Patient.PID)
	observationResults := []entity.ObservationResult{}

	for _, o := range s.Observation {
		observationResults = append(observationResults, h.mapOBXToObservationResultEntity(o.OBX))
	}
	specimen.ObservationResult = observationResults

	return specimen
}

func (h *Handler) mapPIDtoSpecimenEntity(pid *h231.PID) entity.Specimen {
	if pid == nil {
		return entity.Specimen{}
	}

	if len(pid.PatientIdentifierList) == 0 {
		return entity.Specimen{}
	}

	barcode := pid.PatientIdentifierList[0].ID

	return entity.Specimen{
		Barcode: barcode,
	}
}

func (h *Handler) getObservationIdentifier(field h231.CE) (string, string) {
	return field.Identifier, field.Text
}

func (h *Handler) getUnits(field *h231.CE) string {
	if field == nil {
		return ""
	}
	return field.Identifier
}

func (h *Handler) mapObservationValueToValues(values []h231.VARIES) entity.JSONStringArray {
	if values == nil {
		return entity.JSONStringArray{}
	}

	var results entity.JSONStringArray
	for i := range values {
		results = append(results, fmt.Sprintf("%v", values[i]))
	}
	return results
}

func (h *Handler) mapOBXToObservationResultEntity(obx *h231.OBX) entity.ObservationResult {
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

func (h *Handler) mapORUR01PatientToPatientEntity(p h231.ORU_R01_PatientResult) entity.Patient {
	if p.Patient == nil {
		return entity.Patient{}
	}

	return mapPIDToPatientEntity(p.Patient.PID)
}

func (h *Handler) createMSHAck(m entity.MSH, msgControlID h231.ST) *h231.MSH {
	msh := &h231.MSH{
		HL7:                  h231.HL7Name{},
		FieldSeparator:       "|",
		EncodingCharacters:   "^~\\&",
		SendingApplication:   common.SimpleHD231(constant.ThisApplication),
		SendingFacility:      common.SimpleHD231(constant.ThisFacility),
		ReceivingApplication: common.SimpleHD231(m.SendingApplication),
		ReceivingFacility:    common.SimpleHD231(m.SendingFacility),
		DateTimeOfMessage:    time.Now(),
		Security:             "",
		MessageType: h231.MSG{
			HL7:          h231.HL7Name{},
			MessageType:  "ACK",
			TriggerEvent: "R01",
		},
		MessageControlID:    "1",
		ProcessingID:        h231.PT{ProcessingID: "P"},
		VersionID:           h231.VID{VersionID: "2.3.1"},
		SequenceNumber:      "",
		ContinuationPointer: "",
		// AcceptAcknowledgmentType:            "ER",
		// ApplicationAcknowledgmentType:       "AL",
		// CountryCode:                         "ID",
		CharacterSet:                        []string{"UNICODE"},
		PrincipalLanguageOfMessage:          &h231.CE{},
		AlternateCharacterSetHandlingScheme: "",
		// MessageProfileIdentifier: []h231.EI{
		// 	{
		// 		HL7:              h231.HL7Name{},
		// 		EntityIdentifier: "LAB-28",
		// 		NamespaceID:      "IHE",
		// 		UniversalID:      "",
		// 		UniversalIDType:  "",
		// 	},
		// },
	}
	return msh
}
