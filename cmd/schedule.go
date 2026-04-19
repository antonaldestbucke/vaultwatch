package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"github.com/vaultwatch/internal/audit"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Manage audit schedules",
}

var scheduleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all schedule entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		store, err := audit.LoadSchedule(file)
		if err != nil {
			return fmt.Errorf("load schedule: %w", err)
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tPATH\tINTERVAL\tENABLED\tNEXT DUE")
		for _, e := range store.Entries {
			due, err := audit.NextDue(e)
			dueSt := "now"
			if err == nil && due > 0 {
				dueSt = due.Round(time.Second).String()
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\t%s\n", e.Name, e.Path, e.Interval, e.Enabled, dueSt)
		}
		return w.Flush()
	},
}

var scheduleAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a schedule entry",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")
		path, _ := cmd.Flags().GetString("path")
		interval, _ := cmd.Flags().GetString("interval")
		if name == "" || path == "" || interval == "" {
			return fmt.Errorf("--name, --path, and --interval are required")
		}
		store, _ := audit.LoadSchedule(file)
		store.Entries = append(store.Entries, audit.ScheduleEntry{
			Name: name, Path: path, Interval: interval, Enabled: true,
		})
		if err := audit.SaveSchedule(file, store); err != nil {
			return err
		}
		fmt.Printf("Added schedule %q\n", name)
		return nil
	},
}

func init() {
	for _, sub := range []*cobra.Command{scheduleListCmd, scheduleAddCmd} {
		sub.Flags().String("file", "schedule.json", "Path to schedule file")
	}
	scheduleAddCmd.Flags().String("name", "", "Schedule name")
	scheduleAddCmd.Flags().String("path", "", "Vault secret path")
	scheduleAddCmd.Flags().String("interval", "", "Poll interval (e.g. 1h, 30m)")
	scheduleCmd.AddCommand(scheduleListCmd, scheduleAddCmd)
	rootCmd.AddCommand(scheduleCmd)
}
