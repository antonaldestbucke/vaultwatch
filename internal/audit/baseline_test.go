package audit

import (
	"os"
	"testing"
)

func sampleReportsForBaseline() []CompareReport {
	return []CompareReport{
		{Path: "secret/app", EnvA: "dev", EnvB: "prod", OnlyInA: []string{"debug"}, OnlyInB: []string{}},
		{Path: "secret/db", EnvA: "dev", EnvB: "prod", OnlyInA: []string{}, OnlyInB: []string{}},
	}
}

func TestSaveAndLoadBaseline_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	reports := sampleReportsForBaseline()

	if err := SaveBaseline(dir, "v1", "secret/app", reports); err != nil {
		t.Fatalf("SaveBaseline: %v", err)
	}

	b, err := LoadBaseline(dir, "v1")
	if err != nil {
		t.Fatalf("LoadBaseline: %v", err)
	}

	if b.Name != "v1" {
		t.Errorf("expected name v1, got %s", b.Name)
	}
	if len(b.Data["dev"]) != 1 || b.Data["dev"][0] != "debug" {
		t.Errorf("unexpected dev keys: %v", b.Data["dev"])
	}
}

func TestLoadBaseline_MissingFile(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadBaseline(dir, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing baseline")
	}
}

func TestDiffAgainstBaseline_NoChange(t *testing.T) {
	dir := t.TempDir()
	reports := sampleReportsForBaseline()
	_ = SaveBaseline(dir, "v1", "secret/app", reports)
	b, _ := LoadBaseline(dir, "v1")

	changed := DiffAgainstBaseline(b, reports)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
}

func TestDiffAgainstBaseline_WithChange(t *testing.T) {
	dir := t.TempDir()
	reports := sampleReportsForBaseline()
	_ = SaveBaseline(dir, "v1", "secret/app", reports)
	b, _ := LoadBaseline(dir, "v1")

	modified := []CompareReport{
		{Path: "secret/app", EnvA: "dev", EnvB: "prod", OnlyInA: []string{"debug", "newkey"}, OnlyInB: []string{}},
	}
	changed := DiffAgainstBaseline(b, modified)
	if len(changed) != 1 || changed[0] != "secret/app" {
		t.Errorf("expected secret/app changed, got %v", changed)
	}
}

func TestEqualSlices(t *testing.T) {
	if !equalSlices([]string{"a", "b"}, []string{"b", "a"}) {
		t.Error("expected equal")
	}
	if equalSlices([]string{"a"}, []string{"a", "b"}) {
		t.Error("expected not equal")
	}
}

func init() {
	_ = os.Getenv // suppress unused import
}
