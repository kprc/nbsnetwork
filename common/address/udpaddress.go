package address

import (
	"bytes"
	"fmt"
	"github.com/kprc/nbsdht/nbserr"
	"net"
	"strconv"
	"strings"
)

var udpparseerr = nbserr.NbsErr{ErrId: nbserr.UDP_ADDR_PARSE, Errmsg: "Parse ip4 address fault"}
var ipreadableerr = nbserr.NbsErr{ErrId: nbserr.UDP_ADDR_TOSTRING_ERR, Errmsg: "byte of ip4 address to string ip4 address error"}

const (
	IP_TYPE_IP4 int = 0
)

type UdpAddresser interface {
	AddIP4(ipstr string, port uint16) error
	AddIP4Str(ipstr string) error
	DelIP4(ipstr string, port uint16)
	Iterator()
	GetAddr(idx int) (addr []byte, port uint16)
	Next() (addr []byte, port uint16, idx int)
	First() (addr []byte, port uint16)
	GetAddrS(idx int) (saddr string, port uint16)
	NextS() (saddr string, port uint16, idx int)
	FirstS() (saddr string, port uint16)
	PrintAll()
	Clone() UdpAddresser
	Len() int
	Equals(ua UdpAddresser) bool
}

type address struct {
	ip6type int
	addr    []byte //big-endian
	port    uint16
}

type udpAddress struct {
	addrs []address
	pos   int
}

func NewUdpAddress() UdpAddresser {
	var addrs udpAddress

	addrs.addrs = make([]address, 0)

	return &addrs
}

func NewUdpAddressP(addr []byte, port uint16) UdpAddresser {
	var addrs udpAddress

	addrs.addrs = make([]address, 1)

	addrs.append(address{ip6type: IP_TYPE_IP4, addr: addr, port: port})

	return &addrs
}

func (ua *udpAddress) DelIP4(ipstr string, port uint16) {
	byteIp := ipString2Byte(ipstr)
	idx := -1
	for i, addr := range ua.addrs {
		if bytes.Compare(byteIp, addr.addr) == 0 && port == addr.port {
			idx = i
			break
		}
	}

	if idx == -1 {
		return
	}

	l := len(ua.addrs)
	if idx < l-1 {
		ua.addrs[idx] = ua.addrs[l-1]

	}
	ua.addrs = ua.addrs[:l-1]
}

func (ua *udpAddress) Len() int {
	return len(ua.addrs)
}

func NewUdpAddressS(ipstr string, port uint16) (error, UdpAddresser) {

	var addrs udpAddress
	sarr := strings.Split(ipstr, ".")

	if len(sarr) != 4 {
		return udpparseerr, nil
	}

	baddr := address{ip6type: IP_TYPE_IP4}

	baddr.addr = make([]byte, 0)

	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(sarr[i])
		if err != nil || n > 255 || n < 0 {
			return udpparseerr, nil
		}
		baddr.addr = append(baddr.addr)
	}

	baddr.port = port

	addrs.append(baddr)

	return nil, &addrs
}

func (uaddr *udpAddress) Iterator() {
	uaddr.pos = 0
}

func (uaddr *udpAddress) Next() (addr []byte, port uint16, idx int) {
	if len(uaddr.addrs) > uaddr.pos {
		addr = uaddr.addrs[uaddr.pos].addr
		port = uaddr.addrs[uaddr.pos].port
		idx = uaddr.pos
		uaddr.pos++

		return
	}

	return nil, 0, 0
}

func (uaddr *udpAddress) GetAddr(idx int) (addr []byte, port uint16) {
	if len(uaddr.addrs) > idx {
		addr = uaddr.addrs[idx].addr
		port = uaddr.addrs[idx].port
		return
	}

	return nil, 0
}

func stringIP(addr []byte) (error, string) {
	if len(addr) != 4 {
		return ipreadableerr, ""
	}

	saddr := make([]string, 4)

	for i := 0; i < len(addr); i++ {
		saddr[i] = strconv.Itoa(int(addr[i]))
	}

	return nil, strings.Join(saddr, ".")
}

