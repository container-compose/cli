package main

import (
	"github.com/container-compose/cli/cmd/start"
	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use: "container-compose",
	}
)

func main() {
	Execute()
}

func Execute() error {
	return cmd.Execute()
}

func init() {
	start.RegisterCommand(cmd)
}
