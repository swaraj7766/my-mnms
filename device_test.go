package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mnms/pkg/simulator"
	"mnms/pkg/simulator/net"
	atopyaml "mnms/pkg/simulator/yaml"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bitfield/script"
	"github.com/qeof/q"
)

func init() {
	q.O = "stderr"
	q.P = ".*"
}

func TestDevices(t *testing.T) {
	killNmsctlProcesses()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		HTTPMain()
	}()

	myName := "test123"
	adminToken, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	url := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	resp, err := PostWithToken(url, adminToken, bytes.NewBuffer([]byte(myName)))
	if err != nil {
		t.Fatalf("post %v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	resp.Body.Close()
	err = ReadTestData()
	if err != nil {
		t.Fatal(err)
	}
	url = fmt.Sprintf("http://localhost:%d/api/v1/devices", QC.Port)
	resp, err = GetWithToken(url, adminToken)

	if err != nil {
		t.Fatalf("get %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("response status code expect 200 but  %v", resp.StatusCode)
	}
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("reading resp, %v", err)
		}
		q.Q(string(body))
		l, err := script.Echo(string(body)).JQ("[.[]]|length").String()
		if err != nil {
			t.Fatalf("%v", err)
		}
		l = strings.TrimSpace(l)
		q.Q(l)
		lval, err := strconv.Atoi(l)
		if err != nil {
			t.Fatal(err)
		}

		if lval < 45 {
			t.Fatalf("wrong len, %v", l)
		}
		l, err = script.Echo(string(body)).JQ(`.[]|select(.mac=="02-42-C0-A8-64-86" ) | .ipaddress`).String()
		if err != nil {
			t.Fatalf("%v", err)
		}
		l = strings.TrimSpace(l)
		q.Q(l)
		if l != `"192.168.100.134"` {
			t.Fatalf("wrong ipaddr, %v", l)
		}
	}
	resp.Body.Close()
	q.Q("end test devices")
}

func ReadTestData() error {
	dat, err := os.ReadFile("testdata.json")
	if err != nil {
		q.Q("can't read file %v", err)
		return err
	}
	var devinfo map[string]DevInfo
	err = json.Unmarshal(dat, &devinfo)
	if err != nil {
		q.Q("can't unmarshal", err)
		return err
	}

	for _, v := range devinfo {
		if !InsertDev(v) {
			q.Q("can't insert dev", err)
			return err
		}
	}
	return nil
}

