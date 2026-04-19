package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"vaultwatch/internal/audit"
)

func init() {
	var file string

	ownershipCmd := &cobra.Command{
		Use:   "ownership",
		Short: "Manage path ownership mappings",
	}

	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add or update an ownership entry",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("path")
			owner, _ := cmd.Flags().GetString("owner")
			team, _ := cmd.Flags().GetString("team")
			contact, _ := cmd.Flags().GetString("contact")
			if path == "" || owner == "" {
				return fmt.Errorf("--path and --owner are required")
			}
			store, err := audit.LoadOwnership(file)
			if err != nil {
				return err
			}
			store.Owners = append(store.Owners, audit.OwnerEntry{
				Path: path, Owner: owner, Team: team, Contact: contact,
			})
			return audit.SaveOwnership(file, store)
		},
	}
	addCmd.Flags().String("path", "", "Secret path prefix")
	addCmd.Flags().String("owner", "", "Owner name")
	addCmd.Flags().String("team", "", "Team name")
	addCmd.Flags().String("contact", "", "Contact email")

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List ownership entries",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := audit.LoadOwnership(file)
			if err != nil {
				return err
			}
			if len(store.Owners) == 0 {
				fmt.Println("No ownership entries found.")
				return nil
			}
			for _, o := range store.Owners {
				fmt.Fprintf(os.Stdout, "path=%-30s owner=%-15s team=%-15s contact=%s\n",
					o.Path, o.Owner, o.Team, o.Contact)
			}
			return nil
		},
	}

	ownershipCmd.PersistentFlags().StringVar(&file, "file", "ownership.json", "Ownership store file")
	ownershipCmd.AddCommand(addCmd, listCmd)
	rootCmd.AddCommand(ownershipCmd)
}
