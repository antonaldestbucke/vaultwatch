package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func samplePolicyReports() []CompareReport {
	return []CompareReport{
		{Path: "secret/prod/db", OnlyInA: []string{"password", "host"}, OnlyInB: []string{}},
		{Path: "secret/prod/api", OnlyInA: []string{"token"}, OnlyInB: []string{"debug_key"}},
		{Path: "secret/dev/db", OnlyInA: []string{"host"}, OnlyInB: []string{}},
	}
}

func TestSaveAndLoadPolicy_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	store := PolicyStore{
		Rules: []PolicyRule{
			{PathPrefix: "secret/prod", RequiredKeys: []string{"password"}, ForbiddenKeys: []string{"debug_key"}},
		},
	}
	if err := SavePolicy(path, store); err != nil {
		t.Fatalf("save: %v", err)
	}
	loaded, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(loaded.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(loaded.Rules))
	}
}

func TestLoadPolicy_MissingFile(t *testing.T) {
	_, err := LoadPolicy("/nonexistent/policy.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestEvaluatePolicy_RequiredKeyMissing(t *testing.T) {
	store := PolicyStore{
		Rules: []PolicyRule{
			{PathPrefix: "secret/prod", RequiredKeys: []string{"password"}},
		},
	}
	reports := []CompareReport{
		{Path: "secret/prod/api", OnlyInA: []string{"token"}, OnlyInB: []string{}},
	}
	violations := EvaluatePolicy(reports, store)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
}

func TestEvaluatePolicy_ForbiddenKeyPresent(t *testing.T) {
	store := PolicyStore{
		Rules: []PolicyRule{
			{PathPrefix: "secret/prod", ForbiddenKeys: []string{"debug_key"}},
		},
	}
	reports := samplePolicyReports()
	violations := EvaluatePolicy(reports, store)
	if len(violations) != 1 {
		t.Errorf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Path != "secret/prod/api" {
		t.Errorf("unexpected path: %s", violations[0].Path)
	}
}

func TestEvaluatePolicy_NoViolations(t *testing.T) {
	store := PolicyStore{
		Rules: []PolicyRule{
			{PathPrefix: "secret/staging", RequiredKeys: []string{"token"}},
		},
	}
	violations := EvaluatePolicy(samplePolicyReports(), store)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestEvaluatePolicy_PrefixNotMatched(t *testing.T) {
	store := PolicyStore{
		Rules: []PolicyRule{
			{PathPrefix: "secret/prod", RequiredKeys: []string{"password"}},
		},
	}
	reports := []CompareReport{
		{Path: "secret/dev/db", OnlyInA: []string{"host"}, OnlyInB: []string{}},
	}
	violations := EvaluatePolicy(reports, store)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d", len(violations))
	}
}

func TestSavePolicy_InvalidPath(t *testing.T) {
	err := SavePolicy("/nonexistent/dir/policy.json", PolicyStore{})
	if err == nil {
		t.Error("expected error for invalid path")
	}
	_ = os.Remove("/nonexistent/dir/policy.json")
}
