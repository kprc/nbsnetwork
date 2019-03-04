package store

import (
	"github.com/kprc/nbsnetwork/common/hashlist"
	"sync"
)

type block struct {
	blk interface{}
	timeout int
	sendtime int64
	resndCnt int
	sn uint64     // one message one code, one message include many packet
	pos uint64     // one packet one code, small then message
}

type blockstore struct {
	hashlist.HashList
}

type BlockStore interface {
	AddBlock(blk interface{})
	DelBlock(blk interface{})
	FindDo(arg interface{},do hashlist.FDo) (v interface{},err error)
}


var (
	instancelock sync.Mutex
	instance BlockStore
)

var fhash = func(v interface{}) uint {
	blk:=v.(block)

	h:= blk.sn >> 4 + blk.pos << 4

	return uint(h&0x3FF)
}

var fequals = func(v1 interface{},v2 interface{}) int{
	blk1:=v1.(block)
	blk2:=v2.(block)

	if blk1.sn == blk2.sn && blk1.pos == blk2.pos{
		return 0
	}

	return 1
}

func NewBlockStoreInstance()  BlockStore {

	if instance == nil{
		instancelock.Lock()
		defer instancelock.Unlock()

		if instance == nil{
			instance = hashlist.NewHashList(0x400,fhash,fequals).(BlockStore)
		}
	}

	return instance
}

func (bs *blockstore)AddBlock(blk interface{}) {
	bs.Add(blk)
}

func (bs *blockstore)DelBlock(blk interface{}) {
	bs.Del(blk)
}

func (bs *blockstore)FindDo(arg interface{},do hashlist.FDo) (v interface{},err error)  {
	return bs.FindDo(arg,do)
}



