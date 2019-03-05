package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type block struct {
	blk interface{}
	timeout int
	lastsendtime int64
	sn uint64     // one message one code, one message include many packet
}

type blockstore struct {
	hashlist.HashList
	tick chan int64
}

type BlockStore interface {
	AddBlock(blk interface{})
	DelBlock(blk interface{})
	FindBlockDo(v interface{},arg interface{},do hashlist.FDo) (r interface{},err error)
	TimeOut(arg interface{},del list.FDel)
	Run()
}


var (
	instancelock sync.Mutex
	instance BlockStore
	lasttimeout int64
	timeouttv int64 = 5000 //ms
)

var fhash = func(v interface{}) uint {
	blk:=v.(block)

	return uint(blk.sn&0x3FF)
}

var fequals = func(v1 interface{},v2 interface{}) int{
	blk1:=v1.(block)
	blk2:=v2.(block)

	if blk1.sn == blk2.sn {
		return 0
	}

	return 1
}

var fdel = func(arg interface{},v interface{}) bool{
	return false
}

func NewBlockStoreInstance()  BlockStore {

	if instance == nil{
		instancelock.Lock()
		defer instancelock.Unlock()

		if instance == nil{
			instance = newBlockStore()
		}
	}


	return instance
}

func newBlockStore() BlockStore {

	bs:=hashlist.NewHashList(0x400,fhash,fequals).(*blockstore)

	bs.tick = make(chan int64,64)

	t := tools.GetNbsTickerInstance()

	t.Reg(&bs.tick)

	lasttimeout = tools.GetNowMsTime()

	return bs
}


func (bs *blockstore)AddBlock(blk interface{}) {
	bs.Add(blk)
}

func (bs *blockstore)DelBlock(blk interface{}) {
	bs.Del(blk)
}

func (bs *blockstore)FindBlockDo(v interface{},arg interface{},do hashlist.FDo) (r interface{},err error)  {
	return bs.FindDo(v,arg,do)
}

func (bs *blockstore)TimeOut(arg interface{},del list.FDel)  {
	bs.TraversDel(arg,del)
}


func (bs *blockstore)Run()  {
	select {
	case <-bs.tick:
		if tools.GetNowMsTime() - lasttimeout > timeouttv{
			lasttimeout = tools.GetNowMsTime()
			bs.TimeOut(nil,fdel)
		}
	}
}






