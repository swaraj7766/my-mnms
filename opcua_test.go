package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"github.com/pkg/errors"
	"github.com/qeof/q"
)

func TestOpcFindServers(t *testing.T) {

	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	opacuclient := NewOpcuaClient()
	r, err := opacuclient.FindServers(&ua.FindServersRequest{EndpointURL: endpointURL})
	if err != nil {
		t.Fatal(err)
	}
	q.Q(r)
}

func TestOpcGetEndpoints(t *testing.T) {

	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	opacuclient := NewOpcuaClient()
	time.Sleep(time.Second * 1)
	r, err := opacuclient.GetEndpoints(&ua.GetEndpointsRequest{EndpointURL: endpointURL})
	if err != nil {
		t.Fatal(err)
	}
	q.Q(r)
}

func TestOpcWrite(t *testing.T) {
	fv := 42.0
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	id := ua.ParseNodeID("i=1002")
	value := ua.NewDataValue(float32(fv), 0, time.Time{}, 0, time.Time{}, 0)
	req := &ua.WriteRequest{
		NodesToWrite: []ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value:       value,
			},
		},
	}
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()

	res, err := opacuclient.WriteNodeID(req) //client.WithClientCertificateFile("./pki/client.crt", "./pki/client.key"),

	if err != nil {
		t.Fatal(err)
	}
	if res.Results[0].IsBad() {
		t.Error(errors.Wrap(res.Results[0], "Error Write"))
		return
	}
	readreq := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: id, AttributeID: ua.AttributeIDValue},
		},
	}

	r, err := opacuclient.ReadNodeID(readreq)
	if err != nil {
		t.Fatal(err)
	}

	if r.Results[0].StatusCode.IsBad() {
		t.Error(errors.Wrap(r.Results[0].StatusCode, "Error reading"))
		return
	}
	for _, result := range r.Results {
		if result.StatusCode.IsGood() {
			t.Logf("%s: %v", id, result.Value)
			if value.Value == result.Value {
				q.Q(result.Value)
			} else {
				t.Fatalf("expect%v,actual:%v", value.Value, result.Value)
			}

		} else {
			t.Error(errors.Wrap(result.StatusCode, "Error reading node"))
		}
	}

}

func TestOpcBrowse(t *testing.T) {
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()

	req := &ua.BrowseRequest{
		NodesToBrowse: []ua.BrowseDescription{
			{
				NodeID:          ua.ParseNodeID("i=85"),
				BrowseDirection: ua.BrowseDirectionForward,
				ReferenceTypeID: ua.ReferenceTypeIDHierarchicalReferences,
				IncludeSubtypes: true,
				ResultMask:      uint32(ua.BrowseResultMaskAll),
			},
		},
	}
	res, err := opacuclient.Browse(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Results[0].StatusCode.IsBad() {
		t.Error(errors.Wrap(res.Results[0].StatusCode, "Error browsing"))
		return
	}
	for _, r := range res.Results[0].References {
		q.Q(r)
	}
}

func TestOpcReadServerStatus(t *testing.T) {
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	req := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: ua.VariableIDServerServerStatus, AttributeID: ua.AttributeIDValue},
		},
	}
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()
	res, err := opacuclient.ReadNodeID(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Results[0].StatusCode.IsBad() {
		t.Error(errors.Wrap(res.Results[0].StatusCode, "Error reading ServerStatus"))
		return
	}
	status, ok := res.Results[0].Value.(ua.ServerStatusDataType)
	if !ok {
		t.Error(errors.New("Error decoding ServerStatusDataType"))
		return
	}
	q.Q("Server status:")
	q.Q(status.BuildInfo.ProductName)
	q.Q(status.BuildInfo.SoftwareVersion)
	q.Q(status.BuildInfo.ManufacturerName)
	q.Q(status.State)
	q.Q(status.CurrentTime)
}

func TestOpcRead(t *testing.T) {
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	id := ua.ParseNodeID("i=1002")
	req := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: id, AttributeID: ua.AttributeIDValue},
		},
	}
	time.Sleep(time.Second * 1)
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()
	res, err := opacuclient.ReadNodeID(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.Results[0].StatusCode.IsBad() {
		t.Error(errors.Wrap(res.Results[0].StatusCode, "Error reading"))
		return
	}
	for _, result := range res.Results {
		if result.StatusCode.IsGood() {
			t.Logf("%s: %v", id, result.Value)
			q.Q(result.Value)
		} else {
			t.Error(errors.Wrap(result.StatusCode, "Error reading node"))
		}
	}

}

