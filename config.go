package mnms

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/qeof/q"
)

var ErrorDeviceExisted = errors.New("device ip exist")

// There are three types of config: basic, snmp and switch cli
//
// this file implements basic config.
//
// The basic config are :
//    host, ip,  trap server, syslog server, enable ip,
//    username, password, etc.
//
// reset, mtderase and beep are not here. They are basic commands in cmd.go
//
// snmp stuff is in snmp.go
//
// switch cli stuff is in switch.go

func getGwdNetConfig(dev *DevInfo) GwdNetworkConfig {
	return GwdNetworkConfig{
		IPAddress:    dev.IPAddress,
		MACAddress:   dev.Mac,
		NewIPAddress: dev.IPAddress,
		Netmask:      dev.Netmask,
		Gateway:      dev.Gateway,
		Hostname:     dev.Hostname,
		Username:     "admin",   // XXX
		Password:     "default", // XXX
	}
}

// setting some config (like IP ) may have to be device specific.
// some devices can be set via gwd, some may not have gwd.
//  some devices may have snmp working.
//  some may or may not have switch cli.
// start with basic -- GWD.  gradually add support for devices
// that require special support (e.g. device has no gwd, only snmp).

// Use gwd to configure network setting.
//
// Usage : config net [mac address] [current ip] [new ip] [mask] [gateway] [hostname]
//
//	[mac address] : target device mac address
//	[current ip]  : target device current ip address
//	[new ip]      : target device would modify ip address
//	[mask]        : target device network mask
//	[gateway]     : target device gateway
//	[hostname]    : target device host name
//
// Example :
//
//	config net AA-BB-CC-DD-EE-FF 10.0.50.1 10.0.50.2 255.255.255.0 0.0.0.0 switch
func ConfigNet(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 8 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	devId := ws[2]
	currentIp := ws[3]
	newip := ws[4]
	mask := ws[5]
	gateway := ws[6]
	hostname := ws[7]
	// validate current ip
	err := CheckIPAddress(currentIp)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	// validate new ip
	err = CheckIPAddress(newip)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	// validate mask
	err = CheckIPAddress(mask)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	// validate gateway
	err = CheckIPAddress(gateway)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	cmdinfo.DevId = devId
	dev, err := FindDev(devId)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	b, err := DevIsLocked(devId)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	if b {
		cmdinfo.Status = fmt.Sprintf("error:%v", "device is upgrading")
		return cmdinfo
	}
	// check ip existed
	if currentIp != newip && newip != "0.0.0.0" {
		r, err := ArpCheckExisted(newip)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("error:%v", err)
			return cmdinfo
		}
		if r {
			q.Q(ErrorDeviceExisted)
			cmdinfo.Status = fmt.Sprintf("error:%v", ErrorDeviceExisted)
			return cmdinfo
		}
	}
	// TODO add group support if FindDev fails
	n := getGwdNetConfig(dev)
	n.NewIPAddress = newip
	n.Netmask = mask
	n.Gateway = gateway
	n.Hostname = hostname
	err = GwdConfig(n)
	if err != nil {
		q.Q(err)
		cmdinfo.Status = fmt.Sprintf("error: %v", err)
		return cmdinfo
	}
	// check ip config status
	//disable this
	/*if currentIp != newip {
		r, _ := ArpCheckExisted(newip)
		if !r {
			syslogerr := SendSyslog(LOG_ALERT, "ConfigNet", dev.Mac+" set ip fail")
			if syslogerr != nil {
				q.Q(syslogerr)
			}
			cmdinfo.Status = "set ip fail"
			return cmdinfo
		}
	}*/
	cmdinfo.Status = "ok"
	return cmdinfo
}

type snmpBasicSetting struct {
	Oid  string
	Type string
}

func getSnmpBasicSetting(targetIp, cate, kind string) (*snmpBasicSetting, error) {
	syslogFields := map[string]string{
		"status":       ".10.1.2.1.0:Integer",
		"server-ip":    ".10.1.2.6.0:OctetString",
		"server-port":  ".10.1.2.3.0:Integer",
		"server-level": ".10.1.2.4.0:Integer",
		"LogToFlash":   ".10.1.2.5.0:Integer",
	}

	trapServerFields := map[string]string{
		"status":      ".8.6.1.5.0:Integer",
		"server-ip":   ".8.6.1.7.0:OctetString",
		"server-port": ".8.6.1.6.0:Integer",
		"community":   ".8.6.1.3.0:OctetString",
	}
	pointToDefined := map[string]string{}

	switch cate {
	case "syslog":
		pointToDefined = syslogFields
	case "snmp-trap":
		pointToDefined = trapServerFields
	}

	objID, err := SnmpGetObjectID(targetIp)
	if err != nil {
		q.Q(err)
		return nil, fmt.Errorf("dev not support snmp: %v", err)
	}
	q.Q("objID", objID)

	t, ok := pointToDefined[kind]
	if !ok {
		q.Q("error", kind)
		return nil, fmt.Errorf("invalid syslog field %v", kind)
	}
	tt := strings.Split(t, ":")
	oid := objID + tt[0]
	datatype := tt[1]
	return &snmpBasicSetting{Oid: oid, Type: datatype}, nil
}

