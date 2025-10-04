package wondfo

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/oibacidem/lims-hl-seven/internal/constant"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/pkg/mllp/common"
)

// isValidBarcodePrefix checks if the string starts with any valid barcode prefix
func (h *Handler) isValidBarcodePrefix(s string) bool {
	prefixes := []string{"SER", "WBL", "URI", "SEM", "PLM", "LIQ", "CSF"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func (h *Handler) ORUR01(ctx context.Context, m h251.ORU_R01, msgByte []byte) (string, error) {
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

	oruR01, ok := msg.(h251.ORU_R01)
	if !ok {
		return entity.ORU_R01{}, fmt.Errorf("invalid message type, expected ORU_R01, got %T", msg)
	}

	data, err := h.MapORUR01ToEntity(&oruR01)
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

	// For Wondfo, extract the patient ID from PID segment and use it as barcode
	// Based on the HL7: PID|202510020009||PatienNumber^AdmNumber^BedNumber||SER251003001||20251002|M
	// The patient ID "SER251003001" is in Patient Name field (field 5), not Patient Identifier List
	// Barcode can have different prefixes: SER, WBL, URI, SEM, PLM, LIQ, CSF
	if res.Patient != nil && res.Patient.PID != nil {
		slog.Info("Debugging PID structure for Wondfo")

		// First, try Patient Name field (field 5) - this is where SER251003001 should be
		if len(res.Patient.PID.PatientName) > 0 {
			slog.Info("Checking Patient Name fields", "count", len(res.Patient.PID.PatientName))
			for i, name := range res.Patient.PID.PatientName {
				slog.Info("Patient Name field", "index", i, "family_name", name.FamilyName, "given_name", name.GivenName)

				// Check if any name field contains barcode prefixes
				if h.isValidBarcodePrefix(name.FamilyName) {
					specimen.Barcode = name.FamilyName
					slog.Info("Found barcode in FamilyName", "barcode", name.FamilyName)
					break
				}
				if h.isValidBarcodePrefix(name.GivenName) {
					specimen.Barcode = name.GivenName
					slog.Info("Found barcode in GivenName", "barcode", name.GivenName)
					break
				}
			}
		}

		// If not found in Patient Name, check Patient Identifier List (field 3)
		if specimen.Barcode == "" {
			slog.Info("Checking Patient Identifier List", "count", len(res.Patient.PID.PatientIdentifierList))
			for i, id := range res.Patient.PID.PatientIdentifierList {
				slog.Info("PID identifier", "index", i, "id_number", id.IDNumber, "assigning_authority", id.AssigningAuthority)

				if h.isValidBarcodePrefix(id.IDNumber) {
					specimen.Barcode = id.IDNumber
					slog.Info("Found barcode in identifier list", "barcode", id.IDNumber)
					break
				}

				// Check if the ID contains ^ separated values
				if strings.Contains(id.IDNumber, "^") {
					parts := strings.Split(id.IDNumber, "^")
					for j, part := range parts {
						slog.Info("Checking compound part", "part_index", j, "part_value", part)
						if h.isValidBarcodePrefix(part) {
							specimen.Barcode = part
							slog.Info("Found barcode in compound field", "barcode", part)
							break
						}
					}
					if specimen.Barcode != "" {
						break
					}
				}
			}
		}

		// Fallback - use first non-compound identifier
		if specimen.Barcode == "" && len(res.Patient.PID.PatientIdentifierList) > 0 {
			for _, id := range res.Patient.PID.PatientIdentifierList {
				if !strings.Contains(id.IDNumber, "^") && id.IDNumber != "" {
					specimen.Barcode = id.IDNumber
					slog.Info("Using first non-compound identifier as fallback", "barcode", specimen.Barcode)
					break
				}
			}
		}
	}

	// Log the final barcode being used
	slog.Info("Final specimen barcode for Wondfo", "barcode", specimen.Barcode)

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

	// For Wondfo, don't use OBR barcode - we'll use PID patient ID instead
	// This will be overridden in mapORUR01OrderObservationToSpecimenEntity
	return entity.Specimen{
		Barcode: "", // Will be set from PID segment
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
		// For Wondfo, parse the value and unit from the observation value
		// Example: "0.42 ng/mL" -> extract "0.42"
		valueStr := fmt.Sprintf("%v", values[i])
		// Split by space and take the first part (the numeric value)
		parts := strings.Fields(valueStr)
		if len(parts) > 0 {
			results = append(results, parts[0])
		} else {
			results = append(results, valueStr)
		}
	}
	return results
}

func (h *Handler) mapOBXToObservationResultEntity(obx *h251.OBX) entity.ObservationResult {
	if obx == nil {
		return entity.ObservationResult{}
	}

	testCode, description := h.getObservationIdentifier(obx.ObservationIdentifier)

	// Debug the test code mapping for Wondfo
	slog.Info("Wondfo OBX Debug",
		"test_code", testCode,
		"description", description,
		"observation_identifier", obx.ObservationIdentifier.Identifier,
		"observation_text", obx.ObservationIdentifier.Text,
		"observation_sub_id", obx.ObservationSubID)

	// For Wondfo, test code might be in ObservationSubID instead of ObservationIdentifier
	if testCode == "" && obx.ObservationSubID != "" {
		testCode = obx.ObservationSubID
		description = obx.ObservationSubID
		slog.Info("Wondfo: Using ObservationSubID as test code", "test_code", testCode)
	}

	// Extract unit from observation value for Wondfo
	// Example: "0.42 ng/mL" -> extract "ng/mL"
	unit := ""
	if len(obx.ObservationValue) > 0 {
		valueStr := fmt.Sprintf("%v", obx.ObservationValue[0])
		parts := strings.Fields(valueStr)
		if len(parts) > 1 {
			unit = parts[1] // Take the unit part
		}
		slog.Info("Wondfo Value Debug", "raw_value", valueStr, "parsed_unit", unit)
	}

	// If no unit extracted from value, use the Units field
	if unit == "" {
		unit = h.getUnits(obx.Units)
	}

	observationResult := entity.ObservationResult{
		TestCode:       testCode,
		Description:    description,
		Values:         h.mapObservationValueToValues(obx.ObservationValue),
		Type:           obx.ValueType,
		Unit:           unit,
		ReferenceRange: obx.ReferencesRange,
		Date:           obx.DateTimeOfTheObservation,
		AbnormalFlag:   obx.AbnormalFlags,
		Comments:       obx.ObservationResultStatus,
	}

	slog.Info("Wondfo Final ObservationResult",
		"test_code", observationResult.TestCode,
		"values", observationResult.Values,
		"unit", observationResult.Unit,
		"reference_range", observationResult.ReferenceRange)

	return observationResult
}

func (h *Handler) mapORUR01PatientToPatientEntity(p h251.ORU_R01_PatientResult) entity.Patient {
	if p.Patient == nil {
		return entity.Patient{}
	}

	// Use the standard patient mapping - barcode is handled in specimen
	patient := common.MapPIDToPatientEntity(p.Patient.PID)

	return patient
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
		VersionID:                           h251.VID{VersionID: "2.4"},
		SequenceNumber:                      "",
		ContinuationPointer:                 "",
		AcceptAcknowledgmentType:            "NE",
		ApplicationAcknowledgmentType:       "NE",
		CountryCode:                         "CHN",
		CharacterSet:                        []string{"UTF8"},
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
