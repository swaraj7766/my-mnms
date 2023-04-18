package mnms

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/awcullen/opcua/client"
	opcuaserver "github.com/awcullen/opcua/server"
	"github.com/awcullen/opcua/ua"
	"github.com/pkg/errors"
	"github.com/qeof/q"
)

const testnodeset = `
<UANodeSet xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:uax="http://opcfoundation.org/UA/2008/02/Types.xsd" xmlns="http://opcfoundation.org/UA/2011/03/UANodeSet.xsd" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
    <NamespaceUris>
        <Uri>http://github.com/awcullen/opcua/testserver/</Uri>
    </NamespaceUris>
	 <Aliases>
        <Alias Alias="Boolean">i=1</Alias>
        <Alias Alias="SByte">i=2</Alias>
        <Alias Alias="Byte">i=3</Alias>
        <Alias Alias="Int16">i=4</Alias>
        <Alias Alias="UInt16">i=5</Alias>
        <Alias Alias="Int32">i=6</Alias>
        <Alias Alias="UInt32">i=7</Alias>
        <Alias Alias="Int64">i=8</Alias>
        <Alias Alias="UInt64">i=9</Alias>
        <Alias Alias="Float">i=10</Alias>
        <Alias Alias="Double">i=11</Alias>
        <Alias Alias="String">i=12</Alias>
        <Alias Alias="DateTime">i=13</Alias>
        <Alias Alias="Guid">i=14</Alias>
        <Alias Alias="ByteString">i=15</Alias>
        <Alias Alias="XmlElement">i=16</Alias>
        <Alias Alias="NodeId">i=17</Alias>
        <Alias Alias="StatusCode">i=19</Alias>
        <Alias Alias="QualifiedName">i=20</Alias>
        <Alias Alias="LocalizedText">i=21</Alias>
        <Alias Alias="Number">i=26</Alias>
        <Alias Alias="Integer">i=27</Alias>
        <Alias Alias="UInteger">i=28</Alias>
        <Alias Alias="Organizes">i=35</Alias>
        <Alias Alias="HasModellingRule">i=37</Alias>
        <Alias Alias="HasTypeDefinition">i=40</Alias>
        <Alias Alias="HasSubtype">i=45</Alias>
        <Alias Alias="HasProperty">i=46</Alias>
        <Alias Alias="HasComponent">i=47</Alias>
        <Alias Alias="NodeClass">i=257</Alias>
        <Alias Alias="Duration">i=290</Alias>
        <Alias Alias="UtcTime">i=294</Alias>
        <Alias Alias="Argument">i=296</Alias>
        <Alias Alias="Range">i=884</Alias>
        <Alias Alias="EUInformation">i=887</Alias>
        <Alias Alias="EnumValueType">i=7594</Alias>
        <Alias Alias="TimeZoneDataType">i=8912</Alias>
    </Aliases>
<UAObject NodeId="i=1001" BrowseName="0:Boiler" ParentNodeId="i=85">
    <DisplayName>Boiler</DisplayName>
    <Description>The base type for all object nodes.</Description>
    <References>
      <Reference ReferenceType="Organizes" IsForward="false">i=85</Reference>
      <Reference ReferenceType="HasTypeDefinition">i=58</Reference>
      <Reference ReferenceType="HasComponent">i=1002</Reference>
      <Reference ReferenceType="HasComponent">i=1003</Reference>
    </References>
  </UAObject>
  <UAVariable  DataType="Float" NodeId="i=1002" BrowseName="0:Temperature" ParentNodeId="i=1001" UserAccessLevel="3" AccessLevel="3">
    <DisplayName>Temperature</DisplayName>
    <Description>Temperature</Description>
    <References>
      <Reference ReferenceType="HasComponent" IsForward="false">i=1001</Reference>
      <Reference ReferenceType="HasTypeDefinition">i=63</Reference>
    </References>
    <Value>
      <uax:Float>0.5</uax:Float>
    </Value>
  </UAVariable>
  <UAVariable DataType="Float" NodeId="i=1003" BrowseName="0:Pressure" ParentNodeId="i=1001"  UserAccessLevel="3" AccessLevel="3">
    <DisplayName>Pressure</DisplayName>
    <Description>Pressure</Description>
    <References>
      <Reference ReferenceType="HasComponent" IsForward="false">i=1001</Reference>
      <Reference ReferenceType="HasTypeDefinition">i=63</Reference>
    </References>
    <Value>
      <uax:Float>0.99</uax:Float>
    </Value>
  </UAVariable>
</UANodeSet>
`

