package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleOwnershipStore() OwnershipStore {
	return OwnershipStore{
		Owners: []OwnerEntry{
			{Path: "secret/prod", Owner: "alice", Team: "platform", Contact: "alice@example.com"},
			{Path: "secret/staging", Owner: "bob", Team: "dev", Contact: "bob@example.com"},
		},
	}
}

func TestSaveAndLoadOwnership_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ownership.json")
	store := sampleOwnershipStore()
	if err := SaveOwnership(p, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadOwnership(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Owners) != 2 {
		t.Errorf("expected 2 owners, got %d", len(loaded.Owners))
	}
	// Verify individual fields survive the round-trip.
	if loaded.Owners[0].Owner != "alice" {
		t.Errorf("expected alice, got %s", loaded.Owners[0].Owner)
	}
	if loaded.Owners[1].Contact != "bob@example.com" {
		t.Errorf("expected bob@example.com, got %s", loaded.Owners[1].Contact)
	}
}

func TestLoadOwnership_MissingFile(t *testing.T) {
	_, err := LoadOwnership("/nonexistent/ownership.json")
	if err != nil {
		t.Errorf("expected nil for missing file, got %v", err)
	}
}

func TestLoadOwnership_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	os.WriteFile(p, []byte("not-json"), 0644)
	_, err := LoadOwnership(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApplyOwnership_MatchingPath(t *testing.T) {
	store := sampleOwnershipStore()
	reports := []CompareReport{
		{Path: "secret/prod/db", Diffs: []DiffResult{}},
	}
	result := ApplyOwnership(reports, store)
	if result[0].Annotations["owner"] != "alice" {
		t.Errorf("expected owner alice, got %s", result[0].Annotations["owner"])
	}
	if result[0].Annotations["team"] != "platform" {
		t.Errorf("expected team platform, got %s", result[0].Annotations["team"])
	}
}

func TestApplyOwnership_NoMatch(t *testing.T) {
	store := sampleOwnershipStore()
	reports := []CompareReport{
		{Path: "secret/other/path", Diffs: []DiffResult{}},
	}
	result := ApplyOwnership(reports, store)
	if result[0].Annotations["owner"] != "" {
		t.Errorf("expected no owner annotation, got %s", result[0].Annotations["owner"])
	}
}

func TestLookupOwner_Found(t *testing.T) {
	store := sampleOwnershipStore()
	entry, ok := LookupOwner(store, "secret/staging/api")
	if !ok {
		t.Fatal("expected to find owner")
	}
	if entry.Owner != "bob" {
		t.Errorf("expected bob, got %s", entry.Owner)
	}
}

func TestLookupOwner_NotFound(t *testing.T) {
	store := sampleOwnershipStore()
	_, ok := LookupOwner(store, "secret/other/path")
	if ok {
		t.Error("expected not found")
	}
}
