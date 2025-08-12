package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"crypto/tls"
)

type Response struct {
	Message string `json:"message"`
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
	req.Header.Set("Referer", baseURL + "/")
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

	// Check if the response message is "ok"
	if response.Message != "ok" {
		fmt.Printf("Unexpected response message: %s\n", response.Message)
		return
	}

	// Construct the file download URL
	fileURL := fmt.Sprintf("%s/files/%s", baseURL, fileName)
	outputPath := filepath.Join(outputDir, fileName)

	// Download the file
	fmt.Printf("Downloading file from %s...\n", fileURL)
	if err := downloadFile(fileURL, outputPath, client); err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
		return
	}

	fmt.Printf("File successfully downloaded to: %s\n", outputPath)
}

func downloadFile(url, outputPath string, client *http.Client) error {
	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	// Get the file data
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}
