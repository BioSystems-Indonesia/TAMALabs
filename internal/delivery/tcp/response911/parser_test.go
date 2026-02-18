package response911

import (
	"testing"
)

func TestParseResponse911_Simple(t *testing.T) {
	raw := "4O|1|SARTINEM|SARTINEM|||20260217092021|||||||||S||||||||||F\r" +
		"5R|1|^^^^UA|4.21|mg/dL||||V||@@DEFUSER@@\r" +
		"4L|1|N\r"

	oru, err := ParseResponse911(raw)
	if err != nil {
		t.Fatalf("ParseResponse911 returned error: %v", err)
	}

	if len(oru.Patient) == 0 {
		t.Fatalf("expected patient present")
	}

	s := oru.Patient[0].Specimen
	if len(s) == 0 {
		t.Fatalf("expected specimen present")
	}

	if s[0].Barcode != "SARTINEM" {
		t.Fatalf("expected barcode SARTINEM, got %s", s[0].Barcode)
	}

	if len(s[0].ObservationResult) == 0 {
		t.Fatalf("expected observation result")
	}

	tr := s[0].ObservationResult[0]
	if tr.TestCode != "UA" {
		t.Fatalf("expected test code UA, got %s", tr.TestCode)
	}

	if len(tr.Values) == 0 || tr.Values[0] != "4.21" {
		t.Fatalf("expected value 4.21, got %v", tr.Values)
	}

	if tr.Unit != "mg/dL" {
		t.Fatalf("expected unit mg/dL, got %s", tr.Unit)
	}
}
