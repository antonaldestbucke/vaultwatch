package audit

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func sampleHeatmapReports() []ScoredReport {
	return []ScoredReport{
		{
			Report: CompareReport{Path: "secret/app", EnvA: "dev", EnvB: "prod", OnlyInA: []string{"KEY_A"}, OnlyInB: []string{}},
			Score: 60,
			Risk:  "medium",
		},
		{
			Report: CompareReport{Path: "secret/app", EnvA: "dev", EnvB: "prod", OnlyInA: []string{"KEY_B"}, OnlyInB: []string{}},
			Score: 60,
			Risk:  "medium",
		},
		{
			Report: CompareReport{Path: "secret/db", EnvA: "dev", EnvB: "staging", OnlyInA: []string{}, OnlyInB: []string{"PASS"}},
			Score: 40,
			Risk:  "low",
		},
		{
			Report: CompareReport{Path: "secret/clean", EnvA: "dev", EnvB: "prod", OnlyInA: []string{}, OnlyInB: []string{}},
			Score: 100,
			Risk:  "none",
		},
	}
}

func TestBuildHeatmap_CountsDrift(t *testing.T) {
	h := BuildHeatmap(sampleHeatmapReports())
	if len(h.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h.Entries))
	}
	if h.Entries[0].Path != "secret/app" {
		t.Errorf("expected secret/app first, got %s", h.Entries[0].Path)
	}
	if h.Entries[0].DriftCount != 2 {
		t.Errorf("expected drift count 2, got %d", h.Entries[0].DriftCount)
	}
}

func TestBuildHeatmap_Empty(t *testing.T) {
	h := BuildHeatmap([]ScoredReport{})
	if len(h.Entries) != 0 {
		t.Errorf("expected empty heatmap")
	}
}

func TestBuildHeatmap_TotalMatchesInput(t *testing.T) {
	reports := sampleHeatmapReports()
	h := BuildHeatmap(reports)
	if h.Total != len(reports) {
		t.Errorf("expected total %d, got %d", len(reports), h.Total)
	}
}

func TestPrintHeatmap_WithDrift(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintHeatmap(BuildHeatmap(sampleHeatmapReports()))
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("secret/app")) {
		t.Errorf("expected secret/app in output")
	}
}

func TestPrintHeatmap_NoDrift(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintHeatmap(Heatmap{})
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("No drift")) {
		t.Errorf("expected no-drift message")
	}
}
