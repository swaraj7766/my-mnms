# MNMS Opcua guide

## Command

1. ### Connect

   #### request

   ```sh
   opcua connect url
   ```

   example:

   ```sh
   opcua connect opc.tcp://127.0.0.1:4840
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-10T14:27:08+08:00", Command:"opcua connect opc.tcp://127.0.0.1:4840", Result:"", Status:"ok", Name:"", Retries:0}
   ```

   

2. ### Read

   #### request

   ```sh
   opcua read nodid
   ```

   example:

   ```sh
   opcua read i=1002
   ```

   #### response

   ```sh
   &mnms.CmdInfo{Timestamp:"2023-02-10T14:27:43+08:00", Command:"opcua read i=1002", Result:"{\"SpecifiedAttributes\":0,\"DisplayName\":\"Temperature\",\"Description\":\"Temperature\",\"WriteMask\":0,\"UserWriteMask\":0,\"Value\":0.5,\"DataType\":\"i=10\",\"ValueRank\":-1,\"ArrayDimensions\":[],\"AccessLevel\":3,\"UserAccessLevel\":3,\"MinimumSamplingInterval\":0,\"Historizing\":false}", Status:"ok", Name:"", Retries:0}
   ```

3. ### Browse

   #### request

   ```sh
   opcua read nodid
   ```

   example:

   ```sh
   opcua browse i=85
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-10T14:28:06+08:00", Command:"opcua browse i=85", Result:"[{\"ReferenceTypeID\":\"i=35\",\"IsForward\":true,\"NodeID\":{\"ServerIndex\":0,\"NamespaceURI\":\"\",\"NodeID\":\"i=2253\"},\"BrowseName\":\"0:Server\",\"DisplayName\":\"Server\",\"NodeClass\":1,\"TypeDefinition\":{\"ServerIndex\":0,\"NamespaceURI\":\"\",\"NodeID\":\"i=2004\"}},{\"ReferenceTypeID\":\"i=35\",\"IsForward\":true,\"NodeID\":{\"ServerIndex\":0,\"NamespaceURI\":\"\",\"NodeID\":\"i=1001\"},\"BrowseName\":\"0:Boiler\",\"DisplayName\":\"Boiler\",\"NodeClass\":1,\"TypeDefinition\":{\"ServerIndex\":0,\"NamespaceURI\":\"\",\"NodeID\":\"i=58\"}}]", Status:"ok", Name:"", Retries:0}
   ```

4. ### Subscription

   #### request

   ```
   opcua read nodid
   ```

   example:

   ```sh
   opcua sub i=1002
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-10T14:29:34+08:00", Command:"opcua sub i=1002", Result:"SubscriptionID:1,MonitoredItemID:1", Status:"ok", Name:"", Retries:0}
   ```

5. ### Delete Subscription

   #### request

   ```sh
   opcua deletesub SubscriptionID MonitoredItemID
   ```

   example:

   ```sh
   opcua deletesub 1 1
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-10T14:32:49+08:00", Command:"opcua deletesub 1 1", Result:"", Status:"ok", Name:"", Retries:0}
   ```

6. ### Close

   #### request

   ```sh
   opcua close
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-10T14:34:45+08:00", Command:"opcua close", Result:"", Status:"ok", Name:"", Retries:0}
   ```

   

