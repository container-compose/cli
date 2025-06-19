package start

import (
	"log/slog"
	"os"

	"github.com/container-compose/cli/internal/entities"
	"github.com/container-compose/cli/internal/logger"
	"github.com/spf13/cobra"
)

var (
	file string
	cmd  = &cobra.Command{
		Use: "start",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			ctx, logger := logger.New(ctx, os.Stdout, slog.LevelDebug)
			logger.InfoContext(ctx, "starting containers", "file", file)

			// attempt to read the file and get the contents
			contents, err := os.ReadFile(file)
			if err != nil {
				logger.ErrorContext(ctx, err.Error())
				return
			}

			// parse the config
			config, err := entities.Parse(contents)
			if err != nil {
				logger.ErrorContext(ctx, err.Error())
				return
			}

			// start the services
			for _, service := range config.Services {

				// check if the service is already running
				isRunning, err := service.IsRunning(ctx)
				if err != nil {
					logger.ErrorContext(ctx, err.Error())
					return
				}
				if isRunning {
					continue
				}

				// if we already have the service, but it's not running, start it
				exists, err := service.Exists(ctx)
				if err != nil {
					logger.ErrorContext(ctx, err.Error())
					return
				}
				if exists {
					// start it back up
					cmd, err := service.StartCommand(ctx)
					if err != nil {
						logger.ErrorContext(ctx, err.Error())
						return
					}
					err = cmd.Exec(ctx)
					if err != nil {
						logger.ErrorContext(ctx, err.Error())
						return
					}
					logger.InfoContext(ctx, "started service", "name", service.Name)
					continue
				}

				// start the service
				cmd, err := service.RunCommand(ctx)
				if err != nil {
					logger.ErrorContext(ctx, err.Error())
					return
				}
				err = cmd.Exec(ctx)
				if err != nil {
					logger.ErrorContext(ctx, err.Error())
					return
				}
				logger.InfoContext(ctx, "started service", "name", service.Name)
			}

			return
		},
	}
)

func init() {
	cmd.PersistentFlags().StringVarP(&file, "file", "f", "compose.yaml", "the compose file")
}

func RegisterCommand(parent *cobra.Command) {
	parent.AddCommand(cmd)
}
