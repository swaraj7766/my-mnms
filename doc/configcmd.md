# MNMS config guide

## Command

1. ### net

   #### request

   ```sh
   config net mac currentIp newip mask gateway hostname
   ```

   example:

   ```sh
   config net 00-60-E9-18-3C-3C 192.168.4.29 192.168.4.29 255.255.255.0 192.168.4.254 test
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-21T16:44:12+08:00", Command:"config net 00-60-E9-18-3C-3C 192.168.4.29 192.168.4.29 255.255.255.0 192.168.4.254 test", Result:"", Status:"ok", Name:"", Retries:0}
   ```

   

2. ### beep

   #### request

   ```sh
   config beep mac  ip
   ```

   example:

   ```sh
   beep 00-60-E9-18-3C-3C 192.168.4.29
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-21T16:54:17+08:00", Command:"config beep 00-60-E9-18-3C-3C 192.168.4.29", Result:"", Status:"ok", Name:"", Retries:0
   ```

3. ### Snmp

   #### enable

   ```sh
   config snmp enable mac use pwd
   ```

   example:

   ```sh
   config snmp enable 00-60-E9-18-3C-3C admin default
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-21T17:00:14+08:00", Command:"config snmp enable 00-60-E9-18-3C-3C admin default", Result:" \r\ntest(config)# snmp\r\nblah blah\r\n\r\ntest(config)# \r\ntest(config)# blah blah\r\n%", Status:"ok", Name:"", Retries:0
   ```

   

   #### disable

   ```
   config snmp disbale mac use pwd
   ```

   example:

   ```
   config snmp disable 00-60-E9-18-3C-3C admin default
   ```

   #### response

   ```sh
   &mnms.CmdInfo{Timestamp:"2023-02-21T17:01:46+08:00", Command:"config snmp disable 00-60-E9-18-3C-3C admin default", Result:" \r\ntest(config)# no snmp\r\ntest(config)# \r\ntest(config)# blah blah\r\n%", Status:"ok", Name:"", Retries:0}
   ```

   

4. ### syslog path

   #### request

   ```sh
   config local syslog path param
   ```

   example:

   ```sh
   config local syslog path ./testlog
   ```

   #### response

   ```sh
   &mnms.CmdInfo{Timestamp:"2023-02-21T17:03:55+08:00", Command:"config syslogpath ./testlog", Result:"", Status:"ok", Name:"", Retries:0}
   ```

   

5. ### syslog max size

   #### request

   ```sh
   config local syslog maxsize param
   ```

   example:

   ```sh
   config local syslog maxsize 200
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-22T11:12:11+08:00", Command:"config local syslog maxsize 200", Result:"", Status:"ok", Name:"", Retries:0}
   ```

   

6. ### enable /disbale syslog compress 

   #### request

   ```sh
   config local syslog compress param
   ```

   example:

   ```sh
   config local syslog compress false
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-22T11:16:31+08:00", Command:"config local syslog compress true", Result:"", Status:"ok", Name:"", Retries:0}
   ```

   

7. ### syslog read

   #### request with time

   ```sh
   config local syslog read starttime endtime 
   ```

   without time

   ```sh
   config local syslog read 
   ```

   example:

   ```sh
   config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00
   ```

   with max line

   note: if without max line, that mean read all of lines

   ```sh
   config local syslog read 5
   ```

   with time and line

   ```sh
   config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00 5
   ```

   #### response

   ```sh
   mnms.CmdInfo{Timestamp:"2023-02-22T14:06:21+08:00", Command:"config local syslog read 2023/02/21 22:06:00 2023/02/22 22:08:00", Result:"[{\"Facility\":0,\"Severity\":0,\"Priority\":0,\"Timestamp\":\"2023-02-21T22:17:12Z\",\"Hostname\":\"172.18.112.1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"123456\"},{\"Facility\":0,\"Severity\":0,\"Priority\":0,\"Timestamp\":\"2023-02-21T22:17:13Z\",\"Hostname\":\"172.18.112.1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"123456\"},{\"Facility\":0,\"Severity\":0,\"Priority\":0,\"Timestamp\":\"2023-02-21T22:17:16Z\",\"Hostname\":\"172.18.112.1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"99\"},{\"Facility\":0,\"Severity\":0,\"Priority\":0,\"Timestamp\":\"2023-02-21T22:17:20Z\",\"Hostname\":\"172.18.112.1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"1111\"},{\"Facility\":5,\"Severity\":7,\"Priority\":47,\"Timestamp\":\"2023-02-22T10:23:12Z\",\"Hostname\":\"LAPTOP-ERS90EE1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"test\"},{\"Facility\":5,\"Severity\":7,\"Priority\":47,\"Timestamp\":\"2023-02-22T10:23:24Z\",\"Hostname\":\"LAPTOP-ERS90EE1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"test\"},{\"Facility\":0,\"Severity\":6,\"Priority\":6,\"Timestamp\":\"2023-02-22T10:29:45+08:00\",\"Hostname\":\"local\",\"Appname\":\"d:\\\\NMS\\\\mnms\\\\issue169\\\\mnms\\\\__debug_bin.exe\",\"ProcID\":\"96452\",\"MsgID\":\"RFC5424Formatter\",\"Message\":\"hekko\"},{\"Facility\":0,\"Severity\":0,\"Priority\":0,\"Timestamp\":\"2023-02-22T11:09:15Z\",\"Hostname\":\"172.18.112.1\",\"Appname\":null,\"ProcID\":null,\"MsgID\":null,\"Message\":\"1111\"}]", Status:"ok", Name:"", Retries:0}
   ```

