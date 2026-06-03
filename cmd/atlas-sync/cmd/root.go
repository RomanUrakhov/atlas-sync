package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:           "atlas-sync",
	Short:         "Sync Atlassian Goals & Projects to BigQuery",
	SilenceErrors: true,
}

func Execute() error {
	return rootCmd.Execute()
}
