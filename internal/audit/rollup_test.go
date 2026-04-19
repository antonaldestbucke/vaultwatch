package audit

import (
	"testing"
)

func sampleRollupReports() []ScoredReport {
	return []ScoredReport{
		{Path: "/secret/prod/db", Score: 40, Risk: "high"},
		{Path: "/secret/prod/api", Score: 70, Risk: "medium"},
		{Path: "/secret/prod/cache", Score: 90, Risk: "low"},
		{Path: "/secret/staging/db", Score: 85, Risk: "low"},
		{Path: "/secret/staging/api", Score: 30, Risk: "high"},
	}
}

func TestBuildRollup_Depth2(t *testing.T) {
	reports := sampleRollupReports()
	rollup := BuildRollup(reports, 2)

	if len(rollup.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(rollup.Entries))
	}

	for _, e := range rollup.Entries {
		if e.Prefix == "/secret/prod" {
			if e.TotalPaths != 3 {
				t.Errorf("expected 3 total for prod, got %d", e.TotalPaths)
			}
			if e.DriftedPaths != 2 {
				t.Errorf("expected 2 drifted for prod, got %d", e.DriftedPaths)
			}
		}
		if e.Prefix == "/secret/staging" {
			if e.TotalPaths != 2 {
				t.Errorf("expected 2 total for staging, got %d", e.TotalPaths)
			}
			if e.DriftedPaths != 1 {
				t.Errorf("expected 1 drifted for staging, got %d", e.DriftedPaths)
			}
		}
	}
}

func TestBuildRollup_Empty(t *testing.T) {
	rollup := BuildRollup([]ScoredReport{}, 2)
	if len(rollup.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(rollup.Entries))
	}
}

func TestBuildRollup_DriftRate(t *testing.T) {
	reports := []ScoredReport{
		{Path: "/a/b/c", Risk: "high"},
		{Path: "/a/b/d", Risk: "low"},
	}
	rollup := BuildRollup(reports, 2)
	if len(rollup.Entries) != 1 {
		t.Fatalf("expected 1 entry")
	}
	if rollup.Entries[0].DriftRate != 50.0 {
		t.Errorf("expected 50%% drift rate, got %.1f", rollup.Entries[0].DriftRate)
	}
}

func TestPrefixAtDepth(t *testing.T) {
	if got := prefixAtDepth("/secret/prod/db", 2); got != "/secret/prod" {
		t.Errorf("unexpected prefix: %s", got)
	}
	if got := prefixAtDepth("/secret/prod/db", 0); got != "/secret/prod/db" {
		t.Errorf("unexpected prefix: %s", got)
	}
}
