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

// isValidBarcodePrefix checks if the string starts with any valid barcode prefix (case insensitive)
func (h *Handler) isValidBarcodePrefix(s string) bool {
	prefixes := []string{"SER", "WBL", "URI", "SEM", "PLM", "LIQ", "CSF"}
	upperS := strings.ToUpper(s)
	for _, prefix := range prefixes {
		if strings.HasPrefix(upperS, prefix) {
			return true
		}
	}
	return false
}

func (h *Handler) ORUR01(ctx context.Context, m h251.ORU_R01, msgByte []byte) (string, error) {
	msgControlID := m.MSH.MessageControlID

	slog.Info("Wondfo: Starting ORUR01 processing")

	oruR01, err := h.decodeORUR01(msgByte)
	if err != nil {
		slog.Error("Wondfo: Decode failed", "error", err)
		return "", fmt.Errorf("decode failed: %w", err)
	}

	// Log the parsed data structure for debugging
	slog.Info("Wondfo: Decoded ORU_R01",
		"patient_count", len(oruR01.Patient),
		"message_control_id", msgControlID)

	for i, patient := range oruR01.Patient {
		slog.Info("Wondfo: Patient data",
			"patient_index", i,
			"specimen_count", len(patient.Specimen))

		for j, specimen := range patient.Specimen {
			slog.Info("Wondfo: Specimen data",
				"patient_index", i,
				"specimen_index", j,
				"barcode", specimen.Barcode,
				"observation_count", len(specimen.ObservationResult))

			for k, obs := range specimen.ObservationResult {
				slog.Info("Wondfo: Observation data",
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
		slog.Error("Wondfo: Process failed", "error", err)
		return "", fmt.Errorf("process failed: %w", err)
	}

	slog.Info("Wondfo: ORU_R01 processing completed successfully")

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
	// Based on the HL7: PID|202510070001||||ser251007014||20251007|O
	// The patient ID "ser251007014" is in Patient Name field (field 5)
	// Barcode can have different prefixes: SER, WBL, URI, SEM, PLM, LIQ, CSF
	if res.Patient != nil && res.Patient.PID != nil {
		slog.Info("Debugging PID structure for Wondfo")

		// First, try Patient Name field (field 5) - this is where ser251007014 should be
		if len(res.Patient.PID.PatientName) > 0 {
			slog.Info("Checking Patient Name fields", "count", len(res.Patient.PID.PatientName))
			for i, name := range res.Patient.PID.PatientName {
				slog.Info("Patient Name field", "index", i, "family_name", name.FamilyName, "given_name", name.GivenName)

				// For Wondfo, the barcode might be in FamilyName as a complete string
				if h.isValidBarcodePrefix(name.FamilyName) {
					// Convert to uppercase to match database format (SERxxxxxxx vs serxxxxxxx)
					specimen.Barcode = strings.ToUpper(name.FamilyName)
					slog.Info("Found barcode in FamilyName", "original", name.FamilyName, "converted", specimen.Barcode)
					break
				}
				if h.isValidBarcodePrefix(name.GivenName) {
					// Convert to uppercase to match database format
					specimen.Barcode = strings.ToUpper(name.GivenName)
					slog.Info("Found barcode in GivenName", "original", name.GivenName, "converted", specimen.Barcode)
					break
				}

				// For Wondfo: Check if the name field contains the barcode as a single field
				// Sometimes the entire patient name field contains just the barcode
				fullName := strings.TrimSpace(name.FamilyName + name.GivenName)
				if fullName != "" && h.isValidBarcodePrefix(fullName) {
					// Convert to uppercase to match database format
					specimen.Barcode = strings.ToUpper(fullName)
					slog.Info("Found barcode in combined name", "original", fullName, "converted", specimen.Barcode)
					break
				}

				// Check if FamilyName contains the barcode even if it doesn't match prefix exactly (case insensitive)
				if name.FamilyName != "" {
					lowerFamily := strings.ToLower(name.FamilyName)
					if strings.HasPrefix(lowerFamily, "ser") || strings.HasPrefix(lowerFamily, "wbl") ||
						strings.HasPrefix(lowerFamily, "uri") || strings.HasPrefix(lowerFamily, "sem") ||
						strings.HasPrefix(lowerFamily, "plm") || strings.HasPrefix(lowerFamily, "liq") ||
						strings.HasPrefix(lowerFamily, "csf") {
						// Convert to uppercase to match database format
						specimen.Barcode = strings.ToUpper(name.FamilyName)
						slog.Info("Found barcode in FamilyName (case insensitive)", "original", name.FamilyName, "converted", specimen.Barcode)
						break
					}
				}
			}
		}

		// If not found in Patient Name, check Patient Identifier List (field 3)
		if specimen.Barcode == "" {
			slog.Info("Checking Patient Identifier List", "count", len(res.Patient.PID.PatientIdentifierList))
			for i, id := range res.Patient.PID.PatientIdentifierList {
				slog.Info("PID identifier", "index", i, "id_number", id.IDNumber, "assigning_authority", id.AssigningAuthority)

				if h.isValidBarcodePrefix(id.IDNumber) {
					// Convert to uppercase to match database format
					specimen.Barcode = strings.ToUpper(id.IDNumber)
					slog.Info("Found barcode in identifier list", "original", id.IDNumber, "converted", specimen.Barcode)
					break
				}

				// Check if the ID contains ^ separated values
				if strings.Contains(id.IDNumber, "^") {
					parts := strings.Split(id.IDNumber, "^")
					for j, part := range parts {
						slog.Info("Checking compound part", "part_index", j, "part_value", part)
						if h.isValidBarcodePrefix(part) {
							// Convert to uppercase to match database format
							specimen.Barcode = strings.ToUpper(part)
							slog.Info("Found barcode in compound field", "original", part, "converted", specimen.Barcode)
							break
						}
					}
					if specimen.Barcode != "" {
						break
					}
				}
			}
		}

		// Fallback - use first non-compound identifier or any identifier that looks like a barcode
		if specimen.Barcode == "" && len(res.Patient.PID.PatientIdentifierList) > 0 {
			for _, id := range res.Patient.PID.PatientIdentifierList {
				if !strings.Contains(id.IDNumber, "^") && id.IDNumber != "" {
					// Check if it looks like a barcode (contains letters and numbers)
					if len(id.IDNumber) > 3 {
						// Convert to uppercase to match database format
						specimen.Barcode = strings.ToUpper(id.IDNumber)
						slog.Info("Using identifier as fallback barcode", "original", id.IDNumber, "converted", specimen.Barcode)
						break
					}
				}
			}
		}
	}

	// Log the final barcode being used
	if specimen.Barcode == "" {
		slog.Error("Wondfo: No barcode found for specimen - this will cause database lookup to fail")
	} else {
		slog.Info("Wondfo: Final specimen barcode extracted", "barcode", specimen.Barcode)
	}

	observationResults := []entity.ObservationResult{}

	for _, o := range s.Observation {
		observationResult := h.mapOBXToObservationResultEntity(o.OBX)

		// Log each observation result before adding
		slog.Info("Wondfo: Adding observation result",
			"test_code", observationResult.TestCode,
			"values", observationResult.Values,
			"unit", observationResult.Unit)

		observationResults = append(observationResults, observationResult)
	}

	specimen.ObservationResult = observationResults

	slog.Info("Wondfo: Specimen mapping complete",
		"barcode", specimen.Barcode,
		"total_observations", len(observationResults))

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
		// For Wondfo, parse the value and remove unit and abnormal indicators
		// Example: "65.69 pmol/L↑" -> extract "65.69"
		valueStr := fmt.Sprintf("%v", values[i])

		// Split by space and take the first part (the numeric value)
		parts := strings.Fields(valueStr)
		if len(parts) > 0 {
			// Take the first part which should be the numeric value
			numericValue := parts[0]
			// Remove any trailing abnormal indicators that might be attached to the number
			numericValue = strings.TrimRight(numericValue, "↑↓+-<>")
			results = append(results, numericValue)
			slog.Info("Wondfo Value Parsing", "original", valueStr, "parsed", numericValue)
		} else {
			// Fallback: use the original value but clean it
			cleanValue := strings.TrimSpace(valueStr)
			cleanValue = strings.TrimRight(cleanValue, "↑↓+-<>")
			results = append(results, cleanValue)
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
		"observation_sub_id", obx.ObservationSubID,
		"value_type", obx.ValueType)

	// For Wondfo, test code might be in ObservationSubID instead of ObservationIdentifier
	if testCode == "" && obx.ObservationSubID != "" {
		testCode = obx.ObservationSubID
		description = obx.ObservationSubID
		slog.Info("Wondfo: Using ObservationSubID as test code", "test_code", testCode)
	}

	// For Wondfo: Sometimes the test code is in the ObservationIdentifier.Text field
	if testCode == "" && obx.ObservationIdentifier.Text != "" {
		testCode = obx.ObservationIdentifier.Text
		description = obx.ObservationIdentifier.Text
		slog.Info("Wondfo: Using ObservationIdentifier.Text as test code", "test_code", testCode)
	}

	// For Wondfo: Check if test code is in the observation value field 4
	// In the raw message: OBX|202510070001|NM||fT4|65.69 pmol/L↑|pmol/L|12-22
	// The fT4 is actually in field 4 (after the empty field 3)
	if testCode == "" {
		// Try to extract from raw observation value if it's text-based
		if len(obx.ObservationValue) > 0 {
			valueStr := fmt.Sprintf("%v", obx.ObservationValue[0])
			// Check if the value looks like a test code (contains letters, not just numbers)
			if strings.ContainsAny(valueStr, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz") {
				// Split to see if first part is test code
				parts := strings.Fields(valueStr)
				if len(parts) > 1 && !strings.ContainsAny(parts[0], "0123456789.,") {
					testCode = parts[0]
					description = parts[0]
					slog.Info("Wondfo: Extracted test code from observation value", "test_code", testCode)
				}
			}
		}
	}

	// Final fallback: use a default test code if still empty
	if testCode == "" {
		testCode = "UNKNOWN"
		description = "Unknown Test"
		slog.Warn("Wondfo: No test code found, using default", "test_code", testCode)
	}

	// Extract raw observation value for Wondfo parsing
	rawValue := ""
	unit := ""

	if len(obx.ObservationValue) > 0 {
		rawValue = fmt.Sprintf("%v", obx.ObservationValue[0])
		slog.Info("Wondfo Raw Value", "raw_value", rawValue)

		// Parse value and unit from observation value for Wondfo
		// Example: "65.69 pmol/L↑" -> extract "65.69" and "pmol/L"
		parts := strings.Fields(rawValue)
		if len(parts) > 1 {
			// The unit might contain special characters like ↑ or ↓
			unitPart := parts[1]
			// Remove abnormal indicators
			unit = strings.TrimRight(unitPart, "↑↓")
		}
	}

	// If no unit extracted from value, use the Units field
	if unit == "" {
		unit = h.getUnits(obx.Units)
	}

	// For Wondfo, clean the values to remove units and abnormal flags
	cleanValues := h.mapObservationValueToValues(obx.ObservationValue)

	observationResult := entity.ObservationResult{
		TestCode:       testCode,
		Description:    description,
		Values:         cleanValues,
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
		"reference_range", observationResult.ReferenceRange,
		"abnormal_flag", observationResult.AbnormalFlag)

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
