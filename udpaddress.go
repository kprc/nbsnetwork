package nbsnetwork

import (
	"github.com/kprc/nbsdht/nbserr"
	"strings"
)

type UdpAddresser interface {

}


type address struct {
	ip6type bool
	addr []byte		//big-endian
	port uint16
}

type udpAddress struct {
	addrs []address
}

func NewUdpAddress() UdpAddresser  {
	var addrs udpAddress

	addrs.addrs = make([]address,0)

	return &addrs
}

func (uaddr *udpAddress)Append(addr address) {
	uaddr.addrs = append(uaddr.addrs, addr)
}

func (uaddr *udpAddress)AddIP4(ipstr string, port uint16) error {
	sarr := strings.Split(ipstr,".")

	if len(sarr) != 4 {
		return &nbserr.NbsErr{ErrId:nbserr.UDP_ADDR_PARSE,Errmsg:"Parse ip4 address fault"}
	}



}

func (uaddr *udpAddress)Add6(ipstr string, port uint16)  {

}