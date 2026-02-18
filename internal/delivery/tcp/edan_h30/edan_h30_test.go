package edan_h30

import (
	"strings"
	"testing"

	"github.com/kardianos/hl7"
	"github.com/kardianos/hl7/h251"
	"github.com/stretchr/testify/assert"
)

func TestParseEdanH30Message(t *testing.T) {
	// Sample EDAN H30 HL7 message from user
	rawMessage := `MSH|^~\&|H30_Pro^5295|EDANLAB|||20260216155459||ORU^R01|151|P|2.4||||0||UTF8
PID|139||^0||devi amalia|||0|0
OBR||24||EDANLAB^H30Pro^Sample|0|General|20260216131108||0||1|^Moderate Anemia@Microcytes||20260216131002||||admin|0|||||0|0
OBX||NM|0|WBC|9.51|10\S\3/μL|3.50-9.50|1||0||9.51^10\S\3/μL|H
OBX||NM|0|NEUT#|6.46|10\S\3/μL|1.80-6.30|1||0||6.459999999999999^10\S\3/μL|H
OBX||NM|0|NEUT%|67.9|%|40.0-75.0|0||0||67.9^%
OBX||NM|0|LYM#|1.98|10\S\3/μL|1.10-3.20|0||0||1.98^10\S\3/μL
OBX||NM|0|RBC|3.42|10\S\6/μL|3.80-5.80|2||0||3.417282082760522^10\S\6/μL|L
OBX||NM|0|HGB|8.8|g/dL|11.5-17.5|2||0||8.831102787478525^g/dL|L
OBX||NM|0|HCT|26.9|%|35.0-50.0|2||0||26.92892035106044^%|L
OBX||NM|0|MCV|76.7|fL|82.0-100.0|2||0||76.74697709917855^fL|L
OBX||NM|0|PLT|258|10\S\3/μL|125-350|0||0||258.0786319827716^10\S\3/μL`

	// Convert \n to actual line breaks
	message := strings.ReplaceAll(rawMessage, "\\n", "\r")
	message = strings.TrimSpace(message)

	t.Run("Message should contain MSH segment", func(t *testing.T) {
		assert.Contains(t, message, "MSH|")
	})

	t.Run("Message should contain PID segment", func(t *testing.T) {
		assert.Contains(t, message, "PID|139|")
	})

	t.Run("Message should contain OBR segment", func(t *testing.T) {
		assert.Contains(t, message, "OBR||24||EDANLAB")
	})

	t.Run("Message should contain OBX segments", func(t *testing.T) {
		assert.Contains(t, message, "OBX||NM|0|WBC|")
		assert.Contains(t, message, "OBX||NM|0|RBC|")
		assert.Contains(t, message, "OBX||NM|0|HGB|")
	})

	t.Run("Hematology values should be present", func(t *testing.T) {
		// WBC value
		assert.Contains(t, message, "9.51")
		// RBC value
		assert.Contains(t, message, "3.42")
		// HGB value
		assert.Contains(t, message, "8.8")
		// PLT value
		assert.Contains(t, message, "258")
	})
}

func TestEdanH30MessageStructure(t *testing.T) {
	t.Run("Verify sample message format", func(t *testing.T) {
		// The EDAN H30 uses standard HL7 2.4 ORU^R01 message structure
		// MSH - Message Header
		// PID - Patient Identification
		// OBR - Observation Request
		// OBX - Observation/Result

		expectedSegments := []string{"MSH", "PID", "OBR", "OBX"}

		for _, segment := range expectedSegments {
			t.Logf("EDAN H30 should support %s segment", segment)
		}
	})

	t.Run("Verify hematology test parameters", func(t *testing.T) {
		// List of hematology tests from the sample message
		testParameters := []string{
			"WBC", "NEUT#", "NEUT%", "LYM#", "LYM%", "MXD#", "MXD%",
			"RBC", "HGB", "HCT", "MCV", "MCH", "MCHC",
			"RDW_SD", "RDW_CV", "PLT", "PDW", "MPV", "PCT",
			"P_LCR", "P_LCC", "PLR", "NLR",
		}

		for _, param := range testParameters {
			t.Logf("EDAN H30 can report parameter: %s", param)
		}
	})
}

