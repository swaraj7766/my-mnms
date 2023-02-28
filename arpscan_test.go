package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"

	"github.com/qeof/q"
)

func init() {
	q.O = "stderr"
	q.P = ".*"
}

func TestArpscanCheckAllDevicesAlive(t *testing.T) {
	// run simultors
	simmap := map[string]atopyaml.Simulator{}
	simmap["group1"] = atopyaml.Simulator{Number: 1, DeviceType: "EH7508", StartPreFixIp: "192.168.4.1/24", MacAddress: "00-60-E9-28-01-01"}
	name, err := net.GetDefaultInterfaceName()
	if err != nil {
		t.Fatalf("get default interface name %v", err)
	}
	simulators, err := atopyaml.NewSimulator(simmap, name)
	if err != nil {
		t.Fatalf("new simulator %v", err)
	}
	for _, sim := range simulators {
		_ = sim.StartUp()
		defer func(sim *simulator.AtopGwdClient) {
			_ = sim.Shutdown()
		}(sim)
	}

	time.Sleep(time.Second * 1)

	go func() {
		GwdMain()
	}()
	go func() {
		HTTPMain()
	}()

	time.Sleep(1 * time.Second)

	// start test
	err = ArpScan()
	if err != nil {
		t.Fatal(err)
	}
}

func TestArpscanInterval(t *testing.T) {

	cmd := "arp interval 10"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)

	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))

	go func() {
		HTTPMain()
	}()

	myName := "testarpcmd"

	// init
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	var testurl string
	var resp *http.Response

	/*
		testurl := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
		resp, err := PostWithToken(testurl, bytes.NewBuffer([]byte(myName)))
		if err != nil || resp.StatusCode != 200 {
			t.Fatalf("post %v", err)
		}
		resp.Body.Close()
	*/

	testurl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(testurl, token, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	resp.Body.Close()

	testurl = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = GetWithToken(testurl, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	resp.Body.Close()

	err = CheckCommands()
	if err != nil {
		q.Q(err)
	}

	if QC.ArpInterval != 10 {
		t.Errorf("Error : set Arp Interval error.")
	}

	// test interval 0, should not set arp interval
	cmd = "arp interval 0"
	cmdinfo = make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	testurl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(testurl, token, bytes.NewBuffer(jsonBytes))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	resp.Body.Close()
	testurl = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = GetWithToken(testurl, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	resp.Body.Close()
	err = CheckCommands()
	if err != nil {
		q.Q(err)
	}
	if QC.ArpInterval == 0 {
		t.Fatal("command interval should not be set to 0")
	}
}
