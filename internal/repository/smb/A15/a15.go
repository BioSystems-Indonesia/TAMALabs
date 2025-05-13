package a15

import (
	"bytes"
	"context"
	"encoding/csv"
	"log"
	"net"
	"strconv"

	"github.com/hirochachacha/go-smb2"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type Sample struct {
	Class       string
	Type        string
	PatientID   string
	TechniqueID string
	TestTubType string
}

func SendToA15(ctx context.Context, req *entity.SendPayloadRequest) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(req.Device.IPAddress, strconv.Itoa(req.Device.Port)))
	if err != nil {
		return err
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     req.Device.Username,
			Password: req.Device.Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return err
	}
	defer s.Logoff()

	fs, err := s.Mount(req.Device.Path)
	if err != nil {
		return err
	}

	defer fs.Umount()

	err = fs.WriteFile("import.txt", createContentFile(req), 0644)
	if err != nil {
		return err
	}

	log.Println("Send to A15")
	return nil
}

func createContentFile(req *entity.SendPayloadRequest) []byte {
	var samples []Sample
	for _, p := range req.Patients {
		for _, s := range p.Specimen {
			for _, r := range s.ObservationRequest {
				samples = append(samples, row(req, p, s, r))

			}
		}
	}

	w := &bytes.Buffer{}
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
	return w.Bytes()
}

func row(req *entity.SendPayloadRequest, p entity.Patient, s entity.Specimen, r entity.ObservationRequest) Sample {
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
