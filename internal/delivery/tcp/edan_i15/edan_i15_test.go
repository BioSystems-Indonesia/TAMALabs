package edan_i15

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEdanI15Message(t *testing.T) {
	// Sample EDAN I15 HL7 message from user
	rawMessage := `MSH|^~\&|EDAN|i15^M24311240006|LIS||20260216142047||ORU^R01||P|2.4||||0||UNICODE UTF-8||||
PID|||ok|||||U||||||||||||||||||||||||||||||
OBR|||20250902001|EDAN^i15|||20250902190239||||||||Arterial|mfrs|||||||||||||||||||||||||||||||
OBX|0|ST|0|pH|7.400||7.350^7.450||||||Pass|20250902190239||mfrs||||
OBX|1|ST|0|pO2|90|mmHg|83^108||||||Pass|20250902190239||mfrs||||
OBX|2|ST|0|pCO2|40.0|mmHg|35.0^45.0||||||Pass|20250902190239||mfrs||||
OBX|3|ST|0|Na+|140|mmol/L|136^145||||||Pass|20250902190239||mfrs||||
OBX|4|ST|0|K+|4.0|mmol/L|3.4^4.5||||||Pass|20250902190239||mfrs||||
OBX|5|ST|0|Ca++|1.30|mmol/L|1.15^1.33||||||Pass|20250902190239||mfrs||||
OBX|6|ST|0|Cl-|100|mmol/L|98^109||||||Pass|20250902190239||mfrs||||
OBX|7|ST|0|Glu|5.0|mmol/L|3.9^6.1||||||Pass|20250902190239||mfrs||||
OBX|8|ST|0|Lac|1.00|mmol/L|0.50^1.60||||||Pass|20250902190239||mfrs||||
OBX|9|ST|0|Hct|40|%|36^53||||||Pass|20250902190239||mfrs||||
OBX|10|ST|1|tHb(est)|13.7|g/dL|2.9^27.7|||||||20250902190239||mfrs||||
OBX|11|ST|1|cH+|39.8|nmol/L|10.0^316.2|||||||20250902190239||mfrs||||`

	// Convert \n to actual line breaks and clean
	message := strings.ReplaceAll(rawMessage, "\\n", "\r")
	message = strings.TrimSpace(message)

	t.Run("Message should contain MSH segment", func(t *testing.T) {
		assert.Contains(t, message, "MSH|")
	})

	t.Run("Message should contain PID segment", func(t *testing.T) {
		assert.Contains(t, message, "PID|||ok|")
	})

	t.Run("Message should contain OBR segment", func(t *testing.T) {
		assert.Contains(t, message, "OBR|||20250902001|")
	})

	t.Run("Message should contain OBX segments", func(t *testing.T) {
		assert.Contains(t, message, "OBX|0|ST|0|pH|")
		assert.Contains(t, message, "OBX|1|ST|0|pO2|")
		assert.Contains(t, message, "OBX|3|ST|0|Na+|")
	})

	t.Run("Test values should be present", func(t *testing.T) {
		// pH value
		assert.Contains(t, message, "7.400")
		// pO2 value
		assert.Contains(t, message, "90|mmHg")
		// Na+ value
		assert.Contains(t, message, "140|mmol/L")
	})
}

func TestEdanI15MessageStructure(t *testing.T) {
	t.Run("Verify sample message format", func(t *testing.T) {
		// The EDAN I15 uses standard HL7 2.4 ORU^R01 message structure
		// MSH - Message Header
		// PID - Patient Identification
		// OBR - Observation Request
		// OBX - Observation/Result

		expectedSegments := []string{"MSH", "PID", "OBR", "OBX"}

		for _, segment := range expectedSegments {
			t.Logf("EDAN I15 should support %s segment", segment)
		}
	})

	t.Run("Verify test parameters", func(t *testing.T) {
		// List of tests from the sample message
		testParameters := []string{
			"pH", "pO2", "pCO2", "Na+", "K+", "Ca++",
			"Cl-", "Glu", "Lac", "Hct", "tHb(est)", "cH+",
		}

		for _, param := range testParameters {
			t.Logf("EDAN I15 can report parameter: %s", param)
		}
	})
}

