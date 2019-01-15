package recv

import (
	"sync/atomic"
)


type rcvmsg struct {
	refcnt int32
	rcv RcvDataer    //for receive
}


type RcvMsg interface {
	SetRefCnt(cnt int32)
	GetRefCnt() int32
	IncRefCnt()
	DecRefCnt() int32
	GetRecv() RcvDataer
	SetRecv(rcv RcvDataer)
}

func NewRcvMsg() RcvMsg  {
	rm := &rcvmsg{}

	return rm
}




func (rm *rcvmsg)SetRefCnt(cnt int32){
	atomic.StoreInt32(&rm.refcnt,cnt)
}

func (rm *rcvmsg)GetRefCnt() int32{
	return atomic.LoadInt32(&rm.refcnt)
}

func (rm *rcvmsg)IncRefCnt()  {
	atomic.AddInt32(&rm.refcnt,1)
}

func (rm *rcvmsg)DecRefCnt() int32  {
	return atomic.AddInt32(&rm.refcnt,-1)
}

func (rm *rcvmsg)GetRecv() RcvDataer  {
	return rm.rcv
}

func (rm *rcvmsg)SetRecv(dataer RcvDataer)  {
	rm.rcv = dataer
}
