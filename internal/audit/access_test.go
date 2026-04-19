package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleAccessStore() AccessStore {
	return AccessStore{
		Rules: []AccessRule{
			{Path: "secret/prod", Owner: "alice", Team: "platform"},
			{Path: "secret/staging", Owner: "bob", Team: "dev", Expires: time.Now().Add(24 * time.Hour)},
		},
	}
}

func TestSaveAndLoadAccess_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "access.json")
	store := sampleAccessStore()
	if err := SaveAccess(p, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadAccess(p)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != len(store.Rules) {
		t.Errorf("expected %d rules, got %d", len(store.Rules), len(loaded.Rules))
	}
}

func TestLoadAccess_MissingFile(t *testing.T) {
	store, err := LoadAccess("/nonexistent/access.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.Rules) != 0 {
		t.Errorf("expected empty store")
	}
}

func TestLoadAccess_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "access.json")
	os.WriteFile(p, []byte("not json"), 0644)
	_, err := LoadAccess(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApplyAccess_MatchingPath(t *testing.T) {
	store := sampleAccessStore()
	reports := []CompareReport{
		{Path: "secret/prod/db"},
		{Path: "secret/other/key"},
	}
	result := ApplyAccess(reports, store)
	if result[0].Annotations["owner"] != "alice" {
		t.Errorf("expected owner alice, got %s", result[0].Annotations["owner"])
	}
	if result[0].Annotations["team"] != "platform" {
		t.Errorf("expected team platform, got %s", result[0].Annotations["team"])
	}
	if result[1].Annotations["owner"] != "" {
		t.Errorf("expected no owner for unmatched path")
	}
}

func TestLookupAccess_Found(t *testing.T) {
	store := sampleAccessStore()
	rule, ok := LookupAccess(store, "secret/prod/api")
	if !ok {
		t.Fatal("expected rule to be found")
	}
	if rule.Owner != "alice" {
		t.Errorf("expected alice, got %s", rule.Owner)
	}
}

func TestLookupAccess_NotFound(t *testing.T) {
	store := sampleAccessStore()
	_, ok := LookupAccess(store, "secret/unknown/path")
	if ok {
		t.Error("expected no match")
	}
}
