package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultwatch/internal/audit"
)

var accessFile string

var accessCmd = &cobra.Command{
	Use:   "access",
	Short: "Manage access rules for secret paths",
}

var accessAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an access rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, _ := cmd.Flags().GetString("path")
		owner, _ := cmd.Flags().GetString("owner")
		team, _ := cmd.Flags().GetString("team")
		expiry, _ := cmd.Flags().GetString("expires")
		if path == "" || owner == "" {
			return fmt.Errorf("--path and --owner are required")
		}
		store, err := audit.LoadAccess(accessFile)
		if err != nil {
			return err
		}
		rule := audit.AccessRule{Path: path, Owner: owner, Team: team}
		if expiry != "" {
			t, err := time.Parse(time.RFC3339, expiry)
			if err != nil {
				return fmt.Errorf("invalid --expires format, use RFC3339: %w", err)
			}
			rule.Expires = t
		}
		store.Rules = append(store.Rules, rule)
		if err := audit.SaveAccess(accessFile, store); err != nil {
			return err
		}
		fmt.Fprintf(os.Stdout, "access rule added: %s -> %s (%s)\n", path, owner, team)
		return nil
	},
}

var accessListCmd = &cobra.Command{
	Use:   "list",
	Short: "List access rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := audit.LoadAccess(accessFile)
		if err != nil {
			return err
		}
		if len(store.Rules) == 0 {
			fmt.Println("no access rules defined")
			return nil
		}
		for _, r := range store.Rules {
			expiry := ""
			if !r.Expires.IsZero() {
				expiry = fmt.Sprintf(" (expires %s)", r.Expires.Format(time.RFC3339))
			}
			fmt.Printf("  %-40s owner=%-15s team=%s%s\n", r.Path, r.Owner, r.Team, expiry)
		}
		return nil
	},
}

func init() {
	accessCmd.PersistentFlags().StringVar(&accessFile, "file", "access.json", "path to access rules file")
	accessAddCmd.Flags().String("path", "", "secret path prefix")
	accessAddCmd.Flags().String("owner", "", "owner name")
	accessAddCmd.Flags().String("team", "", "team name")
	accessAddCmd.Flags().String("expires", "", "expiry time (RFC3339)")
	accessCmd.AddCommand(accessAddCmd, accessListCmd)
	rootCmd.AddCommand(accessCmd)
}
