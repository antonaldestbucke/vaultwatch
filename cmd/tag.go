package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags for secret paths",
}

var tagAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add tags to a secret path",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		path, _ := cmd.Flags().GetString("path")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		if path == "" || len(tags) == 0 {
			return fmt.Errorf("--path and --tags are required")
		}
		store, err := audit.LoadTags(file)
		if err != nil {
			if os.IsNotExist(err) {
				store = audit.TagStore{}
			} else {
				return fmt.Errorf("load tags: %w", err)
			}
		}
		existing := store[path]
		store[path] = append(existing, tags...)
		if err := audit.SaveTags(file, store); err != nil {
			return fmt.Errorf("save tags: %w", err)
		}
		fmt.Printf("Tagged %s with: %s\n", path, strings.Join(tags, ", "))
		return nil
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tagged paths",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		store, err := audit.LoadTags(file)
		if err != nil {
			return fmt.Errorf("load tags: %w", err)
		}
		for path, tags := range store {
			fmt.Printf("%s: %s\n", path, strings.Join(tags, ", "))
		}
		return nil
	},
}

func init() {
	for _, sub := range []*cobra.Command{tagAddCmd, tagListCmd} {
		sub.Flags().String("file", "tags.json", "Path to tag store JSON file")
	}
	tagAddCmd.Flags().String("path", "", "Secret path to tag")
	tagAddCmd.Flags().StringSlice("tags", nil, "Comma-separated tags to apply")
	tagCmd.AddCommand(tagAddCmd, tagListCmd)
	rootCmd.AddCommand(tagCmd)
}
