package mnms

import (
	"fmt"
	"strings"
)

func HelpCmd(cmd string) string {
	if strings.HasPrefix(cmd, "help mtderase") {
		return `
  Erase target device mtd and restore default settings.

	Usage : mtderase [mac address] [ip address] [username] [password]
		[mac address] : target device mac address
		[ip address]  : target device ip address
		[username]    : target device login user name
		[password]    : target device login passwaord
	Example :
		mtderase AA-BB-CC-DD-EE-FF 10.0.50.1 admin default
		`
	}
	if strings.HasPrefix(cmd, "help beep") {
		return `
  Beep target device.

	Usage : beep [mac address] [ip address]
		[mac address] : target device mac address
		[ip address]  : target device ip address
	Example :
		beep AA-BB-CC-DD-EE-FF 10.0.50.1
		`
	}
	if strings.HasPrefix(cmd, "help reset") {
		return `
  Reset/Reboot target device.

	Usage : reset [mac address] [ip address] [username] [password]
		[mac address] : target device mac address
		[ip address]  : target device ip address
		[username]    : target device login user name
		[password]    : target device login passwaord
	Example :
		reset AA-BB-CC-DD-EE-FF 10.0.50.1 admin default
		`
	}
	if strings.HasPrefix(cmd, "help scan") {
		return `
  Use different protocol to scan all devices.

	Usage : scan [protocol]
		[protocol]    : use gwd/snmp to scan all devices.
	Example :
		scan gwd
		`
	}
	if strings.HasPrefix(cmd, "help config") {
		return `
  Configure device setting.

	Usage : config net [mac address] [current ip] [new ip] [mask] [gateway] [hostname]
		[mac address] : target device mac address
		[current ip]  : target device current ip address
		[new ip]      : target device would modify ip address
		[mask]        : target device network mask
		[gateway]     : target device gateway
		[hostname]    : target device host name
	Example :
		config net AA-BB-CC-DD-EE-FF 10.0.50.1 10.0.50.2 255.255.255.0 0.0.0.0 switch

	Usage : config syslog [mac address] [status] [server ip] [server port] [server level] [log to flash]
		[mac address] : target device mac address
		[status]      : use snmp to configure syslog enable/disable
		[server ip]   : use snmp to configure server ip address
		[server port] : use snmp to configure server port
		[server level]: use snmp to configure server log level
		[log to flash]: use snmp to configure log to flash
	Example :
		config syslog AA-BB-CC-DD-EE-FF 1 10.0.50.2 5514 1 1

	Usage : config beep [mac address] [ip address]
		[mac address] : target device mac address
		[ip address]  : target device ip address
	Example :
		config beep AA-BB-CC-DD-EE-FF 10.0.50.1

	Usage : config getsyslog [mac address]
		[mac address] : target device mac address
	Example :
		config getsyslog 00-60-E9-18-3C-3C

		
	Usage : config mtderase [mac address] [ip address] [username] [password]
		[mac address] : target device mac address
		[ip address]  : target device ip address
		[username]    : target device login user name
		[password]    : target device login passwaord
	Example :
		config mtderase AA-BB-CC-DD-EE-FF 10.0.50.1 admin default
		
	Usage : config snmp [enable]
		[enable]      : enable/disable
	Example :
		config snmp enable AA-BB-CC-DD-EE-FF admin default
		
	Usage : config local syslog path [path]
		[path]        : local syslog path
	Example :
		config local syslog path tmp/log

	Usage : config local syslog maxsize [maxsize]
		[maxsize]     : local syslog file maxsize size
	Example :
		config local syslog maxsize 100

	Usage : config local syslog compress [compress]
		[compress]     : would be compressed
	Example :
		config local syslog compress true

	Usage : config switch save [mac address] [username] [password]
		[mac address] : target device mac address
		[username]    : target device login user name
		[password]    : target device login passwaord
	Example :
		config switch save AA-BB-CC-DD-EE-FF admin default

	Usage : config local syslog read [start date] [start time] [end date] [end time] [max line]
		[start date]   : search syslog start date
		[start time]   : search syslog start time
		[end date]     : search syslog end date
		[end time]     : search syslog end time
		[max line]     : max lines, if without max line, that mean read all of lines
	Example :
		config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00
		config local syslog read 5
		config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00 5

		`
	}
	if strings.HasPrefix(cmd, "help switch") {
		return `
  Use target device CLI configuration commands.

	Usage : switch [mac address] [username] [password] [cli cmd...]
		[mac address] : target device mac address
		[username]    : target device login user name
		[password]    : target device login passwaord
		[cli cmd...]  : target device cli command
	Example :
		switch AA-BB-CC-DD-EE-FF admin default show ip
		`
	}

	if strings.HasPrefix(cmd, "help snmp") {
		return `
  Use snmp get/set/update/communities.

	Usage : snmp get [ip address] [oid]
		[ip address]  : target device ip address
		[oid]         : target oid
	Example :
		snmp get 10.0.50.1 1.3.6.1.2.1.1.1.0

	Usage : snmp set [ip address] [oid] [value] [value type]
		[ip address]  : target device ip address
		[oid]         : target oid
		[value]       : would be set value
		[value type]  : would be set value type.(OctetString, BitString, SnmpNullVar, Counter,
			                Counter64, Gauge, Opaque, Integer, ObjectIdentifier, IpAddress, TimeTicks)
	Example :
		snmp set 10.0.50.1 1.3.6.1.2.1.1.4.0 www.atop.com.tw OctetString
		
	Usage: snmp communities [user] [password] [mac]
 		Read device's SNMP communities and update to system.

		[user]     : Device telnet login user
		[password] : Device telnet login password
		[mac]      : Device mac address

 Example: 
 		snmp communities admin default 00-60-E9-27-E3-39

 Usage: snmp update community [mac] [read community] [write community]
 	Update device's SNMP communities manually.

	[mac]            : Device mac address
	[read community] : Device snmp read community
	[write community]: Device snmp write community

 Example: 
 		snmp update community 00-60-E9-27-E3-39 public private		
		
		
Usage: snmp options [port] [community] [version] [timeout]
 Update global snmp options.

	[port]     : snmp listen port
	[community]: snmp community
	[version]  : snmp version
	[timeout]  : snmp timeout

Example: 
		snmp options 161 public 2c 2	
		`
	}
	if strings.HasPrefix(cmd, "help log") {
		return `
  Configure log setting.

	Usage : log off
	Example :
		log off

	Usage : log pattern [pattern]
		[pattern]     : log pattern
	Example :
		log pattern .*

	Usage : log output [output]
		[output]      : log output
	Example :
		log output stderr

	Usage : log syslog [facility] [severity] [tag] [message]
		[facility]    : syslog facility
		[severity]    : syslog severity
		[tag]         : syslog was sent from tag what feature name
		[message]     : would send messages
	Example :
		log syslog 0 1 InsertDev "new device"
		`
	}
	if strings.HasPrefix(cmd, "help firmware") {
		return `
  Upgrade firmware.

	Usage : firmware [mac address] [file url]
		[mac address] : target device mac address
		[file url]    : file url
	Example :
		firmware AA-BB-CC-DD-EE-FF https://https://www.atoponline.com/.../EHG750X-K770A770.zip
		firmware AA-BB-CC-DD-EE-FF file:///C:/Users/testfile.txt
		`
	}
	if strings.HasPrefix(cmd, "help mqtt") {
		return `
  Use mqtt to publish/subscribe/unsubscribe/list topic.

	Usage : mqtt [mqttcmd] [topic] [data...]
		[mqttcmd]     : pub/sub/unsub/list
		                list is show all subscribe topic
		[tcp address] :	would pub/sub/unsub broker tcp address
		[topic]       : topic name
		[data...]     : data is messages, only publish use it.
	Example :
		mqtt pub 192.168.12.1:1883 topictest "this is messages."
		mqtt sub 192.168.12.1:1883 topictest
		mqtt unsub 192.168.12.1:1883 topictest
		mqtt list
		`
	}
	if strings.HasPrefix(cmd, "help opcua") {
		return `
  Opcua setting.
  
	Usage : opcua connect [url]	
		[url]         : connect to url
	Example :
		opcua connect opc.tcp://127.0.0.1:4840

	Usage : opcua read [node id]
		[node id]     : opcua node id
	Example :
		opcua read i=1002

	Usage : opcua browse [node id]
		[node id]     : opcua node id
	Example :
		opcua browse i=85

	Usage : opcua sub [node id]
		[node id]     : opcua node id
	Example :
		opcua sub i=1002

	Usage : opcua deletesub [sub id] [monitor id]
		[sub id]      : subscribe id
		[monitor id]  : monitored item id
	Example :
		opcua deletesub 1 1

	Usage : opcua close
	Example :
		opcua close
		`
	}
	if strings.HasPrefix(cmd, "help util") {
		return `
Utility command.

Usage: util -mnmspubkey -out [out_file]
	Get default mnms public key. If [out_file] is empty, output to stdout. 
	public key file use for decrypt mnms encrypted config file.

	[out_file]    : output public file name (optional)

Example :
	util -mnmspubkey -out mnms.pub
	util -mnmspubkey > mnms.pub


Usage: util -genrsa
	Generate rsa key pair. default output is 
	$HOME/.mnms/id_rsa and $HOME/.mnms/id_rsa.pub

Example :
	util -genrsa


Usage: util -genrsa -name [file_prefix]
	Generate rsa key pair to [file_prefix].pub and [file_prefix]
	[file_prefix]  : output file prefix (optional)

Example :
	util -genrsa -name ~/mnmskey

Usage: util -encrypt -pubkey [pubkey_file] -in [in_file] -out [out_file]
	Encrypt [in_file] with [pubkey_file] and output to [out_file]. 
	If [out_file] is empty, output to stdout. if [in_file] is empty, 
	input from stdin.

	[pubkey_file]  : public key file
	[in_file]      : input file (optional)
	[out_file]     : output file (optional)

Example :
	util -encrypt -pubkey mnms.pub -in mnms.conf -out mnms.conf.enc


Usage: util -decrypt -privkey [privkey_file] -in [in_file] -out [out_file]
	Decrypt [in_file] with [privkey_file] and output to [out_file]. 
	If [out_file] is empty, output to stdout. if [in_file] is empty, 
	input from stdin.

	[privkey_file] : private key file
	[in_file]      : input file (optional)
	[out_file]     : output file (optional)

Example :
	util -decrypt -privkey mnms.key -in mnms.conf.enc -out mnms.conf


Usage: util -export -configfile -pubkey [pubkey_file] -privkey [privkey_file] -out [out_file]
	Export config file to [out_file]. If [out_file] is empty, output to stdout. 
	if [privkey_file] is empty, use default private key file to decrypt config.json.

	[pubkey_file]  : public key file
	[privkey_file] : private key file (optional)
	[out_file]     : output file (optional)

Example :
	util -export -configfile -pubkey mnms.pub -out mnms.conf
	util -export -configfile -pubkey mnms.pub -privkey mnms.key -out mnms.conf


Usage: util -import -configfile -in [in_file]
	Import config file from [in_file]. If [in_file] is empty, 
	input from stdin. input file must be encrypted by pair public key.

	[in_file]      : input file (optional)
	
Example :
	util -import -configfile -in mnms.conf
	`
	}
	if strings.HasPrefix(cmd, "help") {
		msg := ""
		msg = fmt.Sprintf("  %s\n", "Usage:")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help mtderase", "Erase target device mtd and restore default settings.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help beep", "Beep target device.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help reset", "Reset/Reboot target device.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help scan", "Use different protocol to scan all devices.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help config", "Configure device setting.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help switch", "Use target device CLI configuration commands.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help snmp", "Use snmp get/set.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help log", "Configure log setting.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help firmware", "Upgrade firmware.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help mqtt", "Use mqtt to publish/subscribe topic.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help opcua", "Opcua setting.")
		msg = msg + fmt.Sprintf("\t%-15s %-15s\n", "help util", "Utilities commands.")
		return msg
	}
	return "error: invalid command"
}
