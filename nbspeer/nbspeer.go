package nbspeer

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/address"
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
	"sync"
)

type peer struct {
	addrs address.UdpAddresser
	net netcommon.UdpReaderWriterer
	stationId string
	runningLock sync.Mutex
	runningSend bool
	runningRecv bool
	data2Send chan send.BlockDataer
}


type NbsPeer interface {
	AddIPAddr(ip string,port uint16)
	DelIpAddr(ip string,port uint16)
	GetNet() netcommon.UdpReaderWriterer
	SetNet(net netcommon.UdpReaderWriterer)
	SendAsync(msgid int32,headinfo []byte,data []byte, rcvSn uint64) (*chan int, uint64,  error)
	SendLargeDataAsync(msgid int32,headinfo []byte,rs io.ReadSeeker, rcvSn uint64)(*chan int,uint64, error)
	SendSync(msgid int32, headinfo []byte,data []byte, rcvSn uint64) (uint64,error)
	SendSyncTimeOut(msgid int32,headinfo []byte,data []byte, rcvSn uint64, ms int) (uint64,error)
	//WaitResult(sn uint64) (interface{},error)
	Wait(ch *chan int) error
	Run()
}

func NewNbsPeer(sid string) NbsPeer  {
	return &peer{stationId:sid,data2Send:make(chan send.BlockDataer,32)}
}

func (p *peer)sendbd(bd send.BlockDataer) error {

	return bd.Send2Peer()
}

func (p *peer)Run()  {
	if !p.runningSend  {
		go p.Sendbd()
	}

	if !p.runningRecv && !p.net.IsNeedRemoteAddress(){
		go p.recv()
	}
}


func (p *peer) Sendbd() error{

	if p.runningSend {
		return nbserr.NbsErr{ErrId:nbserr.UDP_SND_RUNNING,Errmsg:"udp Send is running at "+p.net.AddrString()}
	}
	p.runningLock.Lock()
	if p.runningSend {
		p.runningLock.Unlock()
		return nbserr.NbsErr{ErrId:nbserr.UDP_SND_RUNNING,Errmsg:"udp Send is running at "+p.net.AddrString()}
	}else {
		p.runningSend = true
	}
	p.runningLock.Unlock()

	var err error

	for{
		bd:=<-p.data2Send
		if err =p.sendbd(bd); err!=nil{
			break
		}
	}




	return err
}

func toPkt(buf []byte)packet.UdpPacketDataer  {
	pkt:=packet.NewUdpPacketData()
	if err := pkt.DeSerialize(buf);err!=nil{
		fmt.Println("Packet DeSerialize failed!",err)
		return nil
	}

	return pkt
}

