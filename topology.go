package mnms

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	g "github.com/gosnmp/gosnmp"
	"github.com/qeof/q"
	"github.com/sirupsen/logrus"
)

const (
	START = 1
	STOP  = 2
	END   = 3
)

var ch = make(chan int, 1)

type Link struct {
	Source      string `json:"source"`
	Target      string `json:"target"`
	SourcePort  string `json:"sourcePort"`
	TargetPort  string `json:"targetPort"`
	EdgeData    string `json:"edgeData"`
	LinkFlow    bool   `json:"linkFlow"`
	BlockedPort bool   `json:"blockedPort"`
}

type Node struct {
	Id         string `json:"id"`
	IpAddress  string `json:"ipAddress"`
	MacAddress string `json:"macAddress"`
}

type Topology struct {
	LinkData []Link `json:"link_data"`
	NodeData []Node `json:"node_data"`
}

type (
	Links = []Link
	Nodes = []Node
)

func TopologyPollingWithTimer(pollingInterval int) {
	pollingTimer := time.NewTicker(time.Second * time.Duration(pollingInterval))
	start := true
	for range pollingTimer.C {
		select {
		case op := <-ch:
			if op == START {
				start = true
			} else if op == STOP {
				start = false
			} else if op == END {
				return
			}
		default:
			// fetch data here
			if start {
				q.Q("Get Topology data and send through web socket")
				targetes := []string{}
				for _, dev := range QC.DevData {
					targetes = append(targetes, dev.IPAddress)
				}
				LldpScan(targetes)
			}
		}
	}
}

func RemoveIndex(s []Link, index int) []Link {
	return append(s[:index], s[index+1:]...)
}

func LldpScan(targetes []string) {
	// fmt.Println(targetes)
	nodesData := Nodes{}
	linksData := Links{}
	var wg sync.WaitGroup
	for _, ip := range targetes {
		wg.Add(1)
		go func(ip string) {
			localChesIdResult := gatLocalChesisId(ip)
			nodesData = append(nodesData, Node{Id: localChesIdResult, IpAddress: ip, MacAddress: localChesIdResult})

			tempLinksData := getLldpData(ip, localChesIdResult)
			for _, linkData := range tempLinksData {
				isRealTopology := realTopologyData(linksData, linkData)
				if isRealTopology {
					linksData = append(linksData, linkData)
				}
				if linkData.BlockedPort && !isRealTopology {
					for index, linkDataBport := range linksData {
						if linkData.EdgeData == linkDataBport.EdgeData {
							linksData = RemoveIndex(linksData, index)
						}
					}
					linksData = append(linksData, linkData)
				}

				//
			}
			wg.Done()
		}(ip)
	}
	wg.Wait()
	topologyData := Topology{}
	topologyData.NodeData = nodesData
	topologyData.LinkData = linksData
	// res, err := PrettyStruct(topologyData)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("topology Data:  ", res)
	_ = PublishTopology(topologyData)
	// fmt.Printf("nodes Data %v \n", nodesData)
	// res, err = PrettyStruct(linksData)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("links Data:  ", res)
	return
}

