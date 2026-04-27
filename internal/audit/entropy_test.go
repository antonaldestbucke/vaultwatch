package audit

import (
	"testing"
)

func sampleEntropyReports() []CompareReport {
	return []CompareReport{
		{
			Path: "secret/app/prod",
			EnvA: "staging",
			EnvB: "prod",
			Diff: DiffResult{
				OnlyInA: []string{"db_host", "db_port", "db_name", "db_user", "db_pass"},
				OnlyInB: []string{"cache_host", "cache_port", "redis_url", "redis_pass"},
			},
		},
		{
			Path: "secret/app/staging",
			EnvA: "staging",
			EnvB: "prod",
			Diff: DiffResult{
				OnlyInA: []string{"token"},
				OnlyInB: []string{},
			},
		},
		{
			Path: "secret/app/clean",
			EnvA: "staging",
			EnvB: "prod",
			Diff: DiffResult{
				OnlyInA: []string{},
				OnlyInB: []string{},
			},
		},
	}
}

func TestBuildEntropy_ExcludesNoDiff(t *testing.T) {
	results := BuildEntropy(sampleEntropyReports())
	for _, r := range results {
		if r.Path == "secret/app/clean" {
			t.Errorf("expected clean path to be excluded, got %+v", r)
		}
	}
}

func TestBuildEntropy_SortedByEntropy(t *testing.T) {
	results := BuildEntropy(sampleEntropyReports())
	if len(results) < 2 {
		t.Fatal("expected at least 2 results")
	}
	if results[0].Entropy < results[1].Entropy {
		t.Errorf("expected descending entropy order, got %.3f < %.3f", results[0].Entropy, results[1].Entropy)
	}
}

func TestBuildEntropy_RiskLabel(t *testing.T) {
	results := BuildEntropy(sampleEntropyReports())
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	validRisks := map[string]bool{"low": true, "medium": true, "high": true, "critical": true}
	for _, r := range results {
		if !validRisks[r.Risk] {
			t.Errorf("unexpected risk label %q for path %s", r.Risk, r.Path)
		}
	}
}

func TestBuildEntropy_Empty(t *testing.T) {
	results := BuildEntropy([]CompareReport{})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestShannonEntropy_SingleKey(t *testing.T) {
	e := shannonEntropy([]string{"only_key"})
	if e != 0 {
		t.Errorf("expected entropy 0 for single key, got %.3f", e)
	}
}

func TestClassifyEntropy_Boundaries(t *testing.T) {
	cases := []struct {
		input    float64
		expected string
	}{
		{0.5, "low"},
		{1.5, "medium"},
		{2.5, "high"},
		{3.5, "critical"},
	}
	for _, tc := range cases {
		got := classifyEntropy(tc.input)
		if got != tc.expected {
			t.Errorf("classifyEntropy(%.1f) = %q, want %q", tc.input, got, tc.expected)
		}
	}
}