func doAck(pkt packet.UdpPacketDataer)  {

	bs:=send.GetBSInstance()

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

func newMsg(msgid int32,sn uint64, sid string, hi []byte,uw netcommon.UdpReaderWriterer) recv.RcvMsg {
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

func sendAck(w io.Writer, ack packet.UdpResulter, pkt packet.UdpPacketDataer) {
	sndpkt := packet.NewUdpPacketData()
	mc := regcenter.GetMsgCenterInstance()
	msgid, _, hi := mc.GetMsgId(pkt.GetTransInfo())
	sndpkt.SetSerialNo(pkt.GetSerialNo())
	sndpkt.SetACK()
	mh := message.MsgHead{LocalStationId: []byte(nbsid.GetLocalId().String()), MsgId: msgid}
	mh.Headinfo = hi

	bmh, err := proto.Marshal(&mh)
	if err != nil {
		return
	}
	sndpkt.SetTransInfo(bmh)
	back, _ := ack.Serialize()
	sndpkt.SetData(back)
	sndpkt.SetLength(int32(len(back)))

	bsnd, err := sndpkt.Serialize()
	if err != nil {
		return
	}
	w.Write(bsnd)
}

func handleMsg(msgid int32,hi []byte,rcv recv.RcvDataer)  {
	mc:=regcenter.GetMsgCenterInstance()
	h:=mc.GetHandler(msgid)
	defer mc.PutHandler(msgid)
	fdo := h.GetHandler()
	fdo(hi,rcv.GetWs(),rcv.GetUdpReadWriter())
}

func (p *peer)recv() error{
	if p.net.IsNeedRemoteAddress(){
		return nil
	}

	if p.runningRecv {
		return nbserr.NbsErr{ErrId:nbserr.UDP_RCV_RUNNING,Errmsg:"udp receive is running at "+p.net.AddrString()}
	}
	p.runningLock.Lock()
	if p.runningRecv {
		p.runningLock.Unlock()
		return nbserr.NbsErr{ErrId:nbserr.UDP_RCV_RUNNING,Errmsg:"udp receive is running at "+p.net.AddrString()}
	}else {
		p.runningRecv = true
	}
	p.runningLock.Unlock()

	buf := make([]byte,1024)
	var n int
	var remoteAddr *net.UDPAddr
	var err error

	for{
		n,remoteAddr,err=p.net.ReadUdp(buf)
		if err!=nil{
			break
		}
		var pkt packet.UdpPacketDataer
		if pkt=toPkt(buf[:n]);pkt==nil {
			continue
		}
		if pkt.GetTyp() == constant.ACK {
			doAck(pkt)
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
			uw := netcommon.NewReaderWriter(remoteAddr,p.net.GetSock(),p.net.IsNeedRemoteAddress())
			msg = newMsg(msgid,pkt.GetSerialNo(),sid,hi,uw)
			rmr.AddMSG(fk,msg)
			msg = rmr.GetMsg(fk)
		}

		rcv:=msg.GetRecv()
		if rcv.GetKey() == nil{
			rcv.SetKey(fk)
		}

		if ack,_ := rcv.Write(pkt); ack!=nil{
			sendAck(rcv.GetUdpReadWriter(),ack,pkt)
		}

		if rcv.Finish() {
			handleMsg(msgid,hi,rcv)
		}

		rmr.PutMsg(fk)
	}

	p.runningRecv =false

	return err

}


func (p *peer)Close()  {
	//try to close socket

	p.runningSend = false

}

func (p *peer)Dial()  {

}


func (p *peer)AddIPAddr(ip string,port uint16)  {
	p.addrs.AddIP4(ip,port)
}

func (p *peer)DelIpAddr(ip string,port uint16)  {
	p.addrs.DelIP4(ip,port)
}

func (p *peer)GetNet() netcommon.UdpReaderWriterer  {
	return p.net
}

func (p *peer)SetNet(net netcommon.UdpReaderWriterer)  {
	p.net = net
}


func (p *peer)SendAsync(msgid int32,headinfo []byte,data []byte, rcvSn uint64) (*chan int,uint64, error)  {
	rs:=netcommon.NewReadSeeker(data)

	pch,sn:= p.send(msgid,headinfo,rs,rcvSn,0)

	return pch,sn,nil

}

func (p *peer)send(msgid int32,headinfo []byte,rs io.ReadSeeker,rcvSn uint64,ms int) (*chan int,uint64) {
	bd:=send.NewBlockData(rs)
	bd.SetRcvSn(rcvSn)
	bd.SetWriter(p.net)
	bd.SetDeadTime(ms)
	ch := make(chan int,0)
	bd.SetSendResultChan(&ch)
	inn:=nbsid.GetLocalId()
	bd.SetTransInfoOrigin(inn.String(),msgid,headinfo)
	bd.SetDataTyp(constant.DATA_TRANSER)
	p.data2Send <- bd

	return bd.GetSendResultChan(),bd.GetSerialNo()
}

func (p *peer)SendLargeDataAsync(msgid int32,headinfo []byte,rs io.ReadSeeker, rcvSn uint64)(*chan int,uint64, error) {
	pch,sn:= p.send(msgid,headinfo,rs,rcvSn,0)

	return pch,sn,nil
}


func (p *peer)SendSync(msgid int32, headinfo []byte,data []byte, rcvSn uint64) (uint64,error){
	rs:=netcommon.NewReadSeeker(data)
	pch,sn:=p.send(msgid,headinfo,rs,rcvSn,0)

	return sn,p.Wait(pch)
}

func (p *peer)SendSyncTimeOut(msgid int32,headinfo []byte,data []byte, rcvSn uint64, ms int) (uint64,error)  {
	rs:=netcommon.NewReadSeeker(data)
	pch,sn:=p.send(msgid,headinfo,rs,rcvSn,ms)
	return sn,p.Wait(pch)
}

func (p *peer)Wait(ch *chan int) error  {
	code:=<-*ch

	var err error

	switch code {
	case 1:
		err=nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"TimeOut"}
	case 2:
		err=nbserr.NbsErr{ErrId:nbserr.UDP_SND_WRITER_IO_ERR,Errmsg:"Write Error"}
	case 0,3:
		err = nil
		//nothing todo...
	}

	return err
}



