package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vaultwatch/internal/audit"
	"github.com/vaultwatch/internal/vault"
)

var (
	impactScoresFile string
	impactJSONOutput bool
	impactMinLevel   string
)

var impactCmd = &cobra.Command{
	Use:   "impact",
	Short: "Assess the blast radius of secret drift across environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(impactScoresFile)
		if err != nil {
			return fmt.Errorf("reading scores file: %w", err)
		}
		var scored []audit.ScoredReport
		if err := json.Unmarshal(data, &scored); err != nil {
			return fmt.Errorf("parsing scores file: %w", err)
		}

		summary := audit.BuildImpact(scored)

		filtered := filterByLevel(summary.Results, impactMinLevel)

		if impactJSONOutput {
			out := map[string]interface{}{
				"total":      summary.Total,
				"high_count": summary.HighCount,
				"results":    filtered,
			}
			return json.NewEncoder(os.Stdout).Encode(out)
		}

		fmt.Printf("Impact Summary — total: %d, high: %d\n\n", summary.Total, summary.HighCount)
		for _, r := range filtered {
			fmt.Printf("[%s] %s — drifted keys: %d, envs: %v\n",
				r.Level, r.Path, r.DriftedKeys, r.AffectedEnvs)
		}
		return nil
	},
}

func filterByLevel(results []audit.ImpactResult, minLevel string) []audit.ImpactResult {
	if minLevel == "" {
		return results
	}
	weight := map[string]int{"low": 1, "medium": 2, "high": 3}
	min := weight[minLevel]
	var out []audit.ImpactResult
	for _, r := range results {
		if weight[string(r.Level)] >= min {
			out = append(out, r)
		}
	}
	return out
}

func init() {
	_ = vault.NewClient // ensure vault import used via transitive deps
	impactCmd.Flags().StringVar(&impactScoresFile, "scores", "", "Path to scored reports JSON file (required)")
	impactCmd.Flags().BoolVar(&impactJSONOutput, "json", false, "Output results as JSON")
	impactCmd.Flags().StringVar(&impactMinLevel, "min-level", "", "Minimum impact level to display (low|medium|high)")
	_ = impactCmd.MarkFlagRequired("scores")
	rootCmd.AddCommand(impactCmd)
}
