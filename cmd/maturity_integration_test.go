package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"vaultwatch/internal/audit"
)

func writeMaturityScores(t *testing.T, dir string, reports []audit.ScoredReport) string {
	t.Helper()
	data, _ := json.Marshal(reports)
	p := filepath.Join(dir, "scores.json")
	_ = os.WriteFile(p, data, 0644)
	return p
}

func TestMaturityCmd_MissingScores(t *testing.T) {
	_, err := executeCommand(rootCmd, "maturity")
	if err == nil {
		t.Error("expected error when --scores flag is missing")
	}
}

func TestMaturityCmd_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(p, []byte("not-json"), 0644)

	cfgPath := filepath.Join(dir, "cfg.yaml")
	_ = os.WriteFile(cfgPath, []byte("environments:\n  - name: dev\n    address: http://localhost\n    token: t"), 0644)

	_, err := executeCommand(rootCmd, "maturity", "--scores", p, "--config", cfgPath)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestMaturityCmd_ValidInput(t *testing.T) {
	dir := t.TempDir()
	reports := []audit.ScoredReport{
		{Path: "secret/app", Env: "dev", Score: 95, Risk: "low"},
		{Path: "secret/app", Env: "prod", Score: 90, Risk: "low"},
	}
	scoresPath := writeMaturityScores(t, dir, reports)

	cfgPath := filepath.Join(dir, "cfg.yaml")
	_ = os.WriteFile(cfgPath, []byte("environments:\n  - name: dev\n    address: http://localhost\n    token: t\n  - name: prod\n    address: http://localhost\n    token: t"), 0644)

	out, err := executeCommand(rootCmd, "maturity", "--scores", scoresPath, "--config", cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}
