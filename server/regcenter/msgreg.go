package regcenter

import "sync"

type msgCenter struct {
	rwlock sync.RWMutex
	coor map[int32]MsgHandler
}


type MsgCenter interface {
	AddHandler(msgid int32,handler MsgHandler)
	DelHandler(msgid int32)
	GetHandler(msgid int32) MsgHandler

}

func (mc *msgCenter)AddHandler(msgid int32,handler MsgHandler)  {
	mc.rwlock.Lock()
	defer mc.rwlock.Unlock()
	handler.IncRef()

	mc.coor[msgid] = handler
}

func (mc *msgCenter)DelHandler(msgid int32)  {
	mc.rwlock.Lock()
	defer mc.rwlock.Unlock()

	if v,ok:=mc.coor[msgid]; ok {
		v.DecRef()
		if v.GetRef() == 0 {
			delete(mc.coor,msgid)
		}
	}
}

func (mc *msgCenter)GetHandler(msgid int32)  MsgHandler{
	mc.rwlock.RLock()
	defer mc.rwlock.Unlock()

	if v,ok:=mc.coor[msgid]; ok {
		v.IncRef()
		return v
	}

	return nil

}
func (mc *msgCenter)PutHandler(msgid int32) {
	mc.rwlock.RLock()
	defer mc.rwlock.Unlock()

	if v,ok:=mc.coor[msgid]; ok {
		v.DecRef()
	}
}


