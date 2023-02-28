package yaml

import (
	"fmt"
	"net"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/devicetype"
)

type Environment struct {
	Environments map[string]Simulator `yaml:"environments"`
}

type Simulator struct {
	Number        uint   `yaml:"number"`
	DeviceType    string `yaml:"type"`
	MacAddress    string `yaml:"startMacAddress"`
	StartPreFixIp string `yaml:"startPreFixIp"`
}

func NewSimulatorFile(filename, ethname string) ([]*simulator.AtopGwdClient, error) {
	value, err := parsing(filename)
	if err != nil {
		return nil, err
	}
	return NewSimulator(value.Environments, ethname)
}

func NewSimulator(evnironments map[string]Simulator, ethname string) ([]*simulator.AtopGwdClient, error) {
	err := checkSimulator(evnironments)
	if err != nil {
		return nil, err
	}
	count := 0
	simulators := []*simulator.AtopGwdClient{}
	for _, v := range evnironments {
		var cidr string
		var mac string
		for i := 1; i <= int(v.Number); i++ {
			count++
			if i != 1 {
				c, err := NextPrefix(cidr)
				if err != nil {
					return nil, err
				}
				cidr = c

				if mac != "" {
					m, err := NextMac(mac)
					if err != nil {
						return nil, err
					}
					mac = m
				}

			} else {
				cidr = v.StartPreFixIp
				mac = v.MacAddress
			}

			simulator, err := selectedSimulator(uint(count), ethname, cidr, mac, v)
			if err != nil {
				return nil, err
			}
			simulators = append(simulators, simulator)
		}
	}
	if len(simulators) == 0 {
		return nil, fmt.Errorf("device number is 0")
	}

	return simulators, nil
}

func checkSimulator(sims map[string]Simulator) error {
	for _, v := range sims {
		_, b := devicetype.ParseString(v.DeviceType)
		if !b {
			return fmt.Errorf("type error:%v,please input:%v", v.DeviceType, devicetype.ArraySimulator)
		}
		_, _, err := net.ParseCIDR(v.StartPreFixIp)
		if err != nil {
			return fmt.Errorf("Start_prefixip:%v,error:%v", v.StartPreFixIp, err.Error())
		}
		/*_, err = net.ParseMAC(v.MacAddress)
		if err != nil {
			return fmt.Errorf("MacAddress:%v,error:%v", v.MacAddress, err.Error())
		}*/
	}
	return nil
}

func selectedSimulator(id uint, ethname, cidr, mac string, s Simulator) (*simulator.AtopGwdClient, error) {
	device, _ := devicetype.ParseString(s.DeviceType)

	if mac == "" {
		return simulator.NewAtopSimulatorCidrRandom(id, ethname, cidr, device)
	}
	return simulator.NewAtopSimulatorCidr(id, ethname, mac, cidr, device)
}
