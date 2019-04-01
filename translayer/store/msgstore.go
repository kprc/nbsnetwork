package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"github.com/kprc/nbsnetwork/common/list"
	"github.com/kprc/nbsnetwork/tools"
	"sync"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/translayer/message"
)

type block struct {
	blk interface{}
	timeout int32
	step int32
	firstsendtime int64
	lastsendtime int64
	cursendtimes int32
	resendtimes int32
	ackflag bool
	uid string
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
	AddBlock(blk interface{},uid string)
	AddBlockWithParam(data interface{},uid string,tv int32,resndtimes int32,step int32)
	DelBlock(blk interface{})
	FindBlockDo(v interface{},arg interface{},do list.FDo) (r interface{},err error)
	//TimeOut(arg interface{},del list.FDel)
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

func GetBlkAndRefresh(data interface{}) interface{}  {
	if blk,ok:=data.(*block); !ok{
		return nil
	}else{
		blk.lastsendtime = tools.GetNowMsTime()
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


func (bs *blockstore)AddBlock(data interface{},uid string) {
	blk:=&block{blk:data}
	blk.sn = data.(BlockInter).GetSn()
	blk.timeout = 5000  //ms
	blk.resendtimes = 2
	blk.step = 1
	blk.lastsendtime = tools.GetNowMsTime()
	blk.firstsendtime = blk.lastsendtime
	blk.cursendtimes = 1

	bs.Add(blk)
}

func (bs *blockstore)AddBlockWithParam(data interface{},uid string,tv int32,resndtimes int32,step int32){
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
	blk.firstsendtime = blk.lastsendtime
	blk.cursendtimes = 1

	bs.Add(blk)
}

//TODO...

func (bs *blockstore)DelBlock(blk interface{}) {

	bs.Del(blk)
}

func (bs *blockstore)FindBlockDo(v interface{},arg interface{},do list.FDo) (r interface{},err error)  {
	return bs.FindDo(v,arg,do)
}


func (bs *blockstore)doTimeOut()  {

	type blk2del struct{
		arrdel []*block
	}

	arr2del := &blk2del{arrdel:make([]*block,0)}

	fdo:= func(arg interface{}, v interface{}) (ret interface{},err error){
		blk:=v.(block)
		if blk.ackflag {
			return
		}

		l:=arg.(*blk2del)
		curtime := tools.GetNowMsTime()
		tv:=curtime-blk.firstsendtime
		if tv > int64(blk.timeout){
			l.arrdel = append(l.arrdel,v.(*block))
			return
		}

		if blk.cursendtimes >= blk.resendtimes{
			l.arrdel = append(l.arrdel,v.(*block))
			return
		}

		if curtime - blk.lastsendtime > int64(blk.step){

			cs:=netcommon.GetConnStoreInstance()
			conn:=cs.GetConn(blk.uid)
			if conn == nil || blk.blk == nil{
				return
			}
			message.SendUm(blk.blk.(UdpMsg),conn)

			blk.lastsendtime = curtime
			blk.cursendtimes ++
		}

		return
	}

	bs.TraversAll(arr2del,fdo)

	for _,blk:=range arr2del.arrdel{
		bs.DelBlock(blk)
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






