package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleScoredForChangelog() []ScoredReport {
	return []ScoredReport{
		{
			Report: CompareResult{
				Path:    "secret/app/db",
				OnlyInA: []string{"old_password"},
				OnlyInB: []string{"new_password"},
			},
		},
		{
			Report: CompareResult{
				Path:    "secret/app/api",
				OnlyInA: []string{},
				OnlyInB: []string{},
			},
		},
	}
}

func TestSaveAndLoadChangelog_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "changelog.json")

	store := ChangelogStore{
		Entries: []ChangelogEntry{
			{Timestamp: time.Now().UTC(), Path: "secret/app", Env: "prod", OnlyInA: []string{"key1"}},
		},
	}
	if err := SaveChangelog(p, store); err != nil {
		t.Fatalf("SaveChangelog: %v", err)
	}
	loaded, err := LoadChangelog(p)
	if err != nil {
		t.Fatalf("LoadChangelog: %v", err)
	}
	if len(loaded.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(loaded.Entries))
	}
	if loaded.Entries[0].Path != "secret/app" {
		t.Errorf("unexpected path: %s", loaded.Entries[0].Path)
	}
}

func TestLoadChangelog_MissingFile(t *testing.T) {
	store, err := LoadChangelog("/nonexistent/changelog.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store, got %d entries", len(store.Entries))
	}
}

func TestAppendChangelog_OnlyDrifted(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "changelog.json")

	reports := sampleScoredForChangelog()
	store, err := AppendChangelog(p, reports, "staging")
	if err != nil {
		t.Fatalf("AppendChangelog: %v", err)
	}
	// Only the drifted path should be recorded
	if len(store.Entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(store.Entries))
	}
	if store.Entries[0].Path != "secret/app/db" {
		t.Errorf("unexpected path: %s", store.Entries[0].Path)
	}
	if store.Entries[0].Env != "staging" {
		t.Errorf("unexpected env: %s", store.Entries[0].Env)
	}
}

func TestAppendChangelog_Accumulates(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "changelog.json")

	reports := sampleScoredForChangelog()
	_, _ = AppendChangelog(p, reports, "prod")
	store, err := AppendChangelog(p, reports, "prod")
	if err != nil {
		t.Fatalf("second AppendChangelog: %v", err)
	}
	if len(store.Entries) != 2 {
		t.Errorf("expected 2 accumulated entries, got %d", len(store.Entries))
	}
}

func TestLoadChangelog_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(p, []byte("not-json"), 0644)
	_, err := LoadChangelog(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
