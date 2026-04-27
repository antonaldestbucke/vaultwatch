package audit

import (
	"math"
	"testing"
)

func sampleAnomalyReports() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/app/db", Score: 10.0, Reports: []CompareReport{{EnvA: "dev", EnvB: "prod"}}},
		{Path: "secret/app/api", Score: 12.0, Reports: []CompareReport{{EnvA: "dev", EnvB: "prod"}}},
		{Path: "secret/app/cache", Score: 11.0, Reports: []CompareReport{{EnvA: "dev", EnvB: "prod"}}},
		{Path: "secret/app/outlier", Score: 95.0, Reports: []CompareReport{{EnvA: "dev", EnvB: "prod"}}},
	}
}

func TestBuildAnomalies_DetectsOutlier(t *testing.T) {
	reports := sampleAnomalyReports()
	results := BuildAnomalies(reports, 2.0)

	var anomalies []AnomalyResult
	for _, r := range results {
		if r.IsAnomaly {
			anomalies = append(anomalies, r)
		}
	}

	if len(anomalies) != 1 {
		t.Fatalf("expected 1 anomaly, got %d", len(anomalies))
	}
	if anomalies[0].Path != "secret/app/outlier" {
		t.Errorf("expected outlier path, got %s", anomalies[0].Path)
	}
}

func TestBuildAnomalies_Empty(t *testing.T) {
	results := BuildAnomalies(nil, 2.0)
	if results != nil {
		t.Errorf("expected nil for empty input")
	}
}

func TestBuildAnomalies_UniformScores(t *testing.T) {
	reports := []ScoredReport{
		{Path: "a", Score: 50.0},
		{Path: "b", Score: 50.0},
		{Path: "c", Score: 50.0},
	}
	results := BuildAnomalies(reports, 1.0)
	for _, r := range results {
		if r.IsAnomaly {
			t.Errorf("expected no anomalies for uniform scores, got one at %s", r.Path)
		}
	}
}

func TestBuildAnomalies_ZScoreSign(t *testing.T) {
	reports := sampleAnomalyReports()
	results := BuildAnomalies(reports, 2.0)

	for _, r := range results {
		if r.Path == "secret/app/outlier" && r.ZScore <= 0 {
			t.Errorf("expected positive z-score for high outlier, got %f", r.ZScore)
		}
	}
}

func TestBuildAnomalies_SeverityLabel(t *testing.T) {
	reports := sampleAnomalyReports()
	results := BuildAnomalies(reports, 1.0)

	for _, r := range results {
		if r.IsAnomaly && r.Severity == "" {
			t.Errorf("expected non-empty severity for anomaly at %s", r.Path)
		}
	}
}

func TestBuildAnomalies_MeanIsCorrect(t *testing.T) {
	reports := []ScoredReport{
		{Path: "x", Score: 10.0},
		{Path: "y", Score: 20.0},
		{Path: "z", Score: 30.0},
	}
	results := BuildAnomalies(reports, 10.0)
	for _, r := range results {
		if math.Abs(r.MeanScore-20.0) > 0.001 {
			t.Errorf("expected mean 20.0, got %f", r.MeanScore)
		}
	}
}
