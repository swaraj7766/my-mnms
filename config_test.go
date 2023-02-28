package mnms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	"mnms/pkg/simulator/yaml"

	"github.com/qeof/q"
)

func init() {
	q.O = "stderr"
	q.P = ".*"
}

func TestConfigDevData(t *testing.T) {
	err := ReadTestData()
	if err != nil {
		t.Fatal(err)
	}

	dev, err := FindDev("02-42-C0-A8-64-A7")
	if err != nil {
		t.Fatal(err)
	}
	q.Q(dev)

	configData := `{ "ipaddress": "10.10.10.10","hostname":"fancycat"}`

	conf := make(map[string]string)
	err = json.Unmarshal([]byte(configData), &conf)
	if err != nil {
		t.Fatalf("can't unmarshal %v", err)
	}
}

// TestConfigSnmpBasic tests the basic snmp configuration
/*
  simulator yaml file:

	Environments:
  simulator_group: #create once device using macAddress:"00-60-e9-18-01-99"
      type: "EH7520"
      startPrefixip: "192.168.12.188/24"
      startMacAddress: "00-60-e9-18-01-99"

  run simulator:
	simulator.exe run -y config.yaml

	objID .1.3.6.1.4.1.3755.0.0.21.
	.1.3.6.1.4.1.3755.0.0.21.10.1.2.1.0

	acceptable settings:
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

*/

var testSnmpSimulators []*simulator.AtopGwdClient

func runTestSnmpSimulators(ctx context.Context) {
	go func() {
		name, err := net.GetDefaultInterfaceName()
		if err != nil {
			q.Q(err)
			log.Fatal(err)
		}

		// type: "EH7520"
		//   startPrefixip: "192.168.12.188/24"
		//   startMacAddress: "00-60-e9-18-01-99"
		env := map[string]yaml.Simulator{
			"simulator_group": {
				DeviceType:    "EH7520",
				Number:        1,
				StartPreFixIp: "192.168.12.188/24",
				MacAddress:    targetMac,
			},
		}

		newSims, err := yaml.NewSimulator(env, name)

		if err != nil {
			log.Fatal(err)
		}
		for _, v := range newSims {
			_ = v.StartUp()
			testSnmpSimulators = append(testSnmpSimulators, v)
		}
		q.Q(len(testSnmpSimulators))

		<-ctx.Done()
	}()
}

var targetMac = "00-60-E9-18-01-01"

