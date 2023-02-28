//go:build windows
// +build windows

package yaml_test

import (
	"log"
	"testing"
	"time"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"
)

var ethName string

func TestMain(m *testing.M) {
	var err error
	ethName, err = net.GetDefaultInterfaceName()
	if err != nil {
		log.Fatal(err)
	}
	m.Run()

}
func TestSimulatorFile(t *testing.T) {

	simulators, err := atopyaml.NewSimulatorFile("./test.yaml", ethName)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range simulators {
		_ = v.StartUp()
		defer func(v *simulator.AtopGwdClient) {
			_ = v.Shutdown()
		}(v)
	}
	time.Sleep(300 * time.Millisecond)
}

func TestSimulator(t *testing.T) {
	simmap := map[string]atopyaml.Simulator{}
	simmap["1"] = atopyaml.Simulator{Number: 5, DeviceType: "EH7506", StartPreFixIp: "192.168.6.1/24"}
	simulators, err := atopyaml.NewSimulator(simmap, ethName)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range simulators {
		_ = v.StartUp()
		defer func(v *simulator.AtopGwdClient) {
			_ = v.Shutdown()
		}(v)
	}
	time.Sleep(300 * time.Millisecond)
}