var (
	port            = 4840
	host, _         = os.Hostname()
	SoftwareVersion = "0.1.0"
	serverURL       = fmt.Sprintf("opc.tcp://%s:%d", host, port)
)
var ConnectError = errors.New("connection error")

func ensurePKI() error {

	// check if ./pki already exists
	if _, err := os.Stat("./pki"); !os.IsNotExist(err) {
		return nil
	}

	// make a pki directory, if not exist
	if err := os.MkdirAll("./pki", os.ModeDir|0755); err != nil {
		return err
	}
	// create a server cert in ./pki/server.crt
	if err := createNewCertificate("testserver", "./pki/server.crt", "./pki/server.key"); err != nil {
		return err
	}
	// create a client cert in ./pki
	if err := createNewCertificate("test-client", "./pki/client.crt", "./pki/client.key"); err != nil {
		return err
	}

	return nil
}
func createNewCertificate(appName, certFile, keyFile string) error {

	// create a keypair.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return ua.BadCertificateInvalid
	}

	// get local hostname.
	host, _ := os.Hostname()

	// get local ip address.
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return ua.BadCertificateInvalid
	}
	conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)

	// create a certificate.
	applicationURI, _ := url.Parse(fmt.Sprintf("urn:%s:%s", host, appName))
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	subjectKeyHash := sha1.New()
	subjectKeyHash.Write(key.PublicKey.N.Bytes())
	subjectKeyId := subjectKeyHash.Sum(nil)

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{CommonName: appName},
		SubjectKeyId:          subjectKeyId,
		AuthorityKeyId:        subjectKeyId,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageContentCommitment | x509.KeyUsageKeyEncipherment | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{host},
		IPAddresses:           []net.IP{localAddr.IP},
		URIs:                  []*url.URL{applicationURI},
	}

	rawcrt, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return ua.BadCertificateInvalid
	}

	if f, err := os.Create(certFile); err == nil {
		block := &pem.Block{Type: "CERTIFICATE", Bytes: rawcrt}
		if err := pem.Encode(f, block); err != nil {
			f.Close()
			return err
		}
		f.Close()
	} else {
		return err
	}

	if f, err := os.Create(keyFile); err == nil {
		block := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
		if err := pem.Encode(f, block); err != nil {
			f.Close()
			return err
		}
		f.Close()
	} else {
		return err
	}

	return nil
}

func NewOpcuaServer(endpointurl string) *OpcuaServer {
	if err := ensurePKI(); err != nil {
		q.Q("Error creating PKI.")
		DoExit(1)
	}

	srv, err := opcuaserver.New(
		ua.ApplicationDescription{
			ApplicationURI: fmt.Sprintf("urn:%s:testserver", host),
			ProductURI:     "http://github.com/awcullen/opcua",
			ApplicationName: ua.LocalizedText{
				Text:   fmt.Sprintf("testserver@%s", host),
				Locale: "en",
			},
			ApplicationType:     ua.ApplicationTypeServer,
			GatewayServerURI:    "",
			DiscoveryProfileURI: "",
			DiscoveryURLs:       []string{endpointurl},
		},
		"./pki/server.crt",
		"./pki/server.key",
		endpointurl,
		opcuaserver.WithBuildInfo(
			ua.BuildInfo{
				ProductURI:       "http://github.com/awcullen/opcua",
				ManufacturerName: "awcullen",
				ProductName:      "testserver",
				SoftwareVersion:  SoftwareVersion,
			}),
		opcuaserver.WithAuthenticateUserNameIdentityFunc(func(userIdentity ua.UserNameIdentity, applicationURI string, endpointURL string) error {

			// log.Printf("Login user: %s from %s\n", userIdentity.UserName, applicationURI)
			return nil
		}),

		opcuaserver.WithRolePermissions([]ua.RolePermissionType{{RoleID: ua.ObjectIDWellKnownRoleAnonymous, Permissions: (ua.PermissionTypeBrowse | ua.PermissionTypeRead | ua.PermissionTypeReadHistory | ua.PermissionTypeReceiveEvents | ua.PermissionTypeWrite)}}),
		opcuaserver.WithAnonymousIdentity(true),
		opcuaserver.WithSecurityPolicyNone(true),
		opcuaserver.WithInsecureSkipVerify(),
		opcuaserver.WithServerDiagnostics(true),
		// server.WithTrace(),
	)
	if err != nil {
		q.Q(err)
		DoExit(1)
	}
	return &OpcuaServer{srv: srv}
}

