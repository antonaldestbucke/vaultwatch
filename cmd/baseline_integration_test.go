package cmd

import (
	"testing"
)

func TestBaselineCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "baseline", "--name", "v1")
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
}

func TestBaselineCmd_MissingName(t *testing.T) {
	_, err := executeCommand(rootCmd, "baseline", "--path", "secret/app")
	if err == nil {
		t.Fatal("expected error when --name is missing")
	}
}

func TestBaselineCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "baseline",
		"--path", "secret/app",
		"--name", "v1",
		"--action", "diff",
		"--config", "/nonexistent/config.yaml",
	)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}

func TestBaselineCmd_InvalidAction(t *testing.T) {
	_, err := executeCommand(rootCmd, "baseline",
		"--path", "secret/app",
		"--name", "v1",
		"--action", "unknown",
		"--config", "/nonexistent/config.yaml",
	)
	if err == nil {
		t.Fatal("expected error for invalid action")
	}
}
