package send

import (
	"io"
	"sync/atomic"
	"time"
	"sync"
	"github.com/kprc/nbsdht/nbserr"
	"github.com/kprc/nbsnetwork/common/packet"
	"github.com/kprc/nbsnetwork/common/constant"
	"fmt"
)

var blocksnderr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_DEFAULT_ERR,Errmsg:"Send error"}
var blocksndmtuerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_MTU_ERR,Errmsg:"mtu is 0"}
var blocksndreaderioerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_READER_IO_ERR,Errmsg:"Reader io error"}
var blocksndwriterioerr = nbserr.NbsErr{ErrId:nbserr.UDP_SND_WRITER_IO_ERR,Errmsg:"Writer io error"}
var blocksndtimeout = nbserr.NbsErr{ErrId:nbserr.UDP_SND_TIMEOUT_ERR,Errmsg:"send timeout"}

type BlockData struct {
	timeout int32
	r io.ReadSeeker
	w io.Writer
	serialNo uint64
	unixSec int64
	mtu   uint32
	noacklen int32       //
	totalreadlen uint32   //for seek
	maxcache uint32
	curNum uint32
	totalCnt uint32
	totalRetryCnt uint32
	dataType uint16
	chResult chan interface{}
	rwlock sync.RWMutex
	sndData map[uint32]packet.UdpPacketDataer
}

type BlockDataer interface {
	Send() error
}

var gSerialNo uint64 = constant.UDP_SERIAL_MAGIC_NUM

func (uh *BlockData)nextSerialNo() {
	uh.serialNo = atomic.AddUint64(&gSerialNo,1)
}

func NewBlockData(r io.ReadSeeker,mtu uint32) BlockDataer {
	uh := &BlockData{r:r,mtu:mtu}
	uh.nextSerialNo()
	uh.unixSec = time.Now().Unix()
	//uh.mtu = constant.UDP_MTU
	uh.maxcache = constant.UDP_MAX_CACHE
	uh.timeout = constant.UDP_SEND_TIMEOUT
	return uh
}

func (bd *BlockData)send(round uint32) (int,error){
	buf := make([]byte,bd.mtu)

	n,err := bd.r.Read(buf)
	if n > 0 {
		upr := packet.NewUdpPacketData(bd.serialNo,bd.dataType)
		upr.SetData(buf[:n])

		upr.SetTotalCnt(0)
		upr.SetPos(round)
		upr.SetLength(int32(n))
		atomic.AddUint32(&bd.totalreadlen,uint32(n))

		if err==io.EOF {
			upr.SetTotalCnt(round+1)
		}
		bupr,_ := upr.Serialize()
		fmt.Println(len(bupr))
		nw,err := bd.w.Write(bupr)

		atomic.AddUint32(&bd.totalCnt,1)

		if n!=nw ||  err!=nil{
			//need resend
			bd.enqueue(round,upr)
			return 0,blocksndwriterioerr
		}
		atomic.AddInt32(&bd.noacklen,int32(n))

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

	var round uint32 = 1
	bd.rwlock.RLock()
	defer bd.rwlock.RUnlock()

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

		atomic.AddInt32(&bd.noacklen,upr.GetLength())
		upr.SetTryCnt(1)
		atomic.AddUint32(&bd.totalRetryCnt,1)
	}

	return round,nil
}

func (bd *BlockData)Send() error {

	ret := 0
	curtime := time.Now().Unix()

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
		if time.Now().Unix() - curtime > int64(bd.timeout) {
			return blocksndtimeout
		}
		if ret == 0  || atomic.LoadInt32(&bd.noacklen) < int32(bd.maxcache){
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
		if ret == 1 {
			bd.rwlock.RLock()
			if len(bd.sndData)==0 {
				bd.rwlock.RUnlock()
				return nil
			}
			bd.rwlock.RUnlock()
		}


	}

	return nil

}


func (bd *BlockData)doresult(result interface{}) error {

	switch v:=result.(type) {
	case packet.UdpResulter:
		bd.rwlock.RLock()
	    resend := v.GetReSend()
		for _,id:=range resend{
			if upr,ok:=bd.sndData[id];ok {
				bupr,_ := upr.Serialize()
				nw,err := bd.w.Write(bupr)

				if len(bupr)!=nw ||  err!=nil{
					//need resend
					bd.rwlock.RUnlock()
					return blocksndwriterioerr
				}
				atomic.AddInt32(&bd.noacklen,int32(upr.GetLength()))

				upr.SetTryCnt(upr.GetTryCnt()+1)
				atomic.AddUint32(&bd.totalRetryCnt,1)
			}
		}
		bd.rwlock.Unlock()
		bd.rwlock.Lock()
		if v.GetRcved() >0 {}
		    if upr,ok:=bd.sndData[v.GetRcved()];ok{
		    	atomic.AddInt32(&bd.noacklen,upr.GetLength()*(-1))
			}
			delete(bd.sndData,v.GetRcved())
		}

		bd.rwlock.Unlock()

	return nil
}

func (bd *BlockData)enqueue(pos uint32, data packet.UdpPacketDataer)  {
	bd.rwlock.Lock()
	defer bd.rwlock.Unlock()
	bd.sndData[pos] = data
}
