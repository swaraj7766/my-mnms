package mnms

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"mnms/pkg/simulator"
	simnet "mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bitfield/script"
	"github.com/qeof/q"
)

var adminToken string

func init() {
	flag.StringVar(&q.O, "O", "stderr", "debug log output")
	flag.StringVar(&q.P, "P", "", "debug log pattern")
	var err error
	adminToken, err = GetToken("admin")
	if err != nil {
		q.Q("cannot get admin token", err)
	}
}

/*
func startTestSyslogServer() {
	udpsock, err := net.ListenPacket("udp", ":60514") // run a test "remote" syslog server
	if err != nil {
		q.Q("cannot listen udp syslog service at 60514", err)
		return
	}
	defer udpsock.Close()
	buf := make([]byte, 1024*2)

	for {
		mlen, _, err := udpsock.ReadFrom(buf)
		if err != nil {
			q.Q("cannot read udp sock", err)
			return
		}
		q.Q("syslog input:", buf[:mlen])
	}
}*/

var simulatorList []*simulator.AtopGwdClient

func startSimulators() error {
	startup := make(chan bool)
	name, err := simnet.GetDefaultInterfaceName()
	if err != nil {
		return err
	}
	simmap := map[string]atopyaml.Simulator{}
	simmap["group1"] = atopyaml.Simulator{Number: 5, DeviceType: "EH7506", StartPreFixIp: "192.168.11.1/24", MacAddress: "00-60-E9-18-11-11"}
	simmap["group2"] = atopyaml.Simulator{Number: 5, DeviceType: "EH7508", StartPreFixIp: "191.168.22.1/24", MacAddress: "00-60-E9-18-0D-11"}
	simulators, err := atopyaml.NewSimulator(simmap, name)
	if err != nil {
		q.Q("cannot get new simulator", err)
		return err
	}
	go func() {
		for _, v := range simulators {
			_ = v.StartUp()
			simulatorList = append(simulatorList, v)
		}
		startup <- true
	}()

	time.Sleep(time.Millisecond * 300)
	<-startup
	q.Q("created simulators", len(simulators))
	return nil
}

func stopSimulators() {
	for _, v := range simulatorList {
		_ = v.Shutdown()
	}
	simulatorList = nil // ok to call stopSimulators() more than once
}

