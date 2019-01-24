package nbspeer

import (
	"github.com/kprc/nbsnetwork/common/address"
	"github.com/kprc/nbsnetwork/netcommon"
)

type peer struct {
	address.UdpAddresser
	netcommon.UdpReaderWriterer
	stationId string
}


type NbsPeer interface {
	SendAsync(msgid int,headinfo []byte,data []byte) (sn uint64, err error)
	SendSync(msgid int, headinfo []byte,data []byte) (interface{},error)
	SendSyncTime(msgid int,headinfo []byte,data []byte, ms int) (interface{},error)
	WaitResult(sn uint64) (interface{},error)
}

func NewNbsPeer(sid string) NbsPeer  {
	return &peer{stationId:sid}
}


func (p *peer)AddIPAddr(ip string,port uint16)  {
	p.AddIP4(ip,port)
}

func (p *peer)DelIpAddr(ip string,port uint16)  {
	p.DelIP4(ip,port)
}

func (p *peer)SendAsync(msgid int,headinfo []byte,data []byte) (uint64, error)  {

	return 0,nil
}

func (p *peer)SendSync(msgid int, headinfo []byte,data []byte) (interface{},error)  {
	return nil,nil
}

func (p *peer)SendSyncTime(msgid int,headinfo []byte,data []byte, ms int) (interface{},error)  {
	return nil,nil
}

func (p *peer)WaitResult(sn uint64) (interface{},error)  {
	return nil,nil
}





