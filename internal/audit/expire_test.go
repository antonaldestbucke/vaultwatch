package audit

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleExpiryStore() ExpiryStore {
	now := time.Now()
	return ExpiryStore{
		Rules: []ExpiryRule{
			{Pattern: "secret/prod/*", TTL: 24 * time.Hour, NotifyAt: 2 * time.Hour},
			{Pattern: "secret/staging/*", TTL: 48 * time.Hour, NotifyAt: 4 * time.Hour},
		},
		LastSeen: map[string]time.Time{
			"prod:secret/prod/db": now.Add(-25 * time.Hour),
			"staging:secret/staging/api": now.Add(-1 * time.Hour),
		},
	}
}

func TestSaveAndLoadExpiry_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "expiry.json")
	store := sampleExpiryStore()
	if err := SaveExpiry(path, store); err != nil {
		t.Fatalf("SaveExpiry: %v", err)
	}
	loaded, err := LoadExpiry(path)
	if err != nil {
		t.Fatalf("LoadExpiry: %v", err)
	}
	if len(loaded.Rules) != len(store.Rules) {
		t.Errorf("expected %d rules, got %d", len(store.Rules), len(loaded.Rules))
	}
	if len(loaded.LastSeen) != len(store.LastSeen) {
		t.Errorf("expected %d last_seen entries, got %d", len(store.LastSeen), len(loaded.LastSeen))
	}
}

func TestLoadExpiry_MissingFile(t *testing.T) {
	store, err := LoadExpiry("/nonexistent/expiry.json")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if store.LastSeen == nil {
		t.Error("expected initialized LastSeen map")
	}
}

func TestEvaluateExpiry_Expired(t *testing.T) {
	now := time.Now()
	store := sampleExpiryStore()
	reports := []ScoredReport{
		{Path: "secret/prod/db", Env: "prod", Score: 40, Risk: "high"},
	}
	results := EvaluateExpiry(reports, store, now)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "expired" {
		t.Errorf("expected status 'expired', got '%s'", results[0].Status)
	}
}

func TestEvaluateExpiry_Warning(t *testing.T) {
	now := time.Now()
	store := sampleExpiryStore()
	// staging/api was seen 1 hour ago, TTL=48h, notifyAt=4h -> ok (47h remaining)
	// override: set last seen to 45h ago so 3h remain, within notifyAt=4h
	store.LastSeen["staging:secret/staging/api"] = now.Add(-45 * time.Hour)
	reports := []ScoredReport{
		{Path: "secret/staging/api", Env: "staging", Score: 70, Risk: "medium"},
	}
	results := EvaluateExpiry(reports, store, now)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Status != "warning" {
		t.Errorf("expected status 'warning', got '%s'", results[0].Status)
	}
}

func TestEvaluateExpiry_NoMatch(t *testing.T) {
	now := time.Now()
	store := sampleExpiryStore()
	reports := []ScoredReport{
		{Path: "secret/dev/service", Env: "dev", Score: 90, Risk: "low"},
	}
	results := EvaluateExpiry(reports, store, now)
	if len(results) != 0 {
		t.Errorf("expected 0 results for unmatched path, got %d", len(results))
	}
}

func TestLoadExpiry_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "expiry.json")
	_ = os.WriteFile(path, []byte("not-json"), 0644)
	_, err := LoadExpiry(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
