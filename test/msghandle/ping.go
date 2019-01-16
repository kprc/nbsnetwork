package msghandle

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/outer"
	"io"
)

var ws netcommon.UdpBytesWriterSeeker

func RegPingMsg()  {
	mh:=regcenter.NewMsgHandler()

	mh.SetHandler(handlePing)
	mh.SetWSNew(getPingWS)

	mi:=regcenter.GetMsgCenterInstance()

	mi.AddHandler(constant.MSG_PING,mh)

}

func RegPingAckMsg()  {
	mh:=regcenter.NewMsgHandler()

	mh.SetHandler(handlePingAck)
	mh.SetWSNew(getPingAckWS)

	mi:=regcenter.GetMsgCenterInstance()

	mi.AddHandler(constant.MSG_PING_ACK,mh)

}


func getPingAckWS(param interface{}) io.WriteSeeker {
	return getPingWS(param)
}

func getPingWS(param interface{}) io.WriteSeeker {
	ws= netcommon.NewWriteSeeker(param)

	return ws

}

func handlePingAck(head interface{},data interface{},snd io.Writer) error  {

	if head != nil{
		sh:=head.([]byte)
		fmt.Println("Head is :",string(sh))
	}

	if data !=nil {
		data.(netcommon.UdpBytesWriterSeeker).PrintAll()
	}

	//uw:=snd.(netcommon.UdpReaderWriterer)
	//uo:=send.NewUdpOuter(uw.GetAddr(),uw.GetSock(),uw.IsNeedRemoteAddress())
	//uo.SendBytes([]byte("Title Ack"),constant.MSG_PING_ACK,[]byte("Send Pong"))
	//
	//fmt.Println("Send Pong")

	return nil

}

func handlePing(head interface{},data interface{},snd io.Writer) error  {

	if head != nil{
		sh:=head.([]byte)
		fmt.Println("Head is :",string(sh))
	}

	if data !=nil {
		data.(netcommon.UdpBytesWriterSeeker).PrintAll()
	}

	uw:=snd.(netcommon.UdpReaderWriterer)
	uo:=outer.NewUdpOuter(uw.GetAddr(),uw.GetSock(),uw.IsNeedRemoteAddress())
	uo.SendBytes([]byte("Title Ack"),constant.MSG_PING_ACK,[]byte("Send Pong"))

	fmt.Println("Send Pong")

	return nil

}