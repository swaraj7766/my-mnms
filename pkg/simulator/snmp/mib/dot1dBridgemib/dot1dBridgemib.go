package dot1dbridgemib

import (
	"encoding/hex"
	"fmt"

	snmpvalue "mnms/pkg/simulator/snmp/bindvalue"

	"github.com/slayercat/GoSNMPServer"
	"github.com/slayercat/gosnmp"
)

// All oid of dot1dbridge
func All(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	var result []*GoSNMPServer.PDUValueControlItem
	result = append(result, dot1dBase(d)...)
	result = append(result, dot1dStp(d)...)
	result = append(result, dot1dTp(d)...)
	result = append(result, pBridgeMIB(d)...)
	result = append(result, qBridgeMIB(d)...)
	return result
}

func dot1dBase(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:  "1.3.6.1.2.1.17.1.1.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				targetStr := d.GetMac()
				decoded, err := hex.DecodeString(targetStr)
				if err != nil {
					return nil, err
				}
				return GoSNMPServer.Asn1OctetStringWrap(string(decoded)), nil
			},
		},

		{
			OID:   "1.3.6.1.2.1.17.1.2.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(29), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.1.3.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
	}

	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.1.4.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.1.4.1.2.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.1.4.1.3.%d", ifIndex),
				Type: gosnmp.ObjectIdentifier,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1ObjectIdentifierWrap("1.3.6.1.4.1.3755"), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.1.4.1.4.%d", ifIndex),
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.1.4.1.5.%d", ifIndex),
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap(0), nil
				},
			},
		}
		toRet = append(toRet, currentIf...)
	}

	return toRet
}

func dot1dStp(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {

	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.17.2.1.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(3), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.2.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(32768), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.3.0",
			Type:  gosnmp.TimeTicks,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1TimeTicksWrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.4.0",
			Type:  gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
		},
		{
			OID:  "1.3.6.1.2.1.17.2.5.0",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				targetStr := d.GetMac()
				decoded, err := hex.DecodeString(targetStr)
				if err != nil {
					return nil, err
				}
				return GoSNMPServer.Asn1OctetStringWrap(string(decoded)), nil
			},
		},
		{
			OID:   "1.3.6.1.2.1.17.2.6.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.7.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(255), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.8.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2000), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.9.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(200), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.10.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(100), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.11.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1500), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.12.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2000), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.13.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(200), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.2.14.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1500), nil },
		},
	}

	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.2.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(128), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.3.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.4.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.5.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(65535), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.6.%d", ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					targetStr := d.GetMac()
					decoded, err := hex.DecodeString(targetStr)
					if err != nil {
						return nil, err
					}
					return GoSNMPServer.Asn1OctetStringWrap(string(decoded)), nil
				},
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.7.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.8.%d", ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1OctetStringWrap("00-00-00-00-00-00-00-00"), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.9.%d", ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1OctetStringWrap(fmt.Sprintf("0x00 %02d", ifIndex-1)), nil
				},
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.2.15.1.10.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
		}
		toRet = append(toRet, currentIf...)
	}
	return toRet
}

func dot1dTp(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.17.4.1.0",
			Type:  gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
		},

		{
			OID:   "1.3.6.1.2.1.17.4.2.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(300), nil },
		},
		{
			OID:  "1.3.6.1.2.1.17.4.3.1.1",
			Type: gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) {
				return GoSNMPServer.Asn1OctetStringWrap("00-00-00-00-00-00"), nil
			},
		},
		{
			OID:   "1.3.6.1.2.1.17.4.3.1.2",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.4.3.1.3",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(3), nil },
		},
	}

	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.4.4.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(ifIndex - 1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.4.4.1.2.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1514), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.4.4.1.3.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.4.4.1.4.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.4.4.1.5.%d", ifIndex),
				Type: gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter32Wrap(0), nil
				},
			},

			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.4.5.1.1.%d", ifIndex),
				Type: gosnmp.Counter64,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter64Wrap(1), nil
				},
			},

			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.4.5.1.2.%d", ifIndex),
				Type: gosnmp.Counter64,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter64Wrap(1), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.4.5.1.3.%d", ifIndex),
				Type: gosnmp.Counter64,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1Counter64Wrap(1), nil
				},
			},
		}
		toRet = append(toRet, currentIf...)
	}

	return toRet
}

