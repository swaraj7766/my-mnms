package rmonmib

import (
	"fmt"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of rmon
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, statistics(d)...)
	result = append(result, history()...)
	result = append(result, alarm()...)
	result = append(result, event()...)
	return result
}

func statistics(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.2.%d", ifIndex),
				Type: gosnmp.ObjectIdentifier,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1ObjectIdentifierWrap(fmt.Sprintf("1.3.6.1.2.1.2.2.1.1.%d", ifIndex-1)), nil
				},
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.3.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.4.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.5.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.6.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.7.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.8.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.9.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.10.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.11.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.12.%d", ifIndex),
				Type:  gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("Startup Mgmt"), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.1.1.13.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.4.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.16.1.4.1.2.%d", ifIndex),
				Type:  gosnmp.TimeTicks,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1TimeTicksWrap(uint32(5)), nil },
			},
		}
		toRet = append(toRet, currentIf...)
	}
	return toRet
}

func history() []*GoSNMPServer.PDUValueControlItem {

	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.1.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.2.1",
			Type:  gosnmp.ObjectIdentifier,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1ObjectIdentifierWrap("0.0"), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.3.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(30), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.4.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(30), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.5.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(180), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.6.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.2.1.1.7.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},

		{
			OID:   "1.3.6.1.2.1.16.2.2.1.1.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
	}

	return toRet
}

func alarm() []*GoSNMPServer.PDUValueControlItem {

	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.1.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.2.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(10), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.3.1",
			Type:  gosnmp.ObjectIdentifier,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1ObjectIdentifierWrap("0.0"), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.4.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.5.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.6.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.7.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},

		{
			OID:   "1.3.6.1.2.1.16.3.1.1.8.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.3.1.1.9.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},

		{
			OID:   "1.3.6.1.2.1.16.3.1.1.10.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},

		{
			OID:   "1.3.6.1.2.1.16.3.1.1.11.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
		},

		{
			OID:   "1.3.6.1.2.1.16.3.1.1.12.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
	}

	return toRet
}

func event() []*GoSNMPServer.PDUValueControlItem {

	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.1.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.2.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("Alarm"), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.3.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.4.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("public"), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.5.1",
			Type:  gosnmp.TimeTicks,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1TimeTicksWrap(uint32(0)), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.6.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap(""), nil },
		},
		{
			OID:   "1.3.6.1.2.1.16.9.1.1.7.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
	}

	return toRet
}
