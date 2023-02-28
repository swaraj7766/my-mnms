package mnms

import (
	"archive/zip"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qeof/q"
)

const (
	fwport         = 55950
	size           = 512
	upgradetimeout = 300 * time.Second
	resposetimeout = 5 * time.Second
	conntimeout    = 10 * time.Second
)

type fwStatus int

func firmwarePacket() []byte {
	packet := make([]byte, 40)
	def := "name1234passwd12modelname 123456" //not important, just input anychar
	for i, v := range def {
		packet[i] = byte(v)
	}
	packet[36] = 0x72
	return packet
}

const (
	ready fwStatus = iota
	erased
	finish
	going
)

func (p fwStatus) String() string {
	switch p {
	case ready:
		return "E001"
	case erased:
		return "S001"
	case finish:
		return "S002"
	case going:
		return "a"
	}
	return "unknow"
}

type Firmware struct {
	ip         string
	filesize   int64
	firmStatus FirmStatus
	r          io.Reader
}
type FirmStatus struct {
	Status       string
	ErrorMessage string
}

// GetProcessStatus get status of upgrading fw,process ,percent
func (f *Firmware) GetProcessStatus() (string, error) {
	s := f.firmStatus.Status
	m := f.firmStatus.ErrorMessage
	if m == "" {
		return s, nil
	}
	return s, errors.New(m)
}

// Upgrading fw
func (f *Firmware) Upgrading(fileformat string, file string) error {
	// init status
	f.firmStatus.Status = "Uploading"
	f.firmStatus.ErrorMessage = ""

	q.Q("Uploading file to", f.ip)
	// start upgrading fw
	go func() {
		if fileformat == "http" {
			//download url file to data
			data, err := downloadURLFile(file)
			if err != nil {
				f.firmStatus.Status = "Error"
				f.firmStatus.ErrorMessage = err.Error()
				return
			}
			f.filesize = int64(len(data))
			f.r = bytes.NewReader(data)
		} else if fileformat == "file" {
			//open local file
			fd, err := os.Open(file)
			if err != nil {
				f.firmStatus.Status = "Error"
				f.firmStatus.ErrorMessage = err.Error()
				return
			}
			defer fd.Close()
			// file description
			fi, err := fd.Stat()
			if err != nil {
				f.firmStatus.Status = "Error"
				f.firmStatus.ErrorMessage = err.Error()
				return
			}
			f.filesize = fi.Size()
			f.r = bufio.NewReader(fd)
		}
		q.Q(f.ip, f.filesize)
		address := strings.Join([]string{f.ip, strconv.Itoa(fwport)}, ":")

		conn, err := net.DialTimeout("tcp", address, conntimeout)
		if err != nil {
			f.firmStatus.Status = "Error"
			f.firmStatus.ErrorMessage = err.Error()
			return
		}
		defer conn.Close()

		// send fw header packet
		_, err = conn.Write(downloadRequest(f.filesize))
		if err != nil {
			f.firmStatus.Status = "Error"
			f.firmStatus.ErrorMessage = fmt.Sprintf("%v err:%v", "Uploading", err.Error())
			return
		}
		err = f.waitResponse(conn, going)
		if err != nil {
			f.firmStatus.Status = "Error"
			f.firmStatus.ErrorMessage = fmt.Sprintf("%v err:%v", "Uploading", err.Error())
			return
		}

		//send file
		buf := make([]byte, 0, size)
		for {
			// wait uploading fw
			n, readerr := io.ReadFull(f.r, buf[:cap(buf)])
			// uploading process
			buf = buf[:n]
			_, err := conn.Write(buf)
			if err != nil {
				f.firmStatus.Status = "Error"
				f.firmStatus.ErrorMessage = fmt.Sprintf("%v err:%v", "Uploading", err.Error())
				return
			}
			err = f.waitResponse(conn, going)
			if err != nil {
				f.firmStatus.Status = "Error"
				f.firmStatus.ErrorMessage = fmt.Sprintf("%v err:%v", "Uploading", err.Error())
				return
			}
			// wait updraging fw
			if readerr != nil {
				if readerr == io.EOF || readerr == io.ErrUnexpectedEOF {
					err := f.waitResponse(conn, erased)
					if err != nil {
						f.firmStatus.Status = "Error"
						f.firmStatus.ErrorMessage = fmt.Sprintf("%v err:%v", "Upgrading", err.Error())
						return
					}
					break
				}
			}
		}
		// wait finish
		err = f.waitResponse(conn, finish)
		if err != nil {
			f.firmStatus.Status = "Error"
			f.firmStatus.ErrorMessage = fmt.Sprintf("%v err:%v", "Complete", err.Error())
			return
		}
		conn.Close()
	}()
	return nil
}

