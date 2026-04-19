package audit

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func sampleScoredForRemediate() []ScoredReport {
	return []ScoredReport{
		{
			Risk: "high",
			Report: CompareReport{
				Path:         "secret/app",
				Environments: []string{"prod", "staging"},
				Diffs: []DiffResult{
					{Path: "secret/app", OnlyInA: []string{"DB_PASS"}, OnlyInB: []string{"API_KEY"}},
				},
			},
		},
	}
}

func TestBuildRemediationPlan_WithDiffs(t *testing.T) {
	plan := BuildRemediationPlan(sampleScoredForRemediate())
	if len(plan.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(plan.Actions))
	}
}

func TestBuildRemediationPlan_NoDiffs(t *testing.T) {
	reports := []ScoredReport{
		{
			Report: CompareReport{
				Path:         "secret/clean",
				Environments: []string{"prod", "staging"},
				Diffs:        []DiffResult{},
			},
		},
	}
	plan := BuildRemediationPlan(reports)
	if len(plan.Actions) != 0 {
		t.Fatalf("expected 0 actions, got %d", len(plan.Actions))
	}
}

func TestBuildRemediationPlan_ActionEnvironments(t *testing.T) {
	plan := BuildRemediationPlan(sampleScoredForRemediate())
	for _, a := range plan.Actions {
		if a.Environment == "unknown" {
			t.Errorf("expected known environment, got unknown for path %s", a.Path)
		}
	}
}

func TestPrintRemediationPlan_NoDiff(t *testing.T) {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	PrintRemediationPlan(RemediationPlan{})
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("No remediation")) {
		t.Error("expected no-action message")
	}
}

func TestPrintRemediationPlan_WithActions(t *testing.T) {
	plan := BuildRemediationPlan(sampleScoredForRemediate())
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	PrintRemediationPlan(plan)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("add_key")) {
		t.Error("expected add_key action in output")
	}
}
