package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var (
	lifecycleFile  string
	lifecyclePath  string
	lifecycleStage string
	lifecycleNote  string
)

var lifecycleCmd = &cobra.Command{
	Use:   "lifecycle",
	Short: "Manage secret path lifecycle stages (active, deprecated, retired)",
}

var lifecycleSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the lifecycle stage for a secret path",
	RunE: func(cmd *cobra.Command, args []string) error {
		if lifecyclePath == "" {
			return fmt.Errorf("--path is required")
		}
		if lifecycleStage == "" {
			return fmt.Errorf("--stage is required")
		}
		stage := audit.LifecycleStage(lifecycleStage)
		if stage != audit.StageActive && stage != audit.StageDeprecated && stage != audit.StageRetired {
			return fmt.Errorf("invalid stage %q: must be active, deprecated, or retired", lifecycleStage)
		}
		store, err := audit.LoadLifecycle(lifecycleFile)
		if err != nil {
			return fmt.Errorf("load lifecycle: %w", err)
		}
		audit.SetLifecycleStage(&store, lifecyclePath, stage, lifecycleNote)
		if err := audit.SaveLifecycle(lifecycleFile, store); err != nil {
			return fmt.Errorf("save lifecycle: %w", err)
		}
		fmt.Printf("Set %s → %s\n", lifecyclePath, stage)
		return nil
	},
}

var lifecycleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all lifecycle entries",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := audit.LoadLifecycle(lifecycleFile)
		if err != nil {
			return fmt.Errorf("load lifecycle: %w", err)
		}
		if len(store.Entries) == 0 {
			fmt.Println("No lifecycle entries found.")
			return nil
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "PATH\tSTAGE\tUPDATED\tNOTE")
		for _, e := range store.Entries {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.Path, e.Stage, e.UpdatedAt.Format(time.RFC3339), e.Note)
		}
		return w.Flush()
	},
}

func init() {
	lifecycleSetCmd.Flags().StringVar(&lifecyclePath, "path", "", "Secret path")
	lifecycleSetCmd.Flags().StringVar(&lifecycleStage, "stage", "", "Lifecycle stage: active, deprecated, retired")
	lifecycleSetCmd.Flags().StringVar(&lifecycleNote, "note", "", "Optional note")
	lifecycleSetCmd.Flags().StringVar(&lifecycleFile, "lifecycle-file", "lifecycle.json", "Path to lifecycle store")
	lifecycleListCmd.Flags().StringVar(&lifecycleFile, "lifecycle-file", "lifecycle.json", "Path to lifecycle store")
	lifecycleCmd.AddCommand(lifecycleSetCmd, lifecycleListCmd)
	rootCmd.AddCommand(lifecycleCmd)
}
