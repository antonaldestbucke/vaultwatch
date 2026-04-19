package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// LineageEntry records a single observed state of a secret path.
type LineageEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Keys      []string  `json:"keys"`
	Env       string    `json:"env"`
}

// LineageStore holds all lineage entries indexed by path.
type LineageStore struct {
	Entries []LineageEntry `json:"entries"`
}

// SaveLineage writes the lineage store to a file.
func SaveLineage(path string, store LineageStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lineage: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadLineage reads a lineage store from a file.
func LoadLineage(path string) (LineageStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return LineageStore{}, fmt.Errorf("read lineage: %w", err)
	}
	var store LineageStore
	if err := json.Unmarshal(data, &store); err != nil {
		return LineageStore{}, fmt.Errorf("unmarshal lineage: %w", err)
	}
	return store, nil
}

// AddLineageEntry appends an entry and sorts by timestamp.
func AddLineageEntry(store *LineageStore, entry LineageEntry) {
	store.Entries = append(store.Entries, entry)
	sort.Slice(store.Entries, func(i, j int) bool {
		return store.Entries[i].Timestamp.Before(store.Entries[j].Timestamp)
	})
}

// HistoryForPath returns all entries for a given secret path and env.
func HistoryForPath(store LineageStore, secretPath, env string) []LineageEntry {
	var result []LineageEntry
	for _, e := range store.Entries {
		if e.Path == secretPath && e.Env == env {
			result = append(result, e)
		}
	}
	return result
}
