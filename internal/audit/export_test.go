package audit

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleReports() []PathReport {
	return []PathReport{
		{
			Path: "secret/app",
			EnvA: "staging",
			EnvB: "production",
			Diff: DiffResult{OnlyInA: []string{"debug"}, OnlyInB: []string{"api_key"}},
		},
		{
			Path: "secret/db",
			EnvA: "staging",
			EnvB: "production",
			Diff: DiffResult{},
		},
	}
}

func TestExportReport_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := ExportReport(&buf, sampleReports(), FormatJSON)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out jsonExport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(out.Reports) != 2 {
		t.Errorf("expected 2 reports, got %d", len(out.Reports))
	}
	if out.Reports[0].Path != "secret/app" {
		t.Errorf("unexpected path: %s", out.Reports[0].Path)
	}
}

func TestExportReport_CSV(t *testing.T) {
	var buf bytes.Buffer
	err := ExportReport(&buf, sampleReports(), FormatCSV)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "path") {
		t.Errorf("expected CSV header, got: %s", lines[0])
	}
}

func TestExportReport_UnsupportedFormat(t *testing.T) {
	var buf bytes.Buffer
	err := ExportReport(&buf, sampleReports(), ExportFormat("xml"))
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}
