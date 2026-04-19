package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleLineageStore() LineageStore {
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	return LineageStore{
		Entries: []LineageEntry{
			{Timestamp: t1, Path: "secret/app", Keys: []string{"db_pass"}, Env: "prod"},
			{Timestamp: t2, Path: "secret/app", Keys: []string{"db_pass", "api_key"}, Env: "prod"},
		},
	}
}

func TestSaveAndLoadLineage_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lineage.json")
	store := sampleLineageStore()
	if err := SaveLineage(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadLineage(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != len(store.Entries) {
		t.Errorf("expected %d entries, got %d", len(store.Entries), len(loaded.Entries))
	}
}

func TestLoadLineage_MissingFile(t *testing.T) {
	_, err := LoadLineage("/nonexistent/lineage.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAddLineageEntry_SortsByTimestamp(t *testing.T) {
	store := LineageStore{}
	t2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	AddLineageEntry(&store, LineageEntry{Timestamp: t2, Path: "secret/x", Env: "dev"})
	AddLineageEntry(&store, LineageEntry{Timestamp: t1, Path: "secret/x", Env: "dev"})
	if !store.Entries[0].Timestamp.Before(store.Entries[1].Timestamp) {
		t.Error("entries not sorted by timestamp")
	}
}

func TestHistoryForPath_FiltersByPathAndEnv(t *testing.T) {
	store := sampleLineageStore()
	AddLineageEntry(&store, LineageEntry{
		Timestamp: time.Now(), Path: "secret/other", Keys: []string{"x"}, Env: "prod",
	})
	result := HistoryForPath(store, "secret/app", "prod")
	if len(result) != 2 {
		t.Errorf("expected 2 entries, got %d", len(result))
	}
}

func TestLoadLineage_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadLineage(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
