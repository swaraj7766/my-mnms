package yaml

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"mnms/pkg/simulator/devicetype"

	"gopkg.in/yaml.v3"
)

func parsing(filename string) (Environment, error) {
	yfile, err := ioutil.ReadFile(filename)
	v := Environment{}
	if err != nil {

		return v, err
	}

	err = yaml.Unmarshal(yfile, &v)
	if err != nil {
		return v, err
	}
	err = checkParam(v)
	if err != nil {
		return v, err
	}
	return v, nil
}
func checkParam(envir Environment) error {
	for k, v := range envir.Environments {
		if v.Number == 0 {
			s := envir.Environments[k]
			s.Number = 1
			envir.Environments[k] = s
		}
		_, ok := devicetype.ParseString(v.DeviceType)
		if !ok {
			return ErrType
		}
		_, _, err := net.ParseCIDR(v.StartPreFixIp)
		if err != nil {
			return err
		}
	}
	return nil
}

//NextPrefix get nexip in format cidr
func NextPrefix(cidr string) (string, error) {
	ip, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	nip := nextIP(ip)
	m := strings.Split(cidr, "/")
	newcidr := fmt.Sprintf("%s/%s", nip, m[1])
	return newcidr, nil
}

func nextIP(ip net.IP) net.IP {
	ip = ip.To4()
	ip[3]++
	return ip
}

func NextMac(mac string) (string, error) {
	m, err := net.ParseMAC(mac)
	if err != nil {
		return "", err
	}
	m[5]++
	return m.String(), nil
}
