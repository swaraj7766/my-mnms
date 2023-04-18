package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/qeof/q"
	"github.com/ziutek/telnet"
)

// ValidCommands is a list of valid commands
//
// Keep this list up to date
var ValidCommands []string = []string{
	"config", "mtderase", "beep", "reset", "scan", "switch", "snmp",
	"log", "firmware", "mqtt", "opcua", "help", "util",
}

// cmd status
type cmdStats int

const (
	running cmdStats = iota
)

func (c cmdStats) String() string {
	switch c {
	case running:
		return "running"
	}
	return "unknown"
}

// CmdInfo contains a command, a unit of API call
//
// A command is a API call that is executed in a distributed environement.
//
// A command created by the user and inserted at Root may be replicated
// to clients. A client that is capable of executing the command will
// run the command and return the result by reporting back to the Root.
type CmdInfo struct {
	Kind        string `json:"kind"`
	Timestamp   string `json:"timestamp"`
	Command     string `json:"command"`
	Result      string `json:"result"`
	Status      string `json:"status"`
	Name        string `json:"name"`
	Retries     int    `json:"retries"`
	NoOverwrite bool   `json:"nooverwrite"`
	All         bool   `json:"all"`
	NoSyslog    bool   `json:"nosyslog"`
	Client      string `json:"client"`
	DevId       string `json:"devid"`
	Tag         string `json:"tag"`
}

const telnet_timeout = 10 * time.Second // XXX

func init() {
	QC.CmdData = make(map[string]CmdInfo)
}

// InsertCmd inserts command information into command data list.
//
// It is called by  UpdateCmds when commands are posted via http.
func InsertCmd(cmd string, cmdinfo CmdInfo) {
	q.Q(cmd)
	if cmdinfo.Command == "" {
		cmdinfo.Command = cmd
		q.Q("set command", cmdinfo)
	}
	if cmdinfo.All {
		// insert an instance of the special 'all' command for each client
		for cl := range QC.Clients {
			// per client command has prefix @client cmd ...
			kcmd := "@" + cl + " " + cmd
			q.Q(kcmd)
			ci := cmdinfo
			ci.Timestamp = time.Now().Format(time.RFC3339)
			ci.Client = cl // along with @client, this indicates client cmd
			QC.CmdMutex.Lock()
			_, ok := QC.CmdData[kcmd]
			QC.CmdMutex.Unlock()
			if ok {
				if ci.NoOverwrite {
					q.Q("error: cmd exists already", ci)
					continue
				}
			}
			QC.CmdMutex.Lock()
			QC.CmdData[kcmd] = ci
			QC.CmdMutex.Unlock()
		}
		return
	}
	QC.CmdMutex.Lock()
	_, ok := QC.CmdData[cmd]
	QC.CmdMutex.Unlock()
	if ok {
		if cmdinfo.NoOverwrite {
			q.Q("error: cmd exists already", cmd)
			return
		}
	}
	cmdinfo.Timestamp = time.Now().Format(time.RFC3339)
	QC.CmdMutex.Lock()
	QC.CmdData[cmd] = cmdinfo
	QC.CmdMutex.Unlock()
}

// InsertDownCmds puts downloaded command data into local CmdData[].
//
// Root maintains its own command data list.  Each client node service
// downloads commands from the Root periodically.
//
// InsertDownCmds is called by CheckCmds which is periodically called
// from the main service go routine.
func InsertDownCmds(cmddata *map[string]CmdInfo) {
	q.Q(cmddata)
	for k, v := range *cmddata {
		QC.CmdMutex.Lock()
		_, ok := QC.CmdData[k]
		QC.CmdMutex.Unlock()
		if ok {
			if v.NoOverwrite {
				q.Q("error: cmd exists already", v)
				continue
			}
		}
		v.Name = QC.Name
		QC.CmdMutex.Lock()
		QC.CmdData[k] = v
		QC.CmdMutex.Unlock()
		q.Q("inserted cmd", v)
	}
}

