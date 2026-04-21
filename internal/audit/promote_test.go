package audit

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func sampleScoredForPromote() []ScoredReport {
	return []ScoredReport{
		{
			Report: CompareReport{
				Path:    "secret/app/config",
				EnvA:    "staging",
				EnvB:    "production",
				OnlyInA: []string{"NEW_FEATURE_FLAG", "DEBUG_MODE"},
				OnlyInB: []string{},
			},
			RiskLevel: "high",
			Score:     40,
		},
		{
			Report: CompareReport{
				Path:    "secret/db/creds",
				EnvA:    "staging",
				EnvB:    "production",
				OnlyInA: []string{},
				OnlyInB: []string{},
			},
			RiskLevel: "low",
			Score:     100,
		},
	}
}

func TestBuildPromotionPlan_WithDiffs(t *testing.T) {
	plan := BuildPromotionPlan(sampleScoredForPromote(), "staging", "production")

	if plan.FromEnv != "staging" || plan.ToEnv != "production" {
		t.Errorf("unexpected env pair: %s → %s", plan.FromEnv, plan.ToEnv)
	}
	if len(plan.Actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(plan.Actions))
	}
	for _, a := range plan.Actions {
		if a.FromEnv != "staging" || a.ToEnv != "production" {
			t.Errorf("wrong env on action: %+v", a)
		}
		if a.Path != "secret/app/config" {
			t.Errorf("unexpected path: %s", a.Path)
		}
	}
}

func TestBuildPromotionPlan_NoDiffs(t *testing.T) {
	reports := []ScoredReport{
		{
			Report: CompareReport{
				Path: "secret/clean", EnvA: "staging", EnvB: "production",
				OnlyInA: []string{}, OnlyInB: []string{},
			},
		},
	}
	plan := BuildPromotionPlan(reports, "staging", "production")
	if len(plan.Actions) != 0 {
		t.Errorf("expected 0 actions, got %d", len(plan.Actions))
	}
}

func TestPrintPromotionPlan_NoDiff(t *testing.T) {
	plan := PromotionPlan{FromEnv: "dev", ToEnv: "prod", Actions: nil}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	PrintPromotionPlan(plan)
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !bytes.Contains(buf.Bytes(), []byte("No promotion actions")) {
		t.Errorf("expected no-action message, got: %s", buf.String())
	}
}

func TestBuildPromotionPlan_WrongFromEnv(t *testing.T) {
	plan := BuildPromotionPlan(sampleScoredForPromote(), "production", "staging")
	if len(plan.Actions) != 0 {
		t.Errorf("expected 0 actions when fromEnv doesn't match EnvA, got %d", len(plan.Actions))
	}
}
