package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoverageCmd_MissingPath(t *testing.T) {
	_, err := executeCommand(rootCmd, "coverage",
		"--config", "../configs/vaultwatch.example.yaml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path")
}

func TestCoverageCmd_InvalidConfig(t *testing.T) {
	_, err := executeCommand(rootCmd, "coverage",
		"--path", "secret/",
		"--config", "/nonexistent/config.yaml")
	assert.Error(t, err)
}

func TestCoverageCmd_AllFlagAccepted(t *testing.T) {
	// Verify the --all flag is registered and parseable (config load will fail).
	_, err := executeCommand(rootCmd, "coverage",
		"--path", "secret/",
		"--config", "/nonexistent/config.yaml",
		"--all")
	assert.Error(t, err)
	// Error should be about config, not flag parsing.
	assert.NotContains(t, err.Error(), "unknown flag")
}
