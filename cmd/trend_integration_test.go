package cmd

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"vaultwatch/internal/audit"
)

func writeScoredReports(t *testing.T, reports []audit.ScoredReport) string {
	t.Helper()
	f, err := os.CreateTemp("", "scored-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(reports); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestTrendCmd_MissingFile(t *testing.T) {
	err := executeCommand("trend")
	if err == nil {
		t.Error("expected error when --file is missing")
	}
}

func TestTrendCmd_ValidInput(t *testing.T) {
	reports := []audit.ScoredReport{
		{Timestamp: time.Now().Add(-2 * time.Hour), Score: 0.5, Drifted: 5, Total: 10},
		{Timestamp: time.Now().Add(-1 * time.Hour), Score: 0.7, Drifted: 7, Total: 10},
		{Timestamp: time.Now(), Score: 0.3, Drifted: 3, Total: 10},
	}
	path := writeScoredReports(t, reports)
	err := executeCommand("trend", "--file", path)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestTrendCmd_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "bad-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	err := executeCommand("trend", "--file", f.Name())
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
