package a15

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	"github.com/hirochachacha/go-smb2"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type A15 struct{}

func NewA15() *A15 {
	return &A15{}
}

func (a A15) Send(ctx context.Context, req *entity.SendPayloadRequest) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(req.Device.IPAddress, strconv.Itoa(req.Device.SendPort)))
	if err != nil {
		return fmt.Errorf("cannot connect to %s:%d", req.Device.IPAddress, req.Device.SendPort)
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
		return fmt.Errorf("cannot dial smb2 to %s:%d", req.Device.IPAddress, req.Device.SendPort)
	}
	defer func() {
		if err := s.Logoff(); err != nil {
			slog.Error("error logging off SMB session", "error", err)
		}
	}()

	fs, err := s.Mount(req.Device.Path)
	if err != nil {
		return fmt.Errorf("cannot mount SMB session: %v", err)
	}
	defer func() {
		if err := fs.Umount(); err != nil {
			slog.Error("error unmounting SMB session", "error", err)
		}
	}()

	err = fs.WriteFile("import.txt", createContentFile(req), 0644)
	if err != nil {
		return fmt.Errorf("cannot write file to SMB session: %v", err)
	}

	slog.Info("Send to A15")
	return nil
}

func (a A15) CheckConnection(ctx context.Context, device entity.Device) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(device.IPAddress, strconv.Itoa(device.SendPort)))
	if err != nil {
		return fmt.Errorf("cannot connect to %s:%d", device.IPAddress, device.SendPort)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     device.Username,
			Password: device.Password,
		},
	}

	s, err := d.Dial(conn)
	if err != nil {
		return fmt.Errorf("cannot dial smb2 to %s:%d", device.IPAddress, device.SendPort)
	}
	defer func() {
		if err := s.Logoff(); err != nil {
			slog.Info("error logging off SMB session", "error", err)
		}
	}()

	return nil
}

type Sample struct {
	Class       string
	Type        string
	PatientID   string
	TechniqueID string
	TestTubType string
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
