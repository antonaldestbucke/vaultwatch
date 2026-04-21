package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultwatch/internal/audit"
)

var (
	clusterScoreFile string
	clusterThreshold float64
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Group secret paths into clusters based on drift score proximity",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(clusterScoreFile)
		if err != nil {
			return fmt.Errorf("reading score file: %w", err)
		}

		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("parsing score file: %w", err)
		}

		clusters := audit.ClusterReports(reports, clusterThreshold)
		if len(clusters) == 0 {
			fmt.Println("No clusters found.")
			return nil
		}

		for i, c := range clusters {
			fmt.Printf("Cluster %d (centroid: %s, avg_score: %.2f)\n", i+1, c.Centroid, c.AvgScore)
			for _, p := range c.Paths {
				fmt.Printf("  - %s\n", p)
			}
		}
		return nil
	},
}

func init() {
	clusterCmd.Flags().StringVar(&clusterScoreFile, "scores", "", "Path to scored reports JSON file (required)")
	clusterCmd.Flags().Float64Var(&clusterThreshold, "threshold", 0.15, "Max score distance to group paths into the same cluster")
	_ = clusterCmd.MarkFlagRequired("scores")
	rootCmd.AddCommand(clusterCmd)
}
