package ipmib

import (
	"net"
	"strconv"

	atopnet "mnms/pkg/simulator/net"
	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

//All oid of ip
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, ip(d)...)

	return result
}

func ip(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.4.1",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(2), nil
			},
			Document: "ipForwarding",
		},
		{
			OID:  "1.3.6.1.2.1.4.2",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(64), nil
			},
			Document: "ipDefaultTTL",
		},
		{
			OID:  "1.3.6.1.2.1.4.3",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(675659), nil
			},
			Document: "ipInReceives",
		},
		{
			OID:  "1.3.6.1.2.1.4.4",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipInHdrErrors",
		},
		{
			OID:  "1.3.6.1.2.1.4.5",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(75389), nil
			},
			Document: "ipInAddrErrors",
		},
		{
			OID:  "1.3.6.1.2.1.4.6",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipForwDatagrams",
		},
		{
			OID:  "1.3.6.1.2.1.4.7",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipInUnknownProtos",
		},
		{
			OID:  "1.3.6.1.2.1.4.8",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipInDiscards",
		},
		{
			OID:  "1.3.6.1.2.1.4.9",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(186983), nil
			},
			Document: "ipInDelivers",
		},
		{
			OID:  "1.3.6.1.2.1.4.10",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(6590), nil
			},
			Document: "ipOutRequests",
		},
		{
			OID:  "1.3.6.1.2.1.4.11",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipOutDiscards",
		},
		{
			OID:  "1.3.6.1.2.1.4.12",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipOutNoRoutes",
		},
		{
			OID:  "1.3.6.1.2.1.4.13",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(0), nil
			},
			Document: "ipReasmTimeout",
		},
		{
			OID:  "1.3.6.1.2.1.4.14",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipReasmReqds",
		},
		{
			OID:  "1.3.6.1.2.1.4.15",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipReasmOKs",
		},
		{
			OID:  "1.3.6.1.2.1.4.16",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipReasmFails",
		},
		{
			OID:  "1.3.6.1.2.1.4.17",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(13), nil
			},
			Document: "ipFragOKs",
		},
		{
			OID:  "1.3.6.1.2.1.4.18",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(0), nil
			},
			Document: "ipFragFails",
		},
		{
			OID:  "1.3.6.1.2.1.4.19",
			Type: gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1Counter32Wrap(26), nil
			},
			Document: "ipFragCreates",
		},

		{
			OID:  "1.3.6.1.2.1.4.20.1.1",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetIP())), nil
			},
			Document: "ipAdEntAddr",
		},
		{
			OID:  "1.3.6.1.2.1.4.20.1.2",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(10001), nil
			},
			Document: "ipAdEntIfIndex",
		},
		{
			OID:  "1.3.6.1.2.1.4.20.1.3",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetMask())), nil
			},
			Document: "ipAdEntNetMask",
		},
		{
			OID:  "1.3.6.1.2.1.4.20.1.4",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(1), nil
			},
			Document: "ipAdEntBcastAddr",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.1",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				m := atopnet.CovertMaskToLen(d.GetMask())
				prfixip := d.GetIP() + "/" + strconv.Itoa(m)
				_, ip, err := net.ParseCIDR(prfixip)
				v := ip.IP.String()
				return v, err
			},
			Document: "ipRouteDest",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.2",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {

				return GoSNMPServer.Asn1IntegerWrap(10000), nil
			},
			Document: "ipRouteIfIndex",
		},

		{
			OID:  "1.3.6.1.2.1.4.21.1.3",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(0), nil
			},
			Document: "ipRouteMetric1",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.4",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(-1), nil
			},
			Document: "ipRouteMetric2",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.5",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(-1), nil
			},
			Document: "ipRouteMetric3",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.6",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(-1), nil
			},
			Document: "ipRouteMetric4",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.7",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetGateWay())), nil
			},
			Document: "ipRouteNextHop",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.8",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(3), nil
			},
			Document: "ipRouteType",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.9",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(1), nil
			},
			Document: "ipRouteProto",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.10",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(0), nil
			},
			Document: "ipRouteAge",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.11",
			Type: gosnmp.IPAddress,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IPAddressWrap(net.ParseIP(d.GetMask())), nil
			},
			Document: "ipRouteMask",
		},
		{
			OID:  "1.3.6.1.2.1.4.21.1.12",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(-1), nil
			},
			Document: "ipRouteMetric5",
		},

		{
			OID:  "1.3.6.1.2.1.4.21.1.13",
			Type: gosnmp.ObjectIdentifier,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1ObjectIdentifierWrap("0.0"), nil
			},
			Document: "ipRouteInfo",
		},

		{
			OID:  "1.3.6.1.2.1.4.22.1.1",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(10000), nil
			},
			Document: "ipNetToMediaIfIndex",
		},

		{
			OID:  "1.3.6.1.2.1.4.22.1.4",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(3), nil
			},
			Document: "ipNetToMediaType",
		},
	}
	return v
}
