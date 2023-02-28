package atopmib

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"mnms/pkg/simulator/devicetype"
	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of dot1dbridge
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, atopinfo(d)...)
	result = append(result, end()...)
	return result
}
func defaultAtopInfo() atopvalue {
	val := atopvalue{
		syslogStatus:                  2,
		backupAgentBoardFwFileName:    "test_string.dld",
		restoreAgentBoardFwFileName:   "test_string.dld",
		backupStatus:                  2,
		restoreStatus:                 2,
		eventServerPort:               514,
		eventServerLevel:              3,
		eventLogToFlash:               2,
		sntpClientStatus:              2,
		sntpUTCTimezone:               47,
		sntpServer1:                   "time.nist.gov",
		sntpServer2:                   "time-A.timefreq.bldrdoc.gov",
		sntpServerQueryPeriod:         259200,
		agingTimeSetting:              2,
		ptpState:                      1,
		ptpVersion:                    1,
		ptpSyncInterval:               1,
		ptpClockStratum:               1,
		ptpPriority1:                  1,
		ptpPriority2:                  1,
		rstpStatus:                    1,
		qosCOSPriorityQueue:           1,
		qosTOSPriorityQueue:           1,
		eventPortEventEmail:           1,
		eventPortEventRelay:           1,
		eventPowerEventSMTP1:          1,
		eventPowerEventSMTP2:          1,
		syslogEventsSMTP:              8,
		eventEmailAlertAddr:           "test@test.com",
		eventEmailAlertAuthentication: 1,
		eventEmailAlertAccount:        "test_account",
		lldpStatus:                    1,
		trapServerStatus:              1,
		trapServerIP:                  "",
		trapServerPort:                8080,
		trapServerTrapComm:            "",
	}

	return val
}

func atopinfo(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	atopvalue := defaultAtopInfo()

	v, _ := devicetype.ParsingType(d.GetModel())
	oid := OidType(v)
	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.4.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return d.GetModel(), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.5.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return d.GetKernel(), nil },
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.6.0", oid),
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				targetStr := d.GetMac()
				decoded, err := hex.DecodeString(targetStr)
				if err != nil {
					return nil, err
				}
				return GoSNMPServer.Asn1OctetStringWrap(string(decoded)), nil
			},
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.7.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return d.GetKernel(), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.8.1.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.8.2.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.8.3.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.8.4.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.8.5.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.8.6.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.9.1.1.1.2", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.9.1.1.2.1", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.9.1.1.2.2", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.1.10.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return d.GetModel(), nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					d.SetModel(sb.String())
					return nil

				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "systemModelName",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.2.1.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return d.GetUser(), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.2.2.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return "", nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					if len(r) > 6 {
						d.SetPwd(sb.String())
						return nil
					}
					return fmt.Errorf("%s", gosnmp.WrongLength.String())
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.1.1", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.2.1", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.3.1", oid),
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetIP())), nil
			},
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.4.1", oid),
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetMask())), nil
			},
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.5.1", oid),
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetGateWay())), nil
			},
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.6.1", oid),
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP("0.0.0.0")), nil
			},
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.3.1.1.7.1", oid),
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP("0.0.0.0")), nil
			},
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.4.1.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(atopvalue.sntpClientStatus), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.sntpClientStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "sntpClientStatus",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.4.3.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(atopvalue.sntpUTCTimezone), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.sntpUTCTimezone = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "sntpUTCTimezone",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.4.5.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return "2000-08-02,00:30:50", nil },
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.4.9.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.sntpServer1, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.sntpServer1 = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "sntpServer1",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.4.10.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.sntpServer2, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.sntpServer2 = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "sntpServer2",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.4.11.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.sntpServerQueryPeriod, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.sntpServerQueryPeriod = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "sntpServerQueryPeriod",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.1.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.backupServerIP, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.backupServerIP = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "backupServerIP",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.2.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.backupAgentBoardFwFileName, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.backupAgentBoardFwFileName = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "backupAgentBoardFwFileName",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.3.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.backupStatus, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.backupStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "backupStatus",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.4.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.restoreServerIP, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.restoreServerIP = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "restoreServerIP",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.5.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.restoreAgentBoardFwFileName, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.restoreAgentBoardFwFileName = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "restoreAgentBoardFwFileName",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.6.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.restoreStatus, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.restoreStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "restoreStatus",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.6.7.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},

		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.11.1.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.agingTimeSetting, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.agingTimeSetting = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "agingTimeSetting",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.12.2.1.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.ptpState, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.ptpState = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "ptpState",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.12.2.2.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.ptpVersion, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.ptpVersion = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "ptpVersion",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.12.2.3.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.ptpSyncInterval, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.ptpSyncInterval = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "ptpSyncInterval",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.12.2.5.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.ptpClockStratum, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.ptpClockStratum = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "ptpClockStratum",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.12.2.6.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.ptpPriority1, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.ptpPriority1 = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "ptpPriority1",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.12.2.7.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.ptpPriority2, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.ptpPriority2 = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "ptpPriority2",
		},

		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.4.2.1.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.rstpStatus, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.rstpStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "rstpStatus",
		},

		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.6.4.1.3.1", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.qosCOSPriorityQueue, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.qosCOSPriorityQueue = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "qosCOSPriorityQueue",
		},

		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.6.6.1.3.1", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.qosTOSPriorityQueue, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.qosTOSPriorityQueue = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "qosTOSPriorityQueue",
		},
	}

	toRet = append(toRet, eventPortNumber(oid, d)...)
	toRet = append(toRet, eventPortEventEmail(oid, d, atopvalue)...)
	toRet = append(toRet, eventPortEventRelay(oid, d, atopvalue)...)
	toRet = append(toRet, eventPowerNumber(oid)...)
	toRet = append(toRet, eventPowerEventSMTP(oid, atopvalue)...)
	toRet = append(toRet, eventPowerEventRelay(oid)...)
	toRet = append(toRet, syslogEventsSMTP(oid, atopvalue)...)
	toRet = append(toRet, syslogStatus(oid, atopvalue)...)

	toRet = append(toRet, eventEmailAlert(oid, atopvalue)...)
	toRet = append(toRet, swCurrentPortNameListPortName(oid, d)...)
	return toRet
}

