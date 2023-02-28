package mnms

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	snmplib "github.com/deejross/go-snmplib"
	"github.com/gosnmp/gosnmp"
	"github.com/qeof/q"
)

// snmp scan, get, set

func SnmpScan() error {
	cidrs, err := IfnetCidrs()
	if err != nil {
		q.Q("error: cannot figure out CIDRs", err)
		return fmt.Errorf("error: cidr, %v", err)
	}

	for _, cidr := range cidrs {
		err = ScanCIDR(cidr)
		if err != nil {
			q.Q("error: can't scan cidr", cidr, err)
			continue
		}
	}
	return nil
}

func ScanCIDR(cidr string) error {
	ipaddrs, err := GetIpAddrs(cidr)
	if err != nil {
		q.Q("error: cannot get ip addresses from CIDR", cidr)
		return err
	}

	ipChan := make(chan string, 10)
	done := make(chan bool)

	go func() {
		for _, ip := range ipaddrs {
			ipChan <- ip.String()
		}

		done <- true
	}()

	var wg sync.WaitGroup

loop:
	for {
		select {
		case <-done:
			q.Q("done")
			break loop
		case ipaddr := <-ipChan:
			wg.Add(1)
			go func() {
				defer wg.Done()
				err = probeSnmp(ipaddr)
				if err != nil {
					q.Q("error: probe snmp", err)
				}
			}()
		}
	}

	wg.Wait()
	q.Q(numModelsFound)

	return nil
}

func probeSnmp(ipaddr string) error {
	var err error
	params := &gosnmp.GoSNMP{
		Target:                  ipaddr,
		Port:                    161,
		Community:               "public",
		Version:                 gosnmp.Version2c,
		Timeout:                 time.Duration(2) * time.Second,
		UseUnconnectedUDPSocket: true,
	}

	err = params.Connect()
	if err != nil {
		q.Q("error: snmp connect", err)
		return err
	}
	defer params.Conn.Close()

	oids := []string{
		"1.3.6.1.2.1.1.1.0",
		"1.3.6.1.2.1.1.2.0",
		"1.3.6.1.2.1.1.3.0",
		"1.3.6.1.2.1.1.4.0",
		"1.3.6.1.2.1.1.5.0",
		"1.3.6.1.2.1.1.6.0",
		"1.3.6.1.2.1.1.7.0",
		"1.3.6.1.2.1.1.8.0",
	}

	// XXX these may not work on SE400, SE500, CWR5805, ...
	atopOidSuffixes := []string{
		".2.3.1.1.3.1",
		".2.3.1.1.4.1",
		".2.3.1.1.5.1",
		".2.3.1.1.6.1",
		".2.3.1.1.7.1",
		".1.1.0",
		".1.2.0",
		".1.3.0",
		".1.4.0",
		".1.5.0",
		".1.6.0",
		".1.7.0",
		".1.8.0",
		".1.9.0",
		".1.10.0",
		".1.11.0",
		".1.12.0",
		".1.13.0",
		".1.14.0",
	}

	result, err := params.Get(oids)
	if err != nil {
		q.Q("error: snmp get", ipaddr, err)
		return err
	}

	var atopOids []string
	atopOidPrefix := ""

	atopOidPrefix = parseSnmpResults(result, atopOidPrefix)

	if atopOidPrefix != "" {
		for _, Suffix := range atopOidSuffixes {
			atopOids = append(atopOids, atopOidPrefix+Suffix)
		}

		result, err := params.Get(atopOids)
		if err != nil {
			q.Q("error: snmp get", ipaddr, err)
			return err
		}
		parseSnmpResults(result, atopOidPrefix)
	}
	return nil
}

var numModelsFound int

