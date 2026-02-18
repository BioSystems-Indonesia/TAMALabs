package edan_i15

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/mllp/common"
	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
)

func (h *Handler) ORUR01(ctx context.Context, m h251.ORU_R01, msgByte []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID

	slog.Info("EDAN I15: Starting ORUR01 processing")

	oruR01, err := h.decodeORUR01(msgByte)
	if err != nil {
		slog.Error("EDAN I15: Decode failed", "error", err)
		return "", fmt.Errorf("decode failed: %w", err)
	}

	// Log the parsed data structure for debugging
	slog.Info("EDAN I15: Decoded ORU_R01",
		"patient_count", len(oruR01.Patient),
		"message_control_id", msgControlID)

	for i, patient := range oruR01.Patient {
		slog.Info("EDAN I15: Patient data",
			"patient_index", i,
			"specimen_count", len(patient.Specimen))

		for j, specimen := range patient.Specimen {
			slog.Info("EDAN I15: Specimen data",
				"patient_index", i,
				"specimen_index", j,
				"barcode", specimen.Barcode,
				"observation_count", len(specimen.ObservationResult))

			for k, obs := range specimen.ObservationResult {
				slog.Info("EDAN I15: Observation data",
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
		slog.Error("EDAN I15: Process failed", "error", err)
		return "", fmt.Errorf("process failed: %w", err)
	}

	slog.Info("EDAN I15: ORU_R01 processing completed successfully")

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
	// Pre-process message to fix EDAN I15 specific issues
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

// preprocessMessage cleans EDAN I15 specific issues in HL7 message
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

		// Fix OBR segment - field 7 (RequestedDateTime) sometimes contains invalid data
		if strings.HasPrefix(line, "OBR|") {
			fields := strings.Split(line, "|")
			if len(fields) > 6 {
				datetime := fields[6]
				if len(datetime) > 0 && !strings.HasPrefix(strings.TrimSpace(datetime), "20") {
					fields[6] = ""
					line = strings.Join(fields, "|")
				}
			}
		}

		// Fix OBX segment - escape caret in reference range (field 7) and clear field 12
		if strings.HasPrefix(line, "OBX|") {
			fields := strings.SplitN(line, "|", 20)

			// Field 7 (index 7) - Reference Range
			// EDAN I15 sends ranges like "136^145" which must be escaped to "136\S\145"
			// OBX-7: convert caret-separated ranges to hyphen for safe decoding (decoder rejects raw '^' or escape sequences)
			if len(fields) > 7 && strings.TrimSpace(fields[7]) != "" && strings.Contains(fields[7], "^") {
				orig := strings.TrimSpace(fields[7])
				parts := strings.Split(orig, "^")
				if len(parts) == 2 && isNumericString(parts[0]) && isNumericString(parts[1]) {
					fields[7] = strings.TrimSpace(parts[0]) + "-" + strings.TrimSpace(parts[1])
				} else {
					// fallback — replace carets with hyphen so hl7 parser won't reject the field
					fields[7] = strings.ReplaceAll(orig, "^", "-")
				}
				slog.Info("EDAN I15: normalized OBX-7 ReferenceRange for decoding", "original", orig, "normalized", fields[7])
			}

			// Normalize OBX-6 (Units) like EDAN H30: convert HL7 escape \S\ -> 'e', remove stray backslashes
			if len(fields) > 6 && fields[6] != "" {
				oldu := fields[6]
				u := strings.ReplaceAll(fields[6], "\\\\S\\\\", "e")
				u = strings.ReplaceAll(u, "\\S\\", "e")
				u = strings.ReplaceAll(u, "\\", "")
				fields[6] = u
				if oldu != fields[6] {
					slog.Info("EDAN I15: normalized OBX-6 units", "original", oldu, "normalized", fields[6])
				}
			}

			// Field 12 (index 11) - sometimes contains invalid data
			if len(fields) > 11 {
				fields[11] = ""
			}

			line = strings.Join(fields, "|")
		}

		cleanedLines = append(cleanedLines, line)
	}

	return strings.Join(cleanedLines, "\n")
}

// isNumericString returns true when s can be parsed as a float (allowing
// decimal point). Used to decide whether a caret-separated pair looks like
// a numeric range that we can safely convert to 'low-high'.
func isNumericString(s string) bool {
	s = strings.TrimSpace(s)
	if s == "" {
		return false
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}
	return false
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

	// For EDAN I15, prefer PID-based barcode (mirror EDAN H30 behaviour):
	// 1) PID.PatientIdentifierList (PID-3)
	// 2) PID.PatientID (PID-2) — overrides OBR/PID-3 if present
	// 3) Patient family name (fallback)
	if res.Patient != nil && res.Patient.PID != nil {
		slog.Info("EDAN I15: Debugging PID structure")
		pid := res.Patient.PID

		// 1) PatientIdentifierList (PID-3)
		if len(pid.PatientIdentifierList) > 0 {
			for i, id := range pid.PatientIdentifierList {
				slog.Info("EDAN I15: Patient Identifier List",
					"index", i,
					"id", id.IDNumber,
					"assigning_authority", id.AssigningAuthority)

				if id.IDNumber != "" && id.IDNumber != "0" {
					specimen.Barcode = id.IDNumber
					slog.Info("EDAN I15: Using Patient Identifier as barcode", "barcode", specimen.Barcode)
					break
				}
			}
		}

		// 2) PatientID (PID-2) — prefer/override if present
		if pid.PatientID.IDNumber != "" && pid.PatientID.IDNumber != "0" {
			specimen.Barcode = pid.PatientID.IDNumber
			slog.Info("EDAN I15: Using PID.PatientID as barcode", "barcode", specimen.Barcode)
		}

		// 3) fallback to Patient Name
		if specimen.Barcode == "" && len(pid.PatientName) > 0 {
			for i, name := range pid.PatientName {
				slog.Info("EDAN I15: Patient Name",
					"index", i,
					"family_name", name.FamilyName,
					"given_name", name.GivenName)

				fullName := name.FamilyName
				if fullName != "" && fullName != "0" {
					specimen.Barcode = fullName
					slog.Info("EDAN I15: Using Family Name as barcode", "barcode", specimen.Barcode)
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

	return specimen
}

func (h *Handler) mapOBRToSpecimenEntity(obr *h251.OBR) entity.Specimen {
	if obr == nil {
		return entity.Specimen{}
	}

	// For EDAN I15, try common OBR locations for a barcode-like id as fallback:
	// FillerOrderNumber -> PlacerOrderNumber -> UniversalServiceIdentifier
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
	// Present EDAN I15 numeric ranges back in the same escaped form as EDAN H30
	rr := obx.ReferencesRange
	if rr != "" && strings.Contains(rr, "-") {
		parts := strings.Split(rr, "-")
		if len(parts) == 2 && isNumericString(parts[0]) && isNumericString(parts[1]) {
			rr = strings.TrimSpace(parts[0]) + "\\S\\" + strings.TrimSpace(parts[1])
		}
	}
	obs.ReferenceRange = rr

	// OBX-8: Abnormal Flags
	obs.AbnormalFlag = obx.AbnormalFlags

	// OBX-11: Observation Result Status
	obs.Comments = obx.ObservationResultStatus

	slog.Info("EDAN I15: Mapped OBX to Observation",
		"test_code", obs.TestCode,
		"description", obs.Description,
		"value", obs.Values,
		"unit", obs.Unit,
		"reference_range", obs.ReferenceRange)

	return obs
}

func (h *Handler) mapORUR01PatientToPatientEntity(res h251.ORU_R01_PatientResult) entity.Patient {
	if res.Patient == nil || res.Patient.PID == nil {
		slog.Warn("EDAN I15: No patient data in ORU_R01_PatientResult")
		return entity.Patient{}
	}

	// Use the standard patient mapping
	patient := common.MapPIDToPatientEntity(res.Patient.PID)

	slog.Info("EDAN I15: Mapped PID to Patient",
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
