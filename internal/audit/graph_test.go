package audit

import (
	"strings"
	"testing"
)

func sampleGraphReports() []ScoredReport {
	return []ScoredReport{
		{
			Path:      "secret/app/db",
			Score:     45,
			RiskLevel: "high",
			Keys:      []string{"password"},
		},
		{
			Path:      "secret/app/cache",
			Score:     90,
			RiskLevel: "low",
			Keys:      []string{"url"},
		},
		{
			Path:      "secret/app",
			Score:     60,
			RiskLevel: "medium",
			Keys:      []string{"token"},
		},
	}
}

func TestBuildGraph_NodeCount(t *testing.T) {
	reports := sampleGraphReports()
	deps := map[string][]string{
		"secret/app": {"secret/app/db", "secret/app/cache"},
	}
	result := BuildGraph(reports, deps)
	if len(result.Nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(result.Nodes))
	}
}

func TestBuildGraph_EdgesPopulated(t *testing.T) {
	reports := sampleGraphReports()
	deps := map[string][]string{
		"secret/app": {"secret/app/db", "secret/app/cache"},
	}
	result := BuildGraph(reports, deps)
	children, ok := result.Edges["secret/app"]
	if !ok {
		t.Fatal("expected edge for secret/app")
	}
	if len(children) != 2 {
		t.Errorf("expected 2 children, got %d", len(children))
	}
}

func TestBuildGraph_DriftedFlag(t *testing.T) {
	reports := sampleGraphReports()
	result := BuildGraph(reports, nil)
	for _, node := range result.Nodes {
		if node.Path == "secret/app/db" && !node.Drifted {
			t.Error("expected secret/app/db to be marked drifted")
		}
		if node.Path == "secret/app/cache" && node.Drifted {
			t.Error("expected secret/app/cache to not be marked drifted")
		}
	}
}

func TestPrintGraph_ContainsPaths(t *testing.T) {
	reports := sampleGraphReports()
	deps := map[string][]string{
		"secret/app": {"secret/app/db", "secret/app/cache"},
	}
	result := BuildGraph(reports, deps)
	output := PrintGraph(result)
	if !strings.Contains(output, "secret/app") {
		t.Error("expected output to contain secret/app")
	}
	if !strings.Contains(output, "[DRIFTED]") {
		t.Error("expected output to mark drifted nodes")
	}
}

func TestBuildGraph_Empty(t *testing.T) {
	result := BuildGraph(nil, nil)
	if len(result.Nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(result.Nodes))
	}
}
