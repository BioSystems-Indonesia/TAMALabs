package response911

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

// minimal mock Analyzer that captures the last ORU_R01 passed
type mockAnalyzer struct {
	called chan entity.ORU_R01
}

func (m *mockAnalyzer) ProcessORUR01(ctx context.Context, data entity.ORU_R01) error {
	select {
	case m.called <- data:
	default:
	}
	return nil
}

// implement other methods to satisfy interface (not used in test)
func (m *mockAnalyzer) ProcessORMO01(ctx context.Context, data entity.ORM_O01) ([]entity.Specimen, error) {
	return nil, nil
}
func (m *mockAnalyzer) ProcessAbbott(ctx context.Context, data entity.AbbottMessage) error {
	return nil
}
func (m *mockAnalyzer) ProcessOULR22(ctx context.Context, data entity.OUL_R22) error      { return nil }
func (m *mockAnalyzer) ProcessQBPQ11(ctx context.Context, data entity.QBP_Q11) error      { return nil }
func (m *mockAnalyzer) ProcessCoax(ctx context.Context, data entity.CoaxTestResult) error { return nil }
func (m *mockAnalyzer) ProcessDiestro(ctx context.Context, data entity.DiestroResult) error {
	return nil
}
func (m *mockAnalyzer) ProcessVerifyU120Batch(ctx context.Context, results []entity.VerifyResult) error {
	return nil
}
func (m *mockAnalyzer) ProcessORUR01Batch(ctx context.Context, data []entity.ORU_R01) error {
	return nil
}

func TestAggregator_MergesAndFlushes(t *testing.T) {
	m := &mockAnalyzer{called: make(chan entity.ORU_R01, 2)}
	a := newAggregator(m.ProcessORUR01, 50*time.Millisecond)

	// create two partial ORU messages for same barcode
	oru1 := &entity.ORU_R01{
		Patient: []entity.Patient{
			{
				Specimen: []entity.Specimen{
					{
						Barcode: "BC1",
						ObservationResult: []entity.ObservationResult{
							{TestCode: "UREA", Values: entity.JSONStringArray{"92.8"}},
						},
					},
				},
			},
		},
	}

	oru2 := &entity.ORU_R01{
		Patient: []entity.Patient{
			{
				Specimen: []entity.Specimen{
					{
						Barcode: "BC1",
						ObservationResult: []entity.ObservationResult{
							{TestCode: "CREA", Values: entity.JSONStringArray{"1.12"}},
						},
					},
				},
			},
		},
	}

	a.Add(oru1)
	// add second part shortly after
	time.Sleep(10 * time.Millisecond)
	a.Add(oru2)

	// wait for flush
	select {
	case res := <-m.called:
		// verify aggregated
		if len(res.Patient) == 0 || len(res.Patient[0].Specimen) == 0 {
			t.Fatalf("aggregated ORU missing specimen")
		}
		spec := res.Patient[0].Specimen[0]
		if spec.Barcode != "BC1" {
			t.Fatalf("unexpected barcode: %s", spec.Barcode)
		}
		codes := make(map[string]struct{})
		for _, o := range spec.ObservationResult {
			codes[o.TestCode] = struct{}{}
		}
		if !reflect.DeepEqual(codes, map[string]struct{}{"UREA": {}, "CREA": {}}) {
			t.Fatalf("expected merged tests UREA+CREA, got %v", codes)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for aggregator flush")
	}
}
