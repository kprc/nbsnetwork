package outer

import (
	"fmt"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/dispatch"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"net"
)

type udpOut struct {
	uw netcommon.UdpReaderWriterer
	dispatch dispatch.UdpRcvDispather
	cmd chan int
}



type UdpOuter interface {
	Send(headinfo []byte,msgid int32,r io.ReadSeeker) error
	SendBytes(headinfo []byte,msgid int32,data []byte) error
	GetAddr() *net.UDPAddr
	GetSock() *net.UDPConn
	IsListen() bool
	Listen(is bool)
	Rcv() error
	GetDispatch() dispatch.UdpRcvDispather
	Destroy()
}

func NewUdpOuterUW(uw netcommon.UdpReaderWriterer) UdpOuter {
	return &udpOut{uw:uw}
}

func (uo *udpOut)GetDispatch() dispatch.UdpRcvDispather  {
	return uo.dispatch
}


func (uo *udpOut)GetAddr() *net.UDPAddr{
	return uo.uw.GetAddr()
}
func (uo *udpOut)GetSock() *net.UDPConn{
	return uo.uw.GetSock()
}
func (uo *udpOut)IsListen() bool{
	return uo.uw.IsNeedRemoteAddress()
}
func (uo *udpOut)Listen(is bool){
	uo.uw.SetNeedRemote(is)
}
func (uo *udpOut)Rcv() error{
	dispatch:=dispatch.NewUdpDispath(uo.uw)
	uo.dispatch = dispatch


	return dispatch.Dispatch()
}

func (uo *udpOut)SendBytes(headinfo []byte,msgid int32,data []byte) error  {
	rs := netcommon.NewReadSeeker(data)

	return uo.Send(headinfo,msgid,rs)
}

func (uo *udpOut)Send(headinfo []byte,msgid int32,r io.ReadSeeker) error  {
	bd := send.NewBlockData(r,constant.UDP_MTU)

	bd.SetWriter(uo.uw)
	inn:=nbsid.GetLocalId()
	bd.SetTransInfoOrigin(inn.String(),msgid,headinfo)
	bd.SetDataTyp(constant.DATA_TRANSER)

	sd:=send.NewStoreData(bd)

	bs:=send.GetBSInstance()
	bs.AddBlockDataer(bd.GetSerialNo(),sd)

	sd = bs.GetBlockDataer(bd.GetSerialNo())
	bd = sd.GetBlockData()
	fmt.Println("====Begin to receive",bd.GetSerialNo())
	if !uo.uw.IsNeedRemoteAddress() {
		go uo.Rcv()
	}


	bd.SendAll()

	fmt.Println("===begin destroy")

	uo.Destroy()
	fmt.Println("=== wait dispatch quit")
	uo.dispatch.WaitQuit()
	fmt.Println("=== dispatch quit")
	bs.PutBlockDataer(bd.GetSerialNo())

	return nil
}

func (uo *udpOut)Destroy()   {
	if uo.dispatch != nil {
		uo.dispatch.Close()
	}
}