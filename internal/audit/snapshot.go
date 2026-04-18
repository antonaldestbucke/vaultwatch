package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Snapshot represents a point-in-time capture of secret keys at a given path.
type Snapshot struct {
	Path      string            `json:"path"`
	Env       string            `json:"env"`
	Keys      []string          `json:"keys"`
	CapturedAt time.Time        `json:"captured_at"`
	Meta      map[string]string `json:"meta,omitempty"`
}

// SaveSnapshot writes a Snapshot to a JSON file at the given filepath.
func SaveSnapshot(s Snapshot, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating snapshot file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("encoding snapshot: %w", err)
	}
	return nil
}

// LoadSnapshot reads a Snapshot from a JSON file at the given filepath.
func LoadSnapshot(filepath string) (Snapshot, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return Snapshot{}, fmt.Errorf("opening snapshot file: %w", err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return Snapshot{}, fmt.Errorf("decoding snapshot: %w", err)
	}
	return s, nil
}
