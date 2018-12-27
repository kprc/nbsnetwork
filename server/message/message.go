package message

import (
	"io"
	"sync/atomic"
	"time"
)

type MsgKey struct {
	serialNo uint64
	stationId string
}

type rcvmsg struct {
	key *MsgKey
	w io.WriteSeeker
	refcnt int32
	timeout int64
}


type RcvMsg interface {
	SetWS(w io.WriteSeeker)
	GetWS() io.WriteSeeker
	SetRefCnt(cnt int32)
	GetRefCnt() int32
	IncRefCnt()
	DecRefCnt() int32
	SetKey(key *MsgKey)
	GetKey() *MsgKey
}

func NewRcvMsg() RcvMsg  {
	rm := &rcvmsg{}
	rm.timeout = time.Now().Unix()

	return rm
}

func NewMsgKey(sn uint64,stationId string) *MsgKey {
	return &MsgKey{serialNo:sn,stationId:stationId}
}

func (rm *rcvmsg)SetWS(w io.WriteSeeker) {
	rm.w = w
}
func (rm *rcvmsg)GetWS() io.WriteSeeker {
	return rm.w
}
func (mk *MsgKey)SetSerialNo(sn uint64)  {
	mk.serialNo = sn
}
func (mk *MsgKey)GetSerialNo() uint64 {
	return mk.serialNo
}
func (mk *MsgKey)SetStationId(id string)  {
	mk.stationId = id
}
func (rm *MsgKey)GetStationId() string {
	return rm.stationId
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

func (rm *rcvmsg)SetKey(key *MsgKey)  {
	rm.key = key
}

func (rm *rcvmsg)GetKey() *MsgKey {
	return rm.key
}