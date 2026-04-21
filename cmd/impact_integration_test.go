package cmd

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vaultwatch/internal/audit"
)

func writeImpactScores(t *testing.T, reports []audit.ScoredReport) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "scores-*.json")
	require.NoError(t, err)
	require.NoError(t, json.NewEncoder(f).Encode(reports))
	return f.Name()
}

func TestImpactCmd_MissingScores(t *testing.T) {
	_, err := executeCommand(rootCmd, "impact")
	assert.Error(t, err)
}

func TestImpactCmd_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "bad-*.json")
	_, _ = f.WriteString("not-json")
	_, err := executeCommand(rootCmd, "impact", "--scores", f.Name())
	assert.Error(t, err)
}

func TestImpactCmd_ValidInput(t *testing.T) {
	reports := []audit.ScoredReport{
		{
			Score: 50,
			Risk:  "medium",
			Report: audit.CompareReport{
				Path:    "secret/db",
				Envs:    []string{"prod", "staging"},
				OnlyInA: []string{"password"},
				OnlyInB: []string{},
			},
		},
	}
	scoresFile := writeImpactScores(t, reports)
	out, err := executeCommand(rootCmd, "impact", "--scores", scoresFile)
	require.NoError(t, err)
	assert.Contains(t, out, "secret/db")
}

func TestImpactCmd_JSONOutput(t *testing.T) {
	reports := []audit.ScoredReport{
		{
			Score: 30,
			Risk:  "high",
			Report: audit.CompareReport{
				Path:    "secret/api",
				Envs:    []string{"prod", "staging", "dev"},
				OnlyInA: []string{"k1", "k2", "k3", "k4", "k5"},
				OnlyInB: []string{},
			},
		},
	}
	scoresFile := writeImpactScores(t, reports)
	out, err := executeCommand(rootCmd, "impact", "--scores", scoresFile, "--json")
	require.NoError(t, err)
	var result map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(out), &result))
	assert.Equal(t, float64(1), result["high_count"])
}