func gatLocalChesisId(targetIp string) (localChesisId string) {
	var localChesId string
	oids := []string{"1.0.8802.1.1.2.1.3.2.0"}
	result, err2 := GetOids(targetIp, oids, nil)
	if err2 != nil {
		q.Q("Get() err: ", err2)
	}
	for _, variable := range result {
		switch variable.Type {
		case g.OctetString:
			localChesId = ToMacString(variable.Value)
		default:
			q.Q("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
	return localChesId
}

func getLldpData(targetIp string, sourceChesId string) Links {
	lldpRemChassisId := []g.SnmpPDU{}
	lldpRemPortId := []g.SnmpPDU{}
	erpsEnabled := []g.SnmpPDU{}
	linkDataResult := []Link{}
	var blockedPort []string
	sysObjectId := gatSystemObjectId(targetIp)
	lldpRemChassisIdOid := "1.0.8802.1.1.2.1.4.1.1.5"
	lldpRemPortIdOid := "1.0.8802.1.1.2.1.4.1.1.7"
	erpsEnableOid := sysObjectId + ".4.4.1"
	lldpRemChassisId = Bulk(targetIp, lldpRemChassisIdOid, nil)
	lldpRemPortId = Bulk(targetIp, lldpRemPortIdOid, nil)
	erpsEnabled = Bulk(targetIp, erpsEnableOid, nil)
	if len(erpsEnabled) > 0 {
		for _, pdus := range erpsEnabled {
			if pdus.Value == 1 {
				blockedPort = getErpsBlockedPort(sysObjectId, targetIp)
			}
		}
	}
	if len(lldpRemChassisId) == len(lldpRemPortId) {
		for index, pdus := range lldpRemChassisId {
			linkData, _ := createLldp(pdus, lldpRemPortId[index], sourceChesId, blockedPort)
			linkDataResult = append(linkDataResult, linkData)
		}
	}

	return linkDataResult
}

func createLldp(pdu g.SnmpPDU, ppdu g.SnmpPDU, sourceChesId string, blockedPort []string) (Link, error) {
	linkData := Link{}
	remotePortValue := string(ppdu.Value.([]byte))
	if len(remotePortValue) != 8 {
		return linkData, nil
	}
	intVar, _ := strconv.Atoi(remotePortValue[5:])

	sourcePortvalue := "port" + getLocalPort((strings.Replace(ppdu.Name, "1.0.8802.1.1.2.1.4.1.1.7", "", -1)))
	b := ToMacString(pdu.Value)
	isBlockedPort := isSourcePortBlocked(sourcePortvalue, blockedPort)

	var edgeData string
	if sourceChesId < string(b) {
		edgeData = sourceChesId + "_" + string(b)
	} else {
		edgeData = string(b) + "_" + sourceChesId
	}

	linkData = Link{Source: sourceChesId, Target: string(b), SourcePort: sourcePortvalue, TargetPort: "port" + strconv.Itoa(intVar), EdgeData: edgeData, LinkFlow: !isBlockedPort, BlockedPort: isBlockedPort}

	return linkData, nil
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
			buffer.WriteRune(':')
		}
	}
	return buffer.String()
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func getLocalPort(s string) string {
	splitValue := strings.Split(s, ".")[3]
	return splitValue
}

func realTopologyData(arrayLink []Link, inputLink Link) bool {
	linkDataResult := true

	for _, linkData := range arrayLink {
		if linkData.EdgeData == inputLink.EdgeData {
			linkDataResult = false
		}
	}
	return linkDataResult
}

func gatSystemObjectId(targetIp string) (sysObjectId string) {
	var systemObjectId string
	oids := []string{"1.3.6.1.2.1.1.2.0"}
	result, err2 := GetOids(targetIp, oids, nil)
	if err2 != nil {
		q.Q("Get() err: ", err2)
	}
	for _, variable := range result {
		switch variable.Type {
		case g.ObjectIdentifier:
			t, ok := variable.Value.(string)
			if ok {
				systemObjectId = t
			}

		default:
		}
	}
	return systemObjectId
}

func getErpsBlockedPort(systemObjectId string, targetIp string) []string {
	var erpsVlanId int
	var blockedPort []string
	erpsRsapVlanOid := systemObjectId + ".4.4.3.1.1"
	erpsDataOid := systemObjectId + ".4.4.3.1"
	erpsWPortStatusOid := systemObjectId + ".4.4.3.1.4."
	erpsEPortStatusOid := systemObjectId + ".4.4.3.1.5."
	erpsVlanPdus := []g.SnmpPDU{}
	erpsData := []g.SnmpPDU{}
	erpsVlanPdus = Bulk(targetIp, erpsRsapVlanOid, nil)
	erpsData = Bulk(targetIp, erpsDataOid, nil)
	if len(erpsVlanPdus) > 0 {
		erpsVlanId = int(erpsVlanPdus[0].Value.(int))
	}
	if len(erpsData) > 0 {
		for _, pdus := range erpsData {
			if pdus.Name == erpsWPortStatusOid+strconv.Itoa(erpsVlanId) && int(pdus.Value.(int)) == 2 {
				blockedPort = append(blockedPort, "port"+string(erpsData[1].Value.([]byte)))
			}
			if pdus.Name == erpsEPortStatusOid+strconv.Itoa(erpsVlanId) && int(pdus.Value.(int)) == 2 {
				blockedPort = append(blockedPort, "port"+string(erpsData[2].Value.([]byte)))
			}
		}
	}
	return blockedPort
}

func isSourcePortBlocked(sourcePort string, blockedPort []string) bool {
	isblockedPort := false
	for _, port := range blockedPort {
		if port == sourcePort {
			isblockedPort = true
			q.Q("enter in to check blocked port")
		}
	}
	return isblockedPort
}

func PublishTopology(topologyData Topology) error {
	// send all devices info to root
	if QC.Root == "" {
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
	resp, err := PostWithToken(QC.Root+"/api/v1/topology", QC.AdminToken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		q.Q(err, QC.Root)
	}
	if resp != nil {
		res := make(map[string]interface{})
		_ = json.NewDecoder(resp.Body).Decode(&res)
		q.Q(res)
	}
	return nil
}

func InsertTopology(topoKeys string, topoDesc Topology) bool {
	// insert a topology into the topology list
	QC.DevMutex.Lock()
	QC.TopologyData[topoKeys] = topoDesc
	QC.DevMutex.Unlock()

	return true
}

type SnmpOption struct {
	Port      uint16
	Community string
	Version   g.SnmpVersion
	Timeout   time.Duration
	Retries   int
}

var DefaultSnmpOption = SnmpOption{
	Port:      161,
	Community: "public",
	Version:   g.Version2c,
	Timeout:   time.Second,
	Retries:   1,
}

func GetOids(target string, oids []string, opt *SnmpOption) (results []g.SnmpPDU, err error) {
	if opt == nil {
		opt = &DefaultSnmpOption
	}

	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err = client.Connect()
	if err != nil {
		logrus.Errorf("Connect() err: %v", err)
		return
	}
	defer client.Conn.Close()
	pkt, err := client.Get(oids)
	if err != nil {
		return
	}
	return pkt.Variables, err
}

// Bulk get same type below oid
func Bulk(target string, oid string, opt *SnmpOption) []g.SnmpPDU {
	results := []g.SnmpPDU{}
	if opt == nil {
		opt = &DefaultSnmpOption
	}
	client := g.GoSNMP{
		Target:    target,
		Community: opt.Community,
		Port:      opt.Port,
		Version:   opt.Version,
		Timeout:   opt.Timeout,
	}

	err := client.Connect()
	if err != nil {
		logrus.Errorf("Connect() err: %v", err)
		return results
	}
	defer client.Conn.Close()
	err = client.BulkWalk(oid, func(pdu g.SnmpPDU) error {
		results = append(results, pdu)
		return nil
	})
	if err != nil {
		return results
	}
	return results
}
