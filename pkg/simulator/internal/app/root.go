package app

import (
	"fmt"

	"github.com/sirupsen/logrus"
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

// SetLogLevel sets service's log level
func SetLogLevel(cmd *cobra.Command) {

	if cmd == nil {
		logrus.Error("set log level fail: cmd is nil")
		return
	}
	debug, _ := cmd.Flags().GetBool("debug")

	if debug {
		fmt.Println("Running as debug mode")
		logrus.SetLevel(logrus.DebugLevel)
		return
	}

	info, _ := cmd.Flags().GetBool("verb")
	if info {
		fmt.Println("Running as verb mode")
		logrus.SetLevel(logrus.InfoLevel)
		return
	}
	fmt.Println("Running as default mode")
	logrus.SetLevel(logrus.ErrorLevel)

}

// SetLogLevelFlag set log level flag to cmd
func SetLogLevelFlag(cmd *cobra.Command) {
	if cmd == nil {
		logrus.Error("set log level flag fail: cmd is nil")
		return
	}
	cmd.Flags().BoolP("debug", "d", false, "debug level")
	cmd.Flags().BoolP("verb", "v", false, "verb level")
}
