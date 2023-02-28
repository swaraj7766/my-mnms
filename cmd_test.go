package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/qeof/q"
)

func init() {
	q.O = "stderr"
	q.P = ""
}

func insertcmd(cmd string, cmdinfo *map[string]CmdInfo) {
	ci := CmdInfo{
		Timestamp: time.Now().Format(time.RFC3339),
		Command:   cmd,
	}
	(*cmdinfo)[cmd] = ci
}

func TestCmdInfo(t *testing.T) {
	cmd := "all scan snmp"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)

	// The special debug commands that start with scan and devices are
	// for all clients. They should be prepended with all when issued
	// at root. When debugging directly at client, scan and devices commands
	// can be issued directly to a client. Presence of "all" triggers
	// root to set single shot per client to be fetched.

	//insertcmd("all scan gwd", &cmdinfo)
	//insertcmd("all devices save", &cmdinfo)

	insertcmd("snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0", &cmdinfo)

	insertcmd("snmp set 192.168.100.53 1.3.6.1.2.1.1.5.0 aaa123 OctetString", &cmdinfo)

	insertcmd("switch 00-60-E9-2D-91-3E admin default show ip", &cmdinfo)

	//alternatively you can just fill out the key for the cmdinfo map and assign empty struct.
	//the cmdinfo.Command and timestamp will be filled out for you.
	emptyCmd := CmdInfo{}

	cmdinfo["reset 00-11-22-33-44-55 1.1.1.1 user1 pass1"] = emptyCmd
	cmdinfo["reboot 00-11-22-33-44-55 1.1.1.1 user1 pass1"] = emptyCmd
	cmdinfo["beep 00-11-22-33-44-55 1.1.1.1"] = emptyCmd

	//See config_test.go
	cmdinfo["config ipaddr 00-11-22-33-44-55 2.2.2.2 255.255.255.0 2.2.2.1"] = emptyCmd
	cmdinfo["config hostname 00-11-22-33-44-55 uniquehostname"] = emptyCmd

	cmdinfo["config syslog 00-11-22-33-44-55 status 2"] = emptyCmd
	cmdinfo["config syslog 00-11-22-33-44-55 server-port 123"] = emptyCmd
	cmdinfo["config syslog 00-11-22-33-44-55 server-level 2"] = emptyCmd
	cmdinfo["config syslog 00-11-22-33-44-55 LogToFlash 2"] = emptyCmd

	cmdinfo["config snmptrap 00-11-22-33-44-55 status 2"] = emptyCmd
	cmdinfo["config snmptrap 00-11-22-33-44-55 server-ip 1.1.1.1"] = emptyCmd
	cmdinfo["config snmptrap 00-11-22-33-44-55 server-port 123"] = emptyCmd
	cmdinfo["config snmptrap 00-11-22-33-44-55 community test-community"] = emptyCmd

	var jsonBytes []byte
	var err error
	/*
		jsonBytes, err = json.Marshal(cmdinfo)
		if err != nil {
			t.Fatalf("json marshal %v", err)
		}
		q.Q(string(jsonBytes))
	*/

	// The equivalent json string for the above

	cmdJson := `{"all scan snmp":{"timestamp":"2023-01-12T22:55:36-05:00","command":"all scan snmp","result":"","status":"","name":""},"beep 00-11-22-33-44-55 1.1.1.1":{},"config hostname 00-11-22-33-44-55 uniquehostname":{},"config ipaddr 00-11-22-33-44-55 2.2.2.2 255.255.255.0 2.2.2.1":{},"config snmptrap 00-11-22-33-44-55 community test-community":{},"config snmptrap 00-11-22-33-44-55 server-ip 1.1.1.1":{},"config snmptrap 00-11-22-33-44-55 server-port 123":{},"config snmptrap 00-11-22-33-44-55 status 2":{},"config syslog 00-11-22-33-44-55 LogToFlash 2":{},"config syslog 00-11-22-33-44-55 server-level 2":{},"config syslog 00-11-22-33-44-55 server-port 123":{},"config syslog 00-11-22-33-44-55 status 2":{},"config user 00-11-22-33-44-55 user1 pass1":{},"reboot 00-11-22-33-44-55 1.1.1.1 user1 pass1":{},"reset 00-11-22-33-44-55 1.1.1.1 user1 pass1":{},"snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0":{"timestamp":"2023-01-12T22:55:36-05:00","command":"snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0","result":"","status":"","name":""},"snmp set 192.168.100.53 1.3.6.1.2.1.1.5.0 aaa123 OctetString":{"timestamp":"2023-01-12T22:55:36-05:00","command":"snmp set 192.168.100.53 1.3.6.1.2.1.1.5.0 aaa123 OctetString","result":"","status":"","name":""},"switch 00-60-E9-2D-91-3E admin default show ip":{"timestamp":"2023-01-12T22:55:36-05:00","command":"switch 00-60-E9-2D-91-3E admin default show ip","result":"","status":"","name":""}}`

	//simplest cmdJson could be: `{"all scan snmp":{}}` or `{"beep 00-11-22-33-44-55 2.2.2.2":{}}`

	jsonBytes = []byte(cmdJson)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		HTTPMain()
	}()

	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	myName := "test123"
	admintoken, err := GetToken("admin")
	if err != nil {
		t.Fatalf("get token %v", err)
	}
	time.Sleep(1 * time.Second)
	rooturl := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	// resp, err := http.Post(rooturl,
	// 	"application/text",
	// 	bytes.NewBuffer([]byte(myName)))
	// ? The original request's type is "application/text", but PostWithToken() is "application/json"
	// ? It pass the test, but we should check it
	resp, err := PostWithToken(rooturl, admintoken, bytes.NewBuffer([]byte(myName)))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)

	resp, err = PostWithToken(rooturl, admintoken, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=all", QC.Port)
	resp, err = GetWithToken(rooturl, admintoken)

	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&cmdinfo)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("all", cmdinfo)
	}
	resp.Body.Close()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=%s", QC.Port, url.QueryEscape("config user 00-11-22-33-44-55 user1 pass1"))
	resp, err = GetWithToken(rooturl, admintoken)

	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&cmdinfo)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("config user", cmdinfo)
	}
	resp.Body.Close()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, url.QueryEscape(myName))
	resp, err = GetWithToken(rooturl, admintoken)

	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&cmdinfo)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("id", cmdinfo)
	}
	resp.Body.Close()
}

