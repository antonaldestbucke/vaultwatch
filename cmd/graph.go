package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var (
	graphScoresFile string
	graphDepsFile   string
	graphOutputJSON bool
)

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Build and display a secret path dependency graph",
	RunE: func(cmd *cobra.Command, args []string) error {
		if graphScoresFile == "" {
			return fmt.Errorf("--scores is required")
		}

		data, err := os.ReadFile(graphScoresFile)
		if err != nil {
			return fmt.Errorf("reading scores file: %w", err)
		}
		var reports []audit.ScoredReport
		if err := json.Unmarshal(data, &reports); err != nil {
			return fmt.Errorf("parsing scores file: %w", err)
		}

		deps := map[string][]string{}
		if graphDepsFile != "" {
			depData, err := os.ReadFile(graphDepsFile)
			if err != nil {
				return fmt.Errorf("reading deps file: %w", err)
			}
			if err := json.Unmarshal(depData, &deps); err != nil {
				return fmt.Errorf("parsing deps file: %w", err)
			}
		}

		result := audit.BuildGraph(reports, deps)

		if graphOutputJSON {
			out, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(out))
			return nil
		}

		fmt.Print(audit.PrintGraph(result))
		return nil
	},
}

func init() {
	graphCmd.Flags().StringVar(&graphScoresFile, "scores", "", "Path to scored reports JSON file")
	graphCmd.Flags().StringVar(&graphDepsFile, "deps", "", "Path to dependency map JSON file (optional)")
	graphCmd.Flags().BoolVar(&graphOutputJSON, "json", false, "Output graph as JSON")
	rootCmd.AddCommand(graphCmd)
}
