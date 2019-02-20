package nbspeer

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/client"
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
	"sync/atomic"
)


/*

A client to Dial

     |
     |
     |
     |start
     |
     |
_____V______                           ______________
|          |                           |             |
|   Init   |                           |   Dead      |
|          |                           |             |
------------                           ------^--------
     |                                       |
     |                                       |
     |                                       |
     |                                       |
     |Dial                                   |After attempt dial 3 times
     |                                       |
     |                                       |
     |                                       |
     |                                       |
     |                                       |
     |                                       |
     |                                       |
_____V______                            _____|________
|          |  Ping Error/Send Error     |             |
|Connected |--------------------------->|Disconnected |
|          |<---------------------------|             |
------------    Redial 3 times          ---------------

=================================================================================

A client from server listen
				____________                           ______________
				|          |                           |             |
				|   Init   |                           |   Dead      |
				|          |                           |             |
				---------r--                           ------^--------
					 |    \                                  |
					 |     \                                 |
					 |      \                                |
					 |       \                               |
					 |        \                              |
					 |         \                             |
					 |          \                            |
					 |           \                           |
					 |            \                          |
					 |             \if not dial,             |
					 |              \  then init and dial it |
					 |               \                       |
					 |                \                      |
					 |                 \                     |
					 |                  \                    |
					 |                   \                   |
					 |                    \                  |
					 |                     \                 |
					 |                      \                |
					 |                       \               |
					 |Dial                    \              |After attempt dial 3 times
					 |                         \             |
					 |                          \            |
					 |                           \           |
					 |                            \          |
					 |                             \         |
					 |                              \        |
					 |                               \       |
				_____V______                          \ _____|_________
	start		|          |  Ping Error/Send Error    \|             |
---------------->Connected |--------------------------->|Disconnected |
				|          |<---------------------------|             |
				------------    Redial 3 times          ---------------


*/


var (
	PEER_INIT int=1
	PEER_CONNECTED int=2
	PEER_DISCONNECTED int = 3
	PEER_DEAD int =4


	PEER_IDLE int32 = 0
	PEER_RUNNING int32 = 1

)


type peer struct {
	status int
	dialtimes int //
	isdialconn bool
	runningflag int32
	runningLock sync.Mutex
	addrs address.UdpAddresser
	net netcommon.UdpReaderWriterer
	client client.UdpClient
	stationId string

	runningSend bool
	runningRecv bool

	selfAddr bool
	data2Send chan send.BlockDataer
}


type NbsPeer interface {
	AddIPAddr(ip string,port uint16)
	DelIpAddr(ip string,port uint16)
	Dial() netcommon.UdpReaderWriterer
	GetNet() netcommon.UdpReaderWriterer
	SetNet(net netcommon.UdpReaderWriterer)
	SendAsync(msgid int32,headinfo []byte,data []byte, rcvSn uint64) (*chan int, uint64,  error)
	SendLargeDataAsync(msgid int32,headinfo []byte,rs io.ReadSeeker, rcvSn uint64)(*chan int,uint64, error)
	SendSync(msgid int32, headinfo []byte,data []byte, rcvSn uint64) (uint64,error)
	SendSyncTimeOut(msgid int32,headinfo []byte,data []byte, rcvSn uint64, ms int) (uint64,error)
	Wait(ch *chan int) error
	Run() error
	Close()

}

func NewNbsPeer(sid string) NbsPeer  {
	p:= &peer{stationId:sid,data2Send:make(chan send.BlockDataer,32)}
	atomic.StoreInt32(&p.runningflag,PEER_IDLE)

	return p
}



func (p *peer)sendbd(bd send.BlockDataer) error {

	sd:=send.NewStoreData(bd)
	bs:=send.GetBSInstance()
	bs.AddBlockDataer(bd.GetSerialNo(),sd)
	sd = bs.GetBlockDataer(bd.GetSerialNo())
	err:= bd.Send()

	bs.PutBlockDataer(bd.GetSerialNo())

	bs.DelBlockDataer(bd.GetSerialNo())

	return err
}

var errPeerIsRunning = nbserr.NbsErr{ErrId:nbserr.PEER_RUNNING,Errmsg:"Peer is running"}
var errPeerIsDead = nbserr.NbsErr{ErrId:nbserr.PEER_DEAD,Errmsg:"Peer is dead"}

func (p *peer)isRunning() bool  {
	if atomic.LoadInt32(&p.runningflag) == PEER_RUNNING{
		return true
	}
	return false
}

func (p *peer)peerSetFlag(flag int32)  {
	atomic.StoreInt32(&p.runningflag,flag)
}

