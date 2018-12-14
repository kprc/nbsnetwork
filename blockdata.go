package nbsnetwork

import (
	"io"
	"sync/atomic"
	"time"
	"sync"
	"github.com/kprc/nbsdht/nbserr"
)

var blocksnderr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_READER_ERR,Errmsg:"Reader is nil or Writer is nil"}
var blocksndmtuerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_MTU_ERR,Errmsg:"mtu is 0"}

type BlockData struct {
	timeout uint32   //second
	r io.ReadSeeker
	w io.Writer
	serialNo uint64
	unixSec int64
	mtu   uint32
	sndLen uint32
	maxcache uint32
	curNum uint32
	totalRetryCnt uint32
	rwlock sync.RWMutex
	sndData map[uint32]uint32
}



type BlockDataer interface {

}

var gSerialNo uint64 = UDP_SERIAL_MAGIC_NUM

func (uh *BlockData)nextSerialNo() {
	uh.serialNo = atomic.AddUint64(&gSerialNo,1)
}

func NewBlockData(r io.ReadSeeker,mtu uint32) BlockDataer {
	uh := &BlockData{r:r,mtu:mtu}
	uh.nextSerialNo()
	uh.unixSec = time.Now().Unix()
	uh.mtu = UDP_MTU
	uh.maxcache = UDP_MAX_CACHE
	return uh
}

func (bd *BlockData)Send() error {
	if bd.r == nil || bd.w == nil{
		return blocksnderr
	}

	if bd.mtu == 0 {
		return blocksndmtuerr
	}



	for {
		buf := make([]byte,bd.mtu)

		n,err := bd.r.Read(buf)
		if n > 0 && err==nil {

		}
		//select {
		//case
		//}
	}

	return nil

}