func TestOpcReadAttributes(t *testing.T) {
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	id := ua.ParseNodeID("i=1002")
	req := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: id, AttributeID: ua.AttributeIDNodeID},
			{NodeID: id, AttributeID: ua.AttributeIDNodeClass},
			{NodeID: id, AttributeID: ua.AttributeIDBrowseName},
			{NodeID: id, AttributeID: ua.AttributeIDDisplayName},
			{NodeID: id, AttributeID: ua.AttributeIDDescription},
			{NodeID: id, AttributeID: ua.AttributeIDValue},
			{NodeID: id, AttributeID: ua.AttributeIDRolePermissions},
		},
	}
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()
	res, err := opacuclient.ReadNodeID(req)
	if err != nil {
		t.Fatal(err)
	}
	for i, result := range res.Results {
		if result.StatusCode.IsGood() {
			t.Logf("%d: %v", req.NodesToRead[i].AttributeID, result.Value)
			q.Q(req.NodesToRead[i].AttributeID, result.Value)
		} else {
			q.Q("unsupport", req.NodesToRead[i].AttributeID, result.StatusCode)
		}
	}

}

func TestSubscribe(t *testing.T) {
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()
	req := &ua.CreateSubscriptionRequest{
		RequestedPublishingInterval: 1000.0,
		RequestedMaxKeepAliveCount:  30,
		RequestedLifetimeCount:      30 * 3,
		PublishingEnabled:           true,
	}
	res, err := opacuclient.CreateSubscription(req)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating subscription"))
		return
	}
	id := ua.ParseNodeID("i=1002")
	req2 := &ua.CreateMonitoredItemsRequest{
		SubscriptionID:     res.SubscriptionID,
		TimestampsToReturn: ua.TimestampsToReturnBoth,
		ItemsToCreate: []ua.MonitoredItemCreateRequest{
			{
				ItemToMonitor: ua.ReadValueID{
					AttributeID: ua.AttributeIDValue,
					NodeID:      id,
				},
				MonitoringMode: ua.MonitoringModeReporting,
				RequestedParameters: ua.MonitoringParameters{
					ClientHandle: 42, QueueSize: 1, DiscardOldest: true, SamplingInterval: 500.0,
				},
			},
		},
	}
	res2, err := opacuclient.CreateMonitoredItems(req2)
	if err != nil {
		t.Error(errors.Wrap(err, "Error creating item"))
	}
	_ = res2
	req3 := &ua.PublishRequest{
		RequestHeader:                ua.RequestHeader{TimeoutHint: 60000},
		SubscriptionAcknowledgements: []ua.SubscriptionAcknowledgement{},
	}
	numChanges := 0
	for numChanges < 3 {
		res3, err := opacuclient.Publish(req3)
		if err != nil {
			t.Error(errors.Wrap(err, "Error publishing"))
			break
		}

		// loop thru all the notifications.
		for _, data := range res3.NotificationMessage.NotificationData {
			switch body := data.(type) {
			case ua.DataChangeNotification:
				for _, z := range body.MonitoredItems {
					if z.ClientHandle == 42 {
						t.Logf("value: %s", z.Value.Value)
						q.Q("value:", z.Value.Value)
						numChanges++
					}
				}
			}
			err := Opcwrite(id, 50+numChanges)
			if err != nil {
				t.Fatal(err)
			}
		}

	}

}

func Opcwrite(id ua.NodeID, nmuber int) error {
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect(endpointURL, client.WithInsecureSkipVerify() /*,client.WithUserNameIdentity("test", "test")*/)
	if err != nil {
		return err
	}
	defer opacuclient.Close()
	value := ua.NewDataValue(float32(nmuber), 0, time.Time{}, 0, time.Time{}, 0)
	req := &ua.WriteRequest{
		NodesToWrite: []ua.WriteValue{
			{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value:       value,
			},
		},
	}
	_, err = opacuclient.WriteNodeID(req)

	if err != nil {
		return err
	}
	return nil
}

func TestOpcReadVariableAttributes(t *testing.T) {

	id := ua.ParseNodeID("i=2994")
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect("opc.tcp://opcua.123mc.com:4840/", client.WithInsecureSkipVerify(), client.WithClientCertificateFile("./pki/client.crt", "./pki/client.key"))
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()
	res, err := opacuclient.ReadVariableAttributes(id)
	if err != nil {
		t.Fatal(err)
	}
	q.Q(res)

}