// Use snmp to configure syslog setting.
//
// Usage : config syslog [mac address] [status] [server ip] [server port] [server level] [log to flash]
//
//	[mac address] : target device mac address
//	[status]      : use snmp to configure syslog enable/disable
//	[server ip]   : use snmp to configure server ip address
//	[server port] : use snmp to configure server port
//	[server level]: use snmp to configure server log level
//	[log to flash]: use snmp to configure log to flash
//
// Example :
//
//	config syslog AA-BB-CC-DD-EE-FF 1 10.0.50.2 5514 1 1
func ConfigSyslog(cate string, cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 8 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	devId := ws[2]
	status := ws[3]
	serverIp := ws[4]
	serverPort := ws[5]
	serverLavel := ws[6]
	logToFlash := ws[7]
	cmdinfo.DevId = devId
	dev, err := FindDev(devId)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	// validate server ip
	err = CheckIPAddress(serverIp)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	// status setting
	err = SetSnmpOneCommand(dev.IPAddress, cate, status, "status")
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error: snmp set %v", err)
		return cmdinfo
	}
	// server ip setting
	err = SetSnmpOneCommand(dev.IPAddress, cate, serverIp, "server-ip")
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error: snmp set %v", err)
		return cmdinfo
	}
	// server port setting
	err = SetSnmpOneCommand(dev.IPAddress, cate, serverPort, "server-port")
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error: snmp set %v", err)
		return cmdinfo
	}
	// server lavel setting
	err = SetSnmpOneCommand(dev.IPAddress, cate, serverLavel, "server-level")
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error: snmp set %v", err)
		return cmdinfo
	}
	// log to flash setting
	err = SetSnmpOneCommand(dev.IPAddress, cate, logToFlash, "LogToFlash")
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error: snmp set %v", err)
		return cmdinfo
	}

	cmdinfo.Status = "ok"
	return cmdinfo
}

func SetSnmpOneCommand(IPAddress string, cate string, value string, kind string) error {
	setting, err := getSnmpBasicSetting(IPAddress, cate, kind)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	pkt, err := SnmpSet(IPAddress, setting.Oid, value, setting.Type)
	if err != nil {
		return fmt.Errorf("error: snmp set %v", err)
	}
	if uint8(pkt.Error) > 0 {
		return fmt.Errorf("error: snmp set %v", pkt.Error.String())
	}

	return nil
}

// Use snmp to get syslog setting.
//
// Usage : config getsyslog [mac address]
//
//	[mac address] : target device mac address
//
// Example :
//
//	config getsyslog 00-60-E9-18-3C-3C
func ConfigGetSyslog(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	devId := ws[2]
	dev, err := FindDev(devId)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	v, err := GetSnmpSyslogStatus(dev.IPAddress)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	b, err := json.Marshal(v)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Result = string(b)
	cmdinfo.Status = "ok"
	return cmdinfo
}

type SyslogStatus struct {
	Status     string `json:"status"`
	Serverip   string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	Level      string `json:"server_level"`
	LogToFlash string `json:"logToflash"`
}

func GetSnmpSyslogStatus(IPAddress string) (SyslogStatus, error) {
	cate := "syslog"
	result := SyslogStatus{}
	st, err := getSnmpBasicSetting(IPAddress, cate, "status")
	if err != nil {
		return result, err
	}
	ip, err := getSnmpBasicSetting(IPAddress, cate, "server-ip")
	if err != nil {
		return result, err
	}
	port, err := getSnmpBasicSetting(IPAddress, cate, "server-port")
	if err != nil {
		return result, err
	}
	lev, err := getSnmpBasicSetting(IPAddress, cate, "server-level")
	if err != nil {
		return result, err
	}
	flash, err := getSnmpBasicSetting(IPAddress, cate, "LogToFlash")
	if err != nil {
		return result, err
	}
	res, err := SnmpGet(IPAddress, []string{st.Oid, ip.Oid, port.Oid, lev.Oid, flash.Oid})
	if err != nil {
		return result, err
	}
	for _, v := range res.Variables {
		switch v.Name {
		case st.Oid:
			result.Status = PDUToString(v)
		case ip.Oid:
			result.Serverip = PDUToString(v)
		case port.Oid:
			result.ServerPort = PDUToString(v)
		case lev.Oid:
			result.Level = PDUToString(v)
		case flash.Oid:
			result.LogToFlash = PDUToString(v)
		}
	}

	return result, nil
}

