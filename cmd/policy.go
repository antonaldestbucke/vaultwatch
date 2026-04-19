package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var policyCmd = &cobra.Command{
	Use:   "policy",
	Short: "Evaluate secret paths against policy rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile, _ := cmd.Flags().GetString("config")
		policyFile, _ := cmd.Flags().GetString("policy")
		path, _ := cmd.Flags().GetString("path")

		if path == "" {
			return fmt.Errorf("--path is required")
		}
		if policyFile == "" {
			return fmt.Errorf("--policy is required")
		}

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(path, clients)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}

		store, err := audit.LoadPolicy(policyFile)
		if err != nil {
			return fmt.Errorf("load policy: %w", err)
		}

		violations := audit.EvaluatePolicy(reports, store)
		if len(violations) == 0 {
			fmt.Println("✓ No policy violations found.")
			return nil
		}

		fmt.Fprintf(os.Stderr, "✗ %d policy violation(s) found:\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "  [%s] %s\n", v.Path, v.Message)
		}
		os.Exit(1)
		return nil
	},
}

func init() {
	policyCmd.Flags().String("config", "configs/vaultwatch.yaml", "Path to vaultwatch config")
	policyCmd.Flags().String("policy", "", "Path to policy JSON file")
	policyCmd.Flags().String("path", "", "Secret path to evaluate")
	rootCmd.AddCommand(policyCmd)
}