func TestSwitchCli(t *testing.T) {
	hostname, err := os.Hostname()
	if err != nil {
		t.Fatal()
	}
	if hostname == "testbed" {
		err := ReadTestData()
		if err != nil {
			t.Fatal(err)
		}
		devId := "00-60-E9-2D-91-3E"
		dev, err := FindDev(devId)
		if err != nil {
			t.Fatal(err)
		}
		if !CheckSwitchCliModel(dev.ModelName) {
			t.Fatalf("model %s does not support switch cli\n", dev.ModelName)
		}
		cmdinfo := CmdInfo{Command: "admin default show ip"}
		_ = SendSwitch(&cmdinfo, dev, "admin", "default", "show ip")
		q.Q(cmdinfo)
		cmdinfo.Command = "admin default no snmp"
		_ = SendSwitch(&cmdinfo, dev, "admin", "default", "no snmp")
		q.Q(cmdinfo)
		cmdinfo.Command = "admin default show info"
		_ = SendSwitch(&cmdinfo, dev, "admin", "default", "show info")
		q.Q(cmdinfo)
		cmdinfo.Command = "admin default snmp"
		_ = SendSwitch(&cmdinfo, dev, "admin", "default", "snmp")
		q.Q(cmdinfo)
		cmdinfo.Command = "admin default show info"
		_ = SendSwitch(&cmdinfo, dev, "admin", "default", "show info")
		q.Q(cmdinfo)
		fmt.Printf("%s", cmdinfo.Result)
	}
}

func TestCommandDelete(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		HTTPMain()
	}()
	time.Sleep(time.Second * 1)

	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	QC.Root = "http://localhost:27182"
	QC.Name = "testdeletecmd"
	QC.Clients = make(map[string]string)
	QC.CmdData = make(map[string]CmdInfo)
	admintoken, err := GetToken("admin")
	if err != nil {
		t.Fatalf("get token %v", err)
	}
	QC.AdminToken = admintoken
	time.Sleep(1 * time.Second)
	rooturl := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	resp, err := PostWithToken(rooturl, admintoken, bytes.NewBuffer([]byte(QC.Name)))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	cmdinfo := make(map[string]CmdInfo)
	insertcmd("all command interval 4", &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(rooturl, admintoken, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	time.Sleep(5 * time.Second)
	_ = CheckCommands()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=all", QC.Port)
	resp, err = GetWithToken(rooturl, admintoken)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo = make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&cmdinfo)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("all queue", cmdinfo)
	}
	resp.Body.Close()

	cmdinfo = make(map[string]CmdInfo)
	insertcmd("all command delete command interval 4", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(rooturl, admintoken, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}

	time.Sleep(5 * time.Second)
	_ = CheckCommands()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=all", QC.Port)
	resp, err = GetWithToken(rooturl, admintoken)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&cmdinfo)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("queue after delete", cmdinfo)
		for k := range cmdinfo {
			if k == "command interval 4" {
				t.Fatal("command not deleted")
			}
		}
	}
	resp.Body.Close()

	// all scan gwd again
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("all command interval 4", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(rooturl, admintoken, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	time.Sleep(5 * time.Second)
	_ = CheckCommands()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=all", QC.Port)
	resp, err = GetWithToken(rooturl, admintoken)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo = make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&cmdinfo)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("add cmd again", cmdinfo)
		if _, ok := cmdinfo["command interval 4"]; !ok {
			t.Fatal("command not added")
		}
	}
	resp.Body.Close()
}

func TestCommandIntervalCmd(t *testing.T) {
	QC.IsRoot = true
	myName := "root"
	QC.Name = "root"
	q.P = "jwt"
	// run services
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		GwdMain()
	}()
	wg.Add(1)
	go func() {
		// root and non-root run http api service.
		// if both root and non-root runs on same machine (should not happen)
		// then whoever runs first will run http service (port conflict)
		defer wg.Done()
		HTTPMain()
		// TODO root to dump snapshots of devices, logs, commands
	}()

	// init
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	token, err := GetToken("admin")
	if err != nil {
		t.Fatalf("get token %v", err)
	}
	QC.AdminToken = token

	cmdinfo := make(map[string]CmdInfo)
	insertcmd("command interval 10", &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	url := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(url, token, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()
	url = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = GetWithToken(url, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	resp.Body.Close()

	time.Sleep(1 * time.Second)
	err = CheckCommands()
	if err != nil {
		q.Q(err)
	}

	if QC.CmdInterval != 10 {
		t.Fatal("command interval not set")
	}

	// test set command interval to 0, should not work
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("command interval 0", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	url = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(url, token, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()
	url = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = GetWithToken(url, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	resp.Body.Close()

	time.Sleep(1 * time.Second)
	err = CheckCommands()
	if err != nil {
		q.Q(err)
	}

	if QC.CmdInterval == 0 {
		t.Fatal("command interval should not be set to 0")
	}
}