func (uaddr *udpAddress) GetAddrS(idx int) (saddr string, port uint16) {
	if len(uaddr.addrs) > idx {
		addr := uaddr.addrs[idx].addr
		port = uaddr.addrs[idx].port

		var err error
		if err, saddr = stringIP(addr); err == nil {
			return
		}
	}

	return "", 0
}

func (uaddr *udpAddress) NextS() (saddr string, port uint16, idx int) {
	if len(uaddr.addrs) > uaddr.pos {
		addr := uaddr.addrs[uaddr.pos].addr
		port = uaddr.addrs[uaddr.pos].port

		idx = uaddr.pos

		uaddr.pos++

		var err error

		if err, saddr = stringIP(addr); err == nil {
			return
		}

	}

	return "", 0, 0
}

func (uaddr *udpAddress) First() (addr []byte, port uint16) {
	if len(uaddr.addrs) > 0 {
		addr = uaddr.addrs[0].addr
		port = uaddr.addrs[0].port

		return
	}

	return nil, 0
}

func (uaddr *udpAddress) FirstS() (saddr string, port uint16) {
	if len(uaddr.addrs) > 0 {
		addr := uaddr.addrs[0].addr
		port = uaddr.addrs[0].port

		var err error

		if err, saddr = stringIP(addr); err == nil {
			return
		}
	}

	return "", 0
}

func (uaddr *udpAddress) append(addr address) {
	uaddr.addrs = append(uaddr.addrs, addr)
}

func (uaddr *udpAddress) AddIP4(ipstr string, port uint16) error {

	baddr := address{ip6type: IP_TYPE_IP4}

	addr := ipString2Byte(ipstr)

	baddr.addr = addr

	baddr.port = port

	uaddr.append(baddr)

	return nil

}

func ipString2Byte(ipstr string) []byte {
	sarr := strings.Split(ipstr, ".")

	addr := make([]byte, 0)
	if len(sarr) != 4 {
		return addr
	}

	for i := 0; i < 4; i++ {
		n, err := strconv.Atoi(sarr[i])
		if err != nil || n > 255 || n < 0 {
			return addr
		}

		addr = append(addr, byte(n))
	}

	return addr
}

func (uaddr *udpAddress) Clone() UdpAddresser {
	ua := &udpAddress{addrs: make([]address, 0)}

	for _, addr := range uaddr.addrs {
		r := address{}
		r.ip6type = addr.ip6type
		r.port = addr.port
		r.addr = make([]byte, len(addr.addr))
		copy(r.addr, addr.addr)
		ua.addrs = append(ua.addrs, r)
	}

	return ua
}

func GetAllLocalIPAddr(port uint16) UdpAddresser {

	ua := &udpAddress{addrs: make([]address, 0)}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	for _, addr := range addrs {

		//if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
		if ipnet, ok := addr.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil && !strings.Contains(ipnet.IP.String(), "169.254") {
				r := address{}
				r.ip6type = IP_TYPE_IP4
				r.port = port
				r.addr = ipString2Byte(ipnet.IP.String())
				ua.addrs = append(ua.addrs, r)
			}
		}
	}

	return ua
}

func (uaddr *udpAddress) PrintAll() {
	uaddr.Iterator()
	for {
		addr, port, _ := uaddr.NextS()
		if addr == "" {
			break
		}
		fmt.Println(addr, port)
	}
}

func (uaddr *udpAddress) AddIP4Str(ipstr string) error {

	ipstrport := strings.Split(ipstr, ":")

	port, err := strconv.Atoi(ipstrport[1])
	if err != nil {
		return err
	}

	return uaddr.AddIP4(ipstrport[0], uint16(port))

}

func (uaddr *udpAddress) Equals(ua UdpAddresser) bool {
	addr, port := ua.First()
	if addr == nil {
		return false
	}

	addr1, port1 := uaddr.First()
	if addr1 == nil {
		return false
	}

	if bytes.Compare(addr, addr1) == 0 && port == port1 {
		return true
	}

	return false
}

//not support in this version
//func (uaddr *udpAddress)Add6(ipstr string, port uint16)  {
//
//}
