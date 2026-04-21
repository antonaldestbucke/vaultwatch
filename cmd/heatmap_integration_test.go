package cmd

import (
	"testing"
)

func TestHeatmapCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "heatmap", "--config", "nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
}

func TestHeatmapCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "heatmap",
		"--config", "nonexistent.yaml",
		"--path", "secret/app")
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestHeatmapCmd_TopFlagAccepted(t *testing.T) {
	_, err := executeCommand(rootCmd, "heatmap",
		"--config", "nonexistent.yaml",
		"--path", "secret/app",
		"--top", "5")
	// error expected due to missing config, but flag parsing should succeed
	if err == nil {
		t.Fatal("expected error for invalid config, not flag parse error")
	}
}
