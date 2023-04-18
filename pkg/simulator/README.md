## Simulator

### support: 

- snmp  v2c

  Read community

  - public

  Write community

  - private

- gwd

### OS: 

- linux

  ​	Pre Install:

    - libpcap-dev
  ​		

- windows

  ​	Pre Install:

    - npcap

#### Note: Please use admin or root



## How to Import

1. ```sh
   go env -w GOPRIVATE=github.com/Atop-NMS-team/*
   ```

2. ```sh
   git config --global url."git@github.com:".insteadOf "https://github.com/"
   ```

3. ```sh
   go get github.com/Atop-NMS-team/simulator
   ```

   

## Basic Usage

```GO
import (
	"log"
	"os"
	"os/signal"
	"testing"

	"github.com/Atop-NMS-team/simulator/net"
	atopyaml "github.com/Atop-NMS-team/simulator/yaml"
)

func TestSimulatorFile(t *testing.T) {
    ethName, err := net.GetDefaultInterfaceName()
	if err != nil {
		t.Fatal(err)
	}
	simulators, err := atopyaml.NewSimulatorFile("./test.yaml", ethName)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range simulators {
		v.StartUp()
		defer v.Shutdown()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func TestSimulator(t *testing.T) {
    ethName, err := net.GetDefaultInterfaceName()
	if err != nil {
		t.Fatal(err)
	}
	simmap := map[string]atopyaml.Simulator{}
	simmap["1"] = atopyaml.Simulator{Number: 5, DeviceType: "EH7506", Start_prefixip: "192.168.6.1/24"}
	simulators, err := atopyaml.NewSimulator(simmap, ethName)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range simulators {
		v.StartUp()
		defer v.Shutdown()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

```









### Run

##### Note: if n =0 ,exit

```makefile
Usage:
  simulator run [flags]

Flags:
  -d, --debug            debug level
  -e, --ethName string   Network Interface Name (ip bind in Network Interface selected)
                         example:
                         vEthernet (WSL)
                         乙太網路
                         VirtualBox Host-Only Network
                         Wi-Fi (default "Wi-Fi")
  -h, --help             help for run
  -n, --number uint16    number of simulator (default 1)
  -v, --verb             verb level
  -y, --yaml string      path of yaml file,use yaml to decide simulator type and number
```



#### IP format

create new IP in network card

```sh
10.234.xx.xx/16
```





### RUN With yaml

```sh
 ./simulator.exe run -y ../config.yaml -d
```

```yaml
#exmaple array
#type: EH7506,EH7508,EH7512,EH7520,EHG750x

# environments:
  # simulator_group:                     #create one of device,useing fixed ip,fixed mac
      # type: "EH7520"
      # startMacAddress: "00-60-e9-18-01-02"  #if not exit ,create random mac automatically
      # startPreFixIp: "192.168.5.23/24"
  # simulator_group1:                    #create fiv of devices,useing fixed ip,random mac
      # number: 5
      # type: "EH7508"
      # macAddress: "00-60-e9-18-99-99"  #if not exit ,create random mac automatically
      # startPreFixIp: "192.168.6.23/24"
    
# environments:
  # simulator_group: #create once device using macAddress:"00-60-e9-18-01-99"
      # type: "EH7520"
      # startPreFixIp: "192.168.5.23/24"
      # startMacAddress: "00-60-e9-18-01-99"
  # simulator_group1: #create five of devices use random macAddress
      # number: 5
      # type: "EH7508"
      # startPreFixIp: "192.168.6.23/24"
  # simulator_group2: #create three of devices with macAddress increased from "00-60-e9-18-05-01"
      # nmuber: 3
      # type: "EH7508"
      # startPreFixIp: "192.168.7.23/24"
      # startMacAddress: "00-60-e9-18-05-01"
environments:
  simulator_group1: 
      number: 10
      type: "EH7506"
      startPreFixIp: "192.168.6.1/24"
      startMacAddress: "00-60-e9-18-01-01"
  simulator_group2: 
      number: 10
      type: "EH7508"
      startPreFixIp: "192.168.7.1/24"
      startMacAddress: "00-60-e9-18-0a-01"
```





