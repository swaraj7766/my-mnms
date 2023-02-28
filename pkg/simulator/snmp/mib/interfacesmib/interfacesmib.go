package interfacesmib

import (
	"encoding/hex"
	"fmt"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of interfaces
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, ifNumber(d)...)
	result = append(result, othersOIDs(d)...)
	return result
}

func ifNumber(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.2.1",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(d.Port()), nil
			},
			Document: "ifNumber",
		},
	}
	return v
}

func othersOIDs(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.1.%d", ifIndex),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
				Document: "ifIndex",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.2.%d", ifIndex),
				Type:     gosnmp.OctetString,
				OnGet:    func() (value interface{}, err error) { return fmt.Sprintf("Port%d", ifIndex-1), nil },
				Document: "ifDescr",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.3.%d", ifIndex),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return 6, nil },
				Document: "ifType",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.4.%d", ifIndex),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return 1530, nil },
				Document: "ifMtu",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.5.%d", ifIndex),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return 1000, nil },
				Document: "ifSpeed",
			},

			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.2.2.1.6.%d", ifIndex),
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
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.7.%d", ifIndex),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return 1, nil },
				Document: "ifAdminStatus",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.8.%d", ifIndex),
				Type:     gosnmp.Integer,
				OnGet:    func() (value interface{}, err error) { return 2, nil },
				Document: "ifOperStatus",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.9.%d", ifIndex),
				Type:     gosnmp.TimeTicks,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1TimeTicksWrap(uint32(0)), nil },
				Document: "ifLastChange",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.10.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifInOctets",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.11.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifInUcastPkts",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.12.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifInNUcastPkts",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.13.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifInDiscards",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.14.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifInErrors",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.15.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifInUnknownProtos",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.16.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifOutOctets",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.17.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifOutUcastPkts",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.18.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifOutNUcastPkts",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.19.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifOutDiscards",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.20.%d", ifIndex),
				Type:     gosnmp.Counter32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
				Document: "ifOutErrors",
			},
			{
				OID:      fmt.Sprintf("1.3.6.1.2.1.2.2.1.21.%d", ifIndex),
				Type:     gosnmp.Gauge32,
				OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(0), nil },
				Document: "ifOutQLen",
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.2.2.1.22.%d", ifIndex),
				Type: gosnmp.ObjectIdentifier,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1ObjectIdentifierWrap("0.0"), nil
				},
				Document: "ifSpecific",
			},
		}
		toRet = append(toRet, currentIf...)
	}
	return toRet
}
