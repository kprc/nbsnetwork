package dispatch

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/flowkey"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/pb/message"
	"github.com/kprc/nbsnetwork/recv"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"net"
)

type udpRcvDispath struct {
	sock *net.UDPConn
	addr *net.UDPAddr
	isListen bool
	cmd chan int   // 1 for quit
}

type UdpRcvDispather interface {
	Dispatch() error
	SetSock(sock *net.UDPConn)
	SetAddr(addr *net.UDPAddr)
	Close() error
}


func NewUdpDispath(listen bool)  UdpRcvDispather{
	return &udpRcvDispath{isListen:listen}
}


func (rd *udpRcvDispath)SetSock(sock *net.UDPConn){
	rd.sock = sock
}

func (rd *udpRcvDispath)SetAddr(addr *net.UDPAddr)  {
	rd.addr = addr
}

func (rd *udpRcvDispath)read(buf []byte) (int,*net.UDPAddr,error)  {

	if rd.sock == nil{
		return 0,nil,nbserr.NbsErr{Errmsg:"sock is none",ErrId:nbserr.UDP_RCV_DEFAULT_ERR}
	}
	if rd.isListen {
		return rd.sock.ReadFromUDP(buf)
	}

	n,err:=rd.sock.Read(buf)

	return n,rd.addr,err

}

func (rd *udpRcvDispath)Close() error  {
	var err error = nil
	if rd.sock != nil && !rd.isListen{
		err= rd.sock.Close()
		rd.sock = nil
	}

	rd.cmd <- 1

	return err
}

func toPkt(n int,buf []byte)packet.UdpPacketDataer  {
	pkt:=packet.NewUdpPacketData()
	if err := pkt.DeSerialize(buf[:n]);err!=nil{
		fmt.Println("Packet DeSerialize failed!",err)
		return nil
	}

	return pkt
}

func (rd *udpRcvDispath)newMsg(msgid int32,sn uint64, sid string, hi []byte,uw netcommon.UdpReaderWriterer) recv.RcvMsg {
	m:=recv.NewRcvMsg()

	mc:=regcenter.GetMsgCenterInstance()
	h:=mc.GetHandler(msgid)
	defer mc.PutHandler(msgid)
	fNewWS:=h.GetWSNew()
	if fNewWS == nil {
		return nil
	}
	ws:=fNewWS(hi)
	if ws == nil{
		return nil
	}
	rcv:=recv.NewRcvDataer(sid,sn,ws,uw)
	m.SetRecv(rcv)

	return m
}

func (rd *udpRcvDispath)handleMsg(msgid int32,hi []byte,rcv recv.RcvDataer)  {
	mc:=regcenter.GetMsgCenterInstance()
	h:=mc.GetHandler(msgid)
	defer mc.PutHandler(msgid)
	fdo := h.GetHandler()
	fdo(hi,rcv.GetWs(),rcv.GetUdpReadWriter())
}

func (rd *udpRcvDispath)sendAck(w io.Writer, ack packet.UdpResulter, pkt packet.UdpPacketDataer){
	sndpkt:=packet.NewUdpPacketData()
	mc:=regcenter.GetMsgCenterInstance()
	msgid,_,hi:=mc.GetMsgId(pkt.GetTransInfo())
	sndpkt.SetSerialNo(pkt.GetSerialNo())
	sndpkt.SetACK()
	mh:=message.MsgHead{LocalStationId:[]byte(nbsid.GetLocalId().String()),MsgId:msgid}
	mh.Headinfo = hi

	bmh,err := proto.Marshal(&mh)
	if err!=nil {
		return
	}
	sndpkt.SetTransInfo(bmh)
	back,_:=ack.Serialize()
	sndpkt.SetData(back)
	sndpkt.SetLength(int32(len(back)))

	bsnd,err:= sndpkt.Serialize()
	if err!=nil {
		return
	}
	w.Write(bsnd)

}

func (rd *udpRcvDispath)doAck(pkt packet.UdpPacketDataer)  {

	pkt.PrintAll()

	bs:=send.GetInstance()

	sd:=bs.GetBlockDataer(pkt.GetSerialNo())
	if sd == nil {
		return
	}
	bdr:=sd.GetBlockData()
	if bdr!=nil {
		bdr.PushResult(pkt.GetData())
	}

	bs.PutBlockDataer(pkt.GetSerialNo())

}


func (rd *udpRcvDispath)Dispatch() error  {
	if !rd.isListen && rd.addr == nil {
		return nbserr.NbsErr{ErrId:nbserr.UDP_RCV_DEFAULT_ERR,Errmsg:"Address is none"}
	}

	for {
		buf:=make([]byte,1024)
		//fmt.Println("Dispatch to read")
		n,remoteAddr,err:=rd.read(buf)
		//fmt.Println("read byte number",n)

		if err!=nil {
			select {
			case cmd:=<-rd.cmd:
				if cmd == 1{
					break
				}
			default:
				//fmt.Println("error")
				continue
			}
		}

		var pkt packet.UdpPacketDataer
		if pkt=toPkt(n,buf[:n]); pkt==nil{
			continue
		}

		if pkt.GetTyp() == constant.ACK {
			rd.doAck(pkt)
			continue
		}

		mc:=regcenter.GetMsgCenterInstance()
		msgid,sid,hi:=mc.GetMsgId(pkt.GetTransInfo())
		if msgid == constant.MSG_NONE || sid == "" {
			continue
		}
		
		rmr:=recv.GetInstance()
		fk := flowkey.NewFlowKey(sid,pkt.GetSerialNo())
		if pkt.GetTyp() == constant.FINISH_ACK {
			rmr.DelMsg(fk)
			continue
		}
		msg := rmr.GetMsg(fk)

		if msg == nil {
			uw := netcommon.NewReaderWriter(remoteAddr,rd.sock,rd.isListen)
			msg = rd.newMsg(msgid,pkt.GetSerialNo(),sid,hi,uw)
			rmr.AddMSG(fk,msg)
			msg = rmr.GetMsg(fk)
		}

		rcv:=msg.GetRecv()
		if rcv.GetKey() == nil{
			rcv.SetKey(fk)
		}

		if ack,_ := rcv.Write(pkt); ack!=nil{
			rd.sendAck(rcv.GetUdpReadWriter(),ack,pkt)
		}

		if rcv.Finish() {
			rd.handleMsg(msgid,hi,rcv)
		}

		rmr.PutMsg(fk)
	}

	return nil

}
