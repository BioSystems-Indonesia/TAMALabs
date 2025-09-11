package analyzer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

func (u *Usecase) ProcessBTS(ctx context.Context) error {
	devicesFind, err := u.DeviceRepository.FindAll(ctx, &entity.GetManyRequestDevice{
		Type: []entity.DeviceType{entity.DeviceTypeBTS},
	})
	if err != nil {
		return fmt.Errorf("error getting devices: %w", err)
	}

	errors := make([]error, 0)
	results := make([]entity.ObservationResult, 0)
	for _, device := range devicesFind.Data {
		data, err := u.getBTSData(ctx, device)
		if err != nil {
			return err
		}

		specimenCache := make(map[string]entity.Specimen)
		for _, record := range data {
			specimen, err := u.getSpecimenWithCache(ctx, specimenCache, record.ID)
			if err != nil {
				errors = append(errors, fmt.Errorf("error getting specimen: %w, barcode: %s", err, record.ID))
				continue
			}

			date, _ := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", record.Date, record.Time))
			results = append(results, entity.ObservationResult{
				SpecimenID:  int64(specimen.ID),
				TestCode:    record.Analyte,
				Description: record.Analyte,
				Values: entity.JSONStringArray{
					fmt.Sprintf("%.2f", record.Result),
				},
				Type:           record.Analyte,
				Unit:           record.Units,
				ReferenceRange: record.LowerRange + " - " + record.UpperRange,
				Date:           date,
				AbnormalFlag: entity.JSONStringArray{
					record.Status,
				},
				Comments: record.Status,
			})
		}
	}

	for _, result := range results {
		if err := u.ObservationResultRepository.Create(ctx, &result); err != nil {
			errors = append(errors, fmt.Errorf("error creating observation result: %w", err))
		}
	}

	if len(errors) > 0 {
		slog.Error("error process BTS", "errors", errors)
	}

	return nil
}

func (u *Usecase) getSpecimenWithCache(
	ctx context.Context,
	specimenCache map[string]entity.Specimen,
	barcode string,
) (entity.Specimen, error) {
	if _, ok := specimenCache[barcode]; ok {
		return specimenCache[barcode], nil
	}

	specimen, err := u.SpecimenRepository.FindByBarcode(ctx, barcode)
	if err != nil {
		return entity.Specimen{}, errors.New("specimen not found")
	}
	specimenCache[barcode] = specimen

	return specimen, nil
}

func (u *Usecase) getBTSData(
	ctx context.Context,
	d entity.Device,
) ([]entity.BTSMeasurementData, error) {
	baseURL := "http://" + net.JoinHostPort(d.IPAddress, d.ReceivePort)
	commandsURL := baseURL + "/commands"

	randomFilename := strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
	filename := randomFilename + ".csv"
	payload := entity.BTSRequest{
		Command: entity.BTSCommand{
			Type: "ExportDataDownloading",
			Payload: entity.BTSPayload{
				OperationID: uuid.New().String(),
				Filters: entity.BTSFilters{
					ReadingType: []string{"Blank", "Calibration", "QualityControl", "Sample", "AbsorbanceSample"},
					SampleID:    nil,
					TechniqueID: nil,
					StartDate:   time.Now().AddDate(0, 0, -14),
					EndDate:     time.Now(),
					TimePeriod:  "RANGE",
				},
				ExportOptions: entity.BTSExportOptions{
					FileName:        filename,
					Format:          "CSV",
					Target:          "Download",
					DeleteAfterDump: false,
				},
			},
		},
	}

	payloadByte, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, commandsURL, bytes.NewBuffer(payloadByte))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+"/")
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
	)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s Response: %s", resp.Status, string(body))
	}

	var response entity.BTSResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	fileURL := fmt.Sprintf("%s/files/%s", baseURL, filename)
	outputDir := os.TempDir()
	outputPath := filepath.Join(outputDir, filename)

	data, err := downloadAndParseFile(fileURL, outputPath, client)
	if err != nil {
		return nil, fmt.Errorf("error downloading and parsing file: %w", err)
	}

	slog.Info("Received data from BTS", "data", data)

	return data, nil
}

func downloadAndParseFile(url, outputPath string, client *http.Client) ([]entity.BTSMeasurementData, error) {
	// Get the file data
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Save the file
	out, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %w", err)
	}
	defer out.Close()

	_, err = out.Write(body)
	if err != nil {
		return nil, fmt.Errorf("error saving file: %w", err)
	}

	// Parse the CSV data
	data, err := parseCSVData(string(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV data: %w", err)
	}

	return data, nil
}

func parseCSVData(csvContent string) ([]entity.BTSMeasurementData, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.Comma = '\t'         // Set delimiter to tab
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	if len(records) == 0 {
		return nil, errors.New("no data found in CSV")
	}

	// Skip header row
	var data []entity.BTSMeasurementData
	for i, record := range records[1:] {
		if len(record) < 8 {
			slog.Warn("BTS: Warning: Row has insufficient columns, skipping", "row", i+2)
			continue
		}

		// Parse result as float
		result, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			slog.Warn("BTS: Warning: Row has invalid result value, skipping", "row", i+2, "value", record[3])
			continue
		}

		// Parse absorbance as float
		absorbance, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			slog.Warn("BTS: Warning: Row d has invalid absorbance value, skipping", "row", i+2, "value", record[5])
			continue
		}

		measurement := entity.BTSMeasurementData{
			MeasureType: record[0],
			ID:          record[1],
			Analyte:     record[2],
			Result:      result,
			Units:       record[4],
			Absorbance:  absorbance,
			Date:        record[6],
			Time:        record[7],
		}

		// Handle optional fields (status and ranges)
		if len(record) > 8 {
			measurement.Status = record[8]
		}
		if len(record) > 9 {
			measurement.LowerRange = record[9]
		}
		if len(record) > 10 {
			measurement.UpperRange = record[10]
		}

		data = append(data, measurement)
	}

	return data, nil
}
