package rpc

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"sync"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/translayer/store"
	"github.com/kprc/nbsnetwork/tools"
	"github.com/kprc/nbsnetwork/common/constant"
)

type rpcblockdesc struct {
	blk interface{}
	timeoutInterval int32
	lastAccessTime int64
	sn uint64
}


type rpcstore struct {
	hashlist.HashList
	tick chan int64
	quit chan int64
	wg *sync.WaitGroup
}


type RpcStore interface {
	AddRpcBlock(blk interface{})
	AddRpcBlockWithParam(blk interface{},timeInterval int32)
	DelRpcBlock(blk interface{})
	FindRpcBlockDo(blk interface{},arg interface{},do list.FDo) (r interface{},err error)
	Run()
	Stop()
}

var (
	rpcstoreInst RpcStore
	rpcstoreInstLock sync.Mutex
	lastDoTimeOut int64
	timeouttv int64 = 1000
)

func (rbd *rpcblockdesc)GetSn() uint64  {
	return  rbd.sn
}

func GetRpcStore() RpcStore  {
	if rpcstoreInst == nil {
		rpcstoreInstLock.Lock()
		defer rpcstoreInstLock.Unlock()

		if rpcstoreInst == nil{
			rpcstoreInst = newRpcStore()
		}
	}

	return rpcstoreInst
}

func newRpcStore() RpcStore  {
	hl:=hashlist.NewHashList(0x80,store.FBlockHash,store.FBlockEquals)
	rs:=&rpcstore{HashList:hl}
	rs.tick = make(chan int64,64)
	rs.quit = make(chan int64,1)
	rs.wg = &sync.WaitGroup{}

	t:=tools.GetNbsTickerInstance()
	t.Reg(&rs.tick)

	return rs

}

func (rs *rpcstore)addBlock(data interface{},timeinterval int32)  {
	rbd:=&rpcblockdesc{}
	rbd.sn = data.(store.BlockInter).GetSn()
	rbd.timeoutInterval = timeinterval
	rbd.lastAccessTime = tools.GetNowMsTime()

	rs.Add(rbd)

}

func (rs *rpcstore)AddRpcBlock(blk interface{})  {
	rs.addBlock(blk,int32(constant.RPC_STORE_TIMEOUT))
}

func (rs *rpcstore)AddRpcBlockWithParam(blk interface{},timeInterval int32)  {
	if timeInterval == 0{
		timeInterval = int32(constant.RPC_STORE_TIMEOUT)
	}
	rs.addBlock(blk,timeInterval)

}

func (rs *rpcstore)DelRpcBlock(blk interface{})  {
	rs.Del(blk)
}

func (rs *rpcstore)FindRpcBlockDo(blk interface{},arg interface{},do list.FDo) (r interface{},err error)  {
	return rs.FindDo(blk,arg,do)
}

func (rs *rpcstore)doTimeOut()  {

}

func (rs *rpcstore)Run()  {
	if rs.wg != nil{
		rs.wg.Add(1)
	}
	defer func() {
		if rs.wg !=nil{
			rs.wg.Done()
		}
	}()

	for{
		select {
			case <-rs.tick:
				if tools.GetNowMsTime() - lastDoTimeOut > timeouttv {
					lastDoTimeOut = tools.GetNowMsTime()
					rs.doTimeOut()
				}
			case <-rs.quit:
				return
		}
	}

}

func (rs *rpcstore)Stop()  {
	rs.quit <- 1
	if rs.wg != nil{
		rs.wg.Wait()
	}
}