type OpcuaServer struct {
	srv *opcuaserver.Server
}

// OpcuaRun opcua server run
func (o *OpcuaServer) OpcuaRun() {
	// create directory with certificate and key, if not found.
	nm := o.srv.NamespaceManager()
	if err := nm.LoadNodeSetFromBuffer([]byte(testnodeset)); err != nil {
		DoExit(2)
	}
	c := make(chan bool)

	go func() {
		c <- true
		q.Q("Starting server", o.srv.LocalDescription().ApplicationName.Text, o.srv.EndpointURL())
		if err := o.srv.ListenAndServe(); err != ua.BadServerHalted {
			q.Q(errors.Wrap(err, "Error starting server"))
		}
	}()
	<-c
}

// OpcuaShutdown shutdonw opcua server
func (o *OpcuaServer) OpcuaShutdown() error {
	return o.srv.Close()
}

// NewOpcuaClient create opcua client
func NewOpcuaClient() *OpcuaClient {
	if err := ensurePKI(); err != nil {
		q.Q("Error creating PKI.")
		DoExit(1)
	}
	return &OpcuaClient{timeout: time.Second * 3}
}

type OpcuaClient struct {
	ch      *client.Client
	timeout time.Duration
}

// Connect connect to endpointURL
func (o *OpcuaClient) Connect(endpointURL string, opts ...client.Option) error {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()
	ch, err := client.Dial(
		ctx,
		endpointURL,
		opts...,
	)
	if err == nil {
		o.ch = ch
	}
	return err
}
func (o *OpcuaClient) checkConnection() error {
	if o.ch == nil {
		return ConnectError
	}
	return nil
}

// FindServers
func (o *OpcuaClient) FindServers(request *ua.FindServersRequest) (*ua.FindServersResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()
	return client.FindServers(ctx, request)
}

// FindServers
func (o *OpcuaClient) GetEndpoints(request *ua.GetEndpointsRequest) (*ua.GetEndpointsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()
	return client.GetEndpoints(ctx, request)
}

// Close
func (o *OpcuaClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()
	return o.ch.Close(ctx)
}

// WriteNodeID
func (o *OpcuaClient) WriteNodeID(req *ua.WriteRequest) (*ua.WriteResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()
	return o.ch.Write(ctx, req)
}

// Browse
func (o *OpcuaClient) Browse(req *ua.BrowseRequest) (*ua.BrowseResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.ch.Browse(ctx, req)
}

// ReadNodeID
func (o *OpcuaClient) ReadNodeID(req *ua.ReadRequest) (*ua.ReadResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.ch.Read(ctx, req)

}

// CreateSubscription
func (o *OpcuaClient) CreateSubscription(req *ua.CreateSubscriptionRequest) (*ua.CreateSubscriptionResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.ch.CreateSubscription(ctx, req)

}

