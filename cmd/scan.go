package cmd

import (
	"github.com/spf13/cobra"
)

// scanCmd scans for empty repositories in the organization
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans for empty repositories in the organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		return ForEachRepositories(client, ctx, organization, isRepositoryEmpty, printRepositoryName)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
