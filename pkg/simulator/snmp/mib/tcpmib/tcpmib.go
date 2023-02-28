package tcpmib

import (
	"net"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

//All oid of tcp
func All() []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.6.1",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(1), nil
			},
			Document: "tcpRtoAlgorithm",
		},
		{
			OID:  "1.3.6.1.2.1.6.2",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(200), nil
			},
			Document: "tcpRtoMin",
		},
		{
			OID:  "1.3.6.1.2.1.6.3",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(120000), nil
			},
			Document: "tcpRtoMax",
		},
		{
			OID:  "1.3.6.1.2.1.6.4",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(-1), nil
			},
			Document: "tcpMaxConn",
		},
		{
			OID:      "1.3.6.1.2.1.6.5",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "tcpActiveOpens",
		},
		{
			OID:      "1.3.6.1.2.1.6.6",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(81), nil },
			Document: "tcpPassiveOpens",
		},
		{
			OID:      "1.3.6.1.2.1.6.7",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "tcpAttemptFails",
		},
		{
			OID:      "1.3.6.1.2.1.6.8",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(22), nil },
			Document: "tcpEstabResets",
		},
		{
			OID:      "1.3.6.1.2.1.6.9",
			Type:     gosnmp.Gauge32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(0), nil },
			Document: "tcpCurrEstab",
		},
		{
			OID:      "1.3.6.1.2.1.6.10",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(844), nil },
			Document: "tcpInSegs",
		},

		{
			OID:      "1.3.6.1.2.1.6.11",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(720), nil },
			Document: "tcpOutSegs",
		},
		{
			OID:      "1.3.6.1.2.1.6.12",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "tcpRetransSegs",
		},
		{
			OID:      "1.3.6.1.2.1.6.13.1.1",
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
			Document: "tcpConnState",
		},
		{
			OID:  "1.3.6.1.2.1.6.13.1.2",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP("0.0.0.0")), nil
			},
			Document: "tcpConnLocalAddress",
		},
		{
			OID:      "1.3.6.1.2.1.6.13.1.3",
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(22), nil },
			Document: "tcpConnLocalPort",
		},
		{
			OID:  "1.3.6.1.2.1.6.13.1.4",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP("0.0.0.0")), nil
			},
			Document: "tcpConnRemAddress",
		},
		{
			OID:      "1.3.6.1.2.1.6.13.1.5",
			Type:     gosnmp.Integer,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			Document: "tcpConnRemPort",
		},
		{
			OID:      "1.3.6.1.2.1.6.14",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "tcpInErrs",
		},
		{
			OID:      "1.3.6.1.2.1.6.15",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "tcpOutRsts",
		},
	}
	return v
}
