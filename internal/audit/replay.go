package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ReplayEntry represents a single point-in-time snapshot with metadata.
type ReplayEntry struct {
	Timestamp time.Time      `json:"timestamp"`
	Label     string         `json:"label"`
	Reports   []CompareReport `json:"reports"`
}

// ReplayStore holds a series of replay entries.
type ReplayStore struct {
	Entries []ReplayEntry `json:"entries"`
}

// SaveReplay persists a replay store to disk.
func SaveReplay(path string, store ReplayStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal replay: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadReplay reads a replay store from disk.
func LoadReplay(path string) (ReplayStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ReplayStore{}, fmt.Errorf("read replay: %w", err)
	}
	var store ReplayStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ReplayStore{}, fmt.Errorf("unmarshal replay: %w", err)
	}
	return store, nil
}

// AddReplayEntry appends a new entry and re-sorts by timestamp.
func AddReplayEntry(store *ReplayStore, label string, reports []CompareReport) {
	store.Entries = append(store.Entries, ReplayEntry{
		Timestamp: time.Now().UTC(),
		Label:     label,
		Reports:   reports,
	})
	sort.Slice(store.Entries, func(i, j int) bool {
		return store.Entries[i].Timestamp.Before(store.Entries[j].Timestamp)
	})
}

// ReplayAt returns the entry closest to (but not after) the given time.
func ReplayAt(store ReplayStore, at time.Time) (ReplayEntry, bool) {
	var best ReplayEntry
	found := false
	for _, e := range store.Entries {
		if !e.Timestamp.After(at) {
			best = e
			found = true
		}
	}
	return best, found
}
