package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"

	"github.com/qeof/q"
)

var sim_cidr string
var sim_mac string

func TestFirmwareUpgrade(t *testing.T) {
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
	sim_cidr = "192.168.4.1"
	sim_mac = "00-60-E9-28-01-01"

	time.Sleep(time.Second * 1)

	cmd := "scan gwd"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)

	go func() {
		HTTPMain()
	}()
	q.Q("wait for root to become ready...")
	if err := waitForRoot(); err != nil {
		t.Fatal(err)
	}

	myName := "test123"
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	testurl := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	ci := ClientInfo{Name: myName}
	jsonBytes, err := json.Marshal(ci)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	resp, err := PostWithToken(testurl, adminToken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		if resp.StatusCode != 200 {
			t.Fatalf("post StatusCode %d", resp.StatusCode)
		}
		q.Q(resp.Header)
		// save close, in resp != nil block
		resp.Body.Close()
	}

	adminToken, err := GetToken("admin")
	if err != nil {
		t.Fatalf("get token %v", err)
	}

	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))

	testurl = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = PostWithToken(testurl, adminToken,
		bytes.NewBuffer(jsonBytes))
	if resp == nil {
		t.Fatalf("post resp is nil")
	}
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("post %v", err)
	}
	// save close, already check resp != nil
	resp.Body.Close()

	testurl = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = GetWithToken(testurl, adminToken)
	if resp == nil {
		t.Fatalf("get resp is nil")
	}
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get %v", err)
	}
	// save close, already check resp != nil
	resp.Body.Close()

	// start file test
	var files []string = []string{
		"file:///testdata.json",
		"https://www.atoponline.com/wp-content/themes/atoponline/images/logo-new-thinned.svg",
		"https://www.atoponline.com/wp-content/uploads/2017/11/EHG750X-K770A770.zip",
	}
	ip := sim_cidr
	fileformat := ""
	for _, file := range files {
		fileformat = ""
		u, err := url.Parse(file)
		if err != nil {
			t.Fatal("error: url parse load error")
			continue
		}
		q.Q(u.Scheme, u.Path)
		if u.Scheme == "http" || u.Scheme == "https" {
			t.Log("Test http:// upgrade")
			fileformat = "http"
		} else if u.Scheme == "file" {
			t.Log("Test file:/// upgrade")
			fileformat = "file"
			file = strings.TrimPrefix(u.Path, "/")
		} else {
			t.Fatal("error: unknown file format")
			continue
		}
		// create new  device for firmware
		fs := FirmStatus{Status: ""}
		device := Firmware{ip: ip, firmStatus: fs}

		err = device.Upgrading(fileformat, file)
		if err != nil {
			t.Fatalf(err.Error())
		}

		for {
			time.Sleep(time.Duration(time.Second * 1))
			r, err := device.GetProcessStatus()

			if r == "Error" {
				t.Fatal("device:" + ip + ",Process:" + r + ",err:" + err.Error())
				break
			} else if r == "Complete" {
				t.Log(r)
				break
			}
		}
	}
}