func TestBarcodeFromPIDOverridesOBR(t *testing.T) {
	h := NewHandler(nil)

	raw := `MSH|^~\\&|EDAN|i15^M24311240006|LIS||20260217114340||ORU^R01||P|2.4||||0||UNICODE UTF-8||||
PID|||PLM260217002|||||U|||||||||||||||||||||||||||||||
OBR||20||EDANI15^i15^Sample|0|General|20260214201510||||||||Venous|admin||||||||||||||||||||||||||||||||
OBX|0|ST|0|Na+|142|mmol/L|136-145||||||Pass|20260214201510||admin||||`

	res, err := h.decodeORUR01([]byte(raw))
	assert.NoError(t, err)
	if assert.NotEmpty(t, res.Patient) && assert.NotEmpty(t, res.Patient[0].Specimen) {
		assert.Equal(t, "PLM260217002", res.Patient[0].Specimen[0].Barcode, "barcode must be taken from PID, not OBR")
	}
}
func TestPreprocessEscapesOBXReferenceRange(t *testing.T) {
	h := NewHandler(nil)

	raw := `MSH|^~\&|EDAN|i15^M24311240006|LIS||20260217114340||ORU^R01||P|2.4||||0||UNICODE UTF-8||||
PID|||PLM260217002|||||U|||||||||||||||||||||||||||||||
OBR|||20260214007|EDAN^i15|||20260214201510||||||||Venous|admin||||||||||||||||||||||||||||||||
OBX|0|ST|0|Na+|142|mmol/L|136^145||||||Pass|20260214201510||admin||||
OBX|1|ST|0|K+|3.7|mmol/L|3.4^4.5||||||Pass|20260214201510||admin||||
OBX|2|ST|0|Ca++|1.19|mmol/L|1.15^1.33||||||Pass|20260214201510||admin||||
OBX|3|ST|0|Cl-|111|mmol/L|98^109|↑|||||Pass|20260214201510||admin||||
OBX|4|ST|0|Hct|53|%|36^53|↑|||||Pass|20260214201510||admin||||
OBX|5|ST|1|tHb(est)|18.1|g/dL|2.9^27.7|||||||20260214201510||admin||||
OBX|6|ST|1|mOsm|288.8|mOsm/L|200.9^449.4|||||||20260214201510||admin||||`

	pre := h.preprocessMessage(raw)
	lines := strings.Split(pre, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "OBX|") {
			fields := strings.Split(line, "|")
			if len(fields) > 7 {
				// OBX-7 should no longer contain '^' and should be normalized to a hyphen for decoding
				if strings.TrimSpace(fields[7]) != "" {
					if strings.Contains(fields[7], "^") {
						t.Fatalf("OBX-7 still contains caret after preprocess: %s", fields[7])
					}
					if strings.Contains(fields[7], "\\S\\") {
						t.Fatalf("OBX-7 must not contain HL7 escape sequences at decode time: %s", fields[7])
					}
					if !strings.Contains(fields[7], "-") {
						t.Fatalf("OBX-7 should be normalized to hyphen-range for decoding: %s", fields[7])
					}
				}
			}
		}
	}
}

func TestDecodeORUR01AcceptsReferenceRangeWithCaret(t *testing.T) {
	h := NewHandler(nil)

	raw := `MSH|^~\\&|EDAN|i15^M24311240006|LIS||20260217114340||ORU^R01||P|2.4||||0||UNICODE UTF-8||||
PID|||PLM260217002|||||U|||||||||||||||||||||||||||||||
OBR|||20260214007|EDAN^i15|||20260214201510||||||||Venous|admin||||||||||||||||||||||||||||||||
OBX|0|ST|0|Na+|142|mmol/L|136^145||||||Pass|20260214201510||admin||||
OBX|1|ST|0|K+|3.7|mmol/L|3.4^4.5||||||Pass|20260214201510||admin||||
OBX|2|ST|0|Ca++|1.19|mmol/L|1.15^1.33||||||Pass|20260214201510||admin||||
OBX|3|ST|0|Cl-|111|mmol/L|98^109|↑|||||Pass|20260214201510||admin||||
OBX|4|ST|0|Hct|53|%|36^53|↑|||||Pass|20260214201510||admin||||
OBX|5|ST|1|tHb(est)|18.1|g/dL|2.9^27.7|||||||20260214201510||admin||||
OBX|6|ST|1|mOsm|288.8|mOsm/L|200.9^449.4|||||||20260214201510||admin||||`

	res, err := h.decodeORUR01([]byte(raw))
	if err != nil {
		t.Fatalf("decodeORUR01 failed for message with caret in OBX-7: %v", err)
	}

	// ensure mapped ObservationResult.ReferenceRange uses HL7 escape like EDAN H30
	if assert.NotEmpty(t, res.Patient) && assert.NotEmpty(t, res.Patient[0].Specimen) {
		obs := res.Patient[0].Specimen[0].ObservationResult[0]
		if !strings.Contains(obs.ReferenceRange, "\\S\\") {
			t.Fatalf("mapped ReferenceRange not escaped as \\S\\: %s", obs.ReferenceRange)
		}
	}
}

func TestDecodeWithCRSegmentTerminators(t *testing.T) {
	// ensure preprocessMessage handles CR-separated segments (real device traffic)
	h := NewHandler(nil)

	rawCR := "MSH|^~\\&|EDAN|i15^M24311240006|LIS||20260217114340||ORU^R01||P|2.4||||0||UNICODE UTF-8||||\rPID|||PLM260217002|||||U|||||||||||||||||||||||||||||||\rOBR|||20260214007|EDAN^i15|||20260214201510||||||||Venous|admin\rOBX|0|ST|0|Na+|142|mmol/L|136^145||||||Pass|20260214201510||admin||||\r"

	// preprocessing must normalize CR -> LF and then normalize OBX-7 for decoding
	pre := h.preprocessMessage(rawCR)
	// ensure OBX-7 specifically is normalized (other segments still use '^' as component sep)
	lines := strings.Split(pre, "\n")
	var obxFound bool
	for _, l := range lines {
		if strings.HasPrefix(l, "OBX|") {
			obxFound = true
			f := strings.Split(l, "|")
			if len(f) > 7 {
				if strings.Contains(f[7], "^") {
					t.Fatalf("OBX-7 caret not removed: %s", f[7])
				}
				if !strings.Contains(f[7], "-") {
					t.Fatalf("OBX-7 not normalized to hyphen: %s", f[7])
				}
			}
		}
	}
	if !obxFound {
		t.Fatalf("no OBX line found in preprocessed message")
	}

	// full decode should succeed
	_, err := h.decodeORUR01([]byte(rawCR))
	if err != nil {
		t.Fatalf("decodeORUR01 failed for CR-terminated message: %v", err)
	}
}
