package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Manage locked secret paths that should not drift",
}

var lockAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Lock a secret path",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		path, _ := cmd.Flags().GetString("path")
		by, _ := cmd.Flags().GetString("by")
		reason, _ := cmd.Flags().GetString("reason")
		if path == "" || by == "" {
			return fmt.Errorf("--path and --by are required")
		}
		store, err := audit.LoadLocks(file)
		if err != nil {
			return err
		}
		store.Locks = append(store.Locks, audit.LockEntry{
			Path: path, LockedBy: by, Reason: reason, LockedAt: time.Now(),
		})
		if err := audit.SaveLocks(file, store); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "Locked: %s\n", path)
		return nil
	},
}

var lockListCmd = &cobra.Command{
	Use:   "list",
	Short: "List locked paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		store, err := audit.LoadLocks(file)
		if err != nil {
			return err
		}
		if len(store.Locks) == 0 {
			fmt.Println("No locked paths.")
			return nil
		}
		for _, l := range store.Locks {
			fmt.Printf("%-40s  by=%-15s  reason=%s\n", l.Path, l.LockedBy, l.Reason)
		}
		return nil
	},
}

func init() {
	for _, sub := range []*cobra.Command{lockAddCmd, lockListCmd} {
		sub.Flags().String("file", "locks.json", "Path to lock store file")
	}
	lockAddCmd.Flags().String("path", "", "Secret path to lock")
	lockAddCmd.Flags().String("by", "", "Who is locking the path")
	lockAddCmd.Flags().String("reason", "", "Reason for locking")
	lockCmd.AddCommand(lockAddCmd, lockListCmd)
	rootCmd.AddCommand(lockCmd)
}
