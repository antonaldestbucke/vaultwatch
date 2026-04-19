package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var rollupCmd = &cobra.Command{
	Use:   "rollup",
	Short: "Summarize drift by path prefix from a scored report file",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		depth, _ := cmd.Flags().GetInt("depth")
		jsonOut, _ := cmd.Flags().GetBool("json")

		if file == "" {
			return fmt.Errorf("--file is required")
		}

		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		rollup := audit.BuildRollup(reports, depth)

		if jsonOut {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			return enc.Encode(rollup)
		}

		audit.PrintRollup(rollup)
		return nil
	},
}

func init() {
	rollupCmd.Flags().String("file", "", "Path to scored report JSON file")
	rollupCmd.Flags().Int("depth", 2, "Path depth for grouping prefixes")
	rollupCmd.Flags().Bool("json", false, "Output as JSON")
	rootCmd.AddCommand(rollupCmd)
}
