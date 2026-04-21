package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var heatmapCmd = &cobra.Command{
	Use:   "heatmap",
	Short: "Show a drift frequency heatmap across all audited paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile, _ := cmd.Flags().GetString("config")
		path, _ := cmd.Flags().GetString("path")
		topN, _ := cmd.Flags().GetInt("top")

		if path == "" {
			return fmt.Errorf("--path is required")
		}

		cfg, err := vault.LoadConfig(cfgFile)
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

		scored := audit.ScoreReports(reports)
		heatmap := audit.BuildHeatmap(scored)

		if topN > 0 && topN < len(heatmap.Entries) {
			heatmap.Entries = heatmap.Entries[:topN]
		}

		audit.PrintHeatmap(heatmap)
		return nil
	},
}

func init() {
	heatmapCmd.Flags().String("config", "configs/vaultwatch.yaml", "Path to config file")
	heatmapCmd.Flags().String("path", "", "Secret path to audit")
	heatmapCmd.Flags().Int("top", 0, "Show only top N drifted paths (0 = all)")
	if err := heatmapCmd.MarkFlagRequired("path"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	rootCmd.AddCommand(heatmapCmd)
}
