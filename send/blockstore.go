package send

import (
	"github.com/kprc/nbsnetwork/common/flowkey"
	"sync"
)

type bstore struct {
	glock *sync.RWMutex
	sd map[flowkey.FlowKey]StoreDataer
}

type BStorer interface {
	AddBlockDataer(key *flowkey.FlowKey,sdr StoreDataer)
	GetBlockDataer(key *flowkey.FlowKey) StoreDataer
	PutBlockDataer(key *flowkey.FlowKey)
	DelBlockDataer(key *flowkey.FlowKey)
}

var (
	initlock sync.Mutex
	instance BStorer

)

func GetInstance() BStorer {
	if instance == nil{
		initlock.Lock()

		if instance == nil{
			instance = &bstore{sd:make(map[flowkey.FlowKey]StoreDataer),
				glock:&sync.RWMutex{}}
		}
		initlock.Unlock()
	}
	return instance
}

func (bs *bstore)AddBlockDataer(key *flowkey.FlowKey,bdr StoreDataer)  {
	bs.glock.Lock()
	defer bs.glock.Unlock()

	if _,ok := bs.sd[*key];ok {
		return
	}

	//sd := &storeData{bdr:bdr}
	bdr.ReferCntInc()
	bs.sd[*key] = bdr

}

func (bs *bstore)GetBlockDataer(key *flowkey.FlowKey) StoreDataer{
	bs.glock.RLock()
	defer bs.glock.RUnlock()
	if v,ok:=bs.sd[*key];ok {
		v.ReferCntInc()
		return v
	}

	return nil
}

func (bs *bstore)decRefer(key *flowkey.FlowKey){
	bs.glock.Lock()
	defer bs.glock.Unlock()
	if v,ok :=bs.sd[*key];ok {
		v.ReferCntDec()
		if v.GetReferCnt()<=0 {
			delete(bs.sd,*key)
		}
	}
}


func (bs *bstore)PutBlockDataer(key *flowkey.FlowKey) {
	bs.decRefer(key)
}

func (bs *bstore)DelBlockDataer(key *flowkey.FlowKey)  {
	bs.decRefer(key)
}