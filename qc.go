package mnms

import (
	"bytes"
	"encoding/json"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/gosnmp/gosnmp"
	"github.com/qeof/q"
)

// global context holder
type QContext struct {
	DevMutex                  sync.Mutex
	Name                      string
	Port                      int
	RootURL                   string
	IsRoot                    bool
	DumpStackTrace            bool
	DevData                   map[string]DevInfo
	CmdMutex                  sync.Mutex
	CmdData                   map[string]CmdInfo
	ClientMutex               sync.Mutex
	Clients                   map[string]ClientInfo
	Logs                      map[string]Log
	RemoteSyslogServer        net.Conn
	RemoteSyslogServerAddr    string
	SyslogLocalPath           string
	SyslogFileSize            uint
	SyslogCompress            bool
	MqttBrokerAddr            string
	SyslogServerAddr          string
	TrapServerAddr            string
	WebSocketClient           map[*websocket.Conn]bool
	WebSocketMessageBroadcast chan WebSocketMessage
	CmdInterval               int
	RegisterInterval          int
	GwdInterval               int
	Domain                    string
	TopologyData              map[string]Topology
	AdminToken                string
	OwnPublicKeys             []byte
	SnmpOptions               SnmpOptions
}

var QC QContext

func init() {
	QC.Clients = make(map[string]ClientInfo) // list of non-root mnms services registered to root
	QC.Port = 27182                          // euler's number
	QC.MqttBrokerAddr = ":11883"             // ":1883"
	QC.SyslogServerAddr = ":5514"            // ":514"
	QC.TrapServerAddr = ":5162"              // ":162"
	QC.WebSocketMessageBroadcast = make(chan WebSocketMessage, 100)
	QC.WebSocketClient = make(map[*websocket.Conn]bool)
	QC.TopologyData = make(map[string]Topology)
	QC.CmdInterval = 5
	QC.RegisterInterval = 60
	QC.GwdInterval = 60
	QC.SyslogLocalPath = "syslog_mnms.log"
	QC.SyslogFileSize = 100 //megabytes
	QC.SyslogCompress = true
	QC.SnmpOptions = SnmpOptions{
		Community: "private",
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Port:      161,
	}
}

// ClientInfo contains data about a client node instance in a cluster.
type ClientInfo struct {
	Name            string
	NumDevices      int
	NumCmds         int
	NumLogsReceived int
	NumLogsSent     int
	Start           int
	Now             int
	NumGoroutines   int
	IPAddresses     []string
}

func RegisterMain() {
	startTime := time.Now().Unix()

	for {
		ips, err := GetLocalIP()
		if err != nil {
			ips = []string{"Unknown"}
		}
		ci := ClientInfo{
			Name:            QC.Name,
			NumDevices:      len(QC.DevData),
			NumCmds:         len(QC.CmdData),
			NumLogsReceived: TotalLogsReceived,
			NumLogsSent:     TotalLogsSent,
			Start:           int(startTime),
			Now:             int(time.Now().Unix()),
			NumGoroutines:   runtime.NumGoroutine(),
			IPAddresses:     ips,
		}
		jsonBytes, err := json.Marshal(ci)
		if err != nil {
			q.Q(err)
			continue
		}
		resp, err := PostWithToken(QC.RootURL+"/api/v1/register",
			QC.AdminToken, bytes.NewBuffer(jsonBytes))
		if err != nil {
			q.Q(err)
		}
		if resp != nil {
			// save close, resp should not be nil here
			resp.Body.Close()
		}
		time.Sleep(time.Duration(QC.RegisterInterval) * time.Second) // XXX
	}
}
