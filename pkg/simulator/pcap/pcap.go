package pcap

import (
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

type receiveEvent func(b []byte)

var (
	snapshotLen int32         = 65536
	promiscuous bool          = false
	timeout     time.Duration = 1 * time.Microsecond
)

type InstancePcap struct {
	event  map[string]receiveEvent
	filter string
	device string
	handle *pcap.Handle
}

// NewinstancePcap create pcap by ehtname
func NewinstancePcap(ethname string) (*InstancePcap, error) {
	device, err := GetAdaptersName(ethname)
	if err != nil {
		return nil, err
	}
	/*handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return nil, err
	}*/

	return &InstancePcap{device: device, event: make(map[string]receiveEvent)}, nil
}

// Run start to capture
func (i *InstancePcap) Run() error {

	handle, err := pcap.OpenLive(i.device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return err
	}
	err = handle.SetBPFFilter(i.filter)
	if err != nil {
		return err
	}
	i.handle = handle
	go func() {
		packetSource := gopacket.NewPacketSource(i.handle, i.handle.LinkType())
		for packet := range packetSource.Packets() {
			for _, handler := range i.event {
				applicationLayer := packet.ApplicationLayer()
				if applicationLayer != nil {
					go handler(applicationLayer.Payload())
				}
			}
		}
	}()
	return nil
}

// Close close Pcap
func (i *InstancePcap) Close() {
	i.handle.Close()
}

// SetBPFFilter set filter
func (i *InstancePcap) SetBPFFilter(filter string) {
	i.filter = filter

}

// RegisterEvent reigister reveice event
func (i *InstancePcap) RegisterReceiveEvent(name string, fn receiveEvent) {
	i.event[name] = fn
}

// RemoveReceiveEvent remove reveice event
func (i *InstancePcap) RemoveReceiveEvent(name string) {
	delete(i.event, name)
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
