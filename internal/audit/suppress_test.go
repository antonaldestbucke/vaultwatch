package audit

import (
	"os"
	"testing"
	"time"
)

func sampleStore() SuppressStore {
	return SuppressStore{
		Rules: []SuppressRule{
			{Path: "secret/prod", Reason: "planned maintenance", ExpiresAt: time.Now().Add(1 * time.Hour)},
			{Path: "secret/legacy", Reason: "deprecated", ExpiresAt: time.Now().Add(-1 * time.Hour)}, // expired
		},
	}
}

func TestSaveAndLoadSuppressions_RoundTrip(t *testing.T) {
	f, _ := os.CreateTemp("", "suppress*.json")
	f.Close()
	defer os.Remove(f.Name())

	store := sampleStore()
	if err := SaveSuppressions(f.Name(), store); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := LoadSuppressions(f.Name())
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(loaded.Rules) != len(store.Rules) {
		t.Errorf("expected %d rules, got %d", len(store.Rules), len(loaded.Rules))
	}
}

func TestLoadSuppressions_MissingFile(t *testing.T) {
	_, err := LoadSuppressions("/nonexistent/suppress.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestIsSuppressed_ActiveRule(t *testing.T) {
	store := sampleStore()
	if !IsSuppressed(store, "secret/prod/db") {
		t.Error("expected path to be suppressed")
	}
}

func TestIsSuppressed_ExpiredRule(t *testing.T) {
	store := sampleStore()
	if IsSuppressed(store, "secret/legacy/token") {
		t.Error("expected expired rule to not suppress")
	}
}

func TestApplySuppressions_FiltersMatching(t *testing.T) {
	store := sampleStore()
	reports := []ScoredReport{
		{Path: "secret/prod/api", Score: 40},
		{Path: "secret/dev/api", Score: 80},
	}
	result := ApplySuppressions(reports, store)
	if len(result) != 1 {
		t.Errorf("expected 1 report, got %d", len(result))
	}
	if result[0].Path != "secret/dev/api" {
		t.Errorf("unexpected path: %s", result[0].Path)
	}
}
