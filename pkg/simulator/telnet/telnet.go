package telnet

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/qeof/q"
)

type TelnetServer struct {
	modelname string
	ipaddress string
	username  string
	password  string
	listener  net.Listener
}

const port = "23"
const Timeout = 200
const (
	// IAC interpret as command
	IAC = 255
	// SB is subnegotiation of the indicated option follows
	SB = 250
	// SE is end of subnegotiation parameters
	SE = 240
	// WILL indicates the desire to begin
	WILL = 251
	// WONT indicates the refusal to perform,
	// continue performing, the indicated option
	WONT = 252
	// DO indicates the request that the other
	// party perform, or confirmation that you are
	// expecting the other party to perform, the indicated option
	DO = 253
	// DONT indicates the demand that the other
	// party stop performing, or confirmation that you
	// are no longer expecting the other party to
	// perform, the indicated option
	DONT = 254
)

func skipSBSequence(reader *bufio.Reader) (err error) {
	var peeked []byte

	for {
		_, err = reader.Discard(1)
		if err != nil {
			return
		}

		peeked, err = reader.Peek(2)
		if err != nil {
			return
		}

		if peeked[0] == IAC && peeked[1] == SE {
			_, err = reader.Discard(2)
			break
		}
	}

	return
}

func skipCommand(reader *bufio.Reader) (err error) {
	var peeked []byte

	peeked, err = reader.Peek(1)
	if err != nil {
		return
	}

	switch peeked[0] {
	case WILL, WONT, DO, DONT:
		_, err = reader.Discard(2)
	case SB:
		err = skipSBSequence(reader)
	}

	return
}

func ReadByte(reader *bufio.Reader) (b byte, err error) {
	for {
		b, err = reader.ReadByte()
		if err != nil || b != IAC {
			break
		}

		err = skipCommand(reader)
		if err != nil {
			break
		}
	}

	return
}

func ReadUntil(reader *bufio.Reader, conn net.Conn, data *[]byte, delim byte) (n int, err error) {
	var b byte
	var tmp []byte
	for b != delim {
		b, err = ReadByte(reader)
		if err != nil {
			break
		}
		tmp = append(tmp, b)
		n, err := conn.Write(tmp)
		if err != nil {
			return n, err
		}

		*data = append(*data, b)
		//fmt.Println(data)
		n++
		tmp = append(tmp[:0], tmp[1:]...)
	}
	return
}

// handle communication
func handleTtlSession(conn net.Conn, ts *TelnetServer) {
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	var login []byte
	var password []byte
	var adminSuccess int
	reader := bufio.NewReader(conn)
	_, _ = conn.Write([]byte{0xff, 0xfd, 0x03})
	_, _ = conn.Write([]byte{
		0xff, 0xfb, 0x18,
		0xff, 0xfb, 0x1f,
		0xff, 0xfb, 0x20,
		0xff, 0xfb, 0x21,
		0xff, 0xfb, 0x22})
	for {

		_, _ = conn.Write([]byte("\r\n" + "Username: "))
		//n, _ := conn.Read(login)
		var n int
		adminSuccess = 0
		n, _ = ReadUntil(reader, conn, &login, '\n')
		if bytes.Compare(login[:n], []byte(ts.username+"\r\n")) != 0 {
			q.Q("invalid login ", login[:n])
		} else {
			q.Q("Username ", login[:n])
			adminSuccess = 1
		}
		// Password
		_, _ = conn.Write([]byte("\r\n" + "Password: "))
		//n, _ = conn.Read(password)
		n, _ = ReadUntil(reader, conn, &password, '\n')
		if bytes.Compare(password[:n], []byte(ts.password+"\r\n")) != 0 {
			q.Q("invalid password ", password[:n])
		} else if adminSuccess == 1 {
			q.Q("Password ", password[:n])
			break
		}
		login = login[0:0]
		password = password[0:0]
		adminSuccess = 0
	}
	// Write banner
	_, _ = conn.Write([]byte("\r\n\r\nTest " + ts.modelname + " CLI\r\n"))
	_, _ = conn.Write([]byte("switch# "))
	// telnet cmd
	var dirname string
	var inDir string = "switch"
	for {
		var str []byte
		var n int
		//var err error
		//str, err := reader.ReadString('\n')
		//fmt.Printf("%v %d\n", str, n)
		n, _ = ReadUntil(reader, conn, &str, '\n')

		if bytes.Compare(str[:n], []byte("config\r\n")) == 0 {
			dirname = "(config)"
			inDir = "config"
			//fmt.Println(string(str[:n]))
			//fmt.Println("switch" + dirname + "# ")
			_, _ = conn.Write([]byte("\r\n" + "switch" + dirname + "# "))
		} else if bytes.Compare(str[:n], []byte("snmp enable\r\n")) == 0 && inDir == "config" {
			//fmt.Println(string(str[:n]))
			//fmt.Println("switch" + "(config)# ")
			_, _ = conn.Write([]byte("\r\n" + "snmp enable"))
			_, _ = conn.Write([]byte("\r\n" + "switch" + dirname + "# "))
			//conn.Close()
			//break
		} else if bytes.Compare(str[:n], []byte("no snmp enable\r\n")) == 0 && inDir == "config" {
			//fmt.Println(string(str[:n]))
			//fmt.Println("switch" + "(config)# ")
			_, _ = conn.Write([]byte("\r\n" + "snmp disable"))
			_, _ = conn.Write([]byte("\r\n" + "switch" + dirname + "# "))
			//conn.Close()
			//break
		} else if bytes.Compare(str[:n], []byte("exit\r\n")) == 0 {
			inDir = "switch"
			if dirname == "" {
				//fmt.Println(string(str[:n]))
				conn.Close()
				break
			} else {
				dirname = ""
				fmt.Println(string(str[:n]))
				//fmt.Println("switch" + dirname + "# ")
				_, _ = conn.Write([]byte("\r\n" + "switch" + dirname + "# "))
			}
		} else {
			_, _ = conn.Write([]byte("\r\n" + "Command not found"))
			_, _ = conn.Write([]byte("\r\n" + "switch" + dirname + "# "))
		}
		str = str[0:0]
	}
	conn.Close()
}

func NewTelnetServer(modelname, ip, account, pwd string) *TelnetServer {
	telnetserver := &TelnetServer{
		modelname: modelname,
		ipaddress: ip,
		username:  account,
		password:  pwd,
	}
	return telnetserver
}

func (ts *TelnetServer) Run() error {
	go func() {
		addr := net.JoinHostPort(ts.ipaddress, port)
		l, err := net.Listen("tcp", addr)
		if err != nil {
			q.Q("Listen faild ", addr, err.Error())
			return
		}
		q.Q("listen: ", addr)
		ts.listener = l
		defer l.Close()
		for {
			conn, err := l.Accept()
			if err != nil {
				//q.Q("accept faild: ", err.Error())
				continue
			}
			defer conn.Close()
			// server communication, execute
			go handleTtlSession(conn, ts)
		}
	}()
	return nil
}

func (ts *TelnetServer) Shutdown() {
	if ts.listener != nil {
		ts.listener.Close()
	}
}
