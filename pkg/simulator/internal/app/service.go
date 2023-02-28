package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRunCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "simulator run",
		Run: func(cmd *cobra.Command, args []string) {
			SetLogLevel(cmd)
			n, err := cmd.Flags().GetUint16("number")
			if err != nil {
				logrus.Fatal(err)

			}
			if n == 0 {
				logrus.Fatal("no simulator exist")
			}

			ethName, err := cmd.Flags().GetString("ethName")
			if err != nil {
				logrus.Fatal(err)

			}

			yaml, err := cmd.Flags().GetString("yaml")
			if err != nil {
				logrus.Fatal(err)

			}

			if yaml != "" {
				yamlSimulator(ethName, yaml)
			} else {
				normalSimulator(n, ethName)

			}

		},
	}
	interfs, value := net.GetAllInterface()
	name, err := net.GetDefaultInterfaceName()
	if err != nil {
		if len(interfs) == 0 {
			panic("no interface exist")
		}
		name = interfs[0].Name
	}
	SetLogLevelFlag(runCmd)
	runCmd.Flags().Uint16P("number", "n", 1, "number of simulator")
	runCmd.Flags().StringP("ethName", "e", name, fmt.Sprintf("Network Interface Name (ip bind in Network Interface selected)\nexample:%v", value))
	runCmd.Flags().StringP("yaml", "y", "", "path of yaml file,use yaml to decide simulator type and number")
	return runCmd
}

func normalSimulator(n uint16, ethName string) {

	for i := 1; uint16(i) <= n; i++ {

		d, err := simulator.NewAtopSimulator(uint(i), ethName)
		if err != nil {
			log.Fatal(err)
		}
		_ = d.StartUp()
		defer func() {
			_ = d.Shutdown()
		}()
	}
	log.Printf("simulator number:%v", n)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func yamlSimulator(ethName string, path string) {

	simulators, err := atopyaml.NewSimulatorFile(path, ethName)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range simulators {
		_ = v.StartUp()
		defer func(v *simulator.AtopGwdClient) {
			_ = v.Shutdown()
		}(v)
		//pcapServer.RegisterReceiveEvent(v.ModelInfo.MACAddress, v.Receive)
	}
	log.Printf("simulator number:%v", len(simulators))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