func TestDevicesSaveCmd(t *testing.T) {
	q.Q("test devices save cmd")
	// run simultors
	simmap := map[string]atopyaml.Simulator{}
	simmap["group1"] = atopyaml.Simulator{Number: 1, DeviceType: "EH7506", StartPreFixIp: "192.168.10.1/24", MacAddress: "00-60-E9-28-01-01"}
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

	// run mnms
	QC.IsRoot = true
	QC.Name = "root"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		HTTPMain()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		GwdMain()
	}()
	time.Sleep(time.Second * 1)

	// auth steps...
	// http request to localhost:27182/api
	resp, err := http.Get("http://localhost:27182/api")
	if err != nil {
		t.Fatal("mnmsctl is not running, should run mnmsctl first")
	}
	resp.Body.Close()

	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	adminPass, err := GenPassword(QC.Name, "admin")
	if err != nil {
		t.Fatal(err)
	}
	body, err := json.Marshal(map[string]string{
		"user":     "admin",
		"password": adminPass,
	})
	if err != nil {
		t.Fatal(err)
	}
	// make a post request with username = admin and passwrord = adminPass
	res, err := http.Post("http://localhost:27182/api/v1/login", "application/json", bytes.NewBuffer(body))
	if err != nil || res.StatusCode != http.StatusOK {
		resText, _ := ioutil.ReadAll(res.Body)
		t.Log("resText", string(resText))
		t.Fatal("should be able to login with admin and adminPass")
	}
	var recBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&recBody)
	if err != nil {
		t.Fatal(err)
	}

	token, ok := recBody["token"].(string)
	if !ok {
		t.Fatal("token is not string")
	}
	t.Log("Got token:", string(body))
	res.Body.Close()

	// send scan
	cmdinfo := make(map[string]CmdInfo)
	insertcmd("scan gwd", &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+string(token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to send scan command")
	}
	resp.Body.Close()
	req.Body.Close()

	// wait for scan to finish
	time.Sleep(time.Second * 1)
	_ = CheckCommands()
	time.Sleep(time.Second * 10)

	// get devices
	req, err = http.NewRequest("GET", "http://localhost:27182/api/v1/devices", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+string(token))
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to get devices")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	devInfos := make(map[string]DevInfo)
	err = json.Unmarshal(body, &devInfos)
	if err != nil {
		t.Fatal(err)
	}
	if len(devInfos) == 0 {
		t.Fatal("no devices found")
	}
	q.Q("get devices", devInfos)
	res.Body.Close()

	// if not root, save devices should fail
	QC.IsRoot = false
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("devices save", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	req.Body.Close()
	res.Body.Close()

	_ = CheckCommands()
	time.Sleep(time.Second * 1)
	req, err = http.NewRequest("GET", "http://localhost:27182/api/v1/commands?cmd=all", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	cmdinfo = make(map[string]CmdInfo)
	err = json.Unmarshal(body, &cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(cmdinfo["devices save"].Status, "error") {
		t.Fatal("should not be able to save devices")
	}
	res.Body.Close()

	// if root, save devices should succeed
	QC.IsRoot = true
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to save devices")
	}
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	req.Body.Close()
	res.Body.Close()

	cmdinfo = make(map[string]CmdInfo)
	err = json.Unmarshal(body, &cmdinfo)
	if err != nil {
		t.Fatal(err)
	}
	q.Q("devices save cmd response", string(body), cmdinfo)
	if len(cmdinfo) == 0 {
		t.Fatal("no command found")
	}
	fn := cmdinfo["devices save"].Result
	if fn == "" {
		t.Fatal("devices save command result is empty")
	}

	// check if file exists
	_, err = os.Stat(fn)
	if err != nil {
		t.Fatal("devices save file not found")
	}

	// check if file is valid
	filedata, err := os.ReadFile(fn)
	if err != nil {
		t.Fatal("devices save file read error")
	}
	q.Q("devices save file data", string(filedata))
	devInfos = make(map[string]DevInfo)
	err = json.Unmarshal(filedata, &devInfos)
	if err != nil {
		t.Fatal("devices save file json unmarshal error")
	}
	if len(devInfos) == 0 {
		t.Fatal("no devices found in devices save file")
	}

	// remove file
	err = os.Remove(fn)
	if err != nil {
		t.Fatal("devices save file remove error")
	}
	q.Q("end test devices save cmd")
}

func TestDevicesLoad(t *testing.T) {
	q.Q("test devices load")

	// if not root, load devices should fail
	QC.IsRoot = false
	err := LoadDevices("testdata.json")
	if err == nil {
		t.Fatal("should not be able to load devices")
	}

	// if root, load devices should succeed
	QC.IsRoot = true
	err = LoadDevices("testdata.json")
	if err != nil {
		t.Fatal("should be able to load devices")
	}
	q.Q("loaded test devices", QC.DevData)

	// check devices
	if len(QC.DevData) == 0 {
		t.Fatal("no devices found")
	}
}

func TestDevicesFilesList(t *testing.T) {
	q.Q("test devices files list")

	// load fake devices
	QC.IsRoot = true
	err := LoadDevices("testdata.json")
	if err != nil {
		t.Fatal("should be able to load devices")
	}
	// save devices in fn1
	fn1, err := SaveDevices()
	if err != nil {
		t.Fatal("should be able to save devices")
	}
	defer os.Remove(fn1)
	time.Sleep(time.Second * 1)
	// save devices in fn2
	fn2, err := SaveDevices()
	if err != nil {
		t.Fatal("should be able to save devices")
	}
	defer os.Remove(fn2)

	// if not root, list files should fail
	QC.IsRoot = false
	_, err = ListDevicesFiles()
	if err == nil {
		t.Fatal("should not be able to list files")
	}

	// if root, list files should succeed
	QC.IsRoot = true
	list, err := ListDevicesFiles()
	if err != nil {
		t.Fatal("should be able to list files")
	}
	q.Q("devices files list", list)
	if len(list) != 2 {
		t.Fatal("should be 2 files")
	}
	if list[0] != fn1 && list[0] != fn2 {
		t.Fatal("invalid file name")
	}
	if list[1] != fn1 && list[1] != fn2 {
		t.Fatal("invalid file name")
	}
}
