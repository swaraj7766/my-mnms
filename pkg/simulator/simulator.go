package simulator

import (
	"bytes"
	"errors"
	"net"

	"mnms/pkg/simulator/devicetype"
	"mnms/pkg/simulator/firmware"
	"mnms/pkg/simulator/pcap"
	"mnms/pkg/simulator/telnet"

	"os"
	"strconv"
	"strings"
	"time"

	atopnet "mnms/pkg/simulator/net"
	atop "mnms/pkg/simulator/snmp"
	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/qeof/q"
)

const account = "admin"
const pwd = "default"
const public = "public"
const private = "private"
const packetlen = 300

const port = "55954"
const fiter = "udp and port 55954"

// NewAtopSimulator create new atopdevice simulate withe random mac,random ip,random deviceType
//
// max=253 * 253
func NewAtopSimulator(id uint, ethname string) (*AtopGwdClient, error) {
	Model, _ := GetTestParam("test", id)
	ip, err := atopnet.GetRandPrefix(ethname)
	if err != nil {
		return nil, err
	}
	n, g, err := net.ParseCIDR(ip)
	if err != nil {
		return nil, err
	}
	Model.IPAddress = n.String()
	Model.Netmask = net.IP(g.Mask).String()

	Model.MACAddress = GetRandMac()
	Model.Gateway = "10.0.0.254"

	sndata := snmpvalue.NewBindValue(&Model.Model)
	s := atop.NewSnmp([]string{public, private}, sndata)
	ts := telnet.NewTelnetServer(Model.Model, Model.IPAddress, account, pwd)
	fw := firmware.NewFirmwareServer(Model.IPAddress)

	d := &AtopGwdClient{ModelInfo: &Model, dataflag: true, powerflag: false, ethname: ethname, snmp: s, user: account, pwd: pwd, telnetserver: ts, firmware: fw}
	d.bindValue()
	err = checkLenHost(d.ModelInfo.Hostname)
	if err != nil {
		return nil, err
	}
	p, err := createPcap(d.ethname)
	if err != nil {
		return nil, err
	}
	p.RegisterReceiveEvent(d.ModelInfo.MACAddress, d.Receive)
	d.pcap = p

	return d, err
}

// NewAtopSimulatorCidr create AtopGwdClient with fixed mac, fixed ip, fixed deviceType
func NewAtopSimulatorCidr(id uint, ethname, mac, cidr string, device devicetype.Simulator_type) (*AtopGwdClient, error) {
	Model, _ := GetTestParamModel("test", id, device)
	n, g, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	Model.IPAddress = n.String()
	Model.Netmask = net.IP(g.Mask).String()
	_, err = net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}
	Model.MACAddress = mac

	Model.Gateway = "10.0.0.254"

	sndata := snmpvalue.NewBindValue(&Model.Model)
	s := atop.NewSnmp([]string{public, private}, sndata)
	ts := telnet.NewTelnetServer(Model.Model, Model.IPAddress, account, pwd)
	fw := firmware.NewFirmwareServer(Model.IPAddress)

	d := &AtopGwdClient{ModelInfo: &Model, dataflag: true, powerflag: false, ethname: ethname, snmp: s, user: account, pwd: pwd, telnetserver: ts, firmware: fw}
	d.bindValue()
	err = checkLenHost(d.ModelInfo.Hostname)
	if err != nil {
		return nil, err
	}
	p, err := createPcap(d.ethname)
	if err != nil {
		return nil, err
	}
	p.RegisterReceiveEvent(d.ModelInfo.MACAddress, d.Receive)
	d.pcap = p

	return d, err
}

