package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func sampleRetentionReports(now time.Time) []ScoredReport {
	return []ScoredReport{
		{Path: "secret/app/db", EnvA: "dev", EnvB: "prod", Score: 40, ScannedAt: now.Add(-72 * time.Hour)},
		{Path: "secret/app/api", EnvA: "dev", EnvB: "prod", Score: 90, ScannedAt: now.Add(-10 * time.Hour)},
		{Path: "secret/infra/tls", EnvA: "dev", EnvB: "prod", Score: 55, ScannedAt: now.Add(-200 * time.Hour)},
	}
}

func TestSaveAndLoadRetention_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "retention.json")
	store := RetentionStore{
		Rules: []RetentionRule{
			{PathPrefix: "secret/app", MaxAge: 48 * time.Hour},
		},
	}
	if err := SaveRetention(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadRetention(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].PathPrefix != "secret/app" {
		t.Errorf("unexpected rules: %+v", loaded.Rules)
	}
}

func TestLoadRetention_MissingFile(t *testing.T) {
	_, err := LoadRetention("/nonexistent/retention.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadRetention_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte(`{not valid json`), 0644)
	_, err := LoadRetention(path)
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestEvaluateRetention_PrunesOldReports(t *testing.T) {
	now := time.Now()
	reports := sampleRetentionReports(now)
	store := RetentionStore{
		Rules: []RetentionRule{
			{PathPrefix: "secret/app", MaxAge: 48 * time.Hour},
		},
	}
	results := EvaluateRetention(reports, store, now)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if !results[0].Pruned {
		t.Errorf("expected secret/app/db to be pruned (72h old, max 48h)")
	}
	if results[1].Pruned {
		t.Errorf("expected secret/app/api to be retained (10h old)")
	}
	if results[2].Pruned {
		t.Errorf("expected secret/infra/tls to be retained (no matching rule)")
	}
}

func TestEvaluateRetention_NoRules(t *testing.T) {
	now := time.Now()
	reports := sampleRetentionReports(now)
	results := EvaluateRetention(reports, RetentionStore{}, now)
	for _, r := range results {
		if r.Pruned {
			t.Errorf("expected no pruning with empty rules, got pruned: %s", r.Path)
		}
	}
}

func TestEvaluateRetention_ReasonContainsMaxAge(t *testing.T) {
	now := time.Now()
	reports := []ScoredReport{
		{Path: "secret/app/db", EnvA: "dev", EnvB: "prod", Score: 40, ScannedAt: now.Add(-96 * time.Hour)},
	}
	store := RetentionStore{
		Rules: []RetentionRule{
			{PathPrefix: "secret/app", MaxAge: 24 * time.Hour},
		},
	}
	results := EvaluateRetention(reports, store, now)
	if len(results) == 0 || !results[0].Pruned {
		t.Fatal("expected pruned result")
	}
	if results[0].Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestRetentionStore_JSONRoundTrip(t *testing.T) {
	store := RetentionStore{
		Rules: []RetentionRule{
			{PathPrefix: "secret/", MaxAge: 168 * time.Hour},
		},
	}
	data, err := json.Marshal(store)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out RetentionStore
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Rules[0].MaxAge != 168*time.Hour {
		t.Errorf("MaxAge not preserved: %v", out.Rules[0].MaxAge)
	}
}