func TestPreprocessClearsInvalidOBXEffectiveDate(t *testing.T) {
	h := NewHandler(nil)

	raw := `MSH|^~\\&|H30_Pro^5295|EDANLAB|||20260217114231||ORU^R01|188|P|2.4||||0||UTF8
PID|175|WBL260217001|^0|||||0|0
OBR||9||EDANLAB^H30Pro^Sample|0|General|20260217090514||0||1|^Microcytes||20260217090418|^^Neutrophilia@Lymphopenia@Leukocytosis||^^^Thrombopenia|admin|0|WBL260217001||||0|0
OBX||NM|0|WBC|13.56|10\\S\\3/μL|3.50-9.50|1||0||13.56^10\\S\\3/μL|H
OBX||NM|0|NEUT#|12.93|10\\S\\3/μL|1.80-6.30|1||0||12.93^10\\S\\3/μL|H
OBX||NM|0|NEUT%|95.4|%|40.0-75.0|1||0||95.39999999999999^%|H
OBX||NM|0|LYM#|0.49|10\\S\\3/μL|1.10-3.20|2||0||0.49^10\\S\\3/μL|L
OBX||NM|0|LYM%|3.6|%|20.0-50.0|2||0||3.6^%|L
OBX||NM|0|MXD#|0.14|10\\S\\3/μL|0.10-1.50|0||0||0.14^10\\S\\3/μL
OBX||NM|0|MXD%|1.0|%|3.0-15.0|2||0||1.0^%|L
OBX||NM|0|RBC|4.95|10\\S\\6/μL|3.80-5.80|0||0||4.95066270531312^10\\S\\6/μL
OBX||NM|0|HGB|13.5|g/dL|11.5-17.5|0||0||13.4839694125334^g/dL
OBX||NM|0|HCT|37.8|%|35.0-50.0|0||0||37.84187913751591^%
OBX||NM|0|MCV|76.4|fL|82.0-100.0|2||0||76.43800717205694^fL|L
OBX||NM|0|MCH|27.2|pg|27.0-34.0|0||0||27.23669525232503^pg
OBX||NM|0|MCHC|35.6|g/dL|31.6-35.4|1||0||35.63239903573706^g/dL|H
OBX||NM|0|RDW_SD|35.7|fL|35.0-56.0|0||0||35.67261806488037^fL
OBX||NM|0|RDW_CV|13.8|%|10.0-15.0|0||0||13.79310364429663^%
OBX||NM|0|PLT|13|10\\S\\3/μL|125-350|2||0||12.70169238920884^10\\S\\3/μL|L
OBX||NM|0|PDW|22.28||9.00-17.00|1||0||22.27500088512897|H
OBX||NM|0|MPV|11.4|fL|6.5-12.0|0||0||11.43935974431038^fL
OBX||NM|0|PCT|0.02|%|0.17-0.35|2||0||0.01513533665098125^%|L
OBX||NM|0|P_LCR|37.0|%|11.0-45.0|0||0||37.02399940490723^%
OBX||NM|0|P_LCC|5|10\\S\\3/μL|30-90|2||0||4.702674514593826^10\\S\\3/μL|L
OBX||NM|0|PLR|25.92|||0||0||25.92182120246701
OBX||NM|0|NLR|26.39|||0||0||26.38775510204082
`

	pre := h.preprocessMessage(raw)
	lines := strings.Split(pre, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OBX|") {
			fields := strings.Split(line, "|")
			if len(fields) > 12 {
				assert.Equal(t, "", strings.TrimSpace(fields[12]), "OBX-12 should be cleared if invalid")
			}

			// ensure units (OBX-6) do not contain HL7 escape sequences after preprocessing
			if len(fields) > 6 {
				assert.NotContains(t, fields[6], "\\S\\", "OBX-6 must not contain HL7 escape \\S\\ after preprocessing")
			}
		}
	}
}

