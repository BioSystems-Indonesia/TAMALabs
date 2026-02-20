package abbott

import (
	"testing"
)

func TestParseAbbottData(t *testing.T) {
	// Sample data from user
	rawData := `EMERALD;1;030420-008217;FSE
RESULT
DATE;14/01/2026
TIME;07:23:24
MODE;NORMAL
UNIT;1
SEQ;3;0
SID;3
PID;168217
ID;BY BEATRIKS KRISALSA
TYPE;ANAK
TEST;LMG
OPERATOR;ADMIN
WBC;8.5  ;;;4.0  ;4.0  ;12.0 ;12.0 
RBC;5.08 ;;;3.50 ;3.50 ;5.20 ;5.20 
HGB;18.0;;H;12.0;12.0;16.0;16.0
HCT;53.4;;H;35.0;35.0;49.0;49.0
MCV;105.2;;H;80.0 ;80.0 ;100.0;100.0
MCH;35.4 ;;H;27.0 ;27.0 ;34.0 ;34.0 
MCHC;33.7 ;;;31.0 ;31.0 ;37.0 ;37.0 
RDW;12.7 ;;;11.0 ;11.0 ;16.0 ;16.0 
PLT;188  ;;;150  ;150  ;400  ;400  
MPV;8.3  ;;;6.5  ;6.5  ;12.0 ;12.0 
LYM%;29.9;s;H;3.0 ;3.0 ;15.0;15.0
MID%;14.8;s;;3.0 ;3.0 ;15.0;15.0
GRA%;55.3;s;;50.0;50.0;70.0;70.0
LYM;2.5  ;s;;0.8  ;0.8  ;7.0  ;7.0  
MID;1.3  ;s;;0.1  ;0.1  ;1.5  ;1.5  
GRA;4.7  ;s;;2.0  ;2.0  ;8.0  ;8.0  
WBC CURVE;0;0;0;0;6;20;35;40;38;36;35;37;44;58;81;110;141;170;198;216;212;189;158;130;106;88;77;74;78;84;89;94;97;96;93;89;84;82;79;73;66;59;54;50;46;45;45;43;41;40;41;43;47;51;52;51;52;53;54;56;59;60;60;59;62;67;70;71;72;71;70;71;72;73;72;71;69;68;68;65;62;60;59;60;62;62;62;63;61;56;51;46;`

	message, err := ParseAbbottData(rawData)
	if err != nil {
		t.Fatalf("Failed to parse Abbott data: %v", err)
	}

	if message == nil {
		t.Fatal("Parsed message is nil")
	}

	// Test sample info
	if message.SampleInfo.PatientID != "168217" {
		t.Errorf("Expected PatientID '168217', got '%s'", message.SampleInfo.PatientID)
	}

	if message.SampleInfo.SampleID != "3" {
		t.Errorf("Expected SampleID '3', got '%s'", message.SampleInfo.SampleID)
	}

	if message.SampleInfo.PatientName != "BY BEATRIKS KRISALSA" {
		t.Errorf("Expected PatientName 'BY BEATRIKS KRISALSA', got '%s'", message.SampleInfo.PatientName)
	}

	if message.SampleInfo.Date != "14/01/2026" {
		t.Errorf("Expected Date '14/01/2026', got '%s'", message.SampleInfo.Date)
	}

	if message.SampleInfo.Time != "07:23:24" {
		t.Errorf("Expected Time '07:23:24', got '%s'", message.SampleInfo.Time)
	}

	// Test device info
	if message.DeviceInfo.Mode != "NORMAL" {
		t.Errorf("Expected Mode 'NORMAL', got '%s'", message.DeviceInfo.Mode)
	}

	if message.DeviceInfo.Operator != "ADMIN" {
		t.Errorf("Expected Operator 'ADMIN', got '%s'", message.DeviceInfo.Operator)
	}

	// Test that we have test results
	if len(message.TestResults) == 0 {
		t.Fatal("No test results parsed")
	}

	// Test specific results
	foundWBC := false
	foundHGB := false

	for _, result := range message.TestResults {
		switch result.TestCode {
		case "WBC":
			foundWBC = true
			if result.Value != "8.5" {
				t.Errorf("Expected WBC value '8.5', got '%s'", result.Value)
			}
			if result.RefMin != "4.0" {
				t.Errorf("Expected WBC RefMin '4.0', got '%s'", result.RefMin)
			}
			if result.RefMax != "12.0" {
				t.Errorf("Expected WBC RefMax '12.0', got '%s'", result.RefMax)
			}
		case "HGB":
			foundHGB = true
			if result.Value != "18.0" {
				t.Errorf("Expected HGB value '18.0', got '%s'", result.Value)
			}
			if result.Flag != "H" {
				t.Errorf("Expected HGB flag 'H', got '%s'", result.Flag)
			}
		}
	}

	if !foundWBC {
		t.Error("WBC result not found")
	}
	if !foundHGB {
		t.Error("HGB result not found")
	}

	// Check that CURVE data was not parsed as test result
	for _, result := range message.TestResults {
		if result.TestCode == "WBC CURVE" {
			t.Error("WBC CURVE should not be parsed as test result")
		}
	}

	t.Logf("Successfully parsed %d test results", len(message.TestResults))
	for _, result := range message.TestResults {
		t.Logf("  %s: %s (Flag: %s, Ref: %s - %s)",
			result.TestCode, result.Value, result.Flag, result.RefMin, result.RefMax)
	}
}

func TestConvertToAbbottResults(t *testing.T) {
	rawData := `DATE;14/01/2026
TIME;07:23:24
PID;168217
SID;3
WBC;8.5  ;;;4.0  ;4.0  ;12.0 ;12.0 
HGB;18.0;;H;12.0;12.0;16.0;16.0`

	message, err := ParseAbbottData(rawData)
	if err != nil {
		t.Fatalf("Failed to parse Abbott data: %v", err)
	}

	results := ConvertToAbbottResults(message)

	if len(results) == 0 {
		t.Fatal("No results converted")
	}

	for _, result := range results {
		if result.PatientID != "168217" {
			t.Errorf("Expected PatientID '168217', got '%s'", result.PatientID)
		}
		if result.SampleID != "3" {
			t.Errorf("Expected SampleID '3', got '%s'", result.SampleID)
		}

		t.Logf("Result: Test=%s, Value=%s, PatientID=%s, SampleID=%s",
			result.TestName, result.Value, result.PatientID, result.SampleID)
	}
}
