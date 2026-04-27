package cmd

import (
	"bytes"
	"testing"
)

func TestEntropyCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "entropy", "--config", "../configs/vaultwatch.example.yaml")
	if err == nil {
		t.Fatal("expected error when --path is missing")
	}
}

func TestEntropyCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "entropy",
		"--path", "secret/app",
		"--config", "/nonexistent/config.yaml",
	)
	if err == nil {
		t.Fatal("expected error with invalid config path")
	}
}

func TestEntropyCmd_JSONFlagAccepted(t *testing.T) {
	buf := &bytes.Buffer{}
	_ = buf
	// Verify the flag is registered without panicking.
	cmd, _, err := rootCmd.Find([]string{"entropy"})
	if err != nil || cmd == nil {
		t.Fatal("entropy command not found")
	}
	if cmd.Flags().Lookup("json") == nil {
		t.Error("expected --json flag to be registered")
	}
}

func TestEntropyCmd_MinRiskFlagAccepted(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"entropy"})
	if err != nil || cmd == nil {
		t.Fatal("entropy command not found")
	}
	if cmd.Flags().Lookup("min-risk") == nil {
		t.Error("expected --min-risk flag to be registered")
	}
}
