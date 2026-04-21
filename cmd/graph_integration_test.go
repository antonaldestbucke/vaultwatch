package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"vaultwatch/internal/audit"
)

func writeGraphScores(t *testing.T, dir string, reports []audit.ScoredReport) string {
	t.Helper()
	data, _ := json.Marshal(reports)
	p := filepath.Join(dir, "scores.json")
	os.WriteFile(p, data, 0644)
	return p
}

func TestGraphCmd_MissingScores(t *testing.T) {
	_, err := executeCommand(rootCmd, "graph")
	if err == nil {
		t.Error("expected error when --scores is missing")
	}
}

func TestGraphCmd_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	os.WriteFile(p, []byte("not-json"), 0644)
	_, err := executeCommand(rootCmd, "graph", "--scores", p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestGraphCmd_ValidInput(t *testing.T) {
	dir := t.TempDir()
	reports := []audit.ScoredReport{
		{Path: "secret/app/db", Score: 40, RiskLevel: "high", Keys: []string{"pass"}},
		{Path: "secret/app", Score: 70, RiskLevel: "medium", Keys: []string{"token"}},
	}
	p := writeGraphScores(t, dir, reports)
	out, err := executeCommand(rootCmd, "graph", "--scores", p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}

func TestGraphCmd_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	reports := []audit.ScoredReport{
		{Path: "secret/x", Score: 80, RiskLevel: "low", Keys: []string{"key"}},
	}
	p := writeGraphScores(t, dir, reports)
	out, err := executeCommand(rootCmd, "graph", "--scores", p, "--json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Errorf("expected valid JSON output, got: %s", out)
	}
}
