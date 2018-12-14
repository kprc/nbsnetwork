package nbsnetwork

import (
	"sync"
	"sync/atomic"
)


type storeData struct {
	lock *sync.RWMutex
	bdr BlockDataer
	cnt uint32
}


type bstore struct {
	glock *sync.RWMutex
	sd map[uint64]storeData
}

type BStorer interface {
	AddBlockDataer(sn uint64,bdr BlockDataer)
	GetBlockDataer(sn uint64) BlockDataer
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
		if instance != nil{
			return instance
		}
		instance = &bstore{sd:make(map[uint64]storeData),
							glock:&sync.RWMutex{}}
		initlock.Unlock()
	}

	return instance
}

func (bs *bstore)AddBlockDataer(sn uint64,bdr BlockDataer)  {
	bs.glock.Lock()
	defer bs.glock.Unlock()

	if _,ok := bs.sd[sn];ok {
		return
	}

	sd := storeData{lock:&sync.RWMutex{},bdr:bdr}

	bs.sd[sn] = sd

}

func (bs *bstore)GetBlockDataer(sn uint64) BlockDataer{
	bs.glock.RLock()
	defer bs.glock.RUnlock()
	if v,ok:=bs.sd[sn];ok {
		atomic.AddUint32(&v.cnt,1)
		return v
	}

	return nil
}
func (bs *bstore)PutBlockDataer(sn uint64){
	bs.glock.Lock()
	defer bs.glock.Unlock()
	if v,ok :=bs.sd[sn];ok {
		atomic.AddUint32(&v.cnt,-1)
		if v.cnt<=0 {
			delete(bs.sd,sn)
		}
	}
}