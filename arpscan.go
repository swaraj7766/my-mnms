package mnms

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/qeof/q"
)

func searchLostDevices() {
	QC.DevMutex.Lock()
	devlist := QC.DevData
	QC.DevMutex.Unlock()
	for _, devinfo := range devlist {
		if devinfo.ArpMissed == 2 {
			q.Q(devinfo.Mac, "offline")
			syslogerr := SendSyslog(LOG_ALERT, "ArpCheck", devinfo.Mac+" offline")
			if syslogerr != nil {
				q.Q(syslogerr)
			}
		}
	}
}

func CheckAllDevicesAlive() {
	var wg sync.WaitGroup
	for {
		for i := 0; i < QC.ArpInterval; i++ { // XXX
			time.Sleep(1 * time.Second)
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := ArpScan()
			if err != nil {
				q.Q(err)
			}
			searchLostDevices()
		}()
	}
}

func ArpScan() error {
	// Get timestamp
	timestamp := time.Now().Format(time.RFC3339)
	// Get a list of all interfaces.
	ifaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, iface := range ifaces {
		if SkipIf(iface) {
			continue
		}
		wg.Add(1)
		// Start up a scan on each interface.
		go func(iface net.Interface) {
			defer wg.Done()
			if err := scan(&iface, timestamp); err != nil {
				q.Q(iface, err)
			}
		}(iface)
	}
	// Wait for all interfaces' scans to complete.  They'll try to run
	// forever, but will stop on an error, so if we get past this Wait
	// it means all attempts to write have failed.
	wg.Wait()

	// Add ArpMissed count.
	// If arp response, ArpMissed will be set 0.
	QC.DevMutex.Lock()
	devlist := QC.DevData
	QC.DevMutex.Unlock()
	for _, devinfo := range devlist {
		if timestamp != devinfo.Timestamp {
			devinfo.ArpMissed++
			InserAndPublishDevice(devinfo)
		}
	}
	return err
}

func scan(iface *net.Interface, timestamp string) error {
	// We just look for IPv4 addresses, so try to find if the interface has one.
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return err
	} else {
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
	}
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return errors.New("mask means network is too large")
	}
	// log.Printf("Using network range %v for interface %v", addr, iface.Name)
	name, err := GetAdaptersName(iface.Name)
	if err != nil {
		return err
	}
	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(name, 65536, true, timeout)
	if err != nil {
		return err
	}
	defer handle.Close()

	// Start up a goroutine to read in packet data.
	// when we send too many requests, timeout will add more timeout.
	timeout := 3000
	t := time.Duration(time.Duration(timeout) * time.Millisecond)
	QC.DevMutex.Lock()
	devlist := QC.DevData
	QC.DevMutex.Unlock()

	end := make(chan bool)
	go func() {
		readARP(handle, iface, t, timestamp, devlist)
		end <- true
	}()

	index := 0
	for _, ip := range ips(addr) {
		// need to be careful not to send too many requests in short amount of time.
		if index > 10 {
			time.Sleep(100 * time.Millisecond)
			index = 0
		}
		if err := writeARPEach(handle, iface, addr, ip); err != nil {
			q.Q("error writing packets on %v: %v", iface.Name, err)
			return err
		}
		index = index + 1
	}

	<-end

	return nil
}

func readARP(handle *pcap.Handle, iface *net.Interface, timeout time.Duration, timestamp string, devlist map[string]DevInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case <-ctx.Done():
			return
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
			arpSrcIpaddr := net.IP(arp.SourceProtAddress).String()
			arpSrcHwaddr := strings.ReplaceAll(net.HardwareAddr(arp.SourceHwAddress).String(), ":", "-")
			arpSrcHwaddr = strings.ToUpper(arpSrcHwaddr)

			// QC.DevMutex.Lock()
			// devinfo, ok := QC.DevData[arpSrcHwaddr]
			// QC.DevMutex.Unlock()
			devinfo, ok := devlist[arpSrcHwaddr]

			if ok {
				// q.Q("update device", arpSrcIpaddr, arpSrcHwaddr)
				if devinfo.IPAddress != arpSrcIpaddr {
					devinfo.IPAddress = arpSrcIpaddr
					q.Q(arpSrcIpaddr, "new ip", arpSrcIpaddr)
					syslogerr := SendSyslog(LOG_ALERT, "ArpCheck", arpSrcHwaddr+" new IP:"+arpSrcIpaddr)
					if syslogerr != nil {
						q.Q(syslogerr)
					}
				}
				if devinfo.ArpMissed >= 2 {
					q.Q(arpSrcIpaddr, "online")
					syslogerr := SendSyslog(LOG_ALERT, "ArpCheck", arpSrcHwaddr+" online")
					if syslogerr != nil {
						q.Q(syslogerr)
					}
				}
				devinfo.Timestamp = timestamp
				devinfo.ArpMissed = 0
				InserAndPublishDevice(devinfo)
			}
			continue
		}
	}
}

// writeARP writes an ARP request for each address on our local network to the
// pcap handle.
func writeARPEach(handle *pcap.Handle, iface *net.Interface, addr *net.IPNet, dest net.IP) error {
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
	arp.DstProtAddress = dest.To4()
	gopacket.SerializeLayers(buf, opts, &eth, &arp)
	if err := handle.WritePacketData(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// ips is a simple and not very good method for getting all IPv4 addresses from a
// net.IPNet.  It returns all IPs it can over the channel it sends back, closing
// the channel when done.
func ips(n *net.IPNet) (out []net.IP) {
	num := binary.BigEndian.Uint32([]byte(n.IP))
	// mask means network is too large
	maskString := "255.255.255.0"
	mask := binary.BigEndian.Uint32([]byte(net.ParseIP(maskString).To4()))
	// mask := binary.BigEndian.Uint32([]byte(n.Mask))
	network := num & mask
	broadcast := network | ^mask
	for network++; network < broadcast; network++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], network)
		out = append(out, net.IP(buf[:]))
	}
	return
}
