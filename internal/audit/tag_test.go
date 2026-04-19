package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleTagStore() TagStore {
	return TagStore{
		"secret/app/db":  {"critical", "pii"},
		"secret/app/api": {"external"},
	}
}

func sampleTagReports() []CompareReport {
	return []CompareReport{
		{Path: "secret/app/db"},
		{Path: "secret/app/api"},
		{Path: "secret/app/cache"},
	}
}

func TestSaveAndLoadTags_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "tags.json")
	store := sampleTagStore()
	if err := SaveTags(p, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadTags(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded) != len(store) {
		t.Errorf("expected %d entries, got %d", len(store), len(loaded))
	}
}

func TestLoadTags_MissingFile(t *testing.T) {
	_, err := LoadTags("/nonexistent/tags.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadTags_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	os.WriteFile(p, []byte("not-json"), 0644)
	_, err := LoadTags(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApplyTags_MatchingPaths(t *testing.T) {
	reports := sampleTagReports()
	store := sampleTagStore()
	result := ApplyTags(reports, store)
	if len(result["secret/app/db"]) != 2 {
		t.Errorf("expected 2 tags for db, got %d", len(result["secret/app/db"]))
	}
	if _, ok := result["secret/app/cache"]; ok {
		t.Error("cache should have no tags")
	}
}

func TestFilterByTag_ReturnsMatches(t *testing.T) {
	reports := sampleTagReports()
	store := sampleTagStore()
	out := FilterByTag(reports, store, "pii")
	if len(out) != 1 || out[0].Path != "secret/app/db" {
		t.Errorf("unexpected filter result: %+v", out)
	}
}

func TestFilterByTag_NoMatch(t *testing.T) {
	reports := sampleTagReports()
	store := sampleTagStore()
	out := FilterByTag(reports, store, "nonexistent")
	if len(out) != 0 {
		t.Errorf("expected 0 results, got %d", len(out))
	}
}
