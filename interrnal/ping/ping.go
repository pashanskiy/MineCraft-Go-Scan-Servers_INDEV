package ping

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

type PingServer struct {
	Address *net.TCPAddr
	connect *net.TCPConn
}

type ServerData struct {
	Description struct {
		Extra []struct {
			Text string `json:"text"`
		} `json:"extra"`
		Text string `json:"text"`
	} `json:"description"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"sample"`
	} `json:"players"`
	Version struct {
		Name     string `json:"name"`
		Protocol int    `json:"protocol"`
	} `json:"version"`
	Ping int64
}

func (PS *PingServer) NewPing(ip, port string) error {
	var err error

	if validIP4(ip) {
		PS.Address, err = net.ResolveTCPAddr("tcp", ip+":"+port)
		return err
	}
	ips, _ := net.LookupIP(ip)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			fmt.Println("IPv4: ", ipv4)
			PS.Address, err = net.ResolveTCPAddr("tcp", ipv4.String()+":"+port)
		}
	}
	return err
}

func (PS *PingServer) GetConnect() error {
	var err error
	PS.connect, err = getTCPConnect(PS.Address)
	return err
}

func (PS *PingServer) RequestInfoAndUnmarshall() (ServerData, error) {
	bytes, ping, err := requestServerInfo(PS.connect, PS.Address)
	ServerData := ServerData{}
	ServerData.Ping = ping
	if err != nil {
		return ServerData, err
	}

	err = json.Unmarshal(bytes, &ServerData)
	return ServerData, err
}

func validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	return re.MatchString(ipAddress) 
}

func getTCPConnect(tcpAddr *net.TCPAddr) (*net.TCPConn, error) {

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, errors.New("Dial failed: " + err.Error())
	}

	return conn, nil
}

func buildHandshake(tcpAddr *net.TCPAddr) *bytes.Buffer {
	bytes := &bytes.Buffer{}
	// packet id
	bytes.WriteByte(0x00)
	// protocol version
	bytes.WriteByte(0xff)
	//?
	bytes.WriteByte(0x00)
	// varint ip
	// bytes.WriteByte(0x0d)
	bytes.Write(toVarint(uint64(len(tcpAddr.IP.String()))))
	// string ip
	bytes.WriteString(string(tcpAddr.IP.String()))
	// int16 port
	_ = binary.Write(bytes, binary.BigEndian, int16(tcpAddr.Port))
	bytes.WriteByte(0x01)
	return addPacketLen(bytes)
}

func addPacketLen(byt *bytes.Buffer) *bytes.Buffer {
	buffer := &bytes.Buffer{}
	buffer.Write(toVarint(uint64(byt.Len())))
	buffer.Write(byt.Bytes())
	return buffer
}

func toVarint(x uint64) []byte {
	var bytes [10]byte
	var n int
	for n = 0; x > 127; n++ {
		bytes[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	bytes[n] = uint8(x)
	n++
	return bytes[0:n]
}

func requestServerInfo(conn *net.TCPConn, addr *net.TCPAddr) ([]byte, int64, error) {
	defer conn.Close()
	stopwatchStart := time.Now()
	//handshake
	_, err := conn.Write(buildHandshake(addr).Bytes())
	if err != nil {
		return nil, 0, err
	}
	//request data
	_, err = conn.Write([]byte{0x01, 0x00})
	if err != nil {
		return nil, 0, err
	}
	//read response message length
	reply := make([]byte, 5)
	_, err = conn.Read(reply)
	if err != nil {
		return nil, 0, err
	}
	//read response message
	bytesToRead, _ := binary.Uvarint(reply)
	reply = make([]byte, bytesToRead-3)
	_, err = conn.Read(reply)
	return reply, time.Since(stopwatchStart).Milliseconds(), err
}
