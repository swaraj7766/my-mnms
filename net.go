package mnms

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/qeof/q"
	"github.com/tatsushid/go-fastping"
)

// Get local IP list
func GetLocalIP() ([]string, error) {
	ips := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var ifas []net.Interface
	for _, v := range ifaces {
		if SkipIf(v) {
			continue
		}
		if v.Flags&net.FlagUp == net.FlagUp && v.Flags&net.FlagLoopback != net.FlagLoopback {
			ifas = append(ifas, v)

		}
	}
	for _, ifa := range ifas {
		address, err := ifa.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range address {
			// check the address type and if it is not a loopback the display it
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
					break
				}
			}
		}
	}
	return ips, nil
}

// check ip format
func CheckIPAddress(ip string) error {
	if net.ParseIP(ip) == nil {
		return fmt.Errorf("IP Address: %s - Invalid", ip)
	} else {
		return nil
	}
}

// check Mac format
func CheckMacAddress(mac string) error {
	_, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("MAC Address: %s - Invalid", mac)
	} else {
		return nil
	}
}

// check device exist
func CheckDeviceExisted(ip string) bool {
	r := false
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", ip)
	if err != nil {
		return false
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		r = true
	}
	p.OnIdle = func() {
		q.Q("check finish")
	}
	err = p.Run()
	if err != nil {
		q.Q(err)
		return false
	}
	return r
}

// CheckIsInSubnet check ip if in subnet
//
// ip:"192.168.5.25" network:"192.168.5.1/24"
//
// return true
func CheckIsInSubnet(ip, network string) (bool, error) {
	_, subnet, err := net.ParseCIDR(network)
	if err != nil {
		return false, err
	}
	iP := net.ParseIP(ip)
	if subnet.Contains(iP) {
		return true, nil
	} else {
		return false, nil
	}

}

// GetInterfaceIps get ipv4 ips of Interface
func GetInterfaceIps(name string) ([]string, error) {
	ief, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}
	addrs, err := ief.Addrs()
	if err != nil {
		return nil, err
	}
	ips := make([]string, 0)

	for _, addr := range addrs { // get ipv4 address
		if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
			ips = append(ips, ipv4Addr.To4().String())
		}
	}
	return ips, nil
}

// GetAllInterfaceIps get ips of all Interface
func GetAllInterfaceIps() (map[string][]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	m := make(map[string][]string, 0)
	for _, v := range ifaces {
		ips, err := GetInterfaceIps(v.Name)
		if err != nil {
			break
		}
		m[v.Name] = ips
	}

	return m, nil
}

func GetIpAddrs(cidr string) ([]netip.Addr, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, err
	}

	var ips []netip.Addr
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr)
	}

	if len(ips) < 2 {
		return ips, nil
	}

	return ips[1 : len(ips)-1], nil
}

func IfnetAddresses() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error: cannot get net interfaces: %v", err)
	}

	ipaddrs := []string{}

	for _, iface := range ifaces {
		ipaddr, err := GetIfaceIp(iface)
		if err != nil {
			continue
		}

		ipaddrs = append(ipaddrs, ipaddr)
	}
	return ipaddrs, nil
}

func IfnetCidrs() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error: cannot get net interfaces: %v", err)
	}

	ipaddrs := []string{}

	for _, iface := range ifaces {
		ipaddr, err := GetIfaceCidr(iface)
		if err != nil {
			continue
		}

		ipaddrs = append(ipaddrs, ipaddr)
	}
	return ipaddrs, nil
}

func IfnetBroadcasts() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("error: cannot get net interfaces: %v", err)
	}

	ipaddrs := []string{}

	for _, iface := range ifaces {
		ipaddr, err := GetIfaceBroadcast(iface)
		if err != nil {
			continue
		}
		q.Q(ipaddr)

		ipaddrs = append(ipaddrs, ipaddr)
	}
	return ipaddrs, nil
}

