package mnms

import (
	"net"
	"net/netip"
	"strconv"
	"strings"
	"testing"

	"github.com/google/gopacket/pcap"
	"github.com/qeof/q"
)

func TestLocalAddresses(t *testing.T) {
	addrs, err := IfnetAddresses()
	if err != nil {
		t.Fatalf("IfnetAddresses %v", err)
	}
	q.Q(addrs)
}

func TestGetIfnetAddresses(t *testing.T) {
	ifnetAddrs, err := IfnetAddresses()
	if err != nil {
		t.Fatal(err)
	}
	q.Q(ifnetAddrs)
}

//disable beacuse test environment couldn't have test ip
/*
func TestGetInterfaceByIP(t *testing.T) {
	iface, err := GetInterfaceByIP("192.168.100.123")
	if err != nil {
		t.Fatalf("GetInterfaceByIP %v", err)
	}

	q.Q(iface)
	iface, err = GetInterfaceByIP("1.2.3.4")
	if err == nil {
		t.Fatalf("GetInterfaceByIP expected fail")
	}
	if iface != nil {
		t.Fatalf("GetInterfaceByIP expected nil iface")
	}
	q.Q(iface)
}*/

func TestIfnetMulti(t *testing.T) {
	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}

	for _, iface := range ifaces {
		addrs, err := GetIfaceCidrMulti(iface)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(addrs)

		bcastAddrs, err := GetIfaceBroadcastMulti(iface)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(bcastAddrs)
	}
}

func TestAllNetdevs(t *testing.T) {
	allDevs, err := pcap.FindAllDevs()
	if err != nil {
		q.Q(err)
	}
	q.Q(allDevs)

	var ifaceNames []string

	for _, dev := range allDevs {
		if BogusIf(dev.Name, dev.Description) {
			continue
		}
		q.Q(dev)
		ifaceNames = append(ifaceNames, dev.Name)
	}

	q.Q(ifaceNames)
}

func TestParseAddrPort(t *testing.T) {
	addr := ":162"
	addrPort, err := netip.ParseAddrPort(addr)
	if err != nil {
		q.Q(err)
	}
	q.Q(addrPort, addrPort.Addr(), addrPort.Port())

	ws := strings.Split(addr, ":")
	port, err := strconv.Atoi(ws[1])
	if err != nil {
		q.Q(err)
	}
	q.Q(ws)
	q.Q(port)
	if port != 162 {
		t.Fatal("port should be 162")
	}
}