func TestOpcBrowseReference(t *testing.T) {
	opacuclient := NewOpcuaClient()
	err := opacuclient.Connect("opc.tcp://opcua.123mc.com:4840/", client.WithInsecureSkipVerify(), client.WithClientCertificateFile("./pki/client.crt", "./pki/client.key"))
	if err != nil {
		t.Fatal(err)
	}
	defer opacuclient.Close()

	v, err := opacuclient.BrowseReference(ua.ParseNodeID("i=2253"), ua.BrowseDirectionForward)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range v {
		q.Q(value)
	}

}

func opcconnect(url string) error {
	cmd := "opcua connect " + url
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		return fmt.Errorf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	token, err := GetToken("admin")
	if err != nil {
		return err
	}

	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	if resp != nil && resp.StatusCode != 200 {
		return fmt.Errorf("error: post status %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	return nil
}

func TestOpcConnectCmd(t *testing.T) {
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()

	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	go func() {
		HTTPMain()
		GwdMain()
	}()
	err := opcconnect(endpointURL)
	if err != nil {
		t.Fatal(err)
	}

	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}

	_ = CheckCommands()
	time.Sleep(6 * time.Second)

	url := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(fmt.Sprintf("opcua connect %s", endpointURL))
	resp, err := GetWithToken(url, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get error:%v", err)
	}
	defer resp.Body.Close()

	if resp != nil {
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(commands)
		for _, v := range commands {
			if v.Status != "ok" {
				t.Fatal(v.Status)
			}
		}
	}

}

func TestOpcReadCmd(t *testing.T) {
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	go func() {
		HTTPMain()
		GwdMain()
	}()
	err := opcconnect(endpointURL)
	if err != nil {
		t.Fatal(err)
	}

	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	cmd := "opcua read i=1002"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		t.Fatalf("post status %v", resp.StatusCode)
	}
	resp.Body.Close()

	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	url := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(cmd)
	resp, err = GetWithToken(url, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get error:%v", err)
	}

	if resp != nil {
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(commands)
		for _, v := range commands {
			if v.Status != "ok" {
				t.Fatal(v.Status)
			}
		}
	}
	resp.Body.Close()

}

func TestOpcBrowseCmd(t *testing.T) {
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	go func() {
		HTTPMain()
		GwdMain()
	}()
	err := opcconnect(endpointURL)
	if err != nil {
		t.Fatal(err)
	}

	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	cmd := "opcua browse i=85"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		t.Fatalf("post status %v", resp.StatusCode)
	}
	resp.Body.Close()

	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	url := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(cmd)
	resp, err = GetWithToken(url, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get error:%v", err)
	}

	if resp != nil {
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(commands)
		for _, v := range commands {
			if v.Status != "ok" {
				t.Fatal(v.Status)
			}
		}
	}
}