func TestCluster(t *testing.T) {
	defer func() {
		killNmsctlProcesses()
	}()

	//go startTestSyslogServer()

	go func() { //run the root instance
		cmd := exec.Command("./mnmsctl/mnmsctl", "-n", "root", "-R", "-O", "root.log")
		output, err := cmd.CombinedOutput()
		if err != nil {
			q.Q("error running mnmsctl root instance", err)
		}
		q.Q("exit root process", string(output))
	}()

	q.Q("wait for root to become ready...")
	if err := waitForRoot(); err != nil {
		t.Fatal(err)
	}

	clients := []string{"client1", "client2"}

	// run two fake clients
	for _, clientName := range clients {
		go func(clientName string) {
			q.Q("running a client node instance", clientName)
			cmd := exec.Command("./mnmsctl/mnmsctl", "-n", clientName, "-s", "-O", clientName+".log", "-r", "http://localhost:27182", "-nosyslog", "-notrap", "-nomqbr", "-fake")
			output, err := cmd.CombinedOutput()
			if err != nil {
				q.Q("error running a client node instance", clientName, err)
			}
			q.Q("exit client process", clientName, string(output))
		}(clientName)
	}

	// client1 has devices with mac address starting with 00-61
	// client2 has devices with mac address starting with 02-42

	q.Q("ran client 1 and 2")
	time.Sleep(2 * time.Second)

	url := "http://localhost:27182/api/v1/commands"
	cmdinfo := make(map[string]CmdInfo)

	// special "all" command should go to all clients
	// no more all command skip it
	//insertcmd("all scan gwd", &cmdinfo)

	// snmp commands go to ip addresses. both clients will
	// run them.
	insertcmd("snmp get 192.168.100.53 1.3.6.1.2.1.1.5.0", &cmdinfo)
	insertcmd("snmp get 192.168.100.128 1.3.6.1.2.1.1.5.0", &cmdinfo)

	// this command should be executed by client1
	insertcmd("switch 00-60-E9-2D-91-3E admin default show ip", &cmdinfo)

	// this command should be executed by client2
	insertcmd("switch 02-42-C0-A8-64-80 admin default show ip", &cmdinfo)

	// this command for non-existing device  will be retried and cancelled
	insertcmd("config syslog 01-11-22-33-44-55 1 10.10.10.1 5514 1 1", &cmdinfo)

	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}

	//resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	resp, err := PostWithToken(url, adminToken, bytes.NewBuffer(jsonBytes))

	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("no resp")
	}
	//t.Log(resp.Header)
	if resp.StatusCode != 200 {
		q.Q("post returns bad status code", resp.StatusCode)
		body, _ := ioutil.ReadAll(resp.Body)
		t.Log(string(body))
		t.Fatal(resp.StatusCode)
	}
	//safe closed
	resp.Body.Close()

	// let commands download to clients and run
	q.Q("sent commands to root... will fetch command status...")
	time.Sleep(6 * time.Second)

	// fetch all command status
	url = "http://localhost:27182/api/v1/commands?cmd=all"
	//resp, err = http.Get(url)
	resp, err = GetWithToken(url, adminToken)
	if err != nil {
		q.Q("error doing http get", err)
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("no resp")
	}
	if resp.StatusCode != 200 {
		q.Q("error status code from get", resp.StatusCode)
		t.Fatal(resp.StatusCode)
	}
	//t.Log(resp.StatusCode)
	commands := make(map[string]CmdInfo)
	err = json.NewDecoder(resp.Body).Decode(&commands)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	q.Q("command status retrieved", commands)
	// save close
	resp.Body.Close()

	//issue some new commands
	cmdinfo = make(map[string]CmdInfo)

	//this should be done by client1
	insertcmd("switch 00-60-E9-2D-91-3E admin default show info", &cmdinfo)

	//this should be done by client2
	insertcmd("switch 02-42-C0-A8-64-80 admin default show info", &cmdinfo)

	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	//resp, err = http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	resp, err = PostWithToken(url, adminToken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		if resp.StatusCode != 200 {
			t.Fatal()
		}
		//t.Log(resp.Header)
		// save close
		resp.Body.Close()
	}

	//wait for a while to let the commands run
	q.Q("issued switch show info commands")
	time.Sleep(20 * time.Second)

	url = "http://localhost:27182/api/v1/commands?cmd=all"
	//resp, err = http.Get(url)
	resp, err = GetWithToken(url, adminToken)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {
		if resp.StatusCode != 200 {
			t.Fatal()
		}
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("updated command status fetched", commands)
		//save close
		resp.Body.Close()
	}

	// wait a bit more to see if the pending commands get cancelled
	q.Q("wait for cancel of pending commands")
	time.Sleep(6 * time.Second)

	url = "http://localhost:27182/api/v1/commands?cmd=all"
	//resp, err = http.Get(url)
	resp, err = GetWithToken(url, adminToken)
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil {

		if resp.StatusCode != 200 {
			t.Fatal()
		}
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q("fetched command status again", commands)
		//save close
		resp.Body.Close()
	}

	_, err = getDevices()
	if err != nil {
		t.Fatal("cannot get devices", err)
	}

	//wait for devices upload
	q.Q("wait for devices to be published")
	time.Sleep(6 * time.Second)

	_, err = getDevices()
	if err != nil {
		t.Fatal("cannot get devices", err)
	}
	q.Q("starting simulators")
	err = startSimulators()
	if err != nil {
		t.Fatal(err)
	}
	defer stopSimulators()

	// run a real client to scan devices

	go func(clientName string) {
		q.Q("running a real client node instance", clientName)
		cmd := exec.Command("./mnmsctl/mnmsctl", "-n", clientName, "-s", "-O", clientName+".log", "-r", "http://localhost:27182")
		output, err := cmd.CombinedOutput()
		if err != nil {
			q.Q(err)
		}
		q.Q("exit real client process", clientName, string(output))
	}("clientreal")

	q.Q("started real client to scan for devices")
	time.Sleep(7 * time.Second)

	_, err = getDevices()
	if err != nil {
		t.Fatal("cannot get devices", err)
	}

	//intentionally stop simulators
	stopSimulators()
	q.Q("stopped simulators scan for devices again to see simulators not responding")
	//listInterfaces()
	time.Sleep(10 * time.Second)
	q.Q("send gwd invite to refresh devices list")
	_ = GwdInvite()
	time.Sleep(10 * time.Second)

	newDevs, err := getDevices()
	if err != nil {
		t.Fatal("cannot get devices", err)
	}

	var latest, dts int64
	sdev, ok := (*newDevs)[specialMac]
	if !ok {
		t.Fatal("cannot find special MAC device")
	}
	latest, err = strconv.ParseInt(sdev.Timestamp, 10, 64)
	if err != nil {
		q.Q(sdev.Timestamp, err)
		t.Fatal("strconv parseint fail", err)
	}

	for _, dev := range *newDevs {
		if strings.HasPrefix(dev.ModelName, "Simu_") {
			dts, err = strconv.ParseInt(dev.Timestamp, 10, 64)
			if err != nil {
				q.Q(dev.Timestamp, err)
				t.Fatal("strconv parseint fail", err)
			}
			diff := latest - dts
			q.Q(dev, diff, latest, dts)
			if diff < 10 {
				//This only works on testbed with real devices against simulators
				//t.Fatal("expected timestamp diff >= 10")
				q.Q("error: expected timestamp diff >= 10")
			}
		}
	}

	q.Q("done testing cluster")
}

func waitForRoot() error {
	for i := 1; i < 20; i++ {
		time.Sleep(2 * time.Second)
		url := "http://localhost:27182/api"
		//resp, err := http.Get(url)
		resp, err := GetWithToken(url, adminToken)
		if err != nil {
			return err
		}
		if resp != nil {
			if resp.StatusCode == 200 {
				return nil
			}
			q.Q(resp.StatusCode)
			q.Q("root service is ready")
			return nil
		}
		defer resp.Body.Close()
	}
	q.Q("time out waiting for root service")
	return fmt.Errorf("timeout waiting for root")
}

func getDevices() (*map[string]DevInfo, error) {
	url := "http://localhost:27182/api/v1/devices"
	//resp, err := http.Get(url)
	resp, err := GetWithToken(url, adminToken)

	if err != nil {
		return nil, err
	}

	devices := make(map[string]DevInfo)

	if resp != nil {
		//save close
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("error: response status code %v", resp.StatusCode)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &devices)
		if err != nil {
			return nil, err
		}
		q.Q(devices)
		l, err := script.Echo(string(body)).JQ("[.[]]|length").String()
		if err != nil {
			return nil, err
		}
		l = strings.TrimSpace(l)
		q.Q("fetched devices", devices, l, len(devices))
	}
	return &devices, nil
}

/*
func listInterfaces() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ipconfig")
	} else {
		cmd = exec.Command("ifconfig")
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		q.Q("error running cmd", cmd.Args, err)
	}
	q.Q("output of cmd", cmd.Args, string(output))
}
*/
