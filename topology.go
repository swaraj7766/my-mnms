package mnms

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	g "github.com/gosnmp/gosnmp"
	"github.com/qeof/q"
)

type ChassisDetails struct {
	ChassisId  string `json:"chassis_id"`
	MacAddress string `json:"mac_address"`
}
type Link struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	SourcePort  string `json:"sourcePort"`
	TargetPort  string `json:"targetPort"`
	EdgeData    string `json:"edgeData"`
	LinkFlow    bool   `json:"linkFlow"`
	BlockedPort bool   `json:"blockedPort"`
}

func (l Link) Equal(other Link) bool {
	return l.Source == other.Source && l.Target == other.Target && l.SourcePort == other.SourcePort && l.TargetPort == other.TargetPort
}

type Node struct {
	Id         string `json:"id"`
	IpAddress  string `json:"ipAddress"`
	MacAddress string `json:"macAddress"`
	ModelName  string `json:"modelname"`
}

func (n Node) Equal(other Node) bool {
	return n.Id == other.Id && n.IpAddress == other.IpAddress && n.MacAddress == other.MacAddress && n.ModelName == other.ModelName
}

type Topology struct {
	LinkData []Link `json:"link_data"`
	NodeData []Node `json:"node_data"`
}

func (t Topology) Equal(other Topology) bool {
	if len(t.LinkData) != len(other.LinkData) || len(t.NodeData) != len(other.NodeData) {
		return false
	}
	for _, link := range t.LinkData {
		found := false
		for _, otherLink := range other.LinkData {
			if link.Equal(otherLink) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for _, node := range t.NodeData {
		found := false
		for _, otherNode := range other.NodeData {
			if node.Equal(otherNode) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

type TopoDevice struct {
	IpAddress  string `json:"ip_address"`
	MacAddress string `json:"mac_address"`
	ModelName  string `json:"modelname"`
}

var (
	ChassisIdList map[string]string
	templldpData  []Link
)

func TopologyPollingWithTimer(pollingInterval int) {
	pollingTimer := time.NewTicker(time.Second * time.Duration(pollingInterval))
	for range pollingTimer.C {
		ChassisIdList = make(map[string]string)
		templldpData = []Link{}
		targetes := []TopoDevice{}
		QC.DevMutex.Lock()
		currDevs := QC.DevData
		QC.DevMutex.Unlock()
		for _, dev := range currDevs {
			currTime := time.Now().Unix()
			lastTime, err := strconv.ParseInt(dev.Timestamp, 10, 64)
			if err != nil {
				q.Q("error in parsing timestamp")
				continue
			}
			if currTime-lastTime <= int64(QC.GwdInterval) {
				targetes = append(targetes, TopoDevice{IpAddress: dev.IPAddress, MacAddress: dev.Mac, ModelName: dev.ModelName})
			}
		}
		UpdateChassisIdList(targetes)
		UpdateRealTopologyData(targetes)
		CreateAndPublishTopologyData(targetes)
	}
}

func UpdateChassisIdList(devdata []TopoDevice) {
	for _, device := range devdata {
		result, err := GetChassisId(device.IpAddress, device.MacAddress)
		if err != nil {
			q.Q("chesis id not found")
		}
		if len(result.ChassisId) > 0 {
			ChassisIdList[result.ChassisId] = result.MacAddress
		}

	}
}

func UpdateRealTopologyData(devdata []TopoDevice) {
	for _, device := range devdata {
		err := GetLLDPData(device.IpAddress, device.MacAddress, device.ModelName, ChassisIdList)
		if err != nil {
			q.Q("error message: ", err)
		}
	}
}

func CreateAndPublishTopologyData(devdata []TopoDevice) {
	nodeData := []Node{}
	linkData := []Link{}
	inResult := make(map[string]bool)
	for _, device := range devdata {
		nodeData = append(nodeData, Node{Id: device.MacAddress, IpAddress: device.IpAddress, MacAddress: device.MacAddress, ModelName: device.ModelName})
	}
	// first add if blocaked port
	for _, lldpData := range templldpData {
		if lldpData.BlockedPort && !inResult[lldpData.EdgeData] {
			inResult[lldpData.EdgeData] = true
			linkData = append(linkData, lldpData)
		}
	}
	// add all port filtered data
	for _, lldpData := range templldpData {
		if _, ok := inResult[lldpData.EdgeData]; !ok {
			inResult[lldpData.EdgeData] = true
			linkData = append(linkData, lldpData)
		}
	}
	topologyData := Topology{}
	topologyData.NodeData = nodeData
	topologyData.LinkData = linkData
	_ = PublishTopology(topologyData)
}

func GetChassisId(targetIp string, macaddress string) (chl ChassisDetails, err error) {
	localChesId := ChassisDetails{}
	oids := []string{"1.0.8802.1.1.2.1.3.2.0"}
	result, err := GetOids(targetIp, oids, nil)
	if err != nil {
		return localChesId, err
	}
	for _, variable := range result {
		switch variable.Type {
		case g.OctetString:
			localChesId = ChassisDetails{ChassisId: ToMacString(variable.Value), MacAddress: macaddress}
		default:
			return localChesId, fmt.Errorf("value type is not correct")
		}
	}
	return localChesId, nil
}

func GetLLDPData(ipaddress string, macaddress string, modelname string, chesisIdList map[string]string) error {
	var blockedPort []string
	sysObjectId, err := GetSystemObjectId(ipaddress)
	if err != nil {
		return err
	}
	lldpRemChassisIdOid := "1.0.8802.1.1.2.1.4.1.1.5"
	lldpRemPortIdOid := "1.0.8802.1.1.2.1.4.1.1.7"
	erpsEnableOid := sysObjectId + ".4.4.1"
	erpsRsapVlanOid := sysObjectId + ".4.4.3.1.1"
	erpsDataOid := sysObjectId + ".4.4.3.1"
	erpsWPortStatusOid := sysObjectId + ".4.4.3.1.4."
	erpsEPortStatusOid := sysObjectId + ".4.4.3.1.5."
	erpsWPortOid := sysObjectId + ".4.4.3.1.2."
	erpsEPortOid := sysObjectId + ".4.4.3.1.3."
	lldpRemChassisId, err := GetBulk(ipaddress, lldpRemChassisIdOid, nil)
	if err != nil {
		return err
	}
	lldpRemPortId, err := GetBulk(ipaddress, lldpRemPortIdOid, nil)
	if err != nil {
		return err
	}
	erpsEnable, err := GetBulk(ipaddress, erpsEnableOid, nil)
	if err != nil {
		return err
	}
	erpsVlanId, err := GetBulk(ipaddress, erpsRsapVlanOid, nil)
	if err != nil {
		return err
	}
	erpsDataall, err := GetBulk(ipaddress, erpsDataOid, nil)
	if err != nil {
		return err
	}
	for _, element := range erpsEnable {
		if element.Value == 1 && len(erpsVlanId) > 0 && len(erpsDataall) > 0 {
			for _, vlanElement := range erpsVlanId {
				vlanId := fmt.Sprint(vlanElement.Value)
				var eastPort string
				var westPort string
				for _, erpsDataElement := range erpsDataall {
					if erpsDataElement.Name == erpsWPortOid+vlanId {
						westPort = fmt.Sprintf("port%v", string(erpsDataElement.Value.([]byte)))
					}
					if erpsDataElement.Name == erpsEPortOid+vlanId {
						eastPort = fmt.Sprintf("port%v", string(erpsDataElement.Value.([]byte)))
					}
					if erpsDataElement.Name == erpsWPortStatusOid+vlanId && erpsDataElement.Value == 2 {
						blockedPort = append(blockedPort, westPort)
					}
					if erpsDataElement.Name == erpsEPortStatusOid+vlanId && erpsDataElement.Value == 2 {
						blockedPort = append(blockedPort, eastPort)
					}
				}
			}
		}
	}

	QC.DevMutex.Lock()
	currDevs := QC.DevData
	QC.DevMutex.Unlock()
	for index, element := range lldpRemPortId {
		remotePortValue := string(element.Value.([]byte))
		if len(remotePortValue) != 8 && !strings.HasPrefix(remotePortValue, "port") {
			return fmt.Errorf("remote port is not correct")
		}
		sourcePortName := "port" + getLocalPort((strings.Replace(element.Name, lldpRemPortIdOid, "", -1)))
		intVar, _ := strconv.Atoi(remotePortValue[5:])
		remotePortName := "port" + strconv.Itoa(intVar)
		if len(lldpRemChassisId) != len(lldpRemPortId) {
			return fmt.Errorf("remote chassis length and remote port length diffrent")
		}
		remoteChesisId := ToMacString(lldpRemChassisId[index].Value)
		if len(remoteChesisId) == 0 {
			return fmt.Errorf("remote mac is empty")
		}
		found := false
		for _, dev := range currDevs {
			if dev.Mac == remoteChesisId {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("remote mac is not found")
		}
		remoteMacAddress := remoteChesisId
		sourceMacaddress := macaddress
		var edgeData string
		if sourceMacaddress < remoteMacAddress {
			edgeData = sourceMacaddress + "_" + remoteMacAddress
		} else {
			edgeData = remoteMacAddress + "_" + sourceMacaddress
		}
		isSourcePortBlocked := isSourcePortBlocked(sourcePortName, blockedPort)
		templldpData = append(templldpData, Link{
			Source:      sourceMacaddress,
			Target:      remoteMacAddress,
			SourcePort:  sourcePortName,
			TargetPort:  remotePortName,
			EdgeData:    edgeData,
			BlockedPort: isSourcePortBlocked,
			LinkFlow:    true,
		})
	}
	return nil
}

func GetSystemObjectId(targetIp string) (sysObjectId string, err error) {
	var systemObjectId string
	oids := []string{"1.3.6.1.2.1.1.2.0"}
	result, err := GetOids(targetIp, oids, nil)
	if err != nil {
		return systemObjectId, err
	}
	for _, variable := range result {
		switch variable.Type {
		case g.ObjectIdentifier:
			systemObjectId = fmt.Sprint(variable.Value)
		default:
			return systemObjectId, fmt.Errorf("value type is not correct")
		}
	}
	return systemObjectId, nil
}

func isSourcePortBlocked(sourcePort string, blockedPort []string) bool {
	isblockedPort := false
	for _, port := range blockedPort {
		if port == sourcePort {
			isblockedPort = true
		}
	}
	return isblockedPort
}

func ToMacString(value interface{}) string {
	bytes, ok := value.([]byte)
	if ok {
		temp := hex.EncodeToString(bytes)
		return insertdash(temp)
	}
	return ""
}

func insertdash(s string) string {
	upper := strings.ToUpper(s)
	var buffer bytes.Buffer
	n_1 := 2 - 1
	l_1 := len(upper) - 1
	for i, rune := range upper {
		buffer.WriteRune(rune)
		if i%2 == n_1 && i != l_1 {
			buffer.WriteRune('-')
		}
	}
	return buffer.String()
}

func getLocalPort(s string) string {
	splitValue := strings.Split(s, ".")[3]
	return splitValue
}

func PublishTopology(topologyData Topology) error {
	// send all devices info to root
	if QC.RootURL == "" {
		return fmt.Errorf("skip publishing devices, no root")
	}

	pTopologyData := make(map[string]Topology)
	QC.DevMutex.Lock()
	pTopologyData[QC.Name] = topologyData
	jsonBytes, err := json.Marshal(pTopologyData)
	QC.DevMutex.Unlock()

	if err != nil {
		q.Q(err)
		return err
	}
	resp, err := PostWithToken(QC.RootURL+"/api/v1/topology", QC.AdminToken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		q.Q(err, QC.RootURL)
	}
	if resp != nil {
		res := make(map[string]interface{})
		_ = json.NewDecoder(resp.Body).Decode(&res)
		// q.Q(res)
		// save close here
		defer resp.Body.Close()
	}
	return nil
}

func InsertTopology(topoKeys string, topoDesc Topology) bool {
	QC.DevMutex.Lock()
	_, ok := QC.TopologyData[topoKeys]
	if ok {
		if !topoDesc.Equal(QC.TopologyData[topoKeys]) {
			QC.TopologyData[topoKeys] = topoDesc
			// topologies vary all the time due to snmp polling, so we don't send syslog here
		}
	} else {
		QC.TopologyData[topoKeys] = topoDesc
		err := SendSyslog(LOG_ALERT, "InsertTopo", "new topology from: "+topoKeys)
		if err != nil {
			q.Q(err)
		}
	}
	QC.DevMutex.Unlock()

	return true
}

var DefaultSnmpOption = SnmpOptions{
	Port:      161,
	Community: "public",
	Version:   g.Version2c,
	Timeout:   time.Second,
}

func GetOids(target string, oids []string, opt *SnmpOptions) (results []g.SnmpPDU, err error) {
	var community string
	if opt == nil {
		opt = &DefaultSnmpOption
		devInfo, err := FindDevWithIP(target)
		community = opt.Community
		if err == nil {

			if len(devInfo.ReadCommunity) > 0 {
				community = devInfo.ReadCommunity
			}
		}

	}

	client := g.GoSNMP{
		Target:                  target,
		Community:               community,
		Port:                    opt.Port,
		Version:                 opt.Version,
		Timeout:                 opt.Timeout,
		Retries:                 1,
		UseUnconnectedUDPSocket: true,
	}

	err = client.Connect()
	if err != nil {
		return []g.SnmpPDU{}, err
	}
	defer client.Conn.Close()
	pkt, err := client.Get(oids)
	if err != nil {
		return []g.SnmpPDU{}, err
	}
	return pkt.Variables, nil
}

func GetBulk(target string, oid string, opt *SnmpOptions) (results []g.SnmpPDU, errs error) {
	snmpResults := []g.SnmpPDU{}
	var community string
	if opt == nil {
		opt = &DefaultSnmpOption
		devInfo, err := FindDevWithIP(target)
		community = opt.Community
		if err == nil {

			if len(devInfo.ReadCommunity) > 0 {
				community = devInfo.ReadCommunity
			}
		}

	}

	client := g.GoSNMP{
		Target:                  target,
		Community:               community,
		Port:                    opt.Port,
		Version:                 opt.Version,
		Timeout:                 opt.Timeout,
		Retries:                 1,
		UseUnconnectedUDPSocket: true,
	}

	err := client.Connect()
	if err != nil {
		return snmpResults, err
	}
	defer client.Conn.Close()
	err = client.BulkWalk(oid, func(pdu g.SnmpPDU) error {
		snmpResults = append(snmpResults, pdu)
		return nil
	})
	if err != nil {
		return snmpResults, err
	}
	return snmpResults, nil
}
