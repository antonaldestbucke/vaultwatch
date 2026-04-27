package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleLifecycleStore() LifecycleStore {
	return LifecycleStore{
		Entries: []LifecycleEntry{
			{Path: "secret/app/db", Stage: StageActive, UpdatedAt: time.Now().UTC()},
			{Path: "secret/app/old", Stage: StageDeprecated, UpdatedAt: time.Now().UTC(), Note: "to be removed"},
		},
	}
}

func TestSaveAndLoadLifecycle_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lifecycle.json")
	store := sampleLifecycleStore()
	if err := SaveLifecycle(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadLifecycle(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != len(store.Entries) {
		t.Errorf("expected %d entries, got %d", len(store.Entries), len(loaded.Entries))
	}
}

func TestLoadLifecycle_MissingFile(t *testing.T) {
	store, err := LoadLifecycle("/nonexistent/lifecycle.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestLoadLifecycle_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "lifecycle.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)
	_, err := LoadLifecycle(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestSetLifecycleStage_Upsert(t *testing.T) {
	store := sampleLifecycleStore()
	SetLifecycleStage(&store, "secret/app/db", StageRetired, "archived")
	for _, e := range store.Entries {
		if e.Path == "secret/app/db" {
			if e.Stage != StageRetired {
				t.Errorf("expected retired, got %s", e.Stage)
			}
			if e.Note != "archived" {
				t.Errorf("expected note 'archived', got %s", e.Note)
			}
			return
		}
	}
	t.Error("entry not found after upsert")
}

func TestSetLifecycleStage_NewEntry(t *testing.T) {
	store := sampleLifecycleStore()
	origLen := len(store.Entries)
	SetLifecycleStage(&store, "secret/new/path", StageActive, "")
	if len(store.Entries) != origLen+1 {
		t.Errorf("expected %d entries, got %d", origLen+1, len(store.Entries))
	}
}

func TestApplyLifecycle_AnnotatesMatchingPath(t *testing.T) {
	store := sampleLifecycleStore()
	reports := []ScoredReport{
		{Path: "secret/app/db", Score: 90},
		{Path: "secret/app/old", Score: 40},
		{Path: "secret/app/other", Score: 70},
	}
	result := ApplyLifecycle(reports, store)
	if result[0].Note != "lifecycle:active" {
		t.Errorf("expected lifecycle:active, got %q", result[0].Note)
	}
	if result[1].Note != "lifecycle:deprecated" {
		t.Errorf("expected lifecycle:deprecated, got %q", result[1].Note)
	}
	if result[2].Note != "" {
		t.Errorf("expected empty note for unmatched path, got %q", result[2].Note)
	}
}
