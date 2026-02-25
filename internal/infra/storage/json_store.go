package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Load reads a JSON file at path and unmarshals it into v.
func Load(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Save marshals v as indented JSON and writes it to path.
func Save(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// DataDir returns (and creates) the FocusPlay data directory under the OS cache dir.
func DataDir() string {
	base, _ := os.UserCacheDir() // %LOCALAPPDATA% on Windows
	dir := filepath.Join(base, "FocusPlay")
	_ = os.MkdirAll(dir, 0755)
	return dir
}
