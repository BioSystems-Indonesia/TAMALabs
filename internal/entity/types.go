package entity

import (
	"fmt"
	"log/slog"

	"database/sql/driver"
	"encoding/json"
)

type JSONStringArray []string

// Scan implements the sql.Scanner interface to read JSON from the database
func (j *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	var bytes []byte

	// Handle both []byte and string types
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
		// Log when we receive string type (for monitoring)
		slog.Debug("JSONStringArray received string type from database", "value", v)
	default:
		return fmt.Errorf("unsupported data type: %T", value)
	}

	// Unmarshal JSON into the slice
	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface to write JSON to the database
func (j JSONStringArray) Value() (driver.Value, error) {
	// Marshal the slice into JSON
	return json.Marshal(j)
}
