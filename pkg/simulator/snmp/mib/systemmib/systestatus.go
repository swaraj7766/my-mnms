package systemmib

import (
	"errors"
	"fmt"
	"strings"

	"mnms/pkg/simulator/devicetype"
	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"
	"mnms/pkg/simulator/snmp/mib/atopmib"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of system
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {

	v := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  ".1.3.6.1.2.1.1.1.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringWrap(d.GetModel()), nil
			},
			Document: "sysDescr",
		},
		{
			OID:  ".1.3.6.1.2.1.1.2.0",
			Type: gosnmp.ObjectIdentifier,
			OnGet: func() (value interface{}, err error) {
				m, _ := devicetype.ParsingType(d.GetModel())
				oid := atopmib.OidType(m)
				return GoSNMPServer.Asn1ObjectIdentifierWrap(fmt.Sprintf("1.3.6.1.4.1.3755.0.0.%v", oid)), nil
			},
			Document: "sysObjectID",
		},
		{
			OID:  ".1.3.6.1.2.1.1.4.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringWrap("www.atop.com.tw"), nil
			},
			Document: "sysContact",
		},
		{
			OID:  ".1.3.6.1.2.1.1.5.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringWrap(d.GetSystem()), nil
			},
			OnSet: func(value interface{}) (err error) {
				val, ok := value.([]uint8)
				if ok {
					var sb strings.Builder
					for _, v := range val {
						sb.WriteString(string(v))
					}
					d.SetSystem(sb.String())
					return nil
				}
				return errors.New("fromat error")
			},
			Document: "sysName",
		},
		{
			OID:  ".1.3.6.1.2.1.1.6.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringWrap("Switch Location"), nil
			},
			Document: "sysContact",
		},
		{
			OID:  ".1.3.6.1.2.1.1.7.0",
			Type: gosnmp.Integer,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1IntegerWrap(2), nil
			},
			Document: "sysServices",
		},
	}
	return v
}
