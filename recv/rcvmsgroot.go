package recv

import (
	"github.com/kprc/nbsnetwork/common/flowkey"
	"reflect"
	"sync"
	"time"
)

type rcvmsgroot struct {
	rwlock sync.RWMutex
	msgroot map[flowkey.FlowKey]RcvMsg
	cmd chan int //1 for quit
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
	TimeOut()
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
	rmr.cmd = make(chan int,0)

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

func (rmr *rcvmsgroot)TimeOut()  {

	for {

		fkarr := make([]flowkey.FlowKey, 0)

		rmr.rwlock.RLock()

		keys := reflect.ValueOf(rmr.msgroot).MapKeys()

		for _, k := range keys {
			fk := k.Interface().(flowkey.FlowKey)

			rm := rmr.msgroot[fk]

			rm.IncRefCnt()

			rcv := rm.GetRecv()

			if rcv.TimeOut() {
				fkarr = append(fkarr, fk)
			}
			rm.DecRefCnt()

		}
		rmr.rwlock.RUnlock()

		for _, fk := range fkarr {
			rmr.DelMsg(&fk)
		}
		select {
		case cmd:=<-rmr.cmd:
			if cmd==1 {
				return
			}
		default:
			time.Sleep(time.Second*1)
			//nothing to do ...
		}

	}

}
