package analyzer

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"net"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// A15Result represents a single lab result entry.
type A15Result struct {
	PatientName string    `json:"patient_name"`
	TestName    string    `json:"test_name"`
	SampleType  string    `json:"sample_type"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
}

func (u *Usecase) ProcessA15(ctx context.Context) error {
	wol, err := u.WorkOrderRepository.FindByStatus(ctx, entity.WorkOrderStatusPending)
	if err != nil {
		return err
	}

	for i, wo := range wol {
		isA15 := slices.ContainsFunc(wo.Devices, func(d entity.Device) bool {
			return d.Type == string(entity.DeviceTypeA15)
		})

		if !isA15 {
			continue
		}








	}




}

func ConnectToSamba(device entity.Device) error {

	conn, err := net.Dial("tcp", net.JoinHostPort(device.IPAddress, strconv.Itoa(device.Port)))
	if err != nil {
		return err
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
		return err
	}
	defer s.Logoff()

	fs, err := s.Mount(device.Type)
	if err != nil {
		return err
	}

	defer fs.Umount()

	fil, err := fs.ReadDir(".")
	if err != nil {
		return err
	}

	for _, fi := range fil {
		f, err := fs.Open(fi.Name())
		if err != nil {
			slog.Error(err)
		}

		b := io.ReadAll(f)
		lr, err := ParseLabResults(string(b))
		if err != nil {
			return err
		}

	}


}

// String provides a string representation of LabResult, useful for printing.
func (lr A15Result) String() string {
	return fmt.Sprintf(
		"Patient: %s, Test: %s, Sample: %s, Value: %.3f %s, Time: %s",
		lr.PatientName,
		lr.TestName,
		lr.SampleType,
		lr.Value,
		lr.Unit,
		lr.Timestamp.Format("2006-01-02 15:04:05"), // Standard format for display
	)
}

// ParseLabResults parses a raw string containing multiple lab result lines
// using encoding/csv with a tab delimiter.
func ParseLabResults(data string) ([]A15Result, error) {
	var results []A15Result

	stringReader := strings.NewReader(strings.TrimSpace(data))
	csvReader := csv.NewReader(stringReader)

	csvReader.Comma = '\t'
	csvReader.FieldsPerRecord = 6
	csvReader.LazyQuotes = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV data: %w", err)
	}

	for i, record := range records {
		// The csv.Reader already splits the line into fields.
		// We already set FieldsPerRecord = 6, so the reader would have errored
		// if a line didn't have exactly 6 fields.

		patientID := strings.TrimSpace(record[0])
		testName := strings.TrimSpace(record[1])
		sampleType := strings.TrimSpace(record[2])

		valueStr := strings.ReplaceAll(strings.TrimSpace(record[3]), ",", ".")
		value, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: could not parse value '%s': %w", i+1, record[3], err)
		}

		unit := strings.TrimSpace(record[4])

		timestampStr := strings.TrimSpace(record[5])
		timestamp, err := time.Parse("02/01/2006 15:04:05", timestampStr)
		if err != nil {
			return nil, fmt.Errorf("line %d: could not parse timestamp '%s': %w", i+1, timestampStr, err)
		}

		results = append(results, A15Result{
			PatientName: patientID,
			TestName:    testName,
			SampleType:  sampleType,
			Value:       value,
			Unit:        unit,
			Timestamp:   timestamp,
		})
	}

	return results, nil
}