func GetIfaceIp(iface net.Interface) (string, error) {
	cidr, err := GetIfaceCidr(iface)

	if err != nil {
		return "", err
	}
	if cidr == "" {
		return "", fmt.Errorf("can't find cidr for %s", iface.Name)
	}

	q.Q(cidr)

	ipPart, _, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}

	ip, err := netip.ParseAddr(ipPart.String())
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

func GetIfaceBroadcast(iface net.Interface) (string, error) {
	cidr, err := GetIfaceCidr(iface)

	if err != nil {
		return "", err
	}
	if cidr == "" {
		return "", fmt.Errorf("can't find cidr for %s", iface.Name)
	}

	q.Q(cidr)
	_, ipnetPart, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", err
	}
	q.Q(ipnetPart.IP[0])
	q.Q(ipnetPart.IP[3])
	q.Q(ipnetPart.Mask[0])
	q.Q(ipnetPart.Mask[3])

	if ipnetPart.Mask[0] != 0xff || ipnetPart.Mask[1] != 0xff {
		return "", fmt.Errorf("first two bytes not 0xff mask, network too big")
	}

	//XXX we only handle full octet masks case
	if ipnetPart.Mask[2] == 0 {
		ipnetPart.IP[2] = 255
	}
	if ipnetPart.Mask[3] == 0 {
		ipnetPart.IP[3] = 255
	}
	ip, err := netip.ParseAddr(ipnetPart.IP.String())
	if err != nil {
		return "", err
	}

	return ip.String(), nil
}

func GetIfaceCidr(iface net.Interface) (string, error) {
	q.Q(iface)
	addrs, err := iface.Addrs()
	if err != nil {
		q.Q("error: iface addr", err)
		return "", fmt.Errorf("can't get iface addr for %s", iface.Name)
	}
	for _, addr := range addrs {
		switch addr.(type) {
		//switch v := addr.(type) {
		case *net.IPAddr:
			//q.Q("ip addr", iface, v, addr)
		case *net.IPNet:
			//q.Q("ip addr", iface, v, addr)

			if SkipIf(iface) {
				continue
			}
			cidrStr := addr.String()
			ipPart, ipnetPart, err := net.ParseCIDR(cidrStr)
			if err != nil {
				q.Q(err)
			}
			ipStr := ipPart.String()
			ip, err := netip.ParseAddr(ipStr)

			if err != nil {
				q.Q("error: parse addr", err)
				continue
			}
			q.Q(ipStr, ip, ipPart, ipnetPart)

			if ip.IsLoopback() {
				continue
			}
			if ip.IsMulticast() {
				continue
			}
			if ip.IsUnspecified() {
				continue
			}

			if !ip.Is4() {
				continue
			}

			if strings.HasPrefix(ipStr, "169.254.") {
				continue
			}
			q.Q(cidrStr)
			return cidrStr, nil
		default:
			//q.Q("default", iface, v, addr)
		}
	}
	return "", fmt.Errorf("error: can't get ip for %v", iface.Name)
}

func GetIfaceAddr(iface *net.Interface) (*net.IPNet, error) {
	addrs, err := iface.Addrs()

	if err != nil {
		return nil, err
	}
	var addr *net.IPNet

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				addr = &net.IPNet{
					IP:   ip4,
					Mask: ipnet.Mask[len(ipnet.Mask)-4:],
				}
				break
			}
		}
	}
	if addr == nil {
		return nil, errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return nil, errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return nil, errors.New("mask means network is too large")
	}

	return addr, nil
}

func GetInterfaceByIP(ipaddr string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		q.Q(err)
		return nil, err
	}
	ip := net.ParseIP(ipaddr)

	for _, iface := range ifaces {
		var addr *net.IPNet

		addr, err := GetIfaceAddr(&iface)
		if err != nil {
			q.Q(err)
			continue
		}

		if addr.Contains(ip) {
			return &iface, nil
		}
	}
	return nil, fmt.Errorf("error: can't find interface for %s", ipaddr)
}

