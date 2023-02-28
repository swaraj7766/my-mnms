package mnms

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/qeof/q"
)

func init() {
	q.O = "stderr"
	q.P = ".*"
}

func killNmsctlProcesses() {
	var cmd *exec.Cmd
	q.Q("killing nmnsctl processes")
	// XXX terrible dangerous killing of all mnmsctl
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/f", "/im", "mnmsctl.exe")
	} else {
		cmd = exec.Command("killall", "mnmsctl")
	}
	err := cmd.Run()
	if err != nil {
		q.Q(err)
	}
}
func TestSyslog(t *testing.T) {
	//RFC 3164 Page 10, Facility=20 and Severity=5 would have Priority value of 165
	syslogMsg := []byte("<165>Nov 11 12:34:56 myhost mytag: this is a syslog message")

	facility, severity, err := SyslogParsePriority(string(syslogMsg))
	if err != nil {
		t.Fatal(err)
	}
	q.Q(facility, severity)
	if facility != 20 {
		t.Fatal("facility is wrong")
	}
	if severity != 5 {
		t.Fatal("severity is wrong")
	}

	//logger -d -s -n localhost -P 5514 --rfc3164 -p local3.alert local3 alert syslog test
	//<153>Feb 12 02:23:21 cs-186432676255-default bob_bae: local3 alert syslog test

	QC.RemoteSyslogServerAddr = "localhost:5514"
	err = SendSyslog(LOG_ALERT, "testsyslog", "this is a test alert syslog message")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		killNmsctlProcesses()
	}()
	go func() {
		cmd := exec.Command("./mnmsctl/mnmsctl", "-n", "root", "-R", "-O", "root.log", "-s")
		err := cmd.Run()
		if err != nil {
			q.Q(err)
		}
	}()
	err = waitForRoot()
	if err != nil {
		q.Q(err)
		t.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(1 * time.Second)
			q.Q("sending syslog alert for websock")
			err = SendSyslog(LOG_ALERT, "syslog_test", "websock test alert syslog message")
			if err != nil {
				q.Q("Fatal", err)
				break
			}
			q.Q("sent syslog alert for websock")
		}
	}()
	connectWebsock()
}

func connectWebsock() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: "localhost:27182", Path: "/api/v1/ws"}
	q.Q(u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		q.Q(err)
		return
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				q.Q(err)
				return
			}
			q.Q("read websock", string(message))
		}
	}()
	ticker := time.NewTicker(time.Second)
	count := 1
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			msg := fmt.Sprintf(`{ "kind": "syslog_test", "level": 3, "message": "ticker %v" }`, t)
			count++
			if count > 3 {
				return
			}
			err := c.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				q.Q(err)
				return
			}
			q.Q("wrote to websock", msg)
		case <-interrupt:
			q.Q("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "closing at interrupt"))
			if err != nil {
				q.Q(err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
