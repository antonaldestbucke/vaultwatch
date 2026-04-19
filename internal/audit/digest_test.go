package audit

import (
	"testing"
)

func sampleDigestReports() []CompareReport {
	return []CompareReport{
		{
			Path: "secret/app/db",
			Diff: DiffResult{OnlyInA: []string{"password"}, OnlyInB: []string{}},
		},
		{
			Path: "secret/app/api",
			Diff: DiffResult{OnlyInA: []string{}, OnlyInB: []string{}},
		},
	}
}

func TestBuildDigest_ReturnsHash(t *testing.T) {
	reports := sampleDigestReports()
	entry, err := BuildDigest(reports)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Hash == "" {
		t.Error("expected non-empty hash")
	}
	if entry.PathCount != 2 {
		t.Errorf("expected path_count=2, got %d", entry.PathCount)
	}
	if entry.DriftCount != 1 {
		t.Errorf("expected drift_count=1, got %d", entry.DriftCount)
	}
}

func TestBuildDigest_Deterministic(t *testing.T) {
	reports := sampleDigestReports()
	a, _ := BuildDigest(reports)
	b, _ := BuildDigest(reports)
	if a.Hash != b.Hash {
		t.Error("expected same hash for same input")
	}
}

func TestBuildDigest_Empty(t *testing.T) {
	_, err := BuildDigest([]CompareReport{})
	if err == nil {
		t.Error("expected error for empty reports")
	}
}

func TestDigestsMatch_SameHash(t *testing.T) {
	reports := sampleDigestReports()
	a, _ := BuildDigest(reports)
	b, _ := BuildDigest(reports)
	if !DigestsMatch(a, b) {
		t.Error("expected digests to match")
	}
}

func TestDigestsMatch_DifferentHash(t *testing.T) {
	a, _ := BuildDigest(sampleDigestReports())
	b, _ := BuildDigest([]CompareReport{
		{Path: "secret/other", Diff: DiffResult{}},
	})
	if DigestsMatch(a, b) {
		t.Error("expected digests to differ")
	}
}
