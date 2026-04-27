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
	entropyPath      string
	entropyConfig    string
	entropyJSONOut   bool
	entropyMinRisk   string
)

var entropyCmd = &cobra.Command{
	Use:   "entropy",
	Short: "Compute Shannon entropy over drifted key sets across environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(entropyConfig)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("building clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(clients, entropyPath)
		if err != nil {
			return fmt.Errorf("comparing path: %w", err)
		}

		results := audit.BuildEntropy(reports)

		if entropyMinRisk != "" {
			results = filterByEntropyRisk(results, entropyMinRisk)
		}

		if entropyJSONOut {
			return json.NewEncoder(os.Stdout).Encode(results)
		}

		if len(results) == 0 {
			fmt.Println("No entropy drift detected.")
			return nil
		}

		fmt.Printf("%-40s %-10s %-8s %s\n", "PATH", "RISK", "ENTROPY", "KEYS")
		for _, r := range results {
			fmt.Printf("%-40s %-10s %-8.3f %d\n", r.Path, r.Risk, r.Entropy, len(r.Keys))
		}
		return nil
	},
}

func filterByEntropyRisk(results []audit.EntropyResult, minRisk string) []audit.EntropyResult {
	order := map[string]int{"low": 0, "medium": 1, "high": 2, "critical": 3}
	min := order[minRisk]
	out := results[:0]
	for _, r := range results {
		if order[r.Risk] >= min {
			out = append(out, r)
		}
	}
	return out
}

func init() {
	entropyCmd.Flags().StringVarP(&entropyPath, "path", "p", "", "Vault secret path to audit (required)")
	entropyCmd.Flags().StringVarP(&entropyConfig, "config", "c", "configs/vaultwatch.yaml", "Path to config file")
	entropyCmd.Flags().BoolVar(&entropyJSONOut, "json", false, "Output results as JSON")
	entropyCmd.Flags().StringVar(&entropyMinRisk, "min-risk", "", "Minimum risk level to display (low|medium|high|critical)")
	_ = entropyCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(entropyCmd)
}
