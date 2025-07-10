package analyzer

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// A15Result represents a single lab result entry.
type A15Result struct {
	PatientID  string    `json:"patient_id"`
	TestName   string    `json:"test_name"`
	SampleType string    `json:"sample_type"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Timestamp  time.Time `json:"timestamp"`
}

func (u *Usecase) ProcessA15(ctx context.Context) error {
	wol, err := u.WorkOrderRepository.FindByStatus(ctx, entity.WorkOrderStatusPending)
	if err != nil {
		return err
	}

	a15Device := make([]entity.Device, 0)

	for _, wo := range wol {
		for _, device := range wo.Devices {
			if device.Type == entity.DeviceTypeA15 {
				a15Device = append(a15Device, device)
			}
		}
	}
	a15Device = removeDuplicates(a15Device)

	lrs := make([]A15Result, 0)
	for _, device := range a15Device {
		lr, err := connectToSamba(device)
		if err != nil {
			return err
		}
		lrs = append(lrs, lr...)
	}

	for _, lr := range lrs {
		speciment, err := u.SpecimenRepository.FindByBarcode(ctx, lr.PatientID)
		if err != nil {
			slog.Error("specimen not found", "barcode", lr.PatientID, "error", err)
			continue
		}

		observation := entity.ObservationResult{
			SpecimenID:  int64(speciment.ID),
			TestCode:    lr.TestName,
			Description: lr.TestName,
			Values:      []string{fmt.Sprintf("%.2f", lr.Value)},
			Unit:        lr.Unit,
			Date:        lr.Timestamp,
		}

		err = u.ObservationResultRepository.Create(ctx, &observation)
		if err != nil {
			slog.Error("failed to create observation result", "specimen_id", speciment.ID, "test_code", lr.TestName, "error", err)
			continue
		}
	}

	return nil
}

func connectToSamba(device entity.Device) ([]A15Result, error) {
	conn, err := net.Dial("tcp", net.JoinHostPort(device.IPAddress, strconv.Itoa(device.SendPort)))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer s.Logoff()

	//fs, err := s.Mount(device.Path)
	fs, err := s.Mount("Export")
	if err != nil {
		return nil, err
	}

	defer fs.Umount()

	fil, err := fs.ReadDir(".")
	if err != nil {
		return nil, err
	}
	lrs := make([]A15Result, 0)
	for _, fi := range fil {
		f, err := fs.Open(fi.Name())
		if err != nil {
			return nil, err
		}

		b, _ := io.ReadAll(f)
		lr, err := ParseLabResults(string(b))
		if err != nil {
			return nil, err
		}
		lrs = append(lrs, lr...)
	}
	return lrs, nil
}

// String provides a string representation of LabResult, useful for printing.
func (lr A15Result) String() string {
	return fmt.Sprintf(
		"Patient: %s, Test: %s, Sample: %s, Value: %.3f %s, Time: %s",
		lr.PatientID,
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
			PatientID:  patientID,
			TestName:   testName,
			SampleType: sampleType,
			Value:      value,
			Unit:       unit,
			Timestamp:  timestamp,
		})
	}

	return results, nil
}

func removeDuplicates(slice []entity.Device) []entity.Device {
	seen := make(map[int]struct{})
	var result []entity.Device

	for _, item := range slice {
		if _, exists := seen[item.ID]; !exists { // Assuming `ID` is a unique identifier for `entity.Device`
			seen[item.ID] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// save file result to db
func (u *Usecase) SaveFileResult(context context.Context, data string) error {
	results, err := ParseLabResults(data)
	if err != nil {
		return err
	}

	var errs []error
	for _, result := range results {
		speciment, err := u.SpecimenRepository.FindByBarcode(context, result.PatientID)
		if err != nil {
			slog.Error("specimen not found", "barcode", result.PatientID, "error", err)
			continue
		}

		err = u.ObservationResultRepository.Create(context, &entity.ObservationResult{
			SpecimenID:  int64(speciment.ID),
			TestCode:    result.TestName,
			Description: result.TestName,
			Values:      []string{fmt.Sprintf("%.2f", result.Value)},
			Unit:        result.Unit,
			Date:        result.Timestamp,
		})
		if err != nil {
			errs = append(errs, err)
		}

		slog.Info(
			"observation result created",
			"specimen_id",
			speciment.ID,
			"test_code",
			result.TestName,
			"value",
			result.Value,
			"unit",
			result.Unit,
			"date",
			result.Timestamp,
		)
	}

	if len(errs) > 0 {
		slog.Error("failed to create observation result", "errors", errors.Join(errs...))
	}

	return nil
}
