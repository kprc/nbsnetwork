package nbsnetwork

import "sync"

func init()  {
	G_BlockDatasLock  = sync.RWMutex{}
	G_BlockDatas = make(map[uint64]BlockDataer,0)
}
