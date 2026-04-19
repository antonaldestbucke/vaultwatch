package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
	"vaultwatch/internal/vault"
)

var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "Record or replay audit snapshots over time",
}

var replayRecordCmd = &cobra.Command{
	Use:   "record",
	Short: "Record current state into replay history",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfgFile, _ := cmd.Flags().GetString("config")
		path, _ := cmd.Flags().GetString("path")
		label, _ := cmd.Flags().GetString("label")
		replayFile, _ := cmd.Flags().GetString("replay-file")

		cfg, err := vault.LoadConfig(cfgFile)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		clients, err := vault.ClientsFromConfig(cfg)
		if err != nil {
			return fmt.Errorf("build clients: %w", err)
		}
		reports, err := audit.ComparePathAcrossEnvs(clients, path)
		if err != nil {
			return fmt.Errorf("compare: %w", err)
		}
		var store audit.ReplayStore
		if existing, err := audit.LoadReplay(replayFile); err == nil {
			store = existing
		}
		audit.AddReplayEntry(&store, label, reports)
		if err := audit.SaveReplay(replayFile, store); err != nil {
			return fmt.Errorf("save replay: %w", err)
		}
		fmt.Fprintf(os.Stdout, "Recorded entry '%s' to %s\n", label, replayFile)
		return nil
	},
}

var replayShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show replay state at a given timestamp",
	RunE: func(cmd *cobra.Command, args []string) error {
		replayFile, _ := cmd.Flags().GetString("replay-file")
		atStr, _ := cmd.Flags().GetString("at")

		at, err := time.Parse(time.RFC3339, atStr)
		if err != nil {
			return fmt.Errorf("invalid --at timestamp (use RFC3339): %w", err)
		}
		store, err := audit.LoadReplay(replayFile)
		if err != nil {
			return fmt.Errorf("load replay: %w", err)
		}
		entry, ok := audit.ReplayAt(store, at)
		if !ok {
			return fmt.Errorf("no replay entry found at or before %s", atStr)
		}
		fmt.Fprintf(os.Stdout, "Entry: %s (label: %s)\n", entry.Timestamp.Format(time.RFC3339), entry.Label)
		audit.PrintTextReport(os.Stdout, entry.Reports)
		return nil
	},
}

func init() {
	replayRecordCmd.Flags().String("config", "configs/vaultwatch.yaml", "Config file")
	replayRecordCmd.Flags().String("path", "", "Secret path to audit")
	replayRecordCmd.Flags().String("label", "snapshot", "Label for this entry")
	replayRecordCmd.Flags().String("replay-file", "replay.json", "Replay history file")
	replayRecordCmd.MarkFlagRequired("path")

	replayShowCmd.Flags().String("replay-file", "replay.json", "Replay history file")
	replayShowCmd.Flags().String("at", "", "Timestamp in RFC3339 format")
	replayShowCmd.MarkFlagRequired("at")

	replayCmd.AddCommand(replayRecordCmd, replayShowCmd)
	rootCmd.AddCommand(replayCmd)
}
