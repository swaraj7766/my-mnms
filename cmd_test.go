package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/qeof/q"
)

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

	//cmdinfo.All = true
	//insertcmd("scan gwd", &cmdinfo)

	insertcmd("snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0", &cmdinfo)

	insertcmd("snmp set 192.168.100.53 1.3.6.1.2.1.1.5.0 aaa123 OctetString", &cmdinfo)

	insertcmd("switch 00-60-E9-2D-91-3E admin default show ip", &cmdinfo)
	insertcmd("snmp options 162 test 2c 5", &cmdinfo)

	//alternatively you can just fill out the key for the cmdinfo map and assign empty struct.
	//the cmdinfo.Command and timestamp will be filled out for you.
	emptyCmd := CmdInfo{}

	cmdinfo["mtderase 00-11-22-33-44-55 1.1.1.1 user1 pass1"] = emptyCmd
	cmdinfo["reset 00-11-22-33-44-55 1.1.1.1 user1 pass1"] = emptyCmd
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
	cmdinfo["snmp options 162 test 2c 5"] = emptyCmd
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

	cmdJson := `{"scan snmp":{"all":true,"timestamp":"2023-01-12T22:55:36-05:00","command":"scan snmp","result":"","status":"","name":""},"beep 00-11-22-33-44-55 1.1.1.1":{},"config hostname 00-11-22-33-44-55 uniquehostname":{},"config ipaddr 00-11-22-33-44-55 2.2.2.2 255.255.255.0 2.2.2.1":{},"config snmptrap 00-11-22-33-44-55 community test-community":{},"config snmptrap 00-11-22-33-44-55 server-ip 1.1.1.1":{},"config snmptrap 00-11-22-33-44-55 server-port 123":{},"config snmptrap 00-11-22-33-44-55 status 2":{},"config syslog 00-11-22-33-44-55 LogToFlash 2":{},"config syslog 00-11-22-33-44-55 server-level 2":{},"config syslog 00-11-22-33-44-55 server-port 123":{},"config syslog 00-11-22-33-44-55 status 2":{},"config user 00-11-22-33-44-55 user1 pass1":{},"reset 00-11-22-33-44-55 1.1.1.1 user1 pass1":{},"snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0":{"timestamp":"2023-01-12T22:55:36-05:00","command":"snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0","result":"","status":"","name":""},"snmp set 192.168.100.53 1.3.6.1.2.1.1.5.0 aaa123 OctetString":{"timestamp":"2023-01-12T22:55:36-05:00","command":"snmp set 192.168.100.53 1.3.6.1.2.1.1.5.0 aaa123 OctetString","result":"","status":"","name":""},"switch 00-60-E9-2D-91-3E admin default show ip":{"timestamp":"2023-01-12T22:55:36-05:00","command":"switch 00-60-E9-2D-91-3E admin default show ip","result":"","status":"","name":""},"snmp options 162 test 2c 5": {"timestamp": "2023-01-12T22:55:36-05:00","command": "snmp options 162 test 2c 5","result": "","status": "","name": ""}}`

	//simplest cmdJson could be: `{"scan snmp":{"all":true}}` or `{"beep 00-11-22-33-44-55 2.2.2.2":{}}`

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
	q.Q("wait for root to become ready...")
	if err := waitForRoot(); err != nil {
		t.Fatal(err)
	}

	rooturl := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	ci := ClientInfo{Name: myName}
	jsonBytes, err = json.Marshal(ci)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	resp, err := PostWithToken(rooturl, admintoken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		if resp.StatusCode != 200 {
			t.Fatalf("post StatusCode %d", resp.StatusCode)
		}
		q.Q(resp.Header)

		//save close
		resp.Body.Close()
	}

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)

	resp, err = PostWithToken(rooturl, admintoken, bytes.NewBuffer([]byte(cmdJson)))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		t.Log(resp.Header)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("error: reading resp body %v", err)
		}
		t.Log(string(body))
		//save close
		resp.Body.Close()
		if resp.StatusCode != 200 {
			t.Fatalf("post bad status %v", resp)
		}
	}

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
		q.Q("retrieved all cmd info", cmdinfo)
		//save close
		resp.Body.Close()
	}

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=%s", QC.Port, url.QueryEscape("config user 00-11-22-33-44-55 user1 pass1"))
	resp, err = GetWithToken(rooturl, admintoken)

	if resp == nil {
		t.Fatalf("nil response")
	}
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}

	cmdinfo = make(map[string]CmdInfo)
	err = json.NewDecoder(resp.Body).Decode(&cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	q.Q("config user", cmdinfo)
	//save close
	resp.Body.Close()

	// check snmp options
	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?cmd=%s", QC.Port, url.QueryEscape("snmp options 162 test 2c 5"))
	resp, err = GetWithToken(rooturl, admintoken)

	if resp == nil {
		t.Fatalf("nil response")
	}
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}

	cmdinfo = make(map[string]CmdInfo)
	err = json.NewDecoder(resp.Body).Decode(&cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	q.Q("snmp options", cmdinfo)
	//save close
	resp.Body.Close()

	rooturl = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, url.QueryEscape(myName))
	resp, err = GetWithToken(rooturl, admintoken)

	if err != nil || resp == nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}

	cmdinfo = make(map[string]CmdInfo)
	err = json.NewDecoder(resp.Body).Decode(&cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	q.Q("id", cmdinfo)
	// save close
	resp.Body.Close()
}

func TestCmdSwitchCli(t *testing.T) {
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