// CreateMonitoredItems
func (o *OpcuaClient) CreateMonitoredItems(req *ua.CreateMonitoredItemsRequest) (*ua.CreateMonitoredItemsResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.ch.CreateMonitoredItems(ctx, req)

}

// CreateMonitoredItems
func (o *OpcuaClient) Publish(req *ua.PublishRequest) (*ua.PublishResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.ch.Publish(ctx, req)

}

// CreateMonitoredItems
func (o *OpcuaClient) DeleteMonitoredItems(req *ua.DeleteMonitoredItemsRequest) (*ua.DeleteMonitoredItemsResponse, error) {
	err := o.checkConnection()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), o.timeout)
	defer cancel()

	return o.ch.DeleteMonitoredItems(ctx, req)

}

// ReadVariableAttributes read VariableAttributes
func (o *OpcuaClient) ReadVariableAttributes(id ua.NodeID) (ua.VariableAttributes, error) {
	req := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: id, AttributeID: ua.AttributeIDNodeID},
			{NodeID: id, AttributeID: ua.AttributeIDDisplayName},
			{NodeID: id, AttributeID: ua.AttributeIDDescription},
			{NodeID: id, AttributeID: ua.AttributeIDWriteMask},
			{NodeID: id, AttributeID: ua.AttributeIDUserWriteMask},
			{NodeID: id, AttributeID: ua.AttributeIDValue},
			{NodeID: id, AttributeID: ua.AttributeIDDataType},
			{NodeID: id, AttributeID: ua.AttributeIDValueRank},
			{NodeID: id, AttributeID: ua.AttributeIDArrayDimensions},
			{NodeID: id, AttributeID: ua.AttributeIDAccessLevel},
			{NodeID: id, AttributeID: ua.AttributeIDUserAccessLevel},
			{NodeID: id, AttributeID: ua.AttributeIDMinimumSamplingInterval},
			{NodeID: id, AttributeID: ua.AttributeIDHistorizing},
		},
	}
	res, err := o.ReadNodeID(req)
	if err != nil {
		return ua.VariableAttributes{}, err
	}
	value := ua.VariableAttributes{}
	for i := range res.Results {
		if res.Results[i].StatusCode.IsGood() {
			switch i {
			case 0:
				v, ok := res.Results[i].Value.(uint32)
				if ok {
					value.SpecifiedAttributes = v
				}
			case 1:
				v, ok := res.Results[i].Value.(ua.LocalizedText)
				if ok {
					value.DisplayName = v
				}
			case 2:
				v, ok := res.Results[i].Value.(ua.LocalizedText)
				if ok {
					value.Description = v
				}
			case 3:
				v, ok := res.Results[i].Value.(uint32)
				if ok {
					value.WriteMask = v
				}
			case 4:
				v, ok := res.Results[i].Value.(uint32)
				if ok {
					value.UserWriteMask = v
				}
			case 5:
				value.Value = res.Results[i].Value
			case 6:
				v, ok := res.Results[i].Value.(ua.NodeID)
				if ok {
					value.DataType = v
				}
			case 7:
				v, ok := res.Results[i].Value.(int32)
				if ok {
					value.ValueRank = v
				}
			case 8:
				v, ok := res.Results[i].Value.([]uint32)
				if ok {
					value.ArrayDimensions = v
				}
			case 9:
				v, ok := res.Results[i].Value.(uint8)
				if ok {
					value.AccessLevel = v
				}
			case 10:
				v, ok := res.Results[i].Value.(uint8)
				if ok {
					value.UserAccessLevel = v
				}
			case 11:
				v, ok := res.Results[i].Value.(float64)
				if ok {
					value.MinimumSamplingInterval = v
				}
			case 12:
				v, ok := res.Results[i].Value.(bool)
				if ok {
					value.Historizing = v
				}
			}
		}
	}

	return value, nil
}

