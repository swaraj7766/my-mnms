package mnms

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/qeof/q"
)

func init() {
	q.O = "stderr"
	q.P = ".*"
}

func TestCleanStr(t *testing.T) {
	str := CleanStr("\"Managed, 3.2\"")
	expected := "Managed 3.2"
	if str != expected {
		t.Fatalf("CleanStr %s != %s\n", str, expected)
	}
	q.Q(str)
}

func TestCheckFolder(t *testing.T) {
	_, err := CheckMNMSFolder()
	if err != nil {
		t.Fatal(err)
	}
	usersPath, err := checkMNMSConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	// write some data to users file
	err = ioutil.WriteFile(usersPath, []byte("test"), 0600)
	if err != nil {
		t.Fatal(err)
	}
	// read back and compare
	ret, err := ioutil.ReadFile(usersPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(ret) != "test" {
		t.Fatalf("Read file %s != test\n", string(ret))
	}
	// delete users file
	err = os.Remove(usersPath)
	if err != nil {
		t.Fatal(err)
	}
}
