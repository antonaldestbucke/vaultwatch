package audit

import (
	"strings"
	"testing"
)

func sampleDriftReports() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/app/db", DriftScore: 0.9, Risk: "high", Diffs: []string{"only_in_a: password"}},
		{Path: "secret/app/api", DriftScore: 0.5, Risk: "medium", Diffs: []string{"only_in_b: token"}},
		{Path: "secret/app/cache", DriftScore: 0.0, Risk: "low", Diffs: []string{}},
		{Path: "secret/infra/tls", DriftScore: 0.7, Risk: "high", Diffs: []string{"only_in_a: cert"}},
	}
}

func TestBuildDriftSummary_Counts(t *testing.T) {
	reports := sampleDriftReports()
	summary := BuildDriftSummary(reports, 5)

	if summary.TotalPaths != 4 {
		t.Errorf("expected TotalPaths=4, got %d", summary.TotalPaths)
	}
	if summary.DriftedPaths != 3 {
		t.Errorf("expected DriftedPaths=3, got %d", summary.DriftedPaths)
	}
	if summary.CleanPaths != 1 {
		t.Errorf("expected CleanPaths=1, got %d", summary.CleanPaths)
	}
}

func TestBuildDriftSummary_DriftRate(t *testing.T) {
	reports := sampleDriftReports()
	summary := BuildDriftSummary(reports, 5)

	expected := 75.0
	if summary.DriftRate != expected {
		t.Errorf("expected DriftRate=%.1f, got %.1f", expected, summary.DriftRate)
	}
}

func TestBuildDriftSummary_TopDriftedSortedByScore(t *testing.T) {
	reports := sampleDriftReports()
	summary := BuildDriftSummary(reports, 2)

	if len(summary.TopDrifted) != 2 {
		t.Fatalf("expected 2 top drifted, got %d", len(summary.TopDrifted))
	}
	if !strings.Contains(summary.TopDrifted[0], "secret/app/db") {
		t.Errorf("expected highest drifted path first, got %s", summary.TopDrifted[0])
	}
}

func TestBuildDriftSummary_Empty(t *testing.T) {
	summary := BuildDriftSummary([]ScoredReport{}, 5)

	if summary.TotalPaths != 0 {
		t.Errorf("expected TotalPaths=0, got %d", summary.TotalPaths)
	}
	if summary.DriftRate != 0.0 {
		t.Errorf("expected DriftRate=0, got %f", summary.DriftRate)
	}
	if summary.EnvBreakdown == nil {
		t.Error("expected non-nil EnvBreakdown for empty input")
	}
}

func TestBuildDriftSummary_NoDrift(t *testing.T) {
	reports := []ScoredReport{
		{Path: "secret/clean/a", DriftScore: 0.0, Diffs: []string{}},
		{Path: "secret/clean/b", DriftScore: 0.0, Diffs: []string{}},
	}
	summary := BuildDriftSummary(reports, 5)

	if summary.DriftedPaths != 0 {
		t.Errorf("expected 0 drifted paths, got %d", summary.DriftedPaths)
	}
	if len(summary.TopDrifted) != 0 {
		t.Errorf("expected empty TopDrifted, got %v", summary.TopDrifted)
	}
}
