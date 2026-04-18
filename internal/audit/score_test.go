package audit

import (
	"testing"
)

func sampleScoreReports() []CompareReport {
	return []CompareReport{
		{Path: "secret/a", Diffs: []DiffResult{{Path: "secret/a", OnlyInA: []string{"key1"}}}}，
		{Path: "secret/b", Diffs: []DiffResult{}},
		{Path: "secret/c", Diffs: []DiffResult{{Path: "secret/c", OnlyInB: []string{"key2"}}}}，
		{Path: "secret/d", Diffs: []DiffResult{}},
	}
}

func TestScoreReports_Mixed(t *testing.T) {
	reports := sampleScoreReports()
	result := ScoreReports(reports)

	if result.TotalPaths != 4 {
		t.Errorf("expected 4 total paths, got %d", result.TotalPaths)
	}
	if result.DriftedPaths != 2 {
		t.Errorf("expected 2 drifted paths, got %d", result.DriftedPaths)
	}
	if result.Score != 50.0 {
		t.Errorf("expected score 50.0, got %.1f", result.Score)
	}
	if result.Risk != RiskMedium {
		t.Errorf("expected risk medium, got %s", result.Risk)
	}
}

func TestScoreReports_NoDrift(t *testing.T) {
	reports := []CompareReport{
		{Path: "secret/a", Diffs: []DiffResult{}},
	}
	result := ScoreReports(reports)
	if result.Risk != RiskNone {
		t.Errorf("expected none risk, got %s", result.Risk)
	}
	if result.Score != 0 {
		t.Errorf("expected score 0, got %.1f", result.Score)
	}
}

func TestScoreReports_Empty(t *testing.T) {
	result := ScoreReports([]CompareReport{})
	if result.Risk != RiskNone {
		t.Errorf("expected none risk for empty input")
	}
	if result.TotalPaths != 0 {
		t.Errorf("expected 0 total paths")
	}
}

func TestScoreReports_AllDrifted(t *testing.T) {
	reports := []CompareReport{
		{Path: "secret/x", Diffs: []DiffResult{{Path: "secret/x", OnlyInA: []string{"k"}}}},
		{Path: "secret/y", Diffs: []DiffResult{{Path: "secret/y", OnlyInB: []string{"k"}}}},
	}
	result := ScoreReports(reports)
	if result.Risk != RiskHigh {
		t.Errorf("expected high risk, got %s", result.Risk)
	}
	if result.Score != 100.0 {
		t.Errorf("expected score 100.0, got %.1f", result.Score)
	}
}