// waitResponse wait Response and compare to check status
func (f *Firmware) waitResponse(con net.Conn, w fwStatus) error {
	if w == going {
		err := con.SetReadDeadline(time.Now().Add(resposetimeout))
		if err != nil {
			return err
		}
	} else {
		err := con.SetReadDeadline(time.Now().Add(upgradetimeout))
		if err != nil {
			return err
		}
	}
	dst := make([]byte, len(w.String()))
	for {
		_, err := con.Read(dst)
		if err != nil {
			return err
		}
		r := strings.TrimSpace(string(dst))
		if r == w.String() {
			if w == erased {
				f.firmStatus.Status = "Upgrading"
			}
			if w == finish {
				f.firmStatus.Status = "Complete"
			}
			return nil
		}
	}
}

func downloadRequest(filesize int64) []byte {
	dl_request := firmwarePacket()
	//dl_request[32] ~ dl_request[35] :save file size
	for j := 3; j >= 0; j-- {
		dl_request[j+32] = (byte)(filesize / int64(math.Pow(256, float64(j))))
		filesize = filesize - int64(dl_request[j+32])*int64(math.Pow(256, float64(j)))
	}
	return dl_request
}

func downloadURLFile(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//create to file
	//_ = ioutil.WriteFile(filepath, data, 0755)
	//unzip
	num := int64(len(data))
	zipReader, err := zip.NewReader(bytes.NewReader(data), num)
	if err != nil {
		q.Q("warning: not a valid zip file")
		return data, nil
	}
	file := zipReader.File[0]
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	dataunzip, err2 := ioutil.ReadAll(f)
	if err2 != nil {
		return nil, err2
	}
	return dataunzip, nil
}

// commands :
//
//	firmware 00-60-E9-1E-93-D4 https://https://www.atoponline.com/.../logo-new-thinned.svg
//	firmware 00-60-E9-1E-93-D4 file:///C:/Users/testfile.txt
func FirmwareCmd(cmdinfo *CmdInfo) *CmdInfo {
	cmd := cmdinfo.Command
	ws := strings.Split(cmd, " ")
	if len(ws) < 3 {
		q.Q("error", len(ws))
		cmdinfo.Status = "error: invalid command"
		return cmdinfo
	}

	devId := ws[1]
	dev, err := FindDev(devId)
	if err != nil {
		cmdinfo.Status = "pending: device not found"
		return cmdinfo
	}
	ip := dev.IPAddress

	file := ws[2]
	fileformat := ""
	u, err := url.Parse(file)
	if err != nil {
		cmdinfo.Status = "error: url parse load error"
		return cmdinfo
	}
	q.Q(u.Scheme, u.Path)
	if u.Scheme == "http" || u.Scheme == "https" {
		fileformat = "http"
	} else if u.Scheme == "file" {
		fileformat = "file"
		file = strings.TrimPrefix(u.Path, "/")
	} else {
		cmdinfo.Status = "error: unknown file format"
		return cmdinfo
	}

	// create new  device for firmware
	fs := FirmStatus{Status: ""}
	device := Firmware{ip: ip, firmStatus: fs}

	go func() {
		err = device.Upgrading(fileformat, file)
		if err != nil {
			fmt.Println(err.Error())
		}

		var messages string = ""

		for {
			time.Sleep(time.Duration(time.Second * 1))
			r, err := device.GetProcessStatus()

			if r == "Error" {
				messages = "device:" + ip + ",Process:" + r + ",err:" + err.Error()
				q.Q(messages)
				err := SendSyslog(LOG_ALERT, "firmware", r)
				if err != nil {
					q.Q(err)
				}
				return
			}
			if r == "Complete" {
				messages = "device:" + ip + ",Process:" + r
				q.Q(messages)
				err := SendSyslog(LOG_ALERT, "firmware", r)
				if err != nil {
					q.Q(err)
				}
				return
			}
		}
	}()
	cmdinfo.Status = "ok"
	return cmdinfo
}