// BrowseReference Browse ReferenceDescription
func (o *OpcuaClient) BrowseReference(id ua.NodeID, dir ua.BrowseDirection) ([]ua.ReferenceDescription, error) {
	req := &ua.BrowseRequest{
		NodesToBrowse: []ua.BrowseDescription{
			{
				NodeID:          id,
				BrowseDirection: dir,
				ReferenceTypeID: ua.ReferenceTypeIDHierarchicalReferences,
				IncludeSubtypes: true,
				ResultMask:      uint32(ua.BrowseResultMaskAll),
			},
		},
	}
	res, err := o.Browse(req)
	if err != nil {
		return []ua.ReferenceDescription{}, err
	}
	value := []ua.ReferenceDescription{}
	for _, r := range res.Results {
		if r.StatusCode.IsGood() {
			value = append(value, r.References...)
		}
	}

	return value, nil
}

/*========  opcua Cmd ========*/
var opcuaclient *OpcuaClient

func CheckOpcuaClient() error {
	if opcuaclient == nil {
		return ConnectError
	}
	return nil
}

// Opcua connect setting.
//
// Usage : opcua connect [url]
//
//	[url]         : connect to url
//
// Example :
//
//	opcua connect opc.tcp://127.0.0.1:4840
func OpcuaConnectCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	url := ws[2]
	opclient := NewOpcuaClient()
	err := opclient.Connect(url, client.WithInsecureSkipVerify(), client.WithClientCertificateFile("./pki/client.crt", "./pki/client.key"))
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	opcuaclient = opclient
	cmdinfo.Status = "ok"
	return cmdinfo
}

// Opcua read setting.
//
// Usage : opcua read [node id]
//
//	[node id]     : opcua node id
//
// Example :
//
//	opcua read i=1002
func OpcuaReadCmd(cmdinfo *CmdInfo) *CmdInfo {
	err := CheckOpcuaClient()
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	id := ws[2]
	nodeid := ua.ParseNodeID(id)
	v, err := opcuaclient.ReadVariableAttributes(nodeid)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	b, err := json.Marshal(v)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	cmdinfo.Result = string(b)
	q.Q(v)
	return cmdinfo
}

// Opcua browse setting.
//
// Usage : opcua browse [node id]
//
//	[node id]     : opcua node id
//
// Example :
//
//	opcua browse i=85
func OpcuaBrowseReferenceCmd(cmdinfo *CmdInfo) *CmdInfo {
	err := CheckOpcuaClient()
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	id := ws[2]
	nodeid := ua.ParseNodeID(id)
	v, err := opcuaclient.BrowseReference(nodeid, ua.BrowseDirectionForward)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	b, err := json.Marshal(v)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	cmdinfo.Result = string(b)
	q.Q(v)
	return cmdinfo
}

var opcuaNotificationFlag = false
var opcuMap = map[uint32]ua.NodeID{}

// Opcua subcribe setting.
//
// Usage : opcua sub [node id]
//
//	[node id]     : opcua node id
//
// Example :
//
//	opcua sub i=1002
func OpcuaSubscribeCmd(cmdinfo *CmdInfo) *CmdInfo {
	err := CheckOpcuaClient()
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	req := &ua.CreateSubscriptionRequest{
		RequestedPublishingInterval: 1000.0,
		RequestedMaxKeepAliveCount:  30,
		RequestedLifetimeCount:      30 * 3,
		PublishingEnabled:           true,
	}
	res, err := opcuaclient.CreateSubscription(req)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error() + ", Error creating subscription"
		return cmdinfo
	}
	id := ws[2]
	nodeid := ua.ParseNodeID(id)

	req2 := &ua.CreateMonitoredItemsRequest{
		SubscriptionID:     res.SubscriptionID,
		TimestampsToReturn: ua.TimestampsToReturnBoth,
		ItemsToCreate: []ua.MonitoredItemCreateRequest{
			{
				ItemToMonitor: ua.ReadValueID{
					AttributeID: ua.AttributeIDValue,
					NodeID:      nodeid,
				},
				MonitoringMode: ua.MonitoringModeReporting,
				RequestedParameters: ua.MonitoringParameters{
					ClientHandle: 42, QueueSize: 1, DiscardOldest: true, SamplingInterval: 500.0,
				},
			},
		},
	}
	res2, err := opcuaclient.CreateMonitoredItems(req2)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error() + ",Error creating item"
		return cmdinfo
	}
	_ = res2

	cmdinfo.Status = "ok"
	cmdinfo.Result = fmt.Sprintf("SubscriptionID:%v,MonitoredItemID:%v", res.SubscriptionID, res2.Results[0].MonitoredItemID)
	opcuMap[res.SubscriptionID] = nodeid
	if !opcuaNotificationFlag {
		opcuaNotificationFlag = true
		go runOpcuaDataChangeNotification()
	}
	return cmdinfo
}