func pBridgeMIB(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.17.6.1.1.1.0",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("r"), nil },
		},

		{
			OID:   "1.3.6.1.2.1.17.6.1.1.2.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.6.1.1.3.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
	}

	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.6.1.1.4.1.1.%d", ifIndex),
				Type:  gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("à"), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.6.1.2.1.1.1.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.6.1.2.1.1.2.%d", ifIndex),
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap(2), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.6.1.3.1.1.1.%d", ifIndex),
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap(20), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.6.1.3.1.1.2.%d", ifIndex),
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap(60), nil
				},
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.6.1.3.1.1.3.%d", ifIndex),
				Type: gosnmp.Integer,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1IntegerWrap(1000), nil
				},
			},
		}
		toRet = append(toRet, currentIf...)
	}

	return toRet
}

func qBridgeMIB(d *snmpvalue.Value) []*GoSNMPServer.PDUValueControlItem {
	toRet := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.17.7.1.1.1.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},

		{
			OID:   "1.3.6.1.2.1.17.7.1.1.2.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(4094), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.1.3.0",
			Type:  gosnmp.Gauge32,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(30), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.1.4.0",
			Type:  gosnmp.Gauge32,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.1.5.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
		{
			OID:   ".1.3.6.1.2.1.17.7.1.2.2.1.2.0.0.0.0.0.0.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
		{
			OID:   ".1.3.6.1.2.1.17.7.1.2.2.1.3.0.0.0.0.0.0.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(3), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.1.0",
			Type:  gosnmp.Counter32,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.2.1.3.0.1",
			Type:  gosnmp.Gauge32,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(1), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.2.1.4.0.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("0x00 0F FF FF"), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.2.1.5.0.1",
			Type:  gosnmp.OctetString,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1OctetStringWrap("0xFF FF FF FF"), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.2.1.6.0.1",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.2.1.7.0.1",
			Type:  gosnmp.TimeTicks,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1TimeTicksWrap(uint32(1717690000)), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.4.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},
	}

	for ifIndex := 1; ifIndex <= d.Port(); ifIndex++ {
		currentIf := []*GoSNMPServer.PDUValueControlItem{
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.7.1.4.5.1.1.%d", ifIndex),
				Type:  gosnmp.Gauge32,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Gauge32Wrap(1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.7.1.4.5.1.2.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.7.1.4.5.1.3.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.7.1.4.5.1.4.%d", ifIndex),
				Type:  gosnmp.Integer,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(2), nil },
			},
			{
				OID:   fmt.Sprintf("1.3.6.1.2.1.17.7.1.4.5.1.5.%d", ifIndex),
				Type:  gosnmp.Counter32,
				OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1Counter32Wrap(0), nil },
			},
			{
				OID:  fmt.Sprintf("1.3.6.1.2.1.17.7.1.4.5.1.6.%d", ifIndex),
				Type: gosnmp.OctetString,
				OnGet: func() (value interface{}, err error) {
					return GoSNMPServer.Asn1OctetStringWrap("00-00-00-00-00-00"), nil
				},
			},
		}

		toRet = append(toRet, currentIf...)
	}
	toRet1 := []*GoSNMPServer.PDUValueControlItem{
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.9.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(0), nil },
		},
		{
			OID:   "1.3.6.1.2.1.17.7.1.4.10.0",
			Type:  gosnmp.Integer,
			OnGet: func() (value interface{}, err error) { return GoSNMPServer.Asn1IntegerWrap(1), nil },
		},
	}
	toRet = append(toRet, toRet1...)

	return toRet
}
