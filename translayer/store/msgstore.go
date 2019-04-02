package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
)

type block struct {
	blk interface{}
	timeoutInterval int32
	lastAccessTime int64
	ackflag bool
	sn uint64
}


type BlockInter interface {
	GetSn() uint64
}


type blockstore struct {
	hashlist.HashList
	tick chan int64
}

type BlockStore interface {
	AddMessage(blk interface{},uid string)
	AddMessageWithParam(data interface{},timeInterval int32)
	DelMessage(blk interface{})
	FindMessageDo(v interface{},arg interface{},do list.FDo) (r interface{},err error)
	Run()
}


var (
	instancelock sync.Mutex
	instance BlockStore
	lasttimeout int64
	timeouttv int64 = 1000 //ms
)

var fhash = func(v interface{}) uint {
	blk:=v.(BlockInter)

	return uint(blk.GetSn()&0x3FF)
}

var fequals = func(v1 interface{},v2 interface{}) int{
	blk1:=v1.(BlockInter)
	blk2:=v2.(BlockInter)

	if blk1.GetSn() == blk2.GetSn() {
		return 0
	}

	return 1
}

var fdel = func(arg interface{},v interface{}) bool{
	return false
}

func SetAckFlag(data interface{}) interface{}  {
	if blk,ok:=data.(*block); !ok{
		return nil
	}else{
		blk.ackflag = true
		return blk.blk
	}
}

func (blk *block)GetSn() uint64  {
	return blk.sn
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

	bs:=hashlist.NewHashList(0x80,fhash,fequals).(*blockstore)

	bs.tick = make(chan int64,64)

	t := tools.GetNbsTickerInstance()

	t.Reg(&bs.tick)

	lasttimeout = tools.GetNowMsTime()

	return bs
}


func (bs *blockstore)AddMessage(data interface{},uid string) {
	blk:=&block{blk:data}
	blk.sn = data.(BlockInter).GetSn()
	blk.timeoutInterval = 5000  //ms
	blk.lastAccessTime = tools.GetNowMsTime()

	bs.Add(blk)
}

func (bs *blockstore)AddMessageWithParam(data interface{},timeInterval int32){
	blk:=block{blk:data}
	blk.sn = data.(BlockInter).GetSn()
	blk.timeoutInterval = timeInterval  //ms

	blk.lastAccessTime = tools.GetNowMsTime()

	bs.Add(blk)
}

//TODO...

func (bs *blockstore)DelMessage(blk interface{}) {

	bs.Del(blk)
}

func (bs *blockstore)FindMessageDo(v interface{},arg interface{},do list.FDo) (r interface{},err error)  {
	return bs.FindDo(v,arg,do)
}


func (bs *blockstore)doTimeOut()  {

	type blk2del struct{
		arrdel []*block
	}

	arr2del := &blk2del{arrdel:make([]*block,0)}

	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error){
		blk:=v.(block)

		l:=arg.(*blk2del)
		curtime := tools.GetNowMsTime()
		tv:=curtime-blk.lastAccessTime
		if tv > int64(blk.timeoutInterval){
			um:=blk.blk.(UdpMsg)
			um.Inform(UDP_INFORM_TIMEOUT)
			l.arrdel = append(l.arrdel,v.(*block))
		}

		return
	}

	bs.TraversAll(arr2del,fdo)

	for _,blk:=range arr2del.arrdel{
		bs.DelMessage(blk)
	}

}


func (bs *blockstore)Run()  {
	select {
	case <-bs.tick:
		if tools.GetNowMsTime() - lasttimeout > timeouttv{
			lasttimeout = tools.GetNowMsTime()
			bs.doTimeOut()
		}
	}
}






