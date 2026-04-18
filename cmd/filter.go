package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultwatch/internal/audit"
	"github.com/yourusername/vaultwatch/internal/vault"
)

var (
	filterOnlyDiffs  bool
	filterPathPrefix string
)

var filterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Filter and display secret diff results by path prefix or diff status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(configFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("creating clients: %w", err)
		}
		if secretPath == "" {
			return fmt.Errorf("--path is required")
		}
		reports, err := audit.ComparePathAcrossEnvs(clients, secretPath)
		if err != nil {
			return fmt.Errorf("comparing paths: %w", err)
		}
		filtered := audit.FilterReports(reports, audit.FilterOptions{
			OnlyDiffs:  filterOnlyDiffs,
			PathPrefix: filterPathPrefix,
		})
		audit.PrintTextReport(os.Stdout, filtered)
		return nil
	},
}

func init() {
	filterCmd.Flags().StringVar(&configFile, "config", "vaultwatch.yaml", "Path to config file")
	filterCmd.Flags().StringVar(&secretPath, "path", "", "Vault secret path to compare")
	filterCmd.Flags().BoolVar(&filterOnlyDiffs, "only-diffs", false, "Show only paths with differences")
	filterCmd.Flags().StringVar(&filterPathPrefix, "prefix", "", "Filter results by path prefix")
	rootCmd.AddCommand(filterCmd)
}