## How to test  simulator

- ##### Step 1

  Enable simulator

  ```shell
  ./simulator.exe run -n 2 -d=true
  ```

- ##### Step 2

  Enable atopudpscan server

  ```shell
   ./atopudpscan run 
  ```
  
  
  
- ##### Step 3

  atopudpscan scan

  ```shell
  ./atopudpscan scan -r 10.234.0.0/16
  ```
  
  
  
  ##### response
  
  ```json
  2022/11/11 16:52:20 result[0]:map[ID: created_at: device_type: firmware_information:map[ap:ETH TestDevice_2 kernel:0.1] host_name:test2 ip_address:10.234.28.107 last_missing: last_recovered: last_seen: location:map[path:] mac_address:00-60-E9-AE-BB-42 model:ETH_Test2 more:map[gateway:10.0.0.254 isDHCP:false netmask:255.255.0.0] name: owner: support_protocols:[gwd]]
  2022/11/11 16:52:20 result[1]:map[ID: created_at: device_type: firmware_information:map[ap:ETH TestDevice_1 kernel:0.1] host_name:test1 ip_address:10.234.224.146 last_missing: last_recovered: last_seen: location:map[path:] mac_address:00-60-E9-F0-F5-84 model:ETH_Test1 more:map[gateway:10.0.0.254 isDHCP:false netmask:255.255.0.0] name: owner: support_protocols:[gwd]]
  
  ```
  
  
  
  


  - ##### Step 4

      reboot simulator

    ```shell
    $ ./atopudpscan reboot -d 10.234.28.107    -m 00-60-E9-AE-BB-42
    ```

    ##### Response

    ```json
    2022/11/11 16:52:56 session:{id:"1"  state:"success"  endedTime:"2022-11-11 16:52:56.7798112 +0800 CST m=+1207.221603801"}  device:{device_id:"00-60-E9-AE-BB-42"  device_path:"10.234.28.107"}  config_results:{protocol:"gwd"  kind:"reboot"}
    ```

    ##### simultor log

    ```shell
    time="2022-11-11T16:52:56+08:00" level=info msg="snmp:10.234.28.107 Shutdown"
    time="2022-11-11T16:52:56+08:00" level=info msg="device:10.234.28.107 Shutdown"
    time="2022-11-11T16:53:06+08:00" level=info msg="snmp:10.234.28.107 Run"
    time="2022-11-11T16:53:06+08:00" level=info msg="device:10.234.28.107 start up"
    
    ```




  - ##### Step 5

    make  simulator beep
    
    ```shell
    $ ./atopudpscan beep -d 10.234.28.107    -m 00-60-E9-AE-BB-42
    ```
    
    ##### respone
    
    ```json
    2022/11/11 16:54:11 session:{id:"1"  state:"success"  endedTime:"2022-11-11 16:54:11.856565 +0800 CST m=+1282.298357601"}  device:{device_id:"00-60-E9-AE-BB-42"  device_path:"10.234.28.107"}  config_results:{protocol:"gwd"  kind:"beep"}
    ```
    
    ##### simultor log
    
    ```shell
    time="2022-11-11T16:54:11+08:00" level=info msg="device:10.234.28.107 beep ....."
    time="2022-11-11T16:54:11+08:00" level=info msg="device:10.234.28.107 beep ....."
    
    ```
    
    


  - ##### Step 6

      configure simulator(set host name and new ip)

    ```shell
    $ ./atopudpscan config set  -d 10.234.28.107    -m 00-60-E9-AE-BB-42  -H abcd  -n 10.234.1.21
    ```

    

    ##### respone

    ```json
    2022/11/11 16:56:58 session:{id:"1"  state:"success"  endedTime:"2022-11-11 16:56:58.6658542 +0800 CST m=+1449.107646801"}  device:{device_id:"00-60-E9-AE-BB-42"  device_path:"10.234.28.107"}  config_results:{protocol:"gwd"  kind:"network"}
    
    ```

    ##### simultor log

    ```shell
    time="2022-11-11T16:56:58+08:00" level=info msg="SettingDevice ,device ip:10.234.28.107...."
    time="2022-11-11T16:56:58+08:00" level=info msg="device:10.234.1.21 Shutdown"
    time="2022-11-11T16:56:58+08:00" level=info msg="snmp:10.234.1.21 Shutdown"
    time="2022-11-11T16:57:08+08:00" level=info msg="snmp:10.234.1.21 Run"
    time="2022-11-11T16:57:08+08:00" level=info msg="device:10.234.1.21 start up"
    
    ```

    

    ##### verify result

    ```shell
    ./atopudpscan scan -r 10.234.0.0/16
    ```

    ##### response

    ```jsom
    2022/11/11 17:00:00 result[0]:map[ID: created_at: device_type: firmware_information:map[ap:ETH TestDevice_1 kernel:0.1] host_name:test1 ip_address:10.234.224.146 last_missing: last_recovered: last_seen: location:map[path:] mac_address:00-60-E9-F0-F5-84 model:ETH_Test1 more:map[gateway:10.0.0.254 isDHCP:false netmask:255.255.0.0] name: owner: support_protocols:[gwd]]
    2022/11/11 17:00:00 result[1]:map[ID: created_at: device_type: firmware_information:map[ap:ETH TestDevice_2 kernel:0.1] host_name:abcd ip_address:10.234.1.21 last_missing: last_recovered: last_seen: location:map[path:] mac_address:00-60-E9-AE-BB-42 model:ETH_Test2 more:map[gateway:10.0.0.254 isDHCP:false netmask:255.255.0.0] name: owner: support_protocols:[gwd]]
    
    ```

    

  - ##### Step 7

    Get device vaule 

    ```shell
    ./atopudpscan config get  -d 10.234.1.21 -m 00-60-E9-AE-BB-42
    ```

    ##### response

    ```json
    2022/11/11 17:03:20 device:{device_id:"00-60-E9-AE-BB-42"  device_path:"10.234.1.21"}  settings:{success:true  kind:"network"  configs:{fields:{key:"ap"  value:{string_value:"ETH TestDevice_2"}}  fields:{key:"gateway"  value:{string_value:"10.0.0.254"}}  fields:{key:"hostname"  value:{string_value:"abcd"}}  fields:{key:"iPAddress"  value:{string_value:"10.234.1.21"}}  fields:{key:"kernel"  value:{string_value:"0.1"}}  fields:{key:"macAddress"  value:{string_value:"00-60-E9-AE-BB-42"}}  fields:{key:"model"  value:{string_value:"ETH_Test2"}}  fields:{key:"netmask"  value:{string_value:"255.255.0.0"}}}}
    
    ```

    
    
  - ##### Step 8

    snmpscan simulator

    ```shell
    ./snmpscan scan -r 10.234.0.0/16
    ```

    ##### respone

    ```json
    [
      [
        {
          "value": "ETH_Test1",
          "name": ".1.3.6.1.2.1.1.1.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.1.0"
        },
        {
          "value": ".1.3.6.1.4.1.3755.0.0.21",
          "name": ".1.3.6.1.2.1.1.2.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.2.0"
        },
        {
          "value": 83101,
          "name": ".1.3.6.1.2.1.1.3.0",
          "kind": "int",
          "oid": ".1.3.6.1.2.1.1.3.0"
        },
        {
          "value": "www.atop.com.tw",
          "name": ".1.3.6.1.2.1.1.4.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.4.0"
        },
        {
          "value": "test1",
          "name": ".1.3.6.1.2.1.1.5.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.5.0"
        },
        {
          "value": "00:60:E9:F0:F5:84",
          "name": "MAC",
          "kind": "string",
          "oid": ""
        },
        {
          "value": "10.234.224.146",
          "name": "IP",
          "kind": "string",
          "oid": ""
        }
      ],
      [
        {
          "value": "ETH_Test2",
          "name": ".1.3.6.1.2.1.1.1.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.1.0"
        },
        {
          "value": ".1.3.6.1.4.1.3755.0.0.21",
          "name": ".1.3.6.1.2.1.1.2.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.2.0"
        },
        {
          "value": 83101,
          "name": ".1.3.6.1.2.1.1.3.0",
          "kind": "int",
          "oid": ".1.3.6.1.2.1.1.3.0"
        },
        {
          "value": "www.atop.com.tw",
          "name": ".1.3.6.1.2.1.1.4.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.4.0"
        },
        {
          "value": "abcd",
          "name": ".1.3.6.1.2.1.1.5.0",
          "kind": "string",
          "oid": ".1.3.6.1.2.1.1.5.0"
        },
        {
          "value": "00:60:E9:AE:BB:42",
          "name": "MAC",
          "kind": "string",
          "oid": ""
        },
        {
          "value": "10.234.1.21",
          "name": "IP",
          "kind": "string",
          "oid": ""
        }
      ]
    ]
    
    ```

