package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// TestType represents the test_type table structure
type TestType struct {
	ID   int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name"`
	Code string `json:"code" gorm:"column:code"`
}

// TableName specifies the table name for GORM
func (TestType) TableName() string {
	return "test_types"
}

func main() {
	// Database path from the application
	const dbFileName = "./bin/biosystem-lims.db"

	// Connect to database using modernc.org/sqlite driver
	db, err := gorm.Open(sqlite.Dialector{
		DriverName: "sqlite",
		DSN:        dbFileName + "?_parseTime=true",
	}, &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Query all test types with only name and code
	var testTypes []TestType
	result := db.Select("name, code").Find(&testTypes)
	if result.Error != nil {
		log.Fatalf("Failed to query test types: %v", result.Error)
	}

	// Create CSV file
	filename := "test_types_export.csv"
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Name", "Code"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	// Write data rows
	for _, testType := range testTypes {
		record := []string{
			testType.Name,
			testType.Code,
		}
		if err := writer.Write(record); err != nil {
			log.Fatalf("Failed to write CSV record: %v", err)
		}
	}

	fmt.Printf("Successfully exported %d test types to %s\n", len(testTypes), filename)
	fmt.Printf("CSV file created: %s\n", filename)
}