// NewAtopSimulatorCidrRandom create AtopGwdClient with randomMac,fixed fixed ip,fixed deviceType
func NewAtopSimulatorCidrRandom(id uint, ethname, cidr string, device devicetype.Simulator_type) (*AtopGwdClient, error) {
	Model, _ := GetTestParamModel("test", id, device)
	n, g, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	Model.IPAddress = n.String()
	Model.Netmask = net.IP(g.Mask).String()
	Model.MACAddress = GetRandMac()

	Model.Gateway = "10.0.0.254"

	sndata := snmpvalue.NewBindValue(&Model.Model)
	s := atop.NewSnmp([]string{public, private}, sndata)
	ts := telnet.NewTelnetServer(Model.Model, Model.IPAddress, account, pwd)
	fw := firmware.NewFirmwareServer(Model.IPAddress)

	d := &AtopGwdClient{ModelInfo: &Model, dataflag: true, powerflag: false, ethname: ethname, snmp: s, user: account, pwd: pwd, telnetserver: ts, firmware: fw}
	d.bindValue()
	err = checkLenHost(d.ModelInfo.Hostname)
	if err != nil {
		return nil, err
	}
	p, err := createPcap(d.ethname)
	if err != nil {
		return nil, err
	}
	p.RegisterReceiveEvent(d.ModelInfo.MACAddress, d.Receive)
	d.pcap = p

	return d, err
}

type AtopGwdClient struct {
	ModelInfo    *ModelInfo
	user         string
	pwd          string
	dataflag     bool
	powerflag    bool
	hanlder      CompleteSettingEvent
	ethname      string
	snmp         *atop.Snmp
	pcap         *pcap.InstancePcap
	telnetserver *telnet.TelnetServer
	firmware     *firmware.Firmware
}

func createPcap(ethName string) (*pcap.InstancePcap, error) {
	pcapServer, err := pcap.NewinstancePcap(ethName)
	if err != nil {
		return nil, err
	}
	pcapServer.SetBPFFilter(fiter)
	return pcapServer, nil
}

// send device info
func (a *AtopGwdClient) sendDevice() error {

	a.updateDeviceInfo() // for being synchronized with snmp simulation
	err := a.broadcast(a.getInfoPacket())
	if err != nil {
		return err
	}
	return nil

}

func checkLenHost(name string) error {
	if len(name) > 15 {
		return errors.New("name len than 15")
	}
	return nil
}

// SettingDevice setting device info
func (a *AtopGwdClient) settingDevice(msg []byte) {
	a.setDataflag(false)

	newip := make([]string, 4)
	for i := 0; i < 4; i++ {
		newip[i] = strconv.Itoa(int(msg[16+i]))
	}
	newmask := make([]string, 4)
	for i := 0; i < 4; i++ {
		newmask[i] = strconv.Itoa(int(msg[236+i]))
	}

	gate := make([]string, 4)
	for i := 0; i < 4; i++ {
		gate[i] = strconv.Itoa(int(msg[24+i]))
	}

	hostname := make([]string, 0)
	for i := 0; i < 15; i++ {
		v := msg[90+i]
		if v == 0 {
			break
		}
		hostname = append(hostname, string(msg[90+i]))
	}

	a.DeleteIPaddress(a.ModelInfo.IPAddress, a.ModelInfo.Netmask)

	q.Q("SettingDevice ,device ip:", a.ModelInfo.IPAddress)
	a.ModelInfo.IPAddress = strings.Join(newip[:], ".")
	a.ModelInfo.Netmask = strings.Join(newmask[:], ".")
	a.ModelInfo.Gateway = strings.Join(gate[:], ".")
	host := strings.Join(hostname[:], "")
	a.ModelInfo.Hostname = host
	// for being synchronized with snmp simulation if necessary
	_, err := var_readFile("sysName")
	if err == nil {
		err := var_writeFile("sysName", a.ModelInfo.Hostname)
		if err != nil {
			q.Q(err)
		}
	}

	if a.hanlder != nil {
		a.hanlder(a.ModelInfo.IPAddress, a.ModelInfo.Netmask, a.ModelInfo.Gateway, a.ModelInfo.Hostname)

	}
	a.Reboot()

}

// Shutdown device Shutdown
func (a *AtopGwdClient) Shutdown() error {
	if a.getpowerflag() {
		a.pcap.Close()
		a.setpowerflag(false)
		a.setDataflag(false)
		a.SnmpShutdown()
		a.TelnetServerShutdown()
		a.FirmwareShutdown()
		a.DeleteIPaddress(a.ModelInfo.IPAddress, a.ModelInfo.Netmask)
		q.Q("device: Shutdown", a.ModelInfo.IPAddress)
	}
	return nil
}

