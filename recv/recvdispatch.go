package recv

import (
	"fmt"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/flowkey"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"net"
)

type udpRcvDispath struct {
	sock *net.UDPConn
	addr *net.UDPAddr
	isListen bool
}

type UdpRcvDispather interface {
	Dispatch() error
	SetSock(sock *net.UDPConn)
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
	if rd.isListen {
		return rd.sock.ReadFromUDP(buf)
	}

	n,err:=rd.sock.Read(buf)

	return n,rd.addr,err

}




func (rd *udpRcvDispath)Dispatch() error  {
	if !rd.isListen && rd.addr == nil {
		return nbserr.NbsErr{ErrId:nbserr.UDP_RCV_DEFAULT_ERR,Errmsg:"Address is none"}
	}

	for {
		buf:=make([]byte,1024)
		n,remoteAddr,err:=rd.read(buf)

		if err!=nil {
			continue
		}

		pkt:=packet.NewUdpPacketData()

		if err := pkt.DeSerialize(buf[:n]);err!=nil{
			fmt.Println("Packet DeSerialize failed!",err)
			continue
		}

		if pkt.GetTyp() == constant.ACK {
			//send to send block
			//doAck(pkt)
			continue
		}

		if pkt.GetTyp() == constant.FINISH_ACK {
			//we can delete the msg block
		}
		
		mc:=regcenter.GetMsgCenterInstance()
		msgid,sid,hi:=mc.GetMsgId(pkt.GetTransInfo())
		if msgid == constant.MSG_NONE || sid == "" {
			continue
		}
		
		rmr:=GetInstance()
		fk := flowkey.NewFlowKey(sid,pkt.GetSerialNo())
		msg := rmr.GetMsg(fk)

		if msg == nil {
			
		}
		







	}

	return nil

}