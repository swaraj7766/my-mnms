package mnms

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	snmplib "github.com/deejross/go-snmplib"
	"github.com/gosnmp/gosnmp"
	"github.com/qeof/q"
)

type SnmpOptions struct {
	Port      uint16
	Community string
	Version   gosnmp.SnmpVersion
	Timeout   time.Duration
}

// return read-all-only community and read-write-all community
func extractCommunityNames(output string) (string, string, error) {
	lines := strings.Split(output, "\n")
	communityNamePattern := `(\w+)\s+(read-all-only|read-write-all)`
	regex, err := regexp.Compile(communityNamePattern)
	if err != nil {
		return "", "", err
	}
	var readOnly string
	var readWrite string
	for _, line := range lines {
		match := regex.FindStringSubmatch(line)
		if len(match) == 3 {
			communityName := match[1]
			accessRight := match[2]

			if accessRight == "read-all-only" {
				readOnly = communityName
			}
			if accessRight == "read-write-all" {
				readWrite = communityName
			}
		}
	}
	return readOnly, readWrite, nil

}

// GetSNMPCommunity return read-only community and read-write community
func GetSNMPCommunity(user, pass, devIP string) (string, string, error) {
	dev, err := FindDevWithIP(devIP)
	if err != nil {
		return "", "", err
	}
	if dev.ModelName == "" {
		err := fmt.Errorf("error: invalid device model")
		return "", "", err
	}
	if !CheckSwitchCliModel(dev.ModelName) {
		err := fmt.Errorf("error: switch cli not available")
		return "", "", err
	}

	var cmdinfo CmdInfo
	err = SendSwitch(&cmdinfo, dev, user, pass, "show snmp community")
	if err != nil {
		q.Q(err)
		return "", "", err
	}

	return extractCommunityNames(cmdinfo.Result)
}

