package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var (
	snapshotOutput string
	snapshotEnv    string
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot <path>",
	Short: "Capture a snapshot of secret keys at a given path",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secretPath := args[0]

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("building clients: %w", err)
		}

		client, ok := clients[snapshotEnv]
		if !ok {
			return fmt.Errorf("environment %q not found in config", snapshotEnv)
		}

		keys, err := client.ListSecrets(secretPath)
		if err != nil {
			return fmt.Errorf("listing secrets at %q: %w", secretPath, err)
		}

		snap := audit.Snapshot{
			Path:       secretPath,
			Env:        snapshotEnv,
			Keys:       keys,
			CapturedAt: time.Now().UTC(),
		}

		if err := audit.SaveSnapshot(snap, snapshotOutput); err != nil {
			return fmt.Errorf("saving snapshot: %w", err)
		}

		fmt.Printf("Snapshot saved to %s (%d keys)\n", snapshotOutput, len(keys))
		return nil
	},
}

func init() {
	snapshotCmd.Flags().StringVarP(&snapshotOutput, "output", "o", "snapshot.json", "Output file for the snapshot")
	snapshotCmd.Flags().StringVarP(&snapshotEnv, "env", "e", "", "Environment to snapshot (required)")
	_ = snapshotCmd.MarkFlagRequired("env")
	rootCmd.AddCommand(snapshotCmd)
}
