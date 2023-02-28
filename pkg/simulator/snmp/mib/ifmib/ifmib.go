package ifmib

import (
	"fmt"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of interfaces
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, ifXTable(d)...)
	return result
}

func ifXTable(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {

	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.1.%d", ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1OctetStringWrap(fmt.Sprintf("Port%d", ifIndex)), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.2.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.3.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.4.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.5.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.6.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.7.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.8.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.9.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.10.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.11.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.12.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.13.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.31.1.1.1.15.%d", ifIndex),
				Type: gosnmp.Gauge32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Gauge32Wrap(0), nil
				},
			},
		}
		toRet = append(toRet, currentIf...)
	}
	currentIf := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  ".1.3.6.1.2.1.31.1.1.1.18",
			Type: gosnmp.NoSuchObject,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.ErrNoSNMPInstance.Error(), nil
			},
		}}
	toRet = append(toRet, currentIf...)

	return toRet
}
