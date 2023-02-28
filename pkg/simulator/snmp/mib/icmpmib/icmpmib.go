package icmpmib

import (
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

//All oid of icmp
func All() []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, icmp()...)

	return result
}

func icmp() []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:      "1.3.6.1.2.1.5.1",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(20), nil },
			Document: "icmpInMsgs",
		},
		{
			OID:      "1.3.6.1.2.1.5.2",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(5), nil },
			Document: "icmpInErrors",
		},
		{
			OID:      "1.3.6.1.2.1.5.3",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInDestUnreachs",
		},
		{
			OID:      "1.3.6.1.2.1.5.4",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInTimeExcds",
		},
		{
			OID:      "1.3.6.1.2.1.5.5",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInParmProbs",
		},
		{
			OID:      "1.3.6.1.2.1.5.6",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInSrcQuenchs",
		},
		{
			OID:      "1.3.6.1.2.1.5.7",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInRedirects",
		},
		{
			OID:      "1.3.6.1.2.1.5.8",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(21), nil },
			Document: "icmpInEchos",
		},
		{
			OID:      "1.3.6.1.2.1.5.9",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(4), nil },
			Document: "icmpInEchoReps",
		},
		{
			OID:      "1.3.6.1.2.1.5.10",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInTimestamps",
		},
		{
			OID:      "1.3.6.1.2.1.5.11",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInTimestampReps",
		},
		{
			OID:      "1.3.6.1.2.1.5.12",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInAddrMasks",
		},
		{
			OID:      "1.3.6.1.2.1.5.13",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpInAddrMaskReps",
		},
		{
			OID:      "1.3.6.1.2.1.5.14",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutMsgs",
		},
		{
			OID:      "1.3.6.1.2.1.5.15",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutErrors",
		},
		{
			OID:      "1.3.6.1.2.1.5.16",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutDestUnreachs",
		},
		{
			OID:      "1.3.6.1.2.1.5.17",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutTimeExcds",
		},
		{
			OID:      "1.3.6.1.2.1.5.18",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutParmProbs",
		},
		{
			OID:      "1.3.6.1.2.1.5.19",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutSrcQuenchs",
		},
		{
			OID:      "1.3.6.1.2.1.5.20",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutRedirects",
		},
		{
			OID:      "1.3.6.1.2.1.5.21",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(4), nil },
			Document: "icmpOutEchos",
		},
		{
			OID:      "1.3.6.1.2.1.5.22",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(16), nil },
			Document: "icmpOutEchoReps",
		},
		{
			OID:      "1.3.6.1.2.1.5.23",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutTimestamps",
		},
		{
			OID:      "1.3.6.1.2.1.5.24",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutTimestampReps",
		},
		{
			OID:      "1.3.6.1.2.1.5.25",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutAddrMasks",
		},
		{
			OID:      "1.3.6.1.2.1.5.26",
			Type:     gosnmp.Counter32,
			OnGet:    func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			Document: "icmpOutAddrMaskReps",
		},
	}
	return v
}
