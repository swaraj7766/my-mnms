package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"

	"github.com/qeof/q"
	"github.com/sirupsen/logrus"
)

var simulatorPath = "simulator.yaml"

var totalsimulator []*simulator.AtopGwdClient

// pelase run as administrator or root
//
// createSimualtor simulator base on config.yaml
func createSimulatorFile() error {
	startup := make(chan bool)
	name, err := net.GetDefaultInterfaceName()
	if err != nil {
		logrus.Fatal(err)
	}
	simulators, err := atopyaml.NewSimulatorFile(simulatorPath, name)
	if err != nil {
		q.Q(err)
		return err
	}
	go func() {

		for _, v := range simulators {
			_ = v.StartUp()
			totalsimulator = append(totalsimulator, v)
		}
		startup <- true
	}()

	q.Q("simulator number:", len(simulators))
	time.Sleep(time.Millisecond * 300)
	<-startup
	return nil
}

func createSimulator() error {
	startup := make(chan bool)
	name, err := net.GetDefaultInterfaceName()
	if err != nil {
		logrus.Fatal(err)
	}
	simmap := map[string]atopyaml.Simulator{}
	simmap["group1"] = atopyaml.Simulator{Number: 5, DeviceType: "EH7506", StartPreFixIp: "192.168.10.1/24", MacAddress: "00-60-E9-18-01-01"}
	simmap["group2"] = atopyaml.Simulator{Number: 5, DeviceType: "EH7508", StartPreFixIp: "192.168.20.1/24", MacAddress: "00-60-E9-18-0A-01"}
	simulators, err := atopyaml.NewSimulator(simmap, name)
	if err != nil {
		q.Q((err))
		return err
	}
	go func() {
		for _, v := range simulators {
			_ = v.StartUp()
			totalsimulator = append(totalsimulator, v)
		}
		startup <- true
	}()

	time.Sleep(time.Millisecond * 300)
	<-startup
	q.Q("simulator number:", len(simulators))
	return nil
}

// ShutdownSimulator Shutdown simualtor
func ShutdownSimulator() {
	for _, v := range totalsimulator {
		_ = v.Shutdown()
	}
}

// TestExample  get simulator value of oid 1.3.6.1.2.1.1.5.0
//
// simulator detail follow simulator.yaml
func TestSimFileExample(t *testing.T) {
	err := createSimulatorFile()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()

	cmdinfo := make(map[string]CmdInfo)
	insertcmd("snmp get 192.168.10.1 1.3.6.1.2.1.1.5.0", &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		HTTPMain()
	}()

	myName := "test123"

	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	resp, err := http.Post(url,
		"application/text",
		bytes.NewBuffer([]byte(myName)))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	url = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = http.Post(url,
		"application/json",
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	url = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = http.Get(url)

	if err != nil {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo := make(map[string]CmdInfo)
		_ = json.NewDecoder(resp.Body).Decode(&cmdinfo)

		q.Q(cmdinfo)
	}
	resp.Body.Close()

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCommands()
}

// TestSimExample  get simulator value of oid 1.3.6.1.2.1.1.5.0
func TestSimExample(t *testing.T) {
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()

	cmdinfo := make(map[string]CmdInfo)
	insertcmd("snmp get 192.168.10.1 1.3.6.1.2.1.1.5.0", &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		HTTPMain()
	}()

	myName := "test123"

	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	resp, err := http.Post(url,
		"application/text",
		bytes.NewBuffer([]byte(myName)))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	url = fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err = http.Post(url,
		"application/json",
		bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()

	url = fmt.Sprintf("http://localhost:%d/api/v1/commands?id=%s", QC.Port, myName)
	resp, err = http.Get(url)

	if err != nil {
		t.Fatalf("get %v", err)
	}
	if resp != nil {
		cmdinfo := make(map[string]CmdInfo)
		_ = json.NewDecoder(resp.Body).Decode(&cmdinfo)
		q.Q(cmdinfo)
	}
	resp.Body.Close()

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCommands()
}
