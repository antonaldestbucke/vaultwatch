package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Add or view annotations on secret paths",
}

var annotateAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an annotation to a secret path",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		path, _ := cmd.Flags().GetString("path")
		note, _ := cmd.Flags().GetString("note")
		author, _ := cmd.Flags().GetString("author")

		var store audit.AnnotationStore
		loaded, err := audit.LoadAnnotations(file)
		if err != nil && !os.IsNotExist(err) {
			store = make(audit.AnnotationStore)
		} else if err == nil {
			store = loaded
		} else {
			store = make(audit.AnnotationStore)
		}

		store[path] = audit.Annotation{
			Path:      path,
			Note:      note,
			Author:    author,
			CreatedAt: time.Now(),
		}
		if err := audit.SaveAnnotations(file, store); err != nil {
			return fmt.Errorf("save annotations: %w", err)
		}
		fmt.Printf("Annotation saved for path: %s\n", path)
		return nil
	},
}

var annotateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all annotations",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		store, err := audit.LoadAnnotations(file)
		if err != nil {
			return fmt.Errorf("load annotations: %w", err)
		}
		for _, ann := range store {
			fmt.Printf("[%s] %s — %s (by %s)\n", ann.CreatedAt.Format("2006-01-02"), ann.Path, ann.Note, ann.Author)
		}
		return nil
	},
}

func init() {
	annotateAddCmd.Flags().String("file", "annotations.json", "Annotation store file")
	annotateAddCmd.Flags().String("path", "", "Secret path to annotate")
	annotateAddCmd.Flags().String("note", "", "Note text")
	annotateAddCmd.Flags().String("author", "", "Author name")
	annotateAddCmd.MarkFlagRequired("path")
	annotateAddCmd.MarkFlagRequired("note")

	annotateListCmd.Flags().String("file", "annotations.json", "Annotation store file")

	annotateCmd.AddCommand(annotateAddCmd, annotateListCmd)
	rootCmd.AddCommand(annotateCmd)
}