// runOpcuaDataChangeNotification
func runOpcuaDataChangeNotification() {
	req3 := &ua.PublishRequest{
		RequestHeader:                ua.RequestHeader{TimeoutHint: 60000},
		SubscriptionAcknowledgements: []ua.SubscriptionAcknowledgement{},
	}
	for {
		res3, err := opcuaclient.Publish(req3)
		if err != nil {
			opcuaNotificationFlag = false
			q.Q("err:", err)
			m := fmt.Errorf("connection error:%v", err)
			err := SendSyslog(LOG_ALERT, "opcua ", m.Error())
			if err != nil {
				q.Q("syslog err:", err)
			}
			break
		}

		// loop thru all the notifications.
		for _, data := range res3.NotificationMessage.NotificationData {
			switch body := data.(type) {
			case ua.DataChangeNotification:
				for _, z := range body.MonitoredItems {
					if z.ClientHandle == 42 && z.Value.StatusCode.IsGood() {
						if id, ok := opcuMap[res3.SubscriptionID]; ok {
							m := fmt.Sprintf("SubscriptionID:%v,nodeid:%v,value:%v", res3.SubscriptionID, id, z.Value.Value)
							err := SendSyslog(LOG_ALERT, "opcua ", m)
							if err != nil {
								q.Q("syslog err:", err)
							}
							q.Q("SubscriptionID :", res3.SubscriptionID, "nodeid:", id, "value:", z.Value.Value)
						} else {
							m := fmt.Sprintf("SubscriptionID:%v,value:%v", res3.SubscriptionID, z.Value.Value)
							err := SendSyslog(LOG_ALERT, "opcua ", m)
							if err != nil {
								q.Q("syslog err:", err)
							}
							q.Q("SubscriptionID :", res3.SubscriptionID, "value:", z.Value.Value)
						}
					}
				}
			}
		}
	}
}

// DeleteOpcuSubscribeCmd

// Opcua delete subcribe setting.
//
// Usage : opcua deletesub [sub id] [monitor id]
//
//	[sub id]      : subscribe id
//	[monitor id]  : monitored item id
//
// Example :
//
//	opcua deletesub 1 1
func OpcuDeleteSubscribeCmd(cmdinfo *CmdInfo) *CmdInfo {
	err := CheckOpcuaClient()
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 4 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}
	subid, err := strconv.Atoi(ws[2])
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	monid, err := strconv.Atoi(ws[3])
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	req := &ua.DeleteMonitoredItemsRequest{SubscriptionID: uint32(subid), MonitoredItemIDs: []uint32{uint32(monid)}}

	res, err := opcuaclient.DeleteMonitoredItems(req)
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	for _, r := range res.Results {
		if !r.IsGood() {
			cmdinfo.Status = r.Error()
			return cmdinfo
		}
	}

	cmdinfo.Status = "ok"
	return cmdinfo
}

// OpcuCloseCmd

// Close opcua.
//
// Usage : opcua close
//
// Example :
//
//	opcua close
func OpcuCloseCmd(cmdinfo *CmdInfo) *CmdInfo {
	err := CheckOpcuaClient()
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 2 {
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}

	err = opcuaclient.Close()
	if err != nil {
		cmdinfo.Status = "error: " + err.Error()
		return cmdinfo
	}
	cmdinfo.Status = "ok"
	return cmdinfo
}
