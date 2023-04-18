package mnms

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/go-syslog/v3"
	"github.com/influxdata/go-syslog/v3/rfc3164"
	"github.com/influxdata/go-syslog/v3/rfc5424"
	"github.com/qeof/q"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	TotalLogsReceived int
	TotalLogsSent     int
	TotalLogsDropped  int
)

/*
const severityMask = 0x07
const facilityMask = 0xf8
*/
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
	// TODO when need:  tcp syslog service
	for {
		buf := make([]byte, 1024*2)
		mlen, raddr, err := udpsock.ReadFrom(buf)
		q.Q(len(buf))
		if err != nil {
			q.Q(err)
		}
		q.Q("syslog input", raddr, mlen)
		err = syslogInput(mlen, buf)
		if err != nil {
			// Implement saving and rotating logs locally. Currently
			// if there is no remote syslog server specified we drop the logs.
			if QC.IsRoot {
				_, _, err := parsingDataofSyslog(string(buf[:mlen]))
				if err != nil {
					f, b, err := SyslogParsePriority(string(buf[:mlen]))
					if err != nil {
						continue
					}
					p := fmt.Sprintf("<%v>", (f*8)+b)
					m := strings.ReplaceAll(string(buf[:mlen]), p, "")
					message := fmt.Sprintf("%v%v %v %v", p, time.Now().Format("Jan 02 15:04:05"), raddr.String(), m)
					SaveLog(message)
				} else {
					SaveLog(string(buf[:mlen]))
				}
			}
		}
	}
}

func InitRemoteSyslog() error {
	if QC.RemoteSyslogServerAddr == "" {
		return fmt.Errorf("%v", "Missing remote syslog server address")
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
	SendSocketMessage(severity, bufStr)
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

	TotalLogsSent++
	q.Q(TotalLogsReceived, TotalLogsSent, mlen, string(buf[:mlen]))
	return nil
}

func SendSyslog(priority int, tag string, msg string) error {
	timestamp := time.Now().Format(time.Stamp) // XXX not RFC3339
	var name string
	if len(QC.Name) == 0 {
		name, _ = os.Hostname()
	} else {
		name = QC.Name
	}
	syslogmsg := fmt.Sprintf("<%d>%s %s %s: %s", priority, timestamp, name, tag, msg)
	if QC.RemoteSyslogServerAddr == "" {
		q.Q("Missing remote syslog server address, can't send syslog")
		rootSaveLog(syslogmsg)
		return fmt.Errorf("%v", "Missing remote syslog server address")
	}
	// reuse udp socket instead of open/close per message
	if QC.RemoteSyslogServer == nil {
		udpSock, err := net.Dial("udp4", QC.RemoteSyslogServerAddr)
		if err != nil {
			rootSaveLog(syslogmsg)

			return err
		}
		QC.RemoteSyslogServer = udpSock
		// TODO close udpSock when program exits.
	}

	_, err := QC.RemoteSyslogServer.Write([]byte(syslogmsg))
	if err != nil {
		q.Q(err)
		rootSaveLog(syslogmsg)
		return err
	}
	TotalLogsSent++
	q.Q("sent syslog", string(msg))
	return nil
}

func rootSaveLog(syslogmsg string) {
	if QC.IsRoot {
		_, severity, err := SyslogParsePriority(syslogmsg)
		if err != nil {
			q.Q(err)
		}
		SendSocketMessage(severity, syslogmsg)
		SaveLog(syslogmsg)
	} else {
		TotalLogsDropped++
		q.Q(TotalLogsDropped)
	}
}

var Logger *lumberjack.Logger

func initLogger() *lumberjack.Logger {
	filename := path.Join(QC.SyslogLocalPath)
	Logger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    int(QC.SyslogFileSize),
		MaxBackups: 10,
		Compress:   QC.SyslogCompress,
		LocalTime:  true,
	}
	return Logger
}

// Save syslog to file
func SaveLog(data string) {
	// mkdir()
	if Logger == nil {
		Logger = initLogger()
	} else {
		if Logger.Filename != (path.Join(QC.SyslogLocalPath)) || Logger.Compress != QC.SyslogCompress || Logger.MaxSize != int(QC.SyslogFileSize) {
			err := Logger.Close()
			if err != nil {
				q.Q(err)
			}
			Logger = initLogger()
			q.Q("SaveLog,change local syslog paramter:", Logger)
		}
	}
	re := regexp.MustCompile(`\r?\n`)
	data = re.ReplaceAllString(data, " ")
	_, err := Logger.Write([]byte(data))
	if err != nil {
		q.Q(err)
		//remind user if file error
		SendSocketMessage(LOG_ERR, fmt.Sprintf("can open:%v, please check file", Logger.Filename))
		return
	}
	_, err = Logger.Write([]byte("\n"))
	if err != nil {
		q.Q(err)
		//remind user if file error
		SendSocketMessage(LOG_ERR, fmt.Sprintf("can open:%v, please check file", Logger.Filename))
		return
	}
}

// Configure local syslog path.
//
// Usage : config local syslog path [path]
//
//	[path]        : local syslog path
//
// Example :
//
//	config local syslog path tmp/log
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

// Configure local syslog file maximum size.
//
// Usage : config local syslog maxsize [maxsize]
//
//	[maxsize]     : local syslog file maxsize size
//
// Example :
//
//	config local syslog maxsize 100
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
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	QC.SyslogFileSize = uint(s)
	cmdinfo.Status = "ok"
	return cmdinfo
}

// Whether to configure local syslog files to be compressed
//
// Usage : config local syslog compress [compress]
//
//	[compress]     : would be compressed
//
// Example :
//
//	config local syslog compress true
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
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	QC.SyslogCompress = boolValue
	cmdinfo.Status = "ok"
	return cmdinfo
}

const foramt = "2006/01/02 15:04:05"

// Read local syslog.
//
// Usage : config local syslog read [start date] [start time] [end date] [end time] [max line]
//
//	[start date]   : search syslog start date
//	[start time]   : search syslog start time
//	[end date]     : search syslog end date
//	[end time]     : search syslog end time
//	[max line]     : max lines, if without max line, that mean read all of lines
//
// Example :
//
//	config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00
//	config local syslog read 5
//	config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00 5
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
		// if inlcude max line paramters
		if len(ws) == 9 {
			v, err := strconv.Atoi(ws[8])
			if err != nil {
				q.Q(err)
				cmdinfo.Status = "error: " + err.Error()
				return cmdinfo
			} else {
				maxline = v
				q.Q("read syslog max line=", maxline)
			}
		}
	}
	// if inlcude max line paramters
	if len(ws) == 5 {
		v, err := strconv.Atoi(ws[4])
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: " + err.Error()
			return cmdinfo
		} else {
			maxline = v
			q.Q("read syslog max line=", maxline)
		}
	}

	readFile, err := os.Open(QC.SyslogLocalPath)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
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
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Result = string(b)
	q.Q(cmdinfo.Result)
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

// compareTime compare time size with start time and end time
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

func SendSocketMessage(severity int, bufStr string) {
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
}
