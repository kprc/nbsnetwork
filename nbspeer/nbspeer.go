package nbspeer

import (
	"github.com/kprc/nbsdht/dht/nbsid"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/constant"
	"github.com/kprc/nbsnetwork/common/queue"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/send"
	"io"
	"sync"
)

type peer struct {
	addrs address.UdpAddresser
	net netcommon.UdpReaderWriterer
	stationId string
	lock sync.Mutex
	data2Send queue.Queue
}


type NbsPeer interface {
	AddIPAddr(ip string,port uint16)
	DelIpAddr(ip string,port uint16)
	GetNet() netcommon.UdpReaderWriterer
	SetNet(net netcommon.UdpReaderWriterer)
	SendAsync(msgid int32,headinfo []byte,data []byte, rcvSn uint64) (sn uint64, err error)
	SendLargeDataAsync(msgid int32,headinfo []byte,rs io.ReadSeeker, rcvSn uint64)(uint64, error)
	SendSync(msgid int32, headinfo []byte,data []byte, rcvSn uint64) error
	SendSyncTimeOut(msgid int32,headinfo []byte,data []byte, rcvSn uint64, ms int) error
	WaitResult(sn uint64) (interface{},error)
	Wait(sn uint64) error
}

func NewNbsPeer(sid string) NbsPeer  {
	return &peer{stationId:sid,data2Send:queue.NewQueue()}
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


func (p *peer)SendAsync(msgid int32,headinfo []byte,data []byte, rcvSn uint64) (uint64, error)  {
	rs:=netcommon.NewReadSeeker(data)

	return p.SendLargeDataAsync(msgid,headinfo,rs,rcvSn)

}

func (p *peer)SendLargeDataAsync(msgid int32,headinfo []byte,rs io.ReadSeeker, rcvSn uint64)(uint64, error) {
	bd:=send.NewBlockData(rs)
	bd.SetRcvSn(rcvSn)
	bd.SetWriter(p.net)
	inn:=nbsid.GetLocalId()
	bd.SetTransInfoOrigin(inn.String(),msgid,headinfo)
	bd.SetDataTyp(constant.DATA_TRANSER)
	p.lock.Lock()
	p.data2Send.EnQueueValue(bd)
	p.lock.Unlock()
	return 0,nil
}


func (p *peer)SendSync(msgid int32, headinfo []byte,data []byte, rcvSn uint64) error{
	rs:=netcommon.NewReadSeeker(data)
	p.SendLargeDataAsync(msgid,headinfo,rs,rcvSn)



	//begin wait

	return nil
}

func (p *peer)SendSyncTimeOut(msgid int32,headinfo []byte,data []byte, rcvSn uint64, ms int) error  {


	return nil
}

func (p *peer)WaitResult(sn uint64) (interface{},error)  {
	return nil,nil
}

func (p *peer)Wait(sn uint64) error  {
	return nil
}



