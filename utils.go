package mnms

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"

	"github.com/qeof/q"
)

// random misc stuff
func CleanStr(val string) string {
	val = strings.Replace(val, ",", "", -1)
	val = strings.Replace(val, `"`, "", -1)
	val = strings.TrimSpace(val)
	return val
}

func Unpack(src []string, dst ...*string) {
	for ind, val := range dst {
		*val = src[ind]
	}
}

func SetupFakeClient() error {
	q.Q(QC.Name)

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
	if QC.Name == "client1" {
		for _, v := range devinfo {
			if strings.HasPrefix(v.Mac, "00-60-E9") {
				if !InsertDev(v) {
					q.Q("can't insert dev", err)
					return err
				}
			}
		}
	} else if QC.Name == "client2" {
		for _, v := range devinfo {
			if strings.HasPrefix(v.Mac, "02-42-") {
				if !InsertDev(v) {
					q.Q("can't insert dev", err)
					return err
				}
			}
		}
	} else {
		q.Q("invalid fake client name", QC.Name)
		return err
	}
	return nil
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CreateFile(name string) error {
	fo, err := os.Create(name)
	if err != nil {
		return err
	}
	defer func() {
		fo.Close()
	}()
	return nil
}

/*
Do not use this until Windows port is done.

func IsAdminUser() bool {

	currentUser, err := user.Current()
	if err != nil {
		q.Q(err)
	}
	return currentUser.Username == "root"
}
*/

func DoExit(code int) {
	q.Q("exiting mnms", code)
	if QC.DumpStackTrace && code != 0 {
		buf := make([]byte, 1<<16)
		stackSize := runtime.Stack(buf, true)
		q.Q(string(buf[0:stackSize]))
	}
	os.Exit(code)
}
