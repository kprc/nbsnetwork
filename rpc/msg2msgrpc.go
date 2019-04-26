package rpc

import (
	"github.com/kprc/nbsnetwork/translayer/message"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsdht/nbserr"
)

type m2mrpc struct {
	message.ReliableMsg
	response chan interface{}
}


type M2MRpc interface {
	message.ReliableMsg
	RPCSend(data []byte) (r interface{},err error)
}

var (
	rpctimeouterr=nbserr.NbsErr{ErrId:nbserr.RPC_TIMEOUT_ERR,Errmsg:"rpc invoke timeout error"}
)


func NewM2MRpc(conn netcommon.UdpConn,apptyp uint32,arg interface{},do RpcDo) M2MRpc {
	mr:=message.NewReliableMsgWithDelayInit(conn,false)

	mr.SetAppTyp(apptyp)

	mmr:=&m2mrpc{ReliableMsg:mr,response:make(chan interface{},1)}

	rb:=NewRpcBlock()
	rb.SetResponseChan(&mmr.response)
	rb.SetSn(mmr.GetSn())
	rb.SetData(arg)

	if do ==  nil{
		do = DefaultRpcDo
	}

	rb.SetRpcDo(do)

	rs:=GetRpcStore()
	rs.AddRpcBlock(rb)

	return mmr
}

func (mmr *m2mrpc)RPCSend(data []byte)(r interface{},err error)  {

	if err=mmr.ReliableSend(data);err!=nil{
		return nil,err
	}

	r=<-mmr.response

	switch r.(type) {
	case int64:
		return nil,rpctimeouterr
	default:
		return
	}
}