package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/qeof/q"
	"github.com/ziutek/telnet"
)

type CmdInfo struct {
	Timestamp string `json:"timestamp"`
	Command   string `json:"command"`
	Result    string `json:"result"`
	Status    string `json:"status"`
	Name      string `json:"name"`
	Retries   int    `json:"retries"`
}

const telnet_timeout = 10 * time.Second

func init() {
	QC.CmdData = make(map[string]CmdInfo)
}

func InsertCli(cmd string) {
	// command line support to allow mnmsctl cli commands.
	// cli commands are same as api calls.
	q.Q(cmd)
	// special handling for commands that must run on all connected
	// clients.  root may have multiple clients connected.
	// some commands like scan snmp must be issued to all of them.
	if strings.HasPrefix(cmd, "all") {
		ws := strings.Split(cmd, " ")
		if len(ws) < 2 {
			q.Q("error: invalid command")
			return
		}
		cmd = strings.Join(ws[1:], " ")
		// horrible hack to have one-shot slot for commands going to all
		QC.ClientMutex.Lock()
		defer QC.ClientMutex.Unlock()
		for _, v := range QC.Clients {
			if v != "" {
				q.Q("error: one shot filled")
				return
			}
		}
		//TODO: mechanism to forcibly empty the one shot per client
		for k, _ := range QC.Clients {
			QC.Clients[k] = cmd
		}
		q.Q(QC.Clients)
		return
	}

	//normal commands get queued
	cmdinfo := CmdInfo{
		Timestamp: time.Now().Format(time.RFC3339),
		Command:   cmd,
	}
	QC.CmdMutex.Lock()
	QC.CmdData[cmd] = cmdinfo
	QC.CmdMutex.Unlock()
}

func InsertCommands(cmddata *map[string]CmdInfo) {
	// these are commands that are downloaded from root probably.
	// insert into local queue
	for k, v := range *cmddata {
		QC.CmdMutex.Lock()
		_, ok := QC.CmdData[k]
		QC.CmdMutex.Unlock()
		if ok {
			q.Q("cmd exists already", v)
			continue
		}
		v.Name = QC.Name
		QC.CmdMutex.Lock()
		QC.CmdData[k] = v
		QC.CmdMutex.Unlock()
		q.Q("inserted cmd", v)
	}
}

func UpdateCommands(cmddata *map[string]CmdInfo) {
	q.Q(cmddata)
	// if root is collecting command status from clients,
	// updates from each client will be aggregated.
	// update command status in the queue
	for k, v := range *cmddata {
		if v.Status == "" {
			if v.Command == "" {
				v.Command = k
			}
			InsertCli(v.Command)
			continue
		}

		QC.CmdMutex.Lock()
		found, ok := QC.CmdData[k]
		QC.CmdMutex.Unlock()

		// don't override ok result command history
		if ok && found.Status == "ok" {
			continue
		}
		// fill in missing timestamp
		if v.Timestamp == "" {
			v.Timestamp = time.Now().Format(time.RFC3339)
		}
		QC.CmdMutex.Lock()
		QC.CmdData[k] = v
		QC.CmdMutex.Unlock()

		q.Q(found, v)
	}
}

func CheckCommands() error {

	// this runs periodically to download commands, run and update results
	if QC.Root != "" {
		// if there is a URL to the root we download from it
		resp, err := GetWithToken(QC.Root+"/api/v1/commands?id="+QC.Name, QC.AdminToken)
		if err != nil {
			return err
		}
		if resp != nil {
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
				InsertCommands(&cmddata)
			}
		}
	}

	//XXX this mutex lockout can be very long
	QC.CmdMutex.Lock()
	for k, v := range QC.CmdData {
		if v.Status != "" &&
			!strings.HasPrefix(v.Status, "pending:") {
			continue
		}
		q.Q("run cmd", v)
		// XXX cannot go RunCmd() because of ordering
		res := RunCmd(&v)
		QC.CmdData[k] = v
		q.Q(res)
	}
	QC.CmdMutex.Unlock()

	if QC.Root != "" { //always check for root URL to run even when no root
		// update results back to root
		QC.CmdMutex.Lock()
		jsonBytes, err := json.Marshal(QC.CmdData)
		QC.CmdMutex.Unlock()

		if err != nil {
			return err
		}
		resp, err := PostWithToken(QC.Root+"/api/v1/commands", QC.AdminToken, bytes.NewBuffer(jsonBytes))
		if err != nil {
			return err
		}
		//q.Q("updated commands", QC.CmdData)
		if resp != nil {
			q.Q(resp.Header)
		}
		resp.Body.Close()
	}
	return nil
}

