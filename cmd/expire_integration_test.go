package cmd

import (
	"testing"
)

func TestExpireCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "expire",
		"--expiry-file", "some.json",
	)
	if err == nil {
		t.Error("expected error when --path is missing")
	}
}

func TestExpireCmd_MissingExpiryFile(t *testing.T) {
	_, err := executeCommand(rootCmd, "expire",
		"--path", "secret/prod/db",
	)
	if err == nil {
		t.Error("expected error when --expiry-file is missing")
	}
}

func TestExpireCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "expire",
		"--path", "secret/prod/db",
		"--expiry-file", "expiry.json",
		"--config", "/nonexistent/config.yaml",
	)
	if err == nil {
		t.Error("expected error for invalid config path")
	}
}
