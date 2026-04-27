package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var anomalyCmd = &cobra.Command{
	Use:   "anomaly",
	Short: "Detect anomalous drift scores across secret paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		scoresFile, _ := cmd.Flags().GetString("scores")
		zThreshold, _ := cmd.Flags().GetFloat64("z-threshold")
		jsonOut, _ := cmd.Flags().GetBool("json")
		onlyAnomalies, _ := cmd.Flags().GetBool("only-anomalies")

		data, err := os.ReadFile(scoresFile)
		if err != nil {
			return fmt.Errorf("failed to read scores file: %w", err)
		}

		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("invalid scores JSON: %w", err)
		}

		results := audit.BuildAnomalies(reports, zThreshold)

		if onlyAnomalies {
			var filtered []audit.AnomalyResult
			for _, r := range results {
				if r.IsAnomaly {
					filtered = append(filtered, r)
				}
			}
			results = filtered
		}

		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(results)
		}

		if len(results) == 0 {
			fmt.Println("No anomalies detected.")
			return nil
		}

		fmt.Printf("%-40s %-8s %-8s %-8s %-10s\n", "PATH", "SCORE", "ZSCORE", "ANOMALY", "SEVERITY")
		for _, r := range results {
			anomaly := "no"
			if r.IsAnomaly {
				anomaly = "YES"
			}
			fmt.Printf("%-40s %-8.2f %-8.2f %-8s %-10s\n",
				r.Path, r.Score, r.ZScore, anomaly, r.Severity)
		}
		return nil
	},
}

func init() {
	anomalyCmd.Flags().String("scores", "", "Path to scored reports JSON file (required)")
	anomalyCmd.Flags().Float64("z-threshold", 2.0, "Z-score threshold for anomaly detection")
	anomalyCmd.Flags().Bool("json", false, "Output results as JSON")
	anomalyCmd.Flags().Bool("only-anomalies", false, "Only show paths flagged as anomalies")
	_ = anomalyCmd.MarkFlagRequired("scores")
	rootCmd.AddCommand(anomalyCmd)
}
