package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var trendFile string

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Analyze drift score trend from scored snapshot files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if trendFile == "" {
			return fmt.Errorf("--file is required")
		}
		data, err := os.ReadFile(trendFile)
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}
		var snapshots []audit.ScoredReport
		if err := json.Unmarshal(data, &snapshots); err != nil {
			return fmt.Errorf("parsing snapshots: %w", err)
		}
		tr := audit.BuildTrend(snapshots)
		fmt.Printf("Trend points : %d\n", len(tr.Points))
		fmt.Printf("Average score: %.2f\n", tr.AverageScore())
		if w, ok := tr.WorstPoint(); ok {
			fmt.Printf("Worst point  : %s (score=%.2f drifted=%d/%d)\n",
				w.Timestamp.Format("2006-01-02T15:04:05Z"), w.Score, w.Drifted, w.Total)
		}
		return nil
	},
}

func init() {
	trendCmd.Flags().StringVar(&trendFile, "file", "", "Path to JSON file containing []ScoredReport")
	rootCmd.AddCommand(trendCmd)
}
