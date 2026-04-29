package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var (
	maturityScoresFile string
	maturityJSONOut    bool
)

var maturityCmd = &cobra.Command{
	Use:   "maturity",
	Short: "Assess secret path maturity across environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(maturityScoresFile)
		if err != nil {
			return fmt.Errorf("reading scores file: %w", err)
		}
		var scored []audit.ScoredReport
		if err := json.Unmarshal(data, &scored); err != nil {
			return fmt.Errorf("parsing scores file: %w", err)
		}

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		envs := make([]string, 0, len(cfg.Environments))
		for _, e := range cfg.Environments {
			envs = append(envs, e.Name)
		}

		results := audit.BuildMaturity(scored, envs)

		if maturityJSONOut {
			out, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		if len(results) == 0 {
			fmt.Println("No maturity data available.")
			return nil
		}
		fmt.Printf("%-40s %-12s %-12s %s\n", "PATH", "LEVEL", "COVERAGE", "AVG SCORE")
		for _, r := range results {
			fmt.Printf("%-40s %-12s %-12.0f%% %.1f\n",
				r.Path, r.Level, r.EnvCoverage*100, r.AvgDriftScore)
			for _, n := range r.Notes {
				fmt.Printf("  note: %s\n", n)
			}
		}
		return nil
	},
}

func init() {
	maturityCmd.Flags().StringVar(&maturityScoresFile, "scores", "", "Path to scored reports JSON file (required)")
	maturityCmd.Flags().BoolVar(&maturityJSONOut, "json", false, "Output results as JSON")
	_ = maturityCmd.MarkFlagRequired("scores")
	rootCmd.AddCommand(maturityCmd)
}
