package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var (
	alertMinScore  float64
	alertPrefix    string
	alertOnlyDrift bool
	alertFile      string
)

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Evaluate scored reports against alert rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		if alertFile == "" {
			return fmt.Errorf("--file is required")
		}

		data, err := os.ReadFile(alertFile)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}

		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("parsing scored reports: %w", err)
		}

		rule := audit.AlertRule{
			MinRiskScore: alertMinScore,
			PathPrefix:   alertPrefix,
			OnlyDrifted:  alertOnlyDrift,
		}

		alerts := audit.EvaluateAlerts(reports, rule)
		audit.PrintAlertsToStdout(alerts)

		if len(alerts) > 0 {
			os.Exit(2)
		}
		return nil
	},
}

func init() {
	alertCmd.Flags().StringVar(&alertFile, "file", "", "Path to scored reports JSON file")
	alertCmd.Flags().Float64Var(&alertMinScore, "min-score", 0.0, "Minimum risk score to trigger alert")
	alertCmd.Flags().StringVar(&alertPrefix, "path-prefix", "", "Only alert on paths with this prefix")
	alertCmd.Flags().BoolVar(&alertOnlyDrift, "only-drifted", false, "Only alert on drifted paths")
	rootCmd.AddCommand(alertCmd)
}