func parseSnmpResults(result *gosnmp.SnmpPacket, atopOidPrefix string) string {
	model := GwdModelInfo{}
	prefix := ""

	for i, variable := range result.Variables {
		oid := variable.Name[1:]
		q.Q(i, oid)

		switch variable.Type {
		case gosnmp.OctetString:
			if strings.HasPrefix(oid, "1.3.6.1.4.1.3755.") {
				if strings.HasSuffix(oid, ".1.6.0") ||
					strings.HasSuffix(oid, "3755.0.1.1.1.9.0") {
					vv := variable.Value.([]byte)
					model.MACAddress = fmt.Sprintf("%.2X-%.2X-%.2X-%.2X-%.2X-%.2X",
						vv[0], vv[1], vv[2], vv[3], vv[4], vv[5])
					q.Q(model.MACAddress)
				} else if strings.HasSuffix(oid, ".1.5.0") {
					k := string(variable.Value.([]byte))
					model.Kernel = CleanStr(k)
				} else if strings.HasSuffix(oid, ".1.4.0") {
					apStr := string(variable.Value.([]byte))
					model.Ap = CleanStr(apStr)
				} else if strings.HasSuffix(oid, ".1.10.0") {
					m := string(variable.Value.([]byte))
					model.Model = CleanStr(m)
				}
			} else {
				q.Q(variable.Value.([]byte))
			}
		case gosnmp.ObjectIdentifier:
			prefix = variable.Value.(string)
			q.Q(prefix)
			prefix = prefix[1:]
		case gosnmp.IPAddress:
			ipaddr := variable.Value.(string)
			q.Q(ipaddr)
			if strings.HasPrefix(oid, "1.3.6.1.4.1.3755.") {
				if strings.HasSuffix(oid, ".2.3.1.1.3.1") {
					model.IPAddress = ipaddr
				} else if strings.HasSuffix(oid, ".2.3.1.1.4.1") {
					model.Netmask = ipaddr
				} else if strings.HasSuffix(oid, ".2.3.1.1.5.1") {
					model.Gateway = ipaddr
				}
			}
		default:
			q.Q(gosnmp.ToBigInt(variable.Value))
		}
	}

	if atopOidPrefix != "" && model.Model != "" {
		if model.Hostname == "" {
			model.Hostname = "unknown"
		}
		model.ScannedBy = QC.Name
		InsertModel(model, "snmp")
		numModelsFound++
		q.Q(numModelsFound, model)
	}

	return prefix
}

const SystemObjectID = ".1.3.6.1.2.1.1.2.0"

// SnmpGetObjectID get snmp object id from device
func SnmpGetObjectID(path string) (string, error) {
	rets, err := SnmpGet(path, []string{SystemObjectID})
	if err != nil {
		return "", err
	}
	if len(rets.Variables) <= 0 {
		return "", fmt.Errorf("not found")
	}
	return PDUToString(rets.Variables[0]), nil
}

// SnmpGet - get snmp data
func SnmpGet(address string, oids []string) (result *gosnmp.SnmpPacket, err error) {
	params := &gosnmp.GoSNMP{
		Target:                  address,
		Port:                    161,
		Community:               "private",
		Version:                 gosnmp.Version2c,
		Timeout:                 time.Duration(2) * time.Second,
		UseUnconnectedUDPSocket: true,
	}

	err = params.Connect()

	if err != nil {
		q.Q("error: snmp connect", err)
		// cmdinfo.Status = "error: cannot contact snmp target"
		return nil, err
	}
	defer params.Conn.Close()
	result, err = params.Get(oids)
	if err != nil {
		q.Q("error: snmp get", err)
		return nil, err
	}
	return result, nil
}

// SnmpSet - set snmp data
func SnmpSet(address, oid string, value string, valuetype string) (result *gosnmp.SnmpPacket, err error) {
	params := &gosnmp.GoSNMP{
		Target:                  address,
		Port:                    161,
		Community:               "private",
		Version:                 gosnmp.Version2c,
		Timeout:                 time.Duration(2) * time.Second,
		UseUnconnectedUDPSocket: true,
	}
	t := GetType(valuetype)
	q.Q("snmp set type", t, valuetype)
	anyVal := ConvertSetValue(value, valuetype)
	if anyVal == nil {
		return nil, fmt.Errorf("error: value type not supported (%s)", valuetype)
	}
	data := []gosnmp.SnmpPDU{{Name: oid, Value: anyVal, Type: gosnmp.Asn1BER(t)}}
	q.Q(data)
	err = params.Connect()

	if err != nil {
		q.Q("error: snmp connect", err)
		// cmdinfo.Status = "error: cannot contact snmp target"
		return nil, err
	}
	defer params.Conn.Close()
	return params.Set(data)
}

