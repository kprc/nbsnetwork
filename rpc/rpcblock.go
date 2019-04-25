package rpc

import "github.com/kprc/nbsnetwork/translayer/store"

type RpcDo func(data interface{},isTimeOut bool)

type rpcblock struct {
	sn uint64
	data interface{}
	response *chan interface{}
	do RpcDo
}


type RpcBlock interface {
	store.BlockInter
	SetSn(sn uint64)
	SetData(v interface{})
	GetData() interface{}
	SetRpcDo(do RpcDo)
	GetRpcDo() RpcDo
}

func (rb *rpcblock)GetSn() uint64 {
	return rb.sn
}

func (rb *rpcblock)SetSn(sn uint64)  {
	rb.sn = sn
}

func (rb *rpcblock)SetData(v interface{})  {
	rb.data = v
}

func (rb *rpcblock)GetData() interface{}  {
	return rb.data
}

func (rb *rpcblock)SetRpcDo(do RpcDo)  {
	rb.do = do
}

func (rb *rpcblock)GetRpcDo() RpcDo  {
	return rb.do
}

func RpcBlockDo(data interface{},isTimeOut bool)  {
	rb:=data.(RpcBlock)

	rb.GetRpcDo()(rb.GetData(),isTimeOut)
}


