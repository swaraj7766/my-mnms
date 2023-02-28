package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qeof/q"
)

type DevInfo struct {
	Mac       string `json:"mac"`
	ModelName string `json:"modelname"`
	Timestamp string `json:"timestamp"`
	Scanproto string `json:"scanproto"`
	IPAddress string `json:"ipaddress"`
	Netmask   string `json:"netmask"`
	Gateway   string `json:"gateway"`
	Hostname  string `json:"hostname"`
	Kernel    string `json:"kernel"`
	Ap        string `json:"ap"`
	ScannedBy string `json:"scannedby"`
	ArpMissed int    `json:"arpmissed"`
	UnixTime  string `json:"unixtime"`
}

var specialMac = "11-22-33-44-55-66"

var lastUnixTime string

func init() {
	QC.DevData = make(map[string]DevInfo)
	lastUnixTime = strconv.FormatInt(time.Now().Unix(), 10)
}

func InsertModel(model GwdModelInfo, proto string) {
	//discovered device model will be entered into device list
	deviceDesc := DevInfo{
		Mac:       model.MACAddress,
		ModelName: model.Model,
		Scanproto: proto,
		Timestamp: time.Now().Format(time.RFC3339),
		IPAddress: model.IPAddress,
		Netmask:   model.Netmask,
		Gateway:   model.Gateway,
		Hostname:  model.Hostname,
		Kernel:    model.Kernel,
		Ap:        model.Ap,
		ScannedBy: model.ScannedBy,
		ArpMissed: 0,
	}
	InserAndPublishDevice(deviceDesc)
}

func InserAndPublishDevice(deviceDesc DevInfo) {
	if InsertDev(deviceDesc) {
		devinfo := make(map[string]DevInfo)
		devinfo[deviceDesc.Mac] = deviceDesc
		err := PublishDevices(&devinfo)
		if err != nil {
			q.Q(err)
		}
	}
}

func FindDev(Id string) (*DevInfo, error) {
	QC.DevMutex.Lock()
	dev, ok := QC.DevData[Id]
	QC.DevMutex.Unlock()
	if !ok {
		q.Q("device not found", Id)
		return nil, fmt.Errorf("no such device %s", Id)
	}
	q.Q(dev)
	return &dev, nil
}

func InsertDev(deviceDesc DevInfo) bool {
	lastUnixTime = strconv.FormatInt(time.Now().Unix(), 10)
	deviceDesc.UnixTime = lastUnixTime
	QC.DevMutex.Lock()
	dev, ok := QC.DevData[deviceDesc.Mac]
	QC.DevMutex.Unlock()
	if ok { //update existing entry
		// don't override gwd data with snmp data
		if deviceDesc.Scanproto == "snmp" {
			q.Q("do not override with snmp data")
			return false
		}
		// don't override with incomplete info
		if deviceDesc.ModelName != "" && deviceDesc.Ap != "" && deviceDesc.Kernel != "" &&
			deviceDesc.IPAddress != "" && deviceDesc.Netmask != "" &&
			deviceDesc.Gateway != "" && deviceDesc.Hostname != "" {

			QC.DevMutex.Lock()
			QC.DevData[deviceDesc.Mac] = deviceDesc
			QC.DevMutex.Unlock()

			q.Q("override previous entry", dev, deviceDesc, len(QC.DevData))
		}
		q.Q("device seen before", dev, deviceDesc, len(QC.DevData))
		return true
	}
	// new device discovered
	q.Q("new device", deviceDesc, len(QC.DevData))
	err := SendSyslog(LOG_ALERT, "InsertDev", "new device:"+deviceDesc.Mac)
	if err != nil {
		q.Q(err)
	}
	QC.DevMutex.Lock()
	QC.DevData[deviceDesc.Mac] = deviceDesc
	QC.DevMutex.Unlock()
	return true
}

func SaveDevices() (string, error) {
	// generate a file with timestamp
	fn := fmt.Sprintf("devices-%s.json", time.Now().Format("20060102T150405"))

	QC.DevMutex.Lock()
	jsonBytes, err := json.Marshal(QC.DevData)
	QC.DevMutex.Unlock()

	if err != nil {
		q.Q(err)
		return "", err
	}
	err = os.WriteFile(fn, jsonBytes, 0o644)
	if err != nil {
		q.Q(err)
		return "", err
	}
	return fn, nil
}

