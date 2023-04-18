package main

import (
	"bytes"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"mnms"
	_ "net/http/pprof"

	"github.com/qeof/q"
)

var Version string

func main() {
	var wg sync.WaitGroup
	q.O = "stderr"
	q.P = ""

	stop := func() {
		err := mnms.SendSyslog(mnms.LOG_NOTICE, "main", "exiting main() "+mnms.QC.Name)
		if err != nil {
			q.Q(err)
		}
		if err := recover(); err != nil {
			q.Q("recover: exit with exception", err)
			fmt.Fprintf(os.Stderr, "recover: exit with exception: %v\n", err)
			mnms.QC.DumpStackTrace = true
			mnms.DoExit(1)
		}
	}

	flagversion := flag.Bool("version", false, "print version")
	flag.IntVar(&mnms.QC.Port, "p", mnms.QC.Port, "port")
	flag.StringVar(&mnms.QC.RootURL, "r", "", "root URL")
	flag.StringVar(&mnms.QC.Name, "n", "", "name")
	flag.StringVar(&q.O, "O", "stderr", "debug log output")
	dp := flag.String("P", "", "debug log pattern string")
	flag.BoolVar(&mnms.QC.DumpStackTrace, "ds", false, "dump stack trace when exiting with non zero code")
	flag.IntVar(&mnms.QC.CmdInterval, "ic", mnms.QC.CmdInterval, "command processing interval")
	flag.IntVar(&mnms.QC.RegisterInterval, "ir", mnms.QC.RegisterInterval, "client node registration interval")
	flag.IntVar(&mnms.QC.GwdInterval, "ig", mnms.QC.GwdInterval, "device scan interval")
	flag.StringVar(&mnms.QC.Domain, "d", "", "domain")
	svc := flag.Bool("s", false, "run as a service")
	flag.BoolVar(&mnms.QC.IsRoot, "R", false, "run as root")
	nosyslog := flag.Bool("nosyslog", false, "no syslog service")
	notrap := flag.Bool("notrap", false, "no snmp trap service")
	nohttp := flag.Bool("nohttp", false, "no http service")
	fake := flag.Bool("fake", false, "fake mode for testing")
	_ = flag.Bool("M", false, "monitor mode")
	nomqttbroker := flag.Bool("nomqbr", false, "no mqtt broker")
	flag.StringVar(&mnms.QC.MqttBrokerAddr, "mb",
		mnms.QC.MqttBrokerAddr, "mqtt broker address")
	flag.StringVar(&mnms.QC.SyslogServerAddr, "ss",
		mnms.QC.SyslogServerAddr, "syslog server address")
	flag.StringVar(&mnms.QC.TrapServerAddr, "ts",
		mnms.QC.TrapServerAddr, "trap server address")
	flag.StringVar(&mnms.QC.RemoteSyslogServerAddr, "rs",
		mnms.QC.RemoteSyslogServerAddr, "remote syslog server address")
	flag.StringVar(&mnms.QC.SyslogLocalPath, "so", mnms.QC.SyslogLocalPath, "local path of syslog")
	flag.UintVar(&mnms.QC.SyslogFileSize, "sf", mnms.QC.SyslogFileSize, "file size(megabytes) of syslog")
	flag.BoolVar(&mnms.QC.SyslogCompress, "sc", mnms.QC.SyslogCompress, "enable compress file of backup syslog")
	prikeyfile := flag.String("privkey", "", "private key file")
	cmdflagnoow := flag.Bool("cno", false, "command overwrite flag")
	cmdflagall := flag.Bool("ca", false, "command all flag")
	cmdflagnosys := flag.Bool("cns", false, "command syslog flag")
	cmdClient := flag.String("cc", "", "command client specification")
	cmdTag := flag.String("ct", "", "command tag")
	pp := flag.Bool("pprof", false, "enable pprof analysis")
	var daemon string
	flag.StringVar(&daemon, mnms.DaemonFlag, "", mnms.Usage)
	flag.Parse()
	service := func() {
		if *flagversion {
			info, _ := debug.ReadBuildInfo()
			fmt.Fprintln(os.Stderr, Version)
			fmt.Fprintln(os.Stderr, info.GoVersion, info.Main, info.Settings)
			mnms.DoExit(0)
		}
		if *pp {
			go func() {
				q.Q(http.ListenAndServe("localhost:6060", nil))
			}()

		}

		if *dp == "." {
			fmt.Fprintln(os.Stderr, "error: invalid debug pattern")
			mnms.DoExit(1)
		}
		_, err := regexp.Compile(*dp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: invalid regular expression, %v\n", err)
			mnms.DoExit(1)
		}
		q.P = *dp
		q.Q(q.O, q.P)
		args := flag.Args()
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		q.Q(exPath)
		if len(os.Args) > 2 && os.Args[1] == "-M" {
			q.P = ".*"
			q.Q("monitor run mode")
			t0 := time.Now().Unix()
			ix := 0
			for {
				ix++
				runarg := fmt.Sprintf("monitor: run #%d %v", ix, os.Args)
				q.Q("monitor: run", ix, os.Args)
				err = mnms.SendSyslog(mnms.LOG_NOTICE, "monitor", runarg)
				if err != nil {
					q.Q("error: syslog", err)
				}
				ec := exec.Command(os.Args[0], os.Args[2:]...)
				output, err := ec.CombinedOutput()
				t1 := time.Now().Unix()
				diff := t1 - t0
				q.Q("monitor:", string(output))
				if diff < 3 { // XXX
					q.Q("monitor: spinning, exit")
					mnms.DoExit(1)
				}
				t0 = t1
				if err != nil {
					q.Q("monitor:", err)
					errmsg := fmt.Sprintf("monitor: #%d %v",
						ix, err.Error())
					err = mnms.SendSyslog(mnms.LOG_ERR, "monitor", errmsg)
					if err != nil {
						q.Q("error: syslog", err)
					}
					continue
				}
			}
		}
		// check args length > 1 and args[1] is a 'util'
		if len(os.Args) > 2 && os.Args[1] == "util" && os.Args[2] != "help" {
			err := mnms.ProcessDirectCommands()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			mnms.DoExit(0)
			return
		}
		mnms.QC.AdminToken, err = mnms.GetToken("admin")
		if err != nil {
			q.Q(err)
			fmt.Fprintln(os.Stderr, "error: can't get admin token")
			mnms.DoExit(1)
		}
		if mnms.QC.RemoteSyslogServerAddr == "" {
			q.Q("warning: missing remote syslog server address")
		}
		if *svc || mnms.QC.IsRoot {
			if !*nosyslog {
				wg.Add(1)
				go func() {
					defer wg.Done()
					mnms.StartSyslogServer()
				}()
			} else {
				q.Q("skip running syslog server")
			}
		}
		localRootURL := fmt.Sprintf("http://localhost:%d", mnms.QC.Port)

		if !*svc && !mnms.QC.IsRoot {
			q.Q("cli", args)
			CheckArgs(args)
			// implement cli by posting commands via http api
			acmd := args[0]
			found := false
			for _, c := range mnms.ValidCommands {
				if c == acmd {
					found = true
				}
			}

			if !found {
				fmt.Fprintf(os.Stderr, "error: invalid cmd %s\n\n", acmd)
				helpMsg := mnms.HelpCmd("help")
				fmt.Fprintf(os.Stderr, "%s\n", helpMsg)
				mnms.DoExit(1)
			}
			if len(args) < 2 {
				fmt.Fprintf(os.Stderr, "error: insufficient args\n")
				helpMsg := mnms.HelpCmd("help " + acmd)
				fmt.Fprintf(os.Stderr, "%s\n", helpMsg)
				mnms.DoExit(1)
			}
			if args[1] == "help" {
				for _, c := range mnms.ValidCommands {
					if c == acmd {
						fmt.Println(acmd)
						helpMsg := mnms.HelpCmd("help " + acmd)
						fmt.Fprintf(os.Stderr, "%s\n", helpMsg)
						mnms.DoExit(0)
					}
				}
			}

			cmd := strings.Join(args, " ")
			cmdinfo := make(map[string]mnms.CmdInfo)
			kcmd := cmd
			if *cmdClient != "" {
				kcmd = "@" + *cmdClient + " " + cmd
			}
			ci := mnms.CmdInfo{
				Timestamp:   time.Now().Format(time.RFC3339),
				Command:     cmd,
				NoOverwrite: *cmdflagnoow,
				All:         *cmdflagall,
				NoSyslog:    *cmdflagnosys,
				Kind:        "usercommand",
				Client:      *cmdClient,
				Tag:         *cmdTag,
			}
			cmdinfo[kcmd] = ci
			jsonBytes, err := json.Marshal(cmdinfo)
			if err != nil {
				q.Q(err)
				mnms.DoExit(1)
			}
			q.Q("posting cmd", cmdinfo)
			url := fmt.Sprintf(localRootURL + "/api/v1/commands")
			resp, err := mnms.PostWithToken(url, mnms.QC.AdminToken, bytes.NewBuffer(jsonBytes))
			if err != nil {
				q.Q(err.Error())
				fmt.Fprintf(os.Stderr, "error: cannot connect to root server at %v\n", localRootURL)
				mnms.DoExit(1)
			}
			if resp == nil {
				fmt.Fprintf(os.Stderr, "error: no response from root server at %v\n", localRootURL)
				mnms.DoExit(1)
			}

			// save close, resp should be nil here
			defer resp.Body.Close()
			result, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				q.Q(err.Error())
				fmt.Fprintf(os.Stderr, "error: reading response from root server at %v\n", localRootURL)
				mnms.DoExit(1)
			}
			q.Q(string(result))
			mnms.DoExit(0)
		}

		if mnms.QC.Name == "" {
			fmt.Fprintln(os.Stderr, "error: -n name is required")
			mnms.DoExit(1)
		}
		q.Q(mnms.QC.Name)
		q.Q(mnms.QC.Domain, len(mnms.QC.Domain))

		if !*fake {
			wg.Add(1)
			if !*nohttp {
				go func() {
					defer wg.Done()
					mnms.HTTPMain()
				}()
			}
			if mnms.QC.IsRoot {
				wg.Add(1)
				// init publickey
				// setup private key if provided
				if *prikeyfile != "" {
					keyBytes, err := os.ReadFile(*prikeyfile)
					if err == nil {
						// check keyBytes is a private key pem
						block, _ := pem.Decode(keyBytes)
						if block != nil {
							q.Q("use private key file", *prikeyfile)
							mnms.SetPrivateKey(string(keyBytes))
						}
					} else {
						q.Q(err)
						q.Q("use default private key")
					}
				}
				mnms.QC.OwnPublicKeys, err = mnms.GenerateOwnPublickey()
				if err != nil {
					q.Q(err)
					mnms.DoExit(1)
				}
				// init mnms config
				err := mnms.InitDefaultMNMSConfigIfNotExist()
				if err != nil {
					q.Q(err)
					mnms.DoExit(1)
				}
			}
		}

		time.Sleep(1 * time.Second)

		if *svc {
			if mnms.QC.RootURL != "" {
				wg.Add(1)
				q.Q(mnms.QC.RegisterInterval)
				go func() {
					defer wg.Done()
					mnms.RegisterMain()
				}()
			}

			if !*fake {
				wg.Add(1)
				go func() {
					defer wg.Done()
					mnms.GwdMain()
					q.Q("GwdMain returned")
				}()
			} else {
				err := mnms.SetupFakeClient()
				if err != nil {
					q.Q("failed to setup fake client", err)
					mnms.DoExit(1)
				}
			}

			if !*notrap {
				wg.Add(1)
				go func() {
					defer wg.Done()
					mnms.StartTrapServer()
				}()
			} else {
				q.Q("skip running trap server")
			}
			if !*nomqttbroker {
				wg.Add(1)
				go func() {
					defer wg.Done()
					err := mnms.RunMqttBroker(mnms.QC.Name)
					if err != nil {
						q.Q(err)
					}
				}()
			} else {
				q.Q("skip running mqtt broker")
			}
			if mnms.QC.RootURL != "" {
				wg.Add(1)
				go func() {
					defer wg.Done()
					mnms.TopologyPollingWithTimer(10) // XXX
				}()
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				q.Q(mnms.QC.CmdInterval)
				for {
					time.Sleep(time.Duration(mnms.QC.CmdInterval) * time.Second) // XXX
					err := mnms.CheckCmds()
					if err != nil {
						q.Q(err)
					}
				}
			}()
		}
		wg.Wait()
		q.Q("exit normally")
		mnms.DoExit(0)
	}

	//enable Daemon
	s, err := mnms.NewDaemon(mnms.QC.Name, os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		mnms.DoExit(0)
	}
	s.RegisterRunEvent(service)
	s.RegisterStopEvent(stop)
	err = s.RunMode(daemon)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		mnms.DoExit(0)
	}
}
