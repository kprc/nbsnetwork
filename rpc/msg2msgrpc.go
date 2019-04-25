package rpc

import (
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/netcommon"
)

type m2mrpc struct {
	message.ReliableMsg
	response chan interface{}
	do RpcDo
}


type M2MRpc interface {
	message.ReliableMsg
	RPCSend(data []byte,arg interface{}) (r interface{},err error)
}

func NewM2MRpc(conn netcommon.UdpConn) M2MRpc {
	mr:=message.NewReliableMsgWithDelay(conn,false)

	mmr:=&m2mrpc{ReliableMsg:mr,response:make(chan interface{},1)}

	return mmr
}

func (mmr *m2mrpc)RPCSend(data []byte, arg interface{})(r interface{},err error)  {
	return
}