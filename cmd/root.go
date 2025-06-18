package cmd

import (
	"github.com/container-compose/cli/cmd/start"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use: "container-compose",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	start.RegisterCommand(rootCmd)
}
