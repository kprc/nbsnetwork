package nbsnetwork

import (
	"io"
	"sync/atomic"
	"time"
	"sync"
	"github.com/kprc/nbsdht/nbserr"
)

var blocksnderr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_DEFAULT_ERR,Errmsg:"Send error"}
var blocksndmtuerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_MTU_ERR,Errmsg:"mtu is 0"}
var blocksndreaderioerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_READER_IO_ERR,Errmsg:"Reader io error"}
var blocksndwriterioerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_WRITER_IO_ERR,Errmsg:"Writer io error"}


type BlockData struct {
	r io.ReadSeeker
	w io.Writer
	serialNo uint64
	unixSec int64
	mtu   uint32
	noacklen uint32       //
	totalreadlen uint32   //for seek
	maxcache uint32
	curNum uint32
	totalCnt uint32
	totalRetryCnt uint32
	dataType uint16
	chResult chan interface{}
	rwlock sync.RWMutex
	sndData map[uint32]UdpPacketDataer
}



type BlockDataer interface {
	Send() error
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

func (bd *BlockData)send(round uint32) (int,error){
	buf := make([]byte,bd.mtu)

	n,err := bd.r.Read(buf)
	if n > 0 {
		upr := NewUdpPacketData(bd.serialNo,bd.dataType)
		upr.SetData(buf[:n])

		upr.SetTotalCnt(0)
		upr.SetPos(round)
		upr.SetLength(uint32(n))
		atomic.AddUint32(&bd.totalreadlen,uint32(n))

		if err==io.EOF {
			upr.SetTotalCnt(round+1)
		}
		bupr,_ := upr.Serialize()
		nw,err := bd.w.Write(bupr)

		atomic.AddUint32(&bd.totalCnt,1)

		if n!=nw ||  err!=nil{
			//need resend
			bd.enqueue(round,upr)
			return 0,blocksndwriterioerr
		}
		atomic.AddUint32(&bd.noacklen,uint32(n))

		upr.SetTryCnt(1)
		atomic.AddUint32(&bd.totalRetryCnt,1)
		bd.enqueue(round,upr)
	}

	if err == io.EOF {
		return 1,nil
	}else if err!=nil {
		return 0,blocksndreaderioerr
	}

	return 0,nil
}

func (bd *BlockData)nonesend() (uint32,error) {

	var round uint32
	bd.rwlock.RLock()
	defer bd.rwlock.Unlock()

	for _,upr := range bd.sndData{
		bupr,_ := upr.Serialize()
		nw,err := bd.w.Write(bupr)
		if round < upr.GetPos() {
			round = upr.GetPos()
		}
		atomic.AddUint32(&bd.totalCnt,1)

		if len(bupr)!=nw ||  err!=nil{
			//need resend

			return 0,blocksndwriterioerr
		}
		atomic.AddUint32(&bd.noacklen,uint32(upr.GetLength()))

		upr.SetTryCnt(1)
		atomic.AddUint32(&bd.totalRetryCnt,1)
	}

	return round,nil
}

func (bd *BlockData)Send() error {

	ret := 0

	if bd.r == nil || bd.w == nil{
		return blocksnderr
	}

	if bd.mtu == 0 {
		return blocksndmtuerr
	}

	round,err := bd.nonesend()
	if err != nil{
		return err
	}

	if bd.totalreadlen > 0 {
		if _, err := bd.r.Seek(int64(bd.totalreadlen), io.SeekStart); err != nil {
			return blocksndreaderioerr
		}
	}

	for {
		if ret == 0  || atomic.LoadUint32(&bd.noacklen) < bd.maxcache{
			if r,err := bd.send(round); err==nil{
				round++
				ret = r
			}else if err!=nil{
				return err
			}
		}
		select {
		case result := <- bd.chResult:
			bd.doresult(result)
		}
	}

	return nil

}

func (bd *BlockData)doresult(result interface{})  {
	switch v:=result.(type) {
	case udpresult:
		for _,id:=range v.resend{
			//to resend
		}
		for _,id:=range v.rcved{
			//delete
		}
	}
}

func (bd *BlockData)enqueue(pos uint32, data UdpPacketDataer)  {
	bd.rwlock.Lock()
	defer bd.rwlock.Unlock()
	bd.sndData[pos] = data
}