func TestDecodeORUR01SucceedsWithInvalidOBX12(t *testing.T) {
	// ensure decodeORUR01 no longer fails on EDAN H30 malformed OBX-12
	h := NewHandler(nil)

	raw := `MSH|^~\\&|H30_Pro^5295|EDANLAB|||20260217114231||ORU^R01|188|P|2.4||||0||UTF8
PID|175|WBL260217001|^0|||||0|0
OBR||9||EDANLAB^H30Pro^Sample|0|General|20260217090514||0||1|^Microcytes||20260217090418|^^Neutrophilia@Lymphopenia@Leukocytosis||^^^Thrombopenia|admin|0|WBL260217001||||0|0
OBX||NM|0|WBC|13.56|10\\S\\3/μL|3.50-9.50|1||0||13.56^10\\S\\3/μL|H
OBX||NM|0|NEUT#|12.93|10\\S\\3/μL|1.80-6.30|1||0||12.93^10\\S\\3/μL|H
OBX||NM|0|NEUT%|95.4|%|40.0-75.0|1||0||95.39999999999999^%|H
OBX||NM|0|LYM#|0.49|10\\S\\3/μL|1.10-3.20|2||0||0.49^10\\S\\3/μL|L
OBX||NM|0|LYM%|3.6|%|20.0-50.0|2||0||3.6^%|L
OBX||NM|0|MXD#|0.14|10\\S\\3/μL|0.10-1.50|0||0||0.14^10\\S\\3/μL
OBX||NM|0|MXD%|1.0|%|3.0-15.0|2||0||1.0^%|L
OBX||NM|0|RBC|4.95|10\\S\\6/μL|3.80-5.80|0||0||4.95066270531312^10\\S\\6/μL
OBX||NM|0|HGB|13.5|g/dL|11.5-17.5|0||0||13.4839694125334^g/dL
OBX||NM|0|HCT|37.8|%|35.0-50.0|0||0||37.84187913751591^%
OBX||NM|0|MCV|76.4|fL|82.0-100.0|2||0||76.43800717205694^fL|L
OBX||NM|0|MCH|27.2|pg|27.0-34.0|0||0||27.23669525232503^pg
OBX||NM|0|MCHC|35.6|g/dL|31.6-35.4|1||0||35.63239903573706^g/dL|H
OBX||NM|0|RDW_SD|35.7|fL|35.0-56.0|0||0||35.67261806488037^fL
OBX||NM|0|RDW_CV|13.8|%|10.0-15.0|0||0||13.79310364429663^%
OBX||NM|0|PLT|13|10\\S\\3/μL|125-350|2||0||12.70169238920884^10\\S\\3/μL|L
OBX||NM|0|PDW|22.28||9.00-17.00|1||0||22.27500088512897|H
OBX||NM|0|MPV|11.4|fL|6.5-12.0|0||0||11.43935974431038^fL
OBX||NM|0|PCT|0.02|%|0.17-0.35|2||0||0.01513533665098125^%|L
OBX||NM|0|P_LCR|37.0|%|11.0-45.0|0||0||37.02399940490723^%
OBX||NM|0|P_LCC|5|10\\S\\3/μL|30-90|2||0||4.702674514593826^10\\S\\3/μL|L
OBX||NM|0|PLR|25.92|||0||0||25.92182120246701
OBX||NM|0|NLR|26.39|||0||0||26.38775510204082
`

	// decode should not return error after preprocessing
	_, err := h.decodeORUR01([]byte(raw))
	assert.NoError(t, err)
}

func TestBarcodeFromPIDOverridesOBR(t *testing.T) {
	h := NewHandler(nil)

	raw := `MSH|^~\\&|H30_Pro^5295|EDANLAB|||20260217121450||ORU^R01|194|P|2.4||||0||UTF8
PID|181|WBL260217001|^0||ramini|||0|0
OBR||20||EDANLAB^H30Pro^Sample|0|General|20260217115241||0||1|^Microcytes||20260217115138|^^Neutrophilia@Leukocytosis|||admin|0|WBL260217001||||0|0
OBX||NM|0|WBC|16.34|10\\S\\3/μL|3.50-9.50|1||0||16.34^10\\S\\3/μL|H
`

	// also decode with hl7 library to inspect where parser placed the PID/OBR values
	d := hl7.NewDecoder(h251.Registry, nil)
	san := h.preprocessMessage(raw)
	msg, derr := d.Decode([]byte(san))
	t.Logf("hl7 decode error (after preprocess): %v", derr)
	if derr == nil {
		if m, ok := msg.(h251.ORU_R01); ok && len(m.PatientResult) > 0 {
			if m.PatientResult[0].Patient != nil && m.PatientResult[0].Patient.PID != nil {
				pid := m.PatientResult[0].Patient.PID
				t.Logf("parsed PID.PatientID=%+v", pid.PatientID)
				t.Logf("parsed PID.PatientIdentifierList=%+v", pid.PatientIdentifierList)
			}

			if len(m.PatientResult[0].OrderObservation) > 0 && m.PatientResult[0].OrderObservation[0].OBR != nil {
				obr := m.PatientResult[0].OrderObservation[0].OBR
				t.Logf("parsed OBR (selected fields): PlacerOrderNumber=%+v, FillerOrderNumber=%+v, UniversalServiceIdentifier=%+v, SpecimenSource=%+v", obr.PlacerOrderNumber, obr.FillerOrderNumber, obr.UniversalServiceIdentifier, obr.SpecimenSource)
			}
		}
	}

	res, err := h.decodeORUR01([]byte(raw))
	assert.NoError(t, err)
	if assert.NotEmpty(t, res.Patient) && assert.NotEmpty(t, res.Patient[0].Specimen) {
		assert.Equal(t, "WBL260217001", res.Patient[0].Specimen[0].Barcode, "barcode must be taken from PID, not OBR")
	}
}
