package a15

import (
	"context"
	"encoding/csv"
	"strings"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type Sample struct {
	Class       string
	Type        string
	PatientID   string
	TechniqueID string
	TestTubType string
}

func SendToA15(ctx context.Context, req entity.SendPayloadRequest) {
	s := createContentFile(req)
	println(s)
}

func createContentFile(req entity.SendPayloadRequest) string {
	var samples []Sample
	for _, p := range req.Patients {
		for _, s := range p.Specimen {
			for _, r := range s.ObservationRequest {
				samples = append(samples, row(req, p, s, r))

			}
		}
	}

	w := &strings.Builder{}
	cw := csv.NewWriter(w)
	cw.Comma = '\t'

	for _, s := range samples {
		cw.Write([]string{
			s.Class,
			s.Type,
			s.PatientID,
			s.TechniqueID,
			s.TestTubType,
		})
	}

	cw.Flush()
	return w.String()
}

func row(req entity.SendPayloadRequest, p entity.Patient, s entity.Specimen, r entity.ObservationRequest) Sample {
	class := "N"
	if req.Urgent {
		class = "U"
	}

	return Sample{
		Class:       class,
		Type:        s.Type,
		PatientID:   s.Barcode,
		TechniqueID: r.TestCode,
		TestTubType: "T15",
	}
}
