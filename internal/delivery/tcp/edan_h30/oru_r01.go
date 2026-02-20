package edan_h30

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"
	"unicode"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp/common"
	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
)

func (h *Handler) ORUR01(ctx context.Context, m h251.ORU_R01, msgByte []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID

	slog.Info("EDAN H30: Starting ORUR01 processing")

	oruR01, err := h.decodeORUR01(msgByte)
	if err != nil {
		slog.Error("EDAN H30: Decode failed", "error", err)
		return "", fmt.Errorf("decode failed: %w", err)
	}

	// Log the parsed data structure for debugging
	slog.Info("EDAN H30: Decoded ORU_R01",
		"patient_count", len(oruR01.Patient),
		"message_control_id", msgControlID)

	for i, patient := range oruR01.Patient {
		slog.Info("EDAN H30: Patient data",
			"patient_index", i,
			"specimen_count", len(patient.Specimen))

		for j, specimen := range patient.Specimen {
			slog.Info("EDAN H30: Specimen data",
				"patient_index", i,
				"specimen_index", j,
				"barcode", specimen.Barcode,
				"observation_count", len(specimen.ObservationResult))

			for k, obs := range specimen.ObservationResult {
				slog.Info("EDAN H30: Observation data",
					"patient_index", i,
					"specimen_index", j,
					"observation_index", k,
					"test_code", obs.TestCode,
					"description", obs.Description,
					"values", obs.Values,
					"unit", obs.Unit,
					"reference_range", obs.ReferenceRange)
			}
		}
	}

	err = h.analyzerUsecase.ProcessORUR01(ctx, oruR01)
	if err != nil {
		slog.Error("EDAN H30: Process failed", "error", err)
		return "", fmt.Errorf("process failed: %w", err)
	}

	slog.Info("EDAN H30: ORU_R01 processing completed successfully")

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
	// Pre-process message to fix EDAN H30 specific issues
	// EDAN H30 puts invalid data in datetime fields which causes parsing errors
	messageStr := h.preprocessMessage(string(msgByte))

	d := hl7.NewDecoder(h251.Registry, nil)
	msg, err := d.Decode([]byte(messageStr))
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

// preprocessMessage cleans EDAN H30 specific issues in HL7 message
func (h *Handler) preprocessMessage(message string) string {
	// normalize different segment separators to LF for predictable splitting
	message = strings.ReplaceAll(message, "\r\n", "\n")
	message = strings.ReplaceAll(message, "\r", "\n")
	lines := strings.Split(message, "\n")
	var cleanedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// OBR: sanitize RequestedDateTime (field 7, index 6)
		if strings.HasPrefix(line, "OBR|") {
			fields := strings.SplitN(line, "|", 12)
			if len(fields) > 6 {
				datetime := strings.TrimSpace(fields[6])
				// Clear if it contains letters or clearly not a datetime starting with '20'
				hasLetter := false
				for _, r := range datetime {
					if unicode.IsLetter(r) {
						hasLetter = true
						break
					}
				}
				if datetime != "" && (hasLetter || !strings.HasPrefix(datetime, "20")) {
					slog.Warn("EDAN H30: clearing invalid OBR.RequestedDateTime", "original", datetime)
					fields[6] = ""
					line = strings.Join(fields, "|")
				}
			}
		}

		// OBX: 1) escape carets in ReferenceRange (field 7) 2) validate/clear EffectiveDateOfReferenceRange (field 12)
		if strings.HasPrefix(line, "OBX|") {
			fields := strings.SplitN(line, "|", 13)

			// Escape caret in OBX-7 if present (avoid component-sep errors)
			if len(fields) > 7 && fields[7] != "" && strings.Contains(fields[7], "^") {
				old := fields[7]
				fields[7] = strings.ReplaceAll(fields[7], "^", "\\S\\")
				slog.Info("EDAN H30: escaped OBX-7 ReferenceRange", "original", old, "escaped", fields[7])
			}

			// Normalize OBX-6 (Units): some instruments send HL7 escape \S\ inside units
			// which the HL7 decoder rejects when parsing CE strings — convert to a
			// safe textual representation (e.g. 10\S\3/μL -> 10e3/μL).
			if len(fields) > 6 && fields[6] != "" {
				old := fields[6]
				// handle both single- and double-escaped variants coming from different sources
				u := strings.ReplaceAll(fields[6], "\\\\S\\\\", "e") // "\\S\\" in raw -> "\\\\S\\\\" literal
				u = strings.ReplaceAll(u, "\\S\\", "e")
				// collapse any remaining backslashes (safe for units)
				u = strings.ReplaceAll(u, "\\", "")
				fields[6] = u
				slog.Info("EDAN H30: normalized OBX-6 units", "original", old, "normalized", fields[6])
			}

			// Validate OBX-12 (field index 12 in the split slice) — must look like HL7 datetime (digits only, length 4/6/8/12/14)
			// NOTE: fields[0] == "OBX", so OBX-12 corresponds to fields[12]
			if len(fields) > 12 && strings.TrimSpace(fields[12]) != "" {
				v := strings.TrimSpace(fields[12])
				isLikelyDate := false
				switch len(v) {
				case 4, 6, 8, 12, 14:
					allDigits := true
					for _, r := range v {
						if r < '0' || r > '9' {
							allDigits = false
							break
						}
					}
					if allDigits {
						isLikelyDate = true
					}
				default:
					isLikelyDate = false
				}

				if !isLikelyDate {
					slog.Warn("EDAN H30: clearing invalid OBX.EffectiveDateOfReferenceRange", "value", v)
					fields[12] = ""
				}
			}

			line = strings.Join(fields, "|")
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
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

	// For EDAN H30, prefer PID-based barcode (PID can contain the real sample id; OBR often
	// contains analyzer/service identifiers such as "EDANLAB"). Try PID-3, then PID-2,
	// then patient name; fallback to OBR only if PID does not provide a usable value.
	if res.Patient != nil && res.Patient.PID != nil {
		slog.Info("EDAN H30: Debugging PID structure")

		pid := res.Patient.PID

		// 1) PatientIdentifierList (PID-3)
		if len(pid.PatientIdentifierList) > 0 {
			for i, id := range pid.PatientIdentifierList {
				slog.Info("EDAN H30: Patient Identifier List",
					"index", i,
					"id", id.IDNumber,
					"assigning_authority", id.AssigningAuthority)

				if id.IDNumber != "" && id.IDNumber != "0" {
					specimen.Barcode = id.IDNumber
					slog.Info("EDAN H30: Using Patient Identifier as barcode", "barcode", specimen.Barcode)
					break
				}
			}
		}

		// 2) PatientID (PID-2) — prefer PID value even if OBR already provided an identifier
		if pid.PatientID.IDNumber != "" && pid.PatientID.IDNumber != "0" {
			specimen.Barcode = pid.PatientID.IDNumber
			slog.Info("EDAN H30: Using PID.PatientID as barcode", "barcode", specimen.Barcode)
		}

		// 3) fallback to Patient Name
		if specimen.Barcode == "" && len(pid.PatientName) > 0 {
			for i, name := range pid.PatientName {
				slog.Info("EDAN H30: Patient Name",
					"index", i,
					"family_name", name.FamilyName,
					"given_name", name.GivenName)

				fullName := name.FamilyName
				if fullName != "" && fullName != "0" {
					specimen.Barcode = fullName
					slog.Info("EDAN H30: Using Family Name as barcode", "barcode", specimen.Barcode)
					break
				}
			}
		}
	}

	// Observation Mapping (OBX Segment)
	var observationResult []entity.ObservationResult
	for _, observation := range s.Observation {
		ob := h.mapOBXToObservationEntity(observation.OBX)
		observationResult = append(observationResult, ob)
	}
	specimen.ObservationResult = observationResult

	slog.Info("EDAN H30: Specimen mapping complete",
		"barcode", specimen.Barcode,
		"total_observations", len(observationResult))

	return specimen
}

func (h *Handler) mapOBRToSpecimenEntity(obr *h251.OBR) entity.Specimen {
	if obr == nil {
		return entity.Specimen{}
	}

	// For EDAN H30, prefer barcode-like identifiers commonly populated by the
	// analyzer/host: 1) FillerOrderNumber (often contains sample barcode),
	// 2) PlacerOrderNumber, 3) UniversalServiceIdentifier as a last resort.
	barcode := ""
	if obr.FillerOrderNumber != nil && obr.FillerOrderNumber.EntityIdentifier != "" {
		barcode = obr.FillerOrderNumber.EntityIdentifier
	} else if obr.PlacerOrderNumber != nil && obr.PlacerOrderNumber.EntityIdentifier != "" {
		barcode = obr.PlacerOrderNumber.EntityIdentifier
	} else if obr.UniversalServiceIdentifier != (h251.CE{}) {
		barcode = obr.UniversalServiceIdentifier.Identifier
	}

	return entity.Specimen{
		Barcode: barcode,
	}
}

func (h *Handler) mapOBXToObservationEntity(obx *h251.OBX) entity.ObservationResult {
	if obx == nil {
		return entity.ObservationResult{}
	}

	var obs entity.ObservationResult

	// OBX-3: Observation Identifier
	if obx.ObservationIdentifier != (h251.CE{}) {
		obs.TestCode = obx.ObservationIdentifier.Identifier
		obs.Description = obx.ObservationIdentifier.Text
	}

	// OBX-4: Observation Sub-ID - use this as test code if main identifier is empty or "0"
	if obs.TestCode == "" || obs.TestCode == "0" {
		obs.TestCode = obx.ObservationSubID
		obs.Description = obx.ObservationSubID
	}

	// OBX-5: Observation Value
	var values entity.JSONStringArray
	if len(obx.ObservationValue) > 0 {
		for i := range obx.ObservationValue {
			values = append(values, fmt.Sprintf("%v", obx.ObservationValue[i]))
		}
	}
	obs.Values = values

	// OBX-6: Units
	if obx.Units != nil {
		obs.Unit = obx.Units.Identifier
	}

	// OBX-7: Reference Range
	obs.ReferenceRange = obx.ReferencesRange

	// OBX-8: Abnormal Flags
	obs.AbnormalFlag = obx.AbnormalFlags

	// OBX-11: Observation Result Status
	obs.Comments = obx.ObservationResultStatus

	slog.Info("EDAN H30: Mapped OBX to Observation",
		"test_code", obs.TestCode,
		"description", obs.Description,
		"value", obs.Values,
		"unit", obs.Unit,
		"reference_range", obs.ReferenceRange)

	return obs
}

func (h *Handler) mapORUR01PatientToPatientEntity(res h251.ORU_R01_PatientResult) entity.Patient {
	if res.Patient == nil || res.Patient.PID == nil {
		slog.Warn("EDAN H30: No patient data in ORU_R01_PatientResult")
		return entity.Patient{}
	}

	// Use the standard patient mapping
	patient := common.MapPIDToPatientEntity(res.Patient.PID)

	slog.Info("EDAN H30: Mapped PID to Patient",
		"first_name", patient.FirstName,
		"last_name", patient.LastName,
		"sex", patient.Sex)

	return patient
}

func (h *Handler) createMSHAck(msh entity.MSH, msgControlID h251.ST) *h251.MSH {
	return &h251.MSH{
		HL7:                  h251.HL7Name{},
		FieldSeparator:       "|",
		EncodingCharacters:   "^~\\&",
		SendingApplication:   common.SimpleHD(msh.ReceivingApplication),
		SendingFacility:      common.SimpleHD(msh.ReceivingFacility),
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
		CountryCode:                         "ID",
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
}
