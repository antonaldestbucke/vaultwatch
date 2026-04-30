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
	sensitivityConfigFile string
	sensitivityPath       string
	sensitivityJSON       bool
	sensitivityMinRisk    string
)

var sensitivityCmd = &cobra.Command{
	Use:   "sensitivity",
	Short: "Classify secret paths by key sensitivity patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(sensitivityConfigFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}
		reports, err := audit.ComparePathAcrossEnvs(clients, sensitivityPath)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}
		results := audit.BuildSensitivity(reports)

		if sensitivityMinRisk != "" {
			results = filterBySensitivityRisk(results, sensitivityMinRisk)
		}

		if sensitivityJSON {
			return json.NewEncoder(os.Stdout).Encode(results)
		}
		for _, r := range results {
			fmt.Printf("%-50s  %-8s  score=%.2f  matched=%v\n",
				r.Path, r.Label, r.Score, r.MatchedKeys)
		}
		return nil
	},
}

func filterBySensitivityRisk(results []audit.SensitivityResult, minRisk string) []audit.SensitivityResult {
	order := map[string]int{"low": 0, "medium": 1, "high": 2, "critical": 3}
	min := order[minRisk]
	var out []audit.SensitivityResult
	for _, r := range results {
		if order[r.Label] >= min {
			out = append(out, r)
		}
	}
	return out
}

func init() {
	sensitivityCmd.Flags().StringVarP(&sensitivityConfigFile, "config", "c", "configs/vaultwatch.yaml", "Path to config file")
	sensitivityCmd.Flags().StringVarP(&sensitivityPath, "path", "p", "", "Secret path to evaluate (required)")
	sensitivityCmd.Flags().BoolVar(&sensitivityJSON, "json", false, "Output as JSON")
	sensitivityCmd.Flags().StringVar(&sensitivityMinRisk, "min-risk", "", "Minimum risk label to display (low|medium|high|critical)")
	_ = sensitivityCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(sensitivityCmd)
}