// TestConfigSnmpBasic tests the basic snmp configuration such as syslog, trap server
func TestConfigSnmpBasic(t *testing.T) {
	q.P = ""
	//  ./simulator.exe run -y config.yaml
	// mnmsctl scan gwd
	// mnmsctl config syslog  00-60-E9-18-01-99 status 2
	ctx, cancel := context.WithCancel(context.Background())
	// run simulator and wait it start
	runTestSnmpSimulators(ctx)
	defer func() {
		for _, v := range testSnmpSimulators {
			_ = v.Shutdown()
		}
	}()
	time.Sleep(time.Millisecond * 500)
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
	defer cancel()

	// t.Log("mnmsctl is running, waitting 5 secound for scan")
	time.Sleep(time.Second * 5)
	// check mnmsctl is running

	// http request to localhost:27182/api
	resp, err := http.Get("http://localhost:27182/api")
	if err != nil {
		t.Fatal("mnmsctl is not running, should run mnmsctl first")
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("mnms haven't /api endpoint: ", err)
	}
	if string(respBody) != "mnms says hello" {
		t.Fatal("expected mnms says hello, got ", respBody)
	}

	// Run simulator check simulator device exist
	queryUrl := fmt.Sprintf("http://localhost:27182/api/v1/devices?dev=%s", targetMac)

	adminToken, err := GetToken("admin")
	if err != nil {
		t.Fatal("can't get admin token: ", err)
	}
	resp1, err := GetWithToken(queryUrl, adminToken)
	if resp1.StatusCode != 200 || err != nil {
		t.Fatal("can't get device info: ", err)
	}
	defer resp1.Body.Close()
	var devinfo DevInfo
	err = json.NewDecoder(resp1.Body).Decode(&devinfo)
	if err != nil {
		t.Fatal("marshal devinfo fail ", err)
	}

	// Testing syslog
	syslogSettings := map[string]string{
		"status":       "2",
		"server-ip":    "123.122.121.120",
		"server-port":  "123",
		"server-level": "2",
		"LogToFlash":   "2",
	}

	syslogFields := map[string]string{
		"status":       ".10.1.2.1.0:Integer",
		"server-ip":    ".10.1.2.6.0:OctetString",
		"server-port":  ".10.1.2.3.0:OctetString",
		"server-level": ".10.1.2.4.0:Integer",
		"LogToFlash":   ".10.1.2.5.0:Integer",
	}

	for f, s := range syslogSettings {
		// Don't call t.Fatal, t.Error, t.Log, etc. inside go routine
		// because test context 't' may not be there after test is done.
		// If 't' is done, go routine will crash, referring to 't' which
		// is no longer there.
		go func(field string, setting string) {

			cmdinfo := make(map[string]CmdInfo)
			cmd := fmt.Sprintf("config syslog  00-60-E9-18-01-99 %s %s", field, setting)

			insertcmd(cmd, &cmdinfo)
			jsonBytes, err := json.Marshal(cmdinfo)
			if err != nil {
				q.Q("json marshal", err)
				return
			}

			resp, err := http.Post("http://localhost:27182/api/v1/commands", "application/text",
				bytes.NewBuffer([]byte(jsonBytes)))
			if err != nil {
				q.Q("config syslog post error", err)
			}
			if resp != nil && resp.StatusCode != 200 {
				q.Q("config syslog post status code", resp.StatusCode)
			}
			defer resp1.Body.Close()
			// read back
			// TODO: implement read config API
			objID := ".1.3.6.1.4.1.3755.0.0.21."
			oidNType := strings.Split(syslogFields[field], ":")
			oid := objID + oidNType[0]

			res, err := SnmpGet(devinfo.IPAddress, []string{oid})
			// res, err := params.Get(oids)
			if err != nil {
				q.Q("snmp get fail", err)
				return
			}
			if len(res.Variables) < 1 {
				q.Q("expect 1 result but got", len(res.Variables))
				return
			}
			value := PDUToString(res.Variables[0])
			q.Q("expect %s  got %s", setting, value)
			if value != setting {
				q.Q("expect %s but got", setting, value)
			}

			return
		}(f, s)
	}
	trapServerFields := map[string]string{
		"status":      ".8.6.1.5.0:Integer",
		"server-ip":   ".8.6.1.7.0:OctetString",
		"server-port": ".8.6.1.6.0:Integer",
		"community":   ".8.6.1.3.0:OctetString",
	}

	trapSettings := map[string]string{
		"status":      "2",
		"server-ip":   "123.122.121.120",
		"server-port": "123",
		"community":   "test-community",
	}

	for f, s := range trapSettings {
		// Don't call t.Fatal, t.Error, t.Log, etc. inside go routine
		// because test context 't' may not be there after test is done.
		// If 't' is done, go routine will crash, referring to 't' which
		// is no longer there.
		go func(field string, setting string) {

			cmdinfo := make(map[string]CmdInfo)
			cmd := fmt.Sprintf("config snmptrap  00-60-E9-18-01-99 %s %s", field, setting)
			insertcmd(cmd, &cmdinfo)
			jsonBytes, err := json.Marshal(cmdinfo)
			if err != nil {
				q.Q("json marshal", err)
				return
			}

			resp, err := http.Post("http://localhost:27182/api/v1/commands", "application/text",
				bytes.NewBuffer([]byte(jsonBytes)))
			if err != nil {
				q.Q("config snmptrap post fail", err)
			}
			if resp.StatusCode != 200 {
				q.Q("config snmptrap post fail", field, resp.StatusCode)
			}
			defer resp1.Body.Close()
			// read back
			// TODO: implement read config API
			objID := ".1.3.6.1.4.1.3755.0.0.21."
			oidNType := strings.Split(trapServerFields[field], ":")
			oid := objID + oidNType[0]

			res, err := SnmpGet(devinfo.IPAddress, []string{oid})
			// res, err := params.Get(oids)
			if err != nil {
				q.Q("snmp get fail", err)
				return
			}
			if len(res.Variables) < 1 {
				q.Q("expect 1 result but got", len(res.Variables))
				return
			}
			value := PDUToString(res.Variables[0])
			t.Logf("expect %s  got %s", setting, value)
			if value != setting {
				q.Q("expect but got", setting, value)
			}

			return
		}(f, s)
	}
}
