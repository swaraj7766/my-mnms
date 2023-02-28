package firmware

import (
	"math"
	"net"
	"os"
	"time"

	"github.com/qeof/q"
	"github.com/sirupsen/logrus"
	// "time"
)

var file_sizes int64

const port = "55950"

type Firmware struct {
	ipaddress string
	listener  net.Listener
}

func writeFile(content []byte) {
	if len(content) != 0 {
		fp, err := os.OpenFile("upload.dld", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
		if err != nil {
			q.Q("open file faild: ", err)
		}
		defer fp.Close()

		_, err = fp.Write(content)
		if err != nil {
			q.Q("append content to file faild: ", err)
		}
		//log.Printf("append content: %s success\n", string(content))
	}
}

func waitpackethead(con net.Conn) (int64, error) {
	err := con.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return 0, err
	}
	packethead := make([]byte, 40)
	for {
		_, err := con.Read(packethead)
		if err != nil {
			return 0, err
		}
		if packethead[36] == 0x72 {
			//packethead[32] ~ packethead[35] : file size
			var filesize int64 = 0
			for j := 3; j >= 0; j-- {
				tmp := (int64)(packethead[j+32]) * int64(math.Pow(256, float64(j)))
				filesize = filesize + tmp
			}
			//fmt.Println(string(packethead[:]))
			return filesize, nil
		}
	}
}

func uploadfile(con net.Conn, file_size *int64) error {
	err := con.SetReadDeadline(time.Now().Add(300 * time.Second))
	if err != nil {
		return err
	}
	//tmp := file_size
	for {
		buf := make([]byte, 512)
		n, err := con.Read(buf)
		//fmt.Println(n)
		if err != nil {
			return err
		}
		if n > 0 {
			//writeFile(buf[:])
			*file_size = *file_size - (int64)(len(buf))
			//logrus.Printf("file size %d %d", file_sizes, (int64)(len(buf)))
			return nil
		}
	}
}

func serverConn(conn net.Conn) {
	defer conn.Close()
	var err error

	file_sizes, err = waitpackethead(conn)
	if err != nil {
		q.Q("head error\n")
		conn.Close()
		return
	}
	//tmp := file_sizes
	logrus.Printf("file size %d", file_sizes)
	for {
		if file_sizes <= 0 {
			//logrus.Printf("file size %d\n", file_sizes)
			time.Sleep(time.Millisecond * 100)
			//going
			_, err = conn.Write([]byte("a"))
			if err != nil {
				q.Q("write 'a' error")
				conn.Close()
				return
			}
			//erased
			_, err = conn.Write([]byte("S001"))
			if err != nil {
				q.Q("write 'S001' error")
				conn.Close()
				return
			}
			break
		} else {
			//going
			_, err = conn.Write([]byte("a"))
			if err != nil {
				q.Q("write 'a' error")
				conn.Close()
				return
			}
		}
		err = uploadfile(conn, &file_sizes)
		if err != nil {
			q.Q("upload file error. ", err)
			conn.Close()
			return
		}
	}
	//finish
	_, err = conn.Write([]byte("S002"))
	if err != nil {
		q.Q("write 'S002' error")
		conn.Close()
		return
	}
	return
}

func NewFirmwareServer(ip string) *Firmware {
	firmwareserver := &Firmware{
		ipaddress: ip,
	}
	return firmwareserver
}

func (fw *Firmware) Run() error {
	go func() {
		addr := net.JoinHostPort(fw.ipaddress, port)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			q.Q("Listen faild: ", addr, err.Error())
			return
		}

		q.Q("listen: ", addr)
		fw.listener = l
		defer l.Close()

		for {
			conn, err := l.Accept()
			if err != nil {
				//q.Q("accept faild: %v", err)
				continue
			}
			go func() {
				defer conn.Close()
				serverConn(conn)
			}()
		}
	}()
	return nil
}

func (fw *Firmware) Shutdown() {
	if fw.listener != nil {
		fw.listener.Close()
	}
}
