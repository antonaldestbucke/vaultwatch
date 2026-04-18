package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export a diff report to JSON or CSV",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile, _ := cmd.Flags().GetString("config")
		path, _ := cmd.Flags().GetString("path")
		format, _ := cmd.Flags().GetString("format")
		outFile, _ := cmd.Flags().GetString("output")

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(cmd.Context(), clients, path)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}

		w := os.Stdout
		if outFile != "" {
			f, err := os.Create(outFile)
			if err != nil {
				return fmt.Errorf("open output file: %w", err)
			}
			defer f.Close()
			w = f
		}

		return audit.ExportReport(w, reports, audit.ExportFormat(format))
	},
}

func init() {
	exportCmd.Flags().String("config", "vaultwatch.yaml", "Path to config file")
	exportCmd.Flags().String("path", "", "Secret path to compare (required)")
	exportCmd.Flags().String("format", "json", "Export format: json or csv")
	exportCmd.Flags().String("output", "", "Output file path (default: stdout)")
	_ = exportCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(exportCmd)
}
