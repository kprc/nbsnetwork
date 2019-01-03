package client

import (
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"net"
)

type udpClient struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	uw send.UdpReaderWriterer

	processWait chan int
}


type UdpClient interface {
	Send(headinfo []byte,msgid int32,stationId string,r io.ReadSeeker) error
	Dial() error
	ReDial() error
}

func NewUdpClient(rip,lip string,rport,lport uint16) UdpClient {
	uc := &udpClient{dialAddr:address.NewUdpAddress(),processWait:make(chan int,1024)}
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
	ipstr,port := addr.FirstS()
	if ipstr == "" {
		return nil
	}
	return &net.UDPAddr{IP:net.ParseIP(ipstr),Port:int(port)}
}

func (uc *udpClient)Dial() error {
	la := assembleUdpAddr(uc.localAddr)
	ra := assembleUdpAddr(uc.dialAddr)

	conn,err:=net.DialUDP("udp",la,ra)
	if err!=nil{
		return err
	}

	if la != nil{
		uc.realAddr = address.NewUdpAddress()
		uc.realAddr.AddIP4Str(la.String())
	}

	uc.uw = send.NewWriter(ra,conn)

	return nil
}

func (uc *udpClient)ReDial() error  {
	return uc.Dial()
}

func (uc *udpClient)Send(headinfo []byte,msgid int32,stationId string,r io.ReadSeeker) error  {
	bd := send.NewBlockData(r,constant.UDP_MTU)
	bd.SetWriter(uc.uw)

	bd.Rcv()

	bd.Send()

	return nil
}