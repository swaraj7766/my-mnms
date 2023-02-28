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
	"sync"
	"testing"
	"time"

	"github.com/qeof/q"
)

func TestCrontabCmd(t *testing.T) {
	// run simultors
	simmap := map[string]atopyaml.Simulator{}
	simmap["group1"] = atopyaml.Simulator{Number: 1, DeviceType: "EH7506", StartPreFixIp: "192.168.4.1/24", MacAddress: "00-60-E9-28-01-01"}
	name, err := net.GetDefaultInterfaceName()
	if err != nil {
		t.Fatalf("get default interface name %v", err)
	}
	simulators, err := atopyaml.NewSimulator(simmap, name)
	if err != nil {
		t.Fatalf("new simulator %v", err)
	}
	for _, sim := range simulators {
		err := sim.StartUp()
		if err != nil {
			t.Fatal(err)
		}
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

	// add cron jobs
	q.Q("add cron jobs")
	cmdinfo := make(map[string]CmdInfo)
	insertcmd("config crontab add * * * * * scan gwd", &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err := http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}

	resp.Body.Close()
	req.Body.Close()
	time.Sleep(time.Second * 1)

	// get cron jobs
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("config crontab list", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}
	cmdresult := make(map[string]CmdInfo)
	err = json.NewDecoder(res.Body).Decode(&cmdresult)
	if err != nil {
		t.Fatal(err)
	}
	crons := make([]CronInfo, 0)
	for _, v := range cmdresult {
		err := json.Unmarshal([]byte(v.Result), &crons)
		if err != nil {
			t.Fatalf("unmarshal %v", err)
		}
	}
	q.Q("cronlist", cmdresult, "crons", crons)
	if len(crons) != 1 {
		t.Fatalf("jobs not match")
	}
	entryID := crons[0].EntryID

	res.Body.Close()
	req.Body.Close()

	// sleep at least 1 minute
	// t.Log("sleep at least 1 minute")
	// time.Sleep(time.Minute * 1)
	// CheckCommands()
	// get devices
	// req, err = http.NewRequest("GET", "http://localhost:27182/api/v1/devices", nil)
	// if err != nil {
	// 	t.Fatalf("new request %v", err)
	// }
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	// res, err = http.DefaultClient.Do(req)
	// if err != nil || res.StatusCode != http.StatusOK {
	// 	t.Fatalf("get %v", err)
	// }
	// devices := make(map[string]DevInfo)
	// json.NewDecoder(res.Body).Decode(&devices)
	// if len(devices) == 0 {
	// 	t.Fatalf("no devices")
	// }
	// res.Body.Close()

	// dump cron jobs
	q.Q("dump cron jobs")
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("config crontab dump", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}
	defer func() {
		err = RemoveCrontab()
		if err != nil {
			q.Q("remove crontab", err)
		}
	}()
	resp.Body.Close()
	req.Body.Close()
	time.Sleep(time.Second * 1)

	// delete cron jobs
	q.Q("delete cron jobs")
	cmdinfo = make(map[string]CmdInfo)
	insertcmd(fmt.Sprintf("config crontab delete %d", entryID), &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}

	resp.Body.Close()
	req.Body.Close()
	time.Sleep(time.Second * 1)

	// get cron jobs again to check
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("config crontab list", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}
	cmdresult = make(map[string]CmdInfo)
	err = json.NewDecoder(resp.Body).Decode(&cmdresult)
	if err != nil {
		t.Fatal(err)
	}
	crons = make([]CronInfo, 0)
	for _, v := range cmdresult {
		err := json.Unmarshal([]byte(v.Result), &crons)
		if err != nil {
			t.Fatalf("unmarshal %v", err)
		}
	}
	if len(crons) != 0 {
		t.Fatalf("jobs not match")
	}
	req.Body.Close()
	resp.Body.Close()

	// load cron jobs
	q.Q("load cron jobs")
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("config crontab load", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err = http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}

	resp.Body.Close()
	req.Body.Close()
	// time.Sleep(time.Minute*1 + time.Second*10)

	// get cron jobs again to check
	cmdinfo = make(map[string]CmdInfo)
	insertcmd("config crontab list", &cmdinfo)
	jsonBytes, err = json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("marshal %v", err)
	}
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatalf("new request %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("post %v", err)
	}
	cmdresult = make(map[string]CmdInfo)
	_ = json.NewDecoder(res.Body).Decode(&cmdresult)
	crons = make([]CronInfo, 0)
	for _, v := range cmdresult {
		err := json.Unmarshal([]byte(v.Result), &crons)
		if err != nil {
			t.Fatalf("unmarshal %v", err)
		}
	}
	q.Q("cron job list after load", crons)
	if len(crons) != 1 {
		t.Fatalf("jobs not match")
	}
	req.Body.Close()
	res.Body.Close()
}