// StartUp device startup
func (a *AtopGwdClient) StartUp() error {
	if !a.getpowerflag() {
		err := a.pcap.Run()
		if err != nil {
			return err
		}
		a.AddIPaddress(a.ModelInfo.IPAddress, a.ModelInfo.Netmask)
		a.updateDeviceInfo() // for being synchronized with snmp simulation
		err = a.SnmpRun(a.ModelInfo.IPAddress)
		if err != nil {
			q.Q(err)
		}
		a.setpowerflag(true)
		a.setDataflag(true)
		a.TelnetServerRun()
		a.FirmwareRun()
		time.Sleep(time.Millisecond * 300)
		q.Q("device:", a.ModelInfo.IPAddress, " start up")
	}
	return nil
}

// AddIPaddress add ip into interface
func (a *AtopGwdClient) AddIPaddress(ip, mask string) {
	m := atopnet.CovertMaskToLen(mask)
	prfixip := ip + "/" + strconv.Itoa(m)
	q.Q("ethName:", a.ethname, " add ip:", prfixip)
	err := atopnet.AddIPAddress(a.ethname, prfixip)
	if err != nil {
		q.Q(err)
	}

}

// DeleteIPaddress delete ip from interface
func (a *AtopGwdClient) DeleteIPaddress(ip, mask string) {
	m := atopnet.CovertMaskToLen(mask)
	prfixip := ip + "/" + strconv.Itoa(m)
	q.Q("ethName:", a.ethname, " delete ip:", prfixip)
	_ = atopnet.DeleteIPAddress(a.ethname, prfixip)

}

// PowerStatus get device powerstatus ,on:true,off:false
func (a *AtopGwdClient) PowerStatus() bool {
	return a.getpowerflag()
}

func (a *AtopGwdClient) setpowerflag(v bool) {
	a.powerflag = v
}
func (a *AtopGwdClient) getpowerflag() bool {
	v := a.powerflag
	return v
}

func (a *AtopGwdClient) getDataflag() bool {
	v := a.dataflag
	return v
}
func (a *AtopGwdClient) setDataflag(v bool) {
	a.dataflag = v
}

// Reboot device reboot
func (a *AtopGwdClient) Reboot() {
	_ = a.Shutdown()
	time.Sleep(time.Second * 10)
	_ = a.StartUp()

}

// Beep make device beep
func (a *AtopGwdClient) Beep() {
	q.Q("device:", a.ModelInfo.IPAddress, " beep .....")
	q.Q("device:", a.ModelInfo.IPAddress, " beep .....")
}

// broadcast broadcast
func (a *AtopGwdClient) broadcast(msg []byte) error {
	addr := net.JoinHostPort("255.255.255.255", port)
	broadcastAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	//outip := a.outip
	q.Q("IP:", a.ModelInfo.IPAddress)
	q.Q("Mask:", a.ModelInfo.Netmask)
	q.Q("GateWay:", a.ModelInfo.Gateway)
	q.Q("send info.....")
	q.Q("device:", a.ModelInfo.IPAddress, " boradcast to:", addr)
	local, err := net.ResolveUDPAddr("udp", net.JoinHostPort(a.ModelInfo.IPAddress, strconv.Itoa(0)))
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
	return nil
}

// Receive receive data when data incoming
func (a *AtopGwdClient) Receive(b []byte) {
	if len(b) == packetlen {
		q.Q("device:", a.ModelInfo.IPAddress, " receive", a.ModelInfo.IPAddress)
		if a.getDataflag() && a.getpowerflag() {
			r := SelectPacket(b)
			switch r {
			case invite:

				err := a.sendDevice()
				if err != nil {
					q.Q(err)
				}

			case config:
				if a.compareMac(b) && a.comparAccountAndPwd(b) {

					a.settingDevice(b)
				}
			case reboot:
				if a.compareMac(b) && a.comparAccountAndPwd(b) {
					a.Reboot()
				}
			case beep:
				if a.compareMac(b) {
					a.Beep()
				}
			case none:
			}
		}
	}
}

