package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ProvenanceEntry records where a secret path was observed and when.
type ProvenanceEntry struct {
	Path        string    `json:"path"`
	Environment string    `json:"environment"`
	ObservedAt  time.Time `json:"observed_at"`
	Source      string    `json:"source"` // e.g. "vault", "snapshot", "baseline"
}

// ProvenanceStore holds all provenance records keyed by path.
type ProvenanceStore struct {
	Entries []ProvenanceEntry `json:"entries"`
}

// SaveProvenance writes the store to disk as JSON.
func SaveProvenance(path string, store ProvenanceStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("provenance: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadProvenance reads a ProvenanceStore from disk.
func LoadProvenance(path string) (ProvenanceStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ProvenanceStore{}, nil
		}
		return ProvenanceStore{}, fmt.Errorf("provenance: read: %w", err)
	}
	var store ProvenanceStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ProvenanceStore{}, fmt.Errorf("provenance: unmarshal: %w", err)
	}
	return store, nil
}

// AddProvenanceEntry appends a new entry, deduplicating by path+env+source.
func AddProvenanceEntry(store *ProvenanceStore, entry ProvenanceEntry) {
	for i, e := range store.Entries {
		if e.Path == entry.Path && e.Environment == entry.Environment && e.Source == entry.Source {
			store.Entries[i].ObservedAt = entry.ObservedAt
			return
		}
	}
	store.Entries = append(store.Entries, entry)
	sort.Slice(store.Entries, func(i, j int) bool {
		if store.Entries[i].Path != store.Entries[j].Path {
			return store.Entries[i].Path < store.Entries[j].Path
		}
		return store.Entries[i].ObservedAt.Before(store.Entries[j].ObservedAt)
	})
}

// LookupProvenance returns all entries for a given path.
func LookupProvenance(store ProvenanceStore, path string) []ProvenanceEntry {
	var results []ProvenanceEntry
	for _, e := range store.Entries {
		if e.Path == path {
			results = append(results, e)
		}
	}
	return results
}
