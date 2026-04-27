package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// LifecycleStage represents the current stage of a secret path's lifecycle.
type LifecycleStage string

const (
	StageActive     LifecycleStage = "active"
	StageDeprecated LifecycleStage = "deprecated"
	StageRetired    LifecycleStage = "retired"
)

// LifecycleEntry defines the lifecycle state for a secret path.
type LifecycleEntry struct {
	Path      string         `json:"path"`
	Stage     LifecycleStage `json:"stage"`
	UpdatedAt time.Time      `json:"updated_at"`
	Note      string         `json:"note,omitempty"`
}

// LifecycleStore holds all lifecycle entries.
type LifecycleStore struct {
	Entries []LifecycleEntry `json:"entries"`
}

// SaveLifecycle writes the lifecycle store to disk.
func SaveLifecycle(path string, store LifecycleStore) error {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal lifecycle: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadLifecycle reads the lifecycle store from disk.
func LoadLifecycle(path string) (LifecycleStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return LifecycleStore{}, nil
		}
		return LifecycleStore{}, fmt.Errorf("read lifecycle: %w", err)
	}
	var store LifecycleStore
	if err := json.Unmarshal(data, &store); err != nil {
		return LifecycleStore{}, fmt.Errorf("unmarshal lifecycle: %w", err)
	}
	return store, nil
}

// SetLifecycleStage upserts a lifecycle entry for the given path.
func SetLifecycleStage(store *LifecycleStore, path string, stage LifecycleStage, note string) {
	for i, e := range store.Entries {
		if e.Path == path {
			store.Entries[i].Stage = stage
			store.Entries[i].UpdatedAt = time.Now().UTC()
			store.Entries[i].Note = note
			return
		}
	}
	store.Entries = append(store.Entries, LifecycleEntry{
		Path:      path,
		Stage:     stage,
		UpdatedAt: time.Now().UTC(),
		Note:      note,
	})
	sort.Slice(store.Entries, func(i, j int) bool {
		return store.Entries[i].Path < store.Entries[j].Path
	})
}

// ApplyLifecycle annotates scored reports with lifecycle stage metadata.
func ApplyLifecycle(reports []ScoredReport, store LifecycleStore) []ScoredReport {
	stageMap := make(map[string]LifecycleStage, len(store.Entries))
	for _, e := range store.Entries {
		stageMap[e.Path] = e.Stage
	}
	result := make([]ScoredReport, len(reports))
	for i, r := range reports {
		if stage, ok := stageMap[r.Path]; ok {
			r.Note = fmt.Sprintf("lifecycle:%s", stage)
		}
		result[i] = r
	}
	return result
}
