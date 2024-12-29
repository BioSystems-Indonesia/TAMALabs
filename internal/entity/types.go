package entity

import (
	"fmt"

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

	// Convert value to a byte slice
	bytes, ok := value.([]byte)
	if !ok {
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
