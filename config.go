package mnms

import (
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
// reboot, reset and beep are not here. They are basic commands in cmd.go
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

	dev, err := FindDev(devId)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	// check ip existed
	if currentIp != newip {
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
	err = ArpScan()
	if err != nil {
		q.Q(err)
	}
	if currentIp != newip {
		r, _ := ArpCheckExisted(newip)
		if !r {
			syslogerr := SendSyslog(LOG_ALERT, "ConfigNet", dev.Mac+" set ip fail")
			if syslogerr != nil {
				q.Q(syslogerr)
			}
			cmdinfo.Status = "set ip fail"
			return cmdinfo
		}
	}

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
		"server-port":  ".10.1.2.3.0:OctetString",
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
		return nil, fmt.Errorf("error: dev not support snmp: %v", err)
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
	dev, err := FindDev(devId)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
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
	_, err = SnmpSet(IPAddress, setting.Oid, value, setting.Type)
	if err != nil {
		return fmt.Errorf("error: snmp set %v", err)
	}

	return nil
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

	if strings.HasPrefix(cmd, "config reset ") {
		newcmd := strings.ReplaceAll(cmd, "config ", "")
		cmdinfo.Command = newcmd
		rcmd := ResetCmd(cmdinfo)
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
	if strings.HasPrefix(cmd, "config local syslog ") && QC.IsRoot {
		return configOfLocalSyslogCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "config crontab ") && QC.IsRoot {
		return CrontabCmd(cmdinfo)
	}

	q.Q("unrecognized", cmd, len(cmd))
	cmdinfo.Status = "error: unknown command"
	return cmdinfo
}

func configOfLocalSyslogCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command

	if strings.HasPrefix(cmd, "config local syslog path ") && QC.IsRoot {
		return SyslogSetPathCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "config local syslog maxsize ") && QC.IsRoot {
		return SyslogSetMaxSizeCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config local syslog compress ") && QC.IsRoot {
		return SyslogSetCompressCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config local syslog read") && QC.IsRoot {
		return ReadSyslogCmd(cmdinfo)
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
	if strings.HasPrefix(cmds, "config snmp enable ") {
		return fmt.Sprintf("%s %s %s %s %s", "switch", ws[3], ws[4], ws[5], "snmp"), nil
	}
	if strings.HasPrefix(cmds, "config snmp disable ") {
		return fmt.Sprintf("%s %s %s %s %s %s", "switch", ws[3], ws[4], ws[5], "no", "snmp"), nil
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
