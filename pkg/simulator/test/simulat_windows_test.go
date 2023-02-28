//go:build windows
// +build windows

package test

import (
	"log"
	"testing"
	"time"

	"mnms/pkg/simulator"

	atopnet "mnms/pkg/simulator/net"
)

var ethName string

func TestMain(m *testing.M) {
	var err error
	ethName, err = atopnet.GetDefaultInterfaceName()
	if err != nil {
		log.Fatal(err)
	}

	m.Run()

}

const deviceCount = 1

func TestSimulator(t *testing.T) {
	for i := 1; i <= deviceCount; i++ {
		d, err := simulator.NewAtopSimulator(uint(i), ethName)
		if err != nil {
			log.Fatal(err)
		}
		_ = d.StartUp()
		defer func() {
			time.Sleep(time.Second * 1)
			err := d.Shutdown()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}
}
