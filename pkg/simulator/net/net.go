package net

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/liupeidong0620/gateway"
	"github.com/seancfoley/ipaddress-go/ipaddr"
)

func GetInterfaceIPv4(name string) ([]net.IPNet, error) {
	inter, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	addrs, err := inter.Addrs()
	if err != nil {
		return nil, err
	}
	var result []net.IPNet
	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			net, ok := addr.(*net.IPNet)
			if ok {
				result = append(result, *net)
			}
		}
	}
	return result, nil
}

func CovertMaskToLen(mask string) int {
	pref := ipaddr.NewIPAddressString(mask).GetAddress().
		GetBlockMaskPrefixLen(true)
	return pref.Len()
}

const defaultPreip = "10.234."
const mask = "16"
const count = 253 * 253
const max = 254
const min = 1

// GetRandIP get Rand new ip
//
// example,10.234.xxx.xxx/16
func GetRandPrefix(name string) (string, error) {
	for i := 0; i < count; i++ {
		same := false
		rand.Seed(time.Now().UnixNano())
		var sb strings.Builder
		sb.WriteString(defaultPreip)
		rand.Seed(time.Now().UnixNano())
		ip3 := rand.Intn(max-min) + min
		sb.WriteString(strconv.Itoa(ip3))
		ip4 := rand.Intn(max-min) + min
		sb.WriteString(".")
		sb.WriteString(strconv.Itoa(ip4))
		resip := net.ParseIP(sb.String())

		ips, err := GetInterfaceIPv4(name)
		if err != nil {
			return "", err
		}
		for _, v := range ips {
			if v.String() == resip.String() {

				same = true
			}
		}
		if !same {

			return sb.String() + "/" + mask, nil
		}

	}
	return "", fmt.Errorf("can't create ip becasue count than%v:", count)
}

/*
func GetInterFace(name string) (*net.Interface, error) {
	return net.InterfaceByName(name)

}
*/
//GetAllInterface get  all of Interface and name
func GetAllInterface() ([]net.Interface, string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	var ifs []net.Interface
	var summary string
	for _, v := range ifaces {
		if v.Flags&net.FlagUp == net.FlagUp && v.Flags&net.FlagLoopback != net.FlagLoopback {
			summary = addSummary(summary, fmt.Sprintf("%v", v.Name))
			ifs = append(ifs, v)

		}
	}

	return ifs, summary
}

// GetDefaultInterfaceName get defualt gateway name
func GetDefaultInterfaceName() (string, error) {
	data, err := gateway.DiscoverInterface()
	if err == nil {
		return data.Inte.Name, nil
	}
	return "", err
}

func addSummary(summary string, newline string) string {
	return fmt.Sprintf("%s\n%s", summary, newline)
}
