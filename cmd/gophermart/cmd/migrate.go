package cmd

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/vstdy0/go-diploma/cmd/gophermart/cmd/common"
)

// newMigrateCmd creates a new migrate cmd.
func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate DB to the latest version",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)

			st, err := config.BuildPsqlStorage()
			if err != nil {
				return err
			}
			defer func() {
				if err = st.Close(); err != nil {
					log.Error().Err(err).Msg("Shutting down the app")
				}
			}()

			ctx, ctxCancel := context.WithTimeout(context.Background(), config.Timeout)
			defer ctxCancel()

			if err = st.Migrate(ctx); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}
