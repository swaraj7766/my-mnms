package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"mnms"

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/qeof/q"
)

var Usage = func() {
	fmt.Printf("\n")
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func CheckArgs(args []string) {
	if len(args) == 0 {
		q.Q("no args exist")
		Usage()
		mnms.DoExit(1)
	}
}

// direct commands
type directCommands struct {
	GenerateRSAKeyPair bool   // generate rsa key pair
	Decrypt            bool   // decrypt something
	Export             bool   // export mnms config
	Import             bool   // import mnms config
	PublickeyPath      string // public key path
	PrivatekeyPath     string // private key path
	AdminPassword      bool   // admin password
	ConfigFile         bool   //config file
	In                 string // input file
	Out                string // output file
}

// ProcessDirectCommands checks if there is any direct commands
func ProcessDirectCommands() {
	checkerr := func(err error) {
		if err != nil {
			q.Q(err)
			mnms.DoExit(1)
		}
	}
	q.Q(os.Args)
	// check args length > 1 and args[1] is a 'cmd'
	if len(os.Args) < 3 || os.Args[1] != "cmd" {
		return
	}

	dc := directCommands{
		//default values
		GenerateRSAKeyPair: false,
		Decrypt:            false,
		Export:             false,
		Import:             false,
	}

	subsetFlags := flag.NewFlagSet("subset", flag.ExitOnError)
	subsetFlags.BoolVar(&dc.GenerateRSAKeyPair, "genrsa", false, "generate rsa key pair")
	subsetFlags.BoolVar(&dc.Decrypt, "decrypt", false, "decrypt something")
	subsetFlags.BoolVar(&dc.Export, "export", false, "export mnms config")
	subsetFlags.BoolVar(&dc.Import, "import", false, "import mnms config")
	subsetFlags.BoolVar(&dc.AdminPassword, "adminpass", false, "admin password")
	subsetFlags.BoolVar(&dc.ConfigFile, "configfile", false, "config file")
	subsetFlags.StringVar(&dc.PublickeyPath, "pubkey", "", "public key path")
	subsetFlags.StringVar(&dc.PrivatekeyPath, "privkey", "", "private key path")
	subsetFlags.StringVar(&dc.In, "in", "", "input file")
	subsetFlags.StringVar(&dc.Out, "out", "", "output file")
	help := subsetFlags.Bool("help", false, "help")
	subsetFlags.Parse(os.Args[2:])

	// dump usage
	if *help {
		subsetFlags.Usage()
		mnms.DoExit(0)
	}
	q.Q(dc)

	mnmsFolder, err := mnms.CheckMNMSFolder()
	checkerr(err)
	// -genrsa -pubkey {public_key_file} -privkey {private_key_file}
	if dc.GenerateRSAKeyPair {
		q.Q("generating rsa key pair")
		if dc.PrivatekeyPath == "" {
			q.Q("By default, private key will write to $HOME/.mnms/privkey.pem or specify with -privkey flag")
			dc.PrivatekeyPath = path.Join(mnmsFolder, "privkey.pem")
		}
		if dc.PublickeyPath == "" {
			q.Q("By default, public key will write to $HOME/.mnms/pubkey.pem or specify with -pubkey flag")
			dc.PublickeyPath = path.Join(mnmsFolder, "pubkey.pem")
		}

		// generate rsa key pair
		prikey, err := mnms.GenerateRSAKeyPair(4096)
		checkerr(err)
		// generate private and public key bytes
		prikeyBytes := mnms.EndcodePrivateKeyToPEM(prikey)
		// write private key to prikeyPath
		err = ioutil.WriteFile(dc.PrivatekeyPath, prikeyBytes, 0644)
		checkerr(err)

		pubkeyBytes, err := mnms.GenerateRSAPublickey(prikey)
		checkerr(err)

		// write public key to pubkeyPath
		err = ioutil.WriteFile(dc.PublickeyPath, pubkeyBytes, 0644)
		checkerr(err)
		q.Q("done")
		mnms.DoExit(0)
	}

	// mnmsctl cmd -decrypt -in {cipher_file} -out {plain_file} -privkey {private_key_file}
	if dc.Decrypt {
		if dc.In == "" {
			q.Q("input file is required")
			mnms.DoExit(1)
		}
		if dc.Out == "" {
			q.Q("By default output file is saved to $HOME/.mnms/output or specify with -out flag")
			dc.Out = path.Join(mnmsFolder, "output")
		}
		if dc.PrivatekeyPath == "" {
			q.Q("private key is required")
			mnms.DoExit(1)
		}
		// read private key
		prikeyBytes, err := ioutil.ReadFile(dc.PrivatekeyPath)
		checkerr(err)
		// read encrypted pass
		cipherBytes, err := ioutil.ReadFile(dc.In)
		checkerr(err)
		// decrypt pass
		decryptedPass, err := mnms.DecryptWithPrivateKeyPEM(cipherBytes, prikeyBytes)
		checkerr(err)
		// write decrypted pass to output file
		err = ioutil.WriteFile(dc.Out, decryptedPass, 0644)
		checkerr(err)
		q.Q("done")
		mnms.DoExit(0)
	}

	// mnmsctl cmd -export -adminpass -out {output_file} -pubkey {public_key_file}
	// mnmsctl cmd -export -configfile -out {output_file} -pubkey {public_key_file}
	// mnmsctl cmd -import -configfile -in {config_file}
	if dc.Export {
		if dc.Out == "" {
			q.Q("By default output file is saved to $HOME/.mnms/output or specify with -out flag")
			dc.Out = path.Join(mnmsFolder, "output")
		}
		c, err := mnms.GetMNMSConfig()
		checkerr(err)
		if dc.PublickeyPath == "" {
			q.Q("public key is required")
			mnms.DoExit(1)
		}
		// read public key
		pubkeyBytes, err := ioutil.ReadFile(dc.PublickeyPath)
		checkerr(err)

		// export admin password
		if dc.AdminPassword {
			var pass string
			// get admin password
			for _, u := range c.Users {
				if u.Name == "admin" {
					pass = u.Password
					break
				}
			}
			if pass == "" {
				q.Q("error: admin password not found")
				mnms.DoExit(1)
			}

			// encrypt pass with public key
			encryptedPass, err := mnms.EncryptWithPublicKey([]byte(pass), pubkeyBytes)
			checkerr(err)
			// write encrypted pass to passfile
			err = ioutil.WriteFile(dc.Out, encryptedPass, 0644)
			checkerr(err)

			q.Q("done")
			mnms.DoExit(0)
		}

		// export mnms config
		if dc.ConfigFile {
			// export mnms config
			configJSON, err := json.Marshal(c)
			checkerr(err)
			// encrypt with public key
			encryptedConfig, err := mnms.EncryptWithPublicKey(configJSON, pubkeyBytes)
			checkerr(err)
			// write encrypted config to output file

			err = ioutil.WriteFile(dc.Out, encryptedConfig, 0644)
			checkerr(err)
			q.Q("done")
			mnms.DoExit(0)
		}
	}

	// mnmsctl cmd -import -configfile -in {config_file}
	if dc.Import {
		if dc.In == "" {
			q.Q("input file is required")
			mnms.DoExit(1)
		}
		if dc.ConfigFile {
			// read config file
			configBytes, err := ioutil.ReadFile(dc.In)
			checkerr(err)
			c := mnms.MNMSConfig{}
			err = json.Unmarshal(configBytes, &c)
			checkerr(err)
			// save config
			err = mnms.WriteMNMSConfig(&c)
			checkerr(err)
			q.Q("done")
			mnms.DoExit(0)

		}
		q.Q("done")
		mnms.DoExit(0)

	}

}

func main() {
	var wg sync.WaitGroup
	q.O = "stderr"
	q.P = ".*"

	defer func() {
		if err := recover(); err != nil {
			q.Q("exit with exception", err)
			fmt.Fprintf(os.Stderr, "exit with exception: %v\n", err)
			mnms.DoExit(1)
		}
	}()
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	q.Q(exPath)

	ProcessDirectCommands()

	flag.IntVar(&mnms.QC.Port, "p", mnms.QC.Port, "port")
	flag.StringVar(&mnms.QC.Root, "r", "", "root URL")
	flag.StringVar(&mnms.QC.Name, "n", "", "name")
	flag.StringVar(&q.O, "O", "stderr", "debug log output")
	flag.StringVar(&q.P, "P", ".*", "debug log pattern")
	flag.IntVar(&mnms.QC.ArpInterval, "i", mnms.QC.ArpInterval, "ping devices interval")
	flag.IntVar(&mnms.QC.CmdInterval, "ci", mnms.QC.CmdInterval, "command processing interval")
	flag.StringVar(&mnms.QC.Domain, "d", "", "domain")
	svc := flag.Bool("s", false, "run as a service")
	flag.BoolVar(&mnms.QC.IsRoot, "R", false, "run as root")
	nosyslog := flag.Bool("nosyslog", false, "no syslog service")
	notrap := flag.Bool("notrap", false, "no snmp trap service")
	fake := flag.Bool("fake", false, "fake mode for testing")
	nomqttbroker := flag.Bool("nomqbr", false, "no mqtt broker")
	flag.StringVar(&mnms.QC.MqttBrokerAddr, "mb",
		mnms.QC.MqttBrokerAddr, "mqtt broker address")
	flag.StringVar(&mnms.QC.SyslogServerAddr, "ss",
		mnms.QC.SyslogServerAddr, "syslog server address")
	flag.StringVar(&mnms.QC.TrapServerAddr, "ts",
		mnms.QC.TrapServerAddr, "trap server address")
	flag.StringVar(&mnms.QC.RemoteSyslogServerAddr, "rs",
		mnms.QC.RemoteSyslogServerAddr, "remote syslog server address")

	loaddevs := flag.Bool("loaddev", false, "load devices from the last saved file")
	flag.StringVar(&mnms.QC.SyslogLocalPath, "sp", mnms.QC.SyslogLocalPath, "local path of syslog")
	flag.UintVar(&mnms.QC.SyslogFileSize, "sf", mnms.QC.SyslogFileSize, "file size(megabytes) of syslog")
	flag.BoolVar(&mnms.QC.SyslogCompress, "sc", mnms.QC.SyslogCompress, "enable compress file of backup syslog")

	flag.Parse()

	args := flag.Args()
	q.Q(args)
	if mnms.QC.Name == "" {
		rand.Seed(time.Now().UnixNano())
		mnms.QC.Name = petname.Generate(3, "")
		q.Q("No name specified, using the random word")
	}
	q.Q(mnms.QC.Name)
	q.Q(mnms.QC.Domain, len(mnms.QC.Domain))
	mnms.QC.AdminToken, err = mnms.GetToken("admin")
	if err != nil {
		q.Q(err)
		mnms.DoExit(1)
	}

	if mnms.QC.RemoteSyslogServerAddr == "" {
		q.Q("Missing remote syslog server address")
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

	if !*svc && !mnms.QC.IsRoot {
		CheckArgs(args)
		// implement cli by posting commands via http api
		cmd := strings.Join(args, " ")
		cmdinfo := make(map[string]mnms.CmdInfo)
		ci := mnms.CmdInfo{
			Timestamp: time.Now().Format(time.RFC3339),
			Command:   cmd,
		}
		cmdinfo[cmd] = ci
		jsonBytes, err := json.Marshal(cmdinfo)
		if err != nil {
			q.Q(err)
			mnms.DoExit(1)
		}

		url := fmt.Sprintf("http://localhost:%d/api/v1/commands", mnms.QC.Port)
		resp, err := mnms.PostWithToken(url, mnms.QC.AdminToken, bytes.NewBuffer(jsonBytes))
		if err != nil {
			q.Q(err.Error())
			mnms.DoExit(1)
		}
		defer resp.Body.Close()
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			q.Q(err.Error())
			mnms.DoExit(1)
		}
		q.Q(string(result))
		return
	}

	if !*fake {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mnms.HTTPMain()
			// TODO root to dump snapshots of devices, logs, commands
		}()

		if mnms.QC.IsRoot {
			wg.Add(1)
			// init publickey
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

			q.Q("start validate cmds loop")
			go func() {
				defer wg.Done()
				for {
					time.Sleep(time.Duration(mnms.QC.CmdInterval) * time.Second) // XXX
					err := mnms.ValidateCommands()
					if err != nil {
						q.Q(err)
					}
				}
			}()

		}
	}

	time.Sleep(1 * time.Second)

	if *loaddevs {
		err := mnms.LoadDevices()
		if err != nil {
			q.Q(err)
		}
	}

	if *svc {
		if mnms.QC.Root != "" {

			_, err := mnms.PostWithToken(mnms.QC.Root+"/api/v1/register", mnms.QC.AdminToken, bytes.NewBuffer([]byte(mnms.QC.Name)))
			if err != nil {
				q.Q(err)
				mnms.DoExit(1)
			}
			q.Q("name registered to root", mnms.QC.Name)
		}

		if !*fake {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mnms.GwdMain()
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
		if mnms.QC.Root != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mnms.CheckAllDevicesAlive()
			}()
		}

		if mnms.QC.Root != "" {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mnms.TopologyPollingWithTimer(10)
			}()
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			q.Q(mnms.QC.CmdInterval)
			for {
				time.Sleep(time.Duration(mnms.QC.CmdInterval) * time.Second) // XXX
				err := mnms.CheckCommands()
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
