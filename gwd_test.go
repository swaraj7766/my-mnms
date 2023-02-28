package mnms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/qeof/q"
)

var waititime = time.Second * 10

const (
	beepmac = "00-60-E9-18-01-01"
	beepip  = "192.168.10.1"
)

const (
	rebootpmac = "00-60-E9-18-01-02"
	rebootip   = "192.168.10.2"
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
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()

	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("beep %v %v", beepmac, beepip)
	insertcmd(cmd, &cmdinfo)

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		GwdMain()
	}()
	err = waitFindDevice(beepmac, waititime)
	if err != nil {
		t.Fatal(err)
	}
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

	q.Q(QC.Root)
	resp, err := PostWithToken(url, token, bytes.NewBuffer(jsonBytes))
	if resp.StatusCode != 200 || err != nil {
		t.Fatalf("post err:%v", err)
	}
	if resp != nil {
		q.Q(resp.Header)
	}
	defer resp.Body.Close()

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCommands()
}

func TestGwdReboot(t *testing.T) {
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()
	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("reboot %v %v %v %v", rebootpmac, rebootip, loginuser, loginpwd)
	insertcmd(cmd, &cmdinfo)

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		GwdMain()
	}()
	err = waitFindDevice(rebootpmac, waititime)
	if err != nil {
		t.Fatal(err)
	}

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
	}
	defer resp.Body.Close()

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCommands()
}

func TestGwdConfigIP(t *testing.T) {
	err := createSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer ShutdownSimulator()
	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("config net %v %v %v %v %v", configmac, confignnewip, confignmask, configngateway, hostvalue)
	insertcmd(cmd, &cmdinfo)

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
	wg.Add(1)
	go func() {
		defer wg.Done()
		GwdMain()
	}()
	err = waitFindDevice(configmac, waititime)
	if err != nil {
		t.Fatal(err)
	}
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
	q.Q(QC.Root)
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
	defer resp.Body.Close()

	q.Q(QC.CmdData)
	q.Q(QC.Clients)
	_ = CheckCommands()
}

func waitFindDevice(mac string, timeout time.Duration) error {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("time out,device:%v can't find", mac)
		default:
			time.Sleep(500 * time.Millisecond)
			_, err := FindDev(mac)
			if err != nil {
				continue
			}
			return nil
		}

	}

}
