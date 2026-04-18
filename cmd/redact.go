package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/vaultwatch/internal/audit"
	"github.com/user/vaultwatch/internal/vault"
)

var redactPatterns []string

var redactCmd = &cobra.Command{
	Use:   "redact",
	Short: "Run a report with sensitive keys redacted",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			return fmt.Errorf("--path is required")
		}

		cfgFile, _ := cmd.Flags().GetString("config")
		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("failed to create clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(cmd.Context(), path, clients)
		if err != nil {
			return fmt.Errorf("comparison failed: %w", err)
		}

		patterns := normalizePatterns(redactPatterns)
		if len(patterns) == 0 {
			return fmt.Errorf("at least one redact pattern is required")
		}

		opts := audit.RedactOptions{KeyPatterns: patterns}
		redacted := audit.RedactReports(reports, opts)

		audit.PrintTextReport(redacted)
		summary := audit.Summary(redacted)
		fmt.Printf("\nSummary: %d paths checked, %d with diffs\n", summary.Total, summary.WithDiffs)
		return nil
	},
}

// normalizePatterns trims whitespace from each pattern and drops empty entries.
func normalizePatterns(patterns []string) []string {
	result := make([]string, 0, len(patterns))
	for _, p := range patterns {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

func init() {
	redactCmd.Flags().String("path", "", "Vault path to audit (required)")
	redactCmd.Flags().String("config", "vaultwatch.yaml", "Path to config file")
	redactCmd.Flags().StringSliceVar(&redactPatterns, "redact", []string{"password", "token", "secret", "key"},
		"Comma-separated key name patterns to redact (case-insensitive)")
	_ = redactCmd.Flags().SetAnnotation("path", cobra.BashCompOneRequiredFlag, []string{"true"})

	rootCmd.AddCommand(redactCmd)
}
