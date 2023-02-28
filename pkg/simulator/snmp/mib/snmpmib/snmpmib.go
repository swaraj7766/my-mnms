package snmpmib

import (
	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

//All oid of snmp
func All() []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.11.1",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(2124), nil
			},
			Document: "snmpInPkts",
		},
		{
			OID:  "1.3.6.1.2.1.11.2",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(2124), nil
			},
			Document: "snmpOutPkts",
		},
		{
			OID:  "1.3.6.1.2.1.11.3",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInBadVersions",
		},
		{
			OID:  "1.3.6.1.2.1.11.4",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(6), nil
			},
			Document: "snmpInBadCommunityNames",
		},
		{
			OID:  "1.3.6.1.2.1.11.5",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInBadCommunityUses",
		},
		{
			OID:  "1.3.6.1.2.1.11.6",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInASNParseErrs",
		},
		{
			OID:  "1.3.6.1.2.1.11.8",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInTooBigs",
		},
		{
			OID:  "1.3.6.1.2.1.11.9",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInNoSuchNames",
		},

		{
			OID:  "1.3.6.1.2.1.11.10",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInBadValues",
		},
		{
			OID:  "1.3.6.1.2.1.11.11",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInReadOnlys",
		},

		{
			OID:  "1.3.6.1.2.1.11.12",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInGenErrs",
		},

		{
			OID:  "1.3.6.1.2.1.11.13",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(85043), nil
			},
			Document: "snmpInTotalReqVars",
		},
		{
			OID:  "1.3.6.1.2.1.11.14",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInTotalSetVars",
		},
		{
			OID:  "1.3.6.1.2.1.11.15",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(130), nil
			},
			Document: "snmpInGetRequests",
		},
		{
			OID:  "1.3.6.1.2.1.11.16",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(278), nil
			},
			Document: "snmpInGetNexts",
		},

		{
			OID:  "1.3.6.1.2.1.11.17",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(2), nil
			},
			Document: "snmpInSetRequests",
		},
		{
			OID:  "1.3.6.1.2.1.11.18",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInGetResponses",
		},
		{
			OID:  "1.3.6.1.2.1.11.19",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpInTraps",
		},
		{
			OID:  "1.3.6.1.2.1.11.20",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpOutTooBigs",
		},
		{
			OID:  "1.3.6.1.2.1.11.21",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(2), nil
			},
			Document: "snmpOutNoSuchNames",
		},
		{
			OID:  "1.3.6.1.2.1.11.22",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpOutBadValues",
		},

		{
			OID:  "1.3.6.1.2.1.11.24",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpOutGenErrs",
		},
		{
			OID:  "1.3.6.1.2.1.11.25",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpOutGetRequests",
		},
		{
			OID:  "1.3.6.1.2.1.11.26",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpOutGetNexts",
		},
		{
			OID:  "1.3.6.1.2.1.11.27",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "snmpOutSetRequests",
		},
		{
			OID:  "1.3.6.1.2.1.11.28",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(2114), nil
			},
			Document: "snmpOutGetResponses",
		},
		{
			OID:  "1.3.6.1.2.1.11.29",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(6), nil
			},
			Document: "snmpOutTraps",
		},
		{
			OID:  "1.3.6.1.2.1.11.30",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(1), nil
			},
			Document: "snmpEnableAuthenTraps",
		},
		{
			OID:  "1.3.6.1.2.1.11.31.0",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
		},
		{
			OID:  "1.3.6.1.2.1.11.32.0",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
		},
	}

	return v
}
