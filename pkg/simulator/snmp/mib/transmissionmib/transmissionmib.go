package transmissionmib

import (
	"fmt"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of transmission
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, othersOIDs(d)...)
	return result
}

func othersOIDs(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{}
	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.2.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.3.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.4.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.5.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.6.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.7.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.8.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.9.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.10.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.11.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.13.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.16.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.10.7.2.1.17.%d", ifIndex),
				Type: gosnmp.ObjectIdentifier,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1ObjectIdentifierWrap("0.0"), nil
				},
			},
		}
		toRet = append(toRet, currentIf...)
	}
	return toRet
}
