package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var expireCmd = &cobra.Command{
	Use:   "expire",
	Short: "Evaluate secret expiry rules against current vault state",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, _ := cmd.Flags().GetString("config")
		path, _ := cmd.Flags().GetString("path")
		expireFile, _ := cmd.Flags().GetString("expiry-file")
		scoreFile, _ := cmd.Flags().GetString("score-file")

		if path == "" {
			return fmt.Errorf("--path is required")
		}
		if expireFile == "" {
			return fmt.Errorf("--expiry-file is required")
		}

		cfg, err := vault.LoadConfig(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(clients, path)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}

		var scored []audit.ScoredReport
		if scoreFile != "" {
			scored, err = audit.LoadScoredReports(scoreFile)
			if err != nil {
				return fmt.Errorf("load score file: %w", err)
			}
		} else {
			scored = audit.ScoreReports(reports)
		}

		store, err := audit.LoadExpiry(expireFile)
		if err != nil {
			return fmt.Errorf("load expiry store: %w", err)
		}

		results := audit.EvaluateExpiry(scored, store, time.Now())
		if len(results) == 0 {
			fmt.Println("No expiry rules matched any paths.")
			return nil
		}

		return printExpiryResults(results)
	},
}

// printExpiryResults writes expiry evaluation results to stdout in a
// tab-aligned table with columns: ENV, PATH, STATUS, MESSAGE.
func printExpiryResults(results []audit.ExpiryResult) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ENV\tPATH\tSTATUS\tMESSAGE")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Env, r.Path, r.Status, r.Message)
	}
	return w.Flush()
}

func init() {
	expireCmd.Flags().String("config", "configs/vaultwatch.yaml", "Path to vaultwatch config")
	expireCmd.Flags().String("path", "", "Secret path to evaluate")
	expireCmd.Flags().String("expiry-file", "", "Path to expiry rules JSON file")
	expireCmd.Flags().String("score-file", "", "Optional pre-scored report file")
	rootCmd.AddCommand(expireCmd)
}
