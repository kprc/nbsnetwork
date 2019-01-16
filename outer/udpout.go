package outer

import (
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/dispatch"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"net"
	"time"
)

type udpOut struct {
	addr *net.UDPAddr
	sock *net.UDPConn
	isListen bool
	dispatch dispatch.UdpRcvDispather
}



type UdpOuter interface {
	Send(headinfo []byte,msgid int32,r io.ReadSeeker) error
	SendBytes(headinfo []byte,msgid int32,data []byte) error
	SetAddr(addr *net.UDPAddr)
	SetSock(sock *net.UDPConn)
	GetAddr() *net.UDPAddr
	GetSock() *net.UDPConn
	IsListen() bool
	Listen(is bool)
	Rcv() error
	GetDispatch() dispatch.UdpRcvDispather
	Destroy()
}

func NewUdpOuter(addr *net.UDPAddr,sock *net.UDPConn, is bool) UdpOuter {
	return &udpOut{addr,sock,is,nil}
}
func (uo *udpOut)GetDispatch() dispatch.UdpRcvDispather  {
	return uo.dispatch
}

func (uo *udpOut)SetAddr(addr *net.UDPAddr){
	uo.addr = addr
}
func (uo *udpOut)SetSock(sock *net.UDPConn){
	uo.sock = sock
}
func (uo *udpOut)GetAddr() *net.UDPAddr{
	return uo.addr
}
func (uo *udpOut)GetSock() *net.UDPConn{
	return uo.sock
}
func (uo *udpOut)IsListen() bool{
	return uo.isListen
}
func (uo *udpOut)Listen(is bool){
	uo.isListen = is
}
func (uo *udpOut)Rcv() error{
	dispatch:=dispatch.NewUdpDispath(uo.isListen)
	uo.dispatch = dispatch
	dispatch.SetSock(uo.GetSock())
	dispatch.SetAddr(uo.GetAddr())

	return dispatch.Dispatch()
}

func (uo *udpOut)SendBytes(headinfo []byte,msgid int32,data []byte) error  {
	rs := netcommon.NewReadSeeker(data)

	return uo.Send(headinfo,msgid,rs)
}

func (uo *udpOut)Send(headinfo []byte,msgid int32,r io.ReadSeeker) error  {
	bd := send.NewBlockData(r,constant.UDP_MTU)
	uw:=netcommon.NewReaderWriter(uo.GetAddr(),uo.GetSock())
	uw.SetNeedRmAddr(uo.isListen)

	bd.SetWriter(uw)
	inn:=nbsid.GetLocalId()
	bd.SetTransInfoOrigin(inn.String(),msgid,headinfo)
	bd.SetDataTyp(constant.DATA_TRANSER)

	sd:=send.NewStoreData(bd)

	bs:=send.GetInstance()
	bs.AddBlockDataer(bd.GetSerialNo(),sd)

	go uo.Rcv()

	bd.SendAll()

	uo.Destroy()

	time.Sleep(time.Second*1)

	return nil
}

func (uo *udpOut)Destroy()   {
	if uo.dispatch != nil {
		uo.dispatch.Close()
	}
}