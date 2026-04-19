package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleStore() AnnotationStore {
	return AnnotationStore{
		"secret/app/prod": {
			Path:      "secret/app/prod",
			Note:      "reviewed by security team",
			Author:    "alice",
			CreatedAt: time.Now(),
		},
	}
}

func TestSaveAndLoadAnnotations_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "annotations.json")
	store := sampleStore()

	if err := SaveAnnotations(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAnnotations(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if loaded["secret/app/prod"].Note != store["secret/app/prod"].Note {
		t.Errorf("note mismatch")
	}
}

func TestLoadAnnotations_MissingFile(t *testing.T) {
	_, err := LoadAnnotations("/nonexistent/annotations.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadAnnotations_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	os.WriteFile(path, []byte("not json"), 0644)
	_, err := LoadAnnotations(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApplyAnnotations_MatchingPath(t *testing.T) {
	store := sampleStore()
	reports := []CompareResult{
		{Path: "secret/app/prod", OnlyInA: []string{"key1"}},
		{Path: "secret/app/staging"},
	}
	result := ApplyAnnotations(reports, store)
	if result[0].Note != "reviewed by security team" {
		t.Errorf("expected note on matching path, got %q", result[0].Note)
	}
	if result[1].Note != "" {
		t.Errorf("expected empty note on non-matching path")
	}
}
