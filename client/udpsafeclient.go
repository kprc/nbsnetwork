package client

import (
	"fmt"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/message"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/common/regcenter"
	"github.com/kprc/nbsnetwork/recv"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"net"
	"github.com/kprc/nbsnetwork/rw"
)

type udpClient struct {
	dialAddr address.UdpAddresser
	localAddr address.UdpAddresser
	realAddr address.UdpAddresser
	uw rw.UdpReaderWriterer

	processWait chan int
}


type UdpClient interface {
	Send(headinfo []byte,msgid int32,stationId string,r io.ReadSeeker) error
	SendBytes(headinfo []byte,msgid int32,stationId string,data []byte) error
	Rcv() error
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

	uc.uw = rw.NewReaderWriter(ra,conn)

	return nil
}

func (uc *udpClient)ReDial() error  {
	return uc.Dial()
}

func (uc *udpClient)SendBytes(headinfo []byte,msgid int32,stationId string,data []byte) error  {
	rs := rw.NewReadSeeker(data)
	return uc.Send(headinfo,msgid,stationId,rs)
}

func (uc *udpClient)Send(headinfo []byte,msgid int32,stationId string,r io.ReadSeeker) error  {
	bd := send.NewBlockData(r,constant.UDP_MTU)
	bd.SetWriter(uc.uw)
	bd.SetTransInfoOrigin(stationId,msgid,headinfo)

	go uc.Rcv()

	bd.Send()

	return nil
}

func doAck(pkt packet.UdpPacketDataer)  {
	bs:= send.GetInstance()

	sd := bs.GetBlockDataer(pkt.GetSerialNo())
	if sd == nil {
		return
	}

	bdr := sd.GetBlockData()
	bdr.PushResult(pkt.GetData())
	bs.PutBlockDataer(pkt.GetSerialNo())

}

func sendAck(w io.Writer, ack packet.UdpResulter, pkt packet.UdpPacketDataer){
	upd := packet.NewUdpPacketData()
	upd.SetSerialNo(pkt.GetSerialNo())
	upd.SetACK()
	back,_ := ack.Serialize()
	upd.SetLength(int32(len(back)))
	upd.SetData(back)

	w.Write(back)

}


func (uc *udpClient)Rcv() error  {

	for  {
		buf := make([]byte,1024)
		n,err := uc.uw.Read(buf)
		if err!=nil {
			return err
		}
		pkt:=packet.NewUdpPacketData()
		if err = pkt.DeSerialize(buf[:n]);err!=nil{
			fmt.Println("Packet DeSerialize failed!",err)
			continue
		}
		//if type is ack, send to result channel and continue
		if pkt.GetTyp() == constant.ACK {
			//send to send block
			doAck(pkt)
			continue
		}
		mc := regcenter.GetMsgCenterInstance()

		msgid,stationId,headinfo := mc.GetMsgId(pkt.GetTransInfo())
		if msgid == constant.MSG_NONE || stationId == "" {
			fmt.Printf("Receive msgid %d, stationId %s",msgid,stationId)
			continue
		}

		rmr := message.GetInstance()
		k := message.NewMsgKey(pkt.GetSerialNo(),stationId)
		m := rmr.GetMsg(k)
		var rcv recv.RcvDataer
		if m==nil {
			m = message.NewRcvMsg()
			handler := mc.GetHandler(msgid)
			m.SetWS(handler.GetWSNew()(headinfo))
			//uw:=send.NewReaderWriter(remoteAddr,sock)
			m.SetUW(uc.uw)
			m.SetKey(k)
			rcv = recv.NewRcvDataer(pkt.GetSerialNo(),m.GetWS())
			m.SetRecv(rcv)

			rmr.AddMSG(k,m)

			m = rmr.GetMsg(k)   //just to increase the ref count
		}else {
			rcv = m.GetRecv()
		}

		ack,err := rcv.Write(pkt)
		if ack != nil{
			sendAck(m.GetUW(),ack,pkt)
		}
		if rcv.Finish() {
			mc.GetHandler(msgid).GetHandler()(headinfo,m.GetWS(),m.GetUW())
			rmr.PutMsg(k)
			rmr.DelMsg(k)

			break
		}else {
			rmr.PutMsg(k)
		}

	}

	return nil
}