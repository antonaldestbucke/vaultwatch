package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var remediateCmd = &cobra.Command{
	Use:   "remediate",
	Short: "Generate remediation suggestions from scored drift reports",
	RunE: func(cmd *cobra.Command, args []string) error {
		input, _ := cmd.Flags().GetString("input")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if input == "" {
			return fmt.Errorf("--input is required")
		}

		data, err := os.ReadFile(input)
		if err != nil {
			return fmt.Errorf("failed to read input file: %w", err)
		}

		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("invalid JSON in input file: %w", err)
		}

		plan := audit.BuildRemediationPlan(reports)

		if jsonOut {
			out, err := json.MarshalIndent(plan, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal plan: %w", err)
			}
			fmt.Println(string(out))
			return nil
		}

		audit.PrintRemediationPlan(plan)
		return nil
	},
}

func init() {
	remediateCmd.Flags().String("input", "", "Path to scored reports JSON file")
	remediateCmd.Flags().Bool("json", false, "Output remediation plan as JSON")
	rootCmd.AddCommand(remediateCmd)
}
