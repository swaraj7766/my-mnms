package mnms

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/rfc3164"
	"github.com/influxdata/go-syslog/v3/rfc5424"
	"github.com/qeof/q"
	"gopkg.in/natefinch/lumberjack.v2"
)

var TotalLogsReceived int
var TotalLogsWritten int

const severityMask = 0x07
const facilityMask = 0xf8

/*
Severity: RFC3164 page 9 -- numerical code 0 through 7

	0       Emergency: system is unusable
	1       Alert: action must be taken immediately
	2       Critical: critical conditions
	3       Error: error conditions
	4       Warning: warning conditions
	5       Notice: normal but significant condition
	6       Informational: informational messages
	7       Debug: debug-level messages
*/
const (
	LOG_EMERG int = iota
	LOG_ALERT
	LOG_CRIT
	LOG_ERR
	LOG_WARNING
	LOG_NOTICE
	LOG_INFO
	LOG_DEBUG
)

/*
Facility: RFC3164 page 7  -- numerical code 0 through 23

	 0             kernel messages
	 1             user-level messages
	 2             mail system
	 3             system daemons
	 4             security/authorization messages (note 1)
	 5             messages generated internally by syslogd
	 6             line printer subsystem
	 7             network news subsystem
	 8             UUCP subsystem
	 9             clock daemon (note 2)
	10             security/authorization messages (note 1)
	11             FTP daemon
	12             NTP subsystem
	13             log audit (note 1)
	14             log alert (note 1)
	15             clock daemon (note 2)
	16             local use 0  (local0)
	17             local use 1  (local1)
	18             local use 2  (local2)
	19             local use 3  (local3)
	20             local use 4  (local4)
	21             local use 5  (local5)
	22             local use 6  (local6)
	23             local use 7  (local7)
*/
const (
	LOG_KERN int = iota << 3
	LOG_USER
	LOG_MAIL
	LOG_DAEMON
	LOG_AUTH
	LOG_SYSLOG
	LOG_LPR
	LOG_NEWS
	LOG_UUCP
	LOG_CRON
	LOG_AUTHPRIV
	LOG_FTP
	_ // unused
	_ // unused
	_ // unused
	_ // unused
	LOG_LOCAL0
	LOG_LOCAL1
	LOG_LOCAL2
	LOG_LOCAL3
	LOG_LOCAL4
	LOG_LOCAL5
	LOG_LOCAL6
	LOG_LOCAL7
)

func StartSyslogServer() {
	q.Q(QC.SyslogServerAddr)
	udpsock, err := net.ListenPacket("udp4", QC.SyslogServerAddr)
	if err != nil {
		q.Q(err)
		return
	}
	defer udpsock.Close()
	buf := make([]byte, 1024*2)
	// TODO when need:  tcp syslog service
	for {
		mlen, raddr, err := udpsock.ReadFrom(buf)
		if err != nil {
			q.Q(err)
		}
		q.Q("syslog input", raddr, mlen)
		err = syslogInput(mlen, buf)
		if err != nil {
			// Implement saving and rotating logs locally. Currently
			// if there is no remote syslog server specified we drop the logs.
			if QC.IsRoot {
				SaveLog(string(buf[:mlen]))
			}
		}
	}
}

func InitRemoteSyslog() error {
	if QC.RemoteSyslogServerAddr == "" {
		return fmt.Errorf("Missing remote syslog server address")
	}
	if QC.RemoteSyslogServer != nil {
		QC.RemoteSyslogServer.Close()
	}
	udpsock, err := net.Dial("udp4", QC.RemoteSyslogServerAddr)
	if err != nil {
		q.Q(err)
		return err
	}
	QC.RemoteSyslogServer = udpsock

	return nil
}

/*
	func InitRemoteSyslogTcp() error {
		if QC.RemoteSyslogServerAddrTcp == "" {
			return fmt.Errorf("Missing remote syslog serverf of tcp address")
		}
		if QC.RemoteSyslogServerTcp != nil {
			QC.RemoteSyslogServerTcp.Close()
		}
		tcpsock, err := net.Dial("tcp", QC.RemoteSyslogServerAddrTcp)
		if err != nil {
			q.Q(err)
			return err
		}
		QC.RemoteSyslogServerTcp = tcpsock

		return nil
	}
*/
func SyslogParsePriority(buf string) (int, int, error) {
	if !strings.HasPrefix(buf, "<") {
		return 0, 0, fmt.Errorf("no syslog priority start character")
	}
	ix := strings.Index(buf, ">")
	if ix < 0 {
		return 0, 0, fmt.Errorf("no syslog priority end character")
	}
	priority, err := strconv.Atoi(buf[1:ix])
	if err != nil {
		return 0, 0, err
	}
	facility := priority / 8
	severity := priority % 8
	return facility, severity, nil
}

