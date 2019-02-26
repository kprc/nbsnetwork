package client

import (
	"github.com/NBSChain/go-nbs/utils"
	"github.com/kprc/nbsnetwork/common/address"
	"net"
)

var logger = utils.GetLogInstance()

type udpClient struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	addr *net.UDPAddr
	sock *net.UDPConn
}


type UdpClient interface {
	Dial() error
}

func NewUdpClient(rip,lip string,rport,lport uint16) UdpClient {
	uc := &udpClient{dialAddr:address.NewUdpAddress()}
	uc.dialAddr.AddIP4(rip,rport)
	if lip != "" {
		uc.localAddr = address.NewUdpAddress()
		uc.localAddr.AddIP4(lip,lport)
	}else if lport != 0{
		uc.localAddr.AddIP4("0.0.0.0",lport)
	}

	return uc
}



func assembleUdpAddr(addr address.UdpAddresser) *net.UDPAddr  {
	if addr == nil {
		return nil
	}
	ipstr,port := addr.FirstS()
	if ipstr == "" {
		return nil
	}
	return &net.UDPAddr{IP:net.ParseIP(ipstr),Port:int(port)}
}

func (uc *udpClient)Dial() error {
	la := assembleUdpAddr(uc.localAddr)
	ra := assembleUdpAddr(uc.dialAddr)

	conn,err:=net.DialUDP("udp4",la,ra)
	if err!=nil{
		return err
	}

	if la != nil{
		uc.realAddr = address.NewUdpAddress()
		uc.realAddr.AddIP4Str(la.String())
	}

	uc.sock = conn
	uc.addr = ra

	return nil
}


