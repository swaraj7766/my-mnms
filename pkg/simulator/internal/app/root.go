package app

import (
	"github.com/spf13/cobra"
)

//var address = "localhost:" + config.GetgrpcPort()

var RootCmd = &cobra.Command{
	Use: "simulator",
}

func init() {
	RootCmd.AddCommand(NewRunCmd())
}

func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		return err
	}
	return nil
}
