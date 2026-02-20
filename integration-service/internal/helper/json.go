package helper

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0644)
}
