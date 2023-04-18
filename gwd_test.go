package mnms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/qeof/q"
)

const (
	configmac      = "00-60-E9-18-0A-01"
	confignnewip   = "192.168.20.250"
	confignmask    = "255.255.255.0"
	configngateway = "192.168.20.254"
)

const (
	hostmac   = "00-60-E9-18-0A-02"
	hostvalue = "atoptest"
)

const loginuser = "admin"
const loginpwd = "default"

func TestGwdBeep(t *testing.T) {
	t.Skip()
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()

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

	d, err := GetDevData()
	if err != nil {
		t.Fatal(err)
	}

	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("beep %v %v", d.Mac, d.IPAddress)
	insertcmd(cmd, &cmdinfo)

	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	url := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)

	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := PostWithToken(url, token, bytes.NewBuffer(jsonBytes))
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("post err:%v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
		//save close
		defer resp.Body.Close()
	}

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCmds()
}

func TestGwdReset(t *testing.T) {
	t.Skip()
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()

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

	d, err := GetDevData()
	if err != nil {
		t.Fatal(err)
	}
	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("reset %v %v %v %v", d.Mac, d.IPAddress, loginuser, loginpwd)
	insertcmd(cmd, &cmdinfo)

	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))

	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)

	resp, err := PostWithToken(url, token, bytes.NewBuffer(jsonBytes))
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("post err:%v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
		// save close
		defer resp.Body.Close()
	}

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCmds()
}

func TestGwdConfigIP(t *testing.T) {
	t.Skip()
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()

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
	d, err := GetDevData()
	if err != nil {
		t.Fatal(err)
	}
	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("config net %v %v %v %v %v", d.Mac, confignnewip, confignmask, configngateway, hostvalue)
	insertcmd(cmd, &cmdinfo)

	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))

	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(url, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("no resp")
	}
	if resp.StatusCode != 200 {
		t.Fatalf("post status code %v", resp.StatusCode)
	}
	if resp != nil {
		q.Q(resp.Header)

	}

	// save close, already check resp is not nil
	defer resp.Body.Close()
	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCmds()
}

func GetDevData() (DevInfo, error) {
	if len(QC.DevData) == 0 {
		return DevInfo{}, errors.New("no device exist")
	} else {
		for _, v := range QC.DevData {
			return v, nil
		}
	}
	return DevInfo{}, errors.New("no device exist")
}
