package udpmib

import (
	"net"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

//All oid of udp
func All() []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.7.1",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(22222), nil
			},
			Document: "udpInDatagrams",
		},
		{
			OID:  "1.3.6.1.2.1.7.2",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "udpNoPorts",
		},
		{
			OID:  "1.3.6.1.2.1.7.3",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(199), nil
			},
			Document: "udpInErrors",
		},
		{
			OID:  "1.3.6.1.2.1.7.4",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(5990), nil
			},
			Document: "udpOutDatagrams",
		},
		{
			OID:  "1.3.6.1.2.1.7.5.1.1",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP("0.0.0.0")), nil
			},
			Document: "udpLocalAddress",
		},
		{
			OID:      "1.3.6.1.2.1.7.5.1.2",
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(55954), nil },
			Document: "udpLocalPort",
		},
	}
	return v
}
