package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var (
	baselineName   string
	baselineDir    string
	baselineAction string
)

var baselineCmd = &cobra.Command{
	Use:   "baseline",
	Short: "Save or diff against a named baseline snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		if secretPath == "" {
			return fmt.Errorf("--path is required")
		}
		if baselineName == "" {
			return fmt.Errorf("--name is required")
		}

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("init clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(clients, secretPath)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}

		switch baselineAction {
		case "save":
			if err := audit.SaveBaseline(baselineDir, baselineName, secretPath, reports); err != nil {
				return fmt.Errorf("save baseline: %w", err)
			}
			fmt.Fprintf(os.Stdout, "Baseline %q saved to %s\n", baselineName, baselineDir)
		case "diff":
			b, err := audit.LoadBaseline(baselineDir, baselineName)
			if err != nil {
				return fmt.Errorf("load baseline: %w", err)
			}
			changed := audit.DiffAgainstBaseline(b, reports)
			if len(changed) == 0 {
				fmt.Println("No changes since baseline.")
			} else {
				fmt.Println("Changed paths since baseline:")
				for _, p := range changed {
					fmt.Printf("  - %s\n", p)
				}
			}
		default:
			return fmt.Errorf("--action must be 'save' or 'diff'")
		}
		return nil
	},
}

func init() {
	baselineCmd.Flags().StringVar(&secretPath, "path", "", "Secret path to audit")
	baselineCmd.Flags().StringVar(&cfgFile, "config", "configs/vaultwatch.yaml", "Config file")
	baselineCmd.Flags().StringVar(&baselineName, "name", "", "Baseline name")
	baselineCmd.Flags().StringVar(&baselineDir, "dir", ".", "Directory to store baselines")
	baselineCmd.Flags().StringVar(&baselineAction, "action", "diff", "Action: save or diff")
	rootCmd.AddCommand(baselineCmd)
}