// UpdateCmds is called when clients post commands via http.
//
// CLI may post commands via http to allow command line access to API.
// An http client may post commands to send command via rest API.
//
// UpdateCmds will call InsertCmd to record the command
// into local command information list.
func UpdateCmds(cmddata *map[string]CmdInfo) {
	q.Q(cmddata)
	// if root is collecting command status from clients,
	// updates from each client will be aggregated.
	// update command status in the queue
	for k, v := range *cmddata {
		if v.Status == "" {
			InsertCmd(k, v)
			continue
		}
		QC.CmdMutex.Lock()
		found, ok := QC.CmdData[k]
		QC.CmdMutex.Unlock()
		// don't override ok result command history
		if ok {
			if found.Status == "ok" {
				continue
			}
			if strings.HasPrefix(found.Status, "error:") {
				continue
			}
		}
		// fill in missing timestamp
		if v.Timestamp == "" {
			v.Timestamp = time.Now().Format(time.RFC3339)
		}
		QC.CmdMutex.Lock()
		QC.CmdData[k] = v
		QC.CmdMutex.Unlock()
		q.Q("cmd updated", found, v)
	}
}

// CheckCmds runs in client node services and periodically
// download commands from the Root service and run commands and
// update the results back to the Root service.
func CheckCmds() error {
	if QC.RootURL != "" {
		// if there is a URL to the root we download from it
		resp, err := GetWithToken(QC.RootURL+"/api/v1/commands?id="+QC.Name, QC.AdminToken)
		if err != nil {
			return err
		}
		if resp != nil {
			// save close
			defer resp.Body.Close()
			cmddata := make(map[string]CmdInfo)
			//json.NewDecoder(resp.Body).Decode(&cmddata)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &cmddata)
			if err != nil {
				return err
			}
			//q.Q("downloaded cmds", cmddata)
			if len(cmddata) > 0 {
				InsertDownCmds(&cmddata)
			}
		}
	}

	// XXX this mutex lockout can be very long
	QC.CmdMutex.Lock()
	for k, v := range QC.CmdData {
		if v.Status != "" && !strings.HasPrefix(v.Status, "pending:") {
			continue
		}
		// cannot use goroutine because of ordering
		res := RunCmd(&v)
		QC.CmdData[k] = v
		q.Q(res)
	}
	QC.CmdMutex.Unlock()

	if QC.RootURL != "" { //always check for root URL to run even when no root
		// update results back to root
		QC.CmdMutex.Lock()
		jsonBytes, err := json.Marshal(QC.CmdData)
		QC.CmdMutex.Unlock()
		// TODO delete finished and updated commands
		if err != nil {
			return err
		}
		resp, err := PostWithToken(QC.RootURL+"/api/v1/commands", QC.AdminToken, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return err
		}
		//q.Q("updated commands", QC.CmdData)
		if resp != nil {
			q.Q(resp.Header)
			// save close
			resp.Body.Close()
		}
	}
	return nil
}