func RunCmd(cmdinfo *CmdInfo) *CmdInfo {
	//main cmd run dispatch

	// commands == api

	// basic commands: reset, reboot, beep
	// scan commands
	// config commands -- "basic" config commands
	// snmp commands
	// log commands -- debugging
	// devices commands

	//TODO group command handling -- sending to multiple devices
	q.Q(cmdinfo)
	cmd := cmdinfo.Command
	if strings.HasPrefix(cmdinfo.Status, "pending:") {
		cmdinfo.Retries++
	}

	if strings.HasPrefix(cmd, "reset") {
		return ResetCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "beep") {
		return BeepCmd(cmdinfo)
	}
	if strings.HasPrefix(cmd, "reboot") {
		return RebootCmd(cmdinfo)
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

	if strings.HasPrefix(cmd, "devices") {
		return DevicesCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "log") {
		return LogCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "firmware") {
		return FirmwareCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "command") {
		return CommandCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "arp") {
		return ArpCmd(cmdinfo)
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
	_, err := FindDev(macaddr)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	err = GwdBeep(ipaddr, macaddr)
	if err != nil {
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

func RebootCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var macaddr, ipaddr, username, password string
	Unpack(ws[1:], &macaddr, &ipaddr, &username, &password)
	_, err := FindDev(macaddr)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	err = GwdReboot(ipaddr, macaddr, username, password)
	if err != nil {
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

func ResetCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var macaddr, ipaddr, username, password string
	Unpack(ws[1:], &macaddr, &ipaddr, &username, &password)
	_, err := FindDev(macaddr)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	err = GwdReset(ipaddr, macaddr, username, password)
	if err != nil {
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}

func CheckSwitchCliModel(modelname string) bool {
	if strings.HasPrefix(modelname, "EH7") || strings.HasPrefix(modelname, "EHG7") {
		q.Q(modelname)
		return true
	}
	q.Q("switch cli not supported", modelname)
	return false
}

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
	username := ws[2]
	password := ws[3]

	err = SendSwitch(cmdinfo, dev, username, password, strings.Join(ws[4:], " "))
	if err != nil {
		cmdinfo.Status = err.Error()
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
		q.Q(err)
	}
	cmdinfo.Result = string(result)

	return nil
}

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

func ValidateCommands() error {
	q.Q("validate cmds")
	for k, v := range QC.CmdData {
		if v.Retries > 3 {
			q.Q("cancel cmd", v)
			v.Status = "cancelled: " + v.Status
			QC.CmdData[k] = v
		}
	}
	return nil
}

// mqtt pub {topic name} {messages}
// mqtt sub {topic name} {timeout.Second}
func RunMqttCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	selectOption := ws[1]
	topicname := ws[2]
	data := strings.Join(ws[3:], " ")
	if strings.HasPrefix(selectOption, "pub") {
		err := RunMqttPublish(topicname, data)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = err.Error()
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	if strings.HasPrefix(selectOption, "sub") {
		v, _ := strconv.Atoi(data)
		err := RunMqttSubscribe(topicname, v)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = err.Error()
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	cmdinfo.Status = "error: not command pub/sub"
	return cmdinfo
}

func isRootCommand(cmd string) bool {
	if strings.HasPrefix(cmd, "config crontab") {
		return true
	}
	if strings.HasPrefix(cmd, "devices save") {
		return true
	}

	if strings.HasPrefix(cmd, "rsa") {
		return true
	}
	if strings.HasPrefix(cmd, "mnmsconfig") {
		return true
	}
	if strings.HasPrefix(cmd, "mqtt") {
		return true
	}

	if strings.HasPrefix(cmd, "devices load") {
		return true
	}
	if strings.HasPrefix(cmd, "devices files list") {

		return true
	}
	return false
}

func CommandCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command

	if strings.HasPrefix(cmd, "command delete") {
		return CommandDeleteCmd(cmdinfo)
	}

	if strings.HasPrefix(cmd, "command interval") {
		return CommandIntervalCmd(cmdinfo)
	}

	cmdinfo.Status = "error: invalid command"
	return cmdinfo
}

// CommandIntervalCmd set CheckCommand() interval in second(s), ex: command interval 10.
// range 1-3600 seconds.
func CommandIntervalCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command

	ws := strings.Split(cmd, " ")
	if len(ws) != 3 {
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
	QC.CmdInterval = interval
	cmdinfo.Result = fmt.Sprintf("set interval to %v seconds", interval)
	cmdinfo.Status = "ok"
	return cmdinfo
}

// CommandDeleteCmd delete command from queue, ex: command delete switch mac user pass show info.
// Use "all delete command" to issue all clients.
// Also "all delete command" can't be used with other "all" commands in the same time because the QC.Client stores one command at a time.
func CommandDeleteCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	// ex: command delete switch mac user pass show info
	// ex: command delete beep mac ip
	if len(ws) < 3 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	key := strings.Join(ws[2:], " ")
	q.Q("delete cmd", key)
	if _, ok := QC.CmdData[key]; !ok {
		cmdinfo.Status = "error: not found command"
		return cmdinfo
	}
	delete(QC.CmdData, key)
	cmdinfo.Status = "ok"
	return cmdinfo
}