func BogusIf(name string, description string) bool {
	if name == "lo" || name == "bluetooth-monitor" || name == "nflog" ||
		name == "nfqueue" || name == "\\Device\\NPF_Loopback" ||
		strings.HasPrefix(name, "docker") ||
		strings.HasPrefix(name, "br-") ||
		strings.HasPrefix(name, "veth") ||
		strings.HasPrefix(description, "Hyper-V") ||
		strings.HasPrefix(description, "Bluetooth") ||
		strings.HasPrefix(description, "Leaf Networks") ||
		strings.HasPrefix(description, "VMware") ||
		strings.HasPrefix(description, "Microsoft Wi-Fi Direct") ||
		strings.HasPrefix(description, "WireGuard") ||
		strings.HasPrefix(description, "WAN Miniport") {
		return true
	}
	return false
}

func GetAllInterfaces() ([]net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	gIfaces = []net.Interface{}

	for _, iface := range ifaces {
		if SkipIf(iface) {
			continue
		}

		gIfaces = append(gIfaces, iface)
	}
	return gIfaces, nil
}

func SkipIf(iface net.Interface) bool {
	if BogusIf(iface.Name, "") {
		return true
	}
	if iface.Flags&net.FlagUp == 0 {
		return true
	}
	if iface.Flags&net.FlagBroadcast == 0 {
		return true
	}
	return false
}

func GetIfaceBroadcastMulti(iface net.Interface) ([]string, error) {
	var bcastAddrs []string

	cidrs, err := GetIfaceCidrMulti(iface)
	if err != nil {
		return nil, err
	}
	for _, cidr := range cidrs {
		_, ipnetPart, err := net.ParseCIDR(cidr)
		if err != nil {
			q.Q(err)
			continue
		}
		if ipnetPart.Mask[0] != 0xff ||
			ipnetPart.Mask[1] != 0xff {
			q.Q("error: network too big")
			continue
		}
		if ipnetPart.Mask[2] == 0 {
			ipnetPart.IP[2] = 255
		}
		if ipnetPart.Mask[3] == 0 {
			ipnetPart.IP[3] = 255
		}
		ip, err := netip.ParseAddr(ipnetPart.IP.String())
		if err != nil {
			q.Q(err)
			continue
		}
		bcastAddrs = append(bcastAddrs, ip.String())
	}

	return bcastAddrs, nil
}

func GetIfaceCidrMulti(iface net.Interface) ([]string, error) {
	q.Q(iface)
	addrs, err := iface.Addrs()
	if err != nil {
		q.Q("error: iface addr", err)
		return nil, fmt.Errorf("can't get iface addr for %s", iface.Name)
	}
	var cidrs []string
	cidrStr := ""
	for _, addr := range addrs {
		switch addr.(type) {
		//switch v := addr.(type) {
		case *net.IPAddr:
			//q.Q("ip addr", iface, v, addr)
		case *net.IPNet:
			//q.Q("ip addr", iface, v, addr)

			if SkipIf(iface) {
				continue
			}
			cidrStr = addr.String()
			ipPart, ipnetPart, err := net.ParseCIDR(cidrStr)
			if err != nil {
				q.Q(err)
				continue
			}
			ipStr := ipPart.String()
			ip, err := netip.ParseAddr(ipStr)
			q.Q(ipStr, ip, ipPart, ipnetPart)

			if err != nil {
				q.Q("error: parse addr", err)
				continue
			}
			if ip.IsLoopback() {
				continue
			}
			if ip.IsMulticast() {
				continue
			}
			if ip.IsUnspecified() {
				continue
			}

			if !ip.Is4() {
				continue
			}

			if strings.HasPrefix(ipStr, "169.254.") {
				continue
			}
			q.Q(cidrStr)
			cidrs = append(cidrs, cidrStr)
		default:
			//q.Q("default", iface, v, addr)
		}
	}
	return cidrs, nil
}
