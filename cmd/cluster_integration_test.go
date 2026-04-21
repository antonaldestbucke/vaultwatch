package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/vaultwatch/internal/audit"
)

func writeClusterScores(t *testing.T, reports []audit.ScoredReport) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "scores.json")
	data, _ := json.Marshal(reports)
	_ = os.WriteFile(path, data, 0644)
	return path
}

func TestClusterCmd_MissingScores(t *testing.T) {
	_, err := executeCommand(rootCmd, "cluster")
	if err == nil {
		t.Fatal("expected error for missing --scores flag")
	}
}

func TestClusterCmd_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(bad, []byte("not-json"), 0644)
	_, err := executeCommand(rootCmd, "cluster", "--scores", bad)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestClusterCmd_ValidInput(t *testing.T) {
	reports := []audit.ScoredReport{
		{Path: "secret/a", Score: 0.1, Risk: "low"},
		{Path: "secret/b", Score: 0.9, Risk: "high"},
	}
	path := writeClusterScores(t, reports)
	out, err := executeCommand(rootCmd, "cluster", "--scores", path, "--threshold", "0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatal("expected non-empty output")
	}
}
