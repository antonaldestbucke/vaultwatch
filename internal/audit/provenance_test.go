package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleProvenanceStore() ProvenanceStore {
	now := time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
	return ProvenanceStore{
		Entries: []ProvenanceEntry{
			{Path: "secret/db", Environment: "prod", ObservedAt: now, Source: "vault"},
			{Path: "secret/api", Environment: "staging", ObservedAt: now.Add(-time.Hour), Source: "snapshot"},
		},
	}
}

func TestSaveAndLoadProvenance_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "provenance.json")
	store := sampleProvenanceStore()

	if err := SaveProvenance(file, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadProvenance(file)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Entries) != len(store.Entries) {
		t.Errorf("expected %d entries, got %d", len(store.Entries), len(loaded.Entries))
	}
}

func TestLoadProvenance_MissingFile(t *testing.T) {
	store, err := LoadProvenance("/nonexistent/provenance.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(store.Entries) != 0 {
		t.Errorf("expected empty store, got %d entries", len(store.Entries))
	}
}

func TestLoadProvenance_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(file, []byte("not-json"), 0o644)
	_, err := LoadProvenance(file)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestAddProvenanceEntry_DeduplicatesByPathEnvSource(t *testing.T) {
	store := &ProvenanceStore{}
	t1 := time.Now().Add(-time.Hour)
	t2 := time.Now()

	AddProvenanceEntry(store, ProvenanceEntry{Path: "secret/db", Environment: "prod", ObservedAt: t1, Source: "vault"})
	AddProvenanceEntry(store, ProvenanceEntry{Path: "secret/db", Environment: "prod", ObservedAt: t2, Source: "vault"})

	if len(store.Entries) != 1 {
		t.Errorf("expected 1 entry after dedup, got %d", len(store.Entries))
	}
	if !store.Entries[0].ObservedAt.Equal(t2) {
		t.Errorf("expected updated timestamp")
	}
}

func TestAddProvenanceEntry_SortsByPath(t *testing.T) {
	store := &ProvenanceStore{}
	now := time.Now()
	AddProvenanceEntry(store, ProvenanceEntry{Path: "secret/z", Environment: "prod", ObservedAt: now, Source: "vault"})
	AddProvenanceEntry(store, ProvenanceEntry{Path: "secret/a", Environment: "prod", ObservedAt: now, Source: "vault"})

	if store.Entries[0].Path != "secret/a" {
		t.Errorf("expected sorted order, first path = %s", store.Entries[0].Path)
	}
}

func TestLookupProvenance_FiltersByPath(t *testing.T) {
	store := sampleProvenanceStore()
	results := LookupProvenance(store, "secret/db")
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Environment != "prod" {
		t.Errorf("unexpected environment: %s", results[0].Environment)
	}
}
