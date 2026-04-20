package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// ChangelogEntry records a single detected change event for a secret path.
type ChangelogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path"`
	Env       string    `json:"env"`
	OnlyInA   []string  `json:"only_in_a,omitempty"`
	OnlyInB   []string  `json:"only_in_b,omitempty"`
}

// ChangelogStore holds all changelog entries.
type ChangelogStore struct {
	Entries []ChangelogEntry `json:"entries"`
}

// SaveChangelog writes the changelog store to a JSON file.
func SaveChangelog(path string, store ChangelogStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal changelog: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadChangelog reads a changelog store from a JSON file.
func LoadChangelog(path string) (ChangelogStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ChangelogStore{}, nil
		}
		return ChangelogStore{}, fmt.Errorf("read changelog: %w", err)
	}
	var store ChangelogStore
	if err := json.Unmarshal(data, &store); err != nil {
		return ChangelogStore{}, fmt.Errorf("unmarshal changelog: %w", err)
	}
	return store, nil
}

// AppendChangelog adds new entries derived from scored reports and saves the store.
func AppendChangelog(path string, reports []ScoredReport, env string) (ChangelogStore, error) {
	store, err := LoadChangelog(path)
	if err != nil {
		return ChangelogStore{}, err
	}
	now := time.Now().UTC()
	for _, r := range reports {
		if len(r.Report.OnlyInA) == 0 && len(r.Report.OnlyInB) == 0 {
			continue
		}
		store.Entries = append(store.Entries, ChangelogEntry{
			Timestamp: now,
			Path:      r.Report.Path,
			Env:       env,
			OnlyInA:   r.Report.OnlyInA,
			OnlyInB:   r.Report.OnlyInB,
		})
	}
	sort.Slice(store.Entries, func(i, j int) bool {
		return store.Entries[i].Timestamp.Before(store.Entries[j].Timestamp)
	})
	if err := SaveChangelog(path, store); err != nil {
		return ChangelogStore{}, err
	}
	return store, nil
}

// PrintChangelog prints changelog entries to stdout.
func PrintChangelog(store ChangelogStore) {
	if len(store.Entries) == 0 {
		fmt.Println("No changelog entries found.")
		return
	}
	for _, e := range store.Entries {
		fmt.Printf("[%s] %s (%s)\n", e.Timestamp.Format(time.RFC3339), e.Path, e.Env)
		for _, k := range e.OnlyInA {
			fmt.Printf("  - removed: %s\n", k)
		}
		for _, k := range e.OnlyInB {
			fmt.Printf("  + added:   %s\n", k)
		}
	}
}