// PDUToString - convert PDU to string
func PDUToString(pdu gosnmp.SnmpPDU) string {
	var resStr string
	switch pdu.Type {
	case gosnmp.OctetString:
		resStr = string(pdu.Value.([]byte))
	case gosnmp.ObjectIdentifier:
		resStr = pdu.Value.(string)
	case gosnmp.IPAddress:
		resStr = pdu.Value.(string)

	default:
		resStr = fmt.Sprintf("%d", (gosnmp.ToBigInt(pdu.Value)))
	}
	resStr = CleanStr(resStr)
	return resStr
}

func SnmpCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 4 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	address := ws[2]
	oid := ws[3]
	oids := []string{oid}
	if ws[1] == "get" {
		// snmp get {ip_adddress} {oid}
		res, err := SnmpGet(address, oids)
		// res, err := params.Get(oids)
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		var resStr string
		q.Q("snmp get ", res.Variables)
		for _, variable := range res.Variables {
			oid := variable.Name[1:]
			resStr = PDUToString(variable)
			q.Q(oid, resStr)
		}
		cmdinfo.Status = "ok"
		cmdinfo.Result = resStr
		return cmdinfo
	}
	// snmp set {ip_adddress} {oid} {value} {type}
	if ws[1] != "set" {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}

	if len(ws) < 6 {
		cmdinfo.Status = "error: invalid command, too few args"
		return cmdinfo
	}
	value := ws[4]
	valuetype := ws[5]

	pkt, err := SnmpSet(address, oid, value, valuetype)
	if err != nil {
		q.Q(err)
		cmdinfo.Status = fmt.Sprintf("error: %v", err)
		return cmdinfo
	}

	if uint8(pkt.Error) > 0 {
		q.Q(pkt.Error)
	}
	q.Q(pkt.Variables)
	cmdinfo.Status = "ok"
	return cmdinfo
}

func ConvertSetValue(val, t string) any {
	t = strings.TrimSpace(t)
	switch t {
	case "OctetString":
		return val
	case "Integer":
		// val convert to int
		v, err := strconv.Atoi(val)
		if err != nil {
			q.Q("error: convert to int", err)
			return nil
		}
		return v
	default:
		return val
	}
}

func GetType(t string) byte {
	t = strings.TrimSpace(t)
	switch t {
	case "OctetString":
		return byte(gosnmp.OctetString)
	case "BitString":
		return byte(gosnmp.BitString)
	case "SnmpNullVar":
		return byte(gosnmp.Null)
	case "Counter":
		return byte(gosnmp.Counter32)
	case "Counter64":
		return byte(gosnmp.Counter64)
	case "Gauge":
		return byte(gosnmp.Gauge32)
	case "Opaque":
		return byte(gosnmp.Opaque)
	case "Integer":
		return byte(gosnmp.Integer)
	case "ObjectIdentifier":
		return byte(gosnmp.ObjectIdentifier)
	case "IpAddress":
		return byte(gosnmp.IPAddress)
	case "TimeTicks":
		return byte(gosnmp.TimeTicks)
	}

	return byte(0)
}

type snmpHandler struct{}

func (h snmpHandler) OnError(addr net.Addr, err error) {
	q.Q(addr.String(), err)
}

func (h snmpHandler) OnTrap(addr net.Addr, trap snmplib.Trap) {
	prettyPrint, _ := json.MarshalIndent(trap, "", "\t")
	if QC.IsRoot {
		// TODO save data to to q
		q.Q("trapserver :", string(prettyPrint))
	} else {
		// SendSyslogMessage("trap: here trap content will go")
	}
}

func StartTrapServer() {
	ws := strings.Split(QC.TrapServerAddr, ":")
	if len(ws) < 2 {
		q.Q("Invalid trap server address", QC.TrapServerAddr)
		return
	}
	addr := ws[0]
	port, err := strconv.Atoi(ws[1])
	if err != nil {
		q.Q(err)
		return
	}
	q.Q(addr, port)
	server, err := snmplib.NewTrapServer(addr, port)
	if err != nil {
		q.Q(err)
		return
	}
	server.ListenAndServe(snmpHandler{})
}
