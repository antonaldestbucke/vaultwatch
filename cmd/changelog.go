package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultwatch/internal/audit"
)

var (
	changelogFile   string
	changelogEnv    string
	changelogAction string
)

var changelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "View or append to the secret change history log",
	RunE: func(cmd *cobra.Command, args []string) error {
		switch changelogAction {
		case "list":
			return runChangelogList()
		case "append":
			return runChangelogAppend(cmd)
		default:
			return fmt.Errorf("unknown action %q: use 'list' or 'append'", changelogAction)
		}
	},
}

func runChangelogList() error {
	store, err := audit.LoadChangelog(changelogFile)
	if err != nil {
		return fmt.Errorf("load changelog: %w", err)
	}
	audit.PrintChangelog(store)
	return nil
}

func runChangelogAppend(cmd *cobra.Command) error {
	if changelogEnv == "" {
		return fmt.Errorf("--env is required for append action")
	}
	cfgPath, _ := cmd.Flags().GetString("config")
	secretPath, _ := cmd.Flags().GetString("path")
	if secretPath == "" {
		return fmt.Errorf("--path is required for append action")
	}
	cfg, err := audit.LoadConfigFromFile(cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	clients, err := audit.ClientsFromCfg(cfg)
	if err != nil {
		return fmt.Errorf("build clients: %w", err)
	}
	reports, err := audit.ComparePathAcrossEnvs(secretPath, clients)
	if err != nil {
		return fmt.Errorf("compare: %w", err)
	}
	scored := audit.ScoreReports(reports)
	store, err := audit.AppendChangelog(changelogFile, scored, changelogEnv)
	if err != nil {
		return fmt.Errorf("append changelog: %w", err)
	}
	fmt.Fprintf(os.Stdout, "Appended %d change(s) to %s\n", len(store.Entries), changelogFile)
	return nil
}

func init() {
	changelogCmd.Flags().StringVar(&changelogFile, "file", "changelog.json", "Path to changelog JSON file")
	changelogCmd.Flags().StringVar(&changelogEnv, "env", "", "Environment label for appended entries")
	changelogCmd.Flags().StringVar(&changelogAction, "action", "list", "Action to perform: list or append")
	changelogCmd.Flags().String("config", "configs/vaultwatch.yaml", "Path to vaultwatch config")
	changelogCmd.Flags().String("path", "", "Secret path to compare (required for append)")
	rootCmd.AddCommand(changelogCmd)
}
