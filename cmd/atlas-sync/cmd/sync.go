package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanUrakhov/atlas-sync/internal/atlassian"
	"github.com/RomanUrakhov/atlas-sync/internal/bqwriter"
	"github.com/RomanUrakhov/atlas-sync/internal/config"
	"github.com/RomanUrakhov/atlas-sync/internal/syncer"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Atlassian Goals & Projects to BigQuery",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, stop := signal.NotifyContext(
			context.Background(),
			os.Interrupt,
			syscall.SIGTERM,
		)
		defer stop()

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		atlasClient := atlassian.NewClient(
			cfg.AtlassianCloudID,
			cfg.AtlassianEmail,
			cfg.AtlassianAPIToken,
			"https://api.atlassian.com/graphql",
		)
		bqWriter, err := bqwriter.NewBigQueryWriter(
			ctx, cfg.BigQueryProject,
			cfg.BigQueryDataset,
		)
		if err != nil {
			return fmt.Errorf("create bigquery writer: %w", err)
		}
		s := syncer.NewSyncer(atlasClient, bqWriter)
		return s.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
