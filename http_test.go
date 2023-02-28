package mnms

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/qeof/q"
)

// to test https
// 1. start mnms `./mnmsctl/mnmsctl -s`
// 2. check out IP address of the machine `curl ipconfig.io`
// 3. start reverse proxy `sudo caddy reverse-proxy --from 122.147.151.234.sslip.io --to https://localhost:27182`

// Testing jwt with curl
// curl -H 'Accept: application/json' -H "Authorization: Bearer ${TOKEN}" https://{hostname}/api/myresource
// example get token
// curl -X POST -H 'Accept: application/json'  https://localhost:27182/api/v1/login -d '{"user":"admin"}'
// get user password
// curl -H 'Accept: application/json' -H "Authorization: Bearer ${TOKEN}" https://localhost:27182/api/v1/users -d '{"username":"admin"}'

func TestAuthentication(t *testing.T) {
	QC.IsRoot = true
	QC.Name = "root"
	q.P = "jwt"
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

	time.Sleep(time.Second * 5)

	// http request to localhost:27182/api
	resp, err := http.Get("http://localhost:27182/api")
	if err != nil {
		t.Fatal("mnmsctl is not running, should run mnmsctl first")
	}
	resp.Body.Close()

	// init
	err = cleanMNMSConfig()
	if err != nil {
		q.Q(err)
	}
	err = InitDefaultMNMSConfigIfNotExist()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := cleanMNMSConfig()
		if err != nil {
			t.Fatal(err)
		}
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
	res.Body.Close()
	t.Log("/login response: ", recBody)
	token, ok := recBody["token"].(string)
	if !ok {
		t.Fatal("token is not string")
	}

	// Try to get /api/v1/devices without token should fail
	res, err = http.Get("http://localhost:27182/api/v1/devices")
	if err != nil || res.StatusCode != http.StatusUnauthorized {
		t.Fatal("should not be able to get /api/v1/devices without token")
	}
	res.Body.Close()

	// GetWithToken can pass
	admintoken, err := GetToken("admin")
	if err != nil {
		t.Fatal("get token fail", err)
	}
	res, err = GetWithToken("http://localhost:27182/api/v1/devices", admintoken)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to get /api/v1/devices with token")
	}
	res.Body.Close()

	// Try to get /api/v1/devices with jwt token should success
	req, err := http.NewRequest("GET", "http://localhost:27182/api/v1/devices", nil)
	if err != nil {
		t.Fatal("create request fail", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(token))
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to get /api/v1/devices with token")
	}

	res.Body.Close()

	// post command should fail
	cmdinfo := make(map[string]CmdInfo)
	cmd := fmt.Sprintf("scan gwd")
	t.Log("cmd: ", cmd)
	insertcmd(cmd, &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Errorf("json marshal %v", err)
		return
	}

	resp, err = http.Post("http://localhost:27182/api/v1/commands", "application/text",
		bytes.NewBuffer([]byte(jsonBytes)))
	if resp.StatusCode != http.StatusUnauthorized {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Logf("msg: %s", msg)
		t.Errorf("post command should fail err: %v code:%d", err, resp.StatusCode)
	}
	resp.Body.Close()

	// PostWithToken should pass
	resp, err = PostWithToken("http://localhost:27182/api/v1/commands", admintoken, bytes.NewBuffer([]byte(jsonBytes)))
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Logf("msg: %s", msg)
		t.Errorf("post command should pass err: %v code:%d", err, resp.StatusCode)
	}
	resp.Body.Close()

	// post with token
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer([]byte(jsonBytes)))
	if err != nil {
		t.Fatal("create request fail", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(token))
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to get /api/v1/devices with token")
	}
	req.Body.Close()
	res.Body.Close()
}

// TestJWT tests JWT
func TestJWT(t *testing.T) {
	QC.Name = "test"
	pass, err := GenPassword(QC.Name, "admin")
	t.Log("name : ", QC.Name)
	if err != nil {
		t.Fatal(err)
	}
	// reset config
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	// sample token string taken from the New example
	tokenString, err := generateJWT("admin", pass)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tokenString)

	claims, err := parseJWT(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(claims)

	// create a invalied token
	privateKey, _ := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	token := jwt.NewWithClaims(jwt.SigningMethodES512, jwt.MapClaims(claims))
	tokenString, err = token.SignedString(privateKey)
	if err != nil {
		t.Fatal(err)
	}
	_, err = parseJWT(tokenString)
	if err == nil {
		t.Fatal("should be error")
	}
}

// TestgenPassword tests password generation
func TestGenPassword(t *testing.T) {
	expect := "I1wDY1uoQvWEW1Rt"
	password, err := GenPassword("masterPassword", "austin")
	if err != nil {
		t.Fatal(err)
	}

	if password != expect {
		t.Errorf("expected %q, got %q", expect, password)
	}
	// if change master password, generated password should change
	password2, err := GenPassword("masterPassword2", "austin")
	if err != nil {
		t.Fatal(err)
	}
	if password2 == expect {
		t.Errorf("expected not be %q, got %q", expect, password2)
	}

	// every time password shoul be same
	password3, err := GenPassword("12", "austin")
	if err != nil {
		t.Fatal(err)
	}
	password4, err := GenPassword("12", "austin")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("pass3 %s , pass4 %s", password3, password4)
	if password3 != password4 {
		t.Errorf("expected be %q, got %q", password3, password4)
	}
}
