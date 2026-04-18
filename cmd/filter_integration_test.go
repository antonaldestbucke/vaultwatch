package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterCmd_MissingPath(t *testing.T) {
	args := []string{"filter", "--config", "nonexistent.yaml"}
	out, err := executeCommand(rootCmd, args...)
	assert.Error(t, err)
	_ = out
}

func TestFilterCmd_InvalidConfig(t *testing.T) {
	args := []string{"filter", "--config", "nonexistent.yaml", "--path", "secret/app"}
	_, err := executeCommand(rootCmd, args...)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "loading config")
}

func TestFilterCmd_OnlyDiffsFlag(t *testing.T) {
	args := []string{"filter", "--config", "nonexistent.yaml", "--path", "secret/app", "--only-diffs"}
	_, err := executeCommand(rootCmd, args...)
	assert.Error(t, err)
}
