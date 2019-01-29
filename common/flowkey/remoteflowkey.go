package flowkey

import (
	"sync"
)


type flowKeyMap struct {
	fm map[uint64]*FlowKey
}


type FlowKeyMap interface {
	Add(sn uint64,fk *FlowKey)
	Get(sn uint64) *FlowKey
	Del(sn uint64)
}


var (
	globalFlowKeyMap FlowKeyMap
	lock sync.Mutex
)

func newFlowKeyMap() FlowKeyMap {
	fkm := &flowKeyMap{}
	fkm.fm = make(map[uint64]*FlowKey)

	return fkm
}

func GetFlowKeyMapInstance() FlowKeyMap {
	if globalFlowKeyMap != nil{
		return globalFlowKeyMap
	}

	lock.Lock()
	if globalFlowKeyMap != nil{
		lock.Unlock()
		return globalFlowKeyMap
	}

	globalFlowKeyMap = newFlowKeyMap()

	lock.Unlock()

	return globalFlowKeyMap
}

func (rfk *flowKeyMap)Add(sn uint64,fk *FlowKey) {
	rfk.fm[sn] = fk
}

func (rfk *flowKeyMap)Get(sn uint64) *FlowKey  {
	if k,ok:=rfk.fm[sn];ok{
		return k
	}

	return nil
}

func (rfk *flowKeyMap)Del(sn uint64)  {
	delete(rfk.fm,sn)
}