func PublishDevices(devdata *map[string]DevInfo) error {
	// send all devices info to root
	if QC.Root == "" /* || !QC.IsRoot */ {
		return fmt.Errorf("skip publishing devices, no root")
	}

	QC.DevMutex.Lock()
	jsonBytes, err := json.Marshal(devdata)
	QC.DevMutex.Unlock()

	if err != nil {
		q.Q(err)
		return err
	}
	q.Q("publishing", string(jsonBytes))

	url := QC.Root + "/api/v1/devices"
	q.Q(url)

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	resp, err := PostWithToken(url, QC.AdminToken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		q.Q(err, QC.Root)
	}
	if resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("error: publishing to devices, response status code %v", resp.StatusCode)
		}
		res := make(map[string]interface{})
		_ = json.NewDecoder(resp.Body).Decode(&res)
		q.Q(res)
	}
	return nil
}

// LoadDevices loads devices from a file, if file name is empty, it will load from the last file
func LoadDevices(fileName ...string) error {
	// load devices from file
	if !QC.IsRoot {
		return fmt.Errorf("skip loading devices, not root")
	}
	fn := ""
	if len(fileName) == 0 {
		fileslist, err := ListDevicesFiles()
		if err != nil {
			return fmt.Errorf("listing files %v", err)
		}
		if len(fileslist) == 0 {
			return fmt.Errorf("no devices files found")
		}
		// sort files by timestamp, load the last one, the file format is devices-20060102T150405.json
		sort.Slice(fileslist, func(i, j int) bool {
			return fileslist[i] > fileslist[j]
		})
		fn = fileslist[0]
		q.Q("file name empty, loading devices from last file", fn)
	} else {
		fn = fileName[0]
	}

	// read json file as DeviceInfo
	filedata, err := os.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("reading file %v", err)
	}
	devInfos := make(map[string]DevInfo)
	err = json.Unmarshal(filedata, &devInfos)
	if err != nil {
		return fmt.Errorf("unmarshalling file %v", err)
	}

	// insert devices into device list
	QC.DevMutex.Lock()
	for k, dev := range devInfos {
		QC.DevData[k] = dev
	}
	QC.DevMutex.Unlock()
	q.Q("loaded", len(devInfos), "devices from", fn)
	return nil
}

// ListDevicesFiles lists all devices files
func ListDevicesFiles() ([]string, error) {
	// list all devices files
	if !QC.IsRoot {
		return nil, fmt.Errorf("skip listing devices files, not root")
	}

	files, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var devfiles []string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "devices-") && strings.HasSuffix(f.Name(), ".json") {
			devfiles = append(devfiles, f.Name())
		}
	}
	return devfiles, nil
}

func DevicesCmd(cmdinfo *CmdInfo) *CmdInfo {
	// these should be periodically done but also available
	// for manual run
	cmd := cmdinfo.Command
	if cmd == "devices save" && QC.IsRoot {
		fn, err := SaveDevices()
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		cmdinfo.Result = fn
		return cmdinfo
	}

	if strings.HasPrefix(cmd, "devices load") && QC.IsRoot {
		ws := strings.Split(cmdinfo.Command, " ")
		if len(ws) != 3 && len(ws) != 2 {
			cmdinfo.Status = "error: invalid command"
			return cmdinfo
		}
		var err error
		if len(ws) == 2 {
			err = LoadDevices()
		} else {
			err = LoadDevices(ws[2])
		}
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		cmdinfo.Result = fmt.Sprintf("loaded %v devices", len(QC.DevData))
		return cmdinfo
	}

	if cmd == "devices files list" && QC.IsRoot {
		devfiles, err := ListDevicesFiles()
		if err != nil {
			q.Q(err)
			cmdinfo.Status = fmt.Sprintf("error: %v", err)
			return cmdinfo
		}
		cmdinfo.Status = "ok"
		cmdinfo.Result = strings.Join(devfiles, ",")
		return cmdinfo
	}

	cmdinfo.Status = "error: invalid command"
	return cmdinfo
}
