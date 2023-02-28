package mnms

import (
	"net"
	"sync"

	"github.com/gorilla/websocket"
	cron "github.com/robfig/cron/v3"
)

// global context holder
type QContext struct {
	DevMutex                  sync.Mutex
	Name                      string
	Port                      int
	Root                      string
	IsRoot                    bool
	DevData                   map[string]DevInfo
	CmdMutex                  sync.Mutex
	CmdData                   map[string]CmdInfo
	ClientMutex               sync.Mutex
	Clients                   map[string]string
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
	ArpInterval               int
	CmdInterval               int
	CronJobs                  []CronInfo
	Cron                      *cron.Cron
	Domain                    string
	TopologyData              map[string]Topology
	AdminToken                string
	OwnPublicKeys             []byte
}

var QC QContext

func init() {
	QC.Clients = make(map[string]string) // list of non-root mnms services registered to root
	QC.Port = 27182                      // euler's number
	QC.MqttBrokerAddr = ":11883"         // ":1883"
	// QC.SyslogServerAddrTcp = ":1468"     // ":1468"
	QC.SyslogServerAddr = ":5514" // ":514"
	QC.TrapServerAddr = ":5162"   // ":162"
	QC.WebSocketMessageBroadcast = make(chan WebSocketMessage, 100)
	QC.WebSocketClient = make(map[*websocket.Conn]bool)
	QC.TopologyData = make(map[string]Topology)
	QC.ArpInterval = 10 // XXX
	QC.CmdInterval = 5  // XXX
	QC.SyslogLocalPath = SyslogPath
	QC.SyslogFileSize = 100 //megabytes
	QC.SyslogCompress = true
}
