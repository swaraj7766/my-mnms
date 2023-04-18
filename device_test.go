package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/bitfield/script"
	"github.com/qeof/q"
)

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
	q.Q("wait for root to become ready...")
	if err := waitForRoot(); err != nil {
		t.Fatal(err)
	}
	url := fmt.Sprintf("http://localhost:%d/api/v1/register", QC.Port)
	ci := ClientInfo{Name: myName}
	jsonBytes, err := json.Marshal(ci)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	resp, err := PostWithToken(url, adminToken, bytes.NewBuffer(jsonBytes))
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
		//save close, in if resp != nil block
		resp.Body.Close()
	}

	q.Q("end test devices")
}

func ReadTestData() error {
	dat, err := os.ReadFile("testdata.json")
	if err != nil {
		q.Q("can't read file", err)
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
