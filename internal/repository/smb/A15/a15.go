package a15

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net"
	"regexp"

	"github.com/hirochachacha/go-smb2"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type A15 struct{}

func NewA15() *A15 {
	return &A15{}
}

func (a A15) Send(ctx context.Context, req *entity.SendPayloadRequest) error {
	// // For testing - write to local directory instead of SMB
	// localDir := "tmp/a15_test"
	// err := os.MkdirAll(localDir, 0755)
	// if err != nil {
	// 	return fmt.Errorf("cannot create local directory %s: %v", localDir, err)
	// }

	// filePath := filepath.Join(localDir, "import.txt")
	// content := createContentFile(req)

	// err = os.WriteFile(filePath, content, 0644)
	// if err != nil {
	// 	return fmt.Errorf("cannot write file to local directory: %v", err)
	// }

	// slog.Info("Send to A15 (local test)", "file_path", filePath)
	// return nil

	conn, err := net.Dial("tcp", net.JoinHostPort(req.Device.IPAddress, req.Device.SendPort))
	if err != nil {
		return fmt.Errorf("cannot connect to %s:%s", req.Device.IPAddress, req.Device.SendPort)
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
		return fmt.Errorf("cannot dial smb2 to %s:%s", req.Device.IPAddress, req.Device.SendPort)
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
	conn, err := net.Dial("tcp", net.JoinHostPort(device.IPAddress, device.SendPort))
	if err != nil {
		return fmt.Errorf("cannot connect to %s:%s", device.IPAddress, device.SendPort)
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
		return fmt.Errorf("cannot dial smb2 to %s:%s", device.IPAddress, device.SendPort)
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
				if r.TestType.IsCalculatedTest {
					continue
				}
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

	// Remove any alphabetic prefix from barcode to get only the ID number
	// This will handle SER, URI, or any other prefix followed by numbers
	patientID := extractIDFromBarcode(s.Barcode)

	return Sample{
		Class:       class,
		Type:        s.Type,
		PatientID:   patientID,
		TechniqueID: r.TestCode,
		TestTubType: "T15",
	}
}

// extractIDFromBarcode removes any alphabetic prefix and returns only the numeric part
func extractIDFromBarcode(barcode string) string {
	// Use regex to find the first sequence of digits in the barcode
	re := regexp.MustCompile(`\d+`)
	match := re.FindString(barcode)
	if match != "" {
		return match
	}
	// If no digits found, return the original barcode
	return barcode
}