func (p *peer)Run() error {

	if p.isRunning() {
		return errPeerIsRunning
	}

	p.runningLock.Lock()
	if p.isRunning() {
		p.runningLock.Unlock()
		return errPeerIsRunning
	}

	p.peerSetFlag(PEER_RUNNING)

	p.runningLock.Unlock()

	var wg sync.WaitGroup
	var err error
	var sendresultchan chan int
	var rcvresultchan chan int

	if p.status == PEER_DEAD {
		err = errPeerIsDead
	}

	if p.status == PEER_INIT {
		p.Dial()
	}

	if p.status == PEER_DISCONNECTED {
		//p.ReDial()
	}

	if p.status == PEER_CONNECTED {
		if p.isdialconn {
			wg.Add(1)
			rcvresultchan = make(chan int,0)
			go p.recv(&wg, &rcvresultchan)
		}
		sendresultchan = make(chan int,0)
		wg.Add(1)
		go p.Sendbd(&wg,&sendresultchan)
	}


	wg.Wait()

	if sendresultchan != nil {
		code := <-sendresultchan
		if code >0 {

		}
	}

	if rcvresultchan != nil{
		code := <-rcvresultchan
		if code > 0{
			
		}

	}

	p.peerSetFlag(PEER_IDLE)

	return err

}


func (p *peer) Sendbd(wg *sync.WaitGroup, r *chan int) error{

	defer (*wg).Done()

	var err error

	for{
		bd:=<-p.data2Send
		if bd == nil {		//if channel closed
			break
		}
		if err =p.sendbd(bd); err!=nil{   //send error
			break
		}
	}

	p.runningSend = false

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

func (p *peer)recv(wg *sync.WaitGroup, r *chan int) error{

	defer (*wg).Done()

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

	gbuf := make([]byte,1024)
	var n int
	var remoteAddr *net.UDPAddr
	var err error

	for{
		buf := gbuf[:]
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

	return err

}


func (p *peer)Close()  {

	p.runningSend = false
	p.runningRecv =false

}

func (p *peer)Dial() netcommon.UdpReaderWriterer  {
	return nil
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

	pch,sn:= p.sendAsync(msgid,headinfo,rs,rcvSn,0)

	return pch,sn,nil

}

func (p *peer)newBDAsync(msgid int32,headinfo []byte,rs io.ReadSeeker,rcvSn uint64,ms int) send.BlockDataer  {
	bd:=send.NewBlockData(rs)
	bd.SetRcvSn(rcvSn)
	bd.SetWriter(p.net)
	bd.SetTV(ms)

	inn:=nbsid.GetLocalId()
	bd.SetTransInfoOrigin(inn.String(),msgid,headinfo)
	bd.SetDataTyp(constant.DATA_TRANSER)

	return bd
}

func (p *peer)newBDSync(msgid int32,headinfo []byte,rs io.ReadSeeker,rcvSn uint64,ms int) send.BlockDataer {
	bd:=p.newBDAsync(msgid,headinfo,rs,rcvSn,ms)

	ch := make(chan int,0)
	bd.SetSendResultChan(&ch)

	return bd
}

func (p *peer)sendAsync(msgid int32,headinfo []byte,rs io.ReadSeeker,rcvSn uint64,ms int) (*chan int,uint64) {

	bd:=p.newBDAsync(msgid,headinfo,rs,rcvSn,ms)

	p.data2Send <- bd

	return bd.GetSendResultChan(),bd.GetSerialNo()
}

func (p *peer)sendSync(msgid int32,headinfo []byte,rs io.ReadSeeker,rcvSn uint64,ms int) (*chan int,uint64) {

	bd:=p.newBDSync(msgid,headinfo,rs,rcvSn,ms)

	p.data2Send <- bd

	return bd.GetSendResultChan(),bd.GetSerialNo()
}

func (p *peer)SendLargeDataAsync(msgid int32,headinfo []byte,rs io.ReadSeeker, rcvSn uint64)(*chan int,uint64, error) {
	pch,sn:= p.sendAsync(msgid,headinfo,rs,rcvSn,0)

	return pch,sn,nil
}


func (p *peer)SendSync(msgid int32, headinfo []byte,data []byte, rcvSn uint64) (uint64,error){
	rs:=netcommon.NewReadSeeker(data)
	pch,sn:=p.sendSync(msgid,headinfo,rs,rcvSn,0)

	return sn,p.Wait(pch)
}

func (p *peer)SendSyncTimeOut(msgid int32,headinfo []byte,data []byte, rcvSn uint64, ms int) (uint64,error)  {
	rs:=netcommon.NewReadSeeker(data)
	pch,sn:=p.sendSync(msgid,headinfo,rs,rcvSn,ms)
	return sn,p.Wait(pch)
}

func (p *peer)clean()  {
	// nothing to do ...

}

func (p *peer)Wait(ch *chan int) error  {
	code:=<-*ch

	var err error

	switch code {
	case send.SEND_TIME_OUT:
		err=nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"TimeOut"}
	case send.SEND_WRITE_ERR:
		err=nbserr.NbsErr{ErrId:nbserr.UDP_SND_WRITER_IO_ERR,Errmsg:"Write Error"}
	case send.SEND_DEFAULT_RESULT,send.SEND_FINISH:
		err = nil
		//nothing todo...
	case send.CLIENT_CLOSED:
		//p.clean()
		err = nbserr.NbsErr{ErrId:nbserr.UDP_SND_CLOSED}
	}


	return err
}