// swCurrentPortNameListPortName
func swCurrentPortNameListPortName(oid uint, d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		event := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.10.1.2.%d", oid, ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1OctetStringWrap(fmt.Sprintf("Port%d", ifIndex-1)), nil
				},
			},
			{
				OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.2.10.1.3.%d", oid, ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1OctetStringWrap(fmt.Sprintf("%d", ifIndex-1)), nil
				},
			},
		}
		toRet = append(toRet, event...)
	}
	return toRet
}

// eventPortNumber   //
func eventPortNumber(oid uint, d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		event := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.10.1.1.2.1.1.%d", oid, ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
		}
		toRet = append(toRet, event...)
	}
	return toRet
}

// eventPortEventEmail   //
func eventPortEventEmail(oid uint, d *snmpvalue.Value, atopvalue atopvalue) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		event := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.10.1.1.2.1.3.%d", oid, ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return atopvalue.eventPortEventEmail, nil },
				OnSet: func(value interface{}) error {
					val, ok := value.(int)
					if ok {

						atopvalue.eventPortEventEmail = int(val)
						return nil
					}
					return fmt.Errorf("%s", gosnmp.WrongType.String())
				},
				Document: "eventPortEventEmail",
			},
		}
		toRet = append(toRet, event...)
	}
	return toRet
}

// eventPortEventRelay   //
func eventPortEventRelay(oid uint, d *snmpvalue.Value, atopvalue atopvalue) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		event := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%v.10.1.1.2.1.4.%d", oid, ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return atopvalue.eventPortEventRelay, nil },
				OnSet: func(value interface{}) error {
					val, ok := value.(int)
					if ok {

						atopvalue.eventPortEventRelay = int(val)
						return nil
					}
					return fmt.Errorf("%s", gosnmp.WrongType.String())
				},
				Document: "eventPortEventRelay",
			},
		}
		toRet = append(toRet, event...)
	}
	return toRet
}

// Name/OID: eventPowerNumber.1; Value (Integer): 1
func eventPowerNumber(oid uint) []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:      fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.3.1.1.1", oid),
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
			Document: "eventPowerNumber.1",
		},
		{
			OID:      fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.3.1.1.2", oid),
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
			Document: "eventPowerNumber.2",
		},
	}

	return event
}

// Name/OID: eventPowerEventSMTP.1; Value (Integer): disabled (3)
func eventPowerEventSMTP(oid uint, atopvalue atopvalue) []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.3.1.3.1", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.eventPowerEventSMTP1, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.eventPowerEventSMTP1 = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},

			Document: "eventPowerEventSMTP.1",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.3.1.3.2", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.eventPowerEventSMTP2, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.eventPowerEventSMTP2 = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},

			Document: "eventPowerEventSMTP.2",
		},
	}

	return event
}

// Name/OID: eventPowerEventRelay.1; Value (Integer): disabled (3)
func eventPowerEventRelay(oid uint) []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:      fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.3.1.4.1", oid),
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(3), nil },
			Document: "eventPowerEventRelay.1",
		},
		{
			OID:      fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.3.1.4.2", oid),
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(3), nil },
			Document: "eventPowerEventRelay.2",
		},
	}

	return event
}

// Name/OID: syslogEventsSMTP.0; Value (Integer): disabled (8)
func syslogEventsSMTP(oid uint, atopvalue atopvalue) []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.1.4.1.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.syslogEventsSMTP, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.syslogEventsSMTP = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},

			Document: "syslogEventsSMTP",
		},
	}

	return event
}

func trapServer(oid uint, atopvalue atopvalue) []*GoSNMPServer.PDUValueControlItem {

	// snmpTrapServerTrapComm:
	//   oid: ".8.6.1.3.0"
	//   type: 4
	// trapServerStatus:              1,
	// trapServerIP:                  "",
	// trapServerPort:                8080,
	// trapServerTrapComm:            "",
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			// status
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.8.6.1.5.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.trapServerStatus, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.trapServerStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "trapServerStatus",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.8.6.1.7.0", oid),
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return atopvalue.trapServerIP, nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(string)
				if ok {

					atopvalue.trapServerIP = val
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "trapServerIP",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.8.6.1.6.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return atopvalue.trapServerPort, nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {

					atopvalue.trapServerPort = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "trapServerPort",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.8.6.1.3.0", oid),
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return atopvalue.trapServerTrapComm, nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(string)
				if ok {

					atopvalue.trapServerTrapComm = val
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "trapServerTrapComm",
		},
	}
	return event
}

// Name/OID: syslogStatus.0; Value (Integer): disabled (2)
func syslogStatus(oid uint, v atopvalue) []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.2.1.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(v.syslogStatus), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {
					v.syslogStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "syslogStatus",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.2.3.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(v.eventServerPort), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {
					v.eventServerPort = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "eventServerPort",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.2.4.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(v.eventServerLevel), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {
					v.eventServerLevel = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "eventServerLevel",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.2.5.0", oid),
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(v.eventLogToFlash), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {
					v.eventLogToFlash = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "eventLogToFlash",
		},
		{
			OID:  fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.2.6.0", oid),
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringUnwrap(v.eventServerIP), nil
			},
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					v.eventServerIP = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "eventServerIP",
		},
	}

	return event
}

// Name/OID: eventEmailAlert
func eventEmailAlert(oid uint, atopvalue atopvalue) []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.3.2.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.eventEmailAlertAddr, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.eventEmailAlertAddr = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},

			Document: "eventEmailAlertAddr",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.3.3.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.eventEmailAlertAuthentication, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {
					atopvalue.eventEmailAlertAuthentication = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},

			Document: "eventEmailAlertAuthentication",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.10.1.3.4.0", oid),
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return atopvalue.eventEmailAlertAccount, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					r := sb.String()
					atopvalue.eventEmailAlertAccount = r
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "eventEmailAlertAccount",
		},
		{
			OID:   fmt.Sprintf(".1.3.6.1.4.1.3755.0.0.%d.12.1.0", oid),
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return atopvalue.lldpStatus, nil },
			OnSet: func(value interface{}) error {
				val, ok := value.(int)
				if ok {
					atopvalue.lldpStatus = int(val)
					return nil
				}
				return fmt.Errorf("%s", gosnmp.WrongType.String())
			},
			Document: "lldpStatus",
		},
	}

	return event
}

// Name/OID: eventEmailAlert
func end() []*GoSNMPServer.PDUValueControlItem {
	event := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   fmt.Sprint(".1.3.6.1.6.3.16.1.5.2.1.6.4.109.105.98.50.6.1.3.6.1.2.1"),
			Type:  gosnmp.Null,
			OnGet: func() (value interface{}, err error) { return nil, nil },
		},
	}

	return event
}
