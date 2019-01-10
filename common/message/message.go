package message

import (
	"github.com/kprc/nbsnetwork/common/flowkey"
	"github.com/kprc/nbsnetwork/netcommon"
	"github.com/kprc/nbsnetwork/recv"
	"io"
	"sync/atomic"
)


type rcvmsg struct {
	key *flowkey.FlowKey
	w io.WriteSeeker     //used in rcv
	refcnt int32
	uw netcommon.UdpReaderWriterer   //for reply
	rcv recv.RcvDataer    //for receive
}


type RcvMsg interface {
	SetWS(w io.WriteSeeker)
	GetWS() io.WriteSeeker
	SetRefCnt(cnt int32)
	GetRefCnt() int32
	IncRefCnt()
	DecRefCnt() int32
	SetKey(key *flowkey.FlowKey)
	GetKey() *flowkey.FlowKey
	SetUW(uw netcommon.UdpReaderWriterer)
	GetUW() netcommon.UdpReaderWriterer
	GetRecv() recv.RcvDataer
	SetRecv(rcv recv.RcvDataer)
}

func NewRcvMsg() RcvMsg  {
	rm := &rcvmsg{}

	return rm
}


func (rm *rcvmsg)SetWS(w io.WriteSeeker) {
	rm.w = w
}
func (rm *rcvmsg)GetWS() io.WriteSeeker {
	return rm.w
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

func (rm *rcvmsg)SetKey(key *flowkey.FlowKey)  {
	rm.key = key
}

func (rm *rcvmsg)GetKey() *flowkey.FlowKey {
	return rm.key
}

func (rm *rcvmsg)SetUW(uw netcommon.UdpReaderWriterer) {
	rm.uw = uw
}

func (rm *rcvmsg)GetUW() netcommon.UdpReaderWriterer  {
	return rm.uw
}

func (rm *rcvmsg)GetRecv() recv.RcvDataer  {
	return rm.rcv
}

func (rm *rcvmsg)SetRecv(dataer recv.RcvDataer)  {
	rm.rcv = dataer
}
