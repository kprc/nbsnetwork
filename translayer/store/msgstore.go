package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type block struct {
	blk interface{}
	timeout int32
	step int32
	lastsendtime int64
	resendtimes int32
	sn uint64     // one message one code, one message include many packet
}


type BlockInter interface {
	GetSn() uint64
}


type blockstore struct {
	hashlist.HashList
	tick chan int64
}

type BlockStore interface {
	AddBlock(blk interface{})
	AddBlockWithParam(data interface{},tv int32,resndtimes int32,step int32)
	DelBlock(blk interface{})
	FindBlockDo(v interface{},arg interface{},do hashlist.FDo) (r interface{},err error)
	TimeOut(arg interface{},del list.FDel)
	Run()
}


var (
	instancelock sync.Mutex
	instance BlockStore
	lasttimeout int64
	timeouttv int64 = 1000 //ms
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

func GetBlockStoreInstance()  BlockStore {

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


func (bs *blockstore)AddBlock(data interface{}) {
	blk:=&block{blk:data}
	blk.sn = data.(BlockInter).GetSn()
	blk.timeout = 5000  //ms
	blk.resendtimes = 2
	blk.step = 1
	blk.lastsendtime = tools.GetNowMsTime()

	bs.Add(blk)
}

func (bs *blockstore)AddBlockWithParam(data interface{},tv int32,resndtimes int32,step int32){
	blk:=block{blk:data}
	blk.sn = data.(BlockInter).GetSn()
	blk.timeout = 5000  //ms
	blk.resendtimes = 3
	blk.step = 1
	if tv > 0 {
		blk.timeout = tv
	}
	if resndtimes > 0 {
		blk.resendtimes = resndtimes
	}
	if step>0{
		blk.step = step
	}
	blk.lastsendtime = tools.GetNowMsTime()


	bs.Add(blk)
}

//TODO...

func (bs *blockstore)DelBlock(blk interface{}) {

	bs.Del(blk)
}

func (bs *blockstore)FindBlockDo(v interface{},arg interface{},do hashlist.FDo) (r interface{},err error)  {
	return bs.FindDo(v,arg,do)
}

func (bs *blockstore)TimeOut(arg interface{},del list.FDel)  {
	bs.TraversDel(arg,del)
}

func (blk *block)Resend()  {
	
}

func (blk *block)TimeOut()  {
	if tools.GetNowMsTime() - blk.lastsendtime > int64(blk.timeout){
		if blk.resendtimes > 0{
			blk.Resend()
			blk.resendtimes --
		}
		blk.lastsendtime = tools.GetNowMsTime()
	}
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