func GwdConfig(conf GwdNetworkConfig) error {
	b, err := gwdSetConfigPacket(conf)
	if err != nil {
		return err
	}
	return GwdBroadcast(b)
}

func ConfigCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	q.Q("cmd = ", cmd)
	// gwd config
	if strings.HasPrefix(cmd, "config net ") {
		return ConfigNet(cmdinfo)
	}
	// syslog config (via snmp),trap config (snmp)
	if strings.HasPrefix(cmd, "config syslog ") {
		return ConfigSyslog("syslog", cmdinfo)
	}
	// gwd config
	if strings.HasPrefix(cmd, "config beep ") {
		newcmd := strings.ReplaceAll(cmd, "config ", "")
		cmdinfo.Command = newcmd
		rcmd := BeepCmd(cmdinfo)
		rcmd.Command = cmd
		return rcmd
	}
	// config get syslog
	if strings.HasPrefix(cmd, "config getsyslog ") {
		return ConfigGetSyslog(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config mtderase ") {
		newcmd := strings.ReplaceAll(cmd, "config ", "")
		cmdinfo.Command = newcmd
		rcmd := MtdEraseCmd(cmdinfo)
		rcmd.Command = cmd
		return rcmd
	}
	// enable snmp (via telnet)
	if strings.HasPrefix(cmd, "config snmp ") {
		newcmd, err := ConvertSnmpCmd(cmd)
		if err != nil {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		cmdinfo.Command = newcmd
		rcmd := SwitchCmd(cmdinfo)
		rcmd.Command = cmd
		return rcmd
	}
	// config running status to startup
	if strings.HasPrefix(cmd, "config switch save ") {
		return ConfigSwitchSaveCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config local syslog ") && QC.IsRoot {
		return configOfLocalSyslogCmd(cmdinfo)
	}

	q.Q("unrecognized", cmd, len(cmd))
	cmdinfo.Status = "error: unknown command"
	return cmdinfo
}

// ConvertSnmpCmd convert snmp cmd by telnet
func ConvertSnmpCmd(cmds string) (string, error) {
	ws := strings.Split(cmds, " ")
	if len(ws) < 6 {
		return "", errors.New("error: invalid command")
	}
	extra := []string{}
	if len(ws) > 6 {
		extra = ws[6:]
	}

	//example,input: config snmp enable AA-BB-CC-DD-EE-FF admin default
	//return:switch AA-BB-CC-DD-EE-FF admin default snmp enable
	if strings.HasPrefix(cmds, "config snmp enable ") {
		return strings.TrimSpace(fmt.Sprintf("%s %s %s %s %s %s %s", "switch", ws[3], ws[4], ws[5], "snmp", "enable", strings.Join(extra, " "))), nil
	}
	//example,input: config snmp enable AA-BB-CC-DD-EE-FF admin default
	//return:switch AA-BB-CC-DD-EE-FF admin default no snmp enable
	if strings.HasPrefix(cmds, "config snmp disable ") {
		return strings.TrimSpace(fmt.Sprintf("%s %s %s %s %s %s %s %s", "switch", ws[3], ws[4], ws[5], "no", "snmp", "enable", strings.Join(extra, " "))), nil
	}
	/*	if strings.HasPrefix(cmds, "config snmp trap ") {
			if len(ws) < 8 {
				return "", errors.New("error: invalid command")
			}
			return fmt.Sprintf("%s %s %s %s %s %s %s %s", "switch", ws[5], ws[6], ws[7], "snmp", "trap", ws[3], ws[4]), nil
		}

		if strings.HasPrefix(cmds, "config snmp trapmode ") {
			if len(ws) < 7 {
				return "", errors.New("error: invalid command")
			}
			return fmt.Sprintf("%s %s %s %s %s %s %s", "switch", ws[4], ws[5], ws[6], "snmp", "trap-mode", ws[3]), nil
		}*/

	return "", errors.New("error: invalid command")
}
