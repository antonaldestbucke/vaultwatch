package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var (
	coveragePath      string
	coverageConfig    string
	coverageMinPct    float64
	coverageShowFull  bool
)

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "Show environment coverage for each secret path",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(coverageConfig)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("init clients: %w", err)
		}
		reports, err := audit.ComparePathAcrossEnvs(clients, coveragePath)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}
		cov := audit.BuildCoverage(reports)
		if len(cov) == 0 {
			fmt.Println("No paths found.")
			return nil
		}
		for _, c := range cov {
			if c.CoveragePct >= coverageMinPct || coverageShowFull {
				fmt.Printf("%-40s  %.0f%%  present=[%v]  missing=[%v]\n",
					c.Path, c.CoveragePct,
					joinStrings(c.PresentIn),
					joinStrings(c.MissingFrom))
			}
		}
		return nil
	},
}

func joinStrings(ss []string) string {
	out := ""
	for i, s := range ss {
		if i > 0 {
			out += ","
		}
		out += s
	}
	return out
}

func init() {
	coverageCmd.Flags().StringVar(&coveragePath, "path", "", "Secret path prefix to audit (required)")
	coverageCmd.Flags().StringVar(&coverageConfig, "config", "configs/vaultwatch.yaml", "Path to config file")
	coverageCmd.Flags().Float64Var(&coverageMinPct, "min-pct", 0, "Only show paths below this coverage percentage")
	coverageCmd.Flags().BoolVar(&coverageShowFull, "all", false, "Show all paths regardless of coverage")
	_ = coverageCmd.MarkFlagRequired("path")
	if err := coverageCmd.Flags().MarkHidden("min-pct"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	rootCmd.AddCommand(coverageCmd)
}
