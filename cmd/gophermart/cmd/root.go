package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vstdy0/go-diploma/api"
	"github.com/vstdy0/go-diploma/cmd/gophermart/cmd/common"
	"github.com/vstdy0/go-diploma/pkg"
)

const (
	flagConfigPath         = "config"
	flagLogLevel           = "log-level"
	flagTimeout            = "timeout"
	flagRunAddress         = "run_address"
	flagDatabaseURI        = "database_uri"
	flagStorageType        = "storage_type"
	flagAccrualSysAddress  = "accrual_system_address"
	envSecretKey           = "secret_key"
	envUpdaterTimeout      = "updater_timeout"
	envStatusCheckInterval = "status_check_interval"
)

// Execute prepares cobra.Command context and executes root cmd.
func Execute() error {
	return newRootCmd().ExecuteContext(common.NewBaseCmdCtx())
}

// newRootCmd creates a new root cmd.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := setupLogger(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			if err := setupConfig(cmd); err != nil {
				return fmt.Errorf("app initialization: %w", err)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			config := common.GetConfigFromCmdCtx(cmd)
			svcCtx, svcCancel := context.WithCancel(context.Background())
			defer svcCancel()

			svc, err := config.BuildService(svcCtx)
			if err != nil {
				return fmt.Errorf("app initialization: service building: %w", err)
			}

			srv := api.NewServer(svc, config)

			go func() {
				if err = srv.ListenAndServe(); err != http.ErrServerClosed {
					log.Error().Err(err).Msg("HTTP server ListenAndServe")
				}
			}()

			stop := make(chan os.Signal, 1)
			signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
			<-stop

			svcCancel()

			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer shutdownCancel()

			if err = srv.Shutdown(shutdownCtx); err != nil {
				return fmt.Errorf("shutting shutdown server: %w", err)
			}

			if err = svc.Close(); err != nil {
				log.Error().Err(err).Msg("shutting down service")
			}

			log.Info().Msg("server stopped")

			return nil
		},
	}

	config := common.BuildDefaultConfig()
	cmd.PersistentFlags().String(flagConfigPath, "./config.toml", "Config file path")
	cmd.PersistentFlags().String(flagLogLevel, "info", "Logger level [debug,info,warn,error,fatal]")
	cmd.PersistentFlags().Duration(flagTimeout, config.Timeout, "Request timeout")
	cmd.PersistentFlags().StringP(flagDatabaseURI, "d", config.PSQLStorage.URI, "Database URI")
	cmd.Flags().StringP(flagRunAddress, "a", config.RunAddress, "Run address")
	cmd.Flags().StringP(flagStorageType, "s", config.StorageType, "Storage type [psql]")
	cmd.Flags().StringP(flagAccrualSysAddress, "r", config.Provider.AccrualSysAddress, "Accruals system address")

	cmd.AddCommand(newMigrateCmd())

	return cmd
}

// setupLogger configures global logger.
func setupLogger(cmd *cobra.Command) error {
	logLevelBz, err := cmd.Flags().GetString(flagLogLevel)
	if err != nil {
		return fmt.Errorf("%s flag reading: %w", flagLogLevel, err)
	}
	logLevel, err := zerolog.ParseLevel(logLevelBz)
	if err != nil {
		return fmt.Errorf("%s flag parsing: %w", flagLogLevel, err)
	}

	logWriter := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	log.Logger = log.Output(logWriter).Level(logLevel)

	return nil
}

// setupConfig reads app config and stores it to cobra.Command context.
func setupConfig(cmd *cobra.Command) error {
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return fmt.Errorf("flags binding: %w", err)
	}

	if err := viper.BindEnv(envSecretKey); err != nil {
		return fmt.Errorf("%s env binding: %w", envSecretKey, err)
	}
	if err := viper.BindEnv(envUpdaterTimeout); err != nil {
		return fmt.Errorf("%s env binding: %w", envUpdaterTimeout, err)
	}
	if err := viper.BindEnv(envStatusCheckInterval); err != nil {
		return fmt.Errorf("%s env binding: %w", envStatusCheckInterval, err)
	}

	configPath := viper.GetString(flagConfigPath)
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return fmt.Errorf("reading config file: %w", err)
		}
	}

	viper.AutomaticEnv()
	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	config := common.BuildDefaultConfig()
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("config unmarshal: %w", err)
	}

	config.Timeout = viper.GetDuration(flagTimeout)
	common.SetConfigToCmdCtx(cmd, config)

	if config.SecretKey == "" {
		return fmt.Errorf("%s config: %w", envSecretKey, pkg.ErrNoValue)
	}

	return nil
}
