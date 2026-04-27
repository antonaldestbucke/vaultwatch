package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultwatch/internal/audit"
)

var (
	consensusScoresFile string
	consensusThreshold  float64
	consensusJSONOutput bool
)

var consensusCmd = &cobra.Command{
	Use:   "consensus",
	Short: "Analyze key-set agreement across environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(consensusScoresFile)
		if err != nil {
			return fmt.Errorf("reading scores file: %w", err)
		}

		var scored []audit.ScoredReport
		if err := json.Unmarshal(data, &scored); err != nil {
			return fmt.Errorf("parsing scores file: %w", err)
		}

		results := audit.BuildConsensus(scored, consensusThreshold)

		if consensusJSONOutput {
			out, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		if len(results) == 0 {
			fmt.Println("No consensus data available.")
			return nil
		}

		fmt.Printf("%-40s  %-8s  %-10s  %s\n", "PATH", "AGREE%", "CONSENSUS", "DISSENT ENVS")
		for _, r := range results {
			consensusLabel := "no"
			if r.Consensus {
				consensusLabel = "yes"
			}
			dissent := "-"
			if len(r.DissentEnvs) > 0 {
				dissent = joinStrings(r.DissentEnvs, ",")
			}
			fmt.Printf("%-40s  %-8.1f  %-10s  %s\n", r.Path, r.AgreementPct, consensusLabel, dissent)
		}
		return nil
	},
}

func init() {
	consensusCmd.Flags().StringVar(&consensusScoresFile, "scores", "", "Path to scored reports JSON file (required)")
	consensusCmd.Flags().Float64Var(&consensusThreshold, "threshold", 75.0, "Minimum agreement percentage to consider consensus reached")
	consensusCmd.Flags().BoolVar(&consensusJSONOutput, "json", false, "Output results as JSON")
	_ = consensusCmd.MarkFlagRequired("scores")
	rootCmd.AddCommand(consensusCmd)
}
