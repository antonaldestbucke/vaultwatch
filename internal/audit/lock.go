package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// LockEntry represents a locked secret path that should not drift.
type LockEntry struct {
	Path      string    `json:"path"`
	LockedAt  time.Time `json:"locked_at"`
	LockedBy  string    `json:"locked_by"`
	Reason    string    `json:"reason"`
}

// LockStore holds all lock entries.
type LockStore struct {
	Locks []LockEntry `json:"locks"`
}

// SaveLocks writes the lock store to disk.
func SaveLocks(path string, store LockStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal locks: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadLocks reads the lock store from disk.
func LoadLocks(path string) (LockStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return LockStore{}, nil
		}
		return LockStore{}, fmt.Errorf("read locks: %w", err)
	}
	var store LockStore
	if err := json.Unmarshal(data, &store); err != nil {
		return LockStore{}, fmt.Errorf("unmarshal locks: %w", err)
	}
	return store, nil
}

// ApplyLocks annotates reports whose paths are locked.
func ApplyLocks(reports []CompareReport, store LockStore) []CompareReport {
	locked := make(map[string]LockEntry)
	for _, l := range store.Locks {
		locked[l.Path] = l
	}
	result := make([]CompareReport, len(reports))
	copy(result, reports)
	for i, r := range result {
		if entry, ok := locked[r.Path]; ok && len(r.OnlyInA)+len(r.OnlyInB) > 0 {
			r.Notes = fmt.Sprintf("LOCKED by %s: %s", entry.LockedBy, entry.Reason)
			result[i] = r
		}
	}
	return result
}
