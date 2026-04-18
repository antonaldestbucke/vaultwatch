package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultwatch/internal/audit"
	"github.com/yourorg/vaultwatch/internal/vault"
)

var reportPaths []string

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate a diff report for secret paths across environments",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}
		if len(clients) < 2 {
			return fmt.Errorf("at least two environments required")
		}

		envNames := make([]string, 0, len(clients))
		for name := range clients {
			envNames = append(envNames, name)
		}
		envA, envB := envNames[0], envNames[1]
		clientA, clientB := clients[envA], clients[envB]

		var reports []audit.PathReport
		for _, path := range reportPaths {
			result, err := audit.ComparePathAcrossEnvs(cmd.Context(), clientA, clientB, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warn: compare %s: %v\n", path, err)
				continue
			}
			reports = append(reports, audit.PathReport{
				Path: path,
				EnvA: envA,
				EnvB: envB,
				Diff: audit.DiffResult{
					OnlyInA: result.OnlyInA,
					OnlyInB: result.OnlyInB,
				},
			})
		}

		audit.PrintTextReport(os.Stdout, reports)
		fmt.Println(audit.Summary(reports))
		return nil
	},
}

func init() {
	reportCmd.Flags().StringSliceVarP(&reportPaths, "path", "p", nil, "Secret paths to compare (repeatable)")
	_ = reportCmd.MarkFlagRequired("path")
	rootCmd.AddCommand(reportCmd)
}
