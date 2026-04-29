package audit

import (
	"testing"
)

func sampleMaturityReports() []ScoredReport {
	return []ScoredReport{
		{Path: "secret/app", Env: "dev", Score: 95, Risk: "low"},
		{Path: "secret/app", Env: "staging", Score: 92, Risk: "low"},
		{Path: "secret/app", Env: "prod", Score: 91, Risk: "low"},
		{Path: "secret/db", Env: "dev", Score: 60, Risk: "medium"},
		{Path: "secret/db", Env: "staging", Score: 55, Risk: "medium"},
		{Path: "secret/legacy", Env: "dev", Score: 30, Risk: "high"},
	}
}

func TestBuildMaturity_ExemplaryPath(t *testing.T) {
	envs := []string{"dev", "staging", "prod"}
	results := BuildMaturity(sampleMaturityReports(), envs)

	var found *MaturityResult
	for i := range results {
		if results[i].Path == "secret/app" {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected result for secret/app")
	}
	if found.Level != MaturityExemplary {
		t.Errorf("expected exemplary, got %s", found.Level)
	}
	if found.EnvCoverage != 1.0 {
		t.Errorf("expected full coverage, got %f", found.EnvCoverage)
	}
}

func TestBuildMaturity_ImmaturePath(t *testing.T) {
	envs := []string{"dev", "staging", "prod"}
	results := BuildMaturity(sampleMaturityReports(), envs)

	var found *MaturityResult
	for i := range results {
		if results[i].Path == "secret/legacy" {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected result for secret/legacy")
	}
	if found.Level != MaturityImmature {
		t.Errorf("expected immature, got %s", found.Level)
	}
}

func TestBuildMaturity_Empty(t *testing.T) {
	results := BuildMaturity(nil, []string{"dev"})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestBuildMaturity_SortedByPath(t *testing.T) {
	envs := []string{"dev", "staging", "prod"}
	results := BuildMaturity(sampleMaturityReports(), envs)
	for i := 1; i < len(results); i++ {
		if results[i].Path < results[i-1].Path {
			t.Errorf("results not sorted: %s before %s", results[i-1].Path, results[i].Path)
		}
	}
}

func TestBuildMaturity_PartialCoverageNoted(t *testing.T) {
	envs := []string{"dev", "staging", "prod"}
	results := BuildMaturity(sampleMaturityReports(), envs)

	var found *MaturityResult
	for i := range results {
		if results[i].Path == "secret/db" {
			found = &results[i]
		}
	}
	if found == nil {
		t.Fatal("expected result for secret/db")
	}
	hasNote := false
	for _, n := range found.Notes {
		if n == "not present in all environments" {
			hasNote = true
		}
	}
	if !hasNote {
		t.Error("expected coverage note for secret/db")
	}
}

func TestBuildMaturity_NoEnvs(t *testing.T) {
	results := BuildMaturity(sampleMaturityReports(), nil)
	if len(results) != 0 {
		t.Errorf("expected empty results with no envs, got %d", len(results))
	}
}
