package mnms

import (
	"strconv"
	"strings"

	"github.com/qeof/q"
)

type Log struct {
	Kind     string   `json:"kind"`
	Messages []string `json:"messages"`
}

//TODO the log messages are appended to Messages slice which must be dumped
// to files regularly and kept to a limited size.
//TODO need ways to load dumped snapshots to view them, search, via api.

// log kinds: syslog, traps, alerts
//TODO implement syslog server, trap server, alert messaging, integrate q.Q debug logs

func init() {
	QC.Logs = make(map[string]Log)
}

func InsertLogKind(log *Log) {
	QC.Logs[log.Kind] = *log
}

// Configure log setting.
//
// Usage : log off
//
// Usage : log pattern [pattern]
//
//	[pattern]     : log pattern
//
// Usage : log output [output]
//
//	[output]      : log output
//
// Usage : log syslog [facility] [severity] [tag] [message]
//
//	[facility]    : syslog facility
//	[severity]    : syslog severity
//	[tag]         : syslog was sent from tag what feature name
//	[message]     : would send messages
//
// Example :
//
//	log off
//	log pattern .*
//	log output stderr
//	log syslog 0 1 InsertDev "new device"
func LogCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	if cmd == "log off" {
		//use q.Q stuff to control logs for debugging
		q.P = ""
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	if strings.HasPrefix(cmd, "log pattern ") {
		ws := strings.Split(cmd, " ")
		if len(ws) < 3 {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		q.P = ws[2]
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	if strings.HasPrefix(cmd, "log output ") {
		ws := strings.Split(cmd, " ")
		if len(ws) < 3 {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		q.O = ws[2]
		cmdinfo.Status = "ok"
		//TODO ways to retrieve different log output files
		return cmdinfo
	}
	if strings.HasPrefix(cmd, "log syslog ") {
		ws := strings.Split(cmd, " ")
		if len(ws) < 6 {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		val, err := strconv.Atoi(ws[2])
		if err != nil {
			cmdinfo.Status = "error: invalid facility value"
			return cmdinfo
		}
		facility := val
		if facility < LOG_KERN ||
			facility > LOG_LOCAL7 {
			cmdinfo.Status = "error: invalid facility"
			return cmdinfo
		}
		val, err = strconv.Atoi(ws[3])
		if err != nil {
			cmdinfo.Status = "error: invalid severity value"
			return cmdinfo
		}
		severity := val
		if severity < LOG_EMERG ||
			severity > LOG_DEBUG {
			cmdinfo.Status = "error: invalid severity"
			return cmdinfo
		}
		priority := facility*8 | severity
		tag := ws[4]
		msg := strings.Join(ws[5:], " ")
		err = SendSyslog(priority, tag, msg)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = "error: sending syslog"
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}
	cmdinfo.Status = "error: invalid command"
	return cmdinfo
}