func syslogInput(mlen int, buf []byte) error {
	bufStr := string(buf[:mlen])
	q.Q("syslog input", bufStr)
	TotalLogsReceived++
	_, severity, err := SyslogParsePriority(bufStr)
	if err != nil {
		q.Q(err)
		return err
	}

	if severity < 6 {
		ix := strings.Index(bufStr, ">")
		wsMessage := WebSocketMessage{
			Kind:    "mnms_syslog",
			Level:   severity,
			Message: strings.TrimSpace(bufStr[ix+1:]),
		}
		QC.WebSocketMessageBroadcast <- wsMessage
		q.Q("forward to ws", wsMessage)
	}
	if QC.RemoteSyslogServer == nil {
		// First time, initialize client to remote syslog service
		if err := InitRemoteSyslog(); err != nil {
			q.Q(err)
			return err
		}
	}
	_, err = QC.RemoteSyslogServer.Write(buf[:mlen])
	if err != nil {
		q.Q(err)
		// upon failure, re-establish remote client and attemp to write again
		if err := InitRemoteSyslog(); err != nil {
			return err
		}
		_, err := QC.RemoteSyslogServer.Write(buf[:mlen])
		if err != nil {
			return err
		}
	}

	TotalLogsWritten++
	q.Q(TotalLogsReceived, TotalLogsWritten, mlen, string(buf[:mlen]))
	return nil
}

func SendSyslog(priority int, tag string, msg string) error {
	if QC.RemoteSyslogServerAddr == "" {
		q.Q("Missing remote syslog server address, can't send syslog")
		return errors.New("Missing remote syslog server address")
	}
	//reuse udp socket instead of open/close per message
	if QC.RemoteSyslogServer == nil {
		udpSock, err := net.Dial("udp4", QC.RemoteSyslogServerAddr)
		if err != nil {
			q.Q(err)
			return err
		}
		QC.RemoteSyslogServer = udpSock
		// TODO close udpSock when program exits.
	}
	timestamp := time.Now().Format(time.Stamp) // XXX not RFC3339
	syslogmsg := fmt.Sprintf("<%d>%s %s %s: %s", priority, timestamp, QC.Name, tag, msg)
	_, err := QC.RemoteSyslogServer.Write([]byte(syslogmsg))
	if err != nil {
		q.Q(err)
		return err
	}
	q.Q("sent syslog", string(msg))
	return nil
}

var SyslogPath = path.Dir(path.Join(os.TempDir(), dir))

const dir = "mnmslog"

var file = path.Join(dir, "syslog.log")

/*
func mkdir() {



	if _, err := os.Stat("/" + filepath.Join(dir)); os.IsNotExist(err) {
		// make a pki directory, if not exist
		if err := os.MkdirAll(filepath.Join(dir), os.ModeDir|0755); err != nil {
			panic(err)
		}
	}
}*/

var Logger *lumberjack.Logger

func initLogger() *lumberjack.Logger {
	filename := path.Join(QC.SyslogLocalPath, file)
	Logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    int(QC.SyslogFileSize),
		MaxBackups: 10,
		Compress:   QC.SyslogCompress,
		LocalTime:  true,
	}
	return Logger
}

//Save syslog to file
func SaveLog(data string) {
	//mkdir()
	if Logger == nil {
		Logger = initLogger()
	} else {
		if Logger.Filename != (path.Join(QC.SyslogLocalPath, file)) || Logger.Compress != QC.SyslogCompress || Logger.MaxSize != int(QC.SyslogFileSize) {
			err := Logger.Close()
			if err != nil {
				q.Q(err)
			}
			Logger = initLogger()
			q.Q("SaveLog,change local syslog paramter:", Logger)
		}
	}

	_, err := Logger.Write([]byte(data))
	if err != nil {
		q.Q(err)
		return
	}
	_, err = Logger.Write([]byte("\n"))
	if err != nil {
		q.Q(err)
		return
	}

}

func SyslogSetPathCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var path string
	Unpack(ws[4:], &path)
	QC.SyslogLocalPath = path
	cmdinfo.Status = "ok"
	return cmdinfo
}

func SyslogSetMaxSizeCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var size string
	Unpack(ws[4:], &size)
	s, err := strconv.Atoi(size)
	if err != nil {
		q.Q(err)
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	QC.SyslogFileSize = uint(s)
	cmdinfo.Status = "ok"
	return cmdinfo
}

func SyslogSetCompressCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 5 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var enable string
	Unpack(ws[4:], &enable)
	boolValue, err := strconv.ParseBool(enable)
	if err != nil {
		q.Q(err)
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	QC.SyslogCompress = boolValue
	cmdinfo.Status = "ok"
	return cmdinfo
}

const foramt = "2006/01/02 15:04:05"

func ReadSyslogCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	maxline := 0
	ws := strings.Split(cmd, " ")
	if len(ws) < 4 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	var start, end string
	filtertime := false
	if len(ws) == 8 || len(ws) == 9 {
		filtertime = true
		start = strings.Join(ws[4:6], " ")
		end = strings.Join(ws[6:8], " ")
		//if inlcude max line paramters
		if len(ws) == 9 {
			v, err := strconv.Atoi(ws[8])
			if err != nil {
				q.Q(err)
				cmdinfo.Status = err.Error()
				return cmdinfo
			} else {
				maxline = v
				q.Q("read syslog max line=", maxline)
			}
		}
	}
	//if inlcude max line paramters
	if len(ws) == 5 {
		v, err := strconv.Atoi(ws[4])
		if err != nil {
			q.Q(err)
			cmdinfo.Status = err.Error()
			return cmdinfo
		} else {
			maxline = v
			q.Q("read syslog max line=", maxline)
		}
	}

	readFile, err := os.Open(path.Join(QC.SyslogLocalPath, file))
	if err != nil {
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	logs := []syslog.Base{}
	for fileScanner.Scan() {
		b, t, err := parsingDataofSyslog(fileScanner.Text())
		if err != nil {
			continue
		}
		if filtertime {
			r, _ := compareTime(start, end, t.Format(foramt))
			if r {
				logs = append(logs, b)

			}
		} else {
			logs = append(logs, b)
		}
		if maxline != 0 {

			if len(logs) >= maxline {
				break
			}
		}
	}
	_ = readFile.Close()
	b, err := json.Marshal(&logs)
	if err != nil {
		cmdinfo.Status = err.Error()
		return cmdinfo
	}
	cmdinfo.Result = string(b)
	log.Print(cmdinfo.Result)
	cmdinfo.Status = "ok"
	return cmdinfo
}

func parsingDataofSyslog(message string) (syslog.Base, time.Time, error) {
	p := rfc3164.NewParser(rfc3164.WithYear(rfc3164.CurrentYear{}))
	m, err := p.Parse([]byte(message))
	if err != nil {
		p = rfc5424.NewParser()
		m, err = p.Parse([]byte(message))
		if err != nil {
			q.Q(err)
			return syslog.Base{}, time.Time{}, err
		}
	}

	switch v := m.(type) {
	case *rfc5424.SyslogMessage:

		return v.Base, *v.Timestamp, nil
	case *rfc3164.SyslogMessage:
		return v.Base, *v.Timestamp, nil
	}

	return syslog.Base{}, time.Time{}, errors.New("not support type yet")

}

//compareTime compare time size with start time and end time
func compareTime(start, end, target string) (bool, error) {

	s, err := time.Parse(foramt, start)
	if err != nil {
		return false, err
	}
	e, err := time.Parse(foramt, end)
	if err != nil {
		return false, err
	}
	if e.Before(s) {
		return false, fmt.Errorf("end time:%v should than start time:%v ", end, s)
	}
	t, err := time.Parse(foramt, target)
	if err != nil {
		return false, err
	}

	if (s.Before(t) || s.Equal(t)) && (e.After(t) || e.Equal(t)) {
		return true, nil
	}

	return false, nil
}
