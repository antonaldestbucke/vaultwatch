package audit

import (
	"testing"
)

func sampleClusterReports() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/a", Score: 0.1, Risk: "low"},
		{Path: "secret/b", Score: 0.15, Risk: "low"},
		{Path: "secret/c", Score: 0.8, Risk: "high"},
		{Path: "secret/d", Score: 0.85, Risk: "high"},
		{Path: "secret/e", Score: 0.5, Risk: "medium"},
	}
}

func TestClusterReports_GroupsByThreshold(t *testing.T) {
	reports := sampleClusterReports()
	clusters := ClusterReports(reports, 0.1)
	if len(clusters) != 3 {
		t.Fatalf("expected 3 clusters, got %d", len(clusters))
	}
}

func TestClusterReports_Empty(t *testing.T) {
	clusters := ClusterReports(nil, 0.1)
	if clusters != nil {
		t.Fatalf("expected nil for empty input")
	}
}

func TestClusterReports_AllInOne(t *testing.T) {
	reports := []ScoredReport{
		{Path: "secret/a", Score: 0.5},
		{Path: "secret/b", Score: 0.55},
		{Path: "secret/c", Score: 0.45},
	}
	clusters := ClusterReports(reports, 0.2)
	if len(clusters) != 1 {
		t.Fatalf("expected 1 cluster, got %d", len(clusters))
	}
	if len(clusters[0].Paths) != 3 {
		t.Fatalf("expected 3 paths in cluster, got %d", len(clusters[0].Paths))
	}
}

func TestClusterReports_AvgScore(t *testing.T) {
	reports := []ScoredReport{
		{Path: "secret/x", Score: 0.2},
		{Path: "secret/y", Score: 0.4},
	}
	clusters := ClusterReports(reports, 0.5)
	if clusters[0].AvgScore != 0.3 {
		t.Fatalf("expected avg 0.3, got %f", clusters[0].AvgScore)
	}
}

func TestClusterReports_CentroidIsMiddle(t *testing.T) {
	reports := []ScoredReport{
		{Path: "secret/a", Score: 0.1},
		{Path: "secret/b", Score: 0.2},
		{Path: "secret/c", Score: 0.3},
	}
	clusters := ClusterReports(reports, 0.5)
	if clusters[0].Centroid != "secret/b" {
		t.Fatalf("expected centroid secret/b, got %s", clusters[0].Centroid)
	}
}
