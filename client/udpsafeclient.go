package client

import (
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/outer"
	"github.com/kprc/nbsnetwork/netcommon"
	"io"
	"net"
)

type udpClient struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	uw netcommon.UdpReaderWriterer
	uo outer.UdpOuter
}


type UdpClient interface {
	Send(headinfo []byte,msgid int32,r io.ReadSeeker) error
	SendTimeOut(headinfo []byte,msgid int32,r io.ReadSeeker,tvsecond int) error
	SendBytesTimeOut(headinfo []byte,msgid int32,data []byte,tvsecond int) error
	Dial() error
	ReDial() error
	Destroy()
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

func (uc *udpClient)SendTimeOut(headinfo []byte,msgid int32,r io.ReadSeeker,tvsecond int) error {
	return nil
}

func (uc *udpClient)SendBytesTimeOut(headinfo []byte,msgid int32,data []byte,tvsecond int) error  {
	return nil
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

	uc.uw = netcommon.NewReaderWriter(ra,conn,false)

	return nil
}

func (uc *udpClient)ReDial() error  {
	return uc.Dial()
}

func (uc *udpClient)SendBytes(headinfo []byte,msgid int32,data []byte) error  {
	uo:=outer.NewUdpOuter(uc.uw.GetAddr(),uc.uw.GetSock(),false)
	uc.uo = uo

	return uo.SendBytes(headinfo,msgid,data)
}

func (uc *udpClient)Send(headinfo []byte,msgid int32,r io.ReadSeeker) error  {
	uo:=outer.NewUdpOuter(uc.uw.GetAddr(),uc.uw.GetSock(),false)

	uc.uo = uo

	return uo.Send(headinfo,msgid,r)

}



func (uc *udpClient)Destroy()   {
	uc.uo.Destroy()
}

