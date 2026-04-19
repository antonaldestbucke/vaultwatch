package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var digestCmd = &cobra.Command{
	Use:   "digest",
	Short: "Compute a deterministic hash digest over secret path comparisons",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgPath, _ := cmd.Flags().GetString("config")
		secretPath, _ := cmd.Flags().GetString("path")
		outputJSON, _ := cmd.Flags().GetBool("json")

		if secretPath == "" {
			return fmt.Errorf("--path is required")
		}

		cfg, err := vault.LoadConfig(cfgPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}

		reports, err := audit.ComparePathAcrossEnvs(clients, secretPath)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}

		entry, err := audit.BuildDigest(reports)
		if err != nil {
			return fmt.Errorf("digest: %w", err)
		}

		if outputJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(entry)
		}

		fmt.Printf("Digest : %s\n", entry.Hash)
		fmt.Printf("Paths  : %d\n", entry.PathCount)
		fmt.Printf("Drifted: %d\n", entry.DriftCount)
		fmt.Printf("Time   : %s\n", entry.Timestamp.Format("2006-01-02T15:04:05Z"))
		return nil
	},
}

func init() {
	digestCmd.Flags().String("config", "configs/vaultwatch.yaml", "Path to config file")
	digestCmd.Flags().String("path", "", "Secret path to compare across environments")
	digestCmd.Flags().Bool("json", false, "Output digest as JSON")
	rootCmd.AddCommand(digestCmd)
}
