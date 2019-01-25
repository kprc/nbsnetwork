package nbspeer

import (
	"github.com/gogo/protobuf/io"
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/common/queue"
	"github.com/kprc/nbsnetwork/netcommon"
)

type peer struct {
	addrs address.UdpAddresser
	net netcommon.UdpReaderWriterer
	stationId string
	data2Send queue.Queue
}


type NbsPeer interface {
	AddIPAddr(ip string,port uint16)
	DelIpAddr(ip string,port uint16)
	GetNet() netcommon.UdpReaderWriterer
	SetNet(net netcommon.UdpReaderWriterer)
	SendAsync(msgid int,headinfo []byte,data []byte, rcvSn uint64) (sn uint64, err error)
	SendLargeDataAsync(msgid int,headinfo []byte,r io.Reader, rcvSn uint64)(uint64, error)
	SendSync(msgid int, headinfo []byte,data []byte, rcvSn uint64) error
	SendSyncTimeOut(msgid int,headinfo []byte,data []byte, rcvSn uint64, ms int) error
	WaitResult(sn uint64) (interface{},error)
	Wait(sn uint64) error
}

func NewNbsPeer(sid string) NbsPeer  {
	return &peer{stationId:sid}
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


func (p *peer)SendAsync(msgid int,headinfo []byte,data []byte, rcvSn uint64) (uint64, error)  {

	return 0,nil
}

func (p *peer)SendLargeDataAsync(msgid int,headinfo []byte,r io.Reader, rcvSn uint64)(uint64, error) {
	return 0,nil
}


func (p *peer)SendSync(msgid int, headinfo []byte,data []byte, rcvSn uint64) error{
	return nil
}

func (p *peer)SendSyncTimeOut(msgid int,headinfo []byte,data []byte, rcvSn uint64, ms int) error  {
	return nil
}

func (p *peer)WaitResult(sn uint64) (interface{},error)  {
	return nil,nil
}

func (p *peer)Wait(sn uint64) error  {
	return nil
}



