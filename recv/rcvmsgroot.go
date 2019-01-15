package recv

import (
	"github.com/kprc/nbsnetwork/common/flowkey"
	"sync"
)

type rcvmsgroot struct {
	rwlock sync.RWMutex
	msgroot map[flowkey.FlowKey]RcvMsg
}

var (
	instance RcvMsgRoot
	glock sync.Mutex
)


type RcvMsgRoot interface {
	AddMSG(mk *flowkey.FlowKey,rm RcvMsg)
	DelMsg(mk *flowkey.FlowKey)
	GetMsg(mk *flowkey.FlowKey) RcvMsg
	PutMsg(mk *flowkey.FlowKey)
}

func GetInstance() RcvMsgRoot {
	if instance == nil{
		 glock.Lock()
		if instance == nil{
			instance = newRcvMsgRoot()
		}
		 glock.Unlock()
	}
	return instance
}

func newRcvMsgRoot() RcvMsgRoot {
	rmr:=&rcvmsgroot{}
	rmr.msgroot = make(map[flowkey.FlowKey]RcvMsg)

	return rmr
}

func (rmr *rcvmsgroot)AddMSG(mk *flowkey.FlowKey,rm RcvMsg) {
	rmr.rwlock.Lock()
	defer rmr.rwlock.Unlock()
	if _,ok := rmr.msgroot[*mk];ok{
		return
	}
	rm.IncRefCnt()
	rmr.msgroot[*mk] = rm
}


func (rmr *rcvmsgroot)DelMsg(mk *flowkey.FlowKey)  {
	rmr.rwlock.Lock()
	defer rmr.rwlock.Unlock()
	if v,ok := rmr.msgroot[*mk];ok{
		cnt := v.DecRefCnt()
		if cnt <= 0 {
			delete(rmr.msgroot,*mk)
		}
	}
}


func (rmr *rcvmsgroot)GetMsg(mk *flowkey.FlowKey) RcvMsg  {
	rmr.rwlock.RLock()
	defer rmr.rwlock.RUnlock()
	if v,ok := rmr.msgroot[*mk];ok{
		v.IncRefCnt()
		return v
	}
	return nil
}

func (rmr *rcvmsgroot)PutMsg(mk *flowkey.FlowKey)  {
	rmr.rwlock.RLock()

	if v,ok := rmr.msgroot[*mk];ok{
		cnt:=v.DecRefCnt()
		if cnt <= 0{
			rmr.rwlock.RUnlock()
			rmr.rwlock.Lock()
			cnt = v.GetRefCnt()
			if cnt <=0{
				delete(rmr.msgroot,*mk)
			}
			rmr.rwlock.Unlock()
		}else{
			rmr.rwlock.RUnlock()
		}
	}else {
		rmr.rwlock.RUnlock()
	}
}
