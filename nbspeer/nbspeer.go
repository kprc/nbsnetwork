package nbspeer

import (
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/send"
	"io"
)

type peer struct {
	addrs address.UdpAddresser
	net netcommon.UdpReaderWriterer
	stationId string

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
	Run() error
}

func NewNbsPeer(sid string) NbsPeer  {
	return &peer{stationId:sid,data2Send:make(chan send.BlockDataer,32)}
}


func (p *peer) Run() error{

	for{
		bd:=<-p.data2Send




	}


	return nil
}

func (p *peer)recv() {

}


func (p *peer)Close()  {

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



