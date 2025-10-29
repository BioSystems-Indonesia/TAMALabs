package rest

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/logger"
	"github.com/labstack/echo/v4"
)

type LogHandler struct {
	cfg *config.Schema
}

func NewLogHandler(cfg *config.Schema) *LogHandler {
	return &LogHandler{
		cfg: cfg,
	}
}

// getLastNLines reads the last n lines from a file
func getLastNLines(file *os.File, n int) ([]string, error) {
	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := stat.Size()
	if fileSize == 0 {
		return []string{}, nil
	}

	// Read the entire file into memory
	file.Seek(0, os.SEEK_SET)
	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Return the last n lines
	if len(lines) <= n {
		return lines, nil
	}

	return lines[len(lines)-n:], nil
}

func (h *LogHandler) StreamLog(c echo.Context) error {
	writer, flusher, err := createSSEWriter(c)
	if err != nil {
		return handleErrorSSE(c, writer, fmt.Errorf("failed to create SSE writer: %w", err))
	}

	file, err := os.Open(logger.GetDefaultLogFile())
	if err != nil {
		return err
	}
	defer file.Close()

	// Read and send the last 200 lines first
	lines, err := getLastNLines(file, 200)
	if err != nil {
		return fmt.Errorf("failed to read last lines: %w", err)
	}

	// Send the last 200 lines
	for _, line := range lines {
		if line != "" {
			_, err := writer.Write([]byte(entity.NewLogStreamingResponse(line)))
			if err != nil {
				return err
			}
			flusher.Flush()
		}
	}

	// Seek to the end of the file for streaming new lines
	file.Seek(0, os.SEEK_END)
	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			time.Sleep(500 * time.Millisecond) // Wait for new data
			continue
		}
		if line != "" {
			_, err := writer.Write([]byte(entity.NewLogStreamingResponse(line)))
			if err != nil {
				return err
			}
			flusher.Flush()
		}
	}

	return c.NoContent(http.StatusOK)
}

func (h *LogHandler) ExportLog(c echo.Context) error {
	// Prefix and timestamp for zip file
	prefix := "lis_exported_data__"
	timestamp := time.Now().Format("20060102_150405")
	zipFileName := fmt.Sprintf("%s%s.zip", prefix, timestamp)
	zipFilePath := filepath.Join(os.TempDir(), zipFileName)

	srcDir := filepath.Join(".", "logs")
	if err := ZipDir(srcDir, zipFilePath); err != nil {
		return handleError(c, fmt.Errorf("failed to zip tmp dir: %w", err))
	}

	return c.File(zipFilePath)
}