## Writable oid

```sh
sntpClientStatus=.2.4.1.0 						Integer
sntpUTCTimezone=.2.4.3.0 						Integer
sntpServer1=.2.4.9.0 							OctetString
sntpServer2=.2.4.10.0 							OctetString
sntpServerQueryPeriod=.2.4.11.0					Integer
backupServerIP=.2.6.1.0 						OctetString
backupAgentBoardFwFileName=.2.6.2.0 			OctetString
backupStatus=.2.6.3.0 							Integer
restoreServerIP=.2.6.4.0 						OctetString
restoreAgentBoardFwFileName=.2.6.5.0 			OctetString
restoreStatus=.2.6.6.0 							Integer

agingTimeSetting=.2.11.1.0 						Integer
ptpState=.2.12.2.1.0 							Integer
ptpVersion=.2.12.2.2.0 							Integer
ptpSyncInterval=.2.12.2.3.0						Integer
ptpClockStratum=.2.12.2.5.0						Integer
ptpPriority1=.2.12.2.6.0						Integer
ptpPriority2=.2.12.2.7.0						Integer
rstpStatus=.4.2.1.0								Integer
qosCOSPriorityQueue=.6.4.1.3.1					Integer
qosTOSPriorityQueue=.6.6.1.3.1					Integer

trapServerTrapComm=8.6.1.3.0					OctetString
trapServerStatus=.8.6.1.5.0						Integer
trapServerPort=8.6.1.6.0						Integer
trapServerIP=8.6.1.7.0							IPAddress

syslogStatus=.10.1.2.1.0 						Integer
eventServerPort=.10.1.2.3.0 					Integer
eventServerLevel=.10.1.2.4.0 					Integer
eventLogToFlash=.10.1.2.5.0 					Integer
eventServerIP=.10.1.2.6.0 						IpAddress
systemModelName=.1.10.0							OctetString

eventPortEventEmail=.10.1.1.2.1.3.1				Integer
eventPortEventRelay=.10.1.1.2.1.4.1				Integer
eventPowerEventSMTP=.10.1.1.3.1.3.1				Integer
syslogEventsSMTP=.10.1.1.4.1.0					Integer
eventEmailAlertAddr=.10.1.3.2.0					OctetString
eventEmailAlertAuthentication=.10.1.3.3.0		Integer
eventEmailAlertAccount=.10.1.3.4.0				OctetString
lldpStatus=.12.1.0								Integer
```






## QA

#### Make error

```shell
  34 | #include <pcap.h>
      |          ^~~~~~~~
compilation terminated.
make: *** [Makefile:23: bin/simulator] Error 2
```

##### Solve

```shell
apt-get install -y libpcap-dev
```

#### wpcap.dll lose  in win

##### Solve

```sh
install npcap-1.71.exe  
https://npcap.com/
```

