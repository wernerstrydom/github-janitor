/*
Copyright © 2024 Werner Strydom <hello@wernerstrydom.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// archiveCmd represents the archive command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive repositories that are considered empty",
	Long: `This command scans all repositories in the specified organization,
identifies repositories that are considered empty (containing only a README.md,
LICENSE, or .gitignore file, and not updated in the last month), and archives them.`,

	Run: func(cmd *cobra.Command, args []string) {
		if client == nil || ctx == nil {
			log.Fatal("GitHub client is not initialized")
		}

		// Confirm action with the user
		if !confirmAction("Are you sure you want to archive all empty repositories? This action cannot be undone.") {
			return
		}

		err := ForEachRepositories(client, ctx, organization, isRepositoryEmpty, archiveRepository)
		if err != nil {
			log.Fatalf("Error processing repositories: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// archiveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// archiveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
