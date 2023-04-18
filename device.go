package mnms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/qeof/q"
)

type DevInfo struct {

	Mac            string `json:"mac"`
	ModelName      string `json:"modelname"`
	Timestamp      string `json:"timestamp"`
	Scanproto      string `json:"scanproto"`
	IPAddress      string `json:"ipaddress"`
	Netmask        string `json:"netmask"`
	Gateway        string `json:"gateway"`
	Hostname       string `json:"hostname"`
	Kernel         string `json:"kernel"`
	Ap             string `json:"ap"`
	ScannedBy      string `json:"scannedby"`
	ArpMissed      int    `json:"arpmissed"`
	Lock           bool   `json:"lock"`
	ReadCommunity  string `json:"readcommunity"`
	WriteCommunity string `json:"writecommunity"`

}

var specialMac = "11-22-33-44-55-66"

var lastTimestamp string

func init() {
	QC.DevData = make(map[string]DevInfo)
	lastTimestamp = strconv.FormatInt(time.Now().Unix(), 10)
}

// InsertCommunities inserts communities into device list
func InsertCommunities(mac, read, write string) error {
	devinfo, err := FindDev(mac)
	if err != nil {
		return err
	}
	devinfo.ReadCommunity = read
	devinfo.WriteCommunity = write
	InsertAndPublishDevice(*devinfo)
	return nil
}

func InsertModel(model GwdModelInfo, proto string) {

	var deviceDesc DevInfo
	devinfo, err := FindDev(model.MACAddress)
	if err == nil {
		deviceDesc = *devinfo

	}
	// new device give default valuse
	//discovered device model will be entered into device list
	deviceDesc.Mac = model.MACAddress
	deviceDesc.ModelName = model.Model
	deviceDesc.Scanproto = proto
	deviceDesc.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	deviceDesc.IPAddress = model.IPAddress
	deviceDesc.Netmask = model.Netmask
	deviceDesc.Gateway = model.Gateway
	deviceDesc.Hostname = model.Hostname
	deviceDesc.Kernel = model.Kernel
	deviceDesc.Ap = model.Ap
	deviceDesc.ScannedBy = model.ScannedBy
	deviceDesc.ArpMissed = 0
	InsertAndPublishDevice(deviceDesc)
}

func InsertAndPublishDevice(deviceDesc DevInfo) {
	if InsertDev(deviceDesc) {
		devinfo := make(map[string]DevInfo)
		devinfo[deviceDesc.Mac] = deviceDesc
		err := PublishDevices(&devinfo)
		if err != nil {
			q.Q(err)
		}
	}
}

func FindDevWithIP(ip string) (*DevInfo, error) {
	QC.DevMutex.Lock()
	defer QC.DevMutex.Unlock()
	for _, dev := range QC.DevData {
		if dev.IPAddress == ip {
			return &dev, nil
		}
	}
	return nil, fmt.Errorf("no such device %s", ip)
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

func LockDev(Id string) {
	QC.DevMutex.Lock()
	defer QC.DevMutex.Unlock()
	dev, ok := QC.DevData[Id]

	if ok {
		dev.Lock = true
		QC.DevData[Id] = dev
	}
}

func unLockDev(Id string) {
	QC.DevMutex.Lock()
	defer QC.DevMutex.Unlock()
	dev, ok := QC.DevData[Id]
	if ok {
		dev.Lock = false
		QC.DevData[Id] = dev
	}
}

func DevIsLocked(Id string) (bool, error) {
	QC.DevMutex.Lock()
	defer QC.DevMutex.Unlock()
	dev, ok := QC.DevData[Id]
	if ok {
		return dev.Lock, nil
	} else {
		return false, fmt.Errorf("no such device %s", Id)
	}
}

func InsertDev(deviceDesc DevInfo) bool {
	lastTimestamp = strconv.FormatInt(time.Now().Unix(), 10)
	deviceDesc.Timestamp = lastTimestamp
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
		if deviceDesc.ModelName != "" &&
			deviceDesc.Ap != "" &&
			deviceDesc.Kernel != "" &&
			deviceDesc.IPAddress != "" &&
			deviceDesc.Netmask != "" &&
			deviceDesc.Gateway != "" {
			// XXX deviceDesc.Hostname check not done here
			//   to allow for empty hostname on some devices
			QC.DevMutex.Lock()
			QC.DevData[deviceDesc.Mac] = deviceDesc
			QC.DevMutex.Unlock()

			q.Q("override previous entry", dev, deviceDesc, len(QC.DevData))
			return true
		}
		q.Q("incomplete device seen before", dev, deviceDesc, len(QC.DevData))
		return false
	}
	// new device discovered
	q.Q("new device", deviceDesc, len(QC.DevData))
	err := SendSyslog(LOG_ALERT, "InsertDev", "new device: "+deviceDesc.Mac)
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
	if QC.RootURL == "" /* || !QC.IsRoot */ {
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

	url := QC.RootURL + "/api/v1/devices"
	q.Q(url)

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBytes))
	resp, err := PostWithToken(url, QC.AdminToken, bytes.NewBuffer(jsonBytes))
	if err != nil {
		q.Q(err, QC.RootURL)
	}
	if resp != nil {
		//save close, in resp != nil block
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return fmt.Errorf("error: publishing to devices, response status code %v", resp.StatusCode)
		}
		res := make(map[string]interface{})
		err := json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return err
		}
		q.Q(res)
	}

	return nil
}