func RunCmd(cmdinfo *CmdInfo) *CmdInfo {
	defer func() {
		if cmdinfo.Status != "" && !cmdinfo.NoSyslog {
			jsonBytes, err := json.Marshal(cmdinfo)
			if err != nil {
				q.Q(err)
			}
			err = SendSyslog(LOG_NOTICE, "RunCmd", string(jsonBytes))
			q.Q("sending syslog", string(jsonBytes))
			if err != nil {
				q.Q("error: sending syslog", err)
			}
		}
	}()
	q.Q(cmdinfo)
	if strings.HasPrefix(cmdinfo.Status, "error:") {
		q.Q("error: cmd Status already error, will not run")
		return cmdinfo
	}
	cmd := cmdinfo.Command
	if !cmdinfo.All && cmdinfo.Client != "" {
		if cmdinfo.Client != QC.Name {
			q.Q("error: wrong client cmd", cmd, QC.Name)
			return cmdinfo
		}
	}
	cmdinfo.Name = QC.Name
	q.Q(cmd)
	if strings.HasPrefix(cmdinfo.Status, "pending:") {
		cmdinfo.Retries++
	}
	if cmdinfo.Retries > 3 {
		if cmdinfo.DevId != "" {
			devId := cmdinfo.DevId
			dev, err := FindDev(devId)
			if err != nil {
				cmdinfo.Status = "error: cancelled, device does not exist in inventory"
				return cmdinfo
			}
			q.Q("warning: cancelling cmd for device in inventory, may need to issue explicit client specific command", dev.ScannedBy, cmdinfo.Name)
		}

		cmdinfo.Status = "error: cancelled, too many retries"
		return cmdinfo
	}
	if strings.HasPrefix(cmd, "mtderase") {
		return MtdEraseCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "beep") {
		return BeepCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "reset") {
		return ResetCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "scan") {
		return ScanCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "config") {
		return ConfigCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "switch") {
		return SwitchCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "snmp") {
		return SnmpCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "log") {
		return LogCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "firmware") {
		return FirmwareCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "mqtt") {
		return RunMqttCmd(cmdinfo)
	}

	//opcua
	if strings.HasPrefix(cmd, "opcua") {
		ws := strings.Split(cmd, " ")
		if len(ws) < 2 {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		switch ws[1] {
		case "connect":
			return OpcuaConnectCmd(cmdinfo)
		case "read":
			return OpcuaReadCmd(cmdinfo)
		case "browse":
			return OpcuaBrowseReferenceCmd(cmdinfo)
		case "sub":
			return OpcuaSubscribeCmd(cmdinfo)
		case "deletesub":
			return OpcuDeleteSubscribeCmd(cmdinfo)
		case "close":
			return OpcuCloseCmd(cmdinfo)
		}
	}

	q.Q("unrecognized", cmd, len(cmd))

	cmdinfo.Status = "error: invalid command"
	return cmdinfo
}

// Beep target device.
//
// Usage : beep [mac address] [ip address]
//
//	[mac address] : target device mac address
//	[ip address]  : target device ip address
//
// Example :
//
//	beep AA-BB-CC-DD-EE-FF 10.0.50.1
func BeepCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	macaddr := ws[1]
	ipaddr := ws[2]
	q.Q(ipaddr, macaddr)
	cmdinfo.DevId = macaddr
	dev, err := FindDev(macaddr)
	if err != nil || dev.IPAddress != ipaddr {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	// validate ipaddr
	err = CheckIPAddress(ipaddr)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	err = GwdBeep(ipaddr, macaddr)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

// Reset/Reboot target device.
//
// Usage : reset [mac address] [ip address] [username] [password]
//
//	[mac address] : target device mac address
//	[ip address]  : target device ip address
//	[username]    : target device login user name
//	[password]    : target device login passwaord
//
// Example :
//
//	reset AA-BB-CC-DD-EE-FF 10.0.50.1 admin default
func ResetCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var macaddr, ipaddr, username, password string
	Unpack(ws[1:], &macaddr, &ipaddr, &username, &password)
	cmdinfo.DevId = macaddr
	dev, err := FindDev(macaddr)
	if err != nil || dev.IPAddress != ipaddr {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	b, err := DevIsLocked(macaddr)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	if b {
		cmdinfo.Status = fmt.Sprintf("error:%v", "device is upgrading")
		return cmdinfo
	}
	// validate ipaddr
	err = CheckIPAddress(ipaddr)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	err = GwdReset(ipaddr, macaddr, username, password)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

// Erase target device mtd and restore default settings.
//
// Usage : mtderase [mac address] [ip address] [username] [password]
//
//	[mac address] : target device mac address
//	[ip address]  : target device ip address
//	[username]    : target device login user name
//	[password]    : target device login passwaord
//
// Example :
//
//	mtderase AA-BB-CC-DD-EE-FF 10.0.50.1 admin default
func MtdEraseCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var macaddr, ipaddr, username, password string
	Unpack(ws[1:], &macaddr, &ipaddr, &username, &password)
	cmdinfo.DevId = macaddr
	dev, err := FindDev(macaddr)
	if err != nil || dev.IPAddress != ipaddr {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	b, err := DevIsLocked(macaddr)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	if b {
		cmdinfo.Status = fmt.Sprintf("error:%v", "device is upgrading")
		return cmdinfo
	}

	// validate ipaddr
	err = CheckIPAddress(ipaddr)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}
	err = GwdMtdErase(ipaddr, macaddr, username, password)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

func CheckSwitchCliModel(modelname string) bool {

	cliSupportList := []string{
		"EHG7",
		"EHG9",
		"EMG8",
		"RHG9",
		"RHG7",
		"EH7",
		"Simu",
	}
	// check modelname start with cliSupportList,
	for _, v := range cliSupportList {
		// convert v and modelname to lower case
		mo := strings.ToLower(modelname)
		vv := strings.ToLower(v)
		if strings.HasPrefix(mo, vv) {
			q.Q(modelname)
			return true
		}

	}
	// if strings.HasPrefix(modelname, "EH7") || strings.HasPrefix(modelname, "EHG7") {
	// 	q.Q(modelname)
	// 	return true
	// }
	q.Q("switch cli not supported", modelname)
	return false
}

func ConvertSwitchCmd(modelname string, cmd []string) []string {
	if strings.Contains(modelname, "EHG") {
		return cmd
	}
	//judge wether is "no"
	var snmp string
	switch cmd[0] {
	case "no":
		if len(cmd) >= 2 {
			snmp = cmd[1]
		}
	default:
		snmp = cmd[0]
	}

	switch snmp {
	case "snmp":
		cmd := deleteExtraSnmpCmd(cmd)
		return cmd
	}
	return cmd
}

func deleteExtraSnmpCmd(cmd []string) []string {
	rcmd := []string{}
	for _, v := range cmd {
		if v != "enable" {
			rcmd = append(rcmd, v)
		}
	}
	return rcmd
}

// Use target device CLI configuration commands.
//
// Usage : switch [mac address] [username] [password] [cli cmd...]
//
//	[mac address] : target device mac address
//	[username]    : target device login user name
//	[password]    : target device login passwaord
//	[cli cmd...]  : target device cli command
//
// Example :
//
//	switch AA-BB-CC-DD-EE-FF admin default show ip
func SwitchCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	devId := ws[1]
	dev, err := FindDev(devId)
	cmdinfo.DevId = devId
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	if dev.ModelName == "" {
		cmdinfo.Status = "error: invalid device model"
		return cmdinfo
	}
	if !CheckSwitchCliModel(dev.ModelName) {
		cmdinfo.Status = "error: switch cli not available"
		return cmdinfo
	}
	wcmd := ConvertSwitchCmd(dev.ModelName, ws[4:])

	username := ws[2]
	password := ws[3]
	err = SendSwitch(cmdinfo, dev, username, password, strings.Join(wcmd, " "))
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

func expect(t *telnet.Conn, d ...string) {
	q.Q(d)
	err := t.SetReadDeadline(time.Now().Add(telnet_timeout))
	if err != nil {
		q.Q(err)
		return
	}
	err = t.SkipUntil(d...)
	if err != nil {
		q.Q(err)
	}
}

func sendln(t *telnet.Conn, s string) {
	q.Q(s)
	err := t.SetWriteDeadline(time.Now().Add(telnet_timeout))
	if err != nil {
		q.Q(err)
		return
	}

	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
	_, err = t.Write(buf)
	if err != nil {
		q.Q(err)
	}
}

func SendSwitch(cmdinfo *CmdInfo, dev *DevInfo, username, password, cmd string) error {
	t, err := telnet.Dial("tcp", dev.IPAddress+":23")

	if err != nil {
		q.Q(err)
		return err
	}

	defer func() {
		if err := t.Close(); err != nil {
			q.Q(err)
		}
	}()

	t.SetUnixWriteMode(true)

	expect(t, "sername: ")
	sendln(t, username+"\n")
	expect(t, "assword: ")
	sendln(t, password+"\n")
	expect(t, "#")
	sendln(t, "configure\n")
	expect(t, "#")
	sendln(t, cmd+"\n")
	sendln(t, "blah blah\n") //XXX terrible hack
	result, err := t.ReadBytes('%')
	if err != nil {
		fmt.Println(err)
		q.Q(err)
	}

	cmdinfo.Result = string(result)

	return nil
}

func SendSwitchWithoutConfig(cmdinfo *CmdInfo, dev *DevInfo, username, password, cmd string) error {
	t, err := telnet.Dial("tcp", dev.IPAddress+":23")

	if err != nil {
		q.Q(err)
		return err
	}

	defer func() {
		if err := t.Close(); err != nil {
			q.Q(err)
		}
	}()

	t.SetUnixWriteMode(true)

	expect(t, "sername: ")
	sendln(t, username+"\n")
	expect(t, "assword: ")
	sendln(t, password+"\n")
	expect(t, "#")
	sendln(t, cmd+"\n")
	sendln(t, "blah blah\n") //XXX terrible hack
	result, err := t.ReadBytes('%')
	if err != nil {
		q.Q(err)
	}
	cmdinfo.Result = string(result)

	return nil
}

// save config to device.
//
// Usage :config switch save [mac address] [username] [password]
//
//	[mac address] : target device mac address
//	[username]    : target device login user name
//	[password]    : target device login passwaord
//
// Example :
//
func ConfigSwitchSaveCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 6 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	devId := ws[3]
	dev, err := FindDev(devId)
	cmdinfo.DevId = devId
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	if dev.ModelName == "" {
		cmdinfo.Status = "error: invalid device model"
		return cmdinfo
	}
	if !CheckSwitchCliModel(dev.ModelName) {
		cmdinfo.Status = "error: switch cli not available"
		return cmdinfo
	}
	username := ws[4]
	password := ws[5]
	err = SwitchConfigSave(cmdinfo, dev, username, password)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

// SwitchConfigSave save config to device
func SwitchConfigSave(cmdinfo *CmdInfo, dev *DevInfo, username, password string) error {
	switch {
	case strings.Contains(dev.ModelName, "EHG"):
		return SendSwitchWithoutConfig(cmdinfo, dev, username, password, "copy running-config startup-config")
	default:
		return SendSwitch(cmdinfo, dev, username, password, "copy running-config startup-config")

	}
}

// Use different protocol to scan all devices.
//
// Usage : scan [protocol]
//
//	[protocol]    : use gwd/snmp to scan all devices.
//
// Example :
//
//	scan gwd
func ScanCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	if cmd == "scan gwd" {
		err := GwdInvite()
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	if cmd == "scan snmp" {
		go func() {
			err := SnmpScan()
			if err != nil {
				q.Q(err)
			}
		}()
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	cmdinfo.Status = "error: invalid command"
	return cmdinfo
}

// Use mqtt to publish/subscribe/unsubscribe/list topic.
//
// Usage : mqtt [mqttcmd] [tcp address] [topic] [data...]
//
//		[mqttcmd]     : pub/sub/unsub/list
//	                 list is show all subscribe topic
//		[tcp address] : would pub/sub/unsub broker tcp address
//		[topic]       : topic name
//		[data...]     : data is messages, only publish use it.
//
// Example :
//
//	mqtt pub 192.168.12.1:1883 topictest "this is messages."
//	mqtt sub 192.168.12.1:1883 topictest
//	mqtt unsub 192.168.12.1:1883 topictest
//	mqtt list
func RunMqttCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 2 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	selectOption := ws[1]
	if strings.HasPrefix(selectOption, "list") {
		result := DisplayAllSubscribeTopic()
		if result == "" {
			result = "not subscribe topic"
		}
		cmdinfo.Result = result
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	tcpaddr := ws[2]
	checkIP := strings.Split(ws[2], ":")
	// pass ":11883" local address
	if checkIP[0] != "" {
		err := CheckIPAddress(checkIP[0])
		if err != nil {
			cmdinfo.Status = "error: tcp address invalid"
			return cmdinfo
		}
	}
	topicname := ws[3]
	data := strings.Join(ws[4:], " ")
	if strings.HasPrefix(selectOption, "pub") {
		if data == "" {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		err := RunMqttPublish(tcpaddr, topicname, data)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: " + err.Error()
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	if strings.HasPrefix(selectOption, "sub") {
		err := RunMqttSubscribe(tcpaddr, topicname)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: " + err.Error()
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	if strings.HasPrefix(selectOption, "unsub") {
		err := RunMqttUnSubscribe(tcpaddr, topicname)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: " + err.Error()
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	cmdinfo.Status = "error: not command pub/sub"
	return cmdinfo
}
