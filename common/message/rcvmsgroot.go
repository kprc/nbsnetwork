package message

import "sync"

type rcvmsgroot struct {
	rwlock sync.RWMutex
	msgroot map[MsgKey]RcvMsg
}

var (
	instance RcvMsgRoot
	glock sync.Mutex
)


type RcvMsgRoot interface {
	AddMSG(mk *MsgKey,rm RcvMsg)
	DelMsg(mk *MsgKey)
	GetMsg(mk *MsgKey) RcvMsg
	PutMsg(mk *MsgKey)
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
	rmr.msgroot = make(map[MsgKey]RcvMsg)

	return rmr
}

func (rmr *rcvmsgroot)AddMSG(mk *MsgKey,rm RcvMsg) {
	rmr.rwlock.Lock()
	defer rmr.rwlock.Unlock()
	if _,ok := rmr.msgroot[*mk];ok{
		return
	}
	rm.IncRefCnt()
	rmr.msgroot[*mk] = rm
}


func (rmr *rcvmsgroot)DelMsg(mk *MsgKey)  {
	rmr.rwlock.Lock()
	defer rmr.rwlock.Unlock()
	if v,ok := rmr.msgroot[*mk];ok{
		cnt := v.DecRefCnt()
		if cnt <= 0 {
			delete(rmr.msgroot,*mk)
		}
	}
}


func (rmr *rcvmsgroot)GetMsg(mk *MsgKey) RcvMsg  {
	rmr.rwlock.RLock()
	defer rmr.rwlock.RUnlock()
	if v,ok := rmr.msgroot[*mk];ok{
		v.IncRefCnt()
		return v
	}
	return nil
}

func (rmr *rcvmsgroot)PutMsg(mk *MsgKey)  {
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
