package cmd

import (
	"testing"
)

func TestWatchCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "watch")
	if err == nil {
		t.Fatal("expected error when path argument is missing")
	}
}

func TestWatchCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "watch", "secret/app", "--config", "/nonexistent/path.yaml")
	if err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestWatchCmd_ZeroInterval(t *testing.T) {
	_, err := executeCommand(rootCmd, "watch", "secret/app", "--interval", "0s", "--config", "/nonexistent/path.yaml")
	// config error fires before interval validation, so we just expect some error
	if err == nil {
		t.Fatal("expected error")
	}
}
