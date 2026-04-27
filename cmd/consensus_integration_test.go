package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/yourusername/vaultwatch/internal/audit"
)

func writeConsensusScores(t *testing.T, dir string, reports []audit.ScoredReport) string {
	t.Helper()
	path := filepath.Join(dir, "scores.json")
	data, err := json.Marshal(reports)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, data, 0644))
	return path
}

func TestConsensusCmd_MissingScores(t *testing.T) {
	_, err := executeCommand(rootCmd, "consensus")
	assert.Error(t, err)
}

func TestConsensusCmd_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(p, []byte("not-json"), 0644))
	_, err := executeCommand(rootCmd, "consensus", "--scores", p)
	assert.Error(t, err)
}

func TestConsensusCmd_ValidInput(t *testing.T) {
	dir := t.TempDir()
	reports := []audit.ScoredReport{
		{Environment: "prod", Score: 90, Reports: []audit.DiffReport{{Path: "secret/app", OnlyInA: []string{"key1"}}}},
		{Environment: "staging", Score: 90, Reports: []audit.DiffReport{{Path: "secret/app", OnlyInA: []string{"key1"}}}},
	}
	p := writeConsensusScores(t, dir, reports)
	out, err := executeCommand(rootCmd, "consensus", "--scores", p, "--threshold", "50")
	assert.NoError(t, err)
	assert.Contains(t, out, "secret/app")
}

func TestConsensusCmd_JSONOutput(t *testing.T) {
	dir := t.TempDir()
	reports := []audit.ScoredReport{
		{Environment: "prod", Score: 90, Reports: []audit.DiffReport{{Path: "secret/db", OnlyInA: []string{"pass"}}}},
	}
	p := writeConsensusScores(t, dir, reports)
	out, err := executeCommand(rootCmd, "consensus", "--scores", p, "--json")
	assert.NoError(t, err)
	var result []map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &result))
	assert.NotEmpty(t, result)
}
