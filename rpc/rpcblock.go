package rpc

import "github.com/kprc/nbsnetwork/translayer/store"

type RpcDo func(data interface{}, arg interface{}, res *chan interface{}, isTimeOut bool)

type rpcblock struct {
	sn       uint64
	data     interface{}
	response *chan interface{}
	do       RpcDo
}

type RpcBlock interface {
	store.BlockInter
	SetSn(sn uint64)
	SetData(v interface{})
	GetData() interface{}
	GetResponseChan() *chan interface{}
	SetResponseChan(c *chan interface{})
	SetRpcDo(do RpcDo)
	GetRpcDo() RpcDo
}

func NewRpcBlock() RpcBlock {
	return &rpcblock{}
}

func (rb *rpcblock) GetSn() uint64 {
	return rb.sn
}

func (rb *rpcblock) SetSn(sn uint64) {
	rb.sn = sn
}

func (rb *rpcblock) SetData(v interface{}) {
	rb.data = v
}

func (rb *rpcblock) GetData() interface{} {
	return rb.data
}

func (rb *rpcblock) GetResponseChan() *chan interface{} {
	return rb.response
}

func (rb *rpcblock) SetResponseChan(c *chan interface{}) {
	rb.response = c
}

func (rb *rpcblock) SetRpcDo(do RpcDo) {
	rb.do = do
}

func (rb *rpcblock) GetRpcDo() RpcDo {
	return rb.do
}

func RpcBlockDo(data interface{}, arg interface{}, isTimeOut bool) {
	rb := data.(RpcBlock)

	rb.GetRpcDo()(rb.GetData(), arg, rb.GetResponseChan(), isTimeOut)
}

func DefaultRpcDo(data interface{}, arg interface{}, res *chan interface{}, isTimeOut bool) {
	if isTimeOut {
		if res != nil {
			*res <- int64(1)
			return
		}
	}

	if res != nil {
		*res <- arg
	}
	return
}
