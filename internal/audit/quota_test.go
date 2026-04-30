package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func sampleQuotaReports() []ScoredReport {
	return []ScoredReport{
		{
			Report: CompareReport{
				Path:    "secret/app/prod",
				OnlyInA: []string{"KEY_A", "KEY_B", "KEY_C"},
				OnlyInB: []string{},
			},
			Score: 40,
		},
		{
			Report: CompareReport{
				Path:    "secret/app/staging",
				OnlyInA: []string{"KEY_X"},
				OnlyInB: []string{},
			},
			Score: 80,
		},
		{
			Report: CompareReport{
				Path:    "secret/infra/db",
				OnlyInA: []string{},
				OnlyInB: []string{},
			},
			Score: 100,
		},
	}
}

func TestSaveAndLoadQuota_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "quota.json")
	store := QuotaStore{
		Rules: []QuotaRule{
			{PathPrefix: "secret/app", MaxDrifted: 2},
		},
	}
	if err := SaveQuota(path, store); err != nil {
		t.Fatalf("SaveQuota: %v", err)
	}
	loaded, err := LoadQuota(path)
	if err != nil {
		t.Fatalf("LoadQuota: %v", err)
	}
	if len(loaded.Rules) != 1 || loaded.Rules[0].MaxDrifted != 2 {
		t.Errorf("unexpected loaded rules: %+v", loaded.Rules)
	}
}

func TestLoadQuota_MissingFile(t *testing.T) {
	_, err := LoadQuota("/nonexistent/quota.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestEvaluateQuota_DetectsViolation(t *testing.T) {
	reports := sampleQuotaReports()
	store := QuotaStore{
		Rules: []QuotaRule{
			{PathPrefix: "secret/app", MaxDrifted: 2},
		},
	}
	violations := EvaluateQuota(reports, store)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "secret/app/prod" {
		t.Errorf("expected secret/app/prod, got %s", violations[0].Path)
	}
	if violations[0].Drifted != 3 {
		t.Errorf("expected drifted=3, got %d", violations[0].Drifted)
	}
}

func TestEvaluateQuota_NoDrift_NoViolation(t *testing.T) {
	reports := sampleQuotaReports()
	store := QuotaStore{
		Rules: []QuotaRule{
			{PathPrefix: "secret/infra", MaxDrifted: 0},
		},
	}
	violations := EvaluateQuota(reports, store)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d", len(violations))
	}
}

func TestEvaluateQuota_EmptyStore(t *testing.T) {
	reports := sampleQuotaReports()
	violations := EvaluateQuota(reports, QuotaStore{})
	if len(violations) != 0 {
		t.Errorf("expected no violations with empty store, got %d", len(violations))
	}
}

func TestLoadQuota_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := LoadQuota(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
