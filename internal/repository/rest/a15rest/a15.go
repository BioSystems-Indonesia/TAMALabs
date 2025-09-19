package a15rest

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"regexp"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type A15rest struct {
	client *http.Client
}

func NewA15() *A15rest {
	return &A15rest{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (a *A15rest) Send(ctx context.Context, req *entity.SendPayloadRequest) error {
	address := fmt.Sprintf("http://%s:%s/api/v1/request/a15", req.Device.IPAddress, req.Device.SendPort)

	buf := bytes.NewBuffer(nil)
	w := multipart.NewWriter(buf)
	f, err := w.CreateFormFile("file", "import.txt")
	if err != nil {
		return fmt.Errorf("a15: cannot create form file: %w", err)
	}

	content := createContentFile(req)
	_, err = f.Write(content)
	if err != nil {
		return fmt.Errorf("a15: error write to multipart: %w", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("a15: error to close multipart writer: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, address, buf)
	httpReq.Header.Set("Content-Type", w.FormDataContentType())

	res, err := a.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("a15: send failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, errReadBody := io.ReadAll(res.Body)
		return fmt.Errorf("a15: got %s (%d) from server errReadBody: %w, body: %s",
			res.Status,
			res.StatusCode,
			errReadBody,
			body,
		)
	}

	slog.Info("a15: success send to A15rest")
	return nil
}

func (a *A15rest) CheckConnection(ctx context.Context, device entity.Device) error {
	address := fmt.Sprintf("http://%s:%s/", device.IPAddress, device.SendPort)

	res, err := a.client.Get(address)
	if err != nil {
		return fmt.Errorf("a15: check connection error: %w", err)
	}

	if res.StatusCode != http.StatusNotFound {
		return nil
	}

	if res.StatusCode != http.StatusOK {
		return nil
	}

	return fmt.Errorf("a15: check connection got status %d", res.StatusCode)
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
				samples = append(samples, row(req, s, r))
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

func row(req *entity.SendPayloadRequest, s entity.Specimen, r entity.ObservationRequest) Sample {
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