func TestOpcSubscribeCmd(t *testing.T) {
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	go func() {
		HTTPMain()
		GwdMain()
	}()
	err := opcconnect(endpointURL)
	if err != nil {
		t.Fatal(err)
	}

	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}

	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	cmd := "opcua sub i=1002"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	//insertcmd("opcua sub i=1003", &cmdinfo)*/
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		t.Fatalf("post status %v", resp.StatusCode)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	url := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(cmd)
	resp, err = GetWithToken(url, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get error:%v", err)
	}

	if resp != nil {
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(commands)
		for _, v := range commands {
			if v.Status != "ok" {
				t.Fatal(v.Status)
			}
		}
	}
	resp.Body.Close()

	time.Sleep(3 * time.Second)

	numChanges := 0
	for numChanges < 10 {
		numChanges++
		time.Sleep(1 * time.Second)

		err := Opcwrite(ua.ParseNodeID("i=1002"), numChanges)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestOpcDeleteSubscribeCmd(t *testing.T) {
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	go func() {
		HTTPMain()
		GwdMain()
	}()
	err := opcconnect(endpointURL)
	if err != nil {
		t.Fatal(err)
	}

	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	cmd := "opcua sub i=1002"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	//insertcmd("opcua sub i=1003", &cmdinfo)*/
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		t.Fatalf("post status %v", resp.StatusCode)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	urls := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(cmd)
	resp, err = GetWithToken(urls, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get error:%v", err)
	}

	if resp != nil {
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(commands)
		for _, v := range commands {
			if v.Status != "ok" {
				t.Fatal(v.Status)
			}
		}
	}

	time.Sleep(3 * time.Second)

	numChanges := 0
	for numChanges < 10 {
		numChanges++
		time.Sleep(1 * time.Second)

		if numChanges == 6 {
			err := Opcdeletesub("1", "1")
			if err != nil {
				t.Fatal(err)
			}
			_ = CheckCommands()
			time.Sleep(6 * time.Second)
			urls := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(fmt.Sprintf("opcua deletesub %v %v", 1, 1))
			resp, err = GetWithToken(urls, token)
			if err != nil || resp.StatusCode != 200 {
				t.Fatalf("get error:%v", err)
			}

			if resp != nil {
				commands := make(map[string]CmdInfo)
				err := json.NewDecoder(resp.Body).Decode(&commands)
				if err != nil {
					t.Fatal(err)
				}
				q.Q(commands)
			}
		}
		err := Opcwrite(ua.ParseNodeID("i=1002"), numChanges)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func Opcdeletesub(subid, monid string) error {
	cmd := fmt.Sprintf("opcua deletesub %v %v", subid, monid)
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	//insertcmd("opcua sub i=1003", &cmdinfo)*/
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		return fmt.Errorf("json marshal %v", err)
	}
	token, err := GetToken("admin")
	if err != nil {
		return err
	}
	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	if resp != nil && resp.StatusCode != 200 {
		q.Q("post status %v", resp.StatusCode)
		return fmt.Errorf("status code %v", resp.StatusCode)
	}
	return nil
}

func TestOpcCloseCmd(t *testing.T) {
	_ = cleanMNMSConfig()
	_ = InitDefaultMNMSConfigIfNotExist()
	defer func() {
		_ = cleanMNMSConfig()
	}()
	endpointURL := fmt.Sprintf("opc.tcp://localhost:%d", port)
	s := NewOpcuaServer(serverURL)
	s.OpcuaRun()
	defer func() {
		_ = s.OpcuaShutdown()

	}()
	time.Sleep(time.Second * 1)
	go func() {
		HTTPMain()
		GwdMain()
	}()
	err := opcconnect(endpointURL)
	if err != nil {
		t.Fatal(err)
	}
	token, err := GetToken("admin")
	if err != nil {
		t.Fatal(err)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	cmd := "opcua sub i=1002"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	//insertcmd("opcua sub i=1003", &cmdinfo)*/
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		t.Fatalf("json marshal %v", err)
	}
	q.Q(string(jsonBytes))
	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		t.Fatal(err)
	}
	if resp != nil && resp.StatusCode != 200 {
		t.Fatalf("post status %v", resp.StatusCode)
	}
	_ = CheckCommands()
	time.Sleep(6 * time.Second)
	urls := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape(cmd)
	resp, err = GetWithToken(urls, token)
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("get error:%v", err)
	}

	if resp != nil {
		commands := make(map[string]CmdInfo)
		err := json.NewDecoder(resp.Body).Decode(&commands)
		if err != nil {
			t.Fatal(err)
		}
		q.Q(commands)
		for _, v := range commands {
			if v.Status != "ok" {
				t.Fatal(v.Status)
			}
		}
	}

	time.Sleep(3 * time.Second)
	numChanges := 0
	for numChanges < 10 {
		numChanges++
		time.Sleep(1 * time.Second)

		if numChanges == 6 {
			err := OpcuaClose()
			if err != nil {
				t.Fatal(err)
			}
			_ = CheckCommands()
			time.Sleep(6 * time.Second)
			urls := "http://localhost:27182/api/v1/commands?cmd=" + url.QueryEscape("opcua close")
			resp, err = GetWithToken(urls, token)
			if err != nil || resp.StatusCode != 200 {
				t.Fatalf("get error:%v", err)
			}

			if resp != nil {
				commands := make(map[string]CmdInfo)
				err := json.NewDecoder(resp.Body).Decode(&commands)
				if err != nil {
					t.Fatal(err)
				}
				q.Q(commands)
				for _, v := range commands {
					if v.Status != "ok" {
						t.Fatal(v.Status)
					}
				}
			}
		}
		err := Opcwrite(ua.ParseNodeID("i=1002"), numChanges)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func OpcuaClose() error {
	cmd := "opcua close"
	cmdinfo := make(map[string]CmdInfo)
	insertcmd(cmd, &cmdinfo)
	//insertcmd("opcua sub i=1003", &cmdinfo)*/
	jsonBytes, err := json.Marshal(cmdinfo)
	if err != nil {
		return fmt.Errorf("json marshal %v", err)
	}
	token, err := GetToken("admin")
	if err != nil {
		return err
	}

	urlpath := fmt.Sprintf("http://localhost:%d/api/v1/commands", QC.Port)
	resp, err := PostWithToken(urlpath, token, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	if resp != nil && resp.StatusCode != 200 {
		q.Q("post status %v", resp.StatusCode)
		return fmt.Errorf("status code %v", resp.StatusCode)
	}
	return nil
}
