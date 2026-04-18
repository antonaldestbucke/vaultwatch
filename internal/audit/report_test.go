package audit

import (
	"bytes"
	"strings"
	"testing"
)

func makeReport(path, envA, envB string, onlyA, onlyB []string) PathReport {
	return PathReport{
		Path: path,
		EnvA: envA,
		EnvB: envB,
		Diff: DiffResult{OnlyInA: onlyA, OnlyInB: onlyB},
	}
}

func TestPrintTextReport_NoDiff(t *testing.T) {
	var buf bytes.Buffer
	reports := []PathReport{makeReport("secret/app", "dev", "prod", nil, nil)}
	PrintTextReport(&buf, reports)
	if !strings.Contains(buf.String(), "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", buf.String())
	}
}

func TestPrintTextReport_WithDiff(t *testing.T) {
	var buf bytes.Buffer
	reports := []PathReport{
		makeReport("secret/app", "dev", "prod", []string{"DB_PASS"}, []string{"API_KEY"}),
	}
	PrintTextReport(&buf, reports)
	out := buf.String()
	if !strings.Contains(out, "[DIFF]") {
		t.Errorf("expected [DIFF] in output")
	}
	if !strings.Contains(out, "DB_PASS") {
		t.Errorf("expected DB_PASS in output")
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in output")
	}
}

func TestSummary_Mixed(t *testing.T) {
	reports := []PathReport{
		makeReport("secret/a", "dev", "prod", nil, nil),
		makeReport("secret/b", "dev", "prod", []string{"X"}, nil),
	}
	s := Summary(reports)
	if s != "2 path(s) checked, 1 with differences" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestSummary_AllClean(t *testing.T) {
	reports := []PathReport{
		makeReport("secret/a", "dev", "prod", nil, nil),
	}
	s := Summary(reports)
	if s != "1 path(s) checked, 0 with differences" {
		t.Errorf("unexpected summary: %s", s)
	}
}
