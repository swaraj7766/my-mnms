package snmp

import (
	"net"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"mnms/pkg/simulator/snmp/mib/atmib"
	"mnms/pkg/simulator/snmp/mib/atopmib"
	dot1dbridgemib "mnms/pkg/simulator/snmp/mib/dot1dBridgemib"
	"mnms/pkg/simulator/snmp/mib/icmpmib"
	"mnms/pkg/simulator/snmp/mib/ifmib"
	"mnms/pkg/simulator/snmp/mib/interfacesmib"
	"mnms/pkg/simulator/snmp/mib/ipmib"
	"mnms/pkg/simulator/snmp/mib/rmonmib"
	"mnms/pkg/simulator/snmp/mib/snmpmib"
	"mnms/pkg/simulator/snmp/mib/systemmib"
	"mnms/pkg/simulator/snmp/mib/tcpmib"
	"mnms/pkg/simulator/snmp/mib/transmissionmib"
	"mnms/pkg/simulator/snmp/mib/udpmib"

	"github.com/sirupsen/logrus"
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/GoSNMPServer/mibImps/dismanEventMib"
	"github.com/slayercat/GoSNMPServer/mibImps/ucdMib"
)

const port = "161"

type Snmp struct {
	agent  GoSNMPServer.MasterAgent
	server *GoSNMPServer.SNMPServer
	data   *snmpvalue.Value
}

func NewSnmp(community []string, data *snmpvalue.Value) *Snmp {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	toRet = append(toRet, ucdMib.All()...) //dismanEventMib
	toRet = append(toRet, dismanEventMib.All()...)
	toRet = append(toRet, systemmib.All(data)...)
	toRet = append(toRet, interfacesmib.All(data)...)
	toRet = append(toRet, atmib.All(data)...)
	toRet = append(toRet, ipmib.All(data)...)
	toRet = append(toRet, icmpmib.All()...)
	toRet = append(toRet, tcpmib.All()...)
	toRet = append(toRet, udpmib.All()...)
	toRet = append(toRet, transmissionmib.All(data)...)
	toRet = append(toRet, snmpmib.All()...)
	toRet = append(toRet, rmonmib.All(data)...)
	toRet = append(toRet, dot1dbridgemib.All(data)...)
	toRet = append(toRet, ifmib.All(data)...)
	toRet = append(toRet, atopmib.All(data)...)
	master := GoSNMPServer.MasterAgent{
		//Logger: GoSNMPServer.NewDefaultLogger(),
		SecurityConfig: GoSNMPServer.SecurityConfig{
			AuthoritativeEngineBoots: 1,
		},
		SubAgents: []*GoSNMPServer.SubAgent{
			{
				CommunityIDs:        community,
				OIDs:                toRet,
				UserErrorMarkPacket: true,
			},
		},
	}
	snmp := &Snmp{agent: master, data: data}

	return snmp
}

func (s *Snmp) Run(ip string) error {
	addr := net.JoinHostPort(ip, port)
	s.server = GoSNMPServer.NewSNMPServer(s.agent)
	err := s.server.ListenUDP("udp", addr)
	if err != nil {
		logrus.Fatal(err)
	}

	return s.server.ServeForever()
}
func (s *Snmp) SetLogger(log *logrus.Logger) {
	s.agent.Logger = log
}

func (s *Snmp) Shutdown() {
	s.server.Shutdown()
}

func (s *Snmp) GetData() *snmpvalue.Value {
	return s.data
}
