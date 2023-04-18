package mnms

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"testing"

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

	q.Q("wait for root to become ready...")
	if err := waitForRoot(); err != nil {
		t.Fatal(err)
	}

	// http request to localhost:27182/api
	resp, err := http.Get("http://localhost:27182/api")
	if err != nil {
		t.Fatal("mnmsctl is not running, should run mnmsctl first")
	}
	if resp == nil {
		t.Fatal("resp should not be nil")
	}
	// save close, resp should not be nil here
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
			t.Error(err)
		}
	}()

	// check default config.json
	c, err := GetUserConfig("admin")
	if err != nil {
		t.Fatal(err)
	}
	if c.Name != "admin" {
		t.Fatal("default user should be admin")
	}
	if c.Password != AdminDefaultPassword {
		t.Fatal("default password should be default")
	}

	body, err := json.Marshal(map[string]string{
		"user":     "admin",
		"password": AdminDefaultPassword,
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

	if res == nil {
		t.Fatal("res should not be nil")
	}

	var recBody map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&recBody)
	if err != nil {
		t.Fatal(err)
	}
	// save close, res should not be nil here
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
	if res == nil {
		t.Fatal("res should not be nil")
	}
	// save close, res should not be nil here
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
	if res == nil {
		t.Fatal("res should not be nil")
	}
	// save close, res should not be nil here
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
	if res == nil {
		t.Fatal("res should not be nil")
	}
	// save close, res should not be nil here
	res.Body.Close()

	// post command should fail
	cmdinfo := make(map[string]CmdInfo)
	cmd := "scan gwd"
	t.Log("cmd: ", cmd)
	insertcmd(cmd, &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Errorf("json marshal %v", err)
		return
	}

	resp, err = http.Post("http://localhost:27182/api/v1/commands", "application/text",
		bytes.NewBuffer([]byte(jsonBytes)))
	if resp == nil {
		t.Fatal("resp should not be nil")
	}
	if resp.StatusCode != http.StatusUnauthorized {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Logf("msg: %s", msg)
		t.Errorf("post command should fail err: %v code:%d", err, resp.StatusCode)
	}

	// save close, resp should not be nil here
	resp.Body.Close()

	// PostWithToken should pass
	resp, err = PostWithToken("http://localhost:27182/api/v1/commands", admintoken, bytes.NewBuffer([]byte(jsonBytes)))
	if resp == nil {
		t.Fatal("resp should not be nil")
	}
	if resp.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(resp.Body)
		t.Logf("msg: %s", msg)
		t.Errorf("post command should pass err: %v code:%d", err, resp.StatusCode)
	}
	// save close, resp should not be nil here
	resp.Body.Close()

	// post with token
	req, err = http.NewRequest("POST", "http://localhost:27182/api/v1/commands", bytes.NewBuffer([]byte(jsonBytes)))
	if err != nil {
		t.Fatal("create request fail", err)
	}

	req.Header.Set("Authorization", "Bearer "+string(token))
	res, err = http.DefaultClient.Do(req)
	if res == nil {
		t.Fatal("resp should not be nil")
	}
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to get /api/v1/devices with token")
	}
	// save close, res should not be nil here
	res.Body.Close()
}

// TestJWT tests JWT
func TestJWT(t *testing.T) {
	QC.Name = "test"

	// reset config
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	// sample token string taken from the New example
	tokenString, err := generateJWT("admin", AdminDefaultPassword)
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

// TestFileServer tests file server
func TestFileServer(t *testing.T) {
	QC.IsRoot = true
	QC.Name = "root"
	// run services
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// root and non-root run http api service.
		// if both root and non-root runs on same machine (should not happen)
		// then whoever runs first will run http service (port conflict)
		defer wg.Done()
		HTTPMain()
		// TODO root to dump snapshots of devices, logs, commands
	}()

	q.Q("wait for root to become ready...")
	if err := waitForRoot(); err != nil {
		t.Fatal(err)
	}

	// create a file for testing
	fileName := "test.txt"
	fileDir, err := CheckStaticFilesFolder()
	if err != nil {
		t.Fatal("check static files folder fail", err)
	}
	filePath := filepath.Join(fileDir, fileName)
	// fill file with some content to test file server
	// write 3k rand data
	data := make([]byte, 3*1024)
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		t.Fatal("generate rand data fail", err)
	}
	// write data to filePath
	err = ioutil.WriteFile(filePath, data, 0o644)
	if err != nil {
		t.Fatal("write file fail", err)
	}

	// GetWithToken can pass
	admintoken, err := GetToken("admin")
	if err != nil {
		t.Fatal("get token fail", err)
	}
	res, err := GetWithToken("http://localhost:27182/api/v1/files/test.txt", admintoken)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatal("should be able to get /api/v1/devices with token")
	}
	if res == nil {
		t.Fatal("res should not be nil")
	}
	// check content
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("read body fail", err)
	}
	if len(content) == 0 {
		t.Fatal("content should not be empty")
	}
	// compare with data
	if !bytes.Equal(content, data) {
		t.Fatal("content should be equal to data")
	}
	// save close, res should not be nil here
	res.Body.Close()

	// Check without token should fail
	res, err = http.Get("http://localhost:27182/api/v1/files/test.txt")
	if err != nil || res.StatusCode != http.StatusUnauthorized {
		t.Fatal("should not be able to get /api/v1/devices without token")
	}
	if res == nil {
		t.Fatal("res should not be nil")
	}
	res.Body.Close()

	// delete file
	err = os.Remove(filePath)
	if err != nil {
		t.Fatal("delete file fail", err)
	}
}
