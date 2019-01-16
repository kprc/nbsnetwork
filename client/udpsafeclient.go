package client

import (
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/recv"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"net"
	"time"
)

type udpClient struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	uw netcommon.UdpReaderWriterer

	dispatch recv.UdpRcvDispather

}


type UdpClient interface {
	Send(headinfo []byte,msgid int32,r io.ReadSeeker) error
	SendBytes(headinfo []byte,msgid int32,data []byte) error
	Rcv() error
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

	uc.uw = netcommon.NewReaderWriter(ra,conn)

	return nil
}

func (uc *udpClient)ReDial() error  {
	return uc.Dial()
}

func (uc *udpClient)SendBytes(headinfo []byte,msgid int32,data []byte) error  {
	rs := netcommon.NewReadSeeker(data)

	return uc.Send(headinfo,msgid,rs)
}

func (uc *udpClient)Send(headinfo []byte,msgid int32,r io.ReadSeeker) error  {
	bd := send.NewBlockData(r,constant.UDP_MTU)
	bd.SetWriter(uc.uw)
	inn:=nbsid.GetLocalId()
	bd.SetTransInfoOrigin(inn.String(),msgid,headinfo)
	bd.SetDataTyp(constant.DATA_TRANSER)

	sd:=send.NewStoreData(bd)

	bs:=send.GetInstance()
	bs.AddBlockDataer(bd.GetSerialNo(),sd)

	go uc.Rcv()

	bd.SendAll()

	uc.Destroy()

	time.Sleep(time.Second*1000)

	return nil
}

func (uc *udpClient)Rcv() error  {

	dispatch:=recv.NewUdpDispath(false)
	uc.dispatch = dispatch
	dispatch.SetSock(uc.uw.GetSock())
	dispatch.SetAddr(uc.uw.GetAddr())

	return dispatch.Dispatch()

}

func (uc *udpClient)Destroy()   {
	if uc.dispatch != nil {
		uc.dispatch.Close()
	}
}