/*
var trapTypeMessage = []string{
	"cold start",
	"warm start",
	"link down",
	"link up",
	"authentication failure",
	"egp neighbor loss",
	"enterprise specific",
}*/

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
	var community string
	devInfo, err := FindDevWithIP(ipaddr)
	community = QC.SnmpOptions.Community
	if err == nil {

		if len(devInfo.ReadCommunity) > 0 {
			community = devInfo.ReadCommunity
		}
	}
	params := &gosnmp.GoSNMP{
		Target:                  ipaddr,
		Port:                    QC.SnmpOptions.Port,
		Community:               community,
		Version:                 QC.SnmpOptions.Version,
		Timeout:                 QC.SnmpOptions.Timeout,
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
	var community string
	devInfo, err := FindDevWithIP(address)
	community = QC.SnmpOptions.Community
	if err == nil {

		if len(devInfo.ReadCommunity) > 0 {
			community = devInfo.ReadCommunity
		}
	}
	params := &gosnmp.GoSNMP{
		Target:                  address,
		Port:                    QC.SnmpOptions.Port,
		Community:               community,
		Version:                 QC.SnmpOptions.Version,
		Timeout:                 QC.SnmpOptions.Timeout,
		UseUnconnectedUDPSocket: true,
	}

	q.Q("snmp get", params)
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

// SnmpWalk - walk snmp data
func SnmpWalk(address string, oid string) (result []gosnmp.SnmpPDU, err error) {
	var community string
	devInfo, err := FindDevWithIP(address)
	community = QC.SnmpOptions.Community
	if err == nil {

		if len(devInfo.ReadCommunity) > 0 {
			community = devInfo.ReadCommunity
		}
	}
	params := &gosnmp.GoSNMP{
		Target:                  address,
		Port:                    QC.SnmpOptions.Port,
		Community:               community,
		Version:                 QC.SnmpOptions.Version,
		Timeout:                 QC.SnmpOptions.Timeout,
		UseUnconnectedUDPSocket: true,
	}

	err = params.Connect()

	if err != nil {
		q.Q("error: snmp connect", err)
		// cmdinfo.Status = "error: cannot contact snmp target"
		return nil, err
	}
	defer params.Conn.Close()
	result, err = params.WalkAll(oid)
	if err != nil && len(result) == 0 {
		q.Q("error: snmp walk", err)
		return nil, err
	}
	return result, nil
}

// SnmpBulk - Bulk snmp data
func SnmpBulk(address string, oid string) (result []gosnmp.SnmpPDU, err error) {
	var community string
	devInfo, err := FindDevWithIP(address)
	community = QC.SnmpOptions.Community
	if err == nil {

		if len(devInfo.ReadCommunity) > 0 {
			community = devInfo.ReadCommunity
		}
	}
	params := &gosnmp.GoSNMP{
		Target:                  address,
		Port:                    QC.SnmpOptions.Port,
		Community:               community,
		Version:                 QC.SnmpOptions.Version,
		Timeout:                 QC.SnmpOptions.Timeout,
		UseUnconnectedUDPSocket: true,
	}

	err = params.Connect()

	if err != nil {
		q.Q("error: snmp connect", err)
		// cmdinfo.Status = "error: cannot contact snmp target"
		return nil, err
	}
	defer params.Conn.Close()
	result, err = params.BulkWalkAll(oid)
	if err != nil && len(result) == 0 {
		q.Q("error: snmp bulk", err)
		return nil, err
	}
	return result, nil
}

// SnmpSet - set snmp data
func SnmpSet(address, oid string, value string, valuetype string) (result *gosnmp.SnmpPacket, err error) {
	var community string
	devInfo, err := FindDevWithIP(address)
	community = QC.SnmpOptions.Community
	if err == nil {

		if len(devInfo.WriteCommunity) > 0 {
			community = devInfo.WriteCommunity
		}
	}
	params := &gosnmp.GoSNMP{
		Target:                  address,
		Port:                    QC.SnmpOptions.Port,
		Community:               community,
		Version:                 QC.SnmpOptions.Version,
		Timeout:                 QC.SnmpOptions.Timeout,
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
		val, ok := pdu.Value.([]byte)
		if ok {
			addr := net.HardwareAddr(val[:])
			r, err := net.ParseMAC(addr.String())
			if err != nil {
				resStr = string(getValidByte(pdu.Value.([]byte)))
			} else {
				resStr = strings.ToUpper(r.String())
			}
		} else {
			resStr = string(pdu.Value.([]byte))
		}
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

// Use snmp get/set/communities/update.
//
// Usage : snmp get [ip address] [oid]
//
//	[ip address]  : target device ip address
//	[oid]         : target oid
//
// Example : snmp get 10.0.50.1 1.3.6.1.2.1.1.1.0
//
// Usage : snmp set [ip address] [oid] [value] [value type]
//
//	[ip address]  : target device ip address
//	[oid]         : target oid
//	[value]       : would be set value
//	[value type]  : would be set value type.(OctetString, BitString, SnmpNullVar, Counter,
//	                Counter64, Gauge, Opaque, Integer, ObjectIdentifier, IpAddress, TimeTicks)
//
// Example : snmp set 10.0.50.1 1.3.6.1.2.1.1.4.0 www.atop.com.tw OctetString
//
// Usage: snmp communities [user] [password] [mac]
// Read device's SNMP communities and update to system.
//
//	[user]     : Device telnt login user
//	[password] : Device telnt login password
//	[mac]      : Device mac address
//
// Example: snmp communities admin default 00-60-E9-27-E3-39
//
// Usage: snmp update community [mac] [read community] [write community]
// Update device's SNMP communities manually.
//
//	[mac]            : Device mac address
//	[read community] : Device snmp read community
//	[write community]: Device snmp write community
//
// Example: snmp update community 00-60-E9-27-E3-39 public private
//
// Usage: snmp options [port] [community] [version] [timeout]
// Update global snmp options.
//
//	[port]     : snmp listen port
//	[community]: snmp community
//	[version]  : snmp version
//	[timeout]  : snmp timeout
//
// Example: snmp options 161 public 2c 2
func SnmpCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 4 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	// snmp options {port} {community} {version} {timeout}
	// snmp options 161 private 2c 2
	if ws[1] == "options" {
		q.Q(ws)
		if len(ws) < 6 {
			cmdinfo.Status = "error: invalid snmp options command"
			return cmdinfo
		}
		port, err := strconv.Atoi(ws[2])
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: port %v", err)
			return cmdinfo
		}
		timeout, err := strconv.Atoi(ws[5])
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: timeout %v", err)
			return cmdinfo
		}
		var version gosnmp.SnmpVersion
		switch ws[4] {
		case "1":
			version = gosnmp.Version1
		case "2c":
			version = gosnmp.Version2c
		case "3":
			version = gosnmp.Version3
		default:
			cmdinfo.Status = fmt.Sprintf("error: version %v, accept 1|2c|3", err)
			return cmdinfo
		}
		q.Q(port, ws[3], version, timeout)

		QC.SnmpOptions = SnmpOptions{
			Port:      uint16(port),
			Community: ws[3],
			Version:   version,
			Timeout:   time.Duration(timeout) * time.Second,
		}

		cmdinfo.Status = "ok"
		q.Q(QC.SnmpOptions)
		return cmdinfo
	}

	if ws[1] == "communities" {
		// read communities and write to DevInfo
		// snmp communities {user} {password} {mac}

		if len(ws) < 5 {
			cmdinfo.Status = "error: invalid snmp communities command"
			return cmdinfo
		}
		// find device
		devID := ws[4]
		user := ws[2]
		password := ws[3]
		q.Q(devID, user, password)
		dev, err := FindDev(devID)
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		// get communities
		r, rw, err := GetSNMPCommunity(user, password, dev.IPAddress)
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		err = InsertCommunities(dev.Mac, r, rw)
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"

		return cmdinfo
	}
	// snmp update community {mac} {read community} {write community}

	if ws[1] == "update" && ws[2] == "community" {
		if len(ws) < 6 {
			cmdinfo.Status = "error: invalid snmp update community command"
			return cmdinfo
		}
		mac := ws[3]
		r := ws[4]
		w := ws[5]
		dev, err := FindDev(mac)
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		err = InsertCommunities(dev.Mac, r, w)
		if err != nil {
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		return cmdinfo
	}

	address := ws[2]
	oid := ws[3]
	oids := []string{oid}
	// validate address
	err := CheckIPAddress(address)
	if err != nil {
		cmdinfo.Status = fmt.Sprintf("error:%v", err)
		return cmdinfo
	}

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

	if ws[1] == "walk" {
		cmdinfo.Status = running.String()
		go func(cmdinfo CmdInfo) {
			defer func() {
				QC.CmdMutex.Lock()
				QC.CmdData[cmdinfo.Command] = cmdinfo
				QC.CmdMutex.Unlock()
			}()
			type walkResult struct {
				Oid   string
				Value string
			}
			res, err := SnmpWalk(address, oid)
			if err != nil {
				q.Q(err)
				cmdinfo.Status = fmt.Sprintf("error: %v", err)
				return
			}
			q.Q("snmp walk ", res)
			walks := []walkResult{}
			for _, variable := range res {
				oid := variable.Name[1:]
				resStr := PDUToString(variable)
				walks = append(walks, walkResult{Oid: oid, Value: resStr})
			}
			b, err := json.Marshal(&walks)
			if err != nil {
				q.Q(err)
				cmdinfo.Status = fmt.Sprintf("error: %v", err)
				return
			}
			cmdinfo.Status = "ok"
			cmdinfo.Result = string(b)
			// return cmdinfo
		}(*cmdinfo)
		return cmdinfo
	}
	if ws[1] == "bulk" {
		cmdinfo.Status = running.String()
		go func(cmdinfo CmdInfo) {
			defer func() {
				QC.CmdMutex.Lock()
				QC.CmdData[cmdinfo.Command] = cmdinfo
				QC.CmdMutex.Unlock()
			}()

			type walkResult struct {
				Oid   string
				Value string
			}
			res, err := SnmpBulk(address, oid)
			if err != nil {
				q.Q(err)
				cmdinfo.Status = fmt.Sprintf("error: %v", err)
				return
			}
			q.Q("snmp bulk ", res)
			walks := []walkResult{}
			for _, variable := range res {
				oid := variable.Name[1:]
				resStr := PDUToString(variable)
				walks = append(walks, walkResult{Oid: oid, Value: resStr})
			}
			b, err := json.Marshal(&walks)
			if err != nil {
				q.Q(err)
				cmdinfo.Status = fmt.Sprintf("error: %v", err)
				return
			}
			cmdinfo.Status = "ok"
			cmdinfo.Result = string(b)
		}(*cmdinfo)
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
		cmdinfo.Status = fmt.Sprintf("error: %v", pkt.Error.String())
		return cmdinfo
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
	prettyPrint, _ := json.Marshal(trap)
	if QC.IsRoot {
		// TODO save data to to q
		q.Q("trapserver :", string(prettyPrint))
	} else {
		err := SendSyslog(LOG_ALERT, "trapserver", string(prettyPrint))
		if err != nil {
			q.Q("error: sending trap syslog", err)
		}
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
