package rpc

import (
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/netcommon"
)

type m2mrpc struct {
	message.ReliableMsg
	response chan interface{}
}


type M2MRpc interface {
	message.ReliableMsg
	RPCSend(data []byte) (r interface{},err error)
}

func NewM2MRpc(conn netcommon.UdpConn,apptyp uint32,arg interface{},do RpcDo) M2MRpc {
	mr:=message.NewReliableMsgWithDelayInit(conn,false)

	mr.SetAppTyp(apptyp)

	mmr:=&m2mrpc{ReliableMsg:mr,response:make(chan interface{},1)}

	rb:=NewRpcBlock()
	rb.SetResponseChan(&mmr.response)
	rb.SetSn(mmr.GetSn())
	rb.SetData(arg)
	rb.SetRpcDo(do)

	rs:=GetRpcStore()
	rs.AddRpcBlock(rb)

	return mmr
}

func (mmr *m2mrpc)RPCSend(data []byte)(r interface{},err error)  {
	//sn:= mmr.GetSn()

	//rb:=NewRpcBlock()
	//rb.SetSn(sn)
	//rb.SetResponseChan(&mmr.response)
	if err=mmr.ReliableSend(data);err!=nil{
		return nil,err
	}

	<-mmr.response



	return
}