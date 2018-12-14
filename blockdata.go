package nbsnetwork

import (
	"bytes"
	"io"
	"sync/atomic"
	"time"
)


type BlockData struct {
	r io.Reader
	serialNo uint64
	unixSec int64
	mtu   uint32
	curNum uint32
	totalRetryCnt uint32
	sndData map[uint32]*bytes.Buffer
}

type BlockDataer interface {

}

var gSerialNo uint64 = UDP_SERIAL_MAGIC_NUM

func (uh *BlockData)nextSerialNo() {
	uh.serialNo = atomic.AddUint64(&gSerialNo,1)
}

func NewBlockData(r io.Reader,mtu uint32) BlockDataer {
	uh := &BlockData{r:r,mtu:mtu}
	uh.nextSerialNo()
	uh.unixSec = time.Now().Unix()
	uh.mtu = UDP_MTU
	return uh
}

