package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultwatch/internal/audit"
	"github.com/yourusername/vaultwatch/internal/vault"
)

var (
	quotaConfigFile string
	quotaPath       string
	quotaFile       string
	quotaJSON       bool
)

var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "Evaluate drift quota rules against secret paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		if quotaPath == "" {
			return fmt.Errorf("--path is required")
		}
		if quotaFile == "" {
			return fmt.Errorf("--quota-file is required")
		}

		cfg, err := vault.LoadConfig(quotaConfigFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(cmd.Context(), clients, quotaPath)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}

		scored := audit.ScoreReports(reports)

		store, err := audit.LoadQuota(quotaFile)
		if err != nil {
			return fmt.Errorf("load quota file: %w", err)
		}

		violations := audit.EvaluateQuota(scored, store)

		if quotaJSON {
			return json.NewEncoder(os.Stdout).Encode(violations)
		}

		if len(violations) == 0 {
			fmt.Println("No quota violations detected.")
			return nil
		}

		fmt.Printf("%-40s %10s %10s %s\n", "PATH", "DRIFTED", "MAX", "RULE")
		for _, v := range violations {
			fmt.Printf("%-40s %10d %10d %s\n", v.Path, v.Drifted, v.MaxAllowed, v.Rule)
		}
		return nil
	},
}

func init() {
	quotaCmd.Flags().StringVar(&quotaConfigFile, "config", "configs/vaultwatch.yaml", "Path to vaultwatch config")
	quotaCmd.Flags().StringVar(&quotaPath, "path", "", "Secret path to audit")
	quotaCmd.Flags().StringVar(&quotaFile, "quota-file", "", "Path to quota rules JSON file")
	quotaCmd.Flags().BoolVar(&quotaJSON, "json", false, "Output violations as JSON")
	rootCmd.AddCommand(quotaCmd)
}
