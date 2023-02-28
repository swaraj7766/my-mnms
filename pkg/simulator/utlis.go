package simulator

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"mnms/pkg/simulator/devicetype"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

type statue int

const (
	none   statue = 0
	invite statue = 1
	config statue = 2
	reboot statue = 3
	beep   statue = 4
)

//const none = 0

type CompleteSettingEvent func(ip, mask, gateway, hostname string)

const SimulatePort = 55954

func byteToHexString(msg []byte, sep string) string {
	str := make([]string, len(msg))
	for i, v := range msg {
		str[i] = fmt.Sprintf("%02X", v)
	}
	return strings.Join(str, sep)
}
func byteToString(msg []byte, sep string) string {
	str := make([]string, len(msg))
	for i, v := range msg {
		str[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(str, sep)
}

// parsing chinese
func toUtf8(b []byte) string {
	b = GetValidByte(b)
	s, err := DecodeBig5(b)
	if err != nil {
		return string(b)
	}
	return string(s)
}

func DecodeBig5(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// clear '\0' in packet
func GetValidByte(src []byte) []byte {
	var str_buf []byte
	for _, v := range src {
		if v != 0 {
			str_buf = append(str_buf, v)
		} else {

			break
		}
	}
	return str_buf
}

const ipmax = 250
const maxnumber = ipmax * ipmax //250*250
const kernel = "0.1"

// Return TestParam
//
// Exampe:Hostname=device,number=1
//
// MACAddress: "00-60-E9-2C-00-"+nmuber=00-60-E9-2C-00-01
//
// IPAddress:"10.0.1." + number=10.0.1.1
//
// Netmask: "255.255.255.0"
// Gateway: "10.0.1.254"
//
// Hostname:Hostname + number =device1
func GetTestParam(Hostname string, number uint) (ModelInfo, error) {
	if number > maxnumber {
		return ModelInfo{}, fmt.Errorf("number than %v", maxnumber)
	}
	if number == 0 {
		return ModelInfo{}, fmt.Errorf("number can't < %v", 0)
	}

	n := strconv.Itoa(int(number))
	model, ap := devicetype.GetRandomModelAp()
	return ModelInfo{Model: model,
		Ap:       ap,
		Kernel:   kernel,
		Hostname: Hostname + n}, nil
}

func GetTestParamModel(Hostname string, number uint, device_type devicetype.Simulator_type) (ModelInfo, error) {
	if number > maxnumber {
		return ModelInfo{}, fmt.Errorf("number than %v", maxnumber)
	}
	if number == 0 {
		return ModelInfo{}, fmt.Errorf("number can't < %v", 0)
	}

	n := strconv.Itoa(int(number))
	model, ap := devicetype.GetModelAp(device_type)
	return ModelInfo{Model: model,
		Ap:       ap,
		Kernel:   kernel,
		Hostname: Hostname + n}, nil
}

func MacToByte(macaddr string, sep string) []byte {
	macb := make([]byte, 6)
	mac := strings.Split(macaddr, sep)
	for i := 0; i < len(mac); i++ {
		v, _ := strconv.ParseUint(mac[i], 16, 8)
		macb[i] = byte(v)
	}
	return macb
}

// GetRandMac get radom Mac
func GetRandMac() string {
	var macs []string
	macs = append(macs, []string{"00", "60", "E9"}...)
	time.Sleep(time.Nanosecond * time.Duration(1))
	rand.Seed(time.Now().UnixNano())
	for i := 0; i <= 2; i++ {
		macs = append(macs, fmt.Sprintf("%02X", rand.Intn(255)))
	}
	mac := strings.Join(macs, "-")
	return mac
}

const atopmac = "0060e9"

// GetMac get  Mac of atop from nmuber
func GetMac(nmuber int) (string, error) {
	r := fmt.Sprintf("%v%06X", atopmac, nmuber)
	buf, err := hex.DecodeString(r)
	if err != nil {
		return "", err
	}
	return net.HardwareAddr(buf).String(), nil
}

func SelectPacket(packet []byte) statue {

	if packet[0] == 0x02 && packet[1] == 0x01 && packet[2] == 0x06 && packet[4] == 0x92 && packet[5] == 0xDA {
		return invite
	}
	if packet[0] == 0 && packet[1] == 1 && packet[2] == 6 && packet[4] == 0x92 && packet[5] == 0xDA {
		mac := make([]string, 6)
		for i := 0; i < 6; i++ {
			mac[i] = string(packet[28])
		}
		return config

	}
	if packet[0] == 0x05 && packet[1] == 0x01 && packet[2] == 0x06 && packet[4] == 0x92 && packet[5] == 0xDA {
		return reboot
	}
	if packet[0] == 0x07 && packet[1] == 0x01 && packet[2] == 0x06 && packet[4] == 0x92 && packet[5] == 0xDA {
		return beep
	}

	return none
}

// SetLogLevel sets service's log level
func SetLogLevel(cmd *cobra.Command) *logrus.Logger {
	l := logrus.New()
	if cmd == nil {
		logrus.Error("set log level fail: cmd is nil")
		return nil
	}
	debug, _ := cmd.Flags().GetBool("debug")

	if debug {
		fmt.Println("Running as debug mode")
		logrus.SetLevel(logrus.DebugLevel)
		l.SetLevel(logrus.DebugLevel)
		return l
	}

	info, _ := cmd.Flags().GetBool("verb")
	if info {
		fmt.Println("Running as verb mode")
		logrus.SetLevel(logrus.InfoLevel)
		l.SetLevel(logrus.InfoLevel)
		return l
	}
	fmt.Println("Running as default mode")
	logrus.SetLevel(logrus.ErrorLevel)
	l.SetLevel(logrus.ErrorLevel)
	return l
}
