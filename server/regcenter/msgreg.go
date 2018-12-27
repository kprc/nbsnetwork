package regcenter

import (
	"sync"
	"github.com/kprc/nbsnetwork/server/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/kprc/nbsnetwork/common/constant"
)

type msgCenter struct {
	fGetMsgId func([]byte) (msgid int32,stationId string)
	rwlock sync.RWMutex
	coor map[int32]MsgHandler
}



type MsgCenter interface {
	AddHandler(msgid int32,handler MsgHandler)
	DelHandler(msgid int32)
	GetHandler(msgid int32) MsgHandler
	GetMsgId([] byte) (msgid int32,stationId string)
}


var (
	glock sync.Mutex
	instance MsgCenter
)


func GetMsgCenter() MsgCenter{
	if instance == nil{
		glock.Lock()
		if instance == nil {
			mc := &msgCenter{fGetMsgId:getMsgId}
			mc.coor = make(map[int32]MsgHandler)
			instance = mc
		}
		glock.Unlock()
	}
	return instance
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

func (mc *msgCenter) GetMsgId(data [] byte) (msgid int32,stationId string) {
	return mc.fGetMsgId(data)
}

func getMsgId(headData []byte) (msgid int32,stationId string)  {
	mh := message.MsgHead{}

	err := proto.Unmarshal(headData,&mh)

	if err == nil{
		return mh.MessageId,string(mh.StationId)
	}else {
		return constant.MSG_NONE,""
	}
}

