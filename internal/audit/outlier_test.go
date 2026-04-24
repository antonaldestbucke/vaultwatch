package audit

import (
	"testing"
)

func sampleOutlierReports() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/app/db", Envs: []string{"dev", "prod"}, Score: 10.0, Risk: "low"},
		{Path: "secret/app/api", Envs: []string{"dev", "prod"}, Score: 12.0, Risk: "low"},
		{Path: "secret/app/cache", Envs: []string{"dev", "prod"}, Score: 11.0, Risk: "low"},
		{Path: "secret/app/admin", Envs: []string{"dev", "prod"}, Score: 95.0, Risk: "high"},
	}
}

func TestBuildOutliers_DetectsHighZScore(t *testing.T) {
	reports := sampleOutlierReports()
	results := BuildOutliers(reports, 1.5)

	if len(results) != len(reports) {
		t.Fatalf("expected %d results, got %d", len(reports), len(results))
	}

	// The high-score path should be first (sorted by z-score desc)
	if results[0].Path != "secret/app/admin" {
		t.Errorf("expected secret/app/admin as top outlier, got %s", results[0].Path)
	}
	if !results[0].IsOutlier {
		t.Errorf("expected secret/app/admin to be flagged as outlier")
	}
}

func TestBuildOutliers_NormalPathsNotFlagged(t *testing.T) {
	reports := sampleOutlierReports()
	results := BuildOutliers(reports, 1.5)

	for _, r := range results {
		if r.Path != "secret/app/admin" && r.IsOutlier {
			t.Errorf("expected %s not to be an outlier", r.Path)
		}
	}
}

func TestBuildOutliers_Empty(t *testing.T) {
	results := BuildOutliers(nil, 2.0)
	if results != nil {
		t.Errorf("expected nil result for empty input")
	}
}

func TestBuildOutliers_UniformScores(t *testing.T) {
	reports := []ScoredReport{
		{Path: "secret/a", Score: 50.0},
		{Path: "secret/b", Score: 50.0},
		{Path: "secret/c", Score: 50.0},
	}
	results := BuildOutliers(reports, 1.5)
	for _, r := range results {
		if r.IsOutlier {
			t.Errorf("no outliers expected when all scores are equal, got %s", r.Path)
		}
	}
}

func TestBuildOutliers_SortedByZScoreDesc(t *testing.T) {
	reports := sampleOutlierReports()
	results := BuildOutliers(reports, 2.0)

	for i := 1; i < len(results); i++ {
		if results[i].ZScore > results[i-1].ZScore {
			t.Errorf("results not sorted by z-score descending at index %d", i)
		}
	}
}