func (a *AtopGwdClient) compareMac(m []byte) bool {
	if strings.Contains(a.ModelInfo.MACAddress, "-") {
		r := bytes.Compare(MacToByte(a.ModelInfo.MACAddress, "-"), m[28:34])
		return r == 0
	} else {
		r := bytes.Compare(MacToByte(a.ModelInfo.MACAddress, ":"), m[28:34])
		return r == 0
	}
}

func (a *AtopGwdClient) comparAccountAndPwd(m []byte) bool {
	v := make([]byte, 0)
	for i := 71; i <= 88; i++ {
		if m[i] == 0 {
			break
		}
		v = append(v, m[i])
	}
	value := bytes.Split(v, []byte(" "))
	n := string(value[0])
	p := string(value[1])
	if n == a.user && p == a.pwd {
		return true
	}

	return false
}

// updateDeviceInfo  update device info from external files for being synchronized with snmp simulation
func (a *AtopGwdClient) updateDeviceInfo() {
	var_value, err := var_readFile("devModel")
	if err == nil {
		a.ModelInfo.Model = var_value
	}
	var_value, err = var_readFile("devKernel")
	if err == nil {
		a.ModelInfo.Kernel = var_value
	}
	var_value, err = var_readFile("devApplication")
	if err == nil {
		a.ModelInfo.Ap = var_value
	}
	var_value, err = var_readFile("sysName")
	if err == nil {
		a.ModelInfo.Hostname = var_value
	}

	//a.upSnmpData()

}

// getInfoPacket change ModelInfo to packet []byte
func (a *AtopGwdClient) getInfoPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 0x01
	packet[4] = 0x92
	packet[5] = 0xDA

	modinfo := []byte(a.ModelInfo.Model)
	for i := 0; i < len(modinfo); i++ {
		packet[44+i] = modinfo[i]
	}
	ip := strings.Split(a.ModelInfo.IPAddress, ".")
	for i := 0; i < len(ip); i++ {
		v, _ := strconv.Atoi(ip[i])
		packet[12+i] = byte(v)
	}
	mac := spliteMac(a.ModelInfo.MACAddress)
	for i := 0; i < len(mac); i++ {
		v, _ := strconv.ParseUint(mac[i], 16, 8)
		packet[28+i] = byte(v)
	}
	Netmask := strings.Split(a.ModelInfo.Netmask, ".")
	for i := 0; i < len(Netmask); i++ {
		v, _ := strconv.Atoi(Netmask[i])
		packet[236+i] = byte(v)
	}

	Gateway := strings.Split(a.ModelInfo.Gateway, ".")
	for i := 0; i < len(Gateway); i++ {
		v, _ := strconv.Atoi(Gateway[i])
		packet[24+i] = byte(v)
	}

	Hostname := []byte(a.ModelInfo.Hostname)
	for i := 0; i < len(Hostname); i++ {
		packet[90+i] = Hostname[i]
	}
	ap := []byte(a.ModelInfo.Ap)
	for i := 0; i < len(ap); i++ {
		addr := 110 + i
		if addr > 234 {
			break
		}
		packet[addr] = ap[i]
	}

	ker := strings.Split(a.ModelInfo.Kernel, ".")
	one, _ := strconv.Atoi(ker[0])
	two, _ := strconv.Atoi(ker[1])
	if len(ker) == 2 {
		packet[109] = byte(one)
		packet[108] = byte(two)
	}

	return packet
}

func spliteMac(mac string) []string {
	if strings.Contains(mac, "-") {
		m := strings.Split(mac, "-")
		return m
	} else {
		m := strings.Split(mac, ":")
		return m
	}

}

// read/write external files for being synchronized with snmp simulation
const parameter_path = "/tmp/parameters/"

func var_readFile(var_name string) (string, error) {
	filename := parameter_path + var_name

	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func var_writeFile(var_name string, var_value string) error {
	filename := parameter_path + var_name
	data := []byte(var_value)

	err := os.WriteFile(filename, data, 0644)
	//	if err != nil {
	//	}
	return err
}
