package nbsnetwork

import (
	"sync"
)



type bstore struct {
	glock *sync.RWMutex
	sd map[uint64]StoreDataer
}

type BStorer interface {
	AddBlockDataer(sn uint64,sdr StoreDataer)
	GetBlockDataer(sn uint64) StoreDataer
	PutBlockDataer(sn uint64)
}

var (
	initlock sync.Mutex
	instance BStorer

)

func GetInstance() BStorer {

	if instance != nil{
		return instance
	}else {
		initlock.Lock()
		defer initlock.Unlock()
		if instance != nil{
			return instance
		}
		instance = &bstore{sd:make(map[uint64]StoreDataer),
							glock:&sync.RWMutex{}}

	}

	return instance
}

func (bs *bstore)AddBlockDataer(sn uint64,bdr StoreDataer)  {
	bs.glock.Lock()
	defer bs.glock.Unlock()

	if _,ok := bs.sd[sn];ok {
		return
	}

	sd := storeData{lock:&sync.RWMutex{},bdr:bdr}

	bs.sd[sn] = sd

}

func (bs *bstore)GetBlockDataer(sn uint64) StoreDataer{
	bs.glock.RLock()
	defer bs.glock.RUnlock()
	if v,ok:=bs.sd[sn];ok {
		v.ReferCntInc()
		return v
	}

	return nil
}
func (bs *bstore)PutBlockDataer(sn uint64){
	bs.glock.Lock()
	defer bs.glock.Unlock()
	if v,ok :=bs.sd[sn];ok {
		v.ReferCntDec()
		if v.GetReferCnt()<=0 {
			delete(bs.sd,sn)
		}
	}
}