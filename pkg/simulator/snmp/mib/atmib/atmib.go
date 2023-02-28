package atmib

import (
	"encoding/hex"
	"net"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

//All oid of at
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, at(d)...)
	return result
}

func at(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.3.1.1.1",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(10000), nil
			},
			Document: "ifNumber",
		},
		{
			OID:  "1.3.6.1.2.1.3.1.1.2",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				targetStr := d.GetMac()
				decoded, err := hex.DecodeString(targetStr)
				if err != nil {
					return nil, err
				}
				return GoSNMPServer.Asn1OctetStringWrap(string(decoded)), nil
			},
			Document: "ifNumber",
		},
		{
			OID:  "1.3.6.1.2.1.3.1.1.3",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetIP())), nil
			},
			Document: "ifNumber",
		},
	}
	return v
}
