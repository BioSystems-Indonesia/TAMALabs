package analyzer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
		data, err := u.getBTSData(device)
		if err != nil {
			return err
		}

		specimenCache := make(map[string]entity.Specimen)
		for _, record := range data {
			specimen, err := u.getSpecimenWithCache(ctx, specimenCache, record.ID)
			if err != nil {
				errors = append(errors, fmt.Errorf("error getting specimen: %v, barcode: %s", err, record.ID))
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

	if len(errors) > 0 {
		slog.Error("error getting specimen", "errors", errors)
	}

	return nil
}

func (u *Usecase) getSpecimenWithCache(ctx context.Context, specimenCache map[string]entity.Specimen, barcode string) (entity.Specimen, error) {
	if _, ok := specimenCache[barcode]; ok {
		return specimenCache[barcode], nil
	}

	specimen, err := u.SpecimenRepository.FindByBarcode(ctx, barcode)
	if err != nil {
		return entity.Specimen{}, fmt.Errorf("specimen not found")
	}
	specimenCache[barcode] = specimen

	return specimen, nil
}

func (u *Usecase) getBTSData(d entity.Device) ([]entity.BTSMeasurementData, error) {
	baseURL := fmt.Sprintf("http://%s:%s", d.IPAddress, d.ReceivePort)
	commandsURL := baseURL + "/commands"

	payload := `{"command":{"type":"ExportDataDownloading","payload":{"operationId":"14bc1e38-68d7-4087-9e68-5f0818c46087","filters":{"readingType":["Blank","Calibration","QualityControl","Sample","AbsorbanceSample"],"sampleId":null,"techniqueId":null,"startDate":"2025-07-01T00:00:00.000+07:00","endDate":"2025-08-06T23:59:59.999+07:00","timePeriod":"RANGE"},"exportOptions":{"fileName":"file0001.csv","format":"CSV","target":"Download","deleteAfterDump":false}}}}`

	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		return nil, fmt.Errorf("error parsing payload: %v", err)
	}

	command := payloadMap["command"].(map[string]interface{})
	exportOptions := command["payload"].(map[string]interface{})["exportOptions"].(map[string]interface{})
	fileName := exportOptions["fileName"].(string)

	req, err := http.NewRequest("POST", commandsURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %s\nResponse: %s\n", resp.Status, string(body))
	}

	var response entity.BTSResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parsing response: %v", err)
	}

	fileURL := fmt.Sprintf("%s/files/%s", baseURL, fileName)
	outputDir := os.TempDir()
	outputPath := filepath.Join(outputDir, fileName)

	data, err := downloadAndParseFile(fileURL, outputPath, client)
	if err != nil {
		return nil, fmt.Errorf("error downloading and parsing file: %v", err)
	}

	slog.Info("Received data from BTS", "data", data)

	return data, nil
}

func downloadAndParseFile(url, outputPath string, client *http.Client) ([]entity.BTSMeasurementData, error) {
	// Get the file data
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Save the file
	out, err := os.Create(outputPath)
	if err != nil {
		return nil, fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	_, err = out.Write(body)
	if err != nil {
		return nil, fmt.Errorf("error saving file: %v", err)
	}

	// Parse the CSV data
	data, err := parseCSVData(string(body))
	if err != nil {
		return nil, fmt.Errorf("error parsing CSV data: %v", err)
	}

	return data, nil
}

func parseCSVData(csvContent string) ([]entity.BTSMeasurementData, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.Comma = '\t'         // Set delimiter to tab
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no data found in CSV")
	}

	// Skip header row
	var data []entity.BTSMeasurementData
	for i, record := range records[1:] {
		if len(record) < 8 {
			fmt.Printf("Warning: Row %d has insufficient columns, skipping\n", i+2)
			continue
		}

		// Parse result as float
		result, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			fmt.Printf("Warning: Row %d has invalid result value '%s', skipping\n", i+2, record[3])
			continue
		}

		// Parse absorbance as float
		absorbance, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			fmt.Printf("Warning: Row %d has invalid absorbance value '%s', skipping\n", i+2, record[5])
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
