package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleReplayReports() []CompareReport {
	return []CompareReport{
		{Path: "secret/app", OnlyInA: []string{"key1"}, OnlyInB: []string{}},
	}
}

func TestSaveAndLoadReplay_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "replay.json")

	var store ReplayStore
	AddReplayEntry(&store, "v1", sampleReplayReports())

	if err := SaveReplay(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadReplay(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Label != "v1" {
		t.Errorf("expected label v1, got %s", loaded.Entries[0].Label)
	}
}

func TestLoadReplay_MissingFile(t *testing.T) {
	_, err := LoadReplay("/nonexistent/replay.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAddReplayEntry_SortsByTimestamp(t *testing.T) {
	var store ReplayStore
	now := time.Now().UTC()
	store.Entries = append(store.Entries, ReplayEntry{Timestamp: now.Add(time.Hour), Label: "later"})
	AddReplayEntry(&store, "earlier", nil)
	// earlier entry has timestamp ~now, later has now+1h; sort should put earlier first
	if store.Entries[len(store.Entries)-1].Label != "later" {
		t.Errorf("expected 'later' last, got %s", store.Entries[len(store.Entries)-1].Label)
	}
}

func TestReplayAt_ReturnsClosest(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	store := ReplayStore{
		Entries: []ReplayEntry{
			{Timestamp: base, Label: "first"},
			{Timestamp: base.Add(2 * time.Hour), Label: "second"},
		},
	}
	entry, ok := ReplayAt(store, base.Add(time.Hour))
	if !ok {
		t.Fatal("expected entry found")
	}
	if entry.Label != "first" {
		t.Errorf("expected 'first', got %s", entry.Label)
	}
}

func TestReplayAt_NoneFound(t *testing.T) {
	store := ReplayStore{
		Entries: []ReplayEntry{
			{Timestamp: time.Now().Add(time.Hour), Label: "future"},
		},
	}
	_, ok := ReplayAt(store, time.Now().Add(-time.Hour))
	if ok {
		t.Error("expected no entry found")
	}
}

func TestLoadReplay_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadReplay(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
