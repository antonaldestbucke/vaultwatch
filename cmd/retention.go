package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"vaultwatch/internal/audit"
)

var retentionCmd = &cobra.Command{
	Use:   "retention",
	Short: "Evaluate secret report retention against age-based rules",
}

var retentionEvalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate scored reports against a retention policy file",
	RunE: func(cmd *cobra.Command, args []string) error {
		scoresFile, _ := cmd.Flags().GetString("scores")
		policyFile, _ := cmd.Flags().GetString("policy")
		outputJSON, _ := cmd.Flags().GetBool("json")

		if scoresFile == "" {
			return fmt.Errorf("--scores is required")
		}
		if policyFile == "" {
			return fmt.Errorf("--policy is required")
		}

		data, err := os.ReadFile(scoresFile)
		if err != nil {
			return fmt.Errorf("read scores: %w", err)
		}
		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("parse scores: %w", err)
		}

		store, err := audit.LoadRetention(policyFile)
		if err != nil {
			return fmt.Errorf("load retention policy: %w", err)
		}

		results := audit.EvaluateRetention(reports, store, time.Now())

		if outputJSON {
			out, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		pruned := 0
		for _, r := range results {
			status := "RETAIN"
			if r.Pruned {
				status = "PRUNE "
				pruned++
			}
			fmt.Printf("[%s] %s (%s/%s)", status, r.Path, r.EnvA, r.EnvB)
			if r.Reason != "" {
				fmt.Printf(" — %s", r.Reason)
			}
			fmt.Println()
		}
		fmt.Printf("\n%d of %d reports eligible for pruning.\n", pruned, len(results))
		return nil
	},
}

func init() {
	retentionEvalCmd.Flags().String("scores", "", "Path to scored reports JSON file")
	retentionEvalCmd.Flags().String("policy", "", "Path to retention policy JSON file")
	retentionEvalCmd.Flags().Bool("json", false, "Output results as JSON")
	retentionCmd.AddCommand(retentionEvalCmd)
	rootCmd.AddCommand(retentionCmd)
}
