package main

import (
	"bytes"
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Response struct {
	Message string `json:"message"`
}

type MeasurementData struct {
	MeasureType  string  `json:"measureType"`
	ID           string  `json:"id"`
	Analyte      string  `json:"analyte"`
	Result       float64 `json:"result"`
	Units        string  `json:"units"`
	Absorbance   float64 `json:"absorbance"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
	Status       string  `json:"status,omitempty"`
	LowerRange   string  `json:"lowerRange,omitempty"`
	UpperRange   string  `json:"upperRange,omitempty"`
}

func main() {
	baseURL := "http://192.168.1.52:8080"
	commandsURL := baseURL + "/commands"

	// JSON payload
	payload := `{"command":{"type":"ExportDataDownloading","payload":{"operationId":"14bc1e38-68d7-4087-9e68-5f0818c46087","filters":{"readingType":["Blank","Calibration","QualityControl","Sample","AbsorbanceSample"],"sampleId":null,"techniqueId":null,"startDate":"2025-07-01T00:00:00.000+07:00","endDate":"2025-08-06T23:59:59.999+07:00","timePeriod":"RANGE"},"exportOptions":{"fileName":"file0001.csv","format":"CSV","target":"Download","deleteAfterDump":false}}}}`

	// Parse the payload to get the filename
	var payloadMap map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
		fmt.Printf("Error parsing payload: %v\n", err)
		return
	}

	command := payloadMap["command"].(map[string]interface{})
	exportOptions := command["payload"].(map[string]interface{})["exportOptions"].(map[string]interface{})
	fileName := exportOptions["fileName"].(string)

	// Create output directory if it doesn't exist
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Create a new request
	req, err := http.NewRequest("POST", commandsURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", baseURL)
	req.Header.Set("Referer", baseURL+"/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")

	// Create a custom HTTP client that skips TLS verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send the initial request
	fmt.Println("Sending export request...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status: %s\nResponse: %s\n", resp.Status, string(body))
		return
	}

	// Parse the response
	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}

	// Construct the file download URL
	fileURL := fmt.Sprintf("%s/files/%s", baseURL, fileName)
	outputPath := filepath.Join(outputDir, fileName)

	// Download and parse the file
	fmt.Printf("Downloading file from %s...\n", fileURL)
	data, err := downloadAndParseFile(fileURL, outputPath, client)
	if err != nil {
		fmt.Printf("Error downloading and parsing file: %v\n", err)
		return
	}

	fmt.Printf("File successfully downloaded to: %s\n", outputPath)
	fmt.Printf("Parsed %d records\n", len(data))

	// Print first few records as example
	for i, record := range data {
		if i >= 3 { // Show only first 3 records
			break
		}
		fmt.Printf("Record %d: %+v\n", i+1, record)
	}
}

func downloadAndParseFile(url, outputPath string, client *http.Client) ([]MeasurementData, error) {
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

func parseCSVData(csvContent string) ([]MeasurementData, error) {
	reader := csv.NewReader(strings.NewReader(csvContent))
	reader.Comma = '\t' // Set delimiter to tab
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV: %v", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("no data found in CSV")
	}

	// Skip header row
	var data []MeasurementData
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

		measurement := MeasurementData{
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
