package mnms

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"

	"github.com/qeof/q"
)

var (
	GwdPort    string         = ":55954"
	GwdUDPPort layers.UDPPort = 55954
)

type GwdModelInfo struct {
	Model      string `json:"model"`
	MACAddress string `json:"macAddress"`
	IPAddress  string `json:"ipAddress"`
	Netmask    string `json:"netmask"`
	Gateway    string `json:"gateway"`
	Hostname   string `json:"hostname"`
	Kernel     string `json:"kernel"`
	Ap         string `json:"ap"`
	ScannedBy  string `json:"scannedBy"`
	// IsDHCP     bool   `json:"isDHCP"`
}

type GwdNetworkConfig struct {
	MACAddress   string `json:"macAddress"`
	IPAddress    string `json:"ipAddress"`
	NewIPAddress string `json:"newIPAddress"`
	Netmask      string `json:"netmask"`
	Gateway      string `json:"gateway"`
	Hostname     string `json:"hostname"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

var gwdTotalReceived int

func GwdMain() {
	err := GwdProcess()
	if err != nil {
		q.Q(err)
	}
}

var gIfaces []net.Interface

func GwdProcess() error {
	var err error
	gIfaces, err = GetAllInterfaces()
	if err != nil {
		q.Q(err)
		return err
	}

	gwdTotalReceived = 0

	var wg sync.WaitGroup

	allDevs, err := pcap.FindAllDevs()
	if err != nil {
		q.Q(err)
	}

	var ifaceNames []string

	for _, dev := range allDevs {
		if BogusIf(dev.Name, dev.Description) {
			continue
		}
		ifaceNames = append(ifaceNames, dev.Name)
	}

	q.Q(ifaceNames)

	for _, ifaceName := range ifaceNames {
		wg.Add(1)
		go func(ifaceName string) {
			defer wg.Done()
			// pcap input filtering

			err := interfaceInput(ifaceName)
			if err != nil {
				q.Q(err)
			}
		}(ifaceName)
	}

	wg.Wait()

	return nil
}

func interfaceInput(ifaceName string) error {
	// Linux UDP sockets do not get broadcast, capture raw frames

	handle, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		q.Q("error: cannot open interface", ifaceName)
		return err
	}
	q.Q("opened interface", ifaceName)
	defer handle.Close()

	var filter string = "udp and port 55954"
	// var filter string = "ip"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		q.Q(err)
		return err
	}

	stop := make(chan struct{})

	// go readPackets(handle, iface, stop)
	readPackets(handle, stop)

	defer close(stop)

	return nil
}

func readPackets(handle *pcap.Handle, stop chan struct{}) {
	// packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := packetSource.Packets()

	_ = GwdInvite()

	var packet gopacket.Packet
	for {
		select {
		case <-stop:
			return
		case packet = <-in:
			// q.Q(packet)
			processPacket(packet)
		}
	}
}

func processPacket(packet gopacket.Packet) {
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		q.Q(ethernetPacket.SrcMAC, ethernetPacket.DstMAC, ethernetPacket.EthernetType)
	}
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		ip, _ := ipLayer.(*layers.IPv4)
		q.Q(ip.SrcIP, ip.DstIP, ip.Length, ip.Protocol)
	}

	udpLayer := packet.Layer(layers.LayerTypeUDP)
	if udpLayer != nil {
		udp, _ := udpLayer.(*layers.UDP)
		if udp.DstPort != GwdUDPPort {
			return
		}
		q.Q(udp.SrcPort, udp.DstPort, udp.Length)
		_, err := gwdParse(udp.Payload)
		if err != nil {
			q.Q(err)
		}
	}

	if err := packet.ErrorLayer(); err != nil {
		q.Q("error: decoding some part of the packet:", err)
	}
}

func GwdInvite() error {
	q.Q("send gwd invite")
	err := GwdBroadcast(gwdInvitePacket())
	if err != nil {
		q.Q(err)
		return err
	}
	return nil
}

func GwdBroadcast(msg []byte) error {
	ips, _ := GetLocalIP()
	for _, ip := range ips {
		/*bcastAddr, err := GetIfaceBroadcast(iface)
		if err != nil {
			return err
		}*/
		err := SendBcast(ip, "255.255.255.255"+GwdPort, msg)
		if err != nil {
			q.Q(err)
		}
		// Sending to 255.255.255.255 may fail on Linux but try anyway
		//SendBcast("255.255.255.255"+GwdPort, msg)
	}
	return nil
}

func GwdReset(ipaddr string, macaddr string, username string, password string) error {
	n := GwdNetworkConfig{
		IPAddress: ipaddr, MACAddress: macaddr,
		Username: username, Password: password,
	}
	b, err := gwdSetResetDefault(n)
	if err != nil {
		q.Q(err)
		return err
	}
	return GwdBroadcast(b)
}

func GwdReboot(ipaddr string, macaddr string, username string, password string) error {
	n := GwdNetworkConfig{
		IPAddress: ipaddr, MACAddress: macaddr,
		Username: username, Password: password,
	}
	b, err := gwdSetRebootPacket(n)
	if err != nil {
		q.Q(err)
		return err
	}
	return GwdBroadcast(b)
}

func GwdBeep(ipaddr string, macaddr string) error {
	n := GwdNetworkConfig{IPAddress: ipaddr, MACAddress: macaddr}
	b, err := gwdSetBeepPacket(n)
	if err != nil {
		q.Q(err)
		return err
	}
	return GwdBroadcast(b)
}

func SendBcast(laddr string, bcastAddr string, msg []byte) error {
	local, err := net.ResolveUDPAddr("udp", net.JoinHostPort(laddr, strconv.Itoa(0)))
	if err != nil {
		return err
	}
	broadcastAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", local, broadcastAddr)

	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(msg)
	if err != nil {
		return err
	}
	q.Q(laddr, bcastAddr, len(msg))
	return nil
}

func gwdParse(msg []byte) (GwdModelInfo, error) {
	if msg[0] == 1 && msg[4] == 0x92 && msg[5] == 0xDA {
		model := GwdModelInfo{}
		model.Model = CleanStr(toUtf8(msg[44:60]))
		model.MACAddress = byteToHexString(msg[28:34], "-")
		model.IPAddress = byteToString(msg[12:16], ".")
		model.Netmask = byteToString(msg[236:240], ".")
		model.Gateway = byteToString(msg[24:28], ".")
		model.Hostname = CleanStr(toUtf8(msg[90:106]))
		model.Kernel = CleanStr(fmt.Sprintf("%d.%d", msg[109], msg[108]))
		model.Ap = CleanStr(toUtf8(msg[110:235]))
		model.ScannedBy = QC.Name
		model.Model = strings.TrimSpace(model.Model)
		q.Q(model)

		if model.Model != "" {
			InsertModel(model, "gwd")
			gwdTotalReceived++
			q.Q(gwdTotalReceived)
		}
		return model, nil
	}

	return GwdModelInfo{}, errors.New("gwdParse: error, does not look like gwd packet")
}

// copied from yanlin's gwd work

func gwdInvitePacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 2
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func gwdConfigPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 0
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func gwdRebootPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 5
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func gwdBeepPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 7
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func gwdReSetDefaultPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 5
	packet[1] = 1
	packet[2] = 6
	packet[4] = 0x92
	packet[5] = 0xDA
	return packet
}

func gwdSetBeepPacket(config GwdNetworkConfig) ([]byte, error) {
	packet := gwdBeepPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	return packet, nil
}

func gwdSetResetDefault(config GwdNetworkConfig) ([]byte, error) {
	packet := gwdReSetDefaultPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet = addUserAndPwdToPackage(config, packet)
	return packet, nil
}

func parseIp(ip string) ([]byte, error) {
	v := strings.Split(ip, ".")
	if len(v) != 4 {
		return nil, errors.New("IPAddress format error")
	}
	var ips [4]byte
	for i := 0; i < 4; i++ {
		r, err := strconv.Atoi(v[i])
		if err != nil {
			return nil, fmt.Errorf("error: IPAddress format error ,reason:%s", err.Error())
		}
		ips[i] = byte(r)
	}

	return ips[:], nil
}

func parseMacAddress(mac string) ([]byte, error) {
	m := strings.ReplaceAll(mac, "-", "")
	m = strings.ReplaceAll(m, ":", "")
	v, err := hex.DecodeString(m)
	if err != nil {
		return nil, fmt.Errorf("error: MacAddress format error ,reason:%s", err.Error())
	}
	return v, nil
}

func encodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewEncoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func addIpToPackage(config GwdNetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parseIp(config.IPAddress)
	if err != nil {
		return nil, errors.New("IPAddress format error")
	}
	for i := 0; i < 4; i++ {
		packet[12+i] = ip[i]
	}
	return packet, nil
}

func addNewIpToPackage(config GwdNetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parseIp(config.NewIPAddress)
	if err != nil {
		return nil, errors.New("new IPAddress format error")
	}
	for i := 0; i < 4; i++ {
		packet[16+i] = ip[i]
	}
	return packet, nil
}

func addNewNetmaskToPackage(config GwdNetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parseIp(config.Netmask)
	if err != nil {
		return nil, errors.New("netmask format error")
	}
	for i := 0; i < 4; i++ {
		packet[236+i] = ip[i]
	}
	return packet, nil
}

func addNewGatewayToPackage(config GwdNetworkConfig, packet []byte) ([]byte, error) {
	ip, err := parseIp(config.Gateway)
	if err != nil {
		return nil, errors.New("gateway format error")
	}
	for i := 0; i < 4; i++ {
		packet[24+i] = ip[i]
	}
	return packet, nil
}

func addhostNameToPackage(config GwdNetworkConfig, packet []byte) ([]byte, error) {
	i := 90
	host := []byte(config.Hostname)
	h, err := encodeBig5(host)
	if err == nil {
		host = h
	}
	if len(host) >= 16 {
		return nil, fmt.Errorf("HostName len is too long")
	}
	for _, v := range host {
		packet[i] = v
		i++
	}
	return packet, nil
}

func addMacToPackage(config GwdNetworkConfig, packet []byte) ([]byte, error) {
	mac, err := parseMacAddress(config.MACAddress)
	if err != nil {
		return nil, fmt.Errorf("error: MacAddress format error ,reason:%s", err.Error())
	}
	for i := 0; i < 6; i++ {
		packet[28+i] = mac[i]
	}
	return packet, nil
}

func addUserAndPwdToPackage(config GwdNetworkConfig, packet []byte) []byte {
	packet[70] = 2
	i := 71
	user := []byte(config.Username)
	pwd := []byte(config.Password)
	for _, v := range user {
		packet[i] = v
		i++
	}
	packet[i] = 0x20
	i++
	for _, v := range pwd {
		packet[i] = v
		i++
	}
	return packet
}

func byteToHexString(msg []byte, sep string) string {
	str := make([]string, len(msg))
	for i, v := range msg {
		str[i] = fmt.Sprintf("%02X", v)
	}
	return strings.Join(str, sep)
}

func byteToString(msg []byte, sep string) string {
	str := make([]string, len(msg))
	for i, v := range msg {
		str[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(str, sep)
}

func toUtf8(b []byte) string {
	b = getValidByte(b)
	s, err := decodeBig5(b)
	if err != nil {
		return string(b)
	}
	return string(s)
}

func decodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// clear '\0' in packet
func getValidByte(src []byte) []byte {
	var str_buf []byte
	for _, v := range src {
		if v != 0 {
			str_buf = append(str_buf, v)
		} else {
			break
		}
	}
	return str_buf
}

func gwdSetConfigPacket(config GwdNetworkConfig) ([]byte, error) {
	packet := gwdConfigPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addNewIpToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addNewGatewayToPackage(config, packet)
	if err != nil {
		return nil, err
	}
	packet, err = addNewNetmaskToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	packet, err = addhostNameToPackage(config, packet)
	if err != nil {
		return nil, err
	}

	packet = addUserAndPwdToPackage(config, packet)
	return packet, nil
}

func gwdSetRebootPacket(config GwdNetworkConfig) ([]byte, error) {
	packet := gwdRebootPacket()
	packet, err := addIpToPackage(config, packet)
	if err != nil {
		return nil, errors.New("IPAddress format error")
	}
	packet, err = addMacToPackage(config, packet)
	if err != nil {
		return nil, fmt.Errorf("MacAddress format error ,reason:%s", err.Error())
	}
	packet = addUserAndPwdToPackage(config, packet)
	return packet, nil
}
