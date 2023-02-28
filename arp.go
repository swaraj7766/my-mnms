// revive:disable-line:package-comments
package mnms

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/qeof/q"
)

const timeout = 1 * time.Millisecond

// ArpCheckExisted  existent:ture, nonexistent:false
//
// if run on linux please use root
func ArpCheckExisted(ip string) (bool, error) {
	ifaces, _ := GetAllInterface()
	for _, iface := range ifaces {
		if r, err := arpCheck(&iface, net.ParseIP(ip)); err != nil {
			return false, err
		} else {
			if r {
				return r, nil
			}
		}
	}

	// Wait for all interfaces' scans to complete.  They'll try to run
	// forever, but will stop on an error, so if we get past this Wait
	// it means all attempts to write have failed.
	return false, nil
}

// GetAllInterface get  all of Interface and name
func GetAllInterface() ([]net.Interface, string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	var ifs []net.Interface
	var summary string
	for _, v := range ifaces {
		if v.Flags&net.FlagUp == net.FlagUp && v.Flags&net.FlagLoopback != net.FlagLoopback {
			addrs, err := v.Addrs()
			if err != nil { // get addresses
				return nil, ""
			}
			for _, addr := range addrs { // get ipv4 address
				if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
					summary = addSummary(summary, fmt.Sprintf("%v", v.Name))
					ifs = append(ifs, v)
					break
				}
			}

		}
	}

	return ifs, summary
}

func addSummary(summary string, newline string) string {
	return fmt.Sprintf("%s\n%s", summary, newline)
}

// GetAdaptersName get AdaptersName
func GetAdaptersName(name string) (string, error) {
	switch runtime.GOOS {
	case "windows":
		ifter, err := net.InterfaceByName(name)
		if err != nil {
			return "", err
		}
		addrs, err := ifter.Addrs()
		if err != nil { // get addresses
			return "", err
		}
		if (len(addrs)) == 0 {
			return "", fmt.Errorf("addr len:%v", len(addrs))
		}
		var ipv4Addr net.IP
		for _, addr := range addrs { // get ipv4 address
			if ipv4Addr = addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
				break
			}
		}

		if ipv4Addr == nil {
			return "", fmt.Errorf("interface %s don't have an ipv4 address", name)
		}

		devices, err := pcap.FindAllDevs()
		if err != nil {
			return "", err
		}
		for _, device := range devices {
			for _, v := range device.Addresses {
				if net.IP.Equal(v.IP.To4(), ipv4Addr.To4()) {
					return device.Name, nil
				}
			}
		}

	default:
		return name, nil

	}
	return "", fmt.Errorf("can't get AdaptersName%v", name)

}

func arpCheck(iface *net.Interface, ip net.IP) (bool, error) {
	// We just look for IPv4 addresses, so try to find if the interface has one.
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return false, err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					if net.IP.Equal(ip4, ip.To4()) {
						return true, nil
					}
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return false, errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return false, errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return false, errors.New("mask means network is too large")
	}
	//log.Printf("Using network range %v for interface %v", addr, iface.Name)
	name, err := GetAdaptersName(iface.Name)
	if err != nil {
		return false, err
	}
	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(name, 65536, true, timeout)
	if err != nil {
		return false, err
	}
	defer handle.Close()

	// Start up a goroutine to read in packet data.
	//stop := make(chan struct{})
	t := time.Duration(500 * time.Millisecond)
	result := false

	end := make(chan bool)
	go func() {

		result = readARPAndCheck(handle, iface, ip, t)
		end <- true
	}()
	///for {
	// Write our scan packets out to the handle.
	if err := writeARP(handle, iface, addr, ip); err != nil {
		log.Printf("error writing packets on %v: %v", iface.Name, err)
		return false, err
	}
	<-end
	// We don't know exactly how long it'll take for packets to be
	// sent back to us, but 10 seconds should be more than enough
	// time ;)
	//time.Sleep(10 * time.Second)
	//}
	return result, nil
}

func readARPAndCheck(handle *pcap.Handle, iface *net.Interface, ip net.IP, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case <-ctx.Done():
			return false
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
				// This is a packet I sent.
				continue
			}

			//log.Printf("IP %v is at %v", net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
			if bytes.Equal(arp.SourceProtAddress, ip.To4()) {
				return true
			}
			continue

		}
	}
}

func writeARP(handle *pcap.Handle, iface *net.Interface, addr *net.IPNet, dest net.IP) error {
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(addr.IP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	//for _, ip := range ips(addr) {
	arp.DstProtAddress = dest.To4()
	err := gopacket.SerializeLayers(buf, opts, &eth, &arp)
	if err != nil {
		return err
	}
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}
	//}
	return nil
}

// ArpIntervalCmd set arp interval in seconds, default 60
// range 1-3600
func ArpIntervalCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	interval, err := strconv.Atoi(ws[2])
	if err != nil {
		q.Q(err)
		cmdinfo.Status = fmt.Sprintf("error: %v", err)
		return cmdinfo
	}
	if interval < 1 || interval > 3600 {
		cmdinfo.Status = "error: interval range 1-3600 seconds"
		return cmdinfo
	}
	QC.ArpInterval = interval
	cmdinfo.Result = fmt.Sprintf("set interval to %v seconds", interval)
	cmdinfo.Status = "ok"
	return cmdinfo
}

func ArpCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command

	if strings.HasPrefix(cmd, "arp interval") {
		return ArpIntervalCmd(cmdinfo)
	}

	cmdinfo.Status = "error: not found command"
	return cmdinfo